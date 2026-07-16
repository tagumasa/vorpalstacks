import { createConnectTransport } from "@connectrpc/connect-web";
import { createClient, type Interceptor } from "@connectrpc/connect";
import { getToken, getRefreshToken, setTokens, clearTokens } from "./auth";
import { recordRequest, recordError, recordLatency } from "./metrics-store";
import { AdminAuthService } from "@/gen/admin_auth_pb";

/** Current region injected into every outgoing request header. */
let currentRegion = "us-east-1";

/** Update the region used by the region interceptor. */
export function setTransportRegion(r: string) {
  currentRegion = r;
}

/** Returns the application base path by reading the <base> tag or falling back to /webconsole/. */
function getBasePath(): string {
  return document.querySelector("base")?.getAttribute("href") ?? "/webconsole/";
}

/**
 * Shared in-flight refresh promise. When non-null, concurrent 401 responses
 * all await this single promise instead of racing to logout. The promise
 * resolves with the new access token on success or rejects on failure.
 */
let refreshPromise: Promise<string> | null = null;

/**
 * Region interceptor: attaches X-Region header so the backend
 * selects the correct region-specific Pebble storage.
 */
const regionInterceptor: Interceptor = (next) => async (req) => {
  req.header.set("X-Aws-Region", currentRegion);
  return next(req);
};

/**
 * Telemetry interceptor: tracks request count, error count, and latency
 * in the singleton MetricsStore for the StatusBar.
 */
const telemetryInterceptor: Interceptor = (next) => async (req) => {
  recordRequest();
  const start = performance.now();
  try {
    return await next(req);
  } catch (err: unknown) {
    recordError();
    throw err;
  } finally {
    recordLatency(performance.now() - start);
  }
};

/**
 * Transport for the refresh RPC itself — bypasses authInterceptor to
 * prevent self-deadlock when the refresh token is expired and the
 * refresh RPC returns 401. Without this, the interceptor would try to
 * join refreshPromise (which is itself), causing a permanent hang.
 */
const refreshTransport = createConnectTransport({
  baseUrl: "/",
  interceptors: [regionInterceptor, telemetryInterceptor],
});

/**
 * Performs a token refresh RPC and persists the new credentials.
 * Resolves with the new access token. Throws on failure.
 */
async function performRefresh(): Promise<string> {
  const refreshToken = getRefreshToken();
  if (!refreshToken) throw new Error("no refresh token");
  const refreshClient = createClient(AdminAuthService, refreshTransport);
  const res = await refreshClient.refreshToken({ refreshToken });
  setTokens(res.accessToken, res.refreshToken, res.idToken);
  return res.accessToken;
}

/**
 * Auth interceptor: attaches Authorization header with JWT bearer token
 * to every outgoing request. On 401 response, coalesces all concurrent
 * unauthenticated errors onto a single refresh promise. When refresh
 * succeeds, retries the original request with the new token. If refresh
 * fails, clears stored tokens and redirects to the login page.
 */
const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getToken();
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  try {
    return await next(req);
  } catch (err: unknown) {
    if (!(err instanceof Error && "code" in err && (err as { code: string }).code === "unauthenticated")) {
      throw err;
    }
    // 401: start or join the in-flight refresh
    if (!refreshPromise) {
      refreshPromise = performRefresh().finally(() => {
        refreshPromise = null;
      });
    }
    try {
      const newToken = await refreshPromise;
      req.header.set("Authorization", `Bearer ${newToken}`);
      return next(req);
    } catch {
      clearTokens();
      window.location.href = getBasePath() + "login";
      throw err;
    }
  }
};

/**
 * Creates a Connect transport instance with JWT auth and region header interceptors.
 */
export function createTransport() {
  return createConnectTransport({
    baseUrl: "/",
    interceptors: [authInterceptor, regionInterceptor, telemetryInterceptor],
  });
}

/** Shared transport instance used by all Connect RPC clients. */
export const transport = createTransport();
