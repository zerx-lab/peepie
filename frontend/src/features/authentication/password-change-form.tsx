import type { TFunction } from 'i18next';

import { type ComponentProps, useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { toast } from 'sonner';
import * as z from 'zod';

import { Button } from '@/components/ui/button';
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { FormSubmitButton } from '@/components/ui/form-submit-button';
import { InputPassword } from '@/components/ui/input-password';
import { useAppForm } from '@/hooks/use-app-form';
import { api, resolveApiErrorMessage } from '@/lib/axios';
import { cn } from '@/lib/utils';

/**
 * Password policy (mirrors the backend): 8–72 bytes; either 16+ characters of any
 * composition, or 8–15 characters with a digit, lowercase, uppercase and a special
 * character from `!@#$&*`. Only the messages are translated — keep the rules intact.
 */
const createPasswordChangeSchema = (t: TFunction<'auth'>) =>
    z
        .object({
            confirmPassword: z.string().min(1, { message: t('passwordChange.validation.confirmRequired') }),
            currentPassword: z.string().min(1, { message: t('passwordChange.validation.currentRequired') }),
            newPassword: z
                .string()
                .min(8, { message: t('passwordChange.validation.minLength') })
                // bcrypt, which hashes it server-side, refuses anything longer than 72 bytes.
                .refine((password) => new TextEncoder().encode(password).length <= 72, {
                    message: t('passwordChange.validation.maxLength'),
                })
                .refine(
                    (password) => {
                        if (password.length > 15) {
                            return true;
                        }

                        return (
                            password.length >= 8 &&
                            /[0-9]/.test(password) &&
                            /[a-z]/.test(password) &&
                            /[A-Z]/.test(password) &&
                            /[!@#$&*]/.test(password)
                        );
                    },
                    {
                        message: t('passwordChange.validation.complexity'),
                    },
                ),
        })
        .refine((data) => data.newPassword === data.confirmPassword, {
            message: t('passwordChange.validation.mismatch'),
            path: ['confirmPassword'],
        })
        .refine((data) => data.currentPassword !== data.newPassword, {
            message: t('passwordChange.validation.sameAsCurrent'),
            path: ['newPassword'],
        });

interface PasswordChangeFormProps {
    buttonSize?: ComponentProps<typeof Button>['size'];
    layout?: 'horizontal' | 'vertical';
    onCancel?: () => void;
    onSkip?: () => void;
    onSuccess?: () => void;
}

type PasswordChangeFormValues = z.infer<ReturnType<typeof createPasswordChangeSchema>>;

export function PasswordChangeForm({
    buttonSize = 'default',
    layout = 'horizontal',
    onCancel,
    onSkip,
    onSuccess,
}: PasswordChangeFormProps) {
    const { t } = useTranslation(['auth', 'common']);
    const [error, setError] = useState<null | string>(null);
    const passwordChangeSchema = useMemo(() => createPasswordChangeSchema(t), [t]);

    const form = useAppForm<PasswordChangeFormValues>({
        defaultValues: {
            confirmPassword: '',
            currentPassword: '',
            newPassword: '',
        },
        schema: passwordChangeSchema,
    });

    const handleSubmit = async (values: PasswordChangeFormValues) => {
        setError(null);

        try {
            await api.put('/user/password', {
                confirm_password: values.confirmPassword,
                current_password: values.currentPassword,
                password: values.newPassword,
            });

            form.reset();
            toast.success(t('passwordChange.success'));

            onSuccess?.();
        } catch (err: unknown) {
            const errorByCode: Record<string, string> = {
                'Users.ChangePasswordCurrentUser.InvalidCurrentPassword': t(
                    'passwordChange.errors.invalidCurrentPassword',
                ),
                'Users.ChangePasswordCurrentUser.InvalidNewPassword': t('passwordChange.errors.invalidNewPassword'),
                'Users.ChangePasswordCurrentUser.InvalidPassword': t('passwordChange.errors.invalidPassword'),
                'Users.NotFound': t('passwordChange.errors.userNotFound'),
            };

            setError(resolveApiErrorMessage(err, errorByCode, t('passwordChange.errors.failed')));
        }
    };

    const isVertical = layout === 'vertical';

    const skipButton = onSkip && (
        <Button
            className={cn('text-muted-foreground', isVertical && 'w-full')}
            onClick={onSkip}
            size={buttonSize}
            type="button"
            variant="ghost"
        >
            {t('passwordChange.skip')}
        </Button>
    );
    const cancelButton = onCancel && (
        <Button
            className={cn(isVertical && 'w-full')}
            onClick={onCancel}
            size={buttonSize}
            type="button"
            variant="outline"
        >
            {t('common:actions.cancel')}
        </Button>
    );
    const submitButton = (
        <FormSubmitButton
            className={cn(isVertical && 'w-full')}
            size={buttonSize}
        >
            <span>{t('passwordChange.submit')}</span>
        </FormSubmitButton>
    );

    return (
        <Form {...form}>
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
                            <FormLabel>{t('passwordChange.currentLabel')}</FormLabel>
                            <FormControl>
                                <InputPassword
                                    {...field}
                                    placeholder={t('passwordChange.currentPlaceholder')}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="newPassword"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t('passwordChange.newLabel')}</FormLabel>
                            <FormControl>
                                <InputPassword
                                    {...field}
                                    placeholder={t('passwordChange.newPlaceholder')}
                                />
                            </FormControl>
                            <FormDescription className="text-xs">{t('passwordChange.newDescription')}</FormDescription>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                <FormField
                    control={form.control}
                    name="confirmPassword"
                    render={({ field }) => (
                        <FormItem>
                            <FormLabel>{t('passwordChange.confirmLabel')}</FormLabel>
                            <FormControl>
                                <InputPassword
                                    {...field}
                                    placeholder={t('passwordChange.confirmPlaceholder')}
                                />
                            </FormControl>
                            <FormMessage />
                        </FormItem>
                    )}
                />

                {error && <div className="text-destructive text-sm">{error}</div>}

                {isVertical ? (
                    <div className="flex flex-col gap-2 pt-2">
                        {submitButton}
                        {cancelButton}
                        {skipButton}
                    </div>
                ) : (
                    <div className="flex justify-end gap-2 pt-2">
                        {skipButton}
                        {cancelButton}
                        {submitButton}
                    </div>
                )}
            </form>
        </Form>
    );
}
