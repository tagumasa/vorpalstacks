/** localStorage key for the JWT access token. */
const ACCESS_TOKEN_KEY = "vs_access_token";
/** localStorage key for the JWT refresh token. */
const REFRESH_TOKEN_KEY = "vs_refresh_token";
/** localStorage key for the JWT ID token. */
const ID_TOKEN_KEY = "vs_id_token";

/** Returns the stored access token, or null if not authenticated. */
export function getToken(): string | null {
  return localStorage.getItem(ACCESS_TOKEN_KEY);
}

/** Returns the stored refresh token, or null if not available. */
export function getRefreshToken(): string | null {
  return localStorage.getItem(REFRESH_TOKEN_KEY);
}

/**
 * Persists authentication tokens to localStorage.
 * Called after successful Login, RefreshToken, or InitialSetup RPC responses.
 */
export function setTokens(accessToken: string, refreshToken: string, idToken: string) {
  localStorage.setItem(ACCESS_TOKEN_KEY, accessToken);
  localStorage.setItem(REFRESH_TOKEN_KEY, refreshToken);
  localStorage.setItem(ID_TOKEN_KEY, idToken);
}

/** Removes all stored authentication tokens from localStorage. */
export function clearTokens() {
  localStorage.removeItem(ACCESS_TOKEN_KEY);
  localStorage.removeItem(REFRESH_TOKEN_KEY);
  localStorage.removeItem(ID_TOKEN_KEY);
}

/** Returns true if an access token is present in storage. */
export function isAuthenticated(): boolean {
  return !!getToken();
}
