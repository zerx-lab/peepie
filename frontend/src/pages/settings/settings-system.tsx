import { useMutation, useQuery } from '@apollo/client/react';
import { Loader2, Save } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';

import type {
    ExecutionSettingsInput,
    SearchEngineSettingsFragmentFragment,
    SearchEngineSettingsInput,
} from '@/graphql/types';

import { AppHeader, AppHeaderContent, AppHeaderTitle } from '@/components/layouts/app/app-header';
import { ErrorState } from '@/components/shared/error-state';
import { LoadingState } from '@/components/shared/loading-state';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { Input } from '@/components/ui/input';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import {
    SettingsExecutionDocument,
    SettingsSearchEnginesDocument,
    UpdateExecutionSettingsDocument,
    UpdateSearchEngineSettingsDocument,
} from '@/graphql/types';

// Local editable shape: secret fields are always sent as `undefined` unless the
// user types a new value, so an unmodified save never clears an already
// configured key (see backend applySecretSetting). `*Set` mirrors what the
// server reports as already configured, purely for the placeholder/badge.
type SearchForm = Omit<
    SearchEngineSettingsFragmentFragment,
    'firecrawlApiKeySet' | 'googleApiKeySet' | 'perplexityApiKeySet' | 'tavilyApiKeySet' | 'traversaalApiKeySet'
> & {
    firecrawlApiKey: string;
    firecrawlApiKeySet: boolean;
    googleApiKey: string;
    googleApiKeySet: boolean;
    perplexityApiKey: string;
    perplexityApiKeySet: boolean;
    tavilyApiKey: string;
    tavilyApiKeySet: boolean;
    traversaalApiKey: string;
    traversaalApiKeySet: boolean;
};

function ExecutionCard() {
    const { t } = useTranslation('settings');
    const { data, error, loading, refetch } = useQuery(SettingsExecutionDocument);
    const [updateSettings, { loading: isSaving }] = useMutation(UpdateExecutionSettingsDocument);
    // Local edits layered on top of the last known server values. Reset to {}
    // after a successful save, since the mutation result becomes the new
    // server value for the active query (Apollo re-renders `data`).
    const [overrides, setOverrides] = useState<Partial<ExecutionSettingsInput>>({});

    if (loading && !data) {
        return <LoadingState title={t('system.execution.loading')} />;
    }

    if (error) {
        return (
            <ErrorState
                message={error.message}
                onRetry={refetch}
                title={t('system.execution.errorTitle')}
            />
        );
    }

    if (!data?.settingsExecution) {
        return null;
    }

    const form: ExecutionSettingsInput = { ...data.settingsExecution, ...overrides };

    const set = <K extends keyof ExecutionSettingsInput>(key: K, value: ExecutionSettingsInput[K]) =>
        setOverrides((prev) => ({ ...prev, [key]: value }));

    const handleSave = async () => {
        try {
            await updateSettings({ variables: { input: form } });
            setOverrides({});
            toast.success(t('system.saved'));
        } catch (err) {
            toast.error(err instanceof Error ? err.message : t('system.saveFailed'));
        }
    };

    return (
        <Card>
            <CardHeader>
                <CardTitle>{t('system.execution.title')}</CardTitle>
                <CardDescription>{t('system.execution.description')}</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-6">
                <ToggleField
                    checked={form.assistantUseAgents}
                    label={t('system.execution.assistantUseAgents')}
                    onChange={(v) => set('assistantUseAgents', v)}
                />
                <ToggleField
                    checked={form.agentPlanningStepEnabled}
                    label={t('system.execution.agentPlanningStepEnabled')}
                    onChange={(v) => set('agentPlanningStepEnabled', v)}
                />

                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <NumberField
                        label={t('system.execution.maxGeneralAgentToolCalls')}
                        onChange={(v) => set('maxGeneralAgentToolCalls', v)}
                        value={form.maxGeneralAgentToolCalls}
                    />
                    <NumberField
                        label={t('system.execution.maxLimitedAgentToolCalls')}
                        onChange={(v) => set('maxLimitedAgentToolCalls', v)}
                        value={form.maxLimitedAgentToolCalls}
                    />
                </div>

                <div className="flex flex-col gap-3">
                    <ToggleField
                        checked={form.executionMonitorEnabled}
                        label={t('system.execution.monitorEnabled')}
                        onChange={(v) => set('executionMonitorEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <NumberField
                            label={t('system.execution.sameToolLimit')}
                            onChange={(v) => set('executionMonitorSameToolLimit', v)}
                            value={form.executionMonitorSameToolLimit}
                        />
                        <NumberField
                            label={t('system.execution.totalToolLimit')}
                            onChange={(v) => set('executionMonitorTotalToolLimit', v)}
                            value={form.executionMonitorTotalToolLimit}
                        />
                    </div>
                </div>

                <div className="flex justify-end">
                    <Button
                        disabled={isSaving}
                        onClick={handleSave}
                    >
                        {isSaving ? <Loader2 className="animate-spin" /> : <Save />}
                        {t('system.save')}
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}

function NumberField({ label, onChange, value }: { label: string; onChange: (value: number) => void; value: number }) {
    return (
        <div className="flex flex-col gap-1.5">
            <Label>{label}</Label>
            <Input
                onChange={(e) => onChange(Number(e.target.value) || 0)}
                type="number"
                value={value}
            />
        </div>
    );
}

function SearchEnginesCard() {
    const { t } = useTranslation('settings');
    const { data, error, loading, refetch } = useQuery(SettingsSearchEnginesDocument);
    const [updateSettings, { loading: isSaving }] = useMutation(UpdateSearchEngineSettingsDocument);
    // Local edits layered on top of the last known server values. Reset to {}
    // after a successful save, since the mutation result becomes the new
    // server value for the active query (Apollo re-renders `data`).
    const [overrides, setOverrides] = useState<Partial<SearchForm>>({});

    if (loading && !data) {
        return <LoadingState title={t('system.searchEngines.loading')} />;
    }

    if (error) {
        return (
            <ErrorState
                message={error.message}
                onRetry={refetch}
                title={t('system.searchEngines.errorTitle')}
            />
        );
    }

    if (!data?.settingsSearchEngines) {
        return null;
    }

    const form: SearchForm = { ...toSearchForm(data.settingsSearchEngines), ...overrides };

    const set = <K extends keyof SearchForm>(key: K, value: SearchForm[K]) =>
        setOverrides((prev) => ({ ...prev, [key]: value }));

    const handleSave = async () => {
        const input: SearchEngineSettingsInput = {
            duckduckgoEnabled: form.duckduckgoEnabled,
            duckduckgoRegion: form.duckduckgoRegion,
            duckduckgoSafesearch: form.duckduckgoSafesearch,
            duckduckgoTimeRange: form.duckduckgoTimeRange,
            firecrawlApiKey: form.firecrawlApiKey || undefined,
            firecrawlApiUrl: form.firecrawlApiUrl,
            googleApiKey: form.googleApiKey || undefined,
            googleCxKey: form.googleCxKey,
            googleLrKey: form.googleLrKey,
            perplexityApiKey: form.perplexityApiKey || undefined,
            perplexityContextSize: form.perplexityContextSize,
            perplexityModel: form.perplexityModel,
            searxngCategories: form.searxngCategories,
            searxngLanguage: form.searxngLanguage,
            searxngSafesearch: form.searxngSafesearch,
            searxngTimeout: form.searxngTimeout,
            searxngTimeRange: form.searxngTimeRange,
            searxngUrl: form.searxngUrl,
            sploitusEnabled: form.sploitusEnabled,
            tavilyApiKey: form.tavilyApiKey || undefined,
            traversaalApiKey: form.traversaalApiKey || undefined,
            webSearchInternalEnabled: form.webSearchInternalEnabled,
            webSearchInternalMaxSiteBytes: form.webSearchInternalMaxSiteBytes,
            webSearchInternalMaxSites: form.webSearchInternalMaxSites,
        };

        try {
            await updateSettings({ variables: { input } });
            setOverrides({});
            toast.success(t('system.saved'));
        } catch (err) {
            toast.error(err instanceof Error ? err.message : t('system.saveFailed'));
        }
    };

    const secretPlaceholder = (isSet: boolean) =>
        isSet ? t('system.searchEngines.secretConfigured') : t('system.searchEngines.secretEmpty');

    return (
        <Card>
            <CardHeader>
                <CardTitle>{t('system.searchEngines.title')}</CardTitle>
                <CardDescription>{t('system.searchEngines.description')}</CardDescription>
            </CardHeader>
            <CardContent className="flex flex-col gap-6">
                <div className="flex flex-col gap-3">
                    <div className="flex items-center gap-2">
                        <h4 className="text-sm font-medium">{t('system.searchEngines.duckduckgo')}</h4>
                        {form.duckduckgoEnabled && <Badge variant="secondary">{t('system.enabled')}</Badge>}
                    </div>
                    <ToggleField
                        checked={form.duckduckgoEnabled}
                        label={t('system.searchEngines.enabled')}
                        onChange={(v) => set('duckduckgoEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                        <TextField
                            label={t('system.searchEngines.region')}
                            onChange={(v) => set('duckduckgoRegion', v)}
                            value={form.duckduckgoRegion}
                        />
                        <TextField
                            label={t('system.searchEngines.safeSearch')}
                            onChange={(v) => set('duckduckgoSafesearch', v)}
                            value={form.duckduckgoSafesearch}
                        />
                        <TextField
                            label={t('system.searchEngines.timeRange')}
                            onChange={(v) => set('duckduckgoTimeRange', v)}
                            value={form.duckduckgoTimeRange}
                        />
                    </div>
                </div>

                <ToggleField
                    checked={form.sploitusEnabled}
                    label={t('system.searchEngines.sploitus')}
                    onChange={(v) => set('sploitusEnabled', v)}
                />

                <div className="flex flex-col gap-3">
                    <h4 className="text-sm font-medium">{t('system.searchEngines.google')}</h4>
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                        <SecretField
                            label={t('system.searchEngines.apiKey')}
                            onChange={(v) => set('googleApiKey', v)}
                            placeholder={secretPlaceholder(form.googleApiKeySet)}
                            value={form.googleApiKey}
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
                </div>

                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <SecretField
                        label={t('system.searchEngines.traversaal')}
                        onChange={(v) => set('traversaalApiKey', v)}
                        placeholder={secretPlaceholder(form.traversaalApiKeySet)}
                        value={form.traversaalApiKey}
                    />
                    <SecretField
                        label={t('system.searchEngines.tavily')}
                        onChange={(v) => set('tavilyApiKey', v)}
                        placeholder={secretPlaceholder(form.tavilyApiKeySet)}
                        value={form.tavilyApiKey}
                    />
                </div>

                <div className="flex flex-col gap-3">
                    <h4 className="text-sm font-medium">{t('system.searchEngines.firecrawl')}</h4>
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <SecretField
                            label={t('system.searchEngines.apiKey')}
                            onChange={(v) => set('firecrawlApiKey', v)}
                            placeholder={secretPlaceholder(form.firecrawlApiKeySet)}
                            value={form.firecrawlApiKey}
                        />
                        <TextField
                            label={t('system.searchEngines.apiUrl')}
                            onChange={(v) => set('firecrawlApiUrl', v)}
                            value={form.firecrawlApiUrl}
                        />
                    </div>
                </div>

                <div className="flex flex-col gap-3">
                    <h4 className="text-sm font-medium">{t('system.searchEngines.perplexity')}</h4>
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-3">
                        <SecretField
                            label={t('system.searchEngines.apiKey')}
                            onChange={(v) => set('perplexityApiKey', v)}
                            placeholder={secretPlaceholder(form.perplexityApiKeySet)}
                            value={form.perplexityApiKey}
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
                </div>

                <div className="flex flex-col gap-3">
                    <h4 className="text-sm font-medium">{t('system.searchEngines.searxng')}</h4>
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
                </div>

                <div className="flex flex-col gap-3">
                    <h4 className="text-sm font-medium">{t('system.searchEngines.webSearchInternal')}</h4>
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
                </div>

                <div className="flex justify-end">
                    <Button
                        disabled={isSaving}
                        onClick={handleSave}
                    >
                        {isSaving ? <Loader2 className="animate-spin" /> : <Save />}
                        {t('system.save')}
                    </Button>
                </div>
            </CardContent>
        </Card>
    );
}

function SecretField({
    label,
    onChange,
    placeholder,
    value,
}: {
    label: string;
    onChange: (value: string) => void;
    placeholder: string;
    value: string;
}) {
    return (
        <div className="flex flex-col gap-1.5">
            <Label>{label}</Label>
            <Input
                autoComplete="off"
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                type="password"
                value={value}
            />
        </div>
    );
}

function SettingsSystem() {
    return (
        <>
            <SettingsSystemHeader />
            <div className="flex flex-1 flex-col gap-6 p-4">
                <SearchEnginesCard />
                <ExecutionCard />
            </div>
        </>
    );
}

function SettingsSystemHeader() {
    const { t } = useTranslation('settings');

    return (
        <AppHeader>
            <AppHeaderContent>
                <AppHeaderTitle>{t('system.title')}</AppHeaderTitle>
            </AppHeaderContent>
        </AppHeader>
    );
}

function TextField({ label, onChange, value }: { label: string; onChange: (value: string) => void; value: string }) {
    return (
        <div className="flex flex-col gap-1.5">
            <Label>{label}</Label>
            <Input
                onChange={(e) => onChange(e.target.value)}
                value={value}
            />
        </div>
    );
}

function ToggleField({
    checked,
    label,
    onChange,
}: {
    checked: boolean;
    label: string;
    onChange: (value: boolean) => void;
}) {
    return (
        <div className="flex items-center justify-between gap-4 rounded-md border p-3">
            <Label className="font-normal">{label}</Label>
            <Switch
                checked={checked}
                onCheckedChange={onChange}
            />
        </div>
    );
}

function toSearchForm(data: SearchEngineSettingsFragmentFragment): SearchForm {
    return {
        ...data,
        firecrawlApiKey: '',
        googleApiKey: '',
        perplexityApiKey: '',
        tavilyApiKey: '',
        traversaalApiKey: '',
    };
}

export default SettingsSystem;
