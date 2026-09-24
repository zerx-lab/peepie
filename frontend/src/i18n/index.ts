import type { BackendModule } from 'i18next';

import i18n from 'i18next';
import { initReactI18next } from 'react-i18next';

import {
    DEFAULT_LANGUAGE,
    isLanguage,
    type Language,
    LANGUAGE_STORAGE_KEY,
    type LanguagePreference,
    readLanguagePreference,
    resolveLanguage,
    SUPPORTED_LANGUAGES,
    writeLanguagePreference,
} from './languages';
import { DEFAULT_NAMESPACE, loadNamespace, type Namespace, NAMESPACES, sourceResources } from './resources';

const lazyBackend: BackendModule = {
    init: () => undefined,
    read: (language, namespace, callback) => {
        if (!isLanguage(language)) {
            callback(new Error(`Unsupported language: ${language}`), false);

            return;
        }

        loadNamespace(language, namespace).then(
            (data) => callback(null, data),
            (error: Error) => callback(error, false),
        );
    },
    type: 'backend',
};

let preference: LanguagePreference = readLanguagePreference();
const preferenceListeners = new Set<() => void>();

/**
 * Resolves once the initial language's namespaces are loaded (or failed to load, in
 * which case keys fall back to English). Render after it to avoid an English flash.
 */
export const i18nReady: Promise<void> = i18n
    .use(lazyBackend)
    .use(initReactI18next)
    .init({
        defaultNS: DEFAULT_NAMESPACE,
        fallbackLng: DEFAULT_LANGUAGE,
        // English is bundled, so init completes synchronously for it (tests, first paint).
        initAsync: false,
        interpolation: {
            // React already escapes rendered strings.
            escapeValue: false,
        },
        lng: resolveLanguage(preference),
        // Only the exact tag (`zh-CN`), never the bare `zh` which has no files.
        load: 'currentOnly',
        ns: NAMESPACES,
        partialBundledLanguages: true,
        react: {
            // Every namespace of a language is loaded before it becomes active, so
            // components never need to suspend.
            useSuspense: false,
        },
        resources: { [DEFAULT_LANGUAGE]: sourceResources },
        returnNull: false,
        supportedLngs: SUPPORTED_LANGUAGES,
    })
    .then(
        () => undefined,
        () => undefined,
    );

const syncDocumentLanguage = (language: string) => {
    if (typeof document !== 'undefined') {
        document.documentElement.lang = language;
    }
};

syncDocumentLanguage(i18n.language);
i18n.on('languageChanged', syncDocumentLanguage);

const applyLanguagePreference = async (next: LanguagePreference): Promise<void> => {
    preference = next;
    preferenceListeners.forEach((listener) => listener());
    // i18next ignores a load that finishes after a newer changeLanguage call, so rapid
    // switching settles on the last choice.
    await i18n.changeLanguage(resolveLanguage(next));
};

if (typeof window !== 'undefined') {
    // Follow OS/browser language changes live while the user has not picked one explicitly.
    window.addEventListener('languagechange', () => {
        if (preference === 'system') {
            void i18n.changeLanguage(resolveLanguage('system'));
        }
    });

    // Keep every open tab on the choice made in any of them (`storage` fires only in other tabs).
    window.addEventListener('storage', (event) => {
        if (event.key === LANGUAGE_STORAGE_KEY || event.key === null) {
            void applyLanguagePreference(readLanguagePreference());
        }
    });
}

export const getLanguagePreference = (): LanguagePreference => preference;

export const subscribeLanguagePreference = (listener: () => void): (() => void) => {
    preferenceListeners.add(listener);

    return () => preferenceListeners.delete(listener);
};

/** Persists the user's choice and switches the UI once the target language is loaded. */
export async function setLanguagePreference(next: LanguagePreference): Promise<void> {
    writeLanguagePreference(next);
    await applyLanguagePreference(next);
}

/** Active UI language, narrowed to a supported one (i18next may report `cimode` in tests). */
export const getCurrentLanguage = (): Language => (isLanguage(i18n.language) ? i18n.language : DEFAULT_LANGUAGE);

/**
 * Translates a key only known at runtime (e.g. a backend error code), returning
 * `undefined` when no translation exists so callers can fall back to server text.
 * Static UI text must use typed `t('ns:key')` instead — this bypasses key checking.
 * `flat` treats dots as part of the key (namespaces keyed by dotted codes, e.g. `apiErrors`).
 */
export const translateRuntimeKey = (
    namespace: Namespace,
    key: string,
    { flat = false }: { flat?: boolean } = {},
): string | undefined => {
    const options = { ns: namespace, ...(flat && { keySeparator: false as const }) };

    return i18n.exists(key, options) ? String(i18n.t(key as never, options)) : undefined;
};

export default i18n;
