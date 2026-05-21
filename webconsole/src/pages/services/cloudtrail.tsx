/**
 * CloudTrail service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudTrailService, type TrailInfo } from "@/gen/cloudtrail_pb";
import { CreateTrailRequestSchema } from "@/gen/cloudtrail_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<TrailInfo, any>[] => [
  { accessorKey: "name", header: t("services.cloudtrail.trailNameHeader"), cell: MonoCell },
  { accessorKey: "trailarn", header: t("services.cloudtrail.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "homeregion", header: t("services.cloudtrail.homeRegionHeader"), size: 120 },
];

type DetailTab = "detail" | "json";

export function CloudTrailPage() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(CloudTrailService);
  const { queryKey } = useListKey("cloudtrail");
  const columns = getColumns(t);

  const { selected: selectedNames, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<TrailInfo | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formS3Bucket, setFormS3Bucket] = useState("");
  const [formMultiRegion, setFormMultiRegion] = useState(true);
  const [formGlobalEvents, setFormGlobalEvents] = useState(true);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listTrails({}), refetchInterval: REFETCH_INTERVAL });
  const items: TrailInfo[] = dropEmpty(data?.trails ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () => client.createTrail(create(CreateTrailRequestSchema, { name: formName, s3bucketname: formS3Bucket, ismultiregiontrail: formMultiRegion, includeglobalserviceevents: formGlobalEvents })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormS3Bucket(""); setFormMultiRegion(true); setFormGlobalEvents(true); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteTrail({ name }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteTrail({ name: n }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.name) ? null : p)); },
  });

  const handleRowClick = (row: TrailInfo) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.name);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="📋" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td className="detail-label">ARN</td><td className="cell-mono cell-long">{selectedItem.trailarn}</td></tr>
            <tr><td className="detail-label">Home Region</td><td>{selectedItem.homeregion}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="📋" title={t("services.cloudtrail.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cloudtrail.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cloudtrail.create")}</button>
      <button className="btn btn-danger" disabled={selectedNames.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedNames.size > 0 && <span className="batch-count">({selectedNames.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cloudtrail.title") }, { label: t("services.cloudtrail.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedNames.size} label={t("common.selectedCount", { count: selectedNames.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cloudtrail">
          <div className="flex-fill-scroll"><DataTable exportName="cloudtrail" columns={[checkboxColumn<TrailInfo>(selectedNames, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.name), ...columns]} data={items} getRowId={(row) => row.name} onRowClick={handleRowClick} selectedId={selectedItem?.name} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cloudtrail.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName || !formS3Bucket}>
        <label>{t("services.cloudtrail.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.cloudtrail.placeholder")} className="modal-input" /></label>
        <label>{t("services.cloudtrail.s3BucketLabel")}<input value={formS3Bucket} onChange={(e) => setFormS3Bucket(e.target.value)} placeholder={t("services.cloudtrail.s3BucketPlaceholder")} className="modal-input" /></label>
        <label className="checkbox-label"><input type="checkbox" checked={formMultiRegion} onChange={(e) => setFormMultiRegion(e.target.checked)} />{t("services.cloudtrail.multiRegionLabel")}</label>
        <label className="checkbox-label"><input type="checkbox" checked={formGlobalEvents} onChange={(e) => setFormGlobalEvents(e.target.checked)} />{t("services.cloudtrail.globalServiceEventsLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cloudtrail.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedNames.size} ${t("services.cloudtrail.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedNames))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
