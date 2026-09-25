package providers

import (
	"errors"
	"fmt"
	"maps"

	"pentagi/pkg/config"
	"pentagi/pkg/providers/embeddings"
	"pentagi/pkg/providers/provider"

	"github.com/sirupsen/logrus"
)

// providerState is an immutable snapshot of the built-in providers (the ones
// configured through OPEN_AI_KEY & co., from the environment or from the
// llm_providers runtime settings) and everything built from the same inputs.
// The controller swaps whole snapshots, so readers never see a half-rebuilt
// set of providers.
type providerState struct {
	// settings is the llm_providers overrides snapshot this state was built from.
	settings map[string]string
	// cfg is the base config with settings applied; every provider of this
	// state, built-in or user-defined, is constructed from it.
	cfg     *config.Config
	configs provider.ProvidersConfig
	// configErrors records why a provider type has no default config, so
	// patchProviderConfig can report the real cause (e.g. an unreadable
	// *_CONFIG_PATH) instead of an ambiguous "not found".
	configErrors map[provider.ProviderType]error
	// buildErrors records why an enabled (configured) provider type could not
	// be built; the type is then absent from providers.
	buildErrors map[provider.ProviderType]error
	providers   provider.Providers
	embedder    embeddings.Embedder
	// versions changes for a provider type whenever its inputs change, telling
	// live providers of that type (built-in and user-defined) to rebuild.
	versions map[provider.ProviderType]uint64
}

// embedderSettingKeys are the llm_providers settings the embedder falls back
// to when EMBEDDING_KEY / EMBEDDING_URL are not set.
var embedderSettingKeys = []string{config.KeyOpenAIKey, config.KeyOpenAIServerURL}

func settingsChanged(keys []string, prev, next map[string]string) bool {
	for _, key := range keys {
		pv, pok := prev[key]
		nv, nok := next[key]
		if pok != nok || pv != nv {
			return true
		}
	}
	return false
}

// buildProviderState builds the state for settings, reusing from prev every
// provider (and the embedder) whose inputs did not change. It never fails as
// a whole: an enabled provider that cannot be built is left out and its error
// recorded, so one bad credential can neither abort startup nor lock out the
// Web UI that would fix it.
func buildProviderState(base *config.Config, settings map[string]string, prev *providerState) *providerState {
	cfg := base.WithSettings(config.CategoryLLMProviders, settings)
	st := &providerState{
		settings:     maps.Clone(settings),
		cfg:          cfg,
		configs:      make(provider.ProvidersConfig),
		configErrors: make(map[provider.ProviderType]error),
		buildErrors:  make(map[provider.ProviderType]error),
		providers:    make(provider.Providers),
		versions:     make(map[provider.ProviderType]uint64),
	}

	for _, e := range providerRegistry {
		if prev != nil && !settingsChanged(e.Keys, prev.settings, settings) {
			st.versions[e.Type] = prev.versions[e.Type]
			if c, ok := prev.configs[e.Type]; ok {
				st.configs[e.Type] = c
			}
			if err, ok := prev.configErrors[e.Type]; ok {
				st.configErrors[e.Type] = err
			}
			if err, ok := prev.buildErrors[e.Type]; ok {
				st.buildErrors[e.Type] = err
			}
			if p, ok := prev.providers[e.Name]; ok {
				st.providers[e.Name] = p
			}
			continue
		}

		if prev != nil {
			st.versions[e.Type] = prev.versions[e.Type] + 1
		}

		enabled := e.Enabled(cfg)
		pcfg, err := e.NewConfig(cfg)
		if err != nil {
			st.configErrors[e.Type] = err
			if enabled {
				st.buildErrors[e.Type] = fmt.Errorf("failed to create %s provider config: %w", e.Type, err)
				logrus.WithError(err).Errorf("%s provider is configured but its config cannot be loaded", e.Type)
			} else {
				logrus.WithError(err).Warnf("skipping config for disabled %s provider", e.Type)
			}
			continue
		}
		st.configs[e.Type] = pcfg

		if !enabled {
			continue
		}

		p, err := e.New(cfg, e.Name, pcfg)
		if err != nil {
			st.buildErrors[e.Type] = fmt.Errorf("failed to create %s provider: %w", e.Type, err)
			logrus.WithError(err).Errorf("%s provider is configured but cannot be built", e.Type)
			continue
		}
		st.providers[e.Name] = p
	}

	if prev != nil && !settingsChanged(embedderSettingKeys, prev.settings, settings) {
		st.embedder = prev.embedder
	} else {
		embedder, err := embeddings.New(cfg)
		if err != nil {
			logrus.WithError(err).Errorf("failed to create embedder '%s'", cfg.EmbeddingProvider)
		}
		st.embedder = embedder
	}

	return st
}

// current returns the state matching the live llm_providers overrides,
// rebuilding it (only the providers whose settings changed) when the overrides
// moved since the last build — e.g. after the periodic reload picked up a .env
// edit made in the installer TUI.
func (pc *providerController) current() *providerState {
	settings := pc.cfg.Overrides.Snapshot(config.CategoryLLMProviders)
	if st := pc.state.Load(); st != nil && maps.Equal(st.settings, settings) {
		return st
	}

	pc.reloadMu.Lock()
	defer pc.reloadMu.Unlock()

	prev := pc.state.Load()
	if prev != nil && maps.Equal(prev.settings, settings) {
		return prev
	}

	next := buildProviderState(pc.cfg, settings, prev)
	pc.state.Store(next)
	return next
}

// ApplyLLMProviderSettings validates candidate llm_providers settings by
// building every provider whose inputs they change, then runs persist and
// publishes the new providers only if both succeed. On a build failure nothing
// is persisted and the joined build errors are returned. settings must be the
// complete, normalized (config.NormalizeSetting) llm_providers snapshot that
// persist will leave in the overrides, so current() does not rebuild again.
func (pc *providerController) ApplyLLMProviderSettings(settings map[string]string, persist func() error) error {
	pc.reloadMu.Lock()
	defer pc.reloadMu.Unlock()

	prev := pc.state.Load()
	next := buildProviderState(pc.cfg, settings, prev)

	var errs []error
	for _, e := range providerRegistry {
		if prev != nil && next.versions[e.Type] == prev.versions[e.Type] {
			continue // untouched by this change; a pre-existing failure must not block it
		}
		if err, ok := next.buildErrors[e.Type]; ok {
			errs = append(errs, err)
		}
	}
	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	if persist != nil {
		if err := persist(); err != nil {
			return err
		}
	}

	pc.state.Store(next)
	return nil
}

// DefaultProviderStatus reports whether the built-in provider of type
// prvtype is currently built and usable and, when it is configured but could
// not be built, why.
func (pc *providerController) DefaultProviderStatus(prvtype provider.ProviderType) (bool, error) {
	st := pc.current()
	e, ok := entryForType(prvtype)
	if !ok {
		return false, fmt.Errorf("unknown provider type: %s", prvtype)
	}
	_, active := st.providers[e.Name]
	return active, st.buildErrors[prvtype]
}
