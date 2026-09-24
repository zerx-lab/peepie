import { Languages } from 'lucide-react';
import { useTranslation } from 'react-i18next';

import {
    DropdownMenuPortal,
    DropdownMenuRadioGroup,
    DropdownMenuRadioItem,
    DropdownMenuSub,
    DropdownMenuSubContent,
    DropdownMenuSubTrigger,
} from '@/components/ui/dropdown-menu';
import { Select, SelectContent, SelectItem, SelectTrigger, SelectValue } from '@/components/ui/select';
import { useLanguage } from '@/hooks/use-language';
import {
    detectSystemLanguage,
    isLanguagePreference,
    LANGUAGE_NATIVE_NAMES,
    type LanguagePreference,
    SUPPORTED_LANGUAGES,
} from '@/i18n/languages';

/** Language submenu for dropdown menus (e.g. the sidebar user menu). */
export function LanguageMenuSub() {
    const { t } = useTranslation();
    const { preference, setPreference } = useLanguage();
    const options = useLanguageOptions();

    return (
        <DropdownMenuSub>
            <DropdownMenuSubTrigger>
                <Languages />
                {t('language.label')}
            </DropdownMenuSubTrigger>
            <DropdownMenuPortal>
                <DropdownMenuSubContent>
                    <DropdownMenuRadioGroup
                        onValueChange={(value) => {
                            if (isLanguagePreference(value)) {
                                void setPreference(value);
                            }
                        }}
                        value={preference}
                    >
                        {options.map((option) => (
                            <DropdownMenuRadioItem
                                key={option.value}
                                value={option.value}
                            >
                                {option.label}
                            </DropdownMenuRadioItem>
                        ))}
                    </DropdownMenuRadioGroup>
                </DropdownMenuSubContent>
            </DropdownMenuPortal>
        </DropdownMenuSub>
    );
}

/** Compact standalone language picker for pages without the app shell (e.g. login). */
export function LanguageSelect({ className }: { className?: string }) {
    const { t } = useTranslation();
    const { preference, setPreference } = useLanguage();
    const options = useLanguageOptions();

    return (
        <Select
            onValueChange={(value) => {
                if (isLanguagePreference(value)) {
                    void setPreference(value);
                }
            }}
            value={preference}
        >
            <SelectTrigger
                aria-label={t('language.label')}
                className={className}
            >
                <Languages className="size-4" />
                <SelectValue />
            </SelectTrigger>
            <SelectContent>
                {options.map((option) => (
                    <SelectItem
                        key={option.value}
                        value={option.value}
                    >
                        {option.label}
                    </SelectItem>
                ))}
            </SelectContent>
        </Select>
    );
}

function useLanguageOptions(): { label: string; value: LanguagePreference }[] {
    const { t } = useTranslation();

    return [
        {
            label: t('language.systemWithValue', { language: LANGUAGE_NATIVE_NAMES[detectSystemLanguage()] }),
            value: 'system',
        },
        ...SUPPORTED_LANGUAGES.map((language) => ({ label: LANGUAGE_NATIVE_NAMES[language], value: language })),
    ];
}
