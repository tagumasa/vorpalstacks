/**
 * SSM Parameter Store service page — 3-panel inspector layout.
 *
 * Panel 1 (toolbar): Breadcrumb navigation
 * Panel 2 (table):   Parameter list with checkbox multi-select
 * Panel 3 (detail):  Inline parameter detail with View=Edit (Structured + JSON tabs)
 */
import { useState, useCallback } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  SSMService,
  PutParameterRequestSchema,
  DeleteParameterRequestSchema,
  type ParameterMetadata,
  ParameterType,
  ParameterTier,
} from "@/gen/ssm_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  DateCell,
  fmtDate,
  useServiceClient,
} from "@/components/shared/service-page";
import {
  checkboxColumn,
  Breadcrumb,
  SelectionBadge,
  DetailPanel,
  DetailEmpty,
  useSelection,
} from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

// ─── Column Definitions ─────────────────────────────────────────

const PARAM_TYPES = [
  { value: ParameterType.STRING, i18nKey: "services.ssm.paramTypeString" },
  { value: ParameterType.SECURE_STRING, i18nKey: "services.ssm.paramTypeSecureString" },
  { value: ParameterType.STRING_LIST, i18nKey: "services.ssm.paramTypeStringList" },
];

const getColumns = (t: TFunction): ColumnDef<ParameterMetadata, any>[] => [
  { accessorKey: "name", header: t("services.ssm.nameHeader"), cell: MonoCell },
  {
    accessorKey: "type",
    header: t("services.ssm.typeHeader"),
    cell: ({ getValue }) => {
      const val = getValue() as ParameterType;
      const labels: Record<number, string> = {
        [ParameterType.SECURE_STRING]: t("services.ssm.paramTypeSecureString"),
        [ParameterType.STRING_LIST]: t("services.ssm.paramTypeStringList"),
        [ParameterType.STRING]: t("services.ssm.paramTypeString"),
      };
      return <span className="badge">{labels[val] ?? String(val)}</span>;
    },
  },
  { accessorKey: "version", header: t("services.ssm.versionHeader"), size: 50 },
  {
    accessorKey: "tier",
    header: t("services.ssm.tierHeader"),
    cell: ({ getValue }) => {
      const val = getValue() as ParameterTier;
      const labels: Record<number, string> = {
        [ParameterTier.STANDARD]: t("services.ssm.tierStandard"),
        [ParameterTier.ADVANCED]: t("services.ssm.tierAdvanced"),
        [ParameterTier.INTELLIGENT_TIERING]: t("services.ssm.tierIntelligent"),
      };
      return <span>{labels[val] ?? String(val)}</span>;
    },
    size: 90,
  },
  {
    accessorKey: "description",
    header: t("services.ssm.descriptionHeader"),
    cell: ({ getValue }) => (getValue() as string) || "\u2014",
  },
  { accessorKey: "lastmodifieddate", header: t("services.ssm.lastModifiedHeader"), cell: DateCell },
];

// ─── Detail panel tab ───────────────────────────────────────────

type DetailTab = "detail" | "json";

// ─── SSM Page ───────────────────────────────────────────────────

export function SSMPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(SSMService);
  const { queryKey } = useListKey("ssm");
  const columns = getColumns(t);

  // ── Selection state ──────────────────────────────────────────
  const { selected: selectedNames, toggle: toggleName, toggleAll: toggleAllNames, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<ParameterMetadata | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  // ── Create form state ────────────────────────────────────────
  const [formName, setFormName] = useState("");
  const [formValue, setFormValue] = useState("");
  const [formType, setFormType] = useState(ParameterType.STRING);
  const [formDesc, setFormDesc] = useState("");

  // ── Data ─────────────────────────────────────────────────────
  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.describeParameters({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: ParameterMetadata[] = dropEmpty(data?.parameters ?? [], "name");

  // ── Mutations ────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: () =>
      client.putParameter(
        create(PutParameterRequestSchema, {
          name: formName,
          value: formValue,
          type: formType,
          description: formDesc,
          tier: ParameterTier.STANDARD,
          overwrite: false,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormValue("");
      setFormType(ParameterType.STRING);
      setFormDesc("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) =>
      client.deleteParameter(create(DeleteParameterRequestSchema, { name })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
      clearSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => {
      const results = await Promise.allSettled(
        names.map((name) =>
          client.deleteParameter(create(DeleteParameterRequestSchema, { name }))
        ),
      );
      return results;
    },
    onSuccess: (_data, names) => {
      invalidate(queryKey);
      setShowBatchDelete(false);
      clearSelection();
      setSelectedItem((prev) => (prev && names.includes(prev.name) ? null : prev));
    },
  });

  // ── Handlers ─────────────────────────────────────────────────

  const handleRowClick = useCallback((row: ParameterMetadata) => {
    setSelectedItem(row);
    setDetailTab("detail");
  }, []);

  const allIds = items.map((i) => i.name);

  // ── Breadcrumb ───────────────────────────────────────────────

  const breadcrumb = (
    <Breadcrumb parts={[
      { label: t("services.ssm.title") },
      { label: t("services.ssm.countLabel") },
    ]} />
  );

  const selectionInfo = (
    <SelectionBadge count={selectedNames.size} label={t("common.selectedCount", { count: selectedNames.size })} />
  );

  // ── Detail Panel ─────────────────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedItem) {
      return <DetailEmpty message={t("common.noItemSelected")} />;
    }

    const detailTabs = [
      { key: "detail", label: t("services.ssm.title").split(" ")[0] ?? "Detail" },
      { key: "json", label: t("common.rawJson") },
    ];

    return (
      <DetailPanel
        title={selectedItem.name}
        titleIcon="📋"
        tabs={detailTabs}
        activeTab={detailTab}
        onTabChange={(k) => setDetailTab(k as DetailTab)}
        actions={
          <button
            className="btn btn-danger btn-sm"
            onClick={() => setShowDelete(true)}
          >
            {t("common.delete")}
          </button>
        }
      >
        {detailTab === "detail" ? (
          <table className="settings-table">
            <tbody>
              <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
              <tr><td className="detail-label">Type</td><td>{(() => {
                const labels: Record<number, string> = {
                  [ParameterType.SECURE_STRING]: t("services.ssm.paramTypeSecureString"),
                  [ParameterType.STRING_LIST]: t("services.ssm.paramTypeStringList"),
                  [ParameterType.STRING]: t("services.ssm.paramTypeString"),
                };
                return labels[selectedItem.type] ?? String(selectedItem.type);
              })()}</td></tr>
              <tr><td className="detail-label">Version</td><td>{selectedItem.version}</td></tr>
              <tr><td className="detail-label">Tier</td><td>{(() => {
                const labels: Record<number, string> = {
                  [ParameterTier.STANDARD]: t("services.ssm.tierStandard"),
                  [ParameterTier.ADVANCED]: t("services.ssm.tierAdvanced"),
                  [ParameterTier.INTELLIGENT_TIERING]: t("services.ssm.tierIntelligent"),
                };
                return labels[selectedItem.tier] ?? String(selectedItem.tier);
              })()}</td></tr>
              {selectedItem.description && (
                <tr><td className="detail-label">Description</td><td>{selectedItem.description}</td></tr>
              )}
              <tr><td className="detail-label">Data Type</td><td>{selectedItem.datatype || "text"}</td></tr>
              {selectedItem.lastmodifieddate && (
                <tr><td className="detail-label">Last Modified</td><td>{fmtDate(selectedItem.lastmodifieddate, i18n.language)}</td></tr>
              )}
            </tbody>
          </table>
        ) : (
          <JsonViewer data={selectedItem} />
        )}
      </DetailPanel>
    );
  };

  // ── Render ───────────────────────────────────────────────────

  return (
    <ServicePageLayout
      icon="📋"
      title={t("services.ssm.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.ssm.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.ssm.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedNames.size === 0}
            onClick={() => setShowBatchDelete(true)}
          >
            {t("common.deleteSelected")}
            {selectedNames.size > 0 && (
              <span className="batch-count">({selectedNames.size})</span>
            )}
          </button>
        </>
      }
    >
      {/* Inspector toolbar */}
      <div className="inspector-toolbar">
        {breadcrumb}
        <div className="toolbar-selection-info">{selectionInfo}</div>
      </div>

      {/* Table + detail split */}
      {items.length > 0 ? (
        <Splitter
          direction="horizontal"
          initialSize={240}
          minSize={80}
          maxSize={600}
          storageKey="vs-split-ssm"
        >
          <div className="flex-fill-scroll">
            <DataTable
              exportName="ssm"
              columns={[
                checkboxColumn<ParameterMetadata>(selectedNames, toggleName, () => toggleAllNames(allIds), allIds, t, (row) => row.name),
                ...columns,
              ]}
              data={items}
              getRowId={(row) => row.name}
              onRowClick={handleRowClick}
              selectedId={selectedItem?.name}
            />
          </div>
          {renderDetailPanel()}
        </Splitter>
      ) : (
        <div className="empty-state">{t("common.noData")}</div>
      )}

      {/* ── Modals ──────────────────────────────────────────── */}

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.ssm.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formValue}
      >
        <label>
          {t("services.ssm.nameField")}
          <input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.ssm.placeholder")} className="modal-input" />
        </label>
        <label>
          {t("services.ssm.valueLabel")}
          <textarea value={formValue} onChange={(e) => setFormValue(e.target.value)} placeholder={t("services.ssm.valuePlaceholder")} className="modal-textarea" rows={3} />
        </label>
        <label>
          {t("services.ssm.typeLabel")}
          <select value={formType} onChange={(e) => setFormType(Number(e.target.value))} className="modal-select">
            {PARAM_TYPES.map((pt) => (
              <option key={pt.value} value={pt.value}>{t(pt.i18nKey)}</option>
            ))}
          </select>
        </label>
        <label>
          {t("services.ssm.descLabel")}
          <input value={formDesc} onChange={(e) => setFormDesc(e.target.value)} placeholder={t("services.ssm.descPlaceholder")} className="modal-input" />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.ssm.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />

      {/* Batch Delete */}
      <ServiceDeleteDialog
        open={showBatchDelete}
        title={t("common.deleteSelected")}
        name={`${selectedNames.size} ${t("services.ssm.countLabel")}`}
        error={batchDeleteMutation.error}
        isPending={batchDeleteMutation.isPending}
        onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedNames))}
        onClose={() => setShowBatchDelete(false)}
      />
    </ServicePageLayout>
  );
}
