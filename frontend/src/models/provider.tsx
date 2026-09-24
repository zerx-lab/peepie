import { ProviderType } from '@/graphql/types';
import { translateRuntimeKey } from '@/i18n';

export interface Provider {
    name: string;
    type: ProviderType;
}

export const getProviderDisplayName = (provider: Provider): string => {
    return provider.name;
};

export const isProviderValid = (provider: Provider, providers: Provider[]): boolean => {
    return providers.some((p) => p.name === provider.name && p.type === provider.type);
};

export const findProvider = (provider: Provider, providers: Provider[]): Provider | undefined => {
    return providers.find((p) => p.name === provider.name && p.type === provider.type);
};

export const findProviderByName = (providerName: string, providers: Provider[]): Provider | undefined => {
    return providers.find((provider) => provider.name === providerName);
};

export const sortProviders = (providers: Provider[]): Provider[] => {
    return [...providers].sort((a, b) => a.name.localeCompare(b.name));
};

const toTitleCase = (key: string): string =>
    key.replaceAll(/([A-Z])/g, ' $1').replace(/^./, (item) => item.toUpperCase());

/** Display name for an agent config key from the API (`primaryAgent`); unknown keys fall back to title case. */
export const getAgentTypeDisplayName = (agentKey: string): string =>
    translateRuntimeKey('providers', `agentTypes.${agentKey}`) ?? toTitleCase(agentKey);

/** Display name for an agent config field key (`maxTokens`); unknown keys fall back to title case. */
export const getAgentFieldDisplayName = (fieldKey: string): string =>
    translateRuntimeKey('providers', `fieldNames.${fieldKey}`) ?? toTitleCase(fieldKey);
