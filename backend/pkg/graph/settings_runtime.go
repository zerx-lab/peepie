package graph

import (
	"context"
	"fmt"
	"path/filepath"
	"sort"
	"strconv"
	"strings"

	"pentagi/pkg/config"
	"pentagi/pkg/database"
	"pentagi/pkg/graph/model"
	"pentagi/pkg/providers"
	"pentagi/pkg/providers/provider"
)

// providerConfigsDir holds the provider config files shipped in the image
// (and the ones docker-compose mounts), offered as configPath suggestions.
const providerConfigsDir = "/opt/pentagi/conf"

// applySettings publishes one category's Web UI changes through
// config.ApplySettings: the settings .env file first, then the system_settings
// table, then the live overrides, so every subsequent read observes them.
func applySettings(
	ctx context.Context,
	db database.Querier,
	cfg *config.Config,
	uid int64,
	category string,
	values map[string]string,
) error {
	return cfg.ApplySettings(category, values, func() error {
		return persistSettings(ctx, db, uid, category, values)
	})
}

// persistSettings upserts values into system_settings. Callers go through
// config.ApplySettings (directly or via the provider controller), which
// normalizes values before this runs.
func persistSettings(
	ctx context.Context,
	db database.Querier,
	uid int64,
	category string,
	values map[string]string,
) error {
	for key, value := range values {
		if _, err := db.UpsertSystemSetting(ctx, database.UpsertSystemSettingParams{
			Category:  category,
			Key:       key,
			Value:     value,
			IsSecret:  config.SecretKeys[key],
			UpdatedBy: database.Int64ToNullInt64(&uid),
		}); err != nil {
			return err
		}
	}
	return nil
}

// putSecret records a secret input unless it is nil: the frontend never
// re-sends a plaintext secret it only echoed as "configured", so nil means
// "leave unchanged", while an explicit empty string clears the secret.
func putSecret(values map[string]string, key string, value *string) {
	if value != nil {
		values[key] = *value
	}
}

// llmProviderSettings reports the live built-in LLM provider settings (env
// values with llm_providers overrides applied) and each provider's status.
func llmProviderSettings(cfg *config.Config, pc providers.ProviderController) *model.LLMProviderSettings {
	c := cfg.WithOverrides(config.CategoryLLMProviders)

	status := func(prvtype provider.ProviderType) (bool, *string) {
		active, err := pc.DefaultProviderStatus(prvtype)
		if err != nil {
			msg := err.Error()
			return active, &msg
		}
		return active, nil
	}
	keyed := func(
		prvtype provider.ProviderType, apiKey, serverURL string, providerName *string,
	) *model.LLMProviderKeySettings {
		active, errMsg := status(prvtype)
		return &model.LLMProviderKeySettings{
			Active:       active,
			Error:        errMsg,
			APIKeySet:    apiKey != "",
			ServerURL:    serverURL,
			ProviderName: providerName,
		}
	}

	bedrockActive, bedrockErr := status(provider.ProviderBedrock)
	ollamaActive, ollamaErr := status(provider.ProviderOllama)
	customActive, customErr := status(provider.ProviderCustom)

	return &model.LLMProviderSettings{
		Openai:    keyed(provider.ProviderOpenAI, c.OpenAIKey, c.OpenAIServerURL, nil),
		Anthropic: keyed(provider.ProviderAnthropic, c.AnthropicAPIKey, c.AnthropicServerURL, nil),
		Gemini:    keyed(provider.ProviderGemini, c.GeminiAPIKey, c.GeminiServerURL, nil),
		Bedrock: &model.BedrockProviderSettings{
			Active:             bedrockActive,
			Error:              bedrockErr,
			Region:             c.BedrockRegion,
			DefaultAuth:        c.BedrockDefaultAuth,
			BearerTokenSet:     c.BedrockBearerToken != "",
			AccessKeyIDSet:     c.BedrockAccessKey != "",
			SecretAccessKeySet: c.BedrockSecretKey != "",
			SessionTokenSet:    c.BedrockSessionToken != "",
			ServerURL:          c.BedrockServerURL,
		},
		Ollama: &model.OllamaProviderSettings{
			Active:            ollamaActive,
			Error:             ollamaErr,
			ServerURL:         c.OllamaServerURL,
			APIKeySet:         c.OllamaServerAPIKey != "",
			Model:             c.OllamaServerModel,
			ConfigPath:        c.OllamaServerConfig,
			PullModelsEnabled: c.OllamaServerPullModelsEnabled,
			PullModelsTimeout: c.OllamaServerPullModelsTimeout,
			LoadModelsEnabled: c.OllamaServerLoadModelsEnabled,
		},
		Custom: &model.CustomProviderSettings{
			Active:            customActive,
			Error:             customErr,
			ServerURL:         c.LLMServerURL,
			APIKeySet:         c.LLMServerKey != "",
			Model:             c.LLMServerModel,
			ConfigPath:        c.LLMServerConfig,
			ProviderName:      c.LLMServerProvider,
			LegacyReasoning:   c.LLMServerLegacyReasoning,
			PreserveReasoning: c.LLMServerPreserveReasoning,
		},
		Deepseek:    keyed(provider.ProviderDeepSeek, c.DeepSeekAPIKey, c.DeepSeekServerURL, &c.DeepSeekProvider),
		Glm:         keyed(provider.ProviderGLM, c.GLMAPIKey, c.GLMServerURL, &c.GLMProvider),
		Kimi:        keyed(provider.ProviderKimi, c.KimiAPIKey, c.KimiServerURL, &c.KimiProvider),
		Qwen:        keyed(provider.ProviderQwen, c.QwenAPIKey, c.QwenServerURL, &c.QwenProvider),
		Minimax:     keyed(provider.ProviderMiniMax, c.MiniMaxAPIKey, c.MiniMaxServerURL, &c.MiniMaxProvider),
		ConfigPaths: providerConfigPaths(),
	}
}

// providerConfigPaths lists the provider config files available in the
// container; empty outside the Docker image.
func providerConfigPaths() []string {
	paths := make([]string, 0)
	for _, pattern := range []string{"*.yml", "*.yaml", "*.json"} {
		matches, err := filepath.Glob(filepath.Join(providerConfigsDir, pattern))
		if err == nil {
			paths = append(paths, matches...)
		}
	}
	sort.Strings(paths)
	return paths
}

// llmProviderSettingValues flattens the providers present in the input into
// llm_providers setting values; absent providers and nil secrets are left
// unchanged.
func llmProviderSettingValues(in model.LLMProviderSettingsInput) (map[string]string, error) {
	values := make(map[string]string)

	keyed := func(p *model.LLMProviderKeySettingsInput, apiKey, serverURL, providerName string) {
		if p == nil {
			return
		}
		putSecret(values, apiKey, p.APIKey)
		values[serverURL] = p.ServerURL
		if providerName != "" && p.ProviderName != nil {
			values[providerName] = *p.ProviderName
		}
	}

	keyed(in.Openai, config.KeyOpenAIKey, config.KeyOpenAIServerURL, "")
	keyed(in.Anthropic, config.KeyAnthropicAPIKey, config.KeyAnthropicServerURL, "")
	keyed(in.Gemini, config.KeyGeminiAPIKey, config.KeyGeminiServerURL, "")
	keyed(in.Deepseek, config.KeyDeepSeekAPIKey, config.KeyDeepSeekServerURL, config.KeyDeepSeekProvider)
	keyed(in.Glm, config.KeyGLMAPIKey, config.KeyGLMServerURL, config.KeyGLMProvider)
	keyed(in.Kimi, config.KeyKimiAPIKey, config.KeyKimiServerURL, config.KeyKimiProvider)
	keyed(in.Qwen, config.KeyQwenAPIKey, config.KeyQwenServerURL, config.KeyQwenProvider)
	keyed(in.Minimax, config.KeyMiniMaxAPIKey, config.KeyMiniMaxServerURL, config.KeyMiniMaxProvider)

	if b := in.Bedrock; b != nil {
		values[config.KeyBedrockRegion] = b.Region
		values[config.KeyBedrockDefaultAuth] = strconv.FormatBool(b.DefaultAuth)
		putSecret(values, config.KeyBedrockBearerToken, b.BearerToken)
		putSecret(values, config.KeyBedrockAccessKey, b.AccessKeyID)
		putSecret(values, config.KeyBedrockSecretKey, b.SecretAccessKey)
		putSecret(values, config.KeyBedrockSessionToken, b.SessionToken)
		values[config.KeyBedrockServerURL] = b.ServerURL
	}

	if o := in.Ollama; o != nil {
		if err := validateProviderConfigPath(o.ConfigPath); err != nil {
			return nil, err
		}
		if o.PullModelsTimeout < 0 {
			return nil, fmt.Errorf("ollama pull models timeout must not be negative")
		}
		values[config.KeyOllamaServerURL] = o.ServerURL
		putSecret(values, config.KeyOllamaServerAPIKey, o.APIKey)
		values[config.KeyOllamaServerModel] = o.Model
		values[config.KeyOllamaServerConfig] = o.ConfigPath
		values[config.KeyOllamaPullEnabled] = strconv.FormatBool(o.PullModelsEnabled)
		values[config.KeyOllamaPullTimeout] = strconv.Itoa(o.PullModelsTimeout)
		values[config.KeyOllamaLoadEnabled] = strconv.FormatBool(o.LoadModelsEnabled)
	}

	if c := in.Custom; c != nil {
		if err := validateProviderConfigPath(c.ConfigPath); err != nil {
			return nil, err
		}
		values[config.KeyLLMServerURL] = c.ServerURL
		putSecret(values, config.KeyLLMServerKey, c.APIKey)
		values[config.KeyLLMServerModel] = c.Model
		values[config.KeyLLMServerConfig] = c.ConfigPath
		values[config.KeyLLMServerProvider] = c.ProviderName
		values[config.KeyLLMServerLegacyReason] = strconv.FormatBool(c.LegacyReasoning)
		values[config.KeyLLMServerPreserveReas] = strconv.FormatBool(c.PreserveReasoning)
	}

	for key, value := range values {
		values[key] = strings.TrimSpace(value)
	}

	return values, nil
}

// validateProviderConfigPath limits provider config paths to absolute
// YAML/JSON files: the backend reads whatever path is configured, so an
// arbitrary file must not be pulled into provider-config parse errors.
func validateProviderConfigPath(path string) error {
	path = strings.TrimSpace(path)
	if path == "" {
		return nil
	}
	if !filepath.IsAbs(path) {
		return fmt.Errorf("provider config path must be absolute: %s", path)
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".yml", ".yaml", ".json":
		return nil
	default:
		return fmt.Errorf("provider config path must point to a .yml, .yaml or .json file: %s", path)
	}
}
