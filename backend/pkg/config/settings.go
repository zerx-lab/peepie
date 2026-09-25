package config

import (
	"fmt"
	"reflect"
	"strconv"
)

// settingDef ties one hot-reloadable runtime setting (category, key) to the
// environment variable that seeds the matching Config field at boot. The same
// variable name is used to mirror Web UI changes into the settings .env file
// (the file the installer TUI edits), so both editors always agree.
type settingDef struct {
	category string
	key      string
	env      string
}

var settingDefs = []settingDef{
	{CategorySearchEngines, KeyDuckDuckGoEnabled, "DUCKDUCKGO_ENABLED"},
	{CategorySearchEngines, KeyDuckDuckGoRegion, "DUCKDUCKGO_REGION"},
	{CategorySearchEngines, KeyDuckDuckGoSafeSearch, "DUCKDUCKGO_SAFESEARCH"},
	{CategorySearchEngines, KeyDuckDuckGoTimeRange, "DUCKDUCKGO_TIME_RANGE"},
	{CategorySearchEngines, KeySploitusEnabled, "SPLOITUS_ENABLED"},
	{CategorySearchEngines, KeyGoogleAPIKey, "GOOGLE_API_KEY"},
	{CategorySearchEngines, KeyGoogleCXKey, "GOOGLE_CX_KEY"},
	{CategorySearchEngines, KeyGoogleLRKey, "GOOGLE_LR_KEY"},
	{CategorySearchEngines, KeyTraversaalAPIKey, "TRAVERSAAL_API_KEY"},
	{CategorySearchEngines, KeyTavilyAPIKey, "TAVILY_API_KEY"},
	{CategorySearchEngines, KeyFirecrawlAPIKey, "FIRECRAWL_API_KEY"},
	{CategorySearchEngines, KeyFirecrawlAPIURL, "FIRECRAWL_API_URL"},
	{CategorySearchEngines, KeyPerplexityAPIKey, "PERPLEXITY_API_KEY"},
	{CategorySearchEngines, KeyPerplexityModel, "PERPLEXITY_MODEL"},
	{CategorySearchEngines, KeyPerplexityContext, "PERPLEXITY_CONTEXT_SIZE"},
	{CategorySearchEngines, KeySearxngURL, "SEARXNG_URL"},
	{CategorySearchEngines, KeySearxngCategories, "SEARXNG_CATEGORIES"},
	{CategorySearchEngines, KeySearxngLanguage, "SEARXNG_LANGUAGE"},
	{CategorySearchEngines, KeySearxngSafeSearch, "SEARXNG_SAFESEARCH"},
	{CategorySearchEngines, KeySearxngTimeRange, "SEARXNG_TIME_RANGE"},
	{CategorySearchEngines, KeySearxngTimeout, "SEARXNG_TIMEOUT"},
	{CategorySearchEngines, KeyWebSearchIntEnabled, "WEB_SEARCH_INTERNAL_ENABLED"},
	{CategorySearchEngines, KeyWebSearchIntMaxSites, "WEB_SEARCH_INTERNAL_MAX_SITES"},
	{CategorySearchEngines, KeyWebSearchIntMaxBytes, "WEB_SEARCH_INTERNAL_MAX_SITE_BYTES"},

	{CategoryExecution, KeyExecutionMonitorEnabled, "EXECUTION_MONITOR_ENABLED"},
	{CategoryExecution, KeyExecutionSameToolLimit, "EXECUTION_MONITOR_SAME_TOOL_LIMIT"},
	{CategoryExecution, KeyExecutionTotalToolLimit, "EXECUTION_MONITOR_TOTAL_TOOL_LIMIT"},
	{CategoryExecution, KeyMaxGeneralAgentToolCalls, "MAX_GENERAL_AGENT_TOOL_CALLS"},
	{CategoryExecution, KeyMaxLimitedAgentToolCalls, "MAX_LIMITED_AGENT_TOOL_CALLS"},
	{CategoryExecution, KeyAgentPlanningStepEnabled, "AGENT_PLANNING_STEP_ENABLED"},
	{CategoryExecution, KeyAssistantUseAgents, "ASSISTANT_USE_AGENTS"},

	{CategoryLLMProviders, KeyOpenAIKey, "OPEN_AI_KEY"},
	{CategoryLLMProviders, KeyOpenAIServerURL, "OPEN_AI_SERVER_URL"},
	{CategoryLLMProviders, KeyAnthropicAPIKey, "ANTHROPIC_API_KEY"},
	{CategoryLLMProviders, KeyAnthropicServerURL, "ANTHROPIC_SERVER_URL"},
	{CategoryLLMProviders, KeyGeminiAPIKey, "GEMINI_API_KEY"},
	{CategoryLLMProviders, KeyGeminiServerURL, "GEMINI_SERVER_URL"},
	{CategoryLLMProviders, KeyBedrockRegion, "BEDROCK_REGION"},
	{CategoryLLMProviders, KeyBedrockDefaultAuth, "BEDROCK_DEFAULT_AUTH"},
	{CategoryLLMProviders, KeyBedrockBearerToken, "BEDROCK_BEARER_TOKEN"},
	{CategoryLLMProviders, KeyBedrockAccessKey, "BEDROCK_ACCESS_KEY_ID"},
	{CategoryLLMProviders, KeyBedrockSecretKey, "BEDROCK_SECRET_ACCESS_KEY"},
	{CategoryLLMProviders, KeyBedrockSessionToken, "BEDROCK_SESSION_TOKEN"},
	{CategoryLLMProviders, KeyBedrockServerURL, "BEDROCK_SERVER_URL"},
	{CategoryLLMProviders, KeyDeepSeekAPIKey, "DEEPSEEK_API_KEY"},
	{CategoryLLMProviders, KeyDeepSeekServerURL, "DEEPSEEK_SERVER_URL"},
	{CategoryLLMProviders, KeyDeepSeekProvider, "DEEPSEEK_PROVIDER"},
	{CategoryLLMProviders, KeyGLMAPIKey, "GLM_API_KEY"},
	{CategoryLLMProviders, KeyGLMServerURL, "GLM_SERVER_URL"},
	{CategoryLLMProviders, KeyGLMProvider, "GLM_PROVIDER"},
	{CategoryLLMProviders, KeyKimiAPIKey, "KIMI_API_KEY"},
	{CategoryLLMProviders, KeyKimiServerURL, "KIMI_SERVER_URL"},
	{CategoryLLMProviders, KeyKimiProvider, "KIMI_PROVIDER"},
	{CategoryLLMProviders, KeyQwenAPIKey, "QWEN_API_KEY"},
	{CategoryLLMProviders, KeyQwenServerURL, "QWEN_SERVER_URL"},
	{CategoryLLMProviders, KeyQwenProvider, "QWEN_PROVIDER"},
	{CategoryLLMProviders, KeyMiniMaxAPIKey, "MINIMAX_API_KEY"},
	{CategoryLLMProviders, KeyMiniMaxServerURL, "MINIMAX_SERVER_URL"},
	{CategoryLLMProviders, KeyMiniMaxProvider, "MINIMAX_PROVIDER"},
	{CategoryLLMProviders, KeyOllamaServerURL, "OLLAMA_SERVER_URL"},
	{CategoryLLMProviders, KeyOllamaServerAPIKey, "OLLAMA_SERVER_API_KEY"},
	{CategoryLLMProviders, KeyOllamaServerModel, "OLLAMA_SERVER_MODEL"},
	{CategoryLLMProviders, KeyOllamaServerConfig, "OLLAMA_SERVER_CONFIG_PATH"},
	{CategoryLLMProviders, KeyOllamaPullEnabled, "OLLAMA_SERVER_PULL_MODELS_ENABLED"},
	{CategoryLLMProviders, KeyOllamaLoadEnabled, "OLLAMA_SERVER_LOAD_MODELS_ENABLED"},
	{CategoryLLMProviders, KeyOllamaPullTimeout, "OLLAMA_SERVER_PULL_MODELS_TIMEOUT"},
	{CategoryLLMProviders, KeyLLMServerURL, "LLM_SERVER_URL"},
	{CategoryLLMProviders, KeyLLMServerKey, "LLM_SERVER_KEY"},
	{CategoryLLMProviders, KeyLLMServerModel, "LLM_SERVER_MODEL"},
	{CategoryLLMProviders, KeyLLMServerConfig, "LLM_SERVER_CONFIG_PATH"},
	{CategoryLLMProviders, KeyLLMServerProvider, "LLM_SERVER_PROVIDER"},
	{CategoryLLMProviders, KeyLLMServerLegacyReason, "LLM_SERVER_LEGACY_REASONING"},
	{CategoryLLMProviders, KeyLLMServerPreserveReas, "LLM_SERVER_PRESERVE_REASONING"},
}

type settingID struct {
	category string
	key      string
}

// configField describes the Config struct field an environment variable binds
// to, as declared by its `env`/`envDefault` tags.
type configField struct {
	index      int
	defaultVal string
}

var (
	settingsByID   = make(map[settingID]settingDef, len(settingDefs))
	settingsByEnv  = make(map[string]settingDef, len(settingDefs))
	configFieldsBy = configFieldsByEnv()
)

func init() {
	for _, def := range settingDefs {
		if _, ok := configFieldsBy[def.env]; !ok {
			panic(fmt.Sprintf("config: setting %s/%s maps to unknown env var %s", def.category, def.key, def.env))
		}
		settingsByID[settingID{def.category, def.key}] = def
		settingsByEnv[def.env] = def
	}
}

func configFieldsByEnv() map[string]configField {
	t := reflect.TypeOf(Config{})
	fields := make(map[string]configField, t.NumField())
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		name := f.Tag.Get("env")
		if name == "" || name == "-" {
			continue
		}
		fields[name] = configField{index: i, defaultVal: f.Tag.Get("envDefault")}
	}
	return fields
}

// SettingEnvName returns the environment variable backing (category, key).
func SettingEnvName(category, key string) (string, bool) {
	def, ok := settingsByID[settingID{category, key}]
	return def.env, ok
}

// SettingKeys returns the catalogued keys of category, in catalogue order.
func SettingKeys(category string) []string {
	keys := make([]string, 0)
	for _, def := range settingDefs {
		if def.category == category {
			keys = append(keys, def.key)
		}
	}
	return keys
}

// NormalizeSetting maps an empty value to the variable's envDefault, matching
// how env parsing treats an empty variable at boot. Without this, a blank line
// such as `PERPLEXITY_MODEL=` in the .env file (or a cleared Web UI field)
// would override the built-in default with an empty string at runtime while
// the same file yields the default after a restart.
func NormalizeSetting(category, key, value string) string {
	if value != "" {
		return value
	}
	def, ok := settingsByID[settingID{category, key}]
	if !ok {
		return value
	}
	return configFieldsBy[def.env].defaultVal
}

// WithSettings returns a shallow copy of c where every catalogued field of the
// given category that has an entry in values is replaced by that value, parsed
// the same way env parsing does. Callers that construct long-lived objects
// from Config fields (LLM providers, embedders) use it to build from the live
// overrides without mutating the shared *Config other goroutines read.
func (c *Config) WithSettings(category string, values map[string]string) *Config {
	out := *c
	v := reflect.ValueOf(&out).Elem()
	for key, raw := range values {
		def, ok := settingsByID[settingID{category, key}]
		if !ok {
			continue
		}
		field := configFieldsBy[def.env]
		setFieldFromString(v.Field(field.index), NormalizeSetting(category, key, raw))
	}
	return &out
}

// WithOverrides is WithSettings applied to the live overrides of category.
func (c *Config) WithOverrides(category string) *Config {
	return c.WithSettings(category, c.Overrides.Snapshot(category))
}

// setFieldFromString assigns raw to a string/bool/int field; unparsable
// numbers or booleans leave the field at its env-sourced value, mirroring the
// fallback behaviour of Overrides.GetBool/GetInt.
func setFieldFromString(field reflect.Value, raw string) {
	switch field.Kind() {
	case reflect.String:
		field.SetString(raw)
	case reflect.Bool:
		if b, err := strconv.ParseBool(raw); err == nil {
			field.SetBool(b)
		}
	case reflect.Int, reflect.Int64:
		if n, err := strconv.ParseInt(raw, 10, 64); err == nil {
			field.SetInt(n)
		}
	}
}
