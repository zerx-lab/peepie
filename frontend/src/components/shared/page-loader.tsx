import { useTranslation } from 'react-i18next';

function PageLoader() {
    const { t } = useTranslation(['ui', 'common']);

    return (
        <div className="grid h-screen w-full place-items-center">
            <p>{t('common:status.loading')}</p>
        </div>
    );
}

export default PageLoader;
