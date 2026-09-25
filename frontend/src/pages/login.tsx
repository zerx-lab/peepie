import { useLocation, useSearchParams } from 'react-router-dom';

import Logo from '@/components/icons/logo';
import { LanguageSelect } from '@/components/shared/language-switcher';
import { Spinner } from '@/components/ui/spinner';
import LoginForm from '@/features/authentication/login-form';
import { routes } from '@/lib/routes';
import { getSafeReturnUrl } from '@/lib/utils/auth';
import { useUser } from '@/providers/user-provider';

function Login() {
    const [searchParams] = useSearchParams();
    const location = useLocation();
    const { authInfo, isLoading } = useUser();
    const authProviders = authInfo?.providers || [];

    const returnUrl = getSafeReturnUrl(
        (location.state?.from as string) || searchParams.get('returnUrl'),
        routes.newFlow,
    );

    return (
        <div className="relative flex h-dvh w-full items-center justify-center">
            <LanguageSelect className="absolute top-4 right-4 z-10 w-auto" />
            <div className="h-dvh w-full lg:grid lg:grid-cols-2">
                <div className="flex items-center justify-center px-4 py-12">
                    {!isLoading ? (
                        <LoginForm
                            providers={authProviders}
                            returnUrl={returnUrl}
                        />
                    ) : (
                        <Spinner
                            className="size-16"
                            variant="circle"
                        />
                    )}
                </div>
                <div className="from-primary/20 via-primary/10 to-background hidden bg-linear-to-br lg:flex">
                    <Logo
                        aria-hidden
                        className="text-foreground m-auto h-24 w-auto"
                    />
                </div>
            </div>
        </div>
    );
}

export default Login;
