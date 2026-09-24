import type { TFunction } from 'i18next';

import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import * as z from 'zod';

import { Button } from '@/components/ui/button';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { FormSubmitButton } from '@/components/ui/form-submit-button';
import { Input } from '@/components/ui/input';
import { useAppForm } from '@/hooks/use-app-form';
import { api, resolveApiErrorMessage } from '@/lib/axios';
import { useUser } from '@/providers/user-provider';

const createNameChangeSchema = (t: TFunction<'auth'>) =>
    z.object({
        name: z
            .string()
            .trim()
            .min(1, { message: t('nameChange.validation.required') })
            .max(70, { message: t('nameChange.validation.tooLong') }),
    });

interface NameChangeFormProps {
    onCancel?: () => void;
    onSuccess?: () => void;
}

type NameChangeFormValues = z.infer<ReturnType<typeof createNameChangeSchema>>;

export function NameChangeForm({ onCancel, onSuccess }: NameChangeFormProps) {
    const { t } = useTranslation(['auth', 'common']);
    const [error, setError] = useState<null | string>(null);
    const { authInfo, patchUser, refreshAuthInfo } = useUser();
    const nameChangeSchema = useMemo(() => createNameChangeSchema(t), [t]);

    const form = useAppForm<NameChangeFormValues>({
        defaultValues: {
            name: authInfo?.user?.name ?? '',
        },
        schema: nameChangeSchema,
    });

    const handleSubmit = async (values: NameChangeFormValues) => {
        setError(null);

        try {
            await api.put('/user/name', { name: values.name });

            toast.success(t('nameChange.success'));

            patchUser({ name: values.name });
            await refreshAuthInfo();

            onSuccess?.();
        } catch (err: unknown) {
            const errorByCode: Record<string, string> = {
                'Users.ChangeNameCurrentUser.InvalidName': t('nameChange.errors.invalidName'),
                'Users.NotFound': t('nameChange.errors.userNotFound'),
            };

            setError(resolveApiErrorMessage(err, errorByCode, t('nameChange.errors.failed')));
        }
    };

    return (
        <Form {...form}>
            <form
                className="flex flex-col gap-4"
                noValidate
                onSubmit={form.handleSubmit(handleSubmit)}
            >
                <FormField
                    control={form.control}
                    name="name"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t('nameChange.label')}</FormLabel>
                            <FormControl>
                                <Input
                                    {...field}
                                    placeholder={t('nameChange.placeholder')}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {error && <div className="text-destructive text-sm">{error}</div>}

                <div className="flex justify-end gap-2 pt-2">
                    {onCancel && (
                        <Button
                            onClick={onCancel}
                            size="sm"
                            type="button"
                            variant="outline"
                        >
                            {t('common:actions.cancel')}
                        </Button>
                    )}
                    <FormSubmitButton size="sm">
                        <span>{t('nameChange.submit')}</span>
                    </FormSubmitButton>
                </div>
            </form>
        </Form>
    );
}
