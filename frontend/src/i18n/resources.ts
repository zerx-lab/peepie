import type { Language } from './languages';

import apiErrors from './locales/en/apiErrors.json';
import auth from './locales/en/auth.json';
import common from './locales/en/common.json';
import dashboard from './locales/en/dashboard.json';
import editor from './locales/en/editor.json';
import errors from './locales/en/errors.json';
import fileManager from './locales/en/fileManager.json';
import flowDetails from './locales/en/flowDetails.json';
import flows from './locales/en/flows.json';
import knowledges from './locales/en/knowledges.json';
import layout from './locales/en/layout.json';
import providers from './locales/en/providers.json';
import resources from './locales/en/resources.json';
import settings from './locales/en/settings.json';
import templates from './locales/en/templates.json';
import ui from './locales/en/ui.json';

/**
 * English is the source language: bundled eagerly (synchronous first render, fallback
 * for any missing key) and the type source for `t()` key checking (`./i18next.d.ts`).
 * Registering a namespace = add `locales/<lng>/<ns>.json` for every language and import
 * the English file here; `locales.test.ts` fails when a file on disk is not registered.
 * `apiErrors` is keyed by backend error code (dots are literal) and read through
 * `translateRuntimeKey(…, { flat: true })`.
 */
export const sourceResources = {
    apiErrors,
    auth,
    common,
    dashboard,
    editor,
    errors,
    fileManager,
    flowDetails,
    flows,
    knowledges,
    layout,
    providers,
    resources,
    settings,
    templates,
    ui,
} as const;

export type Namespace = keyof typeof sourceResources;

export const NAMESPACES = Object.keys(sourceResources) as Namespace[];

export const DEFAULT_NAMESPACE = 'common' satisfies Namespace;

type NamespaceModule = { default: Record<string, unknown> };

// Non-source languages are code-split per namespace and fetched only when selected.
const lazyLoaders = import.meta.glob<NamespaceModule>(['./locales/*/*.json', '!./locales/en/*.json']);

export async function loadNamespace(language: Language, namespace: string): Promise<Record<string, unknown>> {
    const loader = lazyLoaders[`./locales/${language}/${namespace}.json`];

    if (!loader) {
        throw new Error(`Missing translations for ${language}/${namespace}`);
    }

    return (await loader()).default;
}
