import { afterEach, describe, expect, it } from 'vitest';

import {
    detectSystemLanguage,
    LANGUAGE_STORAGE_KEY,
    normalizeLanguage,
    readLanguagePreference,
    writeLanguagePreference,
} from './languages';

describe('normalizeLanguage', () => {
    it.each([
        ['en', 'en'],
        ['en-GB', 'en'],
        ['zh-CN', 'zh-CN'],
        ['zh_CN.UTF-8', 'zh-CN'],
        ['ZH-cn', 'zh-CN'],
        ['zh', 'zh-CN'],
        ['zh-Hant-TW', 'zh-CN'],
        ['fr-FR', undefined],
        ['', undefined],
        [undefined, undefined],
    ])('%s → %s', (tag, expected) => {
        expect(normalizeLanguage(tag)).toBe(expected);
    });
});

describe('detectSystemLanguage', () => {
    it('picks the first supported language in preference order', () => {
        expect(detectSystemLanguage(['fr-FR', 'zh-TW', 'en-US'])).toBe('zh-CN');
    });

    it('falls back to English when nothing is supported', () => {
        expect(detectSystemLanguage(['fr-FR', 'de'])).toBe('en');
        expect(detectSystemLanguage([])).toBe('en');
    });
});

describe('language preference storage', () => {
    afterEach(() => localStorage.clear());

    it('defaults to following the system', () => {
        expect(readLanguagePreference()).toBe('system');
    });

    it('persists an explicit choice and clears it when returning to system', () => {
        writeLanguagePreference('zh-CN');
        expect(readLanguagePreference()).toBe('zh-CN');

        writeLanguagePreference('system');
        expect(localStorage.getItem(LANGUAGE_STORAGE_KEY)).toBeNull();
        expect(readLanguagePreference()).toBe('system');
    });

    it('ignores unsupported stored values', () => {
        localStorage.setItem(LANGUAGE_STORAGE_KEY, 'xx');
        expect(readLanguagePreference()).toBe('system');
    });
});
