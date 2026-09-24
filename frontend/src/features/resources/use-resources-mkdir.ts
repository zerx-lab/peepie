import type { TFunction } from 'i18next';

import { useCallback, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import { z } from 'zod';

import { api, getApiErrorMessage } from '@/lib/axios';

import { RESOURCES_MKDIR_API_PATH } from './resources-constants';

export const createResourcesMkdirFormSchema = (t: TFunction<'resources'>) =>
    z.object({
        path: z
            .string()
            .trim()
            .min(1, { message: t('validation.pathRequired') })
            .refine((value) => !value.startsWith('/'), { message: t('validation.pathRelative') })
            .refine((value) => !value.split('/').includes('..'), { message: t('validation.pathNoParent') }),
    });

export type ResourcesMkdirFormValues = z.infer<ReturnType<typeof createResourcesMkdirFormSchema>>;

interface MkdirRequestBody {
    path: string;
}

interface UseResourcesMkdirResult {
    isCreating: boolean;
    mkdir: (values: ResourcesMkdirFormValues) => Promise<boolean>;
}

/** Wraps `POST /resources/mkdir` (idempotent — returns existing dir on hit). */
export function useResourcesMkdir(): UseResourcesMkdirResult {
    const [isCreating, setIsCreating] = useState(false);
    const { t } = useTranslation('resources');

    const mkdir = useCallback(
        async ({ path }: ResourcesMkdirFormValues): Promise<boolean> => {
            setIsCreating(true);

            try {
                await api.post<void, MkdirRequestBody>(RESOURCES_MKDIR_API_PATH, { path: path.trim() });

                toast.success(t('mkdir.toasts.created'), {
                    description: t('mkdir.toasts.createdDescription', { path: path.trim() }),
                });

                return true;
            } catch (error) {
                const description = getApiErrorMessage(error, t('mkdir.toasts.failedFallback'));

                toast.error(t('mkdir.toasts.failed'), { description });

                return false;
            } finally {
                setIsCreating(false);
            }
        },
        [t],
    );

    return {
        isCreating,
        mkdir,
    };
}
