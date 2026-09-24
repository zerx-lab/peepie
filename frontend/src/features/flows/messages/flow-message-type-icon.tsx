import type { LucideIcon } from 'lucide-react';

import {
    BotMessageSquare,
    Brain,
    CheckSquare,
    FileText,
    Globe,
    HelpCircle,
    MessageSquareReply,
    NotepadText,
    Search,
    Terminal,
    User as UserIcon,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { MessageLogType } from '@/graphql/types';
import { cn } from '@/lib/utils';

interface MessageTypeIconProps {
    className?: string;
    tooltip?: string;
    type?: MessageLogType;
}

const messageTypeIcons: Record<MessageLogType, LucideIcon> = {
    [MessageLogType.Advice]: BotMessageSquare,
    [MessageLogType.Answer]: MessageSquareReply,
    [MessageLogType.Ask]: HelpCircle,
    [MessageLogType.Browser]: Globe,
    [MessageLogType.Done]: CheckSquare,
    [MessageLogType.File]: FileText,
    [MessageLogType.Input]: UserIcon,
    [MessageLogType.Report]: NotepadText,
    [MessageLogType.Search]: Search,
    [MessageLogType.Terminal]: Terminal,
    [MessageLogType.Thoughts]: Brain,
};
const defaultIcon = Brain;

const messageTypeLabelKeys = {
    [MessageLogType.Advice]: 'messageTypes.advice',
    [MessageLogType.Answer]: 'messageTypes.answer',
    [MessageLogType.Ask]: 'messageTypes.ask',
    [MessageLogType.Browser]: 'messageTypes.browser',
    [MessageLogType.Done]: 'messageTypes.done',
    [MessageLogType.File]: 'messageTypes.file',
    [MessageLogType.Input]: 'messageTypes.input',
    [MessageLogType.Report]: 'messageTypes.report',
    [MessageLogType.Search]: 'messageTypes.search',
    [MessageLogType.Terminal]: 'messageTypes.terminal',
    [MessageLogType.Thoughts]: 'messageTypes.thoughts',
} as const satisfies Record<MessageLogType, string>;

function FlowMessageTypeIcon({ className, tooltip, type }: MessageTypeIconProps) {
    const { t } = useTranslation('flowDetails');
    const Icon = type ? messageTypeIcons[type] || defaultIcon : defaultIcon;
    const iconElement = <Icon className={cn('size-3 shrink-0', className)} />;
    const labelKey = type ? messageTypeLabelKeys[type] : undefined;
    const tooltipText = tooltip ?? (labelKey ? t(labelKey) : type);

    if (!tooltipText) {
        return iconElement;
    }

    return (
        <Tooltip>
            <TooltipTrigger asChild>{iconElement}</TooltipTrigger>
            <TooltipContent>{tooltipText}</TooltipContent>
        </Tooltip>
    );
}

export default FlowMessageTypeIcon;
