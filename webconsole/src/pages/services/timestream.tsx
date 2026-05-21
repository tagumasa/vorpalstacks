/**
 * Timestream service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { TimestreamWriteService } from "@/gen/timestreamwrite_pb";
import { CreateDatabaseRequestSchema } from "@/gen/timestreamwrite_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

interface TableRow {
  databasename: string;
  arn: string;
  tablecount: number;
  creationtime: string;
  lastupdatedtime: string;
}

const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "databasename", header: t("services.timestream.databaseNameHeader"), cell: MonoCell },
  { accessorKey: "tablecount", header: t("services.timestream.tableCountHeader"), size: 80 },
  { accessorKey: "creationtime", header: t("services.timestream.createdHeader"), cell: DateCell },
  { accessorKey: "arn", header: t("services.timestream.arnHeader"), cell: SmallMonoCell },
];

type DetailTab = "detail" | "json";

export function TimestreamPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(TimestreamWriteService);
  const { queryKey } = useListKey("timestream");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");
  const [formMagneticRetentionDays, setFormMagneticRetentionDays] = useState(365);
  const [formMemoryRetentionHours, setFormMemoryRetentionHours] = useState(24);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listDatabases({}), refetchInterval: REFETCH_INTERVAL });
  const items: TableRow[] = dropEmpty((data?.databases ?? []).map((db) => ({ databasename: db.databasename, arn: db.arn, tablecount: Number(db.tablecount ?? 0), creationtime: db.creationtime ?? "", lastupdatedtime: db.lastupdatedtime ?? "" })), "databasename");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createDatabase(create(CreateDatabaseRequestSchema, {
        databasename: formName,
        kmskeyid: formKmsKeyId || undefined,
      })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormKmsKeyId(""); setFormMagneticRetentionDays(365); setFormMemoryRetentionHours(24); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteDatabase({ databasename: name }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteDatabase({ databasename: n }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.databasename) ? null : p)); },
  });

  const handleRowClick = (row: TableRow) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.databasename);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.databasename} titleIcon="⏱" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Database</td><td className="cell-mono">{selectedItem.databasename}</td></tr>
            <tr><td className="detail-label">Tables</td><td>{selectedItem.tablecount}</td></tr>
            {selectedItem.creationtime && <tr><td className="detail-label">Created</td><td>{fmtDate(selectedItem.creationtime, i18n.language)}</td></tr>}
            {selectedItem.lastupdatedtime && <tr><td className="detail-label">Updated</td><td>{fmtDate(selectedItem.lastupdatedtime, i18n.language)}</td></tr>}
            <tr><td className="detail-label">ARN</td><td className="cell-mono cell-long">{selectedItem.arn}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="⏱" title={t("services.timestream.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.timestream.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.timestream.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.timestream.title") }, { label: t("services.timestream.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-ts">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<TableRow>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.databasename), ...columns]} data={items} getRowId={(row) => row.databasename} onRowClick={handleRowClick} selectedId={selectedItem?.databasename} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.timestream.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.timestream.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.timestream.placeholder")} className="modal-input" /></label>
        <label>{t("services.timestream.kmsKeyLabel")}<input value={formKmsKeyId} onChange={(e) => setFormKmsKeyId(e.target.value)} placeholder={t("services.timestream.kmsKeyPlaceholder")} className="modal-input" /></label>
        <label>{t("services.timestream.magneticRetentionLabel")}<input type="number" min={1} max={73000} value={formMagneticRetentionDays} onChange={(e) => setFormMagneticRetentionDays(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.timestream.memoryRetentionLabel")}<input type="number" min={1} max={8766} value={formMemoryRetentionHours} onChange={(e) => setFormMemoryRetentionHours(Number(e.target.value))} className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.timestream.delete")} name={selectedItem?.databasename} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.databasename)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.timestream.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
