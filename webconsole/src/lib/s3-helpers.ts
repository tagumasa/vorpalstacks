/**
 * S3 file utility helpers.
 *
 * MIME type inference, text/image detection, and object row normalisation.
 * Extracted from s3.tsx for reuse and to reduce the page component size.
 */
import type { CommonPrefix, Object$ } from "@/gen/s3_pb";
import { formatBytes } from "@/lib/format";

// ─── Constants ──────────────────────────────────────────────────

export const MAX_UPLOAD_BYTES = 100 * 1024 * 1024;
export const MAX_TEXT_PREVIEW = 256 * 1024;
export const MAX_IMAGE_PREVIEW = 5 * 1024 * 1024;
export const LIST_PAGE_SIZE = 100;

const TEXT_EXTENSIONS = new Set([
  ".txt", ".json", ".csv", ".md", ".xml", ".log", ".yaml", ".yml",
  ".toml", ".ini", ".cfg", ".conf", ".html", ".css", ".js", ".ts",
  ".go", ".py", ".rs", ".java", ".sh", ".bash", ".zsh", ".env",
]);

const IMAGE_EXTENSIONS = new Set([
  ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico", ".bmp",
]);

const MIME_MAP: Record<string, string> = {
  ".json": "application/json", ".csv": "text/csv", ".xml": "application/xml",
  ".html": "text/html", ".css": "text/css", ".js": "application/javascript",
  ".ts": "application/typescript", ".md": "text/markdown", ".yaml": "application/x-yaml",
  ".yml": "application/x-yaml", ".txt": "text/plain", ".png": "image/png",
  ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".gif": "image/gif",
  ".svg": "image/svg+xml", ".webp": "image/webp", ".pdf": "application/pdf",
  ".ico": "image/x-icon", ".bmp": "image/bmp",
};

// ─── Helpers ────────────────────────────────────────────────────

function getExtension(key: string): string {
  const lastDot = key.lastIndexOf(".");
  if (lastDot === -1) return "";
  return key.slice(lastDot).toLowerCase();
}

/** Infer Content-Type from filename extension. */
export function inferContentType(filename: string): string {
  return MIME_MAP[getExtension(filename)] ?? "application/octet-stream";
}

/** Check if a key has a text-like extension. */
export function isTextFile(key: string): boolean {
  return TEXT_EXTENSIONS.has(getExtension(key));
}

/** Check if a key has an image extension. */
export function isImageFile(key: string): boolean {
  return IMAGE_EXTENSIONS.has(getExtension(key));
}

/** Extract file name from a full S3 object key. */
export function fileNameFromKey(key: string): string {
  const idx = key.lastIndexOf("/");
  return idx === -1 ? key : key.slice(idx + 1);
}

/** Remove trailing slash from a prefix. */
export function stripTrailingSlash(s: string): string {
  return s.endsWith("/") ? s.slice(0, -1) : s;
}

// ─── Object row type ────────────────────────────────────────────

/** Unified row type for the S3 object table (folders + files). */
export interface ObjectRow {
  id: string;
  key: string;
  displayKey: string;
  size: string;
  lastModified: string;
  storageClass: string;
  isFolder: boolean;
  rawSize: bigint;
}

/** Convert raw API response to unified ObjectRow list. */
export function toObjectRows(
  prefixes: CommonPrefix[],
  objects: Object$[],
  currentPrefix: string,
): ObjectRow[] {
  const folders: ObjectRow[] = prefixes.map((p) => ({
    id: `folder:${p.prefix}`,
    key: p.prefix,
    displayKey: stripTrailingSlash(p.prefix.slice(currentPrefix.length)),
    size: "—",
    lastModified: "—",
    storageClass: "—",
    isFolder: true,
    rawSize: 0n,
  }));
  const files: ObjectRow[] = objects
    .filter((o) => o.key !== currentPrefix)
    .map((o) => ({
      id: `file:${o.key}`,
      key: o.key,
      displayKey: o.key.slice(currentPrefix.length),
      size: formatBytes(o.size),
      lastModified: o.lastmodified,
      storageClass: String(o.storageclass ?? "STANDARD"),
      isFolder: false,
      rawSize: o.size,
    }));
  return [...folders, ...files];
}
