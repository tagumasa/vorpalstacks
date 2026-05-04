/**
 * Export utility for downloading data as JSON or CSV files.
 * Uses Blob + URL.createObjectURL + hidden <a> click pattern.
 */

export function exportJSON(rows: Record<string, unknown>[], filename: string) {
  const blob = new Blob([JSON.stringify(rows, null, 2)], { type: "application/json" });
  downloadBlob(blob, filename);
}

export function exportCSV(
  rows: Record<string, unknown>[],
  columns: { accessorKey?: string; header?: string }[],
  filename: string,
) {
  const headers = columns.map((c) => c.accessorKey ?? "");
  const headerRow = columns.map((c) => csvEscape(c.header ?? c.accessorKey ?? ""));
  const dataRows = rows.map((row) =>
    headers.map((key) => csvEscape(String(row[key] ?? ""))).join(","),
  );
  const csv = [headerRow.join(","), ...dataRows].join("\n");
  const blob = new Blob([csv], { type: "text/csv" });
  downloadBlob(blob, filename);
}

function csvEscape(val: string): string {
  if (val.includes(",") || val.includes('"') || val.includes("\n")) {
    return `"${val.replace(/"/g, '""')}"`;
  }
  return val;
}

function downloadBlob(blob: Blob, filename: string) {
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}
