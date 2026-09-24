import { FolderPlus } from 'lucide-react';
import { useEffect, useMemo } from 'react';
import { Trans, useTranslation } from 'react-i18next';

import { Button } from '@/components/ui/button';
import { Dialog, DialogContent, DialogDescription, DialogHeader, DialogTitle } from '@/components/ui/dialog';
import { Form, FormControl, FormDescription, FormField, FormItem, FormLabel, FormMessage } from '@/components/ui/form';
import { FormSubmitButton } from '@/components/ui/form-submit-button';
import { Input } from '@/components/ui/input';
import { useAppForm } from '@/hooks/use-app-form';

import {
    createResourcesMkdirFormSchema,
    type ResourcesMkdirFormValues,
    useResourcesMkdir,
} from './use-resources-mkdir';

interface ResourcesMkdirDialogFormProps {
    defaultParentPath: string;
    onClose: () => void;
}

interface ResourcesMkdirDialogProps {
    /** Pre-filled parent path (without leading "/"). Empty string targets the root. */
    defaultParentPath?: string;
    isOpen: boolean;
    onClose: () => void;
}

const buildDefaultPath = (defaultParentPath: string): string =>
    defaultParentPath ? `${defaultParentPath.replace(/\/+$/u, '')}/new-folder` : 'new-folder';

export function ResourcesMkdirDialog({ defaultParentPath = '', isOpen, onClose }: ResourcesMkdirDialogProps) {
    const handleDialogOpenChange = (nextOpen: boolean) => {
        if (!nextOpen) {
            onClose();
        }
    };

    return (
        <Dialog
            onOpenChange={handleDialogOpenChange}
            open={isOpen}
        >
            {isOpen && (
                <ResourcesMkdirDialogForm
                    defaultParentPath={defaultParentPath}
                    onClose={onClose}
                />
            )}
        </Dialog>
    );
}

function ResourcesMkdirDialogForm({ defaultParentPath, onClose }: ResourcesMkdirDialogFormProps) {
    const { t } = useTranslation(['resources', 'common']);
    const { isCreating, mkdir } = useResourcesMkdir();
    const schema = useMemo(() => createResourcesMkdirFormSchema(t), [t]);

    const form = useAppForm<ResourcesMkdirFormValues>({
        defaultValues: { path: buildDefaultPath(defaultParentPath) },
        schema,
    });

    useEffect(() => {
        form.reset({ path: buildDefaultPath(defaultParentPath) });
    }, [defaultParentPath, form]);

    const handleSubmit = form.handleSubmit(async (values) => {
        const wasCreated = await mkdir(values);

        if (wasCreated) {
            onClose();
        }
    });

    return (
        <DialogContent>
            <DialogHeader>
                <DialogTitle className="flex items-center gap-2">
                    <FolderPlus className="size-4" />
                    {t('mkdir.title')}
                </DialogTitle>
                <DialogDescription>{t('mkdir.description')}</DialogDescription>
            </DialogHeader>

            <Form {...form}>
                <form
                    className="flex flex-col gap-4"
                    noValidate
                    onSubmit={handleSubmit}
                >
                    <FormField
                        control={form.control}
                        name="path"
                        render={({ field }) => (
                            <FormItem>
                                <FormLabel>{t('mkdir.pathLabel')}</FormLabel>
                                <FormControl>
                                    <Input
                                        {...field}
                                        autoComplete="off"
                                        autoFocus
                                        disabled={isCreating}
                                        placeholder={t('mkdir.pathPlaceholder')}
                                    />
                                </FormControl>
                                <FormDescription>
                                    <Trans
                                        components={{ code: <code /> }}
                                        i18nKey="mkdir.pathHint"
                                        t={t}
                                    />
                                </FormDescription>
                                <FormMessage />
                            </FormItem>
                        )}
                    />

                    <div className="flex justify-end gap-2">
                        <Button
                            disabled={isCreating}
                            onClick={onClose}
                            type="button"
                            variant="outline"
                        >
                            {t('common:actions.cancel')}
                        </Button>
                        <FormSubmitButton icon={<FolderPlus />}>{t('common:actions.create')}</FormSubmitButton>
                    </div>
                </form>
            </Form>
        </DialogContent>
    );
}
