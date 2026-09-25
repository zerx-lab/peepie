import { cn } from '@/lib/utils';

interface LogoProps extends React.SVGProps<SVGSVGElement> {
    className?: string;
}

const LETTER_P = 'M20 38V108M20 60a22 22 0 1 0 44 0a22 22 0 1 0 -44 0Z';

// The wordmark's first "p" alone, for square slots too small for the full wordmark (collapsed sidebar).
export function LogoMark({ className, ...props }: LogoProps) {
    return (
        <svg
            aria-label="Peepie"
            className={cn(className)}
            fill="currentColor"
            role="img"
            viewBox="-6 25 96 96"
            xmlns="http://www.w3.org/2000/svg"
            {...props}
        >
            <path
                d={LETTER_P}
                fill="none"
                stroke="currentColor"
                strokeLinecap="round"
                strokeWidth={14}
            />
            <circle
                cx={42}
                cy={60}
                r={7}
            />
        </svg>
    );
}

// Wordmark: geometric lowercase "peepie"; the first p's bowl carries a pupil. Size it by height (`h-* w-auto`).
function Logo({ className, ...props }: LogoProps) {
    return (
        <svg
            aria-label="Peepie"
            className={cn(className)}
            fill="currentColor"
            role="img"
            viewBox="13 5 364 111"
            xmlns="http://www.w3.org/2000/svg"
            {...props}
        >
            <path
                d={`${LETTER_P}M90 60H134A22 22 0 1 0 127.56 75.56M160 60H204A22 22 0 1 0 197.56 75.56M230 38V108M230 60a22 22 0 1 0 44 0a22 22 0 1 0 -44 0ZM300 38V82M326 60H370A22 22 0 1 0 363.56 75.56`}
                fill="none"
                stroke="currentColor"
                strokeLinecap="round"
                strokeWidth={14}
            />
            <circle
                cx={42}
                cy={60}
                r={7}
            />
            <circle
                cx={300}
                cy={14}
                r={8.5}
            />
        </svg>
    );
}

export default Logo;
