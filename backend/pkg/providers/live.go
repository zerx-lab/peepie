package providers

import (
	"context"
	"errors"
	"sync"

	"pentagi/pkg/providers/embeddings"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/templates"

	"github.com/vxcontrol/langchaingo/llms"
	"github.com/vxcontrol/langchaingo/llms/streaming"
)

// liveProvider is the provider.Provider handed out by GetProvider,
// GetProviders and NewProvider. Flows and assistants keep their provider for
// their whole lifetime, so instead of a fixed instance they get this wrapper,
// which re-resolves the concrete provider whenever the settings of its type
// change: a credential or endpoint saved in the Web UI (or the .env file) is
// used by the very next LLM call, without restarting anything.
type liveProvider struct {
	pc    *providerController
	name  provider.ProviderName
	ptype provider.ProviderType
	build func(*providerState) (provider.Provider, error)

	mu      sync.Mutex
	version uint64
	inner   provider.Provider
}

// newLiveProvider builds the concrete provider from st right away, so
// construction errors surface to the caller exactly as before.
func (pc *providerController) newLiveProvider(
	st *providerState,
	prvname provider.ProviderName,
	prvtype provider.ProviderType,
	build func(*providerState) (provider.Provider, error),
) (provider.Provider, error) {
	inner, err := build(st)
	if err != nil {
		return nil, err
	}

	return &liveProvider{
		pc:      pc,
		name:    prvname,
		ptype:   prvtype,
		build:   build,
		version: st.versions[prvtype],
		inner:   inner,
	}, nil
}

// resolve returns the concrete provider for the current settings, rebuilding
// it when its type's settings changed since the last call.
func (lp *liveProvider) resolve() (provider.Provider, error) {
	st := lp.pc.current()
	version := st.versions[lp.ptype]

	lp.mu.Lock()
	defer lp.mu.Unlock()

	if lp.inner != nil && lp.version == version {
		return lp.inner, nil
	}

	inner, err := lp.build(st)
	if err != nil {
		return nil, err
	}
	lp.inner, lp.version = inner, version

	return inner, nil
}

// peek is resolve for informational methods that cannot report an error: it
// falls back to the last successfully built provider (nil if none).
func (lp *liveProvider) peek() provider.Provider {
	if inner, err := lp.resolve(); err == nil {
		return inner
	}

	lp.mu.Lock()
	defer lp.mu.Unlock()

	return lp.inner
}

func (lp *liveProvider) Type() provider.ProviderType { return lp.ptype }
func (lp *liveProvider) Name() provider.ProviderName { return lp.name }

func (lp *liveProvider) Model(opt pconfig.ProviderOptionsType) string {
	if inner := lp.peek(); inner != nil {
		return inner.Model(opt)
	}
	return ""
}

func (lp *liveProvider) ModelWithPrefix(opt pconfig.ProviderOptionsType) string {
	if inner := lp.peek(); inner != nil {
		return inner.ModelWithPrefix(opt)
	}
	return ""
}

func (lp *liveProvider) GetUsage(info map[string]any) pconfig.CallUsage {
	if inner := lp.peek(); inner != nil {
		return inner.GetUsage(info)
	}
	return pconfig.CallUsage{}
}

func (lp *liveProvider) Call(ctx context.Context, opt pconfig.ProviderOptionsType, prompt string) (string, error) {
	inner, err := lp.resolve()
	if err != nil {
		return "", err
	}
	return inner.Call(ctx, opt, prompt)
}

func (lp *liveProvider) CallEx(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chain []llms.MessageContent,
	streamCb streaming.Callback,
) (*llms.ContentResponse, error) {
	inner, err := lp.resolve()
	if err != nil {
		return nil, err
	}
	return inner.CallEx(ctx, opt, chain, streamCb)
}

func (lp *liveProvider) CallWithTools(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chain []llms.MessageContent,
	tools []llms.Tool,
	streamCb streaming.Callback,
) (*llms.ContentResponse, error) {
	inner, err := lp.resolve()
	if err != nil {
		return nil, err
	}
	return inner.CallWithTools(ctx, opt, chain, tools, streamCb)
}

func (lp *liveProvider) CallWithExtraOptions(
	ctx context.Context,
	opt pconfig.ProviderOptionsType,
	chain []llms.MessageContent,
	tools []llms.Tool,
	streamCb streaming.Callback,
	extra ...llms.CallOption,
) (*llms.ContentResponse, error) {
	inner, err := lp.resolve()
	if err != nil {
		return nil, err
	}
	return inner.CallWithExtraOptions(ctx, opt, chain, tools, streamCb, extra...)
}

func (lp *liveProvider) GetRawConfig() []byte {
	if inner := lp.peek(); inner != nil {
		return inner.GetRawConfig()
	}
	return nil
}

func (lp *liveProvider) GetProviderConfig() *pconfig.ProviderConfig {
	if inner := lp.peek(); inner != nil {
		return inner.GetProviderConfig()
	}
	return nil
}

func (lp *liveProvider) GetPriceInfo(opt pconfig.ProviderOptionsType) *pconfig.PriceInfo {
	if inner := lp.peek(); inner != nil {
		return inner.GetPriceInfo(opt)
	}
	return nil
}

func (lp *liveProvider) GetModels() pconfig.ModelsConfig {
	if inner := lp.peek(); inner != nil {
		return inner.GetModels()
	}
	return nil
}

func (lp *liveProvider) GetToolCallIDTemplate(ctx context.Context, prompter templates.Prompter) (string, error) {
	inner, err := lp.resolve()
	if err != nil {
		return "", err
	}
	return inner.GetToolCallIDTemplate(ctx, prompter)
}

var errEmbedderUnavailable = errors.New("embedder is not available")

// liveEmbedder delegates to the embedder of the current provider state, which
// is rebuilt when the OpenAI credentials it falls back to change.
type liveEmbedder struct {
	pc *providerController
}

func (le liveEmbedder) current() embeddings.Embedder {
	return le.pc.current().embedder
}

func (le liveEmbedder) EmbedDocuments(ctx context.Context, texts []string) ([][]float32, error) {
	e := le.current()
	if e == nil {
		return nil, errEmbedderUnavailable
	}
	return e.EmbedDocuments(ctx, texts)
}

func (le liveEmbedder) EmbedQuery(ctx context.Context, text string) ([]float32, error) {
	e := le.current()
	if e == nil {
		return nil, errEmbedderUnavailable
	}
	return e.EmbedQuery(ctx, text)
}

func (le liveEmbedder) IsAvailable() bool {
	e := le.current()
	return e != nil && e.IsAvailable()
}
