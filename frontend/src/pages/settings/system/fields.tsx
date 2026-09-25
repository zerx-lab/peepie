import type { ReactNode } from 'react';

import { Loader2, RotateCcw, Save, Undo2, X } from 'lucide-react';
import { useId } from 'react';
import { useTranslation } from 'react-i18next';

import type { BadgeVariant } from '@/components/ui/badge';

import { Badge } from '@/components/ui/badge';
import { Button } from '@/components/ui/button';
import { Input } from '@/components/ui/input';
import { InputPassword } from '@/components/ui/input-password';
import { Label } from '@/components/ui/label';
import { Switch } from '@/components/ui/switch';
import { cn } from '@/lib/utils';

import type { SecretEdit } from './secrets';

// Union of the LLM provider and search engine statuses.
export type ItemStatus = 'active' | 'disabled' | 'enabled' | 'error' | 'inactive' | 'unconfigured';

const statusDotClassNames: Record<ItemStatus, string> = {
    active: 'bg-green-500',
    disabled: 'bg-muted-foreground/40',
    enabled: 'bg-green-500',
    error: 'bg-destructive',
    inactive: 'bg-yellow-500',
    unconfigured: 'bg-muted-foreground/40',
};

const statusBadgeVariants: Record<ItemStatus, BadgeVariant> = {
    active: 'green',
    disabled: 'outline',
    enabled: 'green',
    error: 'red',
    inactive: 'yellow',
    unconfigured: 'outline',
};

export function FieldHint({ children }: { children: ReactNode }) {
    return <p className="text-muted-foreground text-xs">{children}</p>;
}

export function NumberField({
    hint,
    label,
    onChange,
    value,
}: {
    hint?: string;
    label: string;
    onChange: (value: number) => void;
    value: number;
}) {
    const id = useId();

    return (
        <div className="flex flex-col gap-1.5">
            <Label htmlFor={id}>{label}</Label>
            <Input
                id={id}
                min={0}
                onChange={(e) => onChange(Number(e.target.value) || 0)}
                type="number"
                value={value}
            />
            {hint && <FieldHint>{hint}</FieldHint>}
        </div>
    );
}

// Write-only secret input. `isSet` is the server-side state; the pending edit is
// either a typed replacement or an explicit clear (undoable until saved).
export function SecretField({
    edit,
    isSet,
    label,
    onChange,
}: {
    edit: SecretEdit | undefined;
    isSet: boolean;
    label: string;
    onChange: (edit: SecretEdit | undefined) => void;
}) {
    const { t } = useTranslation('settings');
    const id = useId();
    const isCleared = edit?.action === 'clear';
    const value = edit?.action === 'set' ? edit.value : '';

    let placeholder = t('system.secret.empty');

    if (isCleared) {
        placeholder = t('system.secret.willClear');
    } else if (isSet) {
        placeholder = t('system.secret.configured');
    }

    return (
        <div className="flex flex-col gap-1.5">
            <Label htmlFor={id}>{label}</Label>
            <div className="flex items-center gap-2">
                <div className="min-w-0 flex-1">
                    <InputPassword
                        autoComplete="off"
                        className={cn(isCleared && 'placeholder:text-destructive')}
                        id={id}
                        onChange={(e) =>
                            onChange(e.target.value ? { action: 'set', value: e.target.value } : undefined)
                        }
                        placeholder={placeholder}
                        value={value}
                    />
                </div>
                {isCleared ? (
                    <Button
                        onClick={() => onChange(undefined)}
                        type="button"
                        variant="outline"
                    >
                        <Undo2 />
                        {t('system.secret.undoClear')}
                    </Button>
                ) : (
                    isSet && (
                        <Button
                            onClick={() => onChange({ action: 'clear' })}
                            type="button"
                            variant="outline"
                        >
                            <X />
                            {t('system.secret.clear')}
                        </Button>
                    )
                )}
            </div>
        </div>
    );
}

export function SectionSaveBar({
    isSaving,
    onDiscard,
    onSave,
}: {
    isSaving: boolean;
    onDiscard: () => void;
    onSave: () => void;
}) {
    const { t } = useTranslation('settings');

    return (
        <div className="bg-background/95 sticky bottom-4 z-10 flex flex-wrap items-center justify-between gap-3 rounded-lg border p-3 shadow-lg backdrop-blur">
            <span className="text-sm font-medium">{t('system.unsavedChanges')}</span>
            <div className="flex items-center gap-2">
                <Button
                    disabled={isSaving}
                    onClick={onDiscard}
                    variant="outline"
                >
                    <RotateCcw />
                    {t('system.discard')}
                </Button>
                <Button
                    disabled={isSaving}
                    onClick={onSave}
                >
                    {isSaving ? <Loader2 className="animate-spin" /> : <Save />}
                    {t('system.save')}
                </Button>
            </div>
        </div>
    );
}

export function StatusBadge({ status }: { status: ItemStatus }) {
    const { t } = useTranslation('settings');

    return <Badge variant={statusBadgeVariants[status]}>{t(`system.status.${status}`)}</Badge>;
}

export function StatusDot({ className, status }: { className?: string; status: ItemStatus }) {
    return (
        <span
            aria-hidden
            className={cn('inline-block size-2 shrink-0 rounded-full', statusDotClassNames[status], className)}
        />
    );
}

export function TextField({
    hint,
    label,
    list,
    onChange,
    placeholder,
    value,
}: {
    hint?: string;
    label: string;
    list?: string;
    onChange: (value: string) => void;
    placeholder?: string;
    value: string;
}) {
    const id = useId();

    return (
        <div className="flex flex-col gap-1.5">
            <Label htmlFor={id}>{label}</Label>
            <Input
                autoComplete="off"
                id={id}
                list={list}
                onChange={(e) => onChange(e.target.value)}
                placeholder={placeholder}
                value={value}
            />
            {hint && <FieldHint>{hint}</FieldHint>}
        </div>
    );
}

export function ToggleField({
    checked,
    hint,
    label,
    onChange,
}: {
    checked: boolean;
    hint?: string;
    label: string;
    onChange: (value: boolean) => void;
}) {
    const id = useId();

    return (
        <div className="flex items-center justify-between gap-4 rounded-md border p-3">
            <div className="flex flex-col gap-1">
                <Label
                    className="font-normal"
                    htmlFor={id}
                >
                    {label}
                </Label>
                {hint && <FieldHint>{hint}</FieldHint>}
            </div>
            <Switch
                checked={checked}
                id={id}
                onCheckedChange={onChange}
            />
        </div>
    );
}
