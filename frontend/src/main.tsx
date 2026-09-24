import './styles/index.css';

import * as React from 'react';
import ReactDOM from 'react-dom/client';

import App from '@/app';
import { i18nReady } from '@/i18n';
import { reloadOnce } from '@/lib/chunk-reload';

// Vite fires this when a code-split chunk fails to load — almost always a redeploy
// that rotated the hashed filenames this document still references. Reload once to
// pull the current build; RouteErrorBoundary is the fallback for failures that
// bypass this event, and both share reloadOnce's debounce so at most one fires.
window.addEventListener('vite:preloadError', () => {
    reloadOnce();
});

// Wait for the active language's translations so a non-English UI never flashes English.
void i18nReady.then(() => {
    ReactDOM.createRoot(document.querySelector('#root')!).render(
        <React.StrictMode>
            <App />
        </React.StrictMode>,
    );
});
