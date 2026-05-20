/**
 * Cognito IDP service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CognitoIdentityProviderService, VerifiedAttributeType } from "@/gen/cognitoidentityprovider_pb";
import { CreateUserPoolRequestSchema, PasswordPolicyTypeSchema, UserPoolPolicyTypeSchema } from "@/gen/cognitoidentityprovider_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, BadgeCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

interface TableRow { name: string; id: string; status: string; creationdate: string; lastmodifieddate: string; }

const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.cognito.nameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.cognito.idHeader"), cell: SmallMonoCell },
  { accessorKey: "status", header: t("services.cognito.statusHeader"), cell: BadgeCell, size: 90 },
  { accessorKey: "creationdate", header: t("services.cognito.createdHeader"), cell: DateCell },
];

type DetailTab = "detail" | "json";

export function CognitoPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(CognitoIdentityProviderService);
  const { queryKey } = useListKey("cognito");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formPwMin, setFormPwMin] = useState("8");
  const [formPwUpper, setFormPwUpper] = useState(true);
  const [formPwLower, setFormPwLower] = useState(true);
  const [formPwNumbers, setFormPwNumbers] = useState(true);
  const [formPwSymbols, setFormPwSymbols] = useState(false);
  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listUserPools({}), refetchInterval: REFETCH_INTERVAL });
  const items: TableRow[] = dropEmpty((data?.userpools ?? []).map((p) => ({ name: p.name, id: p.id, status: String(p.status), creationdate: p.creationdate, lastmodifieddate: p.lastmodifieddate })), "id");

  const createMutation = useMutation({
    mutationFn: () => {
      return client.createUserPool(create(CreateUserPoolRequestSchema, {
        poolname: formName,
        policies: create(UserPoolPolicyTypeSchema, { passwordpolicy: create(PasswordPolicyTypeSchema, { minimumlength: Number(formPwMin), requireuppercase: formPwUpper, requirelowercase: formPwLower, requirenumbers: formPwNumbers, requiresymbols: formPwSymbols }) }),
        autoverifiedattributes: [VerifiedAttributeType.EMAIL],
      }));
    },
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (userpoolid: string) => client.deleteUserPool({ userpoolid }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteUserPool({ userpoolid: id }))),
    onSuccess: (_d, ids) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.id) ? null : p)); },
  });

  const handleRowClick = (row: TableRow) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.id);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🔐" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Pool ID</td><td className="cell-mono">{selectedItem.id}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Status</td><td><span className="badge">{selectedItem.status}</span></td></tr>
            {selectedItem.creationdate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{fmtDate(selectedItem.creationdate, i18n.language)}</td></tr>}
            {selectedItem.lastmodifieddate && <tr><td style={{ fontWeight: 600 }}>Last Modified</td><td>{fmtDate(selectedItem.lastmodifieddate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🔐" title={t("services.cognito.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cognito.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cognito.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cognito.title") }, { label: t("services.cognito.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cognito">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<TableRow>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...columns]} data={items} getRowId={(row) => row.id} onRowClick={handleRowClick} selectedId={selectedItem?.id} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cognito.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.cognito.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.cognito.placeholder")} className="modal-input" /></label>
        <label>{t("services.cognito.pwMinLengthLabel")}<input type="number" min={6} max={99} value={formPwMin} onChange={(e) => setFormPwMin(e.target.value)} className="modal-input" /></label>
        <label className="checkbox-label"><input type="checkbox" checked={formPwUpper} onChange={(e) => setFormPwUpper(e.target.checked)} />{t("services.cognito.pwRequireUpperLabel")}</label>
        <label className="checkbox-label"><input type="checkbox" checked={formPwLower} onChange={(e) => setFormPwLower(e.target.checked)} />{t("services.cognito.pwRequireLowerLabel")}</label>
        <label className="checkbox-label"><input type="checkbox" checked={formPwNumbers} onChange={(e) => setFormPwNumbers(e.target.checked)} />{t("services.cognito.pwRequireNumbersLabel")}</label>
        <label className="checkbox-label"><input type="checkbox" checked={formPwSymbols} onChange={(e) => setFormPwSymbols(e.target.checked)} />{t("services.cognito.pwRequireSymbolsLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cognito.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.cognito.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
