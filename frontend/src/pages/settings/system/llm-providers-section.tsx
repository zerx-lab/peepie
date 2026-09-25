import type { TFunction } from 'i18next';
import type { Dispatch, ReactNode, SetStateAction } from 'react';

import { AlertCircle } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import type { LlmProviderSettingsFragmentFragment } from '@/graphql/types';

import { providerIcons } from '@/components/icons/provider-icon';
import { Alert, AlertDescription, AlertTitle } from '@/components/ui/alert';
import { Badge } from '@/components/ui/badge';
import { Card, CardContent, CardDescription, CardHeader, CardTitle } from '@/components/ui/card';
import { ProviderType } from '@/graphql/types';
import { cn } from '@/lib/utils';

import type { SecretEdit } from './secrets';

import { NumberField, SecretField, StatusBadge, TextField, ToggleField } from './fields';
import {
    getLlmProviderStatus,
    isLlmProviderDirty,
    type KeyProviderId,
    type LlmProviderEdits,
    type LlmProviderId,
    llmProviderIds,
} from './llm-providers';

// Brand names stay literal; `null` marks the generic provider whose label is translated.
const providerMeta: Record<LlmProviderId, { brand: null | string; type: ProviderType }> = {
    anthropic: { brand: 'Anthropic', type: ProviderType.Anthropic },
    bedrock: { brand: 'AWS Bedrock', type: ProviderType.Bedrock },
    custom: { brand: null, type: ProviderType.Custom },
    deepseek: { brand: 'DeepSeek', type: ProviderType.Deepseek },
    gemini: { brand: 'Gemini', type: ProviderType.Gemini },
    glm: { brand: 'GLM', type: ProviderType.Glm },
    kimi: { brand: 'Kimi', type: ProviderType.Kimi },
    minimax: { brand: 'MiniMax', type: ProviderType.Minimax },
    ollama: { brand: 'Ollama', type: ProviderType.Ollama },
    openai: { brand: 'OpenAI', type: ProviderType.Openai },
    qwen: { brand: 'Qwen', type: ProviderType.Qwen },
};

const bedrockSecrets = [
    { field: 'bearerToken', isSetField: 'bearerTokenSet' },
    { field: 'accessKeyId', isSetField: 'accessKeyIdSet' },
    { field: 'secretAccessKey', isSetField: 'secretAccessKeySet' },
    { field: 'sessionToken', isSetField: 'sessionTokenSet' },
] as const;

const configPathsListId = 'llm-provider-config-paths';

export const llmProviderAnchorId = (id: LlmProviderId) => `llm-provider-${id}`;

export const getLlmProviderLabel = (id: LlmProviderId, t: TFunction<'settings'>): string =>
    providerMeta[id].brand ?? t('system.llm.customName');

interface LlmProvidersSectionProps {
    edits: LlmProviderEdits;
    onEditsChange: Dispatch<SetStateAction<LlmProviderEdits>>;
    server: LlmProviderSettingsFragmentFragment;
}

export function LlmProvidersSection({ edits, onEditsChange, server }: LlmProvidersSectionProps) {
    const { t } = useTranslation('settings');

    const update = <K extends LlmProviderId>(id: K, patch: NonNullable<LlmProviderEdits[K]>) =>
        onEditsChange((prev) => ({ ...prev, [id]: { ...prev[id], ...patch } }));

    const renderKeyProvider = (id: KeyProviderId) => {
        const current = server[id];
        const edit = edits[id] ?? {};

        return (
            <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                <SecretField
                    edit={edit.apiKey}
                    isSet={current.apiKeySet}
                    label={t('system.llm.fields.apiKey')}
                    onChange={(apiKey) => update(id, { apiKey })}
                />
                <TextField
                    hint={t('system.llm.fields.serverUrlHint')}
                    label={t('system.llm.fields.serverUrl')}
                    onChange={(serverUrl) => update(id, { serverUrl })}
                    value={edit.serverUrl ?? current.serverUrl}
                />
                {current.providerName !== null && (
                    <TextField
                        hint={t('system.llm.fields.providerNameHint')}
                        label={t('system.llm.fields.providerName')}
                        onChange={(providerName) => update(id, { providerName })}
                        value={edit.providerName ?? current.providerName}
                    />
                )}
            </div>
        );
    };

    const renderBedrock = () => {
        const current = server.bedrock;
        const edit = edits.bedrock ?? {};

        return (
            <>
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <TextField
                        label={t('system.llm.fields.region')}
                        onChange={(region) => update('bedrock', { region })}
                        value={edit.region ?? current.region}
                    />
                    <TextField
                        hint={t('system.llm.fields.serverUrlHint')}
                        label={t('system.llm.fields.serverUrl')}
                        onChange={(serverUrl) => update('bedrock', { serverUrl })}
                        value={edit.serverUrl ?? current.serverUrl}
                    />
                </div>
                <ToggleField
                    checked={edit.defaultAuth ?? current.defaultAuth}
                    hint={t('system.llm.fields.defaultAuthHint')}
                    label={t('system.llm.fields.defaultAuth')}
                    onChange={(defaultAuth) => update('bedrock', { defaultAuth })}
                />
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    {bedrockSecrets.map(({ field, isSetField }) => (
                        <SecretField
                            edit={edit[field]}
                            isSet={current[isSetField]}
                            key={field}
                            label={t(`system.llm.fields.${field}`)}
                            onChange={(value: SecretEdit | undefined) => update('bedrock', { [field]: value })}
                        />
                    ))}
                </div>
            </>
        );
    };

    const renderOllama = () => {
        const current = server.ollama;
        const edit = edits.ollama ?? {};

        return (
            <>
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <TextField
                        label={t('system.llm.fields.serverUrl')}
                        onChange={(serverUrl) => update('ollama', { serverUrl })}
                        value={edit.serverUrl ?? current.serverUrl}
                    />
                    <SecretField
                        edit={edit.apiKey}
                        isSet={current.apiKeySet}
                        label={t('system.llm.fields.apiKey')}
                        onChange={(apiKey) => update('ollama', { apiKey })}
                    />
                    <TextField
                        label={t('system.llm.fields.model')}
                        onChange={(model) => update('ollama', { model })}
                        value={edit.model ?? current.model}
                    />
                    <TextField
                        hint={t('system.llm.fields.configPathHint')}
                        label={t('system.llm.fields.configPath')}
                        list={configPathsListId}
                        onChange={(configPath) => update('ollama', { configPath })}
                        value={edit.configPath ?? current.configPath}
                    />
                </div>
                <ToggleField
                    checked={edit.pullModelsEnabled ?? current.pullModelsEnabled}
                    label={t('system.llm.fields.pullModelsEnabled')}
                    onChange={(pullModelsEnabled) => update('ollama', { pullModelsEnabled })}
                />
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <NumberField
                        label={t('system.llm.fields.pullModelsTimeout')}
                        onChange={(pullModelsTimeout) => update('ollama', { pullModelsTimeout })}
                        value={edit.pullModelsTimeout ?? current.pullModelsTimeout}
                    />
                </div>
                <ToggleField
                    checked={edit.loadModelsEnabled ?? current.loadModelsEnabled}
                    label={t('system.llm.fields.loadModelsEnabled')}
                    onChange={(loadModelsEnabled) => update('ollama', { loadModelsEnabled })}
                />
            </>
        );
    };

    const renderCustom = () => {
        const current = server.custom;
        const edit = edits.custom ?? {};

        return (
            <>
                <div className="grid grid-cols-1 gap-3 md:grid-cols-2">
                    <TextField
                        label={t('system.llm.fields.serverUrl')}
                        onChange={(serverUrl) => update('custom', { serverUrl })}
                        value={edit.serverUrl ?? current.serverUrl}
                    />
                    <SecretField
                        edit={edit.apiKey}
                        isSet={current.apiKeySet}
                        label={t('system.llm.fields.apiKey')}
                        onChange={(apiKey) => update('custom', { apiKey })}
                    />
                    <TextField
                        label={t('system.llm.fields.model')}
                        onChange={(model) => update('custom', { model })}
                        value={edit.model ?? current.model}
                    />
                    <TextField
                        hint={t('system.llm.fields.configPathHint')}
                        label={t('system.llm.fields.configPath')}
                        list={configPathsListId}
                        onChange={(configPath) => update('custom', { configPath })}
                        value={edit.configPath ?? current.configPath}
                    />
                    <TextField
                        hint={t('system.llm.fields.providerNameHint')}
                        label={t('system.llm.fields.providerName')}
                        onChange={(providerName) => update('custom', { providerName })}
                        value={edit.providerName ?? current.providerName}
                    />
                </div>
                <ToggleField
                    checked={edit.legacyReasoning ?? current.legacyReasoning}
                    hint={t('system.llm.fields.legacyReasoningHint')}
                    label={t('system.llm.fields.legacyReasoning')}
                    onChange={(legacyReasoning) => update('custom', { legacyReasoning })}
                />
                <ToggleField
                    checked={edit.preserveReasoning ?? current.preserveReasoning}
                    hint={t('system.llm.fields.preserveReasoningHint')}
                    label={t('system.llm.fields.preserveReasoning')}
                    onChange={(preserveReasoning) => update('custom', { preserveReasoning })}
                />
            </>
        );
    };

    const renderFields = (id: LlmProviderId): ReactNode => {
        switch (id) {
            case 'bedrock':
                return renderBedrock();
            case 'custom':
                return renderCustom();
            case 'ollama':
                return renderOllama();
            default:
                return renderKeyProvider(id);
        }
    };

    return (
        <>
            <datalist id={configPathsListId}>
                {server.configPaths.map((path) => (
                    <option
                        key={path}
                        value={path}
                    />
                ))}
            </datalist>
            {llmProviderIds.map((id) => {
                const { className: iconClassName, icon: Icon } = providerIcons[providerMeta[id].type];
                const error = server[id].error;

                return (
                    <Card
                        className="scroll-mt-16"
                        id={llmProviderAnchorId(id)}
                        key={id}
                    >
                        <CardHeader>
                            <div className="flex flex-wrap items-center gap-2">
                                <Icon className={cn('size-5 shrink-0', iconClassName)} />
                                <CardTitle>{getLlmProviderLabel(id, t)}</CardTitle>
                                <StatusBadge status={getLlmProviderStatus(server, id)} />
                                {isLlmProviderDirty(server, edits, id) && (
                                    <Badge variant="blue">{t('system.modified')}</Badge>
                                )}
                            </div>
                            <CardDescription>{t(`system.llm.descriptions.${id}`)}</CardDescription>
                        </CardHeader>
                        <CardContent className="flex flex-col gap-4">
                            {error && (
                                <Alert variant="destructive">
                                    <AlertCircle />
                                    <AlertTitle>{t('system.llm.buildFailed')}</AlertTitle>
                                    <AlertDescription className="break-words">{error}</AlertDescription>
                                </Alert>
                            )}
                            {renderFields(id)}
                        </CardContent>
                    </Card>
                );
            })}
        </>
    );
}
