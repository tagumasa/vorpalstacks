/**
 * Inspector components — shared 3-panel inspector pattern used by all service pages.
 *
 * Layout:  Toolbar (breadcrumb + actions)
 *          ──────────────────────────────
 *          Table (with optional checkbox column)
 *          ═══════ draggable splitter ═══════
 *          Detail panel (tabs + body)
 *
 * Extracted from S3/DynamoDB implementations for DRY reuse.
 */
import { type ReactNode, useState, useCallback } from "react";
import type { TFunction } from "i18next";
import type { ColumnDef } from "@tanstack/react-table";

// ─── Checkbox Column Builder ────────────────────────────────────

/**
 * Creates a checkbox selection column for a DataTable.
 * Supports select-all (header checkbox) and individual row selection.
 *
 * @param getRowId - Must return the same row ID string used by the DataTable's
 *                    getRowId prop. This is critical — the checkbox toggle relies
 *                    on matching IDs in the caller's selection Set.
 */
export function checkboxColumn<T>(
  selectedIds: Set<string>,
  onToggle: (id: string) => void,
  onToggleAll: () => void,
  allIds: string[],
  t: TFunction,
  getRowId: (row: T) => string,
): ColumnDef<T, any> {
  const allSelected = allIds.length > 0 && allIds.every((id) => selectedIds.has(id));
  const someSelected = !allSelected && allIds.some((id) => selectedIds.has(id));
  return {
    id: "select",
    size: 40,
    header: () => (
      <input
        type="checkbox"
        checked={allSelected}
        ref={(el) => {
          if (el) el.indeterminate = someSelected;
        }}
        onChange={onToggleAll}
        title={t("common.selectAll")}
        onClick={(e) => e.stopPropagation()}
      />
    ),
    cell: ({ row }) => {
      const rid = getRowId(row.original);
      return (
        <input
          type="checkbox"
          checked={selectedIds.has(rid)}
          onChange={() => onToggle(rid)}
          onClick={(e) => e.stopPropagation()}
        />
      );
    },
  };
}

// ─── Selection Helpers ──────────────────────────────────────────

/** Toggle a single item in a Set (immutable). */
export function toggleInSet<T>(set: Set<T>, item: T): Set<T> {
  const next = new Set(set);
  if (next.has(item)) next.delete(item);
  else next.add(item);
  return next;
}

/** Toggle all items — if all selected, deselect all; otherwise select all. */
export function toggleAll<T>(set: Set<T>, all: T[]): Set<T> {
  if (set.size === all.length && all.length > 0) return new Set();
  return new Set(all);
}

// ─── Breadcrumb ─────────────────────────────────────────────────

interface BreadcrumbPart {
  label: string;
  onClick?: () => void;
}

/** Renders a breadcrumb trail: Home / Level1 / Current */
export function Breadcrumb({ parts }: { parts: BreadcrumbPart[] }) {
  return (
    <div className="toolbar-breadcrumb">
      {parts.map((part, i) => {
        const isLast = i === parts.length - 1;
        return (
          <span key={i}>
            {i > 0 && <span className="breadcrumb-sep">/</span>}
            {isLast || !part.onClick ? (
              <span className="breadcrumb-current">{part.label}</span>
            ) : (
              <button className="breadcrumb-link" onClick={part.onClick}>
                {part.label}
              </button>
            )}
          </span>
        );
      })}
    </div>
  );
}

/** Selection count badge shown in toolbar. */
export function SelectionBadge({ count, label }: { count: number; label: string }) {
  if (count === 0) return null;
  return <span className="selection-count">{label}</span>;
}

// ─── Detail Panel ───────────────────────────────────────────────

export interface DetailTabDef {
  key: string;
  label: string;
}

interface DetailPanelProps {
  title: string;
  titleIcon?: string;
  tabs: DetailTabDef[];
  activeTab: string;
  onTabChange: (key: string) => void;
  actions?: ReactNode;
  children: ReactNode;
  footer?: ReactNode;
}

/** Bottom detail panel with tab bar, content body, and action buttons. */
export function DetailPanel({
  title,
  titleIcon,
  tabs,
  activeTab,
  onTabChange,
  actions,
  children,
  footer,
}: DetailPanelProps) {
  return (
    <div className="detail-panel-content">
      <div className="detail-header">
        <span className="detail-title">
          {titleIcon && `${titleIcon} `}{title}
        </span>
        {tabs.length > 1 && (
          <div className="detail-tabs">
            {tabs.map((tab) => (
              <button
                key={tab.key}
                className={`detail-tab ${activeTab === tab.key ? "active" : ""}`}
                onClick={() => onTabChange(tab.key)}
              >
                {tab.label}
              </button>
            ))}
          </div>
        )}
        {actions && <div className="detail-actions">{actions}</div>}
      </div>
      {footer}
      <div className="detail-body">
        {children}
      </div>
    </div>
  );
}

/** Empty detail panel placeholder. */
export function DetailEmpty({ message }: { message: string }) {
  return <div className="detail-panel-empty">{message}</div>;
}

// ─── useSelection Hook ──────────────────────────────────────────

/**
 * Hook for managing row selection state (checkbox multi-select).
 * Returns selected IDs and toggle helpers.
 */
export function useSelection<T extends string>(initial?: T[]) {
  const [selected, setSelected] = useState<Set<T>>(new Set(initial ?? []));

  const toggle = useCallback((id: T) => {
    setSelected((prev) => toggleInSet(prev, id));
  }, []);

  const toggleAll = useCallback((all: T[]) => {
    setSelected((prev) => {
      if (prev.size === all.length && all.length > 0) return new Set();
      return new Set(all);
    });
  }, []);

  const clear = useCallback(() => setSelected(new Set()), []);

  return { selected, toggle, toggleAll, clear, setSelected };
}
