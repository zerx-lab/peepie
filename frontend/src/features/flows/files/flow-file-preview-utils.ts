export type FlowFilePreviewKind = 'markdown' | 'pdf' | 'text';

/** Text previews above this size are refused: Monaco stays responsive, but fetching and decoding stops being cheap. */
export const MAX_TEXT_PREVIEW_BYTES = 5 * 1024 * 1024;
export const MAX_PDF_PREVIEW_BYTES = 50 * 1024 * 1024;

const MARKDOWN_EXTENSIONS = ['.md', '.markdown', '.mdown', '.mkd'];

/** Everything that is neither PDF nor Markdown is fetched as text and rejected later if it turns out binary. */
export const getFlowFilePreviewKind = (fileName: string): FlowFilePreviewKind => {
    const lowerName = fileName.toLowerCase();

    if (lowerName.endsWith('.pdf')) {
        return 'pdf';
    }

    return MARKDOWN_EXTENSIONS.some((extension) => lowerName.endsWith(extension)) ? 'markdown' : 'text';
};

const BINARY_SNIFF_BYTES = 8000;

/**
 * Decode bytes as UTF-8 text, or return `null` for binary content. Uses the same
 * heuristic as git/VS Code: a NUL byte in the leading chunk means binary; invalid
 * UTF-8 is also treated as binary so the preview never shows mojibake.
 */
export const decodeTextPreview = (buffer: ArrayBuffer): null | string => {
    const bytes = new Uint8Array(buffer);

    if (bytes.subarray(0, BINARY_SNIFF_BYTES).includes(0)) {
        return null;
    }

    try {
        return new TextDecoder('utf-8', { fatal: true }).decode(bytes);
    } catch {
        return null;
    }
};
