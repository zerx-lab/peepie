import { useSyncExternalStore } from 'react';
import { useTranslation } from 'react-i18next';

import { getCurrentLanguage, getLanguagePreference, setLanguagePreference, subscribeLanguagePreference } from '@/i18n';
import { type Language, type LanguagePreference } from '@/i18n/languages';

interface UseLanguageResult {
    /** Language currently rendered. */
    language: Language;
    /** What the user picked; `system` means `language` follows the browser. */
    preference: LanguagePreference;
    setPreference: (preference: LanguagePreference) => Promise<void>;
}

export function useLanguage(): UseLanguageResult {
    // Subscribes to `languageChanged` so the returned language stays current.
    useTranslation();

    const preference = useSyncExternalStore(subscribeLanguagePreference, getLanguagePreference);

    return {
        language: getCurrentLanguage(),
        preference,
        setPreference: setLanguagePreference,
    };
}
