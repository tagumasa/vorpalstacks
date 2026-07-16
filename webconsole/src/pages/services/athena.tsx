/**
 * Athena service page — 3-panel inspector layout.
 *
 * Panel 1 (toolbar): Breadcrumb navigation
 * Panel 2 (table):   Work group list with checkbox multi-select
 * Panel 3 (detail):  Work group detail (state, engine version, description)
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { AthenaService, type WorkGroupSummary, WorkGroupState } from "@/gen/athena_pb";
import { CreateWorkGroupInputSchema } from "@/gen/athena_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  DateCell,
  FallbackCell,
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

// ─── Helpers ────────────────────────────────────────────────────

const WORKGROUP_STATE_I18N: Record<number, string> = {
  [WorkGroupState.DISABLED]: "services.athena.stateDisabled",
  [WorkGroupState.ENABLED]: "services.athena.stateEnabled",
};

// ─── Column Definitions ─────────────────────────────────────────

const getColumns = (t: TFunction): ColumnDef<WorkGroupSummary, any>[] => [
  { accessorKey: "name", header: t("services.athena.workGroupHeader"), cell: MonoCell },
  { accessorKey: "state", header: t("services.athena.stateHeader"), cell: ({ getValue }) => { const v = getValue() as number; return <span className="badge">{WORKGROUP_STATE_I18N[v] ? t(WORKGROUP_STATE_I18N[v]) : String(v)}</span>; }, size: 90 },
  { accessorKey: "engineversion", header: t("services.athena.engineVersionHeader"), cell: ({ getValue }) => { const ev = getValue(); return ev ? String(ev) : "\u2014"; }, size: 100 },
  { accessorKey: "description", header: t("services.athena.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "creationtime", header: t("services.athena.createdHeader"), cell: DateCell },
];

// ─── Detail panel tab ───────────────────────────────────────────

type DetailTab = "detail" | "json";

// ─── Athena Page ────────────────────────────────────────────────

export function AthenaPage() {
  const { t, i18n } = useTranslation();
  const { client } = useServiceClient(AthenaService);
  const { queryKey } = useListKey("athena");
  const columns = getColumns(t);

  // ── Selection state ──────────────────────────────────────────
  const { selected: selectedNames, toggle: toggleName, toggleAll: toggleAllNames, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<WorkGroupSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  // ── Create form state ────────────────────────────────────────
  const [formName, setFormName] = useState("");
  const [formDescription, setFormDescription] = useState("");

  // ── Data ─────────────────────────────────────────────────────
  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<WorkGroupSummary, Awaited<ReturnType<typeof client.listWorkGroups>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listWorkGroups({ nexttoken: token || undefined }),
    getItems: (r) => r.workgroups ?? [],
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "name");

  // ── Mutations ────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: () =>
      client.createWorkGroup(
        create(CreateWorkGroupInputSchema, { name: formName, description: formDescription }),
      ),
    onSuccess: () => {
      invalidateList();
      setShowCreate(false);
      setFormName("");
      setFormDescription("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (workgroupName: string) =>
      client.deleteWorkGroup({ workgroup: workgroupName, recursivedeleteoption: true }),
    onSuccess: () => {
      invalidateList();
      setShowDelete(false);
      setSelectedItem(null);
      clearSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => {
      const results = await Promise.allSettled(
        names.map((name) => client.deleteWorkGroup({ workgroup: name, recursivedeleteoption: true }))
      );
      return results;
    },
    onSuccess: (_data, names) => {
      invalidateList();
      setShowBatchDelete(false);
      clearSelection();
      setSelectedItem((prev) => (prev && names.includes(prev.name) ? null : prev));
    },
  });

  // ── Handlers ─────────────────────────────────────────────────

  const handleRowClick = (row: WorkGroupSummary) => {
    setSelectedItem(row);
    setDetailTab("detail");
  };

  const allIds = items.map((i) => i.name);

  // ── Detail Panel ─────────────────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedItem) {
      return <DetailEmpty message={t("common.noItemSelected")} />;
    }

    const stateLabel = WORKGROUP_STATE_I18N[selectedItem.state] ? t(WORKGROUP_STATE_I18N[selectedItem.state]!) : String(selectedItem.state);
    const detailTabs = [
      { key: "detail", label: t("common.tabDetail") },
      { key: "json", label: t("common.rawJson") },
    ];

    return (
      <DetailPanel
        title={selectedItem.name}
        titleIcon="🔬"
        tabs={detailTabs}
        activeTab={detailTab}
        onTabChange={(k) => setDetailTab(k as DetailTab)}
        actions={
          <button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>
            {t("common.delete")}
          </button>
        }
      >
        {detailTab === "detail" ? (
          <table className="settings-table">
            <tbody>
              <tr><td className="detail-label-fixed">{t("common.fields.name")}</td><td className="cell-mono">{selectedItem.name}</td></tr>
              <tr><td className="detail-label">{t("common.fields.state")}</td><td><span className="badge">{stateLabel}</span></td></tr>
              <tr><td className="detail-label">{t("common.fields.engineVersion")}</td><td>{selectedItem.engineversion ? String(selectedItem.engineversion) : "\u2014"}</td></tr>
              <tr><td className="detail-label">{t("common.fields.description")}</td><td>{selectedItem.description || "\u2014"}</td></tr>
              {selectedItem.creationtime && (
                <tr><td className="detail-label">{t("common.fields.created")}</td><td>{fmtDate(selectedItem.creationtime, i18n.language)}</td></tr>
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
      icon="🔬"
      title={t("services.athena.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.athena.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.athena.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedNames.size === 0}
            onClick={() => setShowBatchDelete(true)}
          >
            {t("common.deleteSelected")}
            {selectedNames.size > 0 && <span className="batch-count">({selectedNames.size})</span>}
          </button>
        </>
      }
    >
      <div className="inspector-toolbar">
        <Breadcrumb parts={[
          { label: t("services.athena.title") },
          { label: t("services.athena.countLabel") },
        ]} />
        <div className="toolbar-selection-info">
          <SelectionBadge count={selectedNames.size} label={t("common.selectedCount", { count: selectedNames.size })} />
        </div>
      </div>

      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-athena">
          <div className="flex-fill-scroll">
            <DataTable
              columns={[
                checkboxColumn<WorkGroupSummary>(selectedNames, toggleName, () => toggleAllNames(allIds), allIds, t, (row) => row.name),
                ...columns,
              ]}
              data={items}
              getRowId={(row) => row.name}
              onRowClick={handleRowClick}
              selectedId={selectedItem?.name}
              hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore}
            />
          </div>
          {renderDetailPanel()}
        </Splitter>
      ) : (
        <div className="empty-state">{t("common.noData")}</div>
      )}

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.athena.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.athena.nameField")}
          <input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.athena.placeholder")} className="modal-input" />
        </label>
        <label>
          {t("services.athena.descLabel")}
          <input value={formDescription} onChange={(e) => setFormDescription(e.target.value)} placeholder={t("services.athena.descPlaceholder")} className="modal-input" />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.athena.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />

      <ServiceDeleteDialog
        open={showBatchDelete}
        title={t("common.deleteSelected")}
        name={`${selectedNames.size} ${t("services.athena.countLabel")}`}
        error={batchDeleteMutation.error}
        isPending={batchDeleteMutation.isPending}
        onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedNames))}
        onClose={() => setShowBatchDelete(false)}
      />
    </ServicePageLayout>
  );
}
