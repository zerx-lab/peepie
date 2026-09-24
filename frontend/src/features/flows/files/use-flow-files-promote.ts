import type { TFunction } from 'i18next';

import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { z } from 'zod';

import type { OverwriteOutcome } from '@/components/shared/overwrite';
import type { RestResourceList } from '@/features/resources/resources-rest';

import { api, getApiErrorMessage, getApiErrorStatusCode } from '@/lib/axios';

import { FLOW_FILES_PROMOTE_API_PATH } from './flow-files-constants';

export const createFlowFilesPromoteFormSchema = (t: TFunction<'fileManager'>) =>
    z.object({
        destination: z
            .string()
            .trim()
            .min(1, { message: t('promoteDialog.validation.required') })
            .refine((value) => !value.startsWith('/'), { message: t('promoteDialog.validation.relative') })
            .refine((value) => !value.split('/').includes('..'), {
                message: t('promoteDialog.validation.noParentSegments'),
            }),
    });

export type FlowFilesPromoteFormValues = z.infer<ReturnType<typeof createFlowFilesPromoteFormSchema>>;

interface PromoteRequestBody {
    destination: string;
    force: boolean;
    sources: readonly string[];
}

interface UseFlowFilesPromoteParams {
    flowId: null | string;
}

interface UseFlowFilesPromoteResult {
    isPromoting: boolean;
    /**
     * Issue a batch promote in a single atomic request and return a discriminated outcome:
     *   - `ok`        — every flow file/dir was promoted (success toast already fired),
     *   - `conflict`  — at least one resource path is occupied (no toast, caller
     *                   resolves via the shared overwrite workflow),
     *   - `error`     — anything else (failure toast already fired).
     *
     * Backend semantics: with one source, `destination` is the exact target
     * path; with multiple sources, `destination` is a base directory and each
     * source lands at `destination/<basename>`. The Apollo cache stays in sync
     * via `resourceAdded` / `resourceUpdated` GraphQL subscriptions.
     */
    promote: (sources: readonly string[], destination: string, force: boolean) => Promise<OverwriteOutcome>;
}

/**
 * Wraps the "promote flow file → user resource" REST call (`POST /files/to-resources`)
 * with toast notifications and a loading flag.
 */
export function useFlowFilesPromote({ flowId }: UseFlowFilesPromoteParams): UseFlowFilesPromoteResult {
    const { t } = useTranslation('fileManager');
    const [isPromoting, setIsPromoting] = useState(false);

    const promote = useCallback(
        async (sources: readonly string[], destination: string, force: boolean): Promise<OverwriteOutcome> => {
            if (!flowId || sources.length === 0) {
                return { kind: 'error' };
            }

            setIsPromoting(true);

            try {
                await api.post<RestResourceList, PromoteRequestBody>(
                    FLOW_FILES_PROMOTE_API_PATH(flowId),
                    {
                        destination: destination.trim(),
                        force,
                        sources,
                    },
                    { timeout: 0 },
                );

                const description =
                    sources.length === 1
                        ? t('flowFiles.toasts.savedOneDescription', { destination: destination.trim() })
                        : t('flowFiles.toasts.savedManyDescription', {
                              count: sources.length,
                              destination: destination.trim(),
                          });

                toast.success(t('flowFiles.toasts.saved'), { description });

                return { kind: 'ok' };
            } catch (error) {
                if (!force && getApiErrorStatusCode(error) === 409) {
                    return { kind: 'conflict' };
                }

                const description = getApiErrorMessage(error, t('flowFiles.toasts.saveFailedFallback'));

                toast.error(t('flowFiles.toasts.saveFailed'), { description });

                return { kind: 'error' };
            } finally {
                setIsPromoting(false);
            }
        },
        [flowId, t],
    );

    return {
        isPromoting,
        promote,
    };
}
