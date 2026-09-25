import { loader } from '@monaco-editor/react';
import * as monaco from 'monaco-editor/editor/editor.api';
import EditorWorker from 'monaco-editor/editor/editor.worker?worker';
import 'monaco-editor/basic-languages/monaco.contribution';

/*
 * Bundles Monaco (the VS Code editor) locally instead of letting
 * `@monaco-editor/react` pull it from a CDN, so previews work in air-gapped
 * deployments. Only the core editor plus the Monarch tokenizers are loaded —
 * the viewer is read-only, so the heavy TS/JSON/CSS/HTML language-service
 * workers are not needed and every worker label falls back to the base worker.
 */
self.MonacoEnvironment = {
    getWorker: () => new EditorWorker(),
};

loader.config({ monaco });

/** Resolve a Monaco language id from a file name via registered extensions/filenames. */
export const getMonacoLanguage = (fileName: string): string => {
    const lowerName = fileName.toLowerCase();

    for (const language of monaco.languages.getLanguages()) {
        if (language.filenames?.some((name) => name.toLowerCase() === lowerName)) {
            return language.id;
        }
    }

    let bestMatch: null | { id: string; length: number } = null;

    for (const language of monaco.languages.getLanguages()) {
        for (const extension of language.extensions ?? []) {
            const lowerExtension = extension.toLowerCase();

            if (lowerName.endsWith(lowerExtension) && (!bestMatch || lowerExtension.length > bestMatch.length)) {
                bestMatch = { id: language.id, length: lowerExtension.length };
            }
        }
    }

    return bestMatch?.id ?? 'plaintext';
};
