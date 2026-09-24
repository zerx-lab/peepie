import type { Editor } from '@tiptap/react';
import type { LucideIcon } from 'lucide-react';

import { Check, ChevronDown, List, ListOrdered, ListTodo } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { Button } from '@/components/ui/button';
import {
    DropdownMenu,
    DropdownMenuContent,
    DropdownMenuItem,
    DropdownMenuTrigger,
} from '@/components/ui/dropdown-menu';
import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { cn } from '@/lib/utils';

import { returnFocusToEditor } from './markdown-editor-focus';

export type ListType = 'bullet' | 'ordered' | 'task';

interface ListOption {
    icon: LucideIcon;
    labelKey: 'list.bullet' | 'list.ordered' | 'list.task';
    value: ListType;
}

const OPTIONS: ListOption[] = [
    { icon: List, labelKey: 'list.bullet', value: 'bullet' },
    { icon: ListOrdered, labelKey: 'list.ordered', value: 'ordered' },
    { icon: ListTodo, labelKey: 'list.task', value: 'task' },
];

interface ListMenuProps {
    activeType: ListType | null;
    disabled?: boolean;
    editor: Editor;
    isInTableCell: boolean;
}

export function ListMenu({ activeType, disabled, editor, isInTableCell }: ListMenuProps) {
    const active = OPTIONS.find((option) => option.value === activeType);
    // The bullet-list icon doubles as the resting affordance, so the active background — not the glyph — is what
    // distinguishes "in a bullet list" from "no list".
    const TriggerIcon = active?.icon ?? List;
    const { t } = useTranslation('editor');

    return (
        <DropdownMenu>
            <Tooltip>
                <TooltipTrigger asChild>
                    <DropdownMenuTrigger asChild>
                        <Button
                            aria-label={t('list.triggerLabel', { list: active ? t(active.labelKey) : t('list.none') })}
                            className={cn('gap-0.5 px-1.5', active && 'bg-accent text-accent-foreground')}
                            data-toolbar-item=""
                            disabled={disabled}
                            size="sm"
                            type="button"
                            variant="ghost"
                        >
                            <TriggerIcon />
                            <ChevronDown className="size-4 opacity-50" />
                        </Button>
                    </DropdownMenuTrigger>
                </TooltipTrigger>
                <TooltipContent>{t('list.tooltip')}</TooltipContent>
            </Tooltip>
            <DropdownMenuContent
                align="start"
                className="min-w-[160px]"
                onCloseAutoFocus={returnFocusToEditor(editor)}
            >
                <ListMenuItems
                    activeType={activeType}
                    editor={editor}
                    isInTableCell={isInTableCell}
                />
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

function ListMenuItems({
    activeType,
    editor,
    isInTableCell,
}: {
    activeType: ListType | null;
    editor: Editor;
    isInTableCell: boolean;
}) {
    const { t } = useTranslation('editor');

    const applyOption = (value: ListType) => {
        const chain = editor.chain().focus();

        if (value === 'bullet') {
            chain.toggleBulletList().run();
        } else if (value === 'ordered') {
            chain.toggleOrderedList().run();
        } else {
            chain.toggleTaskList().run();
        }
    };

    return OPTIONS.map((option) => (
        <DropdownMenuItem
            aria-checked={activeType === option.value}
            disabled={isInTableCell}
            key={option.value}
            onSelect={() => applyOption(option.value)}
            role="menuitemradio"
        >
            <option.icon className="text-muted-foreground size-4 shrink-0" />
            <span>{t(option.labelKey)}</span>
            {activeType === option.value ? <Check className="ml-auto size-4 shrink-0" /> : null}
        </DropdownMenuItem>
    ));
}
