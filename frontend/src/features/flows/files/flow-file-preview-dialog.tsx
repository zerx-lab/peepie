import { FileWarning } from 'lucide-react';
import { lazy, Suspense, useEffect, useState } from 'react';
import { useTranslation } from 'react-i18next';

import type { FileNode } from '@/components/shared/file-manager';

import Markdown from '@/components/shared/markdown';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Empty, EmptyDescription, EmptyHeader, EmptyMedia, EmptyTitle } from '@/components/ui/empty';
import { Spinner } from '@/components/ui/spinner';
import { ToggleGroup, ToggleGroupItem } from '@/components/ui/toggle-group';
import { buildPathsQuery } from '@/features/resources/resources-utils';
import { axios, getApiErrorMessage } from '@/lib/axios';

import {
    decodeTextPreview,
    getFlowFilePreviewKind,
    MAX_PDF_PREVIEW_BYTES,
    MAX_TEXT_PREVIEW_BYTES,
} from './flow-file-preview-utils';

// Monaco is large; only download it when a preview is actually opened.
const CodeViewer = lazy(() => import('@/components/shared/code-viewer/code-viewer'));

interface FlowFilePreviewDialogProps {
    file: FileNode;
    flowId: string;
    onClose: () => void;
}

type MarkdownMode = 'rendered' | 'source';

type PreviewState =
    | { error: unknown; status: 'error' }
    | { kind: 'markdown' | 'text'; status: 'ready'; text: string }
    | { kind: 'pdf'; status: 'ready'; url: string }
    | { status: 'binary' }
    | { status: 'loading' }
    | { status: 'tooLarge' };

/** Mount with `key={file.path}` so every opened file starts from a fresh state. */
export function FlowFilePreviewDialog({ file, flowId, onClose }: FlowFilePreviewDialogProps) {
    const { t } = useTranslation('fileManager');
    const kind = getFlowFilePreviewKind(file.name);
    const isTooLarge =
        file.size != null && file.size > (kind === 'pdf' ? MAX_PDF_PREVIEW_BYTES : MAX_TEXT_PREVIEW_BYTES);
    const [state, setState] = useState<PreviewState>({ status: isTooLarge ? 'tooLarge' : 'loading' });
    const [markdownMode, setMarkdownMode] = useState<MarkdownMode>('rendered');

    useEffect(() => {
        if (isTooLarge) {
            return;
        }

        const controller = new AbortController();
        let objectUrl: null | string = null;

        axios
            .get<ArrayBuffer, ArrayBuffer>(`/flows/${flowId}/files/download?${buildPathsQuery([file.path])}`, {
                responseType: 'arraybuffer',
                signal: controller.signal,
                timeout: 0,
            })
            .then((buffer) => {
                if (kind === 'pdf') {
                    objectUrl = URL.createObjectURL(new Blob([buffer], { type: 'application/pdf' }));
                    setState({ kind, status: 'ready', url: objectUrl });

                    return;
                }

                const text = decodeTextPreview(buffer);
                setState(text === null ? { status: 'binary' } : { kind, status: 'ready', text });
            })
            .catch((error: unknown) => {
                if (controller.signal.aborted) {
                    return;
                }

                setState({ error, status: 'error' });
            });

        return () => {
            controller.abort();

            if (objectUrl) {
                URL.revokeObjectURL(objectUrl);
            }
        };
    }, [file.path, flowId, isTooLarge, kind]);

    const renderNotice = (title: string, description?: string) => (
        <Empty>
            <EmptyHeader>
                <EmptyMedia variant="icon">
                    <FileWarning />
                </EmptyMedia>
                <EmptyTitle>{title}</EmptyTitle>
                {description && <EmptyDescription>{description}</EmptyDescription>}
            </EmptyHeader>
        </Empty>
    );

    const renderBody = () => {
        switch (state.status) {
            case 'binary': {
                return renderNotice(t('flowFiles.preview.binaryTitle'), t('flowFiles.preview.binaryDescription'));
            }

            case 'error': {
                return renderNotice(
                    t('flowFiles.preview.loadFailed'),
                    getApiErrorMessage(state.error, t('flowFiles.preview.loadFailed')),
                );
            }

            case 'loading': {
                return (
                    <div className="flex h-full items-center justify-center">
                        <Spinner variant="circle" />
                    </div>
                );
            }

            case 'tooLarge': {
                return renderNotice(t('flowFiles.preview.tooLargeTitle'), t('flowFiles.preview.tooLargeDescription'));
            }
        }

        if (state.kind === 'pdf') {
            return (
                <iframe
                    className="size-full rounded-md border"
                    src={state.url}
                    title={file.name}
                />
            );
        }

        if (state.kind === 'markdown' && markdownMode === 'rendered') {
            return (
                <div className="size-full overflow-auto rounded-md border p-4">
                    <Markdown>{state.text}</Markdown>
                </div>
            );
        }

        return (
            <Suspense
                fallback={
                    <div className="flex h-full items-center justify-center">
                        <Spinner variant="circle" />
                    </div>
                }
            >
                <CodeViewer
                    className="overflow-hidden rounded-md border"
                    content={state.text}
                    fileName={file.name}
                />
            </Suspense>
        );
    };

    const isMarkdownReady = state.status === 'ready' && state.kind === 'markdown';

    return (
        <Dialog
            onOpenChange={(isOpen) => !isOpen && onClose()}
            open
        >
            <DialogContent className="flex h-[85vh] w-[90vw] flex-col sm:max-w-[90vw]">
                <DialogHeader className="flex-row items-center gap-4 pr-8">
                    <div className="min-w-0 flex-1">
                        <DialogTitle className="truncate">{file.name}</DialogTitle>
                        <DialogDescription className="truncate font-mono text-xs">{file.path}</DialogDescription>
                    </div>
                    {isMarkdownReady && (
                        <ToggleGroup
                            aria-label={t('flowFiles.preview.markdownMode')}
                            onValueChange={(value) => value && setMarkdownMode(value as MarkdownMode)}
                            size="sm"
                            type="single"
                            value={markdownMode}
                            variant="outline"
                        >
                            <ToggleGroupItem value="rendered">{t('flowFiles.preview.rendered')}</ToggleGroupItem>
                            <ToggleGroupItem value="source">{t('flowFiles.preview.source')}</ToggleGroupItem>
                        </ToggleGroup>
                    )}
                </DialogHeader>
                <div className="min-h-0 flex-1">{renderBody()}</div>
            </DialogContent>
        </Dialog>
    );
}
