import type { TFunction } from 'i18next';

import { useMemo, useState } from 'react';
import { useTranslation } from 'react-i18next';
import { useNavigate } from 'react-router-dom';
import { z } from 'zod';

import type { OAuthProvider } from '@/providers/user-provider';

import Github from '@/components/icons/github';
import Google from '@/components/icons/google';
import Logo from '@/components/icons/logo';
import { Button } from '@/components/ui/button';
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { FormSubmitButton } from '@/components/ui/form-submit-button';
import { Input } from '@/components/ui/input';
import { InputPassword } from '@/components/ui/input-password';
import { useAppForm } from '@/hooks/use-app-form';
import { routes } from '@/lib/routes';
import { useUser } from '@/providers/user-provider';

import { PasswordChangeForm } from './password-change-form';

const createLoginSchema = (t: TFunction<'auth'>) =>
    z.object({
        mail: z
            .string()
            .min(1, {
                message: t('login.validation.loginRequired'),
            })
            .refine(
                (value) =>
                    z.string().email().safeParse(value).success || ['admin', 'demo'].includes(value.toLowerCase()),
                {
                    message: t('login.validation.invalidLogin'),
                },
            ),
        password: z.string().min(1, {
            message: t('login.validation.passwordRequired'),
        }),
    });

interface AuthProviderAction {
    icon: React.ReactNode;
    id: OAuthProvider;
    labelKey: 'login.continueWithGithub' | 'login.continueWithGoogle';
}

type LoginFormValues = z.infer<ReturnType<typeof createLoginSchema>>;

const providerActions: AuthProviderAction[] = [
    {
        icon: <Google className="size-5" />,
        id: 'google',
        labelKey: 'login.continueWithGoogle',
    },
    {
        icon: <Github className="size-5" />,
        id: 'github',
        labelKey: 'login.continueWithGithub',
    },
];

interface LoginFormProps {
    providers: string[];
    returnUrl?: string;
}

function LoginForm({ providers, returnUrl = routes.newFlow }: LoginFormProps) {
    const { t } = useTranslation('auth');
    const formSchema = useMemo(() => createLoginSchema(t), [t]);
    const form = useAppForm<LoginFormValues>({
        defaultValues: {
            mail: '',
            password: '',
        },
        schema: formSchema,
    });
    const [isSubmitting, setIsSubmitting] = useState(false);
    const [error, setError] = useState<null | string>(null);
    const [passwordChangeRequired, setPasswordChangeRequired] = useState(false);
    const navigate = useNavigate();
    const { authInfo, isAuthenticated, login, loginWithOAuth, setAuth } = useUser();

    const handleSubmit = async (values: LoginFormValues) => {
        const errorMessage = t('errors.invalidCredentials');

        setError(null);

        try {
            const result = await login(values);

            if (!result.success) {
                setError(result.error || errorMessage);

                return;
            }

            if (result.passwordChangeRequired) {
                setPasswordChangeRequired(true);

                return;
            }

            navigate(returnUrl);
        } catch {
            setError(errorMessage);
        }
    };

    const handleProviderLogin = async (provider: OAuthProvider) => {
        setError(null);
        setIsSubmitting(true);

        try {
            const result = await loginWithOAuth(provider);

            if (!result.success) {
                setError(result.error || t('errors.authFailed'));

                return;
            }

            navigate(returnUrl);
        } catch (error) {
            setError(error instanceof Error ? error.message : t('errors.invalidCredentials'));
        } finally {
            setIsSubmitting(false);
        }
    };

    const handleSkipPasswordChange = () => {
        navigate(returnUrl);
    };

    const handlePasswordChangeSuccess = () => {
        if (authInfo?.user) {
            const updatedAuthData = {
                ...authInfo,
                user: {
                    ...authInfo.user,
                    password_change_required: false,
                },
            };

            setAuth(updatedAuthData);
            navigate(returnUrl);
        }
    };

    // If the session expired and user refreshed the page, the old authInfo may still
    // be in memory (race condition between clearAuth() and navigate()), but we must
    // NOT show the password change form because:
    //   1. The API endpoint /user/password requires authentication (returns 403 if not)
    //   2. The user must first re-login to establish a new valid session
    // Also check authInfo directly to handle page refresh scenarios where passwordChangeRequired
    // local state is lost but authInfo.user.password_change_required is still true.
    const shouldShowPasswordChange =
        (passwordChangeRequired || authInfo?.user?.password_change_required) &&
        authInfo?.user?.type === 'local' &&
        isAuthenticated();

    if (shouldShowPasswordChange) {
        return (
            <div className="mx-auto flex w-[350px] flex-col gap-6">
                <h1 className="text-center text-3xl font-bold">{t('passwordChange.title')}</h1>
                <p className="text-muted-foreground text-center text-sm">{t('passwordChange.requiredDescription')}</p>
                <PasswordChangeForm
                    layout="vertical"
                    onSkip={handleSkipPasswordChange}
                    onSuccess={handlePasswordChangeSuccess}
                />
            </div>
        );
    }

    return (
        <Form {...form}>
            <form
                className="mx-auto grid w-[350px] gap-8"
                noValidate
                onSubmit={form.handleSubmit(handleSubmit)}
            >
                <h1>
                    <Logo className="text-foreground mx-auto h-10 w-auto" />
                </h1>

                {providers?.length > 0 && (
                    <>
                        <div className="flex flex-col gap-4">
                            {providerActions
                                .filter((provider) => providers.includes(provider.id))
                                .map((provider) => (
                                    <Button
                                        disabled={isSubmitting || form.formState.isSubmitting}
                                        key={provider.id}
                                        onClick={() => handleProviderLogin(provider.id)}
                                        type="button"
                                        variant="secondary"
                                    >
                                        {provider.icon}
                                        {t(provider.labelKey)}
                                    </Button>
                                ))}
                        </div>

                        <div className="relative -mb-4">
                            <div className="absolute inset-0 flex items-center">
                                <div className="w-full border-t border-gray-300" />
                            </div>
                            <div className="relative flex justify-center text-sm">
                                <span className="bg-background px-2">{t('login.or')}</span>
                            </div>
                        </div>
                    </>
                )}

                <div className="flex flex-col gap-4">
                    <FormField
                        control={form.control}
                        name="mail"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>{t('login.loginLabel')}</FormLabel>
                                <FormControl>
                                    <Input
                                        {...field}
                                        autoFocus
                                        placeholder={t('login.loginPlaceholder')}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormField
                        control={form.control}
                        name="password"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>{t('login.passwordLabel')}</FormLabel>
                                <FormControl>
                                    <InputPassword
                                        {...field}
                                        placeholder={t('login.passwordPlaceholder')}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <FormSubmitButton className="w-full">
                        <span>{t('login.signIn')}</span>
                    </FormSubmitButton>

                    {error && <FormMessage>{error}</FormMessage>}
                </div>
            </form>
        </Form>
    );
}

export default LoginForm;
