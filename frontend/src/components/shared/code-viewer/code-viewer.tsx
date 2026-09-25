import Editor from '@monaco-editor/react';
import { useSyncExternalStore } from 'react';

import { Spinner } from '@/components/ui/spinner';

import { getMonacoLanguage } from './monaco-setup';

interface CodeViewerProps {
    className?: string;
    content: string;
    fileName: string;
}

const subscribeToThemeClass = (onChange: () => void) => {
    const observer = new MutationObserver(onChange);
    observer.observe(document.documentElement, { attributeFilter: ['class'], attributes: true });

    return () => observer.disconnect();
};

/** Read-only Monaco (VS Code editor) view; language is inferred from `fileName`. */
function CodeViewer({ className, content, fileName }: CodeViewerProps) {
    // ThemeProvider resolves "system" onto the <html> class, so follow the class, not the stored preference.
    const isDark = useSyncExternalStore(subscribeToThemeClass, () =>
        document.documentElement.classList.contains('dark'),
    );

    return (
        <Editor
            className={className}
            language={getMonacoLanguage(fileName)}
            loading={<Spinner variant="circle" />}
            options={{
                automaticLayout: true,
                domReadOnly: true,
                minimap: { enabled: true },
                readOnly: true,
                scrollBeyondLastLine: false,
                wordWrap: 'off',
            }}
            theme={isDark ? 'vs-dark' : 'vs'}
            value={content}
        />
    );
}

export default CodeViewer;
