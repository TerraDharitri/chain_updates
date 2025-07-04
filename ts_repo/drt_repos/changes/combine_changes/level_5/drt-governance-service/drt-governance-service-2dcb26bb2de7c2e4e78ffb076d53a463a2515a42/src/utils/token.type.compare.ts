import { DcdtToken } from 'src/modules/tokens/models/dcdtToken.model';

export function isDcdtToken(
    token: DcdtToken,
): token is DcdtToken {
    return (token as DcdtToken).identifier !== undefined;
}

export function isDcdtTokenValid(token: DcdtToken): boolean {
    if (
        !token.identifier ||
        !token.decimals ||
        token.identifier === undefined ||
        token.decimals === undefined ||
        token.decimals === 0
    ) {
        return false;
    }
    return true;
}
