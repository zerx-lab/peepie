import { describe, expect, it } from 'vitest';

import { decodeTextPreview, getFlowFilePreviewKind } from './flow-file-preview-utils';

const toBuffer = (bytes: number[]) => new Uint8Array(bytes).buffer;

describe('getFlowFilePreviewKind', () => {
    it('detects PDF and Markdown case-insensitively, everything else is text', () => {
        expect(getFlowFilePreviewKind('Report.PDF')).toBe('pdf');
        expect(getFlowFilePreviewKind('README.md')).toBe('markdown');
        expect(getFlowFilePreviewKind('notes.markdown')).toBe('markdown');
        expect(getFlowFilePreviewKind('main.go')).toBe('text');
        expect(getFlowFilePreviewKind('md')).toBe('text');
    });
});

describe('decodeTextPreview', () => {
    it('decodes UTF-8 text including multibyte characters', () => {
        expect(decodeTextPreview(new TextEncoder().encode('héllo 你好').buffer)).toBe('héllo 你好');
    });

    it('rejects content with NUL bytes as binary', () => {
        expect(decodeTextPreview(toBuffer([0x7f, 0x45, 0x4c, 0x46, 0x00, 0x01]))).toBeNull();
    });

    it('rejects invalid UTF-8 as binary', () => {
        expect(decodeTextPreview(toBuffer([0xff, 0xfe, 0x41]))).toBeNull();
    });
});
