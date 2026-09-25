package providers

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"sync"
	"testing"

	"pentagi/pkg/config"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// authRecorder is a minimal OpenAI-compatible endpoint remembering the bearer
// token of every chat completion request.
type authRecorder struct {
	mu    sync.Mutex
	auths []string
}

func (a *authRecorder) handler(w http.ResponseWriter, r *http.Request) {
	a.mu.Lock()
	a.auths = append(a.auths, r.Header.Get("Authorization"))
	a.mu.Unlock()

	w.Header().Set("Content-Type", "application/json")
	_, _ = w.Write([]byte(`{"id":"1","object":"chat.completion","created":0,"model":"m",` +
		`"choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}],` +
		`"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`))
}

func (a *authRecorder) last() string {
	a.mu.Lock()
	defer a.mu.Unlock()
	if len(a.auths) == 0 {
		return ""
	}
	return a.auths[len(a.auths)-1]
}

func newHotReloadController(t *testing.T, serverURL string) *providerController {
	t.Helper()
	cfg := &config.Config{
		OpenAIKey:       "sk-old",
		OpenAIServerURL: serverURL,
		Overrides:       config.NewOverrides(),
	}
	pc := &providerController{cfg: cfg, db: stubProvidersQuerier{}}
	pc.embedder = liveEmbedder{pc: pc}
	pc.current()
	return pc
}

// A flow keeps the provider it got at creation for its whole lifetime; a key
// saved afterwards must be used by that same provider's next call.
func TestLiveProvider_HeldProviderUsesUpdatedCredentials(t *testing.T) {
	rec := &authRecorder{}
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	pc := newHotReloadController(t, srv.URL)
	held, err := pc.defaultProvider(pc.current(), provider.DefaultProviderNameOpenAI)
	require.NoError(t, err)

	_, err = held.Call(context.Background(), pconfig.OptionsTypeSimple, "hi")
	require.NoError(t, err)
	assert.Equal(t, "Bearer sk-old", rec.last())

	pc.cfg.Overrides.Set(config.CategoryLLMProviders, config.KeyOpenAIKey, "sk-new")

	_, err = held.Call(context.Background(), pconfig.OptionsTypeSimple, "hi")
	require.NoError(t, err)
	assert.Equal(t, "Bearer sk-new", rec.last(), "the held provider must switch to the new key without a restart")
}

func TestLiveProvider_ClearedKeyDisablesHeldProvider(t *testing.T) {
	rec := &authRecorder{}
	srv := httptest.NewServer(http.HandlerFunc(rec.handler))
	defer srv.Close()

	pc := newHotReloadController(t, srv.URL)
	held, err := pc.defaultProvider(pc.current(), provider.DefaultProviderNameOpenAI)
	require.NoError(t, err)

	pc.cfg.Overrides.Set(config.CategoryLLMProviders, config.KeyOpenAIKey, "")

	_, err = held.Call(context.Background(), pconfig.OptionsTypeSimple, "hi")
	require.Error(t, err, "a provider whose key was cleared must stop serving calls")
	active, _ := pc.DefaultProviderStatus(provider.ProviderOpenAI)
	assert.False(t, active)
}

// Rebuilding is scoped to the provider whose settings changed: Ollama/custom
// construction can pull models or hit the network, so unrelated saves must not
// rebuild them.
func TestBuildProviderState_ReusesProvidersWithUnchangedSettings(t *testing.T) {
	pc := newHotReloadController(t, "http://127.0.0.1:1/v1")
	before := pc.current()

	pc.cfg.Overrides.Set(config.CategoryLLMProviders, config.KeyAnthropicAPIKey, "sk-ant")
	after := pc.current()

	require.NotSame(t, before, after)
	assert.Same(t, before.providers[provider.DefaultProviderNameOpenAI], after.providers[provider.DefaultProviderNameOpenAI])
	assert.Equal(t, before.versions[provider.ProviderOpenAI], after.versions[provider.ProviderOpenAI])
	assert.Contains(t, after.providers, provider.DefaultProviderNameAnthropic)
	assert.NotEqual(t, before.versions[provider.ProviderAnthropic], after.versions[provider.ProviderAnthropic])
}

func TestApplyLLMProviderSettings_RejectsUnbuildableChangeWithoutPersisting(t *testing.T) {
	pc := newHotReloadController(t, "http://127.0.0.1:1/v1")
	before := pc.current()

	candidate := map[string]string{
		config.KeyLLMServerURL:    "http://127.0.0.1:1/v1",
		config.KeyLLMServerConfig: filepath.Join(t.TempDir(), "missing.yml"),
	}
	persisted := false
	err := pc.ApplyLLMProviderSettings(candidate, func() error {
		persisted = true
		return nil
	})

	require.Error(t, err)
	assert.False(t, persisted, "a change that cannot build its provider must not be saved")
	assert.Same(t, before, pc.state.Load(), "the running providers stay untouched")
}

func TestApplyLLMProviderSettings_PersistFailureKeepsRunningProviders(t *testing.T) {
	pc := newHotReloadController(t, "http://127.0.0.1:1/v1")
	before := pc.current()

	err := pc.ApplyLLMProviderSettings(map[string]string{config.KeyOpenAIKey: "sk-new"}, func() error {
		return errors.New("db down")
	})

	require.Error(t, err)
	assert.Same(t, before, pc.state.Load())
}

// Every llm_providers setting must be declared by the registry entry it
// feeds; otherwise saving it would never rebuild the provider.
func TestProviderRegistryCoversEveryLLMProviderSetting(t *testing.T) {
	declared := make(map[string]bool)
	for _, e := range providerRegistry {
		for _, key := range e.Keys {
			declared[key] = true
		}
	}

	for _, key := range config.SettingKeys(config.CategoryLLMProviders) {
		assert.Truef(t, declared[key], "llm_providers setting %q is not listed in any providerRegistry entry's Keys", key)
	}
	for key := range declared {
		_, ok := config.SettingEnvName(config.CategoryLLMProviders, key)
		assert.Truef(t, ok, "providerRegistry key %q is not a catalogued llm_providers setting", key)
	}
}
