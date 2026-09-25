package config

// Category names for runtime (hot-reloadable) settings stored in
// system_settings, mirrored into the settings .env file and cached in
// Config.Overrides. Every (category, key) exposed through the Web UI must also
// be listed in settingDefs (settings.go) so it maps to its environment
// variable. Keep in sync with backend/pkg/graph/schema.graphqls and
// frontend/src/pages/settings/settings-system.tsx.
const (
	CategorySearchEngines = "search_engines"
	CategoryExecution     = "execution"
	CategorySummarizer    = "summarizer"
	CategoryEmbedder      = "embedder"
	CategoryLLMProviders  = "llm_providers"
)

// Keys within CategorySearchEngines.
const (
	KeyDuckDuckGoEnabled    = "duckduckgo_enabled"
	KeyDuckDuckGoRegion     = "duckduckgo_region"
	KeyDuckDuckGoSafeSearch = "duckduckgo_safesearch"
	KeyDuckDuckGoTimeRange  = "duckduckgo_time_range"
	KeyBraveAPIKey          = "brave_api_key"
	KeySploitusEnabled      = "sploitus_enabled"
	KeyGoogleAPIKey         = "google_api_key"
	KeyGoogleCXKey          = "google_cx_key"
	KeyGoogleLRKey          = "google_lr_key"
	KeyTraversaalAPIKey     = "traversaal_api_key"
	KeyTavilyAPIKey         = "tavily_api_key"
	KeyFirecrawlAPIKey      = "firecrawl_api_key"
	KeyFirecrawlAPIURL      = "firecrawl_api_url"
	KeyPerplexityAPIKey     = "perplexity_api_key"
	KeyPerplexityModel      = "perplexity_model"
	KeyPerplexityContext    = "perplexity_context_size"
	KeySearxngURL           = "searxng_url"
	KeySearxngCategories    = "searxng_categories"
	KeySearxngLanguage      = "searxng_language"
	KeySearxngSafeSearch    = "searxng_safesearch"
	KeySearxngTimeRange     = "searxng_time_range"
	KeySearxngTimeout       = "searxng_timeout"
	KeyWebSearchIntEnabled  = "web_search_internal_enabled"
	KeyWebSearchIntMaxSites = "web_search_internal_max_sites"
	KeyWebSearchIntMaxBytes = "web_search_internal_max_site_bytes"
)

// Keys within CategoryExecution.
const (
	KeyExecutionMonitorEnabled  = "execution_monitor_enabled"
	KeyExecutionSameToolLimit   = "execution_monitor_same_tool_limit"
	KeyExecutionTotalToolLimit  = "execution_monitor_total_tool_limit"
	KeyMaxGeneralAgentToolCalls = "max_general_agent_tool_calls"
	KeyMaxLimitedAgentToolCalls = "max_limited_agent_tool_calls"
	KeyAgentPlanningStepEnabled = "agent_planning_step_enabled"
	KeyAssistantUseAgents       = "assistant_use_agents"
	KeyDockerFlowImage          = "docker_flow_image"
	KeyDockerAssistantImage     = "docker_assistant_image"
)

// Keys within CategorySummarizer. General and assistant summarizers share the
// same field shape but are distinct rows, disambiguated by an "assistant_"
// prefix on the assistant variant.
const (
	KeySumPreserveLast   = "preserve_last"
	KeySumUseQA          = "use_qa"
	KeySumHumanInQA      = "sum_human_in_qa"
	KeySumLastSecBytes   = "last_sec_bytes"
	KeySumMaxBPBytes     = "max_bp_bytes"
	KeySumMaxQASections  = "max_qa_sections"
	KeySumMaxQABytes     = "max_qa_bytes"
	KeySumKeepQASections = "keep_qa_sections"

	KeyAssistantSumPrefix = "assistant_"
)

// Keys within CategoryEmbedder.
const (
	KeyEmbeddingURL           = "embedding_url"
	KeyEmbeddingKey           = "embedding_key"
	KeyEmbeddingModel         = "embedding_model"
	KeyEmbeddingStripNewLines = "embedding_strip_new_lines"
	KeyEmbeddingBatchSize     = "embedding_batch_size"
	KeyEmbeddingProvider      = "embedding_provider"
	KeyEmbeddingMaxTextBytes  = "embedding_max_text_bytes"
)

// Keys within CategoryLLMProviders. One key per built-in provider credential
// field; provider type prefixes keep the flat namespace collision-free.
const (
	KeyOpenAIKey       = "openai_key"
	KeyOpenAIServerURL = "openai_server_url"

	KeyAnthropicAPIKey    = "anthropic_api_key"
	KeyAnthropicServerURL = "anthropic_server_url"

	KeyGeminiAPIKey    = "gemini_api_key"
	KeyGeminiServerURL = "gemini_server_url"

	KeyBedrockRegion       = "bedrock_region"
	KeyBedrockDefaultAuth  = "bedrock_default_auth"
	KeyBedrockBearerToken  = "bedrock_bearer_token"
	KeyBedrockAccessKey    = "bedrock_access_key_id"
	KeyBedrockSecretKey    = "bedrock_secret_access_key"
	KeyBedrockSessionToken = "bedrock_session_token"
	KeyBedrockServerURL    = "bedrock_server_url"

	KeyDeepSeekAPIKey    = "deepseek_api_key"
	KeyDeepSeekServerURL = "deepseek_server_url"
	KeyDeepSeekProvider  = "deepseek_provider"

	KeyGLMAPIKey    = "glm_api_key"
	KeyGLMServerURL = "glm_server_url"
	KeyGLMProvider  = "glm_provider"

	KeyKimiAPIKey    = "kimi_api_key"
	KeyKimiServerURL = "kimi_server_url"
	KeyKimiProvider  = "kimi_provider"

	KeyQwenAPIKey    = "qwen_api_key"
	KeyQwenServerURL = "qwen_server_url"
	KeyQwenProvider  = "qwen_provider"

	KeyMiniMaxAPIKey    = "minimax_api_key"
	KeyMiniMaxServerURL = "minimax_server_url"
	KeyMiniMaxProvider  = "minimax_provider"

	KeyOllamaServerURL       = "ollama_server_url"
	KeyOllamaServerAPIKey    = "ollama_server_api_key"
	KeyOllamaServerModel     = "ollama_server_model"
	KeyOllamaServerConfig    = "ollama_server_config_path"
	KeyOllamaPullEnabled     = "ollama_server_pull_models_enabled"
	KeyOllamaLoadEnabled     = "ollama_server_load_models_enabled"
	KeyOllamaPullTimeout     = "ollama_server_pull_models_timeout"
	KeyLLMServerURL          = "llm_server_url"
	KeyLLMServerKey          = "llm_server_key"
	KeyLLMServerModel        = "llm_server_model"
	KeyLLMServerConfig       = "llm_server_config_path"
	KeyLLMServerProvider     = "llm_server_provider"
	KeyLLMServerLegacyReason = "llm_server_legacy_reasoning"
	KeyLLMServerPreserveReas = "llm_server_preserve_reasoning"
)

// SecretKeys lists every key across all categories that must be masked when
// echoed back to the frontend (never returned in plaintext once set).
var SecretKeys = map[string]bool{
	KeyGoogleAPIKey:        true,
	KeyBraveAPIKey:         true,
	KeyTraversaalAPIKey:    true,
	KeyTavilyAPIKey:        true,
	KeyFirecrawlAPIKey:     true,
	KeyPerplexityAPIKey:    true,
	KeyEmbeddingKey:        true,
	KeyOpenAIKey:           true,
	KeyAnthropicAPIKey:     true,
	KeyGeminiAPIKey:        true,
	KeyBedrockBearerToken:  true,
	KeyBedrockAccessKey:    true,
	KeyBedrockSecretKey:    true,
	KeyBedrockSessionToken: true,
	KeyDeepSeekAPIKey:      true,
	KeyGLMAPIKey:           true,
	KeyKimiAPIKey:          true,
	KeyQwenAPIKey:          true,
	KeyMiniMaxAPIKey:       true,
	KeyOllamaServerAPIKey:  true,
	KeyLLMServerKey:        true,
}
