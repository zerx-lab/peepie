package config

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func writeTempEnv(t *testing.T, content string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), ".env")
	require.NoError(t, os.WriteFile(path, []byte(content), 0600))
	return path
}

func TestWriteEnvFileSettings_RewritesInPlacePreservingLayout(t *testing.T) {
	path := writeTempEnv(t, "# LLM providers\n"+
		"OPEN_AI_KEY=old # inline comment\n"+
		"#ANTHROPIC_API_KEY=commented\n"+
		"DATABASE_URL=postgres://x\n"+
		"export GEMINI_API_KEY=g-old\n")
	before, err := os.Stat(path)
	require.NoError(t, err)

	err = WriteEnvFileSettings(path, map[string]string{
		"OPEN_AI_KEY":       "sk-new",
		"GEMINI_API_KEY":    "g-new",
		"ANTHROPIC_API_KEY": "a-new",
		"TAVILY_API_KEY":    "",
	})
	require.NoError(t, err)

	raw, err := os.ReadFile(path)
	require.NoError(t, err)
	assert.Equal(t, "# LLM providers\n"+
		"OPEN_AI_KEY=sk-new\n"+
		"#ANTHROPIC_API_KEY=commented\n"+
		"DATABASE_URL=postgres://x\n"+
		"GEMINI_API_KEY=g-new\n"+
		"ANTHROPIC_API_KEY=a-new\n"+
		"TAVILY_API_KEY=\n", string(raw))

	after, err := os.Stat(path)
	require.NoError(t, err)
	assert.True(t, os.SameFile(before, after), "the file must be rewritten in place so Docker bind mounts keep seeing it")
	assert.Equal(t, os.FileMode(0600), after.Mode().Perm(), "permissions of the secrets file are preserved")
}

func TestWriteEnvFileSettings_SpecialValuesRoundTrip(t *testing.T) {
	path := writeTempEnv(t, "")
	values := map[string]string{
		"OPEN_AI_KEY":       "sk-with space # not a comment",
		"ANTHROPIC_API_KEY": "has$dollar",
		"GEMINI_API_KEY":    `it's "quoted" $HOME`,
	}

	require.NoError(t, WriteEnvFileSettings(path, values))
	got, err := ReadEnvFileSettings(path)
	require.NoError(t, err)

	assert.Equal(t, values["OPEN_AI_KEY"], got[CategoryLLMProviders][KeyOpenAIKey])
	assert.Equal(t, values["ANTHROPIC_API_KEY"], got[CategoryLLMProviders][KeyAnthropicAPIKey])
	assert.Equal(t, values["GEMINI_API_KEY"], got[CategoryLLMProviders][KeyGeminiAPIKey])
}

func TestWriteEnvFileSettings_RejectsLineBreaks(t *testing.T) {
	path := writeTempEnv(t, "OPEN_AI_KEY=old\n")

	err := WriteEnvFileSettings(path, map[string]string{"OPEN_AI_KEY": "a\nDATABASE_URL=evil"})

	require.Error(t, err)
	raw, _ := os.ReadFile(path)
	assert.Equal(t, "OPEN_AI_KEY=old\n", string(raw))
}

func TestReadEnvFileSettings_EmptyAssignmentResolvesToDefault(t *testing.T) {
	path := writeTempEnv(t, "PERPLEXITY_MODEL=\nOPEN_AI_KEY=\nDATABASE_URL=postgres://x\n")

	got, err := ReadEnvFileSettings(path)
	require.NoError(t, err)

	assert.Equal(t, "sonar-pro", got[CategorySearchEngines][KeyPerplexityModel],
		"an empty variable means the default at boot, so it must at runtime too")
	assert.Equal(t, "", got[CategoryLLMProviders][KeyOpenAIKey])
	assert.Len(t, got[CategoryLLMProviders], 1, "variables outside the settings catalogue are ignored")
}

func TestReloadOverrides_FileWinsOverDatabase(t *testing.T) {
	cfg := &Config{
		SettingsEnvFile: writeTempEnv(t, "TAVILY_API_KEY=from-file\n"),
		Overrides:       NewOverrides(),
	}

	err := cfg.ReloadOverrides(func() (map[string]map[string]string, error) {
		return map[string]map[string]string{
			CategorySearchEngines: {KeyTavilyAPIKey: "from-db", KeyGoogleCXKey: "cx-db"},
		}, nil
	})
	require.NoError(t, err)

	assert.Equal(t, "from-file", cfg.Overrides.GetString(CategorySearchEngines, KeyTavilyAPIKey, ""))
	assert.Equal(t, "cx-db", cfg.Overrides.GetString(CategorySearchEngines, KeyGoogleCXKey, ""),
		"database values for variables absent from the file still apply")
}

func TestReloadOverrides_MissingFileUsesDatabase(t *testing.T) {
	cfg := &Config{
		SettingsEnvFile: filepath.Join(t.TempDir(), "absent.env"),
		Overrides:       NewOverrides(),
	}

	err := cfg.ReloadOverrides(func() (map[string]map[string]string, error) {
		return map[string]map[string]string{CategorySearchEngines: {KeyTavilyAPIKey: "from-db"}}, nil
	})

	require.NoError(t, err)
	assert.Equal(t, "from-db", cfg.Overrides.GetString(CategorySearchEngines, KeyTavilyAPIKey, ""))
}

func TestReloadOverrides_FailureKeepsCurrentOverrides(t *testing.T) {
	cfg := &Config{
		SettingsEnvFile: t.TempDir(), // a directory: unreadable as an env file
		Overrides:       NewOverrides(),
	}
	cfg.Overrides.Set(CategorySearchEngines, KeyTavilyAPIKey, "live")

	err := cfg.ReloadOverrides(func() (map[string]map[string]string, error) { return nil, nil })
	require.Error(t, err)
	assert.Equal(t, "live", cfg.Overrides.GetString(CategorySearchEngines, KeyTavilyAPIKey, ""))

	err = cfg.ReloadOverrides(func() (map[string]map[string]string, error) { return nil, errors.New("db down") })
	require.Error(t, err)
	assert.Equal(t, "live", cfg.Overrides.GetString(CategorySearchEngines, KeyTavilyAPIKey, ""))
}

func TestApplySettings_WritesFileThenPersistsThenPublishes(t *testing.T) {
	path := writeTempEnv(t, "PERPLEXITY_MODEL=custom\n")
	cfg := &Config{SettingsEnvFile: path, Overrides: NewOverrides()}

	var persisted map[string]string
	values := map[string]string{KeyPerplexityModel: "", KeyTavilyAPIKey: "tv"}
	err := cfg.ApplySettings(CategorySearchEngines, values, func() error {
		persisted = map[string]string{KeyPerplexityModel: values[KeyPerplexityModel], KeyTavilyAPIKey: values[KeyTavilyAPIKey]}
		return nil
	})
	require.NoError(t, err)

	assert.Equal(t, "sonar-pro", persisted[KeyPerplexityModel], "a cleared field is persisted as its default")
	assert.Equal(t, "sonar-pro", cfg.Overrides.GetString(CategorySearchEngines, KeyPerplexityModel, ""))
	raw, _ := os.ReadFile(path)
	assert.Equal(t, "PERPLEXITY_MODEL=sonar-pro\nTAVILY_API_KEY=tv\n", string(raw))
}

func TestApplySettings_FileErrorAbortsBeforePersist(t *testing.T) {
	cfg := &Config{SettingsEnvFile: t.TempDir(), Overrides: NewOverrides()}

	persisted := false
	err := cfg.ApplySettings(CategorySearchEngines, map[string]string{KeyTavilyAPIKey: "tv"}, func() error {
		persisted = true
		return nil
	})

	require.Error(t, err)
	assert.False(t, persisted, "the file would revert the change on the next reload, so nothing may be persisted")
	_, ok := cfg.Overrides.Get(CategorySearchEngines, KeyTavilyAPIKey)
	assert.False(t, ok)
}

func TestApplySettings_MissingFileStillPersists(t *testing.T) {
	cfg := &Config{SettingsEnvFile: filepath.Join(t.TempDir(), "absent.env"), Overrides: NewOverrides()}

	persisted := false
	err := cfg.ApplySettings(CategorySearchEngines, map[string]string{KeyTavilyAPIKey: "tv"}, func() error {
		persisted = true
		return nil
	})

	require.NoError(t, err)
	assert.True(t, persisted)
	assert.Equal(t, "tv", cfg.Overrides.GetString(CategorySearchEngines, KeyTavilyAPIKey, ""))
}

func TestWithSettings_ParsesLikeEnvWithoutMutatingBase(t *testing.T) {
	base := &Config{OpenAIKey: "env-key", OpenAIServerURL: "https://env.example/v1", OllamaServerPullModelsTimeout: 600}

	got := base.WithSettings(CategoryLLMProviders, map[string]string{
		KeyOpenAIKey:           "ui-key",
		KeyOpenAIServerURL:     "",
		KeyOllamaPullTimeout:   "30",
		KeyOllamaPullEnabled:   "true",
		KeyOllamaLoadEnabled:   "not-a-bool",
		KeyTavilyAPIKey:        "wrong category is ignored",
		"not_a_catalogued_key": "x",
	})

	assert.Equal(t, "ui-key", got.OpenAIKey)
	assert.Equal(t, "https://api.openai.com/v1", got.OpenAIServerURL, "empty means the envDefault")
	assert.Equal(t, 30, got.OllamaServerPullModelsTimeout)
	assert.True(t, got.OllamaServerPullModelsEnabled)
	assert.False(t, got.OllamaServerLoadModelsEnabled, "unparsable values keep the env value")
	assert.Equal(t, "", got.TavilyAPIKey)
	assert.Equal(t, "env-key", base.OpenAIKey, "the shared base config is never mutated")
}
