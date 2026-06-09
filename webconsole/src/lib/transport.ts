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

/** Set when a refresh is in-flight to prevent concurrent refresh attempts. */
let refreshing = false;

/**
 * Auth interceptor: attaches Authorization header with JWT bearer token
 * to every outgoing request. On 401 response, attempts to refresh the token
 * using the stored refresh token. If refresh succeeds, retries the original
 * request with the new token. If refresh fails, clears stored tokens
 * and redirects to the login page.
 */
const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getToken();
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  try {
    return await next(req);
  } catch (err: unknown) {
    if (err instanceof Error && "code" in err && (err as { code: string }).code === "unauthenticated") {
      if (refreshing) {
        clearTokens();
        window.location.href = getBasePath() + "login";
        throw err;
      }
      const refreshToken = getRefreshToken();
      if (refreshToken) {
        refreshing = true;
        try {
          const refreshClient = createClient(AdminAuthService, transport);
          const res = await refreshClient.refreshToken({ refreshToken });
          setTokens(res.accessToken, res.refreshToken, res.idToken);
          req.header.set("Authorization", `Bearer ${res.accessToken}`);
          refreshing = false;
          return next(req);
        } catch {
          refreshing = false;
        }
      }
      clearTokens();
      window.location.href = getBasePath() + "login";
    }
    throw err;
  }
};

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
