// Edits for a GraphQL input type whose secret fields `S` are tracked as
// SecretEdit instead of their raw nullable string.
export type EditsWithSecrets<T, S extends keyof T> = Partial<Omit<T, S> & Record<S, SecretEdit>>;

// Secret settings are write-only: the server only reports whether a value is
// stored (`*Set`). A pending edit is either a newly typed value or an explicit
// request to clear the stored one; no edit means "keep what is stored".
export type SecretEdit = { action: 'clear' } | { action: 'set'; value: string };

export const clearSecret: SecretEdit = { action: 'clear' };

export const isSecretDirty = (edit: SecretEdit | undefined): boolean =>
    edit?.action === 'clear' || (edit?.action === 'set' && edit.value !== '');

// Mutation value for a secret input: `undefined` keeps the stored secret, `''`
// clears it, anything else replaces it.
export const resolveSecret = (edit: SecretEdit | undefined): string | undefined => {
    if (edit?.action === 'clear') {
        return '';
    }

    if (edit?.action === 'set' && edit.value !== '') {
        return edit.value;
    }

    return undefined;
};

// True when any pending edit differs from the server value. Secret fields have
// no server counterpart, so they only count once typed or marked for clearing.
export const hasPendingEdits = (
    server: object,
    edits: object | undefined,
    secretFields: readonly string[] = [],
): boolean =>
    Object.entries(edits ?? {}).some(([key, value]) => {
        if (secretFields.includes(key)) {
            return isSecretDirty(value as SecretEdit | undefined);
        }

        return value !== undefined && value !== (server as Record<string, unknown>)[key];
    });
