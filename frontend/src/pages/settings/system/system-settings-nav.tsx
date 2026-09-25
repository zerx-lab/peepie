import type { LucideIcon } from 'lucide-react';

import { useTranslation } from 'react-i18next';

import { cn } from '@/lib/utils';

import { type ItemStatus, StatusDot } from './fields';

export const systemSections = ['llm-providers', 'search-engines', 'execution'] as const;

export interface SystemNavCategory {
    hasUnsavedChanges: boolean;
    icon: LucideIcon;
    id: SystemSection;
    items: SystemNavItem[];
    label: string;
    summary: null | string;
}

export interface SystemNavItem {
    anchorId: string;
    id: string;
    label: string;
    status: ItemStatus;
}

export type SystemSection = (typeof systemSections)[number];

interface SystemSettingsNavProps {
    active: SystemSection;
    categories: SystemNavCategory[];
    onSelect: (section: SystemSection) => void;
}

export function SystemSettingsNav({ active, categories, onSelect }: SystemSettingsNavProps) {
    const { t } = useTranslation('settings');

    return (
        <>
            {/* Narrow screens: horizontal strip above the content. */}
            <nav
                aria-label={t('system.nav.label')}
                className="-mx-4 flex gap-2 overflow-x-auto px-4 pb-1 lg:hidden"
            >
                {categories.map(({ hasUnsavedChanges, icon: Icon, id, label, summary }) => (
                    <button
                        aria-current={id === active ? 'page' : undefined}
                        className={cn(
                            'flex shrink-0 items-center gap-2 rounded-md border px-3 py-2 text-left text-sm transition-colors',
                            id === active ? 'bg-accent text-accent-foreground border-primary/40' : 'hover:bg-accent/50',
                        )}
                        key={id}
                        onClick={() => onSelect(id)}
                        type="button"
                    >
                        <Icon className="size-4 shrink-0" />
                        <span className="font-medium">{label}</span>
                        {summary && <span className="text-muted-foreground text-xs">{summary}</span>}
                        {hasUnsavedChanges && <span className="bg-primary size-1.5 rounded-full" />}
                    </button>
                ))}
            </nav>

            {/* Wide screens: sticky side menu doubling as a status overview. */}
            <nav
                aria-label={t('system.nav.label')}
                className="sticky top-16 hidden max-h-[calc(100svh-5rem)] w-56 shrink-0 self-start overflow-y-auto lg:block"
            >
                <ul className="flex flex-col gap-1">
                    {categories.map(({ hasUnsavedChanges, icon: Icon, id, items, label, summary }) => {
                        const isActive = id === active;

                        return (
                            <li key={id}>
                                <button
                                    aria-current={isActive ? 'page' : undefined}
                                    className={cn(
                                        'flex w-full items-start gap-2 rounded-md px-3 py-2 text-left transition-colors',
                                        isActive ? 'bg-accent text-accent-foreground' : 'hover:bg-accent/50',
                                    )}
                                    onClick={() => onSelect(id)}
                                    type="button"
                                >
                                    <Icon className="mt-0.5 size-4 shrink-0" />
                                    <span className="flex min-w-0 flex-1 flex-col">
                                        <span className="flex items-center gap-2 text-sm font-medium">
                                            {label}
                                            {hasUnsavedChanges && (
                                                <span
                                                    className="bg-primary size-1.5 rounded-full"
                                                    title={t('system.unsavedChanges')}
                                                />
                                            )}
                                        </span>
                                        {summary && <span className="text-muted-foreground text-xs">{summary}</span>}
                                    </span>
                                </button>
                                {isActive && items.length > 0 && (
                                    <ul className="mt-1 mb-2 ml-5 flex flex-col border-l pl-2">
                                        {items.map((item) => (
                                            <li key={item.id}>
                                                <button
                                                    className="text-muted-foreground hover:text-foreground hover:bg-accent/50 flex w-full items-center gap-2 rounded-md px-2 py-1 text-left text-sm transition-colors"
                                                    onClick={() =>
                                                        document
                                                            .getElementById(item.anchorId)
                                                            ?.scrollIntoView({ behavior: 'smooth', block: 'start' })
                                                    }
                                                    title={t(`system.status.${item.status}`)}
                                                    type="button"
                                                >
                                                    <StatusDot status={item.status} />
                                                    <span className="truncate">{item.label}</span>
                                                </button>
                                            </li>
                                        ))}
                                    </ul>
                                )}
                            </li>
                        );
                    })}
                </ul>
            </nav>
        </>
    );
}
