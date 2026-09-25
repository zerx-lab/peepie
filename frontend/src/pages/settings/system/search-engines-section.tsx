import type { Dispatch, ReactNode, SetStateAction } from 'react';

import { useTranslation } from 'react-i18next';

import type { SearchEngineSettingsFragmentFragment } from '@/graphql/types';

import { Card, CardContent, CardHeader, CardTitle } from '@/components/ui/card';

import { NumberField, SecretField, SelectField, StatusBadge, TextField, ToggleField } from './fields';
import {
    duckDuckGoRegionOptions,
    duckDuckGoSafeSearchOptions,
    duckDuckGoTimeRangeOptions,
    getSearchEngineStatus,
    type SearchEngineEdits,
    type SearchEngineId,
} from './search-engines';

// Radix Select rejects an empty-string item value, so the "use backend default" option is
// represented by this sentinel and translated to/from '' at the edit boundary.
const DEFAULT_OPTION = '__default__';
const toSelectValue = (value: string) => (value === '' ? DEFAULT_OPTION : value);
const fromSelectValue = (value: string) => (value === DEFAULT_OPTION ? '' : value);

export const searchEngineAnchorId = (id: SearchEngineId) => `search-engine-${id}`;

interface SearchEnginesSectionProps {
    edits: SearchEngineEdits;
    onEditsChange: Dispatch<SetStateAction<SearchEngineEdits>>;
    server: SearchEngineSettingsFragmentFragment;
}

export function SearchEnginesSection({ edits, onEditsChange, server }: SearchEnginesSectionProps) {
    const { t } = useTranslation('settings');
    const form = { ...server, ...edits };

    const set = <K extends keyof SearchEngineEdits>(key: K, value: SearchEngineEdits[K]) =>
        onEditsChange((prev) => ({ ...prev, [key]: value }));

    return (
        <>
            <EngineCard
                engineId="duckduckgo"
                server={server}
            >
                <>
                    <ToggleField
                        checked={form.duckduckgoEnabled}
                        label={t('system.searchEngines.enabled')}
                        onChange={(v) => set('duckduckgoEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                        <SelectField
                            label={t('system.searchEngines.region')}
                            onChange={(v) => set('duckduckgoRegion', fromSelectValue(v))}
                            options={[
                                { label: t('system.searchEngines.regionDefault'), value: DEFAULT_OPTION },
                                ...duckDuckGoRegionOptions.map((v) => ({
                                    label: t(`system.searchEngines.regionOptions.${v}`),
                                    value: v,
                                })),
                            ]}
                            value={toSelectValue(form.duckduckgoRegion)}
                        />
                        <SelectField
                            label={t('system.searchEngines.safeSearch')}
                            onChange={(v) => set('duckduckgoSafesearch', fromSelectValue(v))}
                            options={[
                                { label: t('system.searchEngines.safeSearchDefault'), value: DEFAULT_OPTION },
                                ...duckDuckGoSafeSearchOptions.map((v) => ({
                                    label: t(`system.searchEngines.safeSearchOptions.${v}`),
                                    value: v,
                                })),
                            ]}
                            value={toSelectValue(form.duckduckgoSafesearch)}
                        />
                        <SelectField
                            label={t('system.searchEngines.timeRange')}
                            onChange={(v) => set('duckduckgoTimeRange', fromSelectValue(v))}
                            options={[
                                { label: t('system.searchEngines.timeRangeDefault'), value: DEFAULT_OPTION },
                                ...duckDuckGoTimeRangeOptions.map((v) => ({
                                    label: t(`system.searchEngines.timeRangeOptions.${v}`),
                                    value: v,
                                })),
                            ]}
                            value={toSelectValue(form.duckduckgoTimeRange)}
                        />
                    </div>
                </>
            </EngineCard>

            <EngineCard
                engineId="sploitus"
                server={server}
            >
                <ToggleField
                    checked={form.sploitusEnabled}
                    label={t('system.searchEngines.enabled')}
                    onChange={(v) => set('sploitusEnabled', v)}
                />
            </EngineCard>

            <EngineCard
                engineId="brave"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <SecretField
                        edit={edits.braveApiKey}
                        isSet={server.braveApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('braveApiKey', v)}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="google"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                    <SecretField
                        edit={edits.googleApiKey}
                        isSet={server.googleApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('googleApiKey', v)}
                    />
                    <TextField
                        label={t('system.searchEngines.googleCxKey')}
                        onChange={(v) => set('googleCxKey', v)}
                        value={form.googleCxKey}
                    />
                    <TextField
                        label={t('system.searchEngines.googleLrKey')}
                        onChange={(v) => set('googleLrKey', v)}
                        value={form.googleLrKey}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="traversaal"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <SecretField
                        edit={edits.traversaalApiKey}
                        isSet={server.traversaalApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('traversaalApiKey', v)}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="tavily"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <SecretField
                        edit={edits.tavilyApiKey}
                        isSet={server.tavilyApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('tavilyApiKey', v)}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="firecrawl"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <SecretField
                        edit={edits.firecrawlApiKey}
                        isSet={server.firecrawlApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('firecrawlApiKey', v)}
                    />
                    <TextField
                        label={t('system.searchEngines.apiUrl')}
                        onChange={(v) => set('firecrawlApiUrl', v)}
                        value={form.firecrawlApiUrl}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="perplexity"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                    <SecretField
                        edit={edits.perplexityApiKey}
                        isSet={server.perplexityApiKeySet}
                        label={t('system.searchEngines.apiKey')}
                        onChange={(v) => set('perplexityApiKey', v)}
                    />
                    <TextField
                        label={t('system.searchEngines.model')}
                        onChange={(v) => set('perplexityModel', v)}
                        value={form.perplexityModel}
                    />
                    <TextField
                        label={t('system.searchEngines.contextSize')}
                        onChange={(v) => set('perplexityContextSize', v)}
                        value={form.perplexityContextSize}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="searxng"
                server={server}
            >
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <TextField
                        label={t('system.searchEngines.apiUrl')}
                        onChange={(v) => set('searxngUrl', v)}
                        value={form.searxngUrl}
                    />
                    <TextField
                        label={t('system.searchEngines.categories')}
                        onChange={(v) => set('searxngCategories', v)}
                        value={form.searxngCategories}
                    />
                    <TextField
                        label={t('system.searchEngines.language')}
                        onChange={(v) => set('searxngLanguage', v)}
                        value={form.searxngLanguage}
                    />
                    <TextField
                        label={t('system.searchEngines.safeSearch')}
                        onChange={(v) => set('searxngSafesearch', v)}
                        value={form.searxngSafesearch}
                    />
                    <TextField
                        label={t('system.searchEngines.timeRange')}
                        onChange={(v) => set('searxngTimeRange', v)}
                        value={form.searxngTimeRange}
                    />
                    <NumberField
                        label={t('system.searchEngines.timeoutSeconds')}
                        onChange={(v) => set('searxngTimeout', v)}
                        value={form.searxngTimeout}
                    />
                </div>
            </EngineCard>

            <EngineCard
                engineId="webSearchInternal"
                server={server}
            >
                <>
                    <ToggleField
                        checked={form.webSearchInternalEnabled}
                        label={t('system.searchEngines.enabled')}
                        onChange={(v) => set('webSearchInternalEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <NumberField
                            label={t('system.searchEngines.maxSites')}
                            onChange={(v) => set('webSearchInternalMaxSites', v)}
                            value={form.webSearchInternalMaxSites}
                        />
                        <NumberField
                            label={t('system.searchEngines.maxSiteBytes')}
                            onChange={(v) => set('webSearchInternalMaxSiteBytes', v)}
                            value={form.webSearchInternalMaxSiteBytes}
                        />
                    </div>
                </>
            </EngineCard>
        </>
    );
}

function EngineCard({
    children,
    engineId,
    server,
}: {
    children: ReactNode;
    engineId: SearchEngineId;
    server: SearchEngineSettingsFragmentFragment;
}) {
    const { t } = useTranslation('settings');

    return (
        <Card
            className="scroll-mt-16"
            id={searchEngineAnchorId(engineId)}
        >
            <CardHeader>
                <div className="flex flex-wrap items-center gap-2">
                    <CardTitle>{t(`system.searchEngines.engines.${engineId}`)}</CardTitle>
                    <StatusBadge status={getSearchEngineStatus(server, engineId)} />
                </div>
            </CardHeader>
            <CardContent className="flex flex-col gap-3">{children}</CardContent>
        </Card>
    );
}
