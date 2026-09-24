import { Lock, Mail, User } from 'lucide-react';
import { useState } from 'react';
import { useTranslation } from 'react-i18next';

import { AppHeader, AppHeaderContent, AppHeaderTitle } from '@/components/layouts/app/app-header';
import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { EmailChangeForm } from '@/features/authentication/email-change-form';
import { NameChangeForm } from '@/features/authentication/name-change-form';
import { PasswordChangeForm } from '@/features/authentication/password-change-form';
import { formatLocalizedDate } from '@/i18n/format';
import { useUser } from '@/providers/user-provider';

type EditingSection = 'email' | 'name' | 'password';

const PROVIDER_LABELS: Record<string, string> = {
    github: 'GitHub',
    google: 'Google',
};

function SettingsAccount() {
    const { t } = useTranslation(['settings', 'common']);
    const { authInfo } = useUser();
    const user = authInfo?.user;
    const [editingSections, setEditingSections] = useState<Set<EditingSection>>(new Set());

    if (!user) {
        return null;
    }

    const isLocal = user.type === 'local';
    const startEditing = (section: EditingSection) => setEditingSections((prev) => new Set(prev).add(section));
    const stopEditing = (section: EditingSection) =>
        setEditingSections((prev) => {
            const next = new Set(prev);
            next.delete(section);

            return next;
        });

    const isEditingName = editingSections.has('name');
    const isEditingEmail = editingSections.has('email');
    const isEditingPassword = editingSections.has('password');
    const displayName = user.name?.trim() || user.mail;
    const initial = ([...(displayName || '?')][0] ?? '?').toUpperCase();
    const createdAt = user.created_at ? new Date(user.created_at) : null;
    const memberSince =
        createdAt && !Number.isNaN(createdAt.getTime()) ? formatLocalizedDate(createdAt, 'monthYear') : null;
    const accountLabel = isLocal
        ? t('account.localAccount')
        : user.provider
          ? (PROVIDER_LABELS[user.provider] ?? user.provider)
          : t('account.oauthAccount');

    return (
        <>
            <AppHeader>
                <AppHeaderContent>
                    <AppHeaderTitle icon={<User className="size-4 shrink-0" />}>{t('account.title')}</AppHeaderTitle>
                </AppHeaderContent>
            </AppHeader>
            <div className="mx-auto flex w-full max-w-2xl flex-col gap-4 p-4">
                <Card>
                    <CardHeader className="flex-row items-center gap-4">
                        <div className="bg-muted text-foreground flex size-12 shrink-0 items-center justify-center rounded-lg text-lg font-semibold">
                            {initial}
                        </div>
                        <div className="flex min-w-0 flex-1 flex-col gap-0.5">
                            <CardTitle className="truncate">{displayName}</CardTitle>
                            {memberSince && (
                                <CardDescription className="truncate">
                                    {t('account.memberSince', { date: memberSince })}
                                </CardDescription>
                            )}
                        </div>
                        <Badge
                            className="shrink-0"
                            variant="secondary"
                        >
                            {accountLabel}
                        </Badge>
                    </CardHeader>
                </Card>

                <Card>
                    <CardHeader className="flex-row items-start justify-between gap-4">
                        <div className="grid gap-1.5">
                            <CardTitle>{t('account.displayName.title')}</CardTitle>
                            <CardDescription>{t('account.displayName.description')}</CardDescription>
                        </div>
                        {!isEditingName && (
                            <Button
                                onClick={() => startEditing('name')}
                                size="sm"
                                variant="outline"
                            >
                                {t('account.change')}
                            </Button>
                        )}
                    </CardHeader>
                    <CardContent>
                        {isEditingName ? (
                            <NameChangeForm
                                onCancel={() => stopEditing('name')}
                                onSuccess={() => stopEditing('name')}
                            />
                        ) : (
                            <div className="text-muted-foreground flex items-center gap-2 text-sm">
                                <User className="size-4 shrink-0" />
                                <span className="truncate">{displayName}</span>
                            </div>
                        )}
                    </CardContent>
                </Card>

                <Card>
                    <CardHeader className="flex-row items-start justify-between gap-4">
                        <div className="grid gap-1.5">
                            <CardTitle>{t('account.email.title')}</CardTitle>
                            <CardDescription>
                                {isLocal
                                    ? t('account.email.descriptionLocal')
                                    : t('account.email.descriptionLinked', { account: accountLabel })}
                            </CardDescription>
                        </div>
                        {isLocal && !isEditingEmail && (
                            <Button
                                onClick={() => startEditing('email')}
                                size="sm"
                                variant="outline"
                            >
                                {t('account.change')}
                            </Button>
                        )}
                    </CardHeader>
                    <CardContent>
                        {isLocal && isEditingEmail ? (
                            <EmailChangeForm
                                onCancel={() => stopEditing('email')}
                                onSuccess={() => stopEditing('email')}
                            />
                        ) : (
                            <div className="text-muted-foreground flex items-center gap-2 text-sm">
                                <Mail className="size-4 shrink-0" />
                                <span className="truncate">{user.mail}</span>
                            </div>
                        )}
                    </CardContent>
                </Card>

                {isLocal && (
                    <Card>
                        <CardHeader className="flex-row items-start justify-between gap-4">
                            <div className="grid gap-1.5">
                                <CardTitle>{t('account.password.title')}</CardTitle>
                                <CardDescription>{t('account.password.description')}</CardDescription>
                            </div>
                            {!isEditingPassword && (
                                <Button
                                    onClick={() => startEditing('password')}
                                    size="sm"
                                    variant="outline"
                                >
                                    {t('account.change')}
                                </Button>
                            )}
                        </CardHeader>
                        <CardContent>
                            {isEditingPassword ? (
                                <PasswordChangeForm
                                    buttonSize="sm"
                                    onCancel={() => stopEditing('password')}
                                    onSuccess={() => stopEditing('password')}
                                />
                            ) : (
                                <div className="text-muted-foreground flex items-center gap-2 text-sm">
                                    <Lock className="size-4 shrink-0" />
                                    <span>••••••••••••</span>
                                </div>
                            )}
                        </CardContent>
                    </Card>
                )}
            </div>
        </>
    );
}

export default SettingsAccount;
