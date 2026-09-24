import { describe, expect, it } from 'vitest';

import { DEFAULT_LANGUAGE, SUPPORTED_LANGUAGES } from './languages';
import { NAMESPACES } from './resources';

const files = import.meta.glob<Record<string, unknown>>('./locales/*/*.json', { eager: true, import: 'default' });

const PLURAL_SUFFIX = /_(zero|one|two|few|many|other)$/;
const INTERPOLATION = /\{\{\s*([\w.]+)[^}]*\}\}/g;

const catalog = (language: string, namespace: string) => files[`./locales/${language}/${namespace}.json`];

const namespacesOf = (language: string) =>
    Object.keys(files)
        .filter((path) => path.startsWith(`./locales/${language}/`))
        .map((path) => path.slice(`./locales/${language}/`.length, -'.json'.length))
        .sort();

/** Flattens nested messages to `a.b.c` → text, merging plural forms (`x_one`/`x_other`) into `x`. */
const flatten = (value: unknown, prefix = '', out = new Map<string, string[]>()) => {
    if (typeof value === 'string') {
        const key = prefix.replace(PLURAL_SUFFIX, '');

        out.set(key, [...(out.get(key) ?? []), value]);

        return out;
    }

    if (value && typeof value === 'object' && !Array.isArray(value)) {
        for (const [key, child] of Object.entries(value)) {
            flatten(child, prefix ? `${prefix}.${key}` : key, out);
        }

        return out;
    }

    throw new Error(`Translation "${prefix}" must be a string or an object, got ${JSON.stringify(value)}`);
};

const variablesOf = (texts: string[]) =>
    [...new Set(texts.flatMap((text) => [...text.matchAll(INTERPOLATION)].map((match) => match[1])))].sort();

describe('translation catalogs', () => {
    it('registers every source namespace file', () => {
        expect(namespacesOf(DEFAULT_LANGUAGE)).toEqual([...NAMESPACES].sort());
    });

    describe.each(SUPPORTED_LANGUAGES.filter((language) => language !== DEFAULT_LANGUAGE))('%s', (language) => {
        it('has a file for every namespace', () => {
            expect(namespacesOf(language)).toEqual(namespacesOf(DEFAULT_LANGUAGE));
        });

        it.each(NAMESPACES)('%s matches the source keys and interpolation variables', (namespace) => {
            const source = flatten(catalog(DEFAULT_LANGUAGE, namespace));
            const target = flatten(catalog(language, namespace));

            expect([...target.keys()].sort()).toEqual([...source.keys()].sort());

            for (const [key, texts] of source) {
                const translated = target.get(key) ?? [];

                expect(
                    translated.every((text) => text.trim() !== ''),
                    `${namespace}:${key} is empty`,
                ).toBe(true);
                expect(variablesOf(translated), `${namespace}:${key} interpolation`).toEqual(variablesOf(texts));
            }
        });
    });
});
