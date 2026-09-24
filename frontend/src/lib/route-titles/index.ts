import type { ComponentType } from 'react';

import {
    FlowDocument,
    FlowTemplateDocument,
    KnowledgeDocumentDocument,
    SettingsProvidersDocument,
} from '@/graphql/types';
import i18n from '@/i18n';

import { apolloTitle } from './apollo-title';
import { formatPromptId } from './format-prompt-id';
import { type RouteParams } from './render-title';

export interface RouteTitleHandle {
    title: TitleResolver;
}

/**
 * A `handle.title` value can be one of three forms:
 *   - `string` — fully static, known at build time (not translated; prefer a resolver).
 *   - `(params) => string` — derived synchronously from URL params. Resolved on every
 *     render of `DocumentTitle`, so resolvers calling `i18n.t` follow the UI language.
 *   - `ComponentType<{ params }>` — reactive (e.g. subscribes to Apollo
 *     cache for resource-driven titles). Must be produced by `apolloTitle()`
 *     so the marker it attaches lets `DocumentTitle` distinguish a component
 *     from a `(params) => string` resolver at runtime. A hand-rolled component
 *     function will be misdetected as a resolver and called with raw params —
 *     always route reactive titles through `apolloTitle()`.
 */
export type TitleResolver = ((params: RouteParams) => string) | ComponentType<{ params: RouteParams }> | string;

/**
 * Single source of truth for every route's document `<title>`. `app.tsx`
 * imports nothing from Apollo for title purposes — it only spreads handles
 * from this registry onto the matching <Route>. Titles are translated at
 * render time (never at module scope) so they follow language switches.
 */
export const routeTitles = {
    account: { title: () => i18n.t('layout:titles.account') },
    apiTokens: { title: () => i18n.t('layout:titles.apiTokens') },
    dashboard: { title: () => i18n.t('layout:titles.dashboard') },
    flow: {
        title: apolloTitle({
            document: FlowDocument,
            select: (data, { flowId }) =>
                data?.flow?.title && flowId
                    ? i18n.t('layout:titles.flowWithTitle', { id: flowId, title: data.flow.title })
                    : i18n.t('layout:titles.flow'),
            variables: ({ flowId }) => (flowId ? { id: flowId } : null),
        }),
    },
    flowReport: { title: () => i18n.t('layout:titles.flowReport') },
    flows: { title: () => i18n.t('layout:titles.flows') },
    knowledge: {
        title: apolloTitle({
            document: KnowledgeDocumentDocument,
            select: (data, { knowledgeId }) =>
                knowledgeId === 'new'
                    ? i18n.t('layout:titles.newKnowledge')
                    : data?.knowledgeDocument?.question || i18n.t('layout:titles.knowledge'),
            variables: ({ knowledgeId }) => (!knowledgeId || knowledgeId === 'new' ? null : { id: knowledgeId }),
        }),
    },
    knowledges: { title: () => i18n.t('layout:titles.knowledges') },
    login: { title: () => i18n.t('layout:titles.login') },
    newFlow: { title: () => i18n.t('layout:titles.newFlow') },
    oauth: { title: () => i18n.t('layout:titles.oauth') },
    prompt: {
        title: (params: RouteParams) =>
            params.promptId ? formatPromptId(params.promptId) : i18n.t('layout:titles.prompt'),
    },
    prompts: { title: () => i18n.t('layout:titles.prompts') },

    provider: {
        title: apolloTitle({
            document: SettingsProvidersDocument,
            select: (data, { providerId }) => {
                if (providerId === 'new') {
                    return i18n.t('layout:titles.newProvider');
                }

                const provider = data?.settingsProviders.userDefined?.find(
                    (candidate) => String(candidate.id) === providerId,
                );

                return provider?.name || i18n.t('layout:titles.provider');
            },
            variables: ({ providerId }) => (providerId === 'new' ? null : {}),
        }),
    },

    providers: { title: () => i18n.t('layout:titles.providers') },

    resources: { title: () => i18n.t('layout:titles.resources') },

    system: { title: () => i18n.t('layout:titles.system') },

    template: {
        title: apolloTitle({
            document: FlowTemplateDocument,
            select: (data, { templateId }) =>
                templateId === 'new'
                    ? i18n.t('layout:titles.newTemplate')
                    : data?.flowTemplate?.title || i18n.t('layout:titles.template'),
            variables: ({ templateId }) => (!templateId || templateId === 'new' ? null : { templateId }),
        }),
    },
    templates: { title: () => i18n.t('layout:titles.templates') },
} as const satisfies Record<string, RouteTitleHandle>;
