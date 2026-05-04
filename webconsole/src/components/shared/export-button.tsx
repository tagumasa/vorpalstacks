/**
 * Export dropdown button for DataTable data.
 * Renders JSON/CSV export options in a small dropdown.
 */
import { useState, useRef, useEffect } from "react";
import { useTranslation } from "react-i18next";
import type { ColumnDef } from "@tanstack/react-table";
import { exportJSON, exportCSV } from "@/lib/export";

interface ExportButtonProps<T> {
  rows: T[];
  columns: ColumnDef<T, any>[];
  filenamePrefix: string;
}

export function ExportButton<T>({ rows, columns, filenamePrefix }: ExportButtonProps<T>) {
  const { t } = useTranslation();
  const [open, setOpen] = useState(false);
  const ref = useRef<HTMLDivElement>(null);

  useEffect(() => {
    const handler = (e: MouseEvent) => {
      if (ref.current && !ref.current.contains(e.target as Node)) {
        setOpen(false);
      }
    };
    document.addEventListener("click", handler);
    return () => document.removeEventListener("click", handler);
  }, []);

  const ts = new Date().toISOString().replace(/[:.]/g, "-").slice(0, 19);
  const accessorColumns = columns
    .filter((c) => "accessorKey" in c && c.accessorKey)
    .map((c) => ({
      accessorKey: (c as any).accessorKey as string,
      header:
        typeof c.header === "string"
          ? c.header
          : (c as any).accessorKey as string,
    }));

  const handleJSON = () => {
    const data = rows.map((row) =>
      Object.fromEntries(
        accessorColumns.map((c) => [c.accessorKey, (row as any)[c.accessorKey]]),
      ),
    );
    exportJSON(data, `${filenamePrefix}-${ts}.json`);
    setOpen(false);
  };

  const handleCSV = () => {
    const data = rows.map((row) =>
      Object.fromEntries(
        accessorColumns.map((c) => [c.accessorKey, (row as any)[c.accessorKey]]),
      ),
    );
    exportCSV(data, accessorColumns, `${filenamePrefix}-${ts}.csv`);
    setOpen(false);
  };

  return (
    <div ref={ref} style={{ position: "relative", display: "inline-block" }}>
      <button
        className="btn btn-secondary"
        onClick={() => setOpen((o) => !o)}
        style={{ fontSize: 11 }}
      >
        ⬇ {t("export.title")}
      </button>
      {open && (
        <div
          style={{
            position: "absolute",
            right: 0,
            top: "100%",
            zIndex: 100,
            background: "var(--bg-secondary)",
            border: "1px solid var(--border-dim)",
            borderRadius: 4,
            marginTop: 2,
            minWidth: 100,
          }}
        >
          <div
            className="dropdown-item"
            onClick={handleJSON}
            style={{ fontSize: 11 }}
          >
            JSON
          </div>
          <div
            className="dropdown-item"
            onClick={handleCSV}
            style={{ fontSize: 11 }}
          >
            CSV
          </div>
        </div>
      )}
    </div>
  );
}
