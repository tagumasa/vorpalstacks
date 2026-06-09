/**
 * CloudWatch Logs service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudWatchLogsService, type LogGroupSummary } from "@/gen/cloudwatchlogs_pb";
import { CreateLogGroupRequestSchema, PutRetentionPolicyRequestSchema } from "@/gen/cloudwatchlogs_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<LogGroupSummary, any>[] => [
  { accessorKey: "loggroupname", header: t("services.cloudwatchlogs.logGroupHeader"), cell: MonoCell },
  { accessorKey: "loggroupclass", header: t("services.cloudwatchlogs.classHeader"), size: 100 },
  { accessorKey: "loggrouparn", header: t("services.cloudwatchlogs.arnHeader"), cell: SmallMonoCell },
];

type DetailTab = "detail" | "json";

export function CloudWatchLogsPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(CloudWatchLogsService);
  const { queryKey } = useListKey("cloudwatchlogs");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<LogGroupSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");
  const [formRetentionDays, setFormRetentionDays] = useState("");

  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<LogGroupSummary, Awaited<ReturnType<typeof client.listLogGroups>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listLogGroups({ nexttoken: token || undefined, limit: 1000 }),
    getItems: (r) => r.loggroups ?? [],
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "loggroupname");

  const createMutation = useMutation({
    mutationFn: async () => {
      await client.createLogGroup(create(CreateLogGroupRequestSchema, {
        loggroupname: formName,
        kmskeyid: formKmsKeyId || undefined,
      }));
      if (formRetentionDays) {
        await client.putRetentionPolicy(create(PutRetentionPolicyRequestSchema, { loggroupname: formName, retentionindays: Number(formRetentionDays) }));
      }
    },
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormName(""); setFormKmsKeyId(""); setFormRetentionDays(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteLogGroup({ loggroupname: name }),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((name) => client.deleteLogGroup({ loggroupname: name }))),
    onSuccess: (_d, names) => { invalidateList(); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.loggroupname) ? null : p)); },
  });

  const handleRowClick = (row: LogGroupSummary) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.loggroupname);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.loggroupname} titleIcon="📜" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.loggroupname}</td></tr>
            <tr><td className="detail-label">Class</td><td>{selectedItem.loggroupclass || "\u2014"}</td></tr>
            <tr><td className="detail-label">ARN</td><td className="cell-mono cell-long">{selectedItem.loggrouparn}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="📜" title={t("services.cloudwatchlogs.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cloudwatchlogs.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cloudwatchlogs.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cloudwatchlogs.title") }, { label: t("services.cloudwatchlogs.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cwlogs">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<LogGroupSummary>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.loggroupname), ...columns]} data={items} getRowId={(row) => row.loggroupname} onRowClick={handleRowClick} selectedId={selectedItem?.loggroupname} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cloudwatchlogs.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.cloudwatchlogs.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.cloudwatchlogs.placeholder")} className="modal-input" /></label>
        <label>{t("services.cloudwatchlogs.kmsKeyLabel")}<input value={formKmsKeyId} onChange={(e) => setFormKmsKeyId(e.target.value)} placeholder={t("services.cloudwatchlogs.kmsKeyPlaceholder")} className="modal-input" /></label>
        <label>{t("services.cloudwatchlogs.retentionLabel")}<select value={formRetentionDays} onChange={(e) => setFormRetentionDays(e.target.value)} className="modal-input">
          <option value="">{t("services.cloudwatchlogs.retentionNone")}</option>
          {[1,3,5,7,14,30,60,90,120,150,180,365,400,545,731,1827,3653].map((d) => <option key={d} value={String(d)}>{d}</option>)}
        </select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cloudwatchlogs.delete")} name={selectedItem?.loggroupname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.loggroupname)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.cloudwatchlogs.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
