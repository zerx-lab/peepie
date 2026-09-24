/**
 * Supported UI languages. Adding a language means: append its BCP 47 tag here, add
 * `locales/<tag>/*.json` for every namespace (the parity test enforces completeness),
 * a native name in `LANGUAGE_NATIVE_NAMES`, and a date-fns locale in `./format.ts`.
 */
export const SUPPORTED_LANGUAGES = ['en', 'zh-CN'] as const;

export type Language = (typeof SUPPORTED_LANGUAGES)[number];

/** `system` follows the browser/OS language list; any other value is an explicit user choice. */
export type LanguagePreference = 'system' | Language;

/** Source language of the UI and the fallback for keys missing in other languages. */
export const DEFAULT_LANGUAGE: Language = 'en';

export const LANGUAGE_STORAGE_KEY = 'language';

/**
 * Language names shown in the switcher. Always rendered in their own language so a user
 * can find theirs regardless of the active UI language — never translate these.
 */
export const LANGUAGE_NATIVE_NAMES: Record<Language, string> = {
    en: 'English',
    'zh-CN': '简体中文',
};

/**
 * Primary subtag → supported language used when there is no exact match. Every Chinese
 * variant (zh, zh-TW, zh-Hant, zh_HK…) maps to Simplified Chinese — the only Chinese
 * translation shipped — which beats falling back to English for those readers.
 */
const PRIMARY_SUBTAG_FALLBACKS: Record<string, Language> = {
    en: 'en',
    zh: 'zh-CN',
};

export const isLanguage = (value: unknown): value is Language =>
    typeof value === 'string' && (SUPPORTED_LANGUAGES as readonly string[]).includes(value);

export const isLanguagePreference = (value: unknown): value is LanguagePreference =>
    value === 'system' || isLanguage(value);

/**
 * Maps an arbitrary locale tag (`zh_CN.UTF-8`, `en-GB`, `ZH-hans`) onto a supported
 * language, or `undefined` when nothing matches.
 */
export const normalizeLanguage = (tag: null | string | undefined): Language | undefined => {
    const cleaned = tag?.trim().split('.')[0]?.split('@')[0]?.replaceAll('_', '-').toLowerCase();

    if (!cleaned) {
        return undefined;
    }

    const exact = SUPPORTED_LANGUAGES.find((language) => language.toLowerCase() === cleaned);

    if (exact) {
        return exact;
    }

    const primary = cleaned.split('-')[0] ?? '';

    return PRIMARY_SUBTAG_FALLBACKS[primary];
};

/** First supported language from the browser's preference list, else the default. */
export const detectSystemLanguage = (candidates: readonly string[] = getNavigatorLanguages()): Language => {
    for (const candidate of candidates) {
        const language = normalizeLanguage(candidate);

        if (language) {
            return language;
        }
    }

    return DEFAULT_LANGUAGE;
};

export const resolveLanguage = (preference: LanguagePreference): Language =>
    preference === 'system' ? detectSystemLanguage() : preference;

export const readLanguagePreference = (): LanguagePreference => {
    try {
        const stored = localStorage.getItem(LANGUAGE_STORAGE_KEY);

        return isLanguagePreference(stored) ? stored : 'system';
    } catch {
        return 'system';
    }
};

export const writeLanguagePreference = (preference: LanguagePreference): void => {
    try {
        if (preference === 'system') {
            localStorage.removeItem(LANGUAGE_STORAGE_KEY);
        } else {
            localStorage.setItem(LANGUAGE_STORAGE_KEY, preference);
        }
    } catch {
        // Storage can be unavailable (privacy mode, quota); the choice then lasts for the session.
    }
};

function getNavigatorLanguages(): readonly string[] {
    if (typeof navigator === 'undefined') {
        return [];
    }

    if (navigator.languages?.length) {
        return navigator.languages;
    }

    return navigator.language ? [navigator.language] : [];
}
