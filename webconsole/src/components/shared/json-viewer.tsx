/**
 * Collapsible JSON tree viewer. Renders key-value pairs with expand/collapse
 * for nested objects and arrays. Used in the service detail panel.
 */
import { useState } from "react";

interface JsonViewerProps {
  data: unknown;
}

export function JsonViewer({ data }: JsonViewerProps) {
  if (data === null || data === undefined) {
    return <span className="json-null">null</span>;
  }

  if (typeof data !== "object") {
    return <span className="json-primitive">{String(data)}</span>;
  }

  const entries = Object.entries(data as Record<string, unknown>);

  return (
    <div className="json-viewer">
      {entries.map(([key, value]) => (
        <JsonEntry key={key} keyName={key} value={value} depth={0} />
      ))}
    </div>
  );
}

/** Single entry in the JSON tree. Objects/arrays are collapsible. */
function JsonEntry({ keyName, value, depth }: { keyName: string; value: unknown; depth: number }) {
  const [expanded, setExpanded] = useState(depth < 1);
  const isObject = value !== null && typeof value === "object";
  const isArray = Array.isArray(value);

  return (
    <div className="json-entry" style={{ paddingLeft: depth * 16 }}>
      {isObject ? (
        <>
          <div className="json-key-row" onClick={() => setExpanded(!expanded)}>
            <span className="json-toggle">{expanded ? "▼" : "▶"}</span>
            <span className="json-key">{keyName}</span>
            <span className="json-type">
              {isArray ? `Array[${(value as unknown[]).length}]` : "Object"}
            </span>
          </div>
          {expanded && (
            <div className="json-children">
              {Object.entries(value as Record<string, unknown>).map(([k, v]) => (
                <JsonEntry key={k} keyName={k} value={v} depth={depth + 1} />
              ))}
            </div>
          )}
        </>
      ) : (
        <div className="json-leaf">
          <span className="json-key">{keyName}:</span>
          <span className={`json-value${value === null ? " json-null" : ""}`}>
            {value === null ? "null" : value === true ? "true" : value === false ? "false" : String(value)}
          </span>
        </div>
      )}
    </div>
  );
}
