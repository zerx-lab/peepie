import type { Dispatch, SetStateAction } from 'react';

import { useTranslation } from 'react-i18next';

import type { ExecutionSettingsFragmentFragment, ExecutionSettingsInput } from '@/graphql/types';

import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';

import { NumberField, TextField, ToggleField } from './fields';

export type ExecutionEdits = Partial<ExecutionSettingsInput>;

interface ExecutionSectionProps {
    edits: ExecutionEdits;
    onEditsChange: Dispatch<SetStateAction<ExecutionEdits>>;
    server: ExecutionSettingsFragmentFragment;
}

export function ExecutionSection({ edits, onEditsChange, server }: ExecutionSectionProps) {
    const { t } = useTranslation('settings');
    const form: ExecutionSettingsInput = { ...server, ...edits };

    const set = <K extends keyof ExecutionSettingsInput>(key: K, value: ExecutionSettingsInput[K]) =>
        onEditsChange((prev) => ({ ...prev, [key]: value }));

    return (
        <>
            <Card>
                <CardHeader>
                    <CardTitle>{t('system.execution.agentsTitle')}</CardTitle>
                    <CardDescription>{t('system.execution.description')}</CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col gap-3">
                    <ToggleField
                        checked={form.assistantUseAgents}
                        label={t('system.execution.assistantUseAgents')}
                        onChange={(v) => set('assistantUseAgents', v)}
                    />
                    <ToggleField
                        checked={form.agentPlanningStepEnabled}
                        label={t('system.execution.agentPlanningStepEnabled')}
                        onChange={(v) => set('agentPlanningStepEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <NumberField
                            label={t('system.execution.maxGeneralAgentToolCalls')}
                            onChange={(v) => set('maxGeneralAgentToolCalls', v)}
                            value={form.maxGeneralAgentToolCalls}
                        />
                        <NumberField
                            label={t('system.execution.maxLimitedAgentToolCalls')}
                            onChange={(v) => set('maxLimitedAgentToolCalls', v)}
                            value={form.maxLimitedAgentToolCalls}
                        />
                    </div>
                </CardContent>
            </Card>
            <Card>
                <CardHeader>
                    <CardTitle>{t('system.execution.imagesTitle')}</CardTitle>
                    <CardDescription>{t('system.execution.imagesDescription')}</CardDescription>
                </CardHeader>
                <CardContent className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <TextField
                        hint={t('system.execution.imageHint')}
                        label={t('system.execution.flowDockerImage')}
                        onChange={(v) => set('flowDockerImage', v)}
                        placeholder={t('system.execution.imagePlaceholder')}
                        value={form.flowDockerImage}
                    />
                    <TextField
                        hint={t('system.execution.imageHint')}
                        label={t('system.execution.assistantDockerImage')}
                        onChange={(v) => set('assistantDockerImage', v)}
                        placeholder={t('system.execution.imagePlaceholder')}
                        value={form.assistantDockerImage}
                    />
                </CardContent>
            </Card>
            <Card>
                <CardHeader>
                    <CardTitle>{t('system.execution.monitorTitle')}</CardTitle>
                    <CardDescription>{t('system.execution.monitorDescription')}</CardDescription>
                </CardHeader>
                <CardContent className="flex flex-col gap-3">
                    <ToggleField
                        checked={form.executionMonitorEnabled}
                        label={t('system.execution.monitorEnabled')}
                        onChange={(v) => set('executionMonitorEnabled', v)}
                    />
                    <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                        <NumberField
                            label={t('system.execution.sameToolLimit')}
                            onChange={(v) => set('executionMonitorSameToolLimit', v)}
                            value={form.executionMonitorSameToolLimit}
                        />
                        <NumberField
                            label={t('system.execution.totalToolLimit')}
                            onChange={(v) => set('executionMonitorTotalToolLimit', v)}
                            value={form.executionMonitorTotalToolLimit}
                        />
                    </div>
                </CardContent>
            </Card>
        </>
    );
}
