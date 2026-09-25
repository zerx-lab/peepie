import type { ReactNode } from 'react';

import { useMutation, useQuery } from '@apollo/client/react';
import { Bot, FileText, Gauge, Search, TriangleAlert } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useSearchParams } from 'react-router-dom';
import { toast } from 'sonner';

import { AppHeader, AppHeaderContent, AppHeaderTitle } from '@/components/layouts/app/app-header';
import { ErrorState } from '@/components/shared/error-state';
import { LoadingState } from '@/components/shared/loading-state';
import { Alert, AlertDescription } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import {
    SettingsEnvFileDocument,
    SettingsExecutionDocument,
    SettingsLlmProvidersDocument,
    SettingsProvidersDocument,
    SettingsSearchEnginesDocument,
    UpdateExecutionSettingsDocument,
    UpdateLlmProviderSettingsDocument,
    UpdateSearchEngineSettingsDocument,
} from '@/graphql/types';
import { type ExecutionEdits, ExecutionSection } from '@/pages/settings/system/execution-section';
import { SectionSaveBar } from '@/pages/settings/system/fields';
import {
    buildLlmProviderInput,
    getLlmProviderStatus,
    isLlmProviderDirty,
    type LlmProviderEdits,
    llmProviderIds,
} from '@/pages/settings/system/llm-providers';
import {
    getLlmProviderLabel,
    llmProviderAnchorId,
    LlmProvidersSection,
} from '@/pages/settings/system/llm-providers-section';
import {
    buildSearchEngineInput,
    getSearchEngineStatus,
    isSearchEnginesDirty,
    type SearchEngineEdits,
    searchEngineIds,
} from '@/pages/settings/system/search-engines';
import { searchEngineAnchorId, SearchEnginesSection } from '@/pages/settings/system/search-engines-section';
import { hasPendingEdits } from '@/pages/settings/system/secrets';
import {
    type SystemNavCategory,
    type SystemSection,
    systemSections,
    SystemSettingsNav,
} from '@/pages/settings/system/system-settings-nav';

function EnvFileBanner() {
    const { t } = useTranslation('settings');
    const { data } = useQuery(SettingsEnvFileDocument);
    const envFile = data?.settingsEnvFile;

    if (!envFile) {
        return null;
    }

    if (envFile.writable) {
        return (
            <Alert>
                <FileText />
                <AlertDescription>{t('system.envFile.writable', { path: envFile.path })}</AlertDescription>
            </Alert>
        );
    }

    return (
        <Alert className="border-yellow-500/50 text-yellow-800 dark:text-yellow-400 [&>svg]:text-yellow-600 dark:[&>svg]:text-yellow-400">
            <TriangleAlert />
            <AlertDescription>{t('system.envFile.notWritable', { path: envFile.path })}</AlertDescription>
        </Alert>
    );
}

function SectionHeader({ children, title }: { children?: ReactNode; title: string }) {
    return (
        <div className="flex flex-col gap-1">
            <h2 className="text-lg font-semibold">{title}</h2>
            {children}
        </div>
    );
}

function SettingsSystem() {
    const { t } = useTranslation('settings');
    const [searchParams, setSearchParams] = useSearchParams();
    const sectionParam = searchParams.get('section');
    const section: SystemSection = systemSections.find((value) => value === sectionParam) ?? 'llm-providers';

    const llmQuery = useQuery(SettingsLlmProvidersDocument);
    const searchQuery = useQuery(SettingsSearchEnginesDocument);
    const executionQuery = useQuery(SettingsExecutionDocument);

    // Pending edits live here (not in the sections) so switching categories keeps them.
    const [llmEdits, setLlmEdits] = useState<LlmProviderEdits>({});
    const [searchEdits, setSearchEdits] = useState<SearchEngineEdits>({});
    const [executionEdits, setExecutionEdits] = useState<ExecutionEdits>({});

    // Settings types carry no `id`, so Apollo cannot normalize mutation results
    // into the queries: write them back explicitly.
    const [updateLlmProviders, { loading: isSavingLlm }] = useMutation(UpdateLlmProviderSettingsDocument, {
        // The Providers page lists which provider types are enabled.
        refetchQueries: [{ query: SettingsProvidersDocument }],
        update: (cache, { data }) => {
            if (data) {
                cache.writeQuery({
                    data: { settingsLLMProviders: data.updateLLMProviderSettings },
                    query: SettingsLlmProvidersDocument,
                });
            }
        },
    });
    const [updateSearchEngines, { loading: isSavingSearch }] = useMutation(UpdateSearchEngineSettingsDocument, {
        update: (cache, { data }) => {
            if (data) {
                cache.writeQuery({
                    data: { settingsSearchEngines: data.updateSearchEngineSettings },
                    query: SettingsSearchEnginesDocument,
                });
            }
        },
    });
    const [updateExecution, { loading: isSavingExecution }] = useMutation(UpdateExecutionSettingsDocument, {
        update: (cache, { data }) => {
            if (data) {
                cache.writeQuery({
                    data: { settingsExecution: data.updateExecutionSettings },
                    query: SettingsExecutionDocument,
                });
            }
        },
    });

    const llm = llmQuery.data?.settingsLLMProviders;
    const search = searchQuery.data?.settingsSearchEngines;
    const execution = executionQuery.data?.settingsExecution;

    const isLlmDirty = !!llm && llmProviderIds.some((id) => isLlmProviderDirty(llm, llmEdits, id));
    const isSearchDirty = !!search && isSearchEnginesDirty(search, searchEdits);
    const isExecutionDirty = !!execution && hasPendingEdits(execution, executionEdits);

    const llmItems = llmProviderIds.map((id) => ({
        anchorId: llmProviderAnchorId(id),
        id,
        label: getLlmProviderLabel(id, t),
        status: llm ? getLlmProviderStatus(llm, id) : ('unconfigured' as const),
    }));
    const llmActiveCount = llmItems.filter((item) => item.status === 'active').length;
    const llmFailedCount = llmItems.filter((item) => item.status === 'error').length;

    const searchItems = searchEngineIds.map((id) => ({
        anchorId: searchEngineAnchorId(id),
        id,
        label: t(`system.searchEngines.engines.${id}`),
        status: search ? getSearchEngineStatus(search, id) : ('unconfigured' as const),
    }));
    const searchEnabledCount = searchItems.filter((item) => item.status === 'enabled').length;

    const categories: SystemNavCategory[] = [
        {
            hasUnsavedChanges: isLlmDirty,
            icon: Bot,
            id: 'llm-providers',
            items: llm ? llmItems : [],
            label: t('system.llm.title'),
            summary: llm ? t('system.llm.navSummary', { active: llmActiveCount, total: llmItems.length }) : null,
        },
        {
            hasUnsavedChanges: isSearchDirty,
            icon: Search,
            id: 'search-engines',
            items: search ? searchItems : [],
            label: t('system.searchEngines.title'),
            summary: search ? t('system.searchEngines.navSummary', { count: searchEnabledCount }) : null,
        },
        {
            hasUnsavedChanges: isExecutionDirty,
            icon: Gauge,
            id: 'execution',
            items: [],
            label: t('system.execution.title'),
            summary: execution
                ? t(execution.executionMonitorEnabled ? 'system.execution.monitorOn' : 'system.execution.monitorOff')
                : null,
        },
    ];

    const selectSection = (next: SystemSection) =>
        setSearchParams((prev) => {
            const params = new URLSearchParams(prev);
            params.set('section', next);

            return params;
        });

    const save = async (mutate: () => Promise<unknown>, reset: () => void) => {
        try {
            await mutate();
            reset();
            toast.success(t('system.saved'));
        } catch (err) {
            // Keep the pending edits so the user can fix and retry.
            toast.error(err instanceof Error ? err.message : t('system.saveFailed'));
        }
    };

    const renderContent = () => {
        switch (section) {
            case 'execution': {
                if (executionQuery.loading && !execution) {
                    return <LoadingState title={t('system.execution.loading')} />;
                }

                if (!execution) {
                    return (
                        <ErrorState
                            message={executionQuery.error?.message ?? ''}
                            onRetry={executionQuery.refetch}
                            title={t('system.execution.errorTitle')}
                        />
                    );
                }

                const form = { ...execution, ...executionEdits };

                return (
                    <>
                        <SectionHeader title={t('system.execution.title')}>
                            <div className="flex flex-wrap gap-2">
                                <Badge variant={execution.executionMonitorEnabled ? 'green' : 'outline'}>
                                    {t(
                                        execution.executionMonitorEnabled
                                            ? 'system.execution.monitorOn'
                                            : 'system.execution.monitorOff',
                                    )}
                                </Badge>
                                <Badge variant={execution.agentPlanningStepEnabled ? 'green' : 'outline'}>
                                    {t(
                                        execution.agentPlanningStepEnabled
                                            ? 'system.execution.planningOn'
                                            : 'system.execution.planningOff',
                                    )}
                                </Badge>
                            </div>
                        </SectionHeader>
                        <ExecutionSection
                            edits={executionEdits}
                            onEditsChange={setExecutionEdits}
                            server={execution}
                        />
                        {isExecutionDirty && (
                            <SectionSaveBar
                                isSaving={isSavingExecution}
                                onDiscard={() => setExecutionEdits({})}
                                onSave={() =>
                                    save(
                                        () =>
                                            updateExecution({
                                                variables: {
                                                    input: {
                                                        agentPlanningStepEnabled: form.agentPlanningStepEnabled,
                                                        assistantDockerImage: form.assistantDockerImage,
                                                        assistantUseAgents: form.assistantUseAgents,
                                                        executionMonitorEnabled: form.executionMonitorEnabled,
                                                        executionMonitorSameToolLimit:
                                                            form.executionMonitorSameToolLimit,
                                                        executionMonitorTotalToolLimit:
                                                            form.executionMonitorTotalToolLimit,
                                                        flowDockerImage: form.flowDockerImage,
                                                        maxGeneralAgentToolCalls: form.maxGeneralAgentToolCalls,
                                                        maxLimitedAgentToolCalls: form.maxLimitedAgentToolCalls,
                                                    },
                                                },
                                            }),
                                        () => setExecutionEdits({}),
                                    )
                                }
                            />
                        )}
                    </>
                );
            }

            case 'llm-providers': {
                if (llmQuery.loading && !llm) {
                    return <LoadingState title={t('system.llm.loading')} />;
                }

                if (!llm) {
                    return (
                        <ErrorState
                            message={llmQuery.error?.message ?? ''}
                            onRetry={llmQuery.refetch}
                            title={t('system.llm.errorTitle')}
                        />
                    );
                }

                return (
                    <>
                        <SectionHeader title={t('system.llm.title')}>
                            <p className="text-muted-foreground text-sm">{t('system.llm.description')}</p>
                            <p className="text-sm">
                                {t('system.llm.summary', { active: llmActiveCount, total: llmItems.length })}
                                {llmFailedCount > 0 && (
                                    <span className="text-destructive">
                                        {' · '}
                                        {t('system.llm.failedCount', { count: llmFailedCount })}
                                    </span>
                                )}
                            </p>
                        </SectionHeader>
                        <LlmProvidersSection
                            edits={llmEdits}
                            onEditsChange={setLlmEdits}
                            server={llm}
                        />
                        {isLlmDirty && (
                            <SectionSaveBar
                                isSaving={isSavingLlm}
                                onDiscard={() => setLlmEdits({})}
                                onSave={() =>
                                    save(
                                        () =>
                                            updateLlmProviders({
                                                variables: { input: buildLlmProviderInput(llm, llmEdits) },
                                            }),
                                        () => setLlmEdits({}),
                                    )
                                }
                            />
                        )}
                    </>
                );
            }

            case 'search-engines': {
                if (searchQuery.loading && !search) {
                    return <LoadingState title={t('system.searchEngines.loading')} />;
                }

                if (!search) {
                    return (
                        <ErrorState
                            message={searchQuery.error?.message ?? ''}
                            onRetry={searchQuery.refetch}
                            title={t('system.searchEngines.errorTitle')}
                        />
                    );
                }

                return (
                    <>
                        <SectionHeader title={t('system.searchEngines.title')}>
                            <p className="text-muted-foreground text-sm">{t('system.searchEngines.description')}</p>
                            <p className="text-sm">
                                {t('system.searchEngines.summary', {
                                    count: searchEnabledCount,
                                    total: searchItems.length,
                                })}
                            </p>
                        </SectionHeader>
                        <SearchEnginesSection
                            edits={searchEdits}
                            onEditsChange={setSearchEdits}
                            server={search}
                        />
                        {isSearchDirty && (
                            <SectionSaveBar
                                isSaving={isSavingSearch}
                                onDiscard={() => setSearchEdits({})}
                                onSave={() =>
                                    save(
                                        () =>
                                            updateSearchEngines({
                                                variables: { input: buildSearchEngineInput(search, searchEdits) },
                                            }),
                                        () => setSearchEdits({}),
                                    )
                                }
                            />
                        )}
                    </>
                );
            }
        }
    };

    return (
        <>
            <AppHeader>
                <AppHeaderContent>
                    <AppHeaderTitle>{t('system.title')}</AppHeaderTitle>
                </AppHeaderContent>
            </AppHeader>
            <div className="flex flex-1 flex-col gap-4 p-4 lg:flex-row lg:gap-6">
                <SystemSettingsNav
                    active={section}
                    categories={categories}
                    onSelect={selectSection}
                />
                <div className="flex min-w-0 flex-1 flex-col gap-4">
                    <EnvFileBanner />
                    {renderContent()}
                </div>
            </div>
        </>
    );
}

export default SettingsSystem;
