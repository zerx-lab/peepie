import { describe, expect, it } from 'vitest';

import type { LlmProviderKeySettingsFragmentFragment, LlmProviderSettingsFragmentFragment } from '@/graphql/types';

import { buildLlmProviderInput, getLlmProviderStatus, isLlmProviderConfigured } from './llm-providers';

const keyProvider = (overrides: Partial<LlmProviderKeySettingsFragmentFragment> = {}) => ({
    active: false,
    apiKeySet: false,
    error: null,
    providerName: null,
    serverUrl: '',
    ...overrides,
});

const makeServer = (
    overrides: Partial<LlmProviderSettingsFragmentFragment> = {},
): LlmProviderSettingsFragmentFragment => ({
    anthropic: keyProvider(),
    bedrock: {
        accessKeyIdSet: false,
        active: false,
        bearerTokenSet: false,
        defaultAuth: false,
        error: null,
        region: 'us-east-1',
        secretAccessKeySet: false,
        serverUrl: '',
        sessionTokenSet: false,
    },
    configPaths: [],
    custom: {
        active: false,
        apiKeySet: false,
        configPath: '',
        error: null,
        legacyReasoning: false,
        model: '',
        preserveReasoning: false,
        providerName: '',
        serverUrl: '',
    },
    deepseek: keyProvider({ providerName: '' }),
    gemini: keyProvider(),
    glm: keyProvider({ providerName: '' }),
    kimi: keyProvider({ providerName: '' }),
    minimax: keyProvider({ providerName: '' }),
    ollama: {
        active: false,
        apiKeySet: false,
        configPath: '',
        error: null,
        loadModelsEnabled: false,
        model: '',
        pullModelsEnabled: false,
        pullModelsTimeout: 600,
        serverUrl: '',
    },
    openai: keyProvider({ active: true, apiKeySet: true, serverUrl: 'https://api.openai.com/v1' }),
    qwen: keyProvider({ providerName: '' }),
    ...overrides,
});

describe('buildLlmProviderInput', () => {
    it('sends nothing without edits', () => {
        expect(buildLlmProviderInput(makeServer(), {})).toEqual({});
    });

    it('skips providers whose edits match the server or only hold an empty typed secret', () => {
        const input = buildLlmProviderInput(makeServer(), {
            anthropic: { apiKey: { action: 'set', value: '' } },
            openai: { serverUrl: 'https://api.openai.com/v1' },
        });

        expect(input).toEqual({});
    });

    it('sends only dirty providers as full sub-inputs, keeping untouched secrets undefined', () => {
        const input = buildLlmProviderInput(makeServer(), {
            deepseek: { providerName: 'deepseek' },
            openai: { serverUrl: 'https://proxy.local/v1' },
        });

        expect(Object.keys(input).sort()).toEqual(['deepseek', 'openai']);
        expect(input.openai).toEqual({
            apiKey: undefined,
            providerName: undefined,
            serverUrl: 'https://proxy.local/v1',
        });
        expect(input.deepseek).toEqual({ apiKey: undefined, providerName: 'deepseek', serverUrl: '' });
    });

    it('sends "" for a cleared secret and the value for a typed one', () => {
        const input = buildLlmProviderInput(makeServer(), {
            anthropic: { apiKey: { action: 'set', value: 'sk-ant' } },
            openai: { apiKey: { action: 'clear' } },
        });

        expect(input.openai?.apiKey).toBe('');
        expect(input.anthropic?.apiKey).toBe('sk-ant');
    });

    it('builds bedrock from server values with per-secret resolution', () => {
        const input = buildLlmProviderInput(makeServer(), {
            bedrock: {
                accessKeyId: { action: 'set', value: 'AKIA' },
                region: 'eu-west-1',
                sessionToken: { action: 'clear' },
            },
        });

        expect(input).toEqual({
            bedrock: {
                accessKeyId: 'AKIA',
                bearerToken: undefined,
                defaultAuth: false,
                region: 'eu-west-1',
                secretAccessKey: undefined,
                serverUrl: '',
                sessionToken: '',
            },
        });
    });

    it('builds ollama and custom with non-secret fields from the server', () => {
        const input = buildLlmProviderInput(makeServer(), {
            custom: { preserveReasoning: true },
            ollama: { serverUrl: 'http://ollama:11434' },
        });

        expect(input.ollama).toEqual({
            apiKey: undefined,
            configPath: '',
            loadModelsEnabled: false,
            model: '',
            pullModelsEnabled: false,
            pullModelsTimeout: 600,
            serverUrl: 'http://ollama:11434',
        });
        expect(input.custom).toMatchObject({ apiKey: undefined, preserveReasoning: true, providerName: '' });
    });
});

describe('LLM provider status', () => {
    it('reports active providers as active', () => {
        expect(getLlmProviderStatus(makeServer(), 'openai')).toBe('active');
    });

    it('reports a build error', () => {
        const server = makeServer({ anthropic: keyProvider({ apiKeySet: true, error: 'invalid key' }) });

        expect(getLlmProviderStatus(server, 'anthropic')).toBe('error');
    });

    it('distinguishes configured-but-inactive from unconfigured', () => {
        const server = makeServer({ gemini: keyProvider({ apiKeySet: true }) });

        expect(getLlmProviderStatus(server, 'gemini')).toBe('inactive');
        expect(getLlmProviderStatus(server, 'anthropic')).toBe('unconfigured');
    });

    it.each([
        ['default auth', { defaultAuth: true }, true],
        ['bearer token', { bearerTokenSet: true }, true],
        ['access key pair', { accessKeyIdSet: true, secretAccessKeySet: true }, true],
        ['access key id only', { accessKeyIdSet: true }, false],
        ['nothing', {}, false],
    ])('treats bedrock with %s as configured=%s', (_label, bedrock, expected) => {
        const base = makeServer();
        const server = makeServer({ bedrock: { ...base.bedrock, ...bedrock } });

        expect(isLlmProviderConfigured(server, 'bedrock')).toBe(expected);
    });

    it('treats ollama and custom as configured once a server URL is set', () => {
        const base = makeServer();
        const server = makeServer({
            custom: { ...base.custom, apiKeySet: true },
            ollama: { ...base.ollama, serverUrl: 'http://ollama:11434' },
        });

        expect(isLlmProviderConfigured(server, 'ollama')).toBe(true);
        expect(isLlmProviderConfigured(server, 'custom')).toBe(false);
    });
});
