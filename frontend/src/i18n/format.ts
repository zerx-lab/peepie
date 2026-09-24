import type { Locale } from 'date-fns';

import { format } from 'date-fns';
import { enUS, zhCN } from 'date-fns/locale';

import type { sourceResources } from './resources';

import i18n, { getCurrentLanguage } from '.';
import { type Language } from './languages';

const DATE_FNS_LOCALES: Record<Language, Locale> = {
    en: enUS,
    'zh-CN': zhCN,
};

/** BCP 47 tags passed to `Intl.*`; kept explicit so `en` formats as en-US like before. */
const INTL_LOCALES: Record<Language, string> = {
    en: 'en-US',
    'zh-CN': 'zh-CN',
};

export const getDateFnsLocale = (language: Language = getCurrentLanguage()): Locale => DATE_FNS_LOCALES[language];

export const getIntlLocale = (language: Language = getCurrentLanguage()): string => INTL_LOCALES[language];

export type DateFormatKey = keyof (typeof sourceResources)['common']['formats'];

/**
 * Formats a date with a per-language pattern from `common:formats.*` (date-fns tokens)
 * and the matching date-fns locale, so both word order and month names follow the UI
 * language. Use this instead of hard-coding patterns in components.
 */
export const formatLocalizedDate = (date: Date | number, patternKey: DateFormatKey): string =>
    format(date, i18n.t(`formats.${patternKey}`, { ns: 'common' }), { locale: getDateFnsLocale() });
