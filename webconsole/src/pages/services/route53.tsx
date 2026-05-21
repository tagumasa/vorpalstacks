/**
 * Route 53 service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { Route53Service, type HostedZone } from "@/gen/route53_pb";
import { CreateHostedZoneRequestSchema, HostedZoneConfigSchema } from "@/gen/route53_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<HostedZone, any>[] => [
  { accessorKey: "name", header: t("services.route53.zoneNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.route53.zoneIdHeader"), cell: MonoCell, size: 120 },
  { accessorKey: "resourcerecordsetcount", header: t("services.route53.recordCountHeader"), size: 80 },
  { accessorKey: "config.comment", header: t("services.route53.commentHeader") },
];

type DetailTab = "detail" | "json";

export function Route53Page() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(Route53Service);
  const { queryKey } = useListKey("route53");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<HostedZone | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formComment, setFormComment] = useState("");
  const [formPrivate, setFormPrivate] = useState(false);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listHostedZones({}), refetchInterval: REFETCH_INTERVAL });
  const items: HostedZone[] = dropEmpty(data?.hostedzones ?? [], "id");

  const createMutation = useMutation({
    mutationFn: () => client.createHostedZone(create(CreateHostedZoneRequestSchema, { callerreference: `vs-${Date.now()}`, name: formName, hostedzoneconfig: create(HostedZoneConfigSchema, { comment: formComment, privatezone: formPrivate }) })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormComment(""); setFormPrivate(false); },
  });

  const deleteMutation = useMutation({
    mutationFn: (zoneId: string) => client.deleteHostedZone({ id: zoneId }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteHostedZone({ id }))),
    onSuccess: (_d, ids) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.id) ? null : p)); },
  });

  const handleRowClick = (row: HostedZone) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.id);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🌐" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td className="detail-label">Zone ID</td><td className="cell-mono">{selectedItem.id}</td></tr>
            <tr><td className="detail-label">Records</td><td>{String(selectedItem.resourcerecordsetcount)}</td></tr>
            {selectedItem.config?.comment && <tr><td className="detail-label">Comment</td><td>{selectedItem.config.comment}</td></tr>}
            {selectedItem.config && <tr><td className="detail-label">Private</td><td>{selectedItem.config.privatezone ? "Yes" : "No"}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🌐" title={t("services.route53.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.route53.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.route53.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.route53.title") }, { label: t("services.route53.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-route53">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<HostedZone>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...columns]} data={items} getRowId={(row) => row.id} onRowClick={handleRowClick} selectedId={selectedItem?.id} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.route53.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.route53.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.route53.placeholder")} className="modal-input" /></label>
        <label>{t("services.route53.commentLabel")}<input value={formComment} onChange={(e) => setFormComment(e.target.value)} placeholder={t("services.route53.commentPlaceholder")} className="modal-input" /></label>
        <label className="checkbox-label"><input type="checkbox" checked={formPrivate} onChange={(e) => setFormPrivate(e.target.checked)} />{t("services.route53.privateZoneLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.route53.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.route53.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
