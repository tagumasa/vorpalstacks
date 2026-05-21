/**
 * Generic data table component built on TanStack Table.
 * Supports column sorting, pagination, global text filtering, row selection,
 * and data export (JSON/CSV).
 */
import { useState, useCallback } from "react";
import { useTranslation } from "react-i18next";
import type { ColumnDef, SortingState, FilterFn } from "@tanstack/react-table";
import { useReactTable, getCoreRowModel, getSortedRowModel, getFilteredRowModel, getPaginationRowModel, flexRender } from "@tanstack/react-table";

/** Case-insensitive substring global filter function. */
const globalTextFilter: FilterFn<any> = (row, _columnId, filterValue: string) => {
  const needle = String(filterValue).toLowerCase();
  if (!needle) return true;
  return Object.values(row.original).some((val) => {
    if (val == null) return false;
    return String(val).toLowerCase().includes(needle);
  });
};

/** Props for the DataTable component. */
interface DataTableProps<T> {
  columns: ColumnDef<T, any>[];
  data: T[];
  onRowClick?: (row: T) => void;
  selectedId?: string;
  getRowId?: (row: T) => string;
  pageSize?: number;
  emptyMessage?: string;
  /** Filename prefix for export (e.g. "ssm", "rds"). */
  exportName?: string;
}

type ExportFormat = "json" | "csv";

function downloadBlob(content: string, filename: string, mime: string) {
  const blob = new Blob([content], { type: mime });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  a.click();
  URL.revokeObjectURL(url);
}

/** Renders a paginated, sortable, filterable table with optional row click selection and data export. */
export function DataTable<T>({ columns, data, onRowClick, selectedId, getRowId, pageSize = 50, emptyMessage, exportName }: DataTableProps<T>) {
  const { t } = useTranslation();
  const [sorting, setSorting] = useState<SortingState>([]);
  const [globalFilter, setGlobalFilter] = useState("");
  const [showExportMenu, setShowExportMenu] = useState(false);

  const table = useReactTable({
    data,
    columns,
    state: { sorting, globalFilter },
    onSortingChange: setSorting,
    onGlobalFilterChange: setGlobalFilter,
    getCoreRowModel: getCoreRowModel(),
    getSortedRowModel: getSortedRowModel(),
    getFilteredRowModel: getFilteredRowModel(),
    getPaginationRowModel: getPaginationRowModel(),
    globalFilterFn: globalTextFilter,
    getRowId,
    initialState: { pagination: { pageSize } },
  });

  const rows = table.getRowModel().rows;

  const handleExport = useCallback((fmt: ExportFormat) => {
    const filteredRows = table.getFilteredRowModel().rows;
    const colIds = columns.map((c) => ("accessorKey" in c ? String(c.accessorKey) : "")).filter(Boolean);
    const rowData = filteredRows.map((r) => {
      const obj: Record<string, unknown> = {};
      for (const key of colIds) {
        obj[key] = (r.original as Record<string, unknown>)[key];
      }
      return obj;
    });
    const date = new Date().toISOString().slice(0, 10);
    const prefix = exportName || "export";
    if (fmt === "json") {
      downloadBlob(JSON.stringify(rowData, null, 2), `${prefix}_${date}.json`, "application/json");
    } else {
      const header = colIds.join(",");
      const csvRows = rowData.map((r) => colIds.map((k) => JSON.stringify(r[k] ?? "")).join(","));
      downloadBlob([header, ...csvRows].join("\n"), `${prefix}_${date}.csv`, "text/csv");
    }
    setShowExportMenu(false);
  }, [table, columns, exportName]);

  return (
    <div className="data-table-wrapper">
      {data.length > 0 && (
        <div className="data-table-toolbar">
          <input
            type="text"
            className="data-table-filter"
            value={globalFilter}
            onChange={(e) => setGlobalFilter(e.target.value)}
            placeholder={t("common.filter")}
            aria-label={t("common.filter")}
          />
          {globalFilter && (
            <span className="data-table-filter-count">
              {rows.length}/{data.length}
            </span>
          )}
          {exportName && (
            <div className="data-table-export">
              <button className="btn btn-sm" onClick={() => setShowExportMenu((v) => !v)}>{t("common.export")}</button>
              {showExportMenu && (
                <div className="data-table-export-menu">
                  <button onClick={() => handleExport("json")}>{t("common.exportJson")}</button>
                  <button onClick={() => handleExport("csv")}>{t("common.exportCsv")}</button>
                </div>
              )}
            </div>
          )}
        </div>
      )}
      <table className="data-table">
        <thead>
          {table.getHeaderGroups().map((hg) => (
            <tr key={hg.id}>
              {hg.headers.map((header) => (
                <th
                  key={header.id}
                  onClick={header.column.getToggleSortingHandler()}
                  className={header.column.getIsSorted() ? "sorted" : ""}
                >
                  {header.isPlaceholder ? null : (
                    <>
                      {flexRender(header.column.columnDef.header, header.getContext())}
                      {header.column.getIsSorted() === "asc" && <span className="sort-asc" />}
                      {header.column.getIsSorted() === "desc" && <span className="sort-desc" />}
                    </>
                  )}
                </th>
              ))}
            </tr>
          ))}
        </thead>
        <tbody>
          {rows.length === 0 ? (
            <tr className="empty-row">
              <td colSpan={columns.length} className="empty-cell">
                {emptyMessage ?? t("common.noData")}
              </td>
            </tr>
          ) : (
            rows.map((row) => {
              const rowId = getRowId ? getRowId(row.original) : row.id;
              return (
                <tr
                  key={row.id}
                  className={selectedId === rowId ? "selected" : ""}
                  onClick={() => onRowClick?.(row.original)}
                >
                  {row.getVisibleCells().map((cell) => (
                    <td key={cell.id}>
                      {flexRender(cell.column.columnDef.cell, cell.getContext())}
                    </td>
                  ))}
                </tr>
              );
            })
          )}
        </tbody>
      </table>
      {table.getPageCount() > 1 && (
        <div className="data-table-pagination">
          <button
            onClick={() => table.previousPage()}
            disabled={!table.getCanPreviousPage()}
          >
            {t("common.prev")}
          </button>
          <span className="page-info">
            {table.getState().pagination.pageIndex + 1} / {table.getPageCount()}
          </span>
          <button
            onClick={() => table.nextPage()}
            disabled={!table.getCanNextPage()}
          >
            {t("common.next")}
          </button>
        </div>
      )}
    </div>
  );
}
