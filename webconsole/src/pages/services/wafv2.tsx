/**
 * WAFv2 service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { WAFV2Service, type WebACLSummary, Scope } from "@/gen/wafv2_pb";
import { CreateWebACLRequestSchema } from "@/gen/wafv2_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, FallbackCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

type WebACLItem = WebACLSummary & { _scope: Scope };

const getColumns = (t: TFunction): ColumnDef<WebACLItem, any>[] => [
  { accessorKey: "name", header: t("services.wafv2.webaclNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.wafv2.idHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.wafv2.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "arn", header: t("services.wafv2.arnHeader"), cell: SmallMonoCell },
];

type DetailTab = "detail" | "json";

export function WAFv2Page() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(WAFV2Service);
  const { queryKey } = useListKey("wafv2");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<WebACLItem | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formScope, setFormScope] = useState<Scope>(Scope.REGIONAL);
  const [formDescription, setFormDescription] = useState("");
  const [formDefaultAction, setFormDefaultAction] = useState<"allow" | "block">("allow");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const [regional, cloudfront] = await Promise.all([client.listWebACLs({ scope: Scope.REGIONAL }), client.listWebACLs({ scope: Scope.CLOUDFRONT })]);
      return [...(regional.webacls ?? []).map((acl) => ({ ...acl, _scope: Scope.REGIONAL })), ...(cloudfront.webacls ?? []).map((acl) => ({ ...acl, _scope: Scope.CLOUDFRONT }))];
    },
    refetchInterval: REFETCH_INTERVAL,
  });
  const items: WebACLItem[] = dropEmpty(data ?? [], "id");

  const createMutation = useMutation({
    mutationFn: () => client.createWebACL(create(CreateWebACLRequestSchema, { name: formName, scope: formScope, description: formDescription, defaultaction: formDefaultAction === "allow" ? { allow: {} } : { block: {} }, visibilityconfig: { sampledrequestsenabled: false, cloudwatchmetricsenabled: false, metricname: formName } })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormScope(Scope.REGIONAL); setFormDescription(""); setFormDefaultAction("allow"); },
  });

  const deleteMutation = useMutation({
    mutationFn: (acl: WebACLItem) => client.deleteWebACL({ id: acl.id, name: acl.name, scope: acl._scope, locktoken: acl.locktoken }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (acls: WebACLItem[]) => Promise.allSettled(acls.map((acl) => client.deleteWebACL({ id: acl.id, name: acl.name, scope: acl._scope, locktoken: acl.locktoken }))),
    onSuccess: (_d, acls) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); const deletedIds = new Set(acls.map((a) => a.id)); setSelectedItem((p) => (p && deletedIds.has(p.id) ? null : p)); },
  });

  const handleRowClick = (row: WebACLItem) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.id);

  const selectedAcls = items.filter((i) => selectedIds.has(i.id));

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🛡️" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td className="detail-label">ID</td><td className="cell-mono">{selectedItem.id}</td></tr>
            <tr><td className="detail-label">Scope</td><td>{selectedItem._scope === Scope.CLOUDFRONT ? "CLOUDFRONT" : "REGIONAL"}</td></tr>
            <tr><td className="detail-label">Description</td><td>{selectedItem.description || "\u2014"}</td></tr>
            <tr><td className="detail-label">ARN</td><td className="cell-mono cell-long">{selectedItem.arn}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🛡️" title={t("services.wafv2.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.wafv2.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.wafv2.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.wafv2.title") }, { label: t("services.wafv2.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-wafv2">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<WebACLItem>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...columns]} data={items} getRowId={(row) => row.id} onRowClick={handleRowClick} selectedId={selectedItem?.id} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.wafv2.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.wafv2.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.wafv2.placeholder")} className="modal-input" /></label>
        <label>{t("services.wafv2.scopeLabel")}<select value={formScope} onChange={(e) => setFormScope(Number(e.target.value) as Scope)} className="modal-input"><option value={Scope.REGIONAL}>{t("services.wafv2.scopeRegional")}</option><option value={Scope.CLOUDFRONT}>{t("services.wafv2.scopeCloudfront")}</option></select></label>
        <label>{t("services.wafv2.descriptionLabel")}<input value={formDescription} onChange={(e) => setFormDescription(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label>{t("services.wafv2.defaultActionLabel")}<select value={formDefaultAction} onChange={(e) => setFormDefaultAction(e.target.value as "allow" | "block")} className="modal-input"><option value="allow">{t("services.wafv2.actionAllow")}</option><option value="block">{t("services.wafv2.actionBlock")}</option></select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.wafv2.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.wafv2.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(selectedAcls)} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
