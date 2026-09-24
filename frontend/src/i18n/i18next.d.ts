import 'i18next';

import type { DEFAULT_NAMESPACE, sourceResources } from './resources';

// Type-checks every `t('ns:key')` / `useTranslation('ns')` call against the English
// source files: a typo or a key missing from `locales/en` is a compile error.
declare module 'i18next' {
    interface CustomTypeOptions {
        defaultNS: typeof DEFAULT_NAMESPACE;
        resources: typeof sourceResources;
        returnNull: false;
    }
}
