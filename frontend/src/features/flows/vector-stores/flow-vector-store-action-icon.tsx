import type { LucideIcon } from 'lucide-react';

import { HardDrive, HardDriveDownload, HardDriveUpload } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import { Tooltip, TooltipContent, TooltipTrigger } from '@/components/ui/tooltip';
import { VectorStoreAction } from '@/graphql/types';
import { cn } from '@/lib/utils';

interface FlowVectorStoreActionIconProps {
    action?: VectorStoreAction;
    className?: string;
    tooltip?: string;
}

const icons: Record<VectorStoreAction, LucideIcon> = {
    [VectorStoreAction.Retrieve]: HardDriveUpload,
    [VectorStoreAction.Store]: HardDriveDownload,
};
const defaultIcon = HardDrive;

const actionLabelKeys = {
    [VectorStoreAction.Retrieve]: 'vectorStores.actions.retrieve',
    [VectorStoreAction.Store]: 'vectorStores.actions.store',
} as const satisfies Record<VectorStoreAction, string>;

function FlowVectorStoreActionIcon({ action, className, tooltip }: FlowVectorStoreActionIconProps) {
    const { t } = useTranslation('flowDetails');
    const Icon = action ? icons[action] || defaultIcon : defaultIcon;
    const labelKey = action ? actionLabelKeys[action] : undefined;
    const tooltipText = tooltip ?? (labelKey ? t(labelKey) : action);
    const iconElement = <Icon className={cn('size-3 shrink-0', tooltipText && 'cursor-pointer', className)} />;

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

export default FlowVectorStoreActionIcon;
