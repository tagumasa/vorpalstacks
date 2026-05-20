/**
 * Kinesis service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { KinesisService, type StreamSummary, StreamMode } from "@/gen/kinesis_pb";
import { CreateStreamInputSchema, StreamModeDetailsSchema } from "@/gen/kinesis_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<StreamSummary, any>[] => [
  { accessorKey: "streamname", header: t("services.kinesis.streamNameHeader"), cell: MonoCell },
  { accessorKey: "streamstatus", header: t("services.kinesis.streamStatusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 90 },
  { accessorKey: "streamarn", header: t("services.kinesis.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "streamcreationtimestamp", header: t("services.kinesis.streamCreatedHeader"), cell: DateCell },
];

type DetailTab = "detail" | "json";

export function KinesisPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(KinesisService);
  const { queryKey } = useListKey("kinesis");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<StreamSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formShardCount, setFormShardCount] = useState("1");
  const [formStreamMode, setFormStreamMode] = useState(StreamMode.PROVISIONED);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listStreams({}), refetchInterval: REFETCH_INTERVAL });
  const items: StreamSummary[] = dropEmpty(data?.streamsummaries ?? [], "streamname");

  const createMutation = useMutation({
    mutationFn: () => client.createStream(create(CreateStreamInputSchema, { streamname: formName, shardcount: parseInt(formShardCount, 10) || 1, streammodedetails: create(StreamModeDetailsSchema, { streammode: formStreamMode }) })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormShardCount("1"); setFormStreamMode(StreamMode.PROVISIONED); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteStream({ streamname: name }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((name) => client.deleteStream({ streamname: name }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.streamname) ? null : p)); },
  });

  const handleRowClick = (row: StreamSummary) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.streamname);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.streamname} titleIcon="🌊" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Name</td><td className="cell-mono">{selectedItem.streamname}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Status</td><td><span className="badge">{String(selectedItem.streamstatus)}</span></td></tr>
            <tr><td style={{ fontWeight: 600 }}>ARN</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedItem.streamarn}</td></tr>
            {selectedItem.streamcreationtimestamp && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{fmtDate(selectedItem.streamcreationtimestamp, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🌊" title={t("services.kinesis.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.kinesis.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.kinesis.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.kinesis.title") }, { label: t("services.kinesis.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-kinesis">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<StreamSummary>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.streamname), ...columns]} data={items} getRowId={(row) => row.streamname} onRowClick={handleRowClick} selectedId={selectedItem?.streamname} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.kinesis.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.kinesis.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.kinesis.placeholder")} className="modal-input" /></label>
        <label>{t("services.kinesis.shardCountLabel")}<input type="number" min="1" value={formShardCount} onChange={(e) => setFormShardCount(e.target.value)} className="modal-input" /></label>
        <label>{t("services.kinesis.modeLabel")}<select value={formStreamMode} onChange={(e) => setFormStreamMode(Number(e.target.value))} className="modal-input"><option value={StreamMode.PROVISIONED}>{t("services.kinesis.modeProvisioned")}</option><option value={StreamMode.ON_DEMAND}>{t("services.kinesis.modeOnDemand")}</option></select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.kinesis.delete")} name={selectedItem?.streamname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.streamname)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.kinesis.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
