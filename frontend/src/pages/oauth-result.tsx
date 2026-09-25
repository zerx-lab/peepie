import { useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import Logo from '@/components/icons/logo';
import { routes } from '@/lib/routes';

type OAuthResultStatus = 'complete' | 'directOpen' | 'inProgress' | 'parentError';

const STATUS_MESSAGE_KEYS = {
    complete: 'oauthResult.complete',
    directOpen: 'oauthResult.directOpen',
    inProgress: 'oauthResult.inProgress',
    parentError: 'oauthResult.parentError',
} as const satisfies Record<OAuthResultStatus, string>;

function OAuthResult() {
    const { t } = useTranslation('auth');
    const [resultStatus, setResultStatus] = useState<OAuthResultStatus>('inProgress');

    const successDelay = 2000;
    const errorDelay = 5000;

    useEffect(() => {
        const params = new URLSearchParams(window.location.search);
        const status = params.get('status');
        const error = params.get('error');

        let redirectTimer: null | ReturnType<typeof setTimeout> = null;
        let cleanupTimer: null | ReturnType<typeof setTimeout> = null;
        let closeTimer: null | ReturnType<typeof setTimeout> = null;

        const updateStatus = (next: OAuthResultStatus) => {
            setResultStatus(next);
        };

        const handleClose = (delay: number) => {
            closeTimer = setTimeout(() => {
                try {
                    if (window && !window.closed) {
                        window.close();
                    }
                } catch (e) {
                    console.error('Delayed window close failed:', e);
                }
            }, delay);
        };

        const handleRedirect = (url: string, delay: number) => {
            redirectTimer = setTimeout(() => {
                try {
                    window.location.href = url;
                } catch (e) {
                    console.error('Redirection failed:', e);
                }
            }, delay);

            cleanupTimer = setTimeout(() => {
                if (redirectTimer) {
                    clearTimeout(redirectTimer);
                    redirectTimer = null;
                }
            }, delay + 100);
        };

        if (window.opener) {
            try {
                window.opener.postMessage(
                    {
                        error,
                        status,
                        type: 'oauth-result',
                    },
                    window.location.origin,
                );

                updateStatus('complete');
                handleClose(successDelay);
            } catch (e) {
                console.error('Failed to send message to opener:', e);
                updateStatus('parentError');
                handleClose(errorDelay);
            }
        } else {
            updateStatus('directOpen');
            handleRedirect(routes.login(), errorDelay / 2);
            handleClose(errorDelay);
        }

        return () => {
            if (redirectTimer) {
                clearTimeout(redirectTimer);
            }

            if (cleanupTimer) {
                clearTimeout(cleanupTimer);
            }

            if (closeTimer) {
                clearTimeout(closeTimer);
            }
        };
    }, [successDelay, errorDelay]);

    return (
        <div className="flex h-screen w-full items-center justify-center bg-linear-to-r from-slate-800 to-slate-950">
            <Logo className="m-auto h-20 w-auto text-white" />
            <div className="fixed bottom-4 text-sm text-white">{t(STATUS_MESSAGE_KEYS[resultStatus])}</div>
        </div>
    );
}

export default OAuthResult;
