/**
 * API Gateway service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { APIGatewayService, type RestApi, EndpointType } from "@/gen/apigateway_pb";
import { CreateRestApiRequestSchema, EndpointConfigurationSchema } from "@/gen/apigateway_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, DateCell, FallbackCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<RestApi, any>[] => [
  { accessorKey: "name", header: t("services.apigateway.nameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.apigateway.idHeader"), cell: MonoCell, size: 100 },
  { accessorKey: "description", header: t("services.apigateway.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "createddate", header: t("services.apigateway.createdHeader"), cell: DateCell },
];

type DetailTab = "detail" | "json";

export function APIGatewayPage() {
  const { t, i18n } = useTranslation();
  const { client } = useServiceClient(APIGatewayService);
  const { queryKey } = useListKey("apigateway");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<RestApi | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDesc, setFormDesc] = useState("");
  const [formEndpointType, setFormEndpointType] = useState<EndpointType>(EndpointType.REGIONAL);

  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<RestApi, Awaited<ReturnType<typeof client.getRestApis>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.getRestApis({ position: token || undefined, limit: 1000 }),
    getItems: (r) => r.items ?? [],
    getNextToken: (r) => r.position ?? "",
  });
  const items = dropEmpty(rawItems, "id");
  

  const createMutation = useMutation({
    mutationFn: () => client.createRestApi(create(CreateRestApiRequestSchema, { name: formName, description: formDesc, endpointconfiguration: create(EndpointConfigurationSchema, { types: [formEndpointType] }) })),
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormName(""); setFormDesc(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (id: string) => client.deleteRestApi({ restapiid: id }),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteRestApi({ restapiid: id }))),
    onSuccess: (_d, ids) => { invalidateList(); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.id) ? null : p)); },
  });

  const handleRowClick = (row: RestApi) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.id);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🌐" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td className="detail-label">ID</td><td className="cell-mono">{selectedItem.id}</td></tr>
            <tr><td className="detail-label">Description</td><td>{selectedItem.description || "\u2014"}</td></tr>
            {selectedItem.createddate && <tr><td className="detail-label">Created</td><td>{fmtDate(selectedItem.createddate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🌐" title={t("services.apigateway.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.apigateway.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.apigateway.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.apigateway.title") }, { label: t("services.apigateway.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-apigw">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<RestApi>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...columns]} data={items} getRowId={(row) => row.id} onRowClick={handleRowClick} selectedId={selectedItem?.id} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.apigateway.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.apigateway.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.apigateway.placeholder")} className="modal-input" /></label>
        <label>{t("services.apigateway.descriptionLabel")}<input value={formDesc} onChange={(e) => setFormDesc(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label>{t("services.apigateway.endpointTypeLabel")}<select value={formEndpointType} onChange={(e) => setFormEndpointType(Number(e.target.value))} className="modal-input"><option value={EndpointType.REGIONAL}>Regional</option><option value={EndpointType.PRIVATE}>Private</option><option value={EndpointType.EDGE}>Edge</option></select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.apigateway.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.apigateway.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
