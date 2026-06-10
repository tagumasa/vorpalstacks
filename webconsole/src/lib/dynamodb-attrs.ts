/**
 * DynamoDB AttributeValue conversion helpers.
 *
 * Pure functions for converting between AttributeValue, AttrRow (structured
 * form representation), and plain JSON. Extracted from dynamodb.tsx for reuse
 * and to reduce the page component size.
 */
import { create } from "@bufbuild/protobuf";
import type { AttributeValue, KeySchemaElement } from "@/gen/dynamodb_pb";
import { AttributeValueSchema } from "@/gen/dynamodb_pb";

// ─── Attribute type ──────────────────────────────────────────────

/** Extended attribute type including complex types. */
export type AttrType = "S" | "N" | "B" | "BOOL" | "NULL" | "SS" | "NS" | "L" | "M";

/** Structured form attribute row. */
export interface AttrRow {
  name: string;
  type: AttrType;
  value: string;
  children?: AttrRow[];
}

/** i18n key for each attribute type label. */
export const TYPE_LABEL_KEYS: Record<AttrType, string> = {
  S: "typeLabelS", N: "typeLabelN", B: "typeLabelB",
  BOOL: "typeLabelBOOL", NULL: "typeLabelNULL",
  SS: "typeLabelSS", NS: "typeLabelNS", L: "typeLabelL", M: "typeLabelM",
};

export const ATTR_TYPES: AttrType[] = ["S", "N", "B", "BOOL", "NULL", "SS", "NS", "L", "M"];

// ─── Formatting ──────────────────────────────────────────────────

/** Render an AttributeValue as a short string (for table cells). */
export function fmtAttr(v: AttributeValue): string {
  if (v.s) return v.s;
  if (v.n) return v.n;
  if (v.bool !== undefined) return String(v.bool);
  if (v.null) return "null";
  if (v.ss?.length) return `[${v.ss.join(", ")}]`;
  if (v.ns?.length) return `[${v.ns.join(", ")}]`;
  if (v.l?.length) return `[${v.l.map((e) => fmtAttr(e)).join(", ")}]`;
  if (v.m && Object.keys(v.m).length) return `{${Object.entries(v.m).map(([k, av]) => `${k}: ${fmtAttr(av)}`).join(", ")}}`;
  if (v.b !== undefined && v.b.length > 0) return `B(${v.b.length})`;
  if (v.bs?.length) return `BS(${v.bs.length})`;
  return "—";
}

/** Attribute type label for a single AttributeValue. */
export function attrTypeLabel(v: AttributeValue): string {
  if (v.s) return "S";
  if (v.n) return "N";
  if (v.m && Object.keys(v.m).length) return "M";
  if (v.l?.length) return "L";
  if (v.bool === true || v.bool === false) return "BOOL";
  if (v.null === true) return "NULL";
  if (v.ss?.length) return "SS";
  if (v.ns?.length) return "NS";
  if (v.b && v.b.length > 0) return "B";
  if (v.bs?.length) return "BS";
  return "?";
}

/** Placeholder hint for each type. */
export function typePlaceholder(t: AttrType): string {
  switch (t) {
    case "S": return "hello";
    case "N": return "42";
    case "B": return "binary data";
    case "BOOL": return "true or false";
    case "NULL": return "(no value needed)";
    case "SS": return "a, b, c";
    case "NS": return "1, 2, 3";
    case "L": return '[1, "hello", true]';
    case "M": return '{"key": "value"}';
  }
}

// ─── Key extraction ──────────────────────────────────────────────

/** Extract PK and optionally SK from an item using the key schema. */
export function extractKey(
  item: Record<string, AttributeValue>,
  keySchema: KeySchemaElement[],
): Record<string, AttributeValue> {
  const key: Record<string, AttributeValue> = {};
  for (const ks of keySchema) {
    const v = item[ks.attributename];
    if (v) key[ks.attributename] = v;
  }
  return key;
}

// ─── AttributeValue ↔ JSON ──────────────────────────────────────

/** Convert a plain JSON value to an AttributeValue. */
export function jsonToAV(v: unknown): AttributeValue {
  if (v === null) return create(AttributeValueSchema, { null: true });
  if (typeof v === "string") return create(AttributeValueSchema, { s: v });
  if (typeof v === "number") return create(AttributeValueSchema, { n: String(v) });
  if (typeof v === "boolean") return create(AttributeValueSchema, { bool: v });
  if (Array.isArray(v)) return create(AttributeValueSchema, { l: v.map(jsonToAV) });
  if (typeof v === "object") {
    const m: Record<string, AttributeValue> = {};
    for (const [k, val] of Object.entries(v as Record<string, unknown>)) {
      m[k] = jsonToAV(val);
    }
    return create(AttributeValueSchema, { m });
  }
  return create(AttributeValueSchema, { s: String(v) });
}

/** Convert an AttributeValue back to a plain JSON value. */
export function avToPlain(v: AttributeValue): unknown {
  if (v.s !== undefined && v.s !== "") return v.s;
  if (v.n) return Number(v.n);
  if (v.bool !== undefined) return v.bool;
  if (v.null) return null;
  if (v.ss?.length) return v.ss.slice();
  if (v.ns?.length) return v.ns.map(Number);
  if (v.l?.length) return v.l.map(avToPlain);
  if (v.m && Object.keys(v.m).length) {
    const obj: Record<string, unknown> = {};
    for (const [k, val] of Object.entries(v.m)) obj[k] = avToPlain(val);
    return obj;
  }
  if (v.b !== undefined && v.b.length > 0) return "[binary]";
  return "";
}

// ─── AttrRow conversions ─────────────────────────────────────────

/** Build an AttributeValue from type and string value. */
export function buildAV(type: AttrType, value: string): AttributeValue {
  switch (type) {
    case "S":
      return create(AttributeValueSchema, { s: value });
    case "N":
      return create(AttributeValueSchema, { n: value });
    case "B":
      return create(AttributeValueSchema, { b: new TextEncoder().encode(value) });
    case "BOOL":
      return create(AttributeValueSchema, { bool: value === "true" });
    case "NULL":
      return create(AttributeValueSchema, { null: true });
    case "SS":
      return create(AttributeValueSchema, { ss: value.split(",").map(s => s.trim()).filter(Boolean) });
    case "NS":
      return create(AttributeValueSchema, { ns: value.split(",").map(s => s.trim()).filter(Boolean) });
    case "L": {
      const parsed = JSON.parse(value);
      if (!Array.isArray(parsed)) throw new Error("List value must be a JSON array");
      return create(AttributeValueSchema, { l: parsed.map((item: unknown) => jsonToAV(item)) });
    }
    case "M": {
      const parsed = JSON.parse(value);
      if (typeof parsed !== "object" || Array.isArray(parsed)) throw new Error("Map value must be a JSON object");
      const m: Record<string, AttributeValue> = {};
      for (const [k, v] of Object.entries(parsed)) {
        m[k] = jsonToAV(v);
      }
      return create(AttributeValueSchema, { m });
    }
    default:
      return create(AttributeValueSchema, { s: value });
  }
}

/** Convert an AttrRow (possibly with children) to an AttributeValue. */
export function rowToAV(row: AttrRow): AttributeValue {
  if (row.type === "NULL") return create(AttributeValueSchema, { null: true });
  if (row.type === "L") {
    const items = (row.children ?? []).map(rowToAV);
    return create(AttributeValueSchema, { l: items });
  }
  if (row.type === "M") {
    const m: Record<string, AttributeValue> = {};
    for (const child of row.children ?? []) {
      if (child.name.trim()) m[child.name] = rowToAV(child);
    }
    return create(AttributeValueSchema, { m });
  }
  return buildAV(row.type, row.value);
}

/** Convert an AttributeValue to an AttrRow (for structured editing). */
export function avToRow(name: string, av: AttributeValue): AttrRow {
  const label = attrTypeLabel(av) as AttrType;
  if (label === "L" && av.l?.length) {
    return { name, type: "L", value: "", children: av.l.map((item, i) => avToRow(String(i), item)) };
  }
  if (label === "M" && av.m && Object.keys(av.m).length) {
    return { name, type: "M", value: "", children: Object.entries(av.m).map(([k, v]) => avToRow(k, v)) };
  }
  return { name, type: label, value: avToValueString(av) };
}

/** Convert an AttrRow tree back to a plain JSON value (for tab sync). */
export function rowToPlain(row: AttrRow): unknown {
  if (row.type === "NULL") return null;
  if (row.type === "BOOL") return row.value === "true";
  if (row.type === "L") return (row.children ?? []).map(rowToPlain);
  if (row.type === "M") {
    const obj: Record<string, unknown> = {};
    for (const child of row.children ?? []) {
      if (child.name.trim()) obj[child.name] = rowToPlain(child);
    }
    return obj;
  }
  try { return JSON.parse(row.value); } catch { return row.value; }
}

/** Convert a plain JSON value to an AttrRow (for JSON→structured conversion). */
export function plainToRow(name: string, v: unknown): AttrRow {
  if (v === null) return { name, type: "NULL", value: "" };
  if (typeof v === "boolean") return { name, type: "BOOL", value: String(v) };
  if (typeof v === "number") return { name, type: "N", value: String(v) };
  if (typeof v === "string") return { name, type: "S", value: v };
  if (Array.isArray(v)) return { name, type: "L", value: "", children: v.map((item, i) => plainToRow(String(i), item)) };
  if (typeof v === "object") {
    return { name, type: "M", value: "", children: Object.entries(v as Record<string, unknown>).map(([k, val]) => plainToRow(k, val)) };
  }
  return { name, type: "S", value: String(v) };
}

/** Convert an AttributeValue to a plain string (for structured form value field). */
export function avToValueString(v: AttributeValue): string {
  if (v.s !== undefined && v.s !== "") return v.s;
  if (v.n) return v.n;
  if (v.bool !== undefined) return String(v.bool);
  if (v.null) return "";
  if (v.ss?.length) return v.ss.join(", ");
  if (v.ns?.length) return v.ns.join(", ");
  if (v.l?.length) return JSON.stringify(v.l.map(avToPlain));
  if (v.m && Object.keys(v.m).length) return JSON.stringify(Object.fromEntries(Object.entries(v.m).map(([k, av]) => [k, avToPlain(av)])));
  return "";
}
