import { act, render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { createMemoryRouter, RouterProvider } from 'react-router-dom';
import { afterEach, expect, it, vi } from 'vitest';

import { SidebarProvider } from '@/components/ui/sidebar';
import i18n from '@/i18n';

vi.mock('@apollo/client/react', async (importOriginal) => ({
    ...(await importOriginal<Record<string, unknown>>()),
    useQuery: () => ({ data: undefined, error: undefined, loading: false, refetch: vi.fn() }),
}));

vi.mock('@/hooks/use-breakpoint', () => ({
    useBreakpoint: () => ({ isDesktop: true, isMobile: false }),
}));

vi.mock('@/providers/templates-provider', () => ({
    useTemplates: () => ({
        createTemplate: vi.fn(),
        deleteTemplate: vi.fn(),
        templates: [],
        updateTemplate: vi.fn(),
    }),
}));

vi.mock('@/components/shared/markdown-editor', () => ({
    EditorViewModeToggle: () => null,
    MarkdownEditorField: ({
        'aria-label': label,
        onChange,
        value,
    }: {
        'aria-label': string;
        onChange: (value: string) => void;
        value: string;
    }) => (
        <textarea
            aria-label={label}
            onChange={(event) => onChange(event.target.value)}
            value={value}
        />
    ),
}));

import TemplatePage from './template';

afterEach(async () => {
    await act(async () => {
        await i18n.changeLanguage('en');
    });
});

it('localizes the preset chooser while preserving the preset payload sent to the form', async () => {
    const user = userEvent.setup();
    const router = createMemoryRouter(
        [
            {
                element: (
                    <SidebarProvider>
                        <TemplatePage />
                    </SidebarProvider>
                ),
                path: '/templates/:templateId',
            },
        ],
        { initialEntries: ['/templates/new'] },
    );

    render(<RouterProvider router={router} />);
    expect(screen.getByRole('heading', { name: 'Create a new template' })).toBeInTheDocument();

    await act(async () => {
        await i18n.changeLanguage('zh-CN');
    });

    expect(screen.getByRole('heading', { name: '创建新模板' })).toBeInTheDocument();
    await user.click(screen.getByRole('button', { name: 'Web 应用安全评估' }));
    expect(screen.getByRole('textbox', { name: '标题' })).toHaveValue('Web Application Security Assessment');
});
