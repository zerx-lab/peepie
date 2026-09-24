import type { TFunction } from 'i18next';

import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import * as z from 'zod';

import { Button } from '@/components/ui/button';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { FormSubmitButton } from '@/components/ui/form-submit-button';
import { Input } from '@/components/ui/input';
import { InputPassword } from '@/components/ui/input-password';
import { useAppForm } from '@/hooks/use-app-form';
import { api, resolveApiErrorMessage } from '@/lib/axios';
import { useUser } from '@/providers/user-provider';

const createEmailChangeSchema = (t: TFunction<'auth'>) =>
    z.object({
        currentPassword: z.string().min(1, { message: t('emailChange.validation.currentPasswordRequired') }),
        newEmail: z
            .string()
            .trim()
            .toLowerCase()
            .min(1, { message: t('emailChange.validation.emailRequired') })
            .email({ message: t('emailChange.validation.invalidEmail') })
            .max(50, { message: t('emailChange.validation.emailTooLong') }),
    });

interface EmailChangeFormProps {
    onCancel?: () => void;
    onSuccess?: () => void;
}

type EmailChangeFormValues = z.infer<ReturnType<typeof createEmailChangeSchema>>;

export function EmailChangeForm({ onCancel, onSuccess }: EmailChangeFormProps) {
    const { t } = useTranslation(['auth', 'common']);
    const [error, setError] = useState<null | string>(null);
    const { patchUser, refreshAuthInfo } = useUser();
    const emailChangeSchema = useMemo(() => createEmailChangeSchema(t), [t]);

    const form = useAppForm<EmailChangeFormValues>({
        defaultValues: {
            currentPassword: '',
            newEmail: '',
        },
        schema: emailChangeSchema,
    });

    const handleSubmit = async (values: EmailChangeFormValues) => {
        setError(null);

        try {
            await api.put('/user/email', {
                current_password: values.currentPassword,
                mail: values.newEmail,
            });

            form.reset();
            toast.success(t('emailChange.success'));

            patchUser({ mail: values.newEmail });
            await refreshAuthInfo();

            onSuccess?.();
        } catch (err: unknown) {
            const errorByCode: Record<string, string> = {
                'Users.ChangeEmailCurrentUser.EmailAlreadyExists': t('emailChange.errors.emailAlreadyExists'),
                'Users.ChangeEmailCurrentUser.InvalidCurrentPassword': t('emailChange.errors.invalidCurrentPassword'),
                'Users.ChangeEmailCurrentUser.InvalidEmail': t('emailChange.errors.invalidEmail'),
                'Users.NotFound': t('emailChange.errors.userNotFound'),
            };

            setError(resolveApiErrorMessage(err, errorByCode, t('emailChange.errors.failed')));
        }
    };

    return (
        <Form {...form}>
            {/* noValidate: the type="email" field would otherwise fire the browser's native (locale-styled)
                validation popup on submit, pre-empting our zod message. Validation runs through zod instead. */}
            <form
                className="flex flex-col gap-4"
                noValidate
                onSubmit={form.handleSubmit(handleSubmit)}
            >
                <FormField
                    control={form.control}
                    name="currentPassword"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t('emailChange.currentPasswordLabel')}</FormLabel>
                            <FormControl>
                                <InputPassword
                                    {...field}
                                    placeholder={t('emailChange.currentPasswordPlaceholder')}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="newEmail"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t('emailChange.newEmailLabel')}</FormLabel>
                            <FormControl>
                                <Input
                                    {...field}
                                    placeholder={t('emailChange.newEmailPlaceholder')}
                                    type="email"
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
                        <span>{t('emailChange.submit')}</span>
                    </FormSubmitButton>
                </div>
            </form>
        </Form>
    );
}
