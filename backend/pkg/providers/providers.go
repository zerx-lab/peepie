package providers

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"pentagi/pkg/config"
	"pentagi/pkg/csum"
	"pentagi/pkg/database"
	"pentagi/pkg/docker"
	"pentagi/pkg/graphiti"
	obs "pentagi/pkg/observability"
	"pentagi/pkg/providers/embeddings"
	"pentagi/pkg/providers/pconfig"
	"pentagi/pkg/providers/provider"
	"pentagi/pkg/providers/tester"
	"pentagi/pkg/templates"
	"pentagi/pkg/tools"

	lru "github.com/hashicorp/golang-lru/v2"
	"github.com/sirupsen/logrus"
)

const deltaCallCounter = 10000

const defaultTestParallelWorkersNumber = 16

const summarizerCacheMaxSize = 100

// newSummarizerCache creates a fixed-size LRU cache for summarizer results.
// Using lru.Cache (non-expirable) instead of expirable.LRU avoids spawning a
// per-instance background goroutine that would otherwise accumulate for every
// short-lived flow/assistant provider.
func newSummarizerCache() *lru.Cache[[32]byte, string] {
	c, _ := lru.New[[32]byte, string](summarizerCacheMaxSize)
	return c
}

const pentestDockerImage = "vxcontrol/kali-linux"

type ProviderController interface {
	NewFlowProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		flowID, userID int64,
		askUser bool,
		input string,
	) (FlowProvider, error)
	LoadFlowProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		flowID, userID int64,
		askUser bool,
		image, language, title, tcIDTemplate string,
	) (FlowProvider, error)
	NewAssistantProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		assistantID, flowID, userID int64,
		image, input string,
		streamCb StreamMessageHandler,
	) (AssistantProvider, error)
	LoadAssistantProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		prompter templates.Prompter,
		executor tools.FlowToolsExecutor,
		assistantID, flowID, userID int64,
		image, language, title, tcIDTemplate string,
		streamCb StreamMessageHandler,
	) (AssistantProvider, error)

	Embedder() embeddings.Embedder
	GraphitiClient() *graphiti.Client
	DefaultProviders() provider.Providers
	DefaultProvidersConfig() provider.ProvidersConfig
	GetProvider(
		ctx context.Context,
		prvname provider.ProviderName,
		userID int64,
	) (provider.Provider, error)
	GetProviders(
		ctx context.Context,
		userID int64,
	) (provider.Providers, error)

	NewProvider(prv database.Provider) (provider.Provider, error)
	CreateProvider(
		ctx context.Context,
		userID int64,
		prvname provider.ProviderName,
		prvtype provider.ProviderType,
		config *pconfig.ProviderConfig,
	) (database.Provider, error)
	UpdateProvider(
		ctx context.Context,
		userID int64,
		prvID int64,
		prvname provider.ProviderName,
		config *pconfig.ProviderConfig,
	) (database.Provider, error)
	DeleteProvider(
		ctx context.Context,
		userID int64,
		prvID int64,
	) (database.Provider, error)

	SeedDefaultProviders(ctx context.Context, userID int64) error

	// ApplyLLMProviderSettings validates and publishes a new llm_providers
	// settings snapshot; see state.go.
	ApplyLLMProviderSettings(settings map[string]string, persist func() error) error
	DefaultProviderStatus(prvtype provider.ProviderType) (bool, error)

	TestAgent(
		ctx context.Context,
		prvtype provider.ProviderType,
		agentType pconfig.ProviderOptionsType,
		config *pconfig.AgentConfig,
	) (tester.AgentTestResults, error)
	TestProvider(
		ctx context.Context,
		prvtype provider.ProviderType,
		config *pconfig.ProviderConfig,
	) (tester.ProviderTestResults, error)
}

type providerController struct {
	db             database.Querier
	cfg            *config.Config
	docker         docker.DockerClient
	embedder       embeddings.Embedder
	graphitiClient *graphiti.Client

	startCallNumber *atomic.Int64

	summarizerAgent     csum.Summarizer
	summarizerAssistant csum.Summarizer

	// state holds the built-in providers built from the live llm_providers
	// settings; reloadMu serializes its rebuilds (see state.go).
	state    atomic.Pointer[providerState]
	reloadMu sync.Mutex
}

func NewProviderController(
	cfg *config.Config,
	db database.Querier,
	docker docker.DockerClient,
) (ProviderController, error) {
	if cfg == nil {
		return nil, fmt.Errorf("config is required")
	}

	summarizerAgent := csum.NewSummarizer(csum.SummarizerConfig{
		PreserveLast:   cfg.SummarizerPreserveLast,
		UseQA:          cfg.SummarizerUseQA,
		SummHumanInQA:  cfg.SummarizerSumHumanInQA,
		LastSecBytes:   cfg.SummarizerLastSecBytes,
		MaxBPBytes:     cfg.SummarizerMaxBPBytes,
		MaxQASections:  cfg.SummarizerMaxQASections,
		MaxQABytes:     cfg.SummarizerMaxQABytes,
		KeepQASections: cfg.SummarizerKeepQASections,
	})

	summarizerAssistant := csum.NewSummarizer(csum.SummarizerConfig{
		PreserveLast:   cfg.AssistantSummarizerPreserveLast,
		UseQA:          true,
		SummHumanInQA:  false,
		LastSecBytes:   cfg.AssistantSummarizerLastSecBytes,
		MaxBPBytes:     cfg.AssistantSummarizerMaxBPBytes,
		MaxQASections:  cfg.AssistantSummarizerMaxQASections,
		MaxQABytes:     cfg.AssistantSummarizerMaxQABytes,
		KeepQASections: cfg.AssistantSummarizerKeepQASections,
	})

	graphitiClient, err := graphiti.NewClient(
		cfg.GraphitiURL,
		time.Duration(cfg.GraphitiTimeout)*time.Second,
		cfg.GraphitiEnabled && cfg.GraphitiURL != "",
	)
	if err != nil {
		logrus.WithError(err).Warn("failed to initialize graphiti client, continuing without it")
		graphitiClient = &graphiti.Client{}
	}

	pc := &providerController{
		db:             db,
		cfg:            cfg,
		docker:         docker,
		graphitiClient: graphitiClient,

		startCallNumber: newAtomicInt64(0), // 0 means to make it random

		summarizerAgent:     summarizerAgent,
		summarizerAssistant: summarizerAssistant,
	}
	pc.embedder = liveEmbedder{pc: pc}
	// Build eagerly so configuration problems are logged at startup; a
	// configured provider that fails to build is reported through
	// DefaultProviderStatus instead of aborting the process.
	pc.current()

	// Seed configured system providers into the DB for all existing users so
	// they are immediately available without any UI interaction. This runs on
	// every startup, so editing a YAML config file and restarting PentAGI
	// automatically propagates new values.
	ctx := context.Background()
	if users, err := db.GetUsers(ctx); err != nil {
		logrus.WithError(err).Warn("failed to fetch users for provider seeding")
	} else {
		for _, u := range users {
			if err := pc.SeedDefaultProviders(ctx, u.ID); err != nil {
				logrus.WithError(err).Warnf("failed to seed default providers for user %d", u.ID)
			}
		}
	}

	return pc, nil
}

// newExecutionMonitor builds an executionMonitor from the live (override-aware)
// execution settings. Called fresh at the start of every agent chain iteration,
// so a Web UI change to these fields takes effect on the very next iteration
// without restarting the process or the flow.
func (pc *providerController) newExecutionMonitor() *executionMonitor {
	return &executionMonitor{
		enabled: pc.cfg.Overrides.GetBool(
			config.CategoryExecution, config.KeyExecutionMonitorEnabled, pc.cfg.ExecutionMonitorEnabled),
		sameThreshold: pc.cfg.Overrides.GetInt(
			config.CategoryExecution, config.KeyExecutionSameToolLimit, pc.cfg.ExecutionMonitorSameToolLimit),
		totalThreshold: pc.cfg.Overrides.GetInt(
			config.CategoryExecution, config.KeyExecutionTotalToolLimit, pc.cfg.ExecutionMonitorTotalToolLimit),
	}
}

// planningEnabled, maxGeneralAgentToolCalls and maxLimitedAgentToolCalls are
// read at FlowProvider/AssistantProvider construction time (i.e. once per
// flow/task/assistant), so a Web UI change takes effect for the next one
// created without a restart.
func (pc *providerController) planningEnabled() bool {
	return pc.cfg.Overrides.GetBool(
		config.CategoryExecution, config.KeyAgentPlanningStepEnabled, pc.cfg.AgentPlanningStepEnabled)
}

func (pc *providerController) maxGeneralAgentToolCalls() int {
	return pc.cfg.Overrides.GetInt(
		config.CategoryExecution, config.KeyMaxGeneralAgentToolCalls, pc.cfg.MaxGeneralAgentToolCalls)
}

func (pc *providerController) maxLimitedAgentToolCalls() int {
	return pc.cfg.Overrides.GetInt(
		config.CategoryExecution, config.KeyMaxLimitedAgentToolCalls, pc.cfg.MaxLimitedAgentToolCalls)
}

func (pc *providerController) NewFlowProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	flowID, userID int64,
	askUser bool,
	input string,
) (FlowProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.NewFlowProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	imageTmpl, err := prompter.RenderTemplate(templates.PromptTypeImageChooser, map[string]any{
		"DefaultImage":           pc.docker.GetDefaultImage(),
		"DefaultImageForPentest": pc.cfg.DockerDefaultImageForPentest,
		"Input":                  input,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get primary docker image template: %w", err)
	}

	image, err := callWithSetupRetries(ctx, prv, pconfig.OptionsTypeSimple, imageTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to select primary docker image via llm call: %w", err)
	}
	image = strings.ToLower(strings.TrimSpace(image))

	languageTmpl, err := prompter.RenderTemplate(templates.PromptTypeLanguageChooser, map[string]any{
		"Input": input,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get language template: %w", err)
	}

	language, err := callWithSetupRetries(ctx, prv, pconfig.OptionsTypeSimple, languageTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to get language: %w", err)
	}
	language = strings.TrimSpace(language)

	titleTmpl, err := prompter.RenderTemplate(templates.PromptTypeFlowDescriptor, map[string]any{
		"Input":       input,
		"Lang":        language,
		"CurrentTime": getCurrentTime(),
		"N":           20,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get flow title template: %w", err)
	}

	title, err := callWithSetupRetries(ctx, prv, pconfig.OptionsTypeSimple, titleTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow title: %w", err)
	}
	title = strings.TrimSpace(title)

	tcIDTemplate, err := prv.GetToolCallIDTemplate(ctx, prompter)
	if err != nil {
		return nil, wrapToolCallIDTemplateError(err)
	}

	fp := &flowProvider{
		db:              pc.db,
		mx:              &sync.RWMutex{},
		cfg:             pc.cfg,
		embedder:        pc.embedder,
		graphitiClient:  pc.graphitiClient,
		flowID:          flowID,
		callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
		image:           image,
		title:           title,
		language:        language,
		askUser:         askUser,
		planning:        pc.planningEnabled(),
		tcIDTemplate:    tcIDTemplate,
		prompter:        prompter,
		executor:        executor,
		summarizer:      pc.summarizerAgent,
		summarizerCache: newSummarizerCache(),
		Provider:        prv,
		maxGACallsLimit: pc.maxGeneralAgentToolCalls(),
		maxLACallsLimit: pc.maxLimitedAgentToolCalls(),
		buildMonitor:    pc.newExecutionMonitor,
	}

	return fp, nil
}

func (pc *providerController) LoadFlowProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	flowID, userID int64,
	askUser bool,
	image, language, title, tcIDTemplate string,
) (FlowProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.LoadFlowProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	fp := &flowProvider{
		db:              pc.db,
		mx:              &sync.RWMutex{},
		cfg:             pc.cfg,
		embedder:        pc.embedder,
		graphitiClient:  pc.graphitiClient,
		flowID:          flowID,
		callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
		image:           image,
		title:           title,
		language:        language,
		askUser:         askUser,
		planning:        pc.planningEnabled(),
		tcIDTemplate:    tcIDTemplate,
		prompter:        prompter,
		executor:        executor,
		summarizer:      pc.summarizerAgent,
		summarizerCache: newSummarizerCache(),
		Provider:        prv,
		maxGACallsLimit: pc.maxGeneralAgentToolCalls(),
		maxLACallsLimit: pc.maxLimitedAgentToolCalls(),
		buildMonitor:    pc.newExecutionMonitor,
	}

	return fp, nil
}

func (pc *providerController) Embedder() embeddings.Embedder {
	return pc.embedder
}

func (pc *providerController) GraphitiClient() *graphiti.Client {
	return pc.graphitiClient
}

func (pc *providerController) NewAssistantProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	assistantID, flowID, userID int64,
	image, input string,
	streamCb StreamMessageHandler,
) (AssistantProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.NewAssistantProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	languageTmpl, err := prompter.RenderTemplate(templates.PromptTypeLanguageChooser, map[string]any{
		"Input": input,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get language template: %w", err)
	}

	language, err := callWithSetupRetries(ctx, prv, pconfig.OptionsTypeSimple, languageTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to get language: %w", err)
	}
	language = strings.TrimSpace(language)

	titleTmpl, err := prompter.RenderTemplate(templates.PromptTypeFlowDescriptor, map[string]any{
		"Input":       input,
		"Lang":        language,
		"CurrentTime": getCurrentTime(),
		"N":           20,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to get flow title template: %w", err)
	}

	title, err := callWithSetupRetries(ctx, prv, pconfig.OptionsTypeSimple, titleTmpl)
	if err != nil {
		return nil, fmt.Errorf("failed to get flow title: %w", err)
	}
	title = strings.TrimSpace(title)

	tcIDTemplate, err := prv.GetToolCallIDTemplate(ctx, prompter)
	if err != nil {
		return nil, wrapToolCallIDTemplateError(err)
	}

	ap := &assistantProvider{
		id:         assistantID,
		summarizer: pc.summarizerAssistant,
		fp: flowProvider{
			db:              pc.db,
			mx:              &sync.RWMutex{},
			cfg:             pc.cfg,
			embedder:        pc.embedder,
			graphitiClient:  pc.graphitiClient,
			flowID:          flowID,
			callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
			image:           image,
			title:           title,
			language:        language,
			tcIDTemplate:    tcIDTemplate,
			prompter:        prompter,
			executor:        executor,
			streamCb:        streamCb,
			summarizer:      pc.summarizerAgent,
			summarizerCache: newSummarizerCache(),
			Provider:        prv,
			maxGACallsLimit: pc.maxGeneralAgentToolCalls(),
			maxLACallsLimit: pc.maxLimitedAgentToolCalls(),
			buildMonitor:    pc.newExecutionMonitor,
		},
	}

	return ap, nil
}

func (pc *providerController) LoadAssistantProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	prompter templates.Prompter,
	executor tools.FlowToolsExecutor,
	assistantID, flowID, userID int64,
	image, language, title, tcIDTemplate string,
	streamCb StreamMessageHandler,
) (AssistantProvider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.LoadAssistantProvider")
	defer span.End()

	prv, err := pc.GetProvider(ctx, prvname, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get provider: %w", err)
	}

	ap := &assistantProvider{
		id:         assistantID,
		summarizer: pc.summarizerAssistant,
		fp: flowProvider{
			db:              pc.db,
			mx:              &sync.RWMutex{},
			cfg:             pc.cfg,
			embedder:        pc.embedder,
			graphitiClient:  pc.graphitiClient,
			flowID:          flowID,
			callCounter:     newAtomicInt64(pc.startCallNumber.Add(deltaCallCounter)),
			image:           image,
			title:           title,
			language:        language,
			tcIDTemplate:    tcIDTemplate,
			prompter:        prompter,
			executor:        executor,
			streamCb:        streamCb,
			summarizer:      pc.summarizerAgent,
			summarizerCache: newSummarizerCache(),
			Provider:        prv,
			maxGACallsLimit: pc.maxGeneralAgentToolCalls(),
			maxLACallsLimit: pc.maxLimitedAgentToolCalls(),
			buildMonitor:    pc.newExecutionMonitor,
		},
	}

	return ap, nil
}

func (pc *providerController) DefaultProviders() provider.Providers {
	return pc.current().providers
}

func (pc *providerController) DefaultProvidersConfig() provider.ProvidersConfig {
	return pc.current().configs
}

func (pc *providerController) GetProvider(
	ctx context.Context,
	prvname provider.ProviderName,
	userID int64,
) (provider.Provider, error) {
	// Lookup user defined providers first so they take precedence over built-in providers
	prv, err := pc.db.GetUserProviderByName(ctx, database.GetUserProviderByNameParams{
		Name:   string(prvname),
		UserID: userID,
	})
	if err != nil && err != sql.ErrNoRows {
		return nil, fmt.Errorf("failed to get provider '%s' from database: %w", prvname, err)
	}
	if err == nil {
		return pc.NewProvider(prv)
	}

	// Fall back to built-in default providers
	return pc.defaultProvider(pc.current(), prvname)
}

func (pc *providerController) GetProviders(
	ctx context.Context,
	userID int64,
) (provider.Providers, error) {
	st := pc.current()
	providersMap := make(provider.Providers, len(st.providers))

	// Copy default providers
	for prvname := range st.providers {
		p, err := pc.defaultProvider(st, prvname)
		if err != nil {
			return nil, err
		}
		providersMap[prvname] = p
	}

	// Copy user providers
	providers, err := pc.db.GetUserProviders(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user providers: %w", err)
	}

	for _, prv := range providers {
		p, err := pc.newUserProvider(st, prv)
		if err != nil {
			// Any unbuildable saved provider (its type is disabled, or its stored
			// config is stale/invalid) is skipped, not propagated — one bad row must
			// not take down the whole list. Logged at Warn (not Error) since this is
			// an expected, recoverable condition (e.g. a type disabled via env vars),
			// but it must stay visible in production logs rather than silently
			// disappearing the provider from the list on every fetch.
			logrus.WithError(err).Warnf("skipping unusable user provider '%s' (type '%s')", prv.Name, prv.Type)
			continue
		}
		providersMap[provider.ProviderName(prv.Name)] = p
	}

	return providersMap, nil
}

func (pc *providerController) NewProvider(prv database.Provider) (provider.Provider, error) {
	return pc.newUserProvider(pc.current(), prv)
}

// defaultProvider wraps the built-in provider prvname of st in a liveProvider.
func (pc *providerController) defaultProvider(
	st *providerState,
	prvname provider.ProviderName,
) (provider.Provider, error) {
	prv, err := st.providers.Get(prvname)
	if err != nil {
		return nil, err
	}

	return pc.newLiveProvider(st, prvname, prv.Type(), func(st *providerState) (provider.Provider, error) {
		return st.providers.Get(prvname)
	})
}

// newUserProvider wraps a user-defined provider row in a liveProvider, so it
// is rebuilt with new credentials when the settings of its type change.
func (pc *providerController) newUserProvider(st *providerState, prv database.Provider) (provider.Provider, error) {
	rawConfig := prv.Config
	if len(rawConfig) == 0 {
		rawConfig = []byte(pconfig.EmptyProviderConfigRaw)
	}
	providerName := provider.ProviderName(prv.Name)
	providerType := provider.ProviderType(prv.Type)

	return pc.newLiveProvider(st, providerName, providerType, func(st *providerState) (provider.Provider, error) {
		return buildUserProvider(st, providerName, providerType, rawConfig)
	})
}

// buildUserProvider builds a user-defined provider from its stored agent
// config and the credentials/endpoints of st.
func buildUserProvider(
	st *providerState,
	prvname provider.ProviderName,
	prvtype provider.ProviderType,
	rawConfig []byte,
) (provider.Provider, error) {
	// Check if the provider type is available via check default one
	if !st.providers.ListTypes().Contains(prvtype) {
		return nil, fmt.Errorf("provider type '%s' is not available", prvtype)
	}

	e, ok := entryForType(prvtype)
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", prvtype)
	}

	config, err := e.BuildConfig(st.cfg, rawConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build %s provider config: %w", prvtype, err)
	}

	return e.New(st.cfg, prvname, config)
}

func (pc *providerController) SeedDefaultProviders(ctx context.Context, userID int64) error {
	st := pc.current()
	if st.cfg.BedrockConfig == "" {
		return nil
	}
	if !st.cfg.BedrockDefaultAuth && st.cfg.BedrockBearerToken == "" &&
		(st.cfg.BedrockAccessKey == "" || st.cfg.BedrockSecretKey == "") {
		return nil
	}

	bedrockCfg, ok := st.configs[provider.ProviderBedrock]
	if !ok {
		return nil
	}

	rawConfig, err := json.Marshal(bedrockCfg)
	if err != nil {
		return fmt.Errorf("failed to marshal bedrock config: %w", err)
	}

	prvname := bedrockCfg.Name
	if prvname == "" {
		prvname = string(provider.DefaultProviderNameBedrock)
	}
	existing, err := pc.db.GetUserProviderByName(ctx, database.GetUserProviderByNameParams{
		Name:   prvname,
		UserID: userID,
	})
	if err != nil && !errors.Is(err, sql.ErrNoRows) {
		return fmt.Errorf("failed to get provider '%s' from database: %w", prvname, err)
	}
	if err != nil {
		_, err = pc.db.CreateProvider(ctx, database.CreateProviderParams{
			UserID: userID,
			Type:   database.ProviderType(provider.ProviderBedrock),
			Name:   prvname,
			Config: rawConfig,
		})
		return err
	}

	_, err = pc.db.UpdateUserProvider(ctx, database.UpdateUserProviderParams{
		ID:     existing.ID,
		UserID: userID,
		Config: rawConfig,
		Name:   existing.Name,
	})
	return err
}

func (pc *providerController) CreateProvider(
	ctx context.Context,
	userID int64,
	prvname provider.ProviderName,
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.CreateProvider")
	defer span.End()

	var (
		err    error
		result database.Provider
	)

	if config, err = pc.patchProviderConfig(prvtype, config); err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	if err = config.Validate(); err != nil {
		return result, fmt.Errorf("invalid provider config: %w", err)
	}

	rawConfig, err := json.Marshal(config)
	if err != nil {
		return result, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	result, err = pc.db.CreateProvider(ctx, database.CreateProviderParams{
		UserID: userID,
		Type:   database.ProviderType(prvtype),
		Name:   string(prvname),
		Config: rawConfig,
	})
	if err != nil {
		return result, fmt.Errorf("failed to create provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) UpdateProvider(
	ctx context.Context,
	userID int64,
	prvID int64,
	prvname provider.ProviderName,
	config *pconfig.ProviderConfig,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.UpdateProvider")
	defer span.End()

	var (
		err    error
		result database.Provider
	)

	prv, err := pc.db.GetUserProvider(ctx, database.GetUserProviderParams{
		ID:     prvID,
		UserID: userID,
	})
	if err != nil {
		return result, fmt.Errorf("failed to get provider: %w", err)
	}
	prvtype := provider.ProviderType(prv.Type)

	if config, err = pc.patchProviderConfig(prvtype, config); err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	if err = config.Validate(); err != nil {
		return result, fmt.Errorf("invalid provider config: %w", err)
	}

	rawConfig, err := json.Marshal(config)
	if err != nil {
		return result, fmt.Errorf("failed to marshal provider config: %w", err)
	}

	result, err = pc.db.UpdateUserProvider(ctx, database.UpdateUserProviderParams{
		ID:     prvID,
		UserID: userID,
		Name:   string(prvname),
		Config: rawConfig,
	})
	if err != nil {
		return result, fmt.Errorf("failed to update provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) DeleteProvider(
	ctx context.Context,
	userID int64,
	prvID int64,
) (database.Provider, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.DeleteProvider")
	defer span.End()

	result, err := pc.db.DeleteUserProvider(ctx, database.DeleteUserProviderParams{
		ID:     prvID,
		UserID: userID,
	})
	if err != nil {
		return result, fmt.Errorf("failed to delete provider: %w", err)
	}

	return result, nil
}

func (pc *providerController) TestAgent(
	ctx context.Context,
	prvtype provider.ProviderType,
	agentType pconfig.ProviderOptionsType,
	config *pconfig.AgentConfig,
) (tester.AgentTestResults, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.TestAgent")
	defer span.End()

	var result tester.AgentTestResults

	// Create provider config with single agent configuration
	testConfig := &pconfig.ProviderConfig{}

	// Set the agent config to the appropriate field based on agent type
	switch agentType {
	case pconfig.OptionsTypeSimple:
		testConfig.Simple = config
	case pconfig.OptionsTypeSimpleJSON:
		testConfig.SimpleJSON = config
	case pconfig.OptionsTypePrimaryAgent:
		testConfig.PrimaryAgent = config
	case pconfig.OptionsTypeAssistant:
		testConfig.Assistant = config
	case pconfig.OptionsTypeGenerator:
		testConfig.Generator = config
	case pconfig.OptionsTypeRefiner:
		testConfig.Refiner = config
	case pconfig.OptionsTypeAdviser:
		testConfig.Adviser = config
	case pconfig.OptionsTypeReflector:
		testConfig.Reflector = config
	case pconfig.OptionsTypeSearcher:
		testConfig.Searcher = config
	case pconfig.OptionsTypeEnricher:
		testConfig.Enricher = config
	case pconfig.OptionsTypeCoder:
		testConfig.Coder = config
	case pconfig.OptionsTypeInstaller:
		testConfig.Installer = config
	case pconfig.OptionsTypePentester:
		testConfig.Pentester = config
	default:
		return result, fmt.Errorf("unsupported agent type: %s", agentType)
	}

	// Patch with defaults
	patchedConfig, err := pc.patchProviderConfig(prvtype, testConfig)
	if err != nil {
		return result, fmt.Errorf("failed to patch provider config: %w", err)
	}

	// Create temporary provider for testing using existing provider logic
	providerName := provider.ProviderName("test-provider")
	tempProvider, err := pc.buildProviderFromConfig(prvtype, providerName, patchedConfig)
	if err != nil {
		return result, fmt.Errorf("failed to create provider for testing: %w", err)
	}

	// Run tests for specific agent type only
	results, err := tester.TestProvider(
		ctx,
		tempProvider,
		tester.WithAgentTypes(agentType),
		tester.WithVerbose(false),
		tester.WithParallelWorkers(defaultTestParallelWorkersNumber),
	)
	if err != nil {
		return result, fmt.Errorf("failed to test agent: %w", err)
	}

	// Extract results for the specific agent type
	switch agentType {
	case pconfig.OptionsTypeSimple:
		result = results.Simple
	case pconfig.OptionsTypeSimpleJSON:
		result = results.SimpleJSON
	case pconfig.OptionsTypePrimaryAgent:
		result = results.PrimaryAgent
	case pconfig.OptionsTypeAssistant:
		result = results.Assistant
	case pconfig.OptionsTypeGenerator:
		result = results.Generator
	case pconfig.OptionsTypeRefiner:
		result = results.Refiner
	case pconfig.OptionsTypeAdviser:
		result = results.Adviser
	case pconfig.OptionsTypeReflector:
		result = results.Reflector
	case pconfig.OptionsTypeSearcher:
		result = results.Searcher
	case pconfig.OptionsTypeEnricher:
		result = results.Enricher
	case pconfig.OptionsTypeCoder:
		result = results.Coder
	case pconfig.OptionsTypeInstaller:
		result = results.Installer
	case pconfig.OptionsTypePentester:
		result = results.Pentester
	default:
		return result, fmt.Errorf("unexpected agent type: %s", agentType)
	}

	return result, nil
}

func (pc *providerController) TestProvider(
	ctx context.Context,
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (tester.ProviderTestResults, error) {
	ctx, span := obs.Observer.NewSpan(ctx, obs.SpanKindInternal, "providers.TestProvider")
	defer span.End()

	var results tester.ProviderTestResults

	// Patch config with defaults
	patchedConfig, err := pc.patchProviderConfig(prvtype, config)
	if err != nil {
		return results, fmt.Errorf("failed to patch provider config: %w", err)
	}

	// Create provider for testing
	providerName := provider.ProviderName("test-provider")
	testProvider, err := pc.buildProviderFromConfig(prvtype, providerName, patchedConfig)
	if err != nil {
		return results, fmt.Errorf("failed to create provider for testing: %w", err)
	}

	// Run full provider testing
	results, err = tester.TestProvider(
		ctx,
		testProvider,
		tester.WithVerbose(false),
		tester.WithParallelWorkers(defaultTestParallelWorkersNumber),
	)
	if err != nil {
		return results, fmt.Errorf("failed to test provider: %w", err)
	}

	return results, nil
}

func (pc *providerController) patchProviderConfig(
	prvtype provider.ProviderType,
	config *pconfig.ProviderConfig,
) (*pconfig.ProviderConfig, error) {
	st := pc.current()
	defaultCfg, ok := st.configs[prvtype]
	if !ok {
		if reason, hasReason := st.configErrors[prvtype]; hasReason {
			return nil, fmt.Errorf(
				"provider type '%s' has no default config because it failed to load: %w",
				prvtype.String(), reason,
			)
		}
		return nil, fmt.Errorf("default provider config not found for type: %s", prvtype.String())
	}

	if config == nil {
		return defaultCfg, nil
	}

	if config.Simple == nil {
		config.Simple = defaultCfg.Simple
	}
	if config.SimpleJSON == nil {
		config.SimpleJSON = defaultCfg.SimpleJSON
	}
	if config.PrimaryAgent == nil {
		config.PrimaryAgent = defaultCfg.PrimaryAgent
	}
	if config.Assistant == nil {
		config.Assistant = defaultCfg.Assistant
	}
	if config.Generator == nil {
		config.Generator = defaultCfg.Generator
	}
	if config.Refiner == nil {
		config.Refiner = defaultCfg.Refiner
	}
	if config.Adviser == nil {
		config.Adviser = defaultCfg.Adviser
	}
	if config.Reflector == nil {
		config.Reflector = defaultCfg.Reflector
	}
	if config.Searcher == nil {
		config.Searcher = defaultCfg.Searcher
	}
	if config.Enricher == nil {
		config.Enricher = defaultCfg.Enricher
	}
	if config.Coder == nil {
		config.Coder = defaultCfg.Coder
	}
	if config.Installer == nil {
		config.Installer = defaultCfg.Installer
	}
	if config.Pentester == nil {
		config.Pentester = defaultCfg.Pentester
	}

	config.SetDefaultOptions(defaultCfg.GetDefaultOptions())

	return config, nil
}

func (pc *providerController) buildProviderFromConfig(
	prvtype provider.ProviderType,
	prvname provider.ProviderName,
	config *pconfig.ProviderConfig,
) (provider.Provider, error) {
	e, ok := entryForType(prvtype)
	if !ok {
		return nil, fmt.Errorf("unknown provider type: %s", prvtype)
	}

	return e.New(pc.current().cfg, prvname, config)
}

func newAtomicInt64(seed int64) *atomic.Int64 {
	var number atomic.Int64

	if seed == 0 {
		bigID, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt64))
		if err != nil {
			return &number
		}
		seed = bigID.Int64()
	}

	number.Store(seed)
	return &number
}

// callWithSetupRetries wraps a single-shot LLM prompt call used during flow/
// assistant bootstrap (docker image, language, and title selection) with the
// same short retry-with-backoff already used for the agent execution loop
// (see performSimpleChain/callWithRetries), so one transient error from the
// LLM gateway (e.g. a bad gateway from a litellm proxy) does not fail flow or
// assistant creation outright.
func callWithSetupRetries(
	ctx context.Context,
	prv provider.Provider,
	opt pconfig.ProviderOptionsType,
	prompt string,
) (string, error) {
	var (
		result string
		err    error
	)

	for idx := 0; idx <= maxRetriesToCallSimpleChain; idx++ {
		if idx == maxRetriesToCallSimpleChain {
			return "", fmt.Errorf("failed to call llm after %d retries: %w", idx, err)
		}

		result, err = prv.Call(ctx, opt, prompt)
		if err == nil {
			return result, nil
		}

		if errors.Is(err, context.Canceled) {
			return "", err
		}

		select {
		case <-ctx.Done():
			return "", ctx.Err()
		case <-time.After(delayBetweenRetries):
		}
	}

	return "", err
}
