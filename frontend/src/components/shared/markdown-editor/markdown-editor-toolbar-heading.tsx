import type { Editor } from '@tiptap/react';
import type { LucideIcon } from 'lucide-react';

import { Check, ChevronDown, Heading1, Heading2, Heading3, Heading4, Heading5, Heading6, Type } from 'lucide-react';
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

export type HeadingLevel = 1 | 2 | 3 | 4 | 5 | 6;

type HeadingLabelKey =
    | 'heading.h1'
    | 'heading.h2'
    | 'heading.h3'
    | 'heading.h4'
    | 'heading.h5'
    | 'heading.h6'
    | 'heading.text';

interface HeadingOption {
    icon: LucideIcon;
    // Optical trim for the dropdown list ONLY (icons sit next to each other there): lucide's Type glyph is drawn
    // taller (16u) than the Heading glyphs (12u), so at the same 16px box it reads bigger — scale it down to match
    // without shrinking the box (which would misalign labels). The trigger shows one icon alone, so it stays full size.
    iconClassName?: string;
    labelKey: HeadingLabelKey;
    value: 'paragraph' | HeadingLevel;
}

const OPTIONS: HeadingOption[] = [
    { icon: Heading1, labelKey: 'heading.h1', value: 1 },
    { icon: Heading2, labelKey: 'heading.h2', value: 2 },
    { icon: Heading3, labelKey: 'heading.h3', value: 3 },
    { icon: Heading4, labelKey: 'heading.h4', value: 4 },
    { icon: Heading5, labelKey: 'heading.h5', value: 5 },
    { icon: Heading6, labelKey: 'heading.h6', value: 6 },
    { icon: Type, iconClassName: 'scale-[0.75]', labelKey: 'heading.text', value: 'paragraph' },
];

interface HeadingMenuProps {
    // 0 = paragraph / any non-heading block; 1-6 = the active heading level.
    activeLevel: 0 | HeadingLevel;
    disabled?: boolean;
    editor: Editor;
    isInTableCell: boolean;
}

const isSelectedIn = (activeLevel: 0 | HeadingLevel, value: HeadingOption['value']) =>
    value === 'paragraph' ? activeLevel === 0 : value === activeLevel;

export function HeadingMenu({ activeLevel, disabled, editor, isInTableCell }: HeadingMenuProps) {
    const active = OPTIONS.find((option) => isSelectedIn(activeLevel, option.value)) ?? OPTIONS[0];
    const ActiveIcon = active?.icon ?? Type;
    const { t } = useTranslation('editor');

    return (
        <DropdownMenu>
            <Tooltip>
                <TooltipTrigger asChild>
                    <DropdownMenuTrigger asChild>
                        <Button
                            aria-label={t('heading.triggerLabel', { style: t(active?.labelKey ?? 'heading.text') })}
                            className="gap-0.5 px-1.5"
                            data-toolbar-item=""
                            disabled={disabled}
                            size="sm"
                            type="button"
                            variant="ghost"
                        >
                            <ActiveIcon />
                            <ChevronDown className="size-4 opacity-50" />
                        </Button>
                    </DropdownMenuTrigger>
                </TooltipTrigger>
                <TooltipContent>{t('heading.tooltip')}</TooltipContent>
            </Tooltip>
            <DropdownMenuContent
                align="start"
                className="min-w-[140px]"
                onCloseAutoFocus={returnFocusToEditor(editor)}
            >
                <HeadingMenuItems
                    activeLevel={activeLevel}
                    editor={editor}
                    isInTableCell={isInTableCell}
                />
            </DropdownMenuContent>
        </DropdownMenu>
    );
}

function HeadingMenuItems({
    activeLevel,
    editor,
    isInTableCell,
}: {
    activeLevel: 0 | HeadingLevel;
    editor: Editor;
    isInTableCell: boolean;
}) {
    const { t } = useTranslation('editor');

    const applyOption = (value: HeadingOption['value']) => {
        if (value === 'paragraph') {
            editor.chain().focus().setParagraph().run();

            return;
        }

        editor.chain().focus().toggleHeading({ level: value }).run();
    };

    return OPTIONS.map((option) => (
        <DropdownMenuItem
            aria-checked={isSelectedIn(activeLevel, option.value)}
            disabled={isInTableCell && option.value !== 'paragraph'}
            key={option.value}
            onSelect={() => applyOption(option.value)}
            role="menuitemradio"
        >
            <option.icon className={cn('text-muted-foreground size-4 shrink-0', option.iconClassName)} />
            <span>{t(option.labelKey)}</span>
            {isSelectedIn(activeLevel, option.value) ? <Check className="ml-auto size-4 shrink-0" /> : null}
        </DropdownMenuItem>
    ));
}
