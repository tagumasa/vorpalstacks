import { createConnectTransport } from "@connectrpc/connect-web";
import type { Interceptor } from "@connectrpc/connect";
import { getToken, clearTokens } from "./auth";

/** Current region injected into every outgoing request header. */
let currentRegion = "us-east-1";

/** Update the region used by the region interceptor. */
export function setTransportRegion(r: string) {
  currentRegion = r;
}

/**
 * Creates a Connect transport instance with JWT auth and region header interceptors.
 */
export function createTransport() {
  return createConnectTransport({
    baseUrl: "/",
    interceptors: [authInterceptor, regionInterceptor],
  });
}

/**
 * Auth interceptor: attaches Authorization header with JWT bearer token
 * to every outgoing request. On 401 response, clears stored tokens
 * and redirects to the login page.
 */
const authInterceptor: Interceptor = (next) => async (req) => {
  const token = getToken();
  if (token) {
    req.header.set("Authorization", `Bearer ${token}`);
  }
  try {
    return await next(req);
  } catch (err: any) {
    if (err?.code === "unauthenticated") {
      clearTokens();
      window.location.href = "/login";
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

/** Shared transport instance used by all Connect RPC clients. */
export const transport = createTransport();
