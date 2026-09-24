import { StatusType } from '@/graphql/types';

/** `flows` namespace keys for flow status labels; resolve with `t()` at render time. */
export const flowStatusLabelKeys = {
    [StatusType.Created]: 'status.created',
    [StatusType.Failed]: 'status.failed',
    [StatusType.Finished]: 'status.finished',
    [StatusType.Running]: 'status.running',
    [StatusType.Waiting]: 'status.waiting',
} as const satisfies Record<StatusType, `status.${StatusType}`>;
