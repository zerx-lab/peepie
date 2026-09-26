import { Copy, RefreshCw } from 'lucide-react';
import { memo, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import type { AssistantLogFragmentFragment, MessageLogFragmentFragment } from '@/graphql/types';

import Markdown from '@/components/shared/markdown';
import Terminal from '@/components/shared/terminal';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { MessageLogType, ResultFormat } from '@/graphql/types';
import { copyMessageToClipboard } from '@/lib/clipboard';
import { cn } from '@/lib/utils';
import { formatDate } from '@/lib/utils/format';

import FlowMessageTypeIcon from './flow-message-type-icon';

interface FlowMessageProps {
    log: AssistantLogFragmentFragment | MessageLogFragmentFragment;
    occurrenceTimestamps?: string[];
    retryInfo?: { lastError: string; retryCount: number };
    searchValue?: string;
}

const containsSearchValue = (text: null | string | undefined, searchValue: string): boolean => {
    if (!text || !searchValue.trim()) {
        return false;
    }

    return text.toLowerCase().includes(searchValue.toLowerCase().trim());
};

function FlowMessage({ log, occurrenceTimestamps, retryInfo, searchValue = '' }: FlowMessageProps) {
    const { t } = useTranslation(['flowDetails', 'common']);
    const { createdAt, message, result, resultFormat = ResultFormat.Plain, thinking, type } = log;
    const isReportMessage = type === MessageLogType.Report;

    const searchChecks = useMemo(() => {
        const trimmedSearch = searchValue.trim();

        if (!trimmedSearch) {
            return { hasResultMatch: false, hasThinkingMatch: false };
        }

        return {
            hasResultMatch: containsSearchValue(result, trimmedSearch),
            hasThinkingMatch: containsSearchValue(thinking, trimmedSearch),
        };
    }, [searchValue, thinking, result]);

    const [isDetailsVisible, setIsDetailsVisible] = useState(isReportMessage);
    const [isThinkingVisible, setIsThinkingVisible] = useState(false);

    const [prevSearchValue, setPrevSearchValue] = useState(searchValue);
    const [prevHasThinkingMatch, setPrevHasThinkingMatch] = useState(searchChecks.hasThinkingMatch);
    const [prevHasResultMatch, setPrevHasResultMatch] = useState(searchChecks.hasResultMatch);

    if (
        searchValue !== prevSearchValue ||
        searchChecks.hasThinkingMatch !== prevHasThinkingMatch ||
        searchChecks.hasResultMatch !== prevHasResultMatch
    ) {
        setPrevSearchValue(searchValue);
        setPrevHasThinkingMatch(searchChecks.hasThinkingMatch);
        setPrevHasResultMatch(searchChecks.hasResultMatch);

        const trimmedSearch = searchValue.trim();

        if (trimmedSearch) {
            if (searchChecks.hasThinkingMatch) {
                setIsThinkingVisible(true);
            }

            if (searchChecks.hasResultMatch) {
                setIsDetailsVisible(true);
            }
        } else {
            setIsDetailsVisible(isReportMessage);
            setIsThinkingVisible(false);
        }
    }

    const toggleDetails = useCallback(() => {
        setIsDetailsVisible((prev) => !prev);
    }, []);

    const toggleThinking = useCallback(() => {
        setIsThinkingVisible((prev) => !prev);
    }, []);

    const handleCopy = useCallback(async () => {
        await copyMessageToClipboard({
            message,
            result,
            resultFormat,
            thinking,
        });
    }, [thinking, message, result, resultFormat]);

    const shouldShowThinking = thinking && (!message || isThinkingVisible);

    const shouldShowThinkingToggle = thinking && message;

    // Gate by visibility so the (potentially heavy) details subtree isn't in the DOM when collapsed.
    const renderDetailsContent = () => {
        if (!isDetailsVisible) {
            return null;
        }

        return (
            <>
                <div className="my-3 border-t" />
                {resultFormat === ResultFormat.Plain && (
                    <Markdown
                        className="prose-xs prose-fixed text-accent-foreground text-sm wrap-break-word"
                        searchValue={searchValue}
                    >
                        {result}
                    </Markdown>
                )}
                {resultFormat === ResultFormat.Markdown && (
                    <Markdown
                        className="prose-xs prose-fixed wrap-break-word"
                        searchValue={searchValue}
                    >
                        {result}
                    </Markdown>
                )}
                {resultFormat === ResultFormat.Terminal && (
                    <Terminal
                        className="bg-card h-[240px] w-full py-1 pl-1"
                        logs={[result as string]}
                    />
                )}
            </>
        );
    };

    const renderThinkingContent = () => {
        if (!shouldShowThinking) {
            return null;
        }

        return (
            <>
                <div className="border-muted mb-3 border-l-2 pl-3">
                    <Markdown
                        className="prose-xs prose-fixed text-muted-foreground/80 wrap-break-word"
                        searchValue={searchValue}
                    >
                        {thinking}
                    </Markdown>
                </div>
            </>
        );
    };

    return (
        <div className={`flex flex-col ${type === MessageLogType.Input ? 'items-end' : 'items-start'}`}>
            <div
                className={cn(
                    'bg-card text-card-foreground max-w-[90%] rounded-xl border p-3 shadow-sm',
                    resultFormat === ResultFormat.Terminal && isDetailsVisible ? 'w-full' : '',
                )}
            >
                {shouldShowThinkingToggle && (
                    <div className="text-muted-foreground mb-2 text-xs">
                        <div
                            className="cursor-pointer"
                            onClick={toggleThinking}
                        >
                            {isThinkingVisible ? t('messages.hideThinking') : t('messages.showThinking')}
                        </div>
                    </div>
                )}

                {renderThinkingContent()}

                {message && (
                    <Markdown
                        className="prose-xs prose-fixed wrap-break-word"
                        searchValue={searchValue}
                    >
                        {message}
                    </Markdown>
                )}

                {result && (
                    <div className="text-muted-foreground mt-2 text-xs">
                        <div
                            className="cursor-pointer"
                            onClick={toggleDetails}
                        >
                            {isDetailsVisible ? t('details.hide') : t('details.show')}
                        </div>
                        {renderDetailsContent()}
                    </div>
                )}
            </div>
            <div
                className={`text-muted-foreground mt-1 flex items-center gap-1 px-1 text-xs ${
                    type === MessageLogType.Input ? 'flex-row-reverse' : 'flex-row'
                }`}
            >
                <FlowMessageTypeIcon type={type} />
                {occurrenceTimestamps && occurrenceTimestamps.length > 1 && (
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <span className="border-amber-500/40 bg-amber-500/10 text-amber-600 dark:text-amber-400 mx-0.5 inline-flex items-center gap-0.5 rounded-full border px-1.5 py-0.5 text-[10px] font-medium leading-none">
                                <RefreshCw className="size-2.5" />
                                {t('messages.retryBadge', { count: occurrenceTimestamps.length })}
                            </span>
                        </TooltipTrigger>
                        <TooltipContent className="max-w-64 whitespace-pre-line">
                            {t('messages.retryTooltip', {
                                times: occurrenceTimestamps.map((ts) => formatDate(new Date(ts))).join('\n'),
                            })}
                        </TooltipContent>
                    </Tooltip>
                )}
                {retryInfo && retryInfo.retryCount > 0 && (
                    <Tooltip>
                        <TooltipTrigger asChild>
                            <span className="border-destructive/40 bg-destructive/10 text-destructive mx-0.5 inline-flex items-center gap-0.5 rounded-full border px-1.5 py-0.5 text-[10px] font-medium leading-none">
                                <RefreshCw className="size-2.5" />
                                {t('messages.failureBadge', { count: retryInfo.retryCount })}
                            </span>
                        </TooltipTrigger>
                        <TooltipContent className="max-w-80 whitespace-pre-line">{retryInfo.lastError}</TooltipContent>
                    </Tooltip>
                )}
                <Tooltip>
                    <TooltipTrigger asChild>
                        <Copy
                            className="hover:text-foreground mx-1 size-3 shrink-0 cursor-pointer transition-colors"
                            onClick={handleCopy}
                        />
                    </TooltipTrigger>
                    <TooltipContent>{t('common:actions.copy')}</TooltipContent>
                </Tooltip>
                <span className="text-muted-foreground/50">{formatDate(new Date(createdAt))}</span>
                <span
                    className="text-muted-foreground/50"
                    data-slot="flow-message-id"
                >
                    {log.id}
                </span>
            </div>
        </div>
    );
}

export default memo(FlowMessage);
