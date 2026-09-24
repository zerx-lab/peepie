package main

import (
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"pentagi/cmd/installer/checker"
	"pentagi/cmd/installer/files"
	"pentagi/cmd/installer/hardening"
	"pentagi/cmd/installer/state"
	"pentagi/cmd/installer/wizard"
	"pentagi/cmd/installer/wizard/locale"
	"pentagi/pkg/version"
)

type Config struct {
	envPath     string
	showVersion bool
	language    string
}

func main() {
	languagePreferencePath := initLanguage()

	config := parseFlags(os.Args)

	if config.showVersion {
		fmt.Println(version.GetBinaryVersion())
		os.Exit(0)
	}

	// parseFlags already switched to a supported -l language; reject anything else
	if config.language != "" {
		if _, ok := locale.Normalize(config.language); !ok {
			log.Fatalf(locale.CLIError, fmt.Errorf(locale.CLIUnsupportedLanguage, config.language, locale.SupportedLanguageTags()))
		}
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	setupSignalHandler(cancel)

	envPath, err := validateEnvPath(config.envPath)
	if err != nil {
		log.Fatalf(locale.CLIError, err)
	}

	appState, err := initializeState(envPath)
	if err != nil {
		log.Fatalf(locale.CLIFailedInitState, err)
	}

	if err := hardening.DoMigrateSettings(appState); err != nil {
		log.Fatalf(locale.CLIFailedMigrate, err)
	}

	if err := hardening.DoSyncNetworkSettings(appState); err != nil {
		log.Fatalf(locale.CLIFailedSyncNetwork, err)
	}

	checkResult, err := gatherSystemFacts(ctx, appState)
	if err != nil {
		log.Fatalf(locale.CLIFailedGatherFacts, err)
	}

	printStartupInfo(envPath, checkResult)

	if err := hardening.DoHardening(appState, checkResult); err != nil {
		log.Fatalf(locale.CLIFailedHardening, err)
	}

	if err := runApplication(ctx, appState, checkResult, languagePreferencePath); err != nil {
		log.Fatalf(locale.CLIApplicationError, err)
	}

	cleanup(appState)
}

// initLanguage applies the language chosen earlier with Ctrl+L, else the system
// language, before anything is printed. It returns where Ctrl+L saves the choice
// ("" when the user config dir is unavailable, which disables saving).
func initLanguage() string {
	lang := locale.DetectSystemLanguage(os.Getenv)

	path, err := locale.PreferencePath()
	if err != nil {
		path = ""
	} else if saved, ok := locale.LoadPreference(path); ok {
		lang = saved
	}

	// detected and loaded languages are always supported
	_ = locale.SetLanguage(lang)

	return path
}

func parseFlags(args []string) Config {
	name := "installer"
	if len(args) > 0 {
		args, name = args[1:], filepath.Base(args[0])
	}

	// A silent first pass only resolves -l, so the flag descriptions and usage text
	// of the real pass below are already in the requested language.
	var probe Config
	probeSet := newFlagSet(name, &probe)
	probeSet.SetOutput(io.Discard)
	probeSet.Usage = func() {}
	_ = probeSet.Parse(args)
	if lang, ok := locale.Normalize(probe.language); ok {
		_ = locale.SetLanguage(lang) // Normalize only returns supported languages
	}

	var config Config
	flagSet := newFlagSet(name, &config)
	flagSet.Usage = func() {
		fmt.Fprintf(os.Stderr, locale.CLIUsageTitle+"\n\n", version.GetBinaryVersion())
		fmt.Fprintf(os.Stderr, locale.CLIUsageLine+"\n\n", name)
		fmt.Fprintln(os.Stderr, locale.CLIUsageOptions)
		flagSet.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n%s\n", locale.CLIUsageExamples)
		fmt.Fprintf(os.Stderr, "  %s                    # %s\n", name, locale.CLIExampleDefault)
		fmt.Fprintf(os.Stderr, "  %s -e config/.env     # %s\n", name, locale.CLIExampleEnvFile)
		fmt.Fprintf(os.Stderr, "  %s -l zh-CN           # %s\n", name, locale.CLIExampleLanguage)
		fmt.Fprintf(os.Stderr, "  %s -v                 # %s\n", name, locale.CLIExampleVersion)
	}

	flagSet.Parse(args)
	return config
}

func newFlagSet(name string, config *Config) *flag.FlagSet {
	flagSet := flag.NewFlagSet(name, flag.ContinueOnError)
	flagSet.BoolVar(&config.showVersion, "v", false, locale.CLIFlagVersion)
	flagSet.StringVar(&config.envPath, "e", ".env", locale.CLIFlagEnvFile)
	flagSet.StringVar(&config.language, "l", "", fmt.Sprintf(locale.CLIFlagLanguage, locale.SupportedLanguageTags()))
	return flagSet
}

func setupSignalHandler(cancel context.CancelFunc) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		sig := <-sigChan
		log.Printf("Received signal: %v, initiating graceful shutdown...", sig)
		cancel()
	}()
}

func validateEnvPath(envPath string) (string, error) {
	// convert to absolute path
	absPath, err := filepath.Abs(envPath)
	if err != nil {
		return "", fmt.Errorf("invalid path '%s': %w", envPath, err)
	}

	// check if file exists
	if info, err := os.Stat(absPath); os.IsNotExist(err) {
		// file doesn't exist, check if we can create it in the directory
		dir := filepath.Dir(absPath)
		if _, err := os.Stat(dir); os.IsNotExist(err) {
			if err := os.MkdirAll(dir, 0755); err != nil {
				return "", fmt.Errorf("cannot create directory '%s': %w", dir, err)
			}
		} else if err != nil {
			return "", fmt.Errorf("cannot access directory '%s': %w", dir, err)
		}

		// try to create initial env file
		if err := createInitialEnvFile(absPath); err != nil {
			return "", fmt.Errorf("cannot create env file '%s': %w", absPath, err)
		}
	} else if info.IsDir() {
		return "", fmt.Errorf("'%s' is a directory", absPath)
	} else if err != nil {
		return "", fmt.Errorf("cannot access file '%s': %w", absPath, err)
	}

	return absPath, nil
}

func createInitialEnvFile(path string) error {
	f := files.NewFiles()

	content, err := f.GetContent(".env")
	if err != nil {
		return fmt.Errorf("cannot read .env file: %w", err)
	}

	content = fmt.Appendf(nil, `# PentAGI Environment Configuration
# Generated by PentAGI Installer v%s
#
# This file contains environment variables for PentAGI configuration.
# You can modify these values through the installer interface.
#
%s`, version.GetBinaryVersion(), string(content))

	if err := os.WriteFile(path, content, 0600); err != nil {
		return fmt.Errorf("cannot write .env file: %w", err)
	}

	return nil
}

func initializeState(envPath string) (state.State, error) {
	appState, err := state.NewState(envPath)
	if err != nil {
		return nil, fmt.Errorf("failed to create state manager: %w", err)
	}

	return appState, nil
}

func gatherSystemFacts(ctx context.Context, appState state.State) (checker.CheckResult, error) {
	result, err := checker.Gather(ctx, appState)
	if err != nil {
		return result, fmt.Errorf("failed to gather system facts: %w", err)
	}

	return result, nil
}

func printStartupInfo(envPath string, checkResult checker.CheckResult) {
	fmt.Printf(locale.CLIStartupTitle+"\n", version.GetBinaryVersion())
	fmt.Printf(locale.CLIStartupEnvFile+"\n", envPath)

	if !checkResult.IsReadyToContinue() {
		fmt.Println(locale.CLISystemNotReady)
	} else {
		fmt.Println(locale.CLISystemReady)
	}
}

func runApplication(
	ctx context.Context, appState state.State, checkResult checker.CheckResult, languagePreferencePath string,
) error {
	return wizard.Run(ctx, appState, checkResult, files.NewFiles(), languagePreferencePath)
}

func cleanup(appState state.State) {
	if appState.IsDirty() {
		fmt.Println(locale.CLIPendingChanges)
		fmt.Println(locale.CLIPendingChangesHint)
	}
}
