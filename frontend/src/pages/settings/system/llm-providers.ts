import type {
    BedrockProviderSettingsInput,
    CustomProviderSettingsInput,
    LlmProviderKeySettingsInput,
    LlmProviderSettingsFragmentFragment,
    LlmProviderSettingsInput,
    OllamaProviderSettingsInput,
} from '@/graphql/types';

import { type EditsWithSecrets, hasPendingEdits, resolveSecret } from './secrets';

// Display order of the built-in providers (matches the installer TUI).
export const llmProviderIds = [
    'openai',
    'anthropic',
    'gemini',
    'bedrock',
    'ollama',
    'custom',
    'deepseek',
    'glm',
    'kimi',
    'qwen',
    'minimax',
] as const;

export type BedrockEdits = EditsWithSecrets<
    BedrockProviderSettingsInput,
    'accessKeyId' | 'bearerToken' | 'secretAccessKey' | 'sessionToken'
>;
export type CustomEdits = EditsWithSecrets<CustomProviderSettingsInput, 'apiKey'>;

export type KeyProviderEdits = EditsWithSecrets<LlmProviderKeySettingsInput, 'apiKey'>;
export type KeyProviderId = Exclude<LlmProviderId, 'bedrock' | 'custom' | 'ollama'>;
export type LlmProviderEdits = { [K in KeyProviderId]?: KeyProviderEdits } & {
    bedrock?: BedrockEdits;
    custom?: CustomEdits;
    ollama?: OllamaEdits;
};
export type LlmProviderId = (typeof llmProviderIds)[number];

export type LlmProviderStatus = 'active' | 'error' | 'inactive' | 'unconfigured';

export type OllamaEdits = EditsWithSecrets<OllamaProviderSettingsInput, 'apiKey'>;

const secretFields: Record<LlmProviderId, readonly string[]> = {
    anthropic: ['apiKey'],
    bedrock: ['accessKeyId', 'bearerToken', 'secretAccessKey', 'sessionToken'],
    custom: ['apiKey'],
    deepseek: ['apiKey'],
    gemini: ['apiKey'],
    glm: ['apiKey'],
    kimi: ['apiKey'],
    minimax: ['apiKey'],
    ollama: ['apiKey'],
    openai: ['apiKey'],
    qwen: ['apiKey'],
};

export const isLlmProviderDirty = (
    server: LlmProviderSettingsFragmentFragment,
    edits: LlmProviderEdits,
    id: LlmProviderId,
): boolean => hasPendingEdits(server[id], edits[id], secretFields[id]);

// Mutation input containing ONLY providers with pending edits. Each one is a
// complete sub-input (server values overlaid with edits) because the backend
// replaces a provider's settings as a whole; secrets stay `undefined` (keep)
// unless typed or explicitly cleared.
export function buildLlmProviderInput(
    server: LlmProviderSettingsFragmentFragment,
    edits: LlmProviderEdits,
): LlmProviderSettingsInput {
    const input: LlmProviderSettingsInput = {};

    for (const id of llmProviderIds) {
        if (!isLlmProviderDirty(server, edits, id)) {
            continue;
        }

        switch (id) {
            case 'bedrock': {
                const current = server.bedrock;
                const edit = edits.bedrock ?? {};
                input.bedrock = {
                    accessKeyId: resolveSecret(edit.accessKeyId),
                    bearerToken: resolveSecret(edit.bearerToken),
                    defaultAuth: edit.defaultAuth ?? current.defaultAuth,
                    region: edit.region ?? current.region,
                    secretAccessKey: resolveSecret(edit.secretAccessKey),
                    serverUrl: edit.serverUrl ?? current.serverUrl,
                    sessionToken: resolveSecret(edit.sessionToken),
                };
                break;
            }

            case 'custom': {
                const current = server.custom;
                const edit = edits.custom ?? {};
                input.custom = {
                    apiKey: resolveSecret(edit.apiKey),
                    configPath: edit.configPath ?? current.configPath,
                    legacyReasoning: edit.legacyReasoning ?? current.legacyReasoning,
                    model: edit.model ?? current.model,
                    preserveReasoning: edit.preserveReasoning ?? current.preserveReasoning,
                    providerName: edit.providerName ?? current.providerName,
                    serverUrl: edit.serverUrl ?? current.serverUrl,
                };
                break;
            }

            case 'ollama': {
                const current = server.ollama;
                const edit = edits.ollama ?? {};
                input.ollama = {
                    apiKey: resolveSecret(edit.apiKey),
                    configPath: edit.configPath ?? current.configPath,
                    loadModelsEnabled: edit.loadModelsEnabled ?? current.loadModelsEnabled,
                    model: edit.model ?? current.model,
                    pullModelsEnabled: edit.pullModelsEnabled ?? current.pullModelsEnabled,
                    pullModelsTimeout: edit.pullModelsTimeout ?? current.pullModelsTimeout,
                    serverUrl: edit.serverUrl ?? current.serverUrl,
                };
                break;
            }

            default: {
                const current = server[id];
                const edit = edits[id] ?? {};
                input[id] = {
                    apiKey: resolveSecret(edit.apiKey),
                    // Only providers that support a model prefix report a non-null name.
                    providerName:
                        current.providerName === null ? undefined : (edit.providerName ?? current.providerName),
                    serverUrl: edit.serverUrl ?? current.serverUrl,
                };
            }
        }
    }

    return input;
}

export function getLlmProviderStatus(
    server: LlmProviderSettingsFragmentFragment,
    id: LlmProviderId,
): LlmProviderStatus {
    const provider = server[id];

    if (provider.active) {
        return 'active';
    }

    if (provider.error) {
        return 'error';
    }

    return isLlmProviderConfigured(server, id) ? 'inactive' : 'unconfigured';
}

// A provider is configured once it has credentials (or, for the self-hosted
// ones, an endpoint); only then can it be built.
export function isLlmProviderConfigured(server: LlmProviderSettingsFragmentFragment, id: LlmProviderId): boolean {
    switch (id) {
        case 'bedrock': {
            const { accessKeyIdSet, bearerTokenSet, defaultAuth, secretAccessKeySet } = server.bedrock;

            return defaultAuth || bearerTokenSet || (accessKeyIdSet && secretAccessKeySet);
        }

        case 'custom':
        case 'ollama':
            return server[id].serverUrl.trim() !== '';

        default:
            return server[id].apiKeySet;
    }
}
