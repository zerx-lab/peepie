import { afterEach, describe, expect, it, vi } from 'vitest';

import i18n, { getLanguagePreference, setLanguagePreference, translateRuntimeKey } from '.';
import { formatLocalizedDate } from './format';
import { LANGUAGE_STORAGE_KEY } from './languages';

describe('language switching', () => {
    afterEach(async () => {
        await setLanguagePreference('system');
    });

    it('lazily loads the chosen language, persists it and updates <html lang>', async () => {
        await setLanguagePreference('zh-CN');

        expect(i18n.language).toBe('zh-CN');
        expect(i18n.t('actions.save')).toBe('保存');
        expect(document.documentElement.lang).toBe('zh-CN');
        expect(localStorage.getItem(LANGUAGE_STORAGE_KEY)).toBe('zh-CN');
        expect(getLanguagePreference()).toBe('zh-CN');
        expect(formatLocalizedDate(new Date(2024, 2, 5, 9, 7), 'dateTime')).toBe('2024年3月5日 09:07');
    });

    it('returns to the browser language when following the system', async () => {
        await setLanguagePreference('zh-CN');
        await setLanguagePreference('system');

        expect(i18n.language).toBe('en');
        expect(i18n.t('actions.save')).toBe('Save');
        expect(localStorage.getItem(LANGUAGE_STORAGE_KEY)).toBeNull();
        expect(formatLocalizedDate(new Date(2024, 2, 5, 9, 7), 'dateTime')).toBe('09:07, 5 Mar 2024');
    });

    it('settles on the last choice when switching rapidly', async () => {
        await Promise.all([setLanguagePreference('zh-CN'), setLanguagePreference('en')]);

        expect(i18n.language).toBe('en');
        expect(getLanguagePreference()).toBe('en');
    });

    it('follows a choice made in another tab', async () => {
        // Another tab wrote the preference; this tab only receives the storage event.
        localStorage.setItem(LANGUAGE_STORAGE_KEY, 'zh-CN');
        window.dispatchEvent(new StorageEvent('storage', { key: LANGUAGE_STORAGE_KEY, newValue: 'zh-CN' }));

        await vi.waitFor(() => expect(i18n.language).toBe('zh-CN'));
        expect(getLanguagePreference()).toBe('zh-CN');
    });
});

describe('translateRuntimeKey', () => {
    it('resolves dotted backend codes literally and reports unknown ones', () => {
        expect(translateRuntimeKey('apiErrors', 'Auth.InvalidCredentials', { flat: true })).toBe(
            'invalid login or password',
        );
        expect(translateRuntimeKey('apiErrors', 'Auth.DoesNotExist', { flat: true })).toBeUndefined();
    });
});
