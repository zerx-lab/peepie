import type { LucideIcon } from 'lucide-react';

import {
    Bot,
    Brain,
    Code2,
    FileText,
    HardDrive,
    HardDriveDownload,
    HelpCircle,
    LayoutList,
    MessagesSquare,
    RefreshCw,
    Search,
    Settings,
    Sigma,
    Skull,
    Wrench,
} from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { AgentType } from '@/graphql/types';
import { cn } from '@/lib/utils';
import { formatName } from '@/lib/utils/format';

interface FlowAgentIconProps {
    className?: string;
    tooltip?: string;
    type?: AgentType;
}

const icons: Record<AgentType, LucideIcon> = {
    [AgentType.Adviser]: HelpCircle,
    [AgentType.Assistant]: Bot,
    [AgentType.Coder]: Code2,
    [AgentType.Enricher]: HardDriveDownload,
    [AgentType.Generator]: LayoutList,
    [AgentType.Installer]: Settings,
    [AgentType.Memorist]: HardDrive,
    [AgentType.Pentester]: Skull,
    [AgentType.PrimaryAgent]: Brain,
    [AgentType.Refiner]: RefreshCw,
    [AgentType.Reflector]: MessagesSquare,
    [AgentType.Reporter]: FileText,
    [AgentType.Searcher]: Search,
    [AgentType.Summarizer]: Sigma,
    [AgentType.ToolCallFixer]: Wrench,
};
const defaultIcon = HelpCircle;

export const agentTypeLabelKeys = {
    [AgentType.Adviser]: 'agentTypes.adviser',
    [AgentType.Assistant]: 'agentTypes.assistant',
    [AgentType.Coder]: 'agentTypes.coder',
    [AgentType.Enricher]: 'agentTypes.enricher',
    [AgentType.Generator]: 'agentTypes.generator',
    [AgentType.Installer]: 'agentTypes.installer',
    [AgentType.Memorist]: 'agentTypes.memorist',
    [AgentType.Pentester]: 'agentTypes.pentester',
    [AgentType.PrimaryAgent]: 'agentTypes.primaryAgent',
    [AgentType.Refiner]: 'agentTypes.refiner',
    [AgentType.Reflector]: 'agentTypes.reflector',
    [AgentType.Reporter]: 'agentTypes.reporter',
    [AgentType.Searcher]: 'agentTypes.searcher',
    [AgentType.Summarizer]: 'agentTypes.summarizer',
    [AgentType.ToolCallFixer]: 'agentTypes.toolCallFixer',
} as const satisfies Record<AgentType, string>;

export function useAgentTypeLabel() {
    const { t } = useTranslation('flowDetails');

    return (type?: string) => {
        const key = type ? agentTypeLabelKeys[type as AgentType] : undefined;

        return key ? t(key) : formatName(type);
    };
}

function FlowAgentIcon({ className, tooltip, type }: FlowAgentIconProps) {
    const getAgentTypeLabel = useAgentTypeLabel();
    const Icon = type ? icons[type] || defaultIcon : defaultIcon;
    const iconElement = <Icon className={cn('size-3 shrink-0', className)} />;
    const tooltipText = tooltip ?? (type ? getAgentTypeLabel(type) : undefined);

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

export default FlowAgentIcon;
