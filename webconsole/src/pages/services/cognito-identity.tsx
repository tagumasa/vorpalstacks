/**
 * Cognito Identity service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CognitoIdentityService, type IdentityPoolShortDescription } from "@/gen/cognitoidentity_pb";
import { CreateIdentityPoolInputSchema } from "@/gen/cognitoidentity_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<IdentityPoolShortDescription, any>[] => [
  { accessorKey: "identitypoolid", header: t("services.cognitoIdentity.poolIdHeader"), cell: MonoCell },
  { accessorKey: "identitypoolname", header: t("services.cognitoIdentity.poolNameHeader"), cell: MonoCell },
];

type DetailTab = "detail" | "json";

export function CognitoIdentityPage() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(CognitoIdentityService);
  const { queryKey } = useListKey("cognitoIdentity");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<IdentityPoolShortDescription | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formTags, setFormTags] = useState("");

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listIdentityPools({}), refetchInterval: REFETCH_INTERVAL });
  const items: IdentityPoolShortDescription[] = dropEmpty(data?.identitypools ?? [], "identitypoolid");

  const createMutation = useMutation({
    mutationFn: () => {
      let identitypooltags: Record<string, string> = {};
      if (formTags.trim()) { try { identitypooltags = JSON.parse(formTags); } catch { throw new Error("Invalid tags JSON"); } }
      return client.createIdentityPool(create(CreateIdentityPoolInputSchema, { identitypoolname: formName, ...(Object.keys(identitypooltags).length > 0 ? { identitypooltags } : {}) }));
    },
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormTags(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (poolId: string) => client.deleteIdentityPool({ identitypoolid: poolId }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteIdentityPool({ identitypoolid: id }))),
    onSuccess: (_d, ids) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.identitypoolid) ? null : p)); },
  });

  const handleRowClick = (row: IdentityPoolShortDescription) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.identitypoolid);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.identitypoolname || selectedItem.identitypoolid} titleIcon="👤" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Pool ID</td><td className="cell-mono">{selectedItem.identitypoolid}</td></tr>
            <tr><td className="detail-label">Pool Name</td><td>{selectedItem.identitypoolname || "\u2014"}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="👤" title={t("services.cognitoIdentity.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cognitoIdentity.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cognitoIdentity.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cognitoIdentity.title") }, { label: t("services.cognitoIdentity.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cognito-id">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<IdentityPoolShortDescription>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.identitypoolid), ...columns]} data={items} getRowId={(row) => row.identitypoolid} onRowClick={handleRowClick} selectedId={selectedItem?.identitypoolid} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cognitoIdentity.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.cognitoIdentity.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.cognitoIdentity.placeholder")} className="modal-input" /></label>
        <label>{t("services.cognitoIdentity.tagsLabel")}<textarea value={formTags} onChange={(e) => setFormTags(e.target.value)} placeholder='{"key":"value"}' className="modal-input" rows={2} /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cognitoIdentity.delete")} name={selectedItem?.identitypoolname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.identitypoolid)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.cognitoIdentity.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
