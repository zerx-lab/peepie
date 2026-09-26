import { zodResolver } from '@hookform/resolvers/zod';
import { ChevronDown, Inbox, ListFilter, Search, X } from 'lucide-react';
import { useEffect, useMemo, useState } from 'react';
import { useForm } from 'react-hook-form';
import { useTranslation } from 'react-i18next';
import { useDebouncedCallback } from 'use-debounce';
import { z } from 'zod';

import { Button } from '@/components/ui/button';
import { Empty, EmptyContent, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty';
import { Form, FormControl, FormField } from '@/components/ui/form';
import { InputGroup, InputGroupAddon, InputGroupButton, InputGroupInput } from '@/components/ui/input-group';
import FlowAgentIcon, { useAgentTypeLabel } from '@/features/flows/agents/flow-agent-icon';
import { MessageLogType, StatusType } from '@/graphql/types';
import { useAutoScroll } from '@/hooks/use-auto-scroll';
import { cn } from '@/lib/utils';
import { useFlow } from '@/providers/flow-provider';

import { FlowForm, type FlowFormValues } from '../flow-form';
import FlowTasksDropdown from '../flow-tasks-dropdown';
import FlowMessage from './flow-message';

interface FlowAutomationMessagesProps {
    className?: string;
}

const searchFormSchema = z.object({
    filter: z
        .object({
            subtaskIds: z.array(z.string()),
            taskIds: z.array(z.string()),
        })
        .optional(),
    search: z.string(),
});

function FlowAutomationMessages({ className }: FlowAutomationMessagesProps) {
    const { t } = useTranslation('flowDetails');
    const getAgentTypeLabel = useAgentTypeLabel();
    const { flowData, flowId, flowStatus, stopAutomation, submitAutomationMessage } = useFlow();

    const logs = useMemo(() => flowData?.messageLogs ?? [], [flowData?.messageLogs]);

    // Real failure/retry info persisted by the backend for the task/subtask that
    // produced each message, keyed by id so the timeline can show *why* a step
    // stalled instead of the previous heuristic-only "repeated N times" badge.
    const retryInfoById = useMemo(() => {
        const info: Record<string, { lastError: string; retryCount: number }> = {};

        for (const task of flowData?.tasks ?? []) {
            if (task.retryCount > 0) {
                info[task.id] = { lastError: task.lastError, retryCount: task.retryCount };
            }

            for (const subtask of task.subtasks ?? []) {
                if (subtask.retryCount > 0) {
                    info[subtask.id] = { lastError: subtask.lastError, retryCount: subtask.retryCount };
                }
            }
        }

        return info;
    }, [flowData?.tasks]);

    // The most recently created agent log tells us which agent role is currently
    // handling the flow; shown next to the input so the user isn't guessing what
    // is happening while waiting (agent logs carry no live "in progress" flag,
    // only completed steps, so this reflects the last agent that reported back).
    const currentAgent = useMemo(() => {
        const agentLogs = flowData?.agentLogs;

        if (!agentLogs || agentLogs.length === 0) {
            return undefined;
        }

        return agentLogs.reduce((latest, log) =>
            new Date(log.createdAt).getTime() > new Date(latest.createdAt).getTime() ? log : latest,
        ).executor;
    }, [flowData?.agentLogs]);

    const [debouncedSearchValue, setDebouncedSearchValue] = useState('');
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [isCanceling, setIsCanceling] = useState(false);

    const { containerRef, endRef, hasNewMessages, isScrolledToBottom, scrollToEnd } = useAutoScroll(logs, flowId);

    const form = useForm<z.infer<typeof searchFormSchema>>({
        defaultValues: {
            filter: {
                subtaskIds: [],
                taskIds: [],
            },
            search: '',
        },
        resolver: zodResolver(searchFormSchema),
    });

    const searchValue = form.watch('search');
    const filter = form.watch('filter');

    const debouncedUpdateSearch = useDebouncedCallback((value: string) => {
        setDebouncedSearchValue(value);
    }, 500);

    useEffect(() => {
        debouncedUpdateSearch(searchValue);

        return () => {
            debouncedUpdateSearch.cancel();
        };
    }, [searchValue, debouncedUpdateSearch]);

    useEffect(() => {
        return () => {
            debouncedUpdateSearch.cancel();
        };
    }, [debouncedUpdateSearch]);

    useEffect(() => {
        form.reset({
            filter: {
                subtaskIds: [],
                taskIds: [],
            },
            search: '',
        });
        setDebouncedSearchValue('');
        debouncedUpdateSearch.cancel();
    }, [flowId, form, debouncedUpdateSearch]);

    const hasActiveFilters = useMemo(() => {
        const hasSearch = !!searchValue.trim();
        const hasTaskFilters = !!(filter?.taskIds?.length || filter?.subtaskIds?.length);

        return hasSearch || hasTaskFilters;
    }, [searchValue, filter]);

    const filteredLogs = useMemo(() => {
        const search = debouncedSearchValue.toLowerCase().trim();

        let filtered = logs || [];

        if (search) {
            filtered = filtered.filter(
                (log) =>
                    log.message.toLowerCase().includes(search) ||
                    (log.result && log.result.toLowerCase().includes(search)) ||
                    (log.thinking && log.thinking.toLowerCase().includes(search)),
            );
        }

        if (filter?.taskIds?.length || filter?.subtaskIds?.length) {
            const selectedTaskIds = new Set(filter.taskIds ?? []);
            const selectedSubtaskIds = new Set(filter.subtaskIds ?? []);

            filtered = filtered.filter((log) => {
                if (log.taskId && selectedTaskIds.has(log.taskId)) {
                    return true;
                }

                if (log.subtaskId && selectedSubtaskIds.has(log.subtaskId)) {
                    return true;
                }

                return false;
            });
        }

        return filtered;
    }, [logs, debouncedSearchValue, filter]);

    // Repeated tool-call failures (e.g. parsing errors) make the agent re-emit the
    // same "thinking" title over and over with no result to show. Collapsing those
    // consecutive duplicates into one card with a retry badge keeps the timeline
    // readable instead of showing N near-identical, content-less cards.
    const groupedLogs = useMemo(() => {
        const groups: Array<{ occurrences: string[]; representative: (typeof filteredLogs)[number] }> = [];

        for (const log of filteredLogs) {
            const isRetryCandidate = log.type === MessageLogType.Thoughts && !!log.message && !log.result;
            const previous = groups[groups.length - 1];

            if (
                isRetryCandidate &&
                previous &&
                previous.representative.type === log.type &&
                previous.representative.message === log.message &&
                !previous.representative.result
            ) {
                previous.representative = log;
                previous.occurrences.push(log.createdAt);
            } else {
                groups.push({ occurrences: [log.createdAt], representative: log });
            }
        }

        return groups;
    }, [filteredLogs]);

    const placeholder = useMemo(() => {
        if (!flowId) {
            return t('messages.placeholders.selectFlow');
        }

        switch (flowStatus) {
            case StatusType.Created: {
                return t('messages.placeholders.automationStarting');
            }

            case StatusType.Failed:
            case StatusType.Finished: {
                return t('messages.placeholders.automationEnded');
            }

            case StatusType.Running: {
                return t('messages.placeholders.automationRunning');
            }

            case StatusType.Waiting: {
                return t('messages.placeholders.automationWaiting');
            }

            default: {
                return t('messages.placeholders.default');
            }
        }
    }, [flowId, flowStatus, t]);

    const handleSubmitMessage = async (values: FlowFormValues) => {
        setIsSubmitting(true);

        try {
            await submitAutomationMessage(values);
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleStopAutomation = async () => {
        setIsCanceling(true);

        try {
            await stopAutomation();
        } finally {
            setIsCanceling(false);
        }
    };

    const handleResetFilters = () => {
        form.reset({
            filter: {
                subtaskIds: [],
                taskIds: [],
            },
            search: '',
        });
        setDebouncedSearchValue('');
        debouncedUpdateSearch.cancel();
    };

    const isFormDisabled = flowStatus === StatusType.Finished || flowStatus === StatusType.Failed;
    const isFormLoading = flowStatus === StatusType.Created || flowStatus === StatusType.Running;
    const isProviderChangeAllowed = flowStatus === StatusType.Waiting;

    return (
        <div className={cn('flex h-full flex-col', className)}>
            <div className="bg-background sticky top-0 z-10 pb-4">
                <Form {...form}>
                    <div className="flex gap-2 p-px">
                        <FormField
                            control={form.control}
                            name="search"
                            render={({ field }) => (
                                <FormControl>
                                    <InputGroup className="flex-1">
                                        <InputGroupAddon>
                                            <Search />
                                        </InputGroupAddon>
                                        <InputGroupInput
                                            {...field}
                                            autoComplete="off"
                                            placeholder={t('messages.searchPlaceholder')}
                                            type="text"
                                        />
                                        {field.value && (
                                            <InputGroupAddon align="inline-end">
                                                <InputGroupButton
                                                    aria-label={t('messages.clearSearch')}
                                                    onClick={() => {
                                                        form.reset({ search: '' });
                                                        setDebouncedSearchValue('');
                                                        debouncedUpdateSearch.cancel();
                                                    }}
                                                    type="button"
                                                >
                                                    <X />
                                                </InputGroupButton>
                                            </InputGroupAddon>
                                        )}
                                    </InputGroup>
                                </FormControl>
                            )}
                        />
                        <FormField
                            control={form.control}
                            name="filter"
                            render={({ field }) => (
                                <FormControl>
                                    <FlowTasksDropdown
                                        onChange={field.onChange}
                                        value={field.value}
                                    />
                                </FormControl>
                            )}
                        />
                    </div>
                </Form>
            </div>

            {filteredLogs.length > 0 ? (
                <div className="relative h-full overflow-y-hidden">
                    <div
                        className="flex h-full flex-col gap-4 overflow-y-auto"
                        ref={containerRef}
                    >
                        {groupedLogs.map((group) => (
                            <FlowMessage
                                key={group.representative.id}
                                log={group.representative}
                                occurrenceTimestamps={group.occurrences}
                                retryInfo={
                                    (group.representative.subtaskId &&
                                        retryInfoById[group.representative.subtaskId]) ||
                                    (group.representative.taskId && retryInfoById[group.representative.taskId]) ||
                                    undefined
                                }
                                searchValue={debouncedSearchValue}
                            />
                        ))}
                        <div ref={endRef} />
                    </div>

                    {!isScrolledToBottom && (
                        <Button
                            aria-label={t('messages.scrollToLatest')}
                            className="absolute right-4 bottom-4 z-10 shadow-md hover:shadow-lg"
                            onClick={() => scrollToEnd()}
                            size="icon-sm"
                            type="button"
                            variant="outline"
                        >
                            <ChevronDown />
                            {hasNewMessages && (
                                <span className="bg-primary absolute -top-1.5 -right-1.5 size-3 rounded-full" />
                            )}
                        </Button>
                    )}
                </div>
            ) : hasActiveFilters ? (
                <Empty>
                    <EmptyHeader>
                        <EmptyMedia variant="icon">
                            <ListFilter />
                        </EmptyMedia>
                        <EmptyTitle>{t('messages.notFound')}</EmptyTitle>
                        <EmptyDescription>{t('filters.adjustHint')}</EmptyDescription>
                    </EmptyHeader>
                    <EmptyContent>
                        <Button
                            onClick={handleResetFilters}
                            variant="outline"
                        >
                            <X />
                            {t('filters.reset')}
                        </Button>
                    </EmptyContent>
                </Empty>
            ) : (
                <Empty>
                    <EmptyHeader>
                        <EmptyMedia variant="icon">
                            <Inbox />
                        </EmptyMedia>
                        <EmptyTitle>{t('messages.automationEmpty.title')}</EmptyTitle>
                        <EmptyDescription>{t('messages.automationEmpty.description')}</EmptyDescription>
                    </EmptyHeader>
                </Empty>
            )}

            <div className="bg-background sticky bottom-0 p-px">
                {flowStatus === StatusType.Running && currentAgent && (
                    <div className="text-muted-foreground mb-2 flex w-fit items-center gap-1.5 rounded-full border px-2 py-1 text-xs">
                        <span className="relative flex size-1.5">
                            <span className="bg-primary absolute inline-flex h-full w-full animate-ping rounded-full opacity-75" />
                            <span className="bg-primary relative inline-flex size-1.5 rounded-full" />
                        </span>
                        <span>{t('messages.currentAgent')}</span>
                        <FlowAgentIcon
                            className="text-foreground"
                            type={currentAgent}
                        />
                        <span className="text-foreground font-medium">{getAgentTypeLabel(currentAgent)}</span>
                    </div>
                )}
                <FlowForm
                    defaultValues={{
                        providerName: flowData?.flow?.provider?.name ?? '',
                    }}
                    isCanceling={isCanceling}
                    isDisabled={isFormDisabled}
                    isLoading={isFormLoading}
                    isProviderDisabled={!isProviderChangeAllowed}
                    isSubmitting={isSubmitting}
                    onCancel={handleStopAutomation}
                    onSubmit={handleSubmitMessage}
                    placeholder={placeholder}
                    type={'automation'}
                />
            </div>
        </div>
    );
}

export default FlowAutomationMessages;
