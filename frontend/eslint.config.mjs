// @ts-check
import { FlatCompat } from '@eslint/eslintrc';
import js from '@eslint/js';
import i18next from 'eslint-plugin-i18next';
import perfectionist from 'eslint-plugin-perfectionist';
import playwright from 'eslint-plugin-playwright';

const compat = new FlatCompat({
    baseDirectory: import.meta.dirname,
    recommendedConfig: js.configs.recommended,
});

const eslintConfig = [
    ...compat.config({
        extends: [
            'eslint:recommended',
            'plugin:@typescript-eslint/recommended',
            'plugin:react/recommended',
            'plugin:react/jsx-runtime',
            'plugin:react-hooks/recommended',
            'prettier',
        ],
        settings: {
            react: {
                version: 'detect',
            },
        },
    }),
    {
        rules: {
            '@typescript-eslint/no-explicit-any': 'warn',
            '@typescript-eslint/no-unused-vars': [
                'error',
                {
                    argsIgnorePattern: '^_',
                    varsIgnorePattern: '^_',
                },
            ],
            curly: ['error', 'all'],
            'no-fallthrough': 'off',
            'no-restricted-syntax': [
                'error',
                {
                    message:
                        'Do not set mode/reValidateMode on useForm — use useAppForm, which owns the silent-until-submit form timing.',
                    selector:
                        'CallExpression[callee.name="useForm"] > ObjectExpression > Property[key.name="mode"], CallExpression[callee.name="useForm"] > ObjectExpression > Property[key.name="reValidateMode"]',
                },
            ],
            'padding-line-between-statements': [
                'error',
                {
                    blankLine: 'always',
                    next: 'return',
                    prev: '*',
                },
                {
                    blankLine: 'always',
                    next: 'block-like',
                    prev: '*',
                },
                {
                    blankLine: 'any',
                    next: 'block-like',
                    prev: 'case',
                },
                {
                    blankLine: 'always',
                    next: '*',
                    prev: 'block-like',
                },
                {
                    blankLine: 'always',
                    next: 'block-like',
                    prev: 'block-like',
                },
                {
                    blankLine: 'any',
                    next: 'while',
                    prev: 'do',
                },
            ],
            // Off until React Compiler lands: its compatibility lint flags every RHF
            // watch() / useReactTable as unactionable noise (see OPEN-ISSUES §3).
            'react-hooks/incompatible-library': 'off',
            'react/no-unescaped-entities': 'off', // Allow quotes in JSX
            'react/prop-types': 'off', // TypeScript provides type checking
        },
    },
    {
        // useAppForm is the single owner of the form timing — the one place
        // allowed to set mode/reValidateMode on useForm.
        files: ['src/hooks/use-app-form.ts'],
        rules: { 'no-restricted-syntax': 'off' },
    },
    perfectionist.configs['recommended-natural'],
    {
        // Every user-visible string goes through i18next (see docs/i18n.md): JSX text and
        // human-facing props must come from `t()`. Technical literals (variants, routes,
        // ids) are not in the attribute list, so they stay allowed.
        files: ['src/**/*.tsx'],
        ignores: ['src/**/*.test.tsx', 'src/components/icons/**', 'src/lib/report/**'],
        plugins: { i18next },
        rules: {
            'i18next/no-literal-string': [
                'error',
                {
                    'jsx-attributes': {
                        include: [
                            'alt',
                            'aria-description',
                            'aria-label',
                            'aria-placeholder',
                            'aria-roledescription',
                            'aria-valuetext',
                            'cancelText',
                            'confirmText',
                            'description',
                            'emptyMessage',
                            'emptyText',
                            'heading',
                            'label',
                            'message',
                            'placeholder',
                            'subtitle',
                            'title',
                            'tooltip',
                        ],
                    },
                    mode: 'jsx-only',
                    words: {
                        // Symbols, numbers, all-caps tokens, and brand/product names are not translatable.
                        exclude: [
                            '[0-9!-/:-@[-`{-~\\s·•—–…→←↑↓×✓✗]+',
                            '[A-Z_-]+',
                            'Peepie',
                        ],
                    },
                },
            ],
        },
    },
    {
        ...playwright.configs['flat/recommended'],
        // *.unit.test.ts are vitest, not Playwright — the plugin's rules
        // (no-standalone-expect) misfire on vitest's `it`.
        files: ['e2e/**/*.ts'],
        ignores: ['e2e/**/*.unit.test.ts'],
    },
    {
        // Playwright fixtures take a `use` callback that the React hooks rule
        // mistakes for a hook call; there is no React under e2e/.
        files: ['e2e/**/*.ts'],
        rules: { 'react-hooks/rules-of-hooks': 'off' },
    },
    {
        // The dependency-free .mjs tools run under plain Node — declare its
        // globals so no-undef doesn't misfire.
        files: ['e2e/**/*.mjs'],
        languageOptions: {
            globals: {
                console: 'readonly',
                fetch: 'readonly',
                process: 'readonly',
                setTimeout: 'readonly',
                URL: 'readonly',
            },
        },
    },
    {
        ignores: [
            'node_modules/**',
            'dist/**',
            'build/**',
            'public/mockServiceWorker.js',
            'src/graphql/types.ts',
            'e2e/test-results/**',
            'e2e/playwright-report/**',
        ],
    },
];

export default eslintConfig;
