import type { TFunction } from 'i18next';

import { Copy } from 'lucide-react';
import { memo, useCallback, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';

import type { VectorStoreLogFragmentFragment } from '@/graphql/types';

import Markdown from '@/components/shared/markdown';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import FlowAgentIcon from '@/features/flows/agents/flow-agent-icon';
import { VectorStoreAction } from '@/graphql/types';
import { copyMessageToClipboard } from '@/lib/clipboard';
import { formatDate } from '@/lib/utils/format';

import FlowVectorStoreActionIcon from './flow-vector-store-action-icon';

const getDescription = (log: VectorStoreLogFragmentFragment, t: TFunction<['flowDetails', 'common']>) => {
    const { action, filter } = log;
    const {
        answer_type: answerType,
        code_lang: codeLang,
        doc_type: docType,
        guide_type: guideType,
        tool_name: toolName,
    } = JSON.parse(filter) || {};

    const isStore = action === VectorStoreAction.Store;
    const parts: string[] = [];

    if (docType) {
        if (docType === 'memory') {
            parts.push(
                isStore
                    ? t('vectorStores.description.storedInMemory')
                    : t('vectorStores.description.retrievedFromMemory'),
            );
        } else {
            parts.push(
                isStore
                    ? t('vectorStores.description.storedDocType', { docType })
                    : t('vectorStores.description.retrievedDocType', { docType }),
            );
        }
    }

    if (codeLang) {
        parts.push(t('vectorStores.description.codeLang', { codeLang }));
    }

    if (toolName) {
        parts.push(t('vectorStores.description.toolName', { toolName }));
    }

    if (guideType) {
        parts.push(t('vectorStores.description.guideType', { guideType }));
    }

    if (answerType) {
        parts.push(t('vectorStores.description.answerType', { answerType }));
    }

    const description = parts.join(t('vectorStores.description.separator'));

    return description.charAt(0).toUpperCase() + description.slice(1);
};

interface FlowVectorStoreProps {
    log: VectorStoreLogFragmentFragment;
    searchValue?: string;
}

const containsSearchValue = (text: null | string | undefined, searchValue: string): boolean => {
    if (!text || !searchValue.trim()) {
        return false;
    }

    return text.toLowerCase().includes(searchValue.toLowerCase().trim());
};

function FlowVectorStore({ log, searchValue = '' }: FlowVectorStoreProps) {
    const { t } = useTranslation(['flowDetails', 'common']);
    const { action, createdAt, executor, initiator, query, result, subtaskId, taskId } = log;

    const searchChecks = useMemo(() => {
        const trimmedSearch = searchValue.trim();

        if (!trimmedSearch) {
            return { hasQueryMatch: false, hasResultMatch: false };
        }

        return {
            hasQueryMatch: containsSearchValue(query, trimmedSearch),
            hasResultMatch: containsSearchValue(result, trimmedSearch),
        };
    }, [searchValue, query, result]);

    const [isDetailsVisible, setIsDetailsVisible] = useState(false);
    const [prevSearchValue, setPrevSearchValue] = useState(searchValue);
    const [prevHasResultMatch, setPrevHasResultMatch] = useState(searchChecks.hasResultMatch);

    if (searchValue !== prevSearchValue || searchChecks.hasResultMatch !== prevHasResultMatch) {
        setPrevSearchValue(searchValue);
        setPrevHasResultMatch(searchChecks.hasResultMatch);

        const trimmedSearch = searchValue.trim();

        if (trimmedSearch) {
            if (searchChecks.hasResultMatch) {
                setIsDetailsVisible(true);
            }
        } else {
            setIsDetailsVisible(false);
        }
    }

    const description = getDescription(log, t);

    const handleCopy = useCallback(async () => {
        await copyMessageToClipboard({
            message: query,
            result: result || undefined,
        });
    }, [query, result]);

    return (
        <div className="flex flex-col items-start">
            <div className="bg-card text-card-foreground max-w-full rounded-xl border p-3 shadow-sm">
                <div className="flex flex-col">
                    <div className="cursor-pointer text-sm font-semibold">
                        <span className="inline-flex items-center gap-1">
                            <FlowVectorStoreActionIcon action={action} />
                            <span>{description}</span>
                        </span>
                    </div>

                    <Markdown
                        className="prose-xs prose-fixed wrap-break-word"
                        searchValue={searchValue}
                    >
                        {query}
                    </Markdown>
                </div>
                {result && (
                    <div className="text-muted-foreground mt-2 text-xs">
                        <div
                            className="cursor-pointer"
                            onClick={() => setIsDetailsVisible(!isDetailsVisible)}
                        >
                            {isDetailsVisible ? t('details.hide') : t('details.show')}
                        </div>
                        {isDetailsVisible && (
                            <>
                                <div className="my-3 border-t" />
                                <Markdown
                                    className="prose-xs prose-fixed wrap-break-word"
                                    searchValue={searchValue}
                                >
                                    {result}
                                </Markdown>
                            </>
                        )}
                    </div>
                )}
            </div>
            <div className="text-muted-foreground mt-1 flex items-center gap-1 px-1 text-xs">
                <span className="flex items-center gap-0.5">
                    <FlowAgentIcon
                        className="text-muted-foreground"
                        type={initiator}
                    />
                    <span className="text-muted-foreground/50">→</span>
                    <FlowAgentIcon
                        className="text-muted-foreground"
                        type={executor}
                    />
                </span>
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
                {taskId && (
                    <>
                        <span className="text-muted-foreground/50">|</span>
                        <span className="text-muted-foreground/50">{t('meta.taskId', { id: taskId })}</span>
                    </>
                )}
                {subtaskId && (
                    <>
                        <span className="text-muted-foreground/50">|</span>
                        <span className="text-muted-foreground/50">{t('meta.subtaskId', { id: subtaskId })}</span>
                    </>
                )}
            </div>
        </div>
    );
}

export default memo(FlowVectorStore);
