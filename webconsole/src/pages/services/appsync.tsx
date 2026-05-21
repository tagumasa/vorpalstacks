/**
 * AppSync service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { AppSyncService, type GraphqlApi, AuthenticationType } from "@/gen/appsync_pb";
import { CreateGraphqlApiRequestSchema, DeleteGraphqlApiRequestSchema } from "@/gen/appsync_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const AUTH_TYPE_LABELS: Record<number, string> = {
  [AuthenticationType.API_KEY]: "API_KEY",
  [AuthenticationType.AWS_IAM]: "AWS_IAM",
  [AuthenticationType.OPENID_CONNECT]: "OPENID_CONNECT",
  [AuthenticationType.AMAZON_COGNITO_USER_POOLS]: "COGNITO_USER_POOLS",
};

const getColumns = (t: TFunction): ColumnDef<GraphqlApi, any>[] => [
  { accessorKey: "name", header: t("services.appsync.apiNameHeader"), cell: MonoCell },
  { accessorKey: "apiid", header: t("services.appsync.apiIdHeader"), cell: SmallMonoCell, size: 120 },
  { accessorKey: "authenticationtype", header: t("services.appsync.authTypeHeader"), cell: ({ getValue }) => { const v = getValue() as number; return <span className="badge">{AUTH_TYPE_LABELS[v] ?? String(v)}</span>; }, size: 100 },
];

type DetailTab = "detail" | "json";

export function AppSyncPage() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(AppSyncService);
  const { queryKey } = useListKey("appsync");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<GraphqlApi | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formAuthType, setFormAuthType] = useState<AuthenticationType>(AuthenticationType.API_KEY);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listGraphqlApis({}), refetchInterval: REFETCH_INTERVAL });
  const items: GraphqlApi[] = dropEmpty(data?.graphqlapis ?? [], "apiid");

  const createMutation = useMutation({
    mutationFn: () => client.createGraphqlApi(create(CreateGraphqlApiRequestSchema, { name: formName, authenticationtype: formAuthType })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (apiId: string) => client.deleteGraphqlApi(create(DeleteGraphqlApiRequestSchema, { apiid: apiId })),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteGraphqlApi(create(DeleteGraphqlApiRequestSchema, { apiid: id })))),
    onSuccess: (_d, ids) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.apiid) ? null : p)); },
  });

  const handleRowClick = (row: GraphqlApi) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.apiid);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const authLabel = AUTH_TYPE_LABELS[selectedItem.authenticationtype] ?? String(selectedItem.authenticationtype);
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🔮" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td className="detail-label">API ID</td><td className="cell-mono">{selectedItem.apiid}</td></tr>
            <tr><td className="detail-label">Auth Type</td><td><span className="badge">{authLabel}</span></td></tr>
            {selectedItem.uris && Object.keys(selectedItem.uris).length > 0 && <tr><td className="detail-label">URIs</td><td className="cell-mono cell-long">{Object.entries(selectedItem.uris).map(([k, v]) => `${k}: ${v}`).join(", ")}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🔮" title={t("services.appsync.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.appsync.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.appsync.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.appsync.title") }, { label: t("services.appsync.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-appsync">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<GraphqlApi>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.apiid), ...columns]} data={items} getRowId={(row) => row.apiid} onRowClick={handleRowClick} selectedId={selectedItem?.apiid} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.appsync.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.appsync.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.appsync.placeholder")} className="modal-input" /></label>
        <label>{t("services.appsync.authTypeLabel")}<select value={formAuthType} onChange={(e) => setFormAuthType(Number(e.target.value))} className="modal-select">
          <option value={AuthenticationType.API_KEY}>API Key</option>
          <option value={AuthenticationType.AWS_IAM}>AWS IAM</option>
          <option value={AuthenticationType.OPENID_CONNECT}>OpenID Connect</option>
          <option value={AuthenticationType.AMAZON_COGNITO_USER_POOLS}>Cognito User Pools</option>
        </select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.appsync.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.apiid)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.appsync.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
