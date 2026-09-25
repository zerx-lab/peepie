import type { SearchEngineSettingsFragmentFragment, SearchEngineSettingsInput } from '@/graphql/types';

import { type EditsWithSecrets, hasPendingEdits, resolveSecret } from './secrets';

export const searchEngineIds = [
    'duckduckgo',
    'sploitus',
    'brave',
    'google',
    'traversaal',
    'tavily',
    'firecrawl',
    'perplexity',
    'searxng',
    'webSearchInternal',
] as const;

export type SearchEngineId = (typeof searchEngineIds)[number];

// Mirrors backend/pkg/tools/searchers/duckduckgo.go's Region*/DuckDuckGoSafeSearch*/TimeRange*
// constants. An empty string means "use the backend default" for each field; keep these in
// sync if the backend adds/removes a supported value.
export const duckDuckGoRegionOptions = [
    'us-en',
    'uk-en',
    'de-de',
    'fr-fr',
    'es-es',
    'it-it',
    'jp-jp',
    'cn-zh',
    'ru-ru',
] as const;
export const duckDuckGoSafeSearchOptions = ['strict', 'moderate', 'off'] as const;
export const duckDuckGoTimeRangeOptions = ['d', 'w', 'm', 'y'] as const;

const searchSecretFields = [
    'braveApiKey',
    'firecrawlApiKey',
    'googleApiKey',
    'perplexityApiKey',
    'tavilyApiKey',
    'traversaalApiKey',
] as const satisfies readonly (keyof SearchEngineSettingsInput)[];

export type SearchEngineEdits = EditsWithSecrets<SearchEngineSettingsInput, (typeof searchSecretFields)[number]>;

// `disabled` = switched off by a toggle; `unconfigured` = missing credentials/endpoint.
export type SearchEngineStatus = 'disabled' | 'enabled' | 'unconfigured';

export const isSearchEnginesDirty = (server: SearchEngineSettingsFragmentFragment, edits: SearchEngineEdits): boolean =>
    hasPendingEdits(server, edits, searchSecretFields);

export function buildSearchEngineInput(
    server: SearchEngineSettingsFragmentFragment,
    edits: SearchEngineEdits,
): SearchEngineSettingsInput {
    const {
        braveApiKey,
        firecrawlApiKey,
        googleApiKey,
        perplexityApiKey,
        tavilyApiKey,
        traversaalApiKey,
        ...plainEdits
    } = edits;

    return {
        braveApiKey: resolveSecret(braveApiKey),
        duckduckgoEnabled: plainEdits.duckduckgoEnabled ?? server.duckduckgoEnabled,
        duckduckgoRegion: plainEdits.duckduckgoRegion ?? server.duckduckgoRegion,
        duckduckgoSafesearch: plainEdits.duckduckgoSafesearch ?? server.duckduckgoSafesearch,
        duckduckgoTimeRange: plainEdits.duckduckgoTimeRange ?? server.duckduckgoTimeRange,
        firecrawlApiKey: resolveSecret(firecrawlApiKey),
        firecrawlApiUrl: plainEdits.firecrawlApiUrl ?? server.firecrawlApiUrl,
        googleApiKey: resolveSecret(googleApiKey),
        googleCxKey: plainEdits.googleCxKey ?? server.googleCxKey,
        googleLrKey: plainEdits.googleLrKey ?? server.googleLrKey,
        perplexityApiKey: resolveSecret(perplexityApiKey),
        perplexityContextSize: plainEdits.perplexityContextSize ?? server.perplexityContextSize,
        perplexityModel: plainEdits.perplexityModel ?? server.perplexityModel,
        searxngCategories: plainEdits.searxngCategories ?? server.searxngCategories,
        searxngLanguage: plainEdits.searxngLanguage ?? server.searxngLanguage,
        searxngSafesearch: plainEdits.searxngSafesearch ?? server.searxngSafesearch,
        searxngTimeout: plainEdits.searxngTimeout ?? server.searxngTimeout,
        searxngTimeRange: plainEdits.searxngTimeRange ?? server.searxngTimeRange,
        searxngUrl: plainEdits.searxngUrl ?? server.searxngUrl,
        sploitusEnabled: plainEdits.sploitusEnabled ?? server.sploitusEnabled,
        tavilyApiKey: resolveSecret(tavilyApiKey),
        traversaalApiKey: resolveSecret(traversaalApiKey),
        webSearchInternalEnabled: plainEdits.webSearchInternalEnabled ?? server.webSearchInternalEnabled,
        webSearchInternalMaxSiteBytes: plainEdits.webSearchInternalMaxSiteBytes ?? server.webSearchInternalMaxSiteBytes,
        webSearchInternalMaxSites: plainEdits.webSearchInternalMaxSites ?? server.webSearchInternalMaxSites,
    };
}

// Mirrors the backend searchers' IsAvailable checks for the saved settings.
export function getSearchEngineStatus(
    server: SearchEngineSettingsFragmentFragment,
    id: SearchEngineId,
): SearchEngineStatus {
    switch (id) {
        case 'brave':
            return server.braveApiKeySet ? 'enabled' : 'unconfigured';
        case 'duckduckgo':
            return server.duckduckgoEnabled ? 'enabled' : 'disabled';
        case 'firecrawl':
            return server.firecrawlApiKeySet ? 'enabled' : 'unconfigured';
        case 'google':
            return server.googleApiKeySet && server.googleCxKey !== '' ? 'enabled' : 'unconfigured';
        case 'perplexity':
            return server.perplexityApiKeySet ? 'enabled' : 'unconfigured';
        case 'searxng':
            return server.searxngUrl.trim() !== '' ? 'enabled' : 'unconfigured';
        case 'sploitus':
            return server.sploitusEnabled ? 'enabled' : 'disabled';
        case 'tavily':
            return server.tavilyApiKeySet ? 'enabled' : 'unconfigured';
        case 'traversaal':
            return server.traversaalApiKeySet ? 'enabled' : 'unconfigured';
        case 'webSearchInternal':
            return server.webSearchInternalEnabled ? 'enabled' : 'disabled';
    }
}
