/**
 * Shared layout and utility components for service pages.
 * Eliminates boilerplate across 29 service page implementations.
 */
import { type ReactNode, useMemo } from "react";
import { useTranslation } from "react-i18next";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import type { DescService } from "@bufbuild/protobuf";
import { DataTable } from "./data-table";
import { JsonViewer } from "./json-viewer";
import { Modal } from "./modal";
import { Splitter } from "./splitter";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator } from "@/lib/use-service-list";

// ─── Helpers ────────────────────────────────────────────────────

/** Extract a human-readable message from Connect RPC errors or plain Error. */
function formatError(err: unknown): string {
  if (err instanceof Error) return err.message;
  if (typeof err === "object" && err !== null && "message" in err) {
    const msg = (err as { message: unknown }).message;
    if (typeof msg === "string" && msg) return msg;
  }
  return String(err);
}

// ─── Layout ─────────────────────────────────────────────────────

interface ServicePageLayoutProps {
  icon: string;
  title: string;
  isLoading?: boolean;
  error?: unknown;
  count?: number;
  countLabel?: string;
  actions?: ReactNode;
  tabs?: { key: string; label: string; count: number }[];
  activeTab?: string;
  onTabChange?: (key: string) => void;
  children: ReactNode;
}

export function ServicePageLayout({
  icon,
  title,
  isLoading,
  error,
  count,
  countLabel,
  actions,
  tabs,
  activeTab,
  onTabChange,
  children,
}: ServicePageLayoutProps) {
  const { t } = useTranslation();

  if (isLoading) {
    return (
      <div className="content-area">
        <PageHeader icon={icon} title={title} />
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <PageHeader icon={icon} title={title} />
        <div className="error-state">{t("common.failedToLoad", { error: formatError(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <div className="page-header-row">
          <span className="page-icon">{icon}</span>
          <h1>{title}</h1>
          {count !== undefined && !tabs && (
            <span className="resource-count">
              {count} {countLabel ?? t("common.resources", { count })}
            </span>
          )}
        </div>
        {actions && <div className="page-actions">{actions}</div>}
        {tabs && (
          <div className="tab-bar">
            {tabs.map((tb) => (
              <button
                key={tb.key}
                className={`tab-btn ${activeTab === tb.key ? "active" : ""}`}
                onClick={() => onTabChange?.(tb.key)}
              >
                {tb.count > 0 ? `${tb.label} (${tb.count})` : tb.label}
              </button>
            ))}
          </div>
        )}
      </div>
      {children}
    </div>
  );
}

function PageHeader({ icon, title }: { icon: string; title: string }) {
  return (
    <div className="page-header">
      <span className="page-icon">{icon}</span>
      <h1>{title}</h1>
    </div>
  );
}

// ─── Split Pane ─────────────────────────────────────────────────

interface SplitPaneProps<T> {
  columns: ColumnDef<T, any>[];
  data: T[];
  getRowId: (row: T) => string;
  onRowClick: (row: T) => void;
  selectedId?: string;
  selected?: T | null;
  detailTitle?: string;
  onDetailClose: () => void;
  DetailComponent?: React.ComponentType<{ item: T }>;
}

export function SplitPane<T>({
  columns,
  data,
  getRowId,
  onRowClick,
  selectedId,
  selected,
  detailTitle,
  onDetailClose,
  DetailComponent,
}: SplitPaneProps<T>) {
  const serviceId = typeof window !== "undefined"
    ? window.location.pathname.split("/").pop() ?? "default"
    : "default";

  const table = (
    <DataTable
      columns={columns}
      data={data}
      getRowId={getRowId}
      onRowClick={onRowClick}
      selectedId={selectedId}
    />
  );

  if (!selected) {
    return table;
  }

  return (
    <Splitter
      initialSize={360}
      minSize={200}
      maxSize={600}
      storageKey={`vs-split-${serviceId}`}
    >
      {table}
      <div className="split-detail">
        <div className="detail-header">
          <h2>{detailTitle ?? ""}</h2>
          <button className="detail-close icon-close" onClick={onDetailClose} />
        </div>
        {DetailComponent ? (
          <DetailComponent item={selected} />
        ) : (
          <JsonViewer data={selected} />
        )}
      </div>
    </Splitter>
  );
}

// ─── Create Modal ───────────────────────────────────────────────

interface ServiceCreateModalProps {
  open: boolean;
  onClose: () => void;
  title: string;
  error?: Error | null;
  isPending?: boolean;
  onCreate: () => void;
  disabled?: boolean;
  children: ReactNode;
}

export function ServiceCreateModal({
  open,
  onClose,
  title,
  error,
  isPending,
  onCreate,
  disabled,
  children,
}: ServiceCreateModalProps) {
  const { t } = useTranslation();

  return (
    <Modal open={open} onClose={onClose}>
      <h2>{title}</h2>
      {error && <div className="modal-error">{formatError(error)}</div>}
      {children}
      <div className="modal-actions">
        <button className="btn btn-secondary" onClick={onClose}>
          {t("common.cancel")}
        </button>
        <button
          className="btn btn-primary"
          disabled={disabled || isPending}
          onClick={onCreate}
        >
          {isPending ? t("common.creating") : t("common.create")}
        </button>
      </div>
    </Modal>
  );
}

// ─── Delete Dialog ──────────────────────────────────────────────

interface ServiceDeleteDialogProps {
  open: boolean;
  title: string;
  name?: string;
  error?: Error | null;
  isPending?: boolean;
  onConfirm: () => void;
  onClose: () => void;
}

export function ServiceDeleteDialog({
  open,
  title,
  name,
  error,
  isPending,
  onConfirm,
  onClose,
}: ServiceDeleteDialogProps) {
  const { t } = useTranslation();

  return (
    <Modal open={open} onClose={onClose}>
      <h2>{title}</h2>
      <p>
        {t("confirm.confirmDelete")} <strong>{name}</strong>?
      </p>
      {error && <div className="modal-error">{formatError(error)}</div>}
      <div className="modal-actions">
        <button className="btn btn-secondary" onClick={onClose}>
          {t("common.cancel")}
        </button>
        <button
          className="btn btn-danger"
          disabled={isPending}
          onClick={onConfirm}
        >
          {isPending ? t("common.deleting") : t("common.delete")}
        </button>
      </div>
    </Modal>
  );
}

// ─── Cell Renderers ─────────────────────────────────────────────

export function MonoCell({ getValue }: { getValue: () => unknown }) {
  return <span className="cell-mono">{getValue() as string}</span>;
}

export function SmallMonoCell({ getValue }: { getValue: () => unknown }) {
  return (
    <span className="cell-mono cell-long">
      {getValue() as string}
    </span>
  );
}

export function DateCell({ getValue }: { getValue: () => unknown }) {
  const { i18n } = useTranslation();
  const v = getValue() as string;
  if (!v) return <>{"\u2014"}</>;
  try {
    return <>{new Date(v).toLocaleString(i18n.language)}</>;
  } catch {
    return <>{v}</>;
  }
}

/** Formats a date string for display, returning an em-dash when empty. */
export function fmtDate(v: string | undefined, locale?: string): string {
  if (!v) return "\u2014";
  try {
    return new Date(v).toLocaleString(locale || undefined);
  } catch {
    return v;
  }
}

/** Renders a string value or an em-dash if empty/undefined. */
export function FallbackCell({ getValue }: { getValue: () => unknown }) {
  const v = getValue() as string;
  return <>{v || "\u2014"}</>;
}

interface BadgeCellProps {
  getValue: () => unknown;
  positive?: string[];
  negative?: string[];
  labels?: Record<string, string>;
}

export function BadgeCell({
  getValue,
  positive = [],
  negative = [],
  labels = {},
}: BadgeCellProps) {
  const v = String(getValue());
  const label = labels[v] ?? v;
  const cls = positive.includes(v)
    ? "badge-green"
    : negative.includes(v)
      ? "badge-red"
      : "";
  return <span className={`badge ${cls}`.trim()}>{label}</span>;
}

/** Renders a boolean value as a green/grey badge with i18n labels. */
export function BooleanBadge({ value, trueLabel, falseLabel }: { value: boolean; trueLabel?: string; falseLabel?: string }) {
  const { t } = useTranslation();
  return value
    ? <span className="badge badge-green">{trueLabel ?? t("common.yes")}</span>
    : <span>{falseLabel ?? t("common.no")}</span>;
}

export function BooleanCell({ getValue }: { getValue: () => unknown }) {
  return <BooleanBadge value={!!getValue()} />;
}

// ─── Hook ───────────────────────────────────────────────────────

export function useServiceClient<S extends DescService>(service: S) {
  const client = useMemo(() => createClient(service, transport), [service]);
  const invalidate = useListInvalidator();
  return { client, invalidate, useListKey };
}
