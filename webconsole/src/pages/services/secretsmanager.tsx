/**
 * Secrets Manager service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { SecretsManagerService, type SecretListEntry } from "@/gen/secretsmanager_pb";
import { CreateSecretRequestSchema } from "@/gen/secretsmanager_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, FallbackCell, BooleanCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getColumns = (t: TFunction): ColumnDef<SecretListEntry, any>[] => [
  { accessorKey: "name", header: t("services.secretsmanager.secretNameHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.secretsmanager.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "rotationenabled", header: t("services.secretsmanager.rotationHeader"), cell: BooleanCell, size: 80 },
  { accessorKey: "createddate", header: t("services.secretsmanager.createdHeader"), cell: DateCell },
  { accessorKey: "deleteddate", header: t("services.secretsmanager.deletedHeader"), cell: DateCell },
  { accessorKey: "kmskeyid", header: t("services.secretsmanager.kmsKeyHeader"), cell: SmallMonoCell },
];

type DetailTab = "detail" | "json";

export function SecretsManagerPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(SecretsManagerService);
  const { queryKey } = useListKey("secretsmanager");
  const columns = getColumns(t);

  const { selected: selectedNames, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<SecretListEntry | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDesc, setFormDesc] = useState("");
  const [formSecretValue, setFormSecretValue] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listSecrets({}), refetchInterval: REFETCH_INTERVAL });
  const items: SecretListEntry[] = dropEmpty(data?.secretlist ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createSecret(create(CreateSecretRequestSchema, {
        name: formName,
        description: formDesc,
        secretstring: formSecretValue,
        kmskeyid: formKmsKeyId || undefined,
      })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormDesc(""); setFormSecretValue(""); setFormKmsKeyId(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (secretName: string) => client.deleteSecret({ secretid: secretName, forcedeletewithoutrecovery: true }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteSecret({ secretid: n, forcedeletewithoutrecovery: true }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.name) ? null : p)); },
  });

  const handleRowClick = (row: SecretListEntry) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.name);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🗝️" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Description</td><td>{selectedItem.description || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Rotation</td><td>{selectedItem.rotationenabled ? "Enabled" : "Disabled"}</td></tr>
            {selectedItem.kmskeyid && <tr><td style={{ fontWeight: 600 }}>KMS Key</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedItem.kmskeyid}</td></tr>}
            {selectedItem.createddate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{fmtDate(selectedItem.createddate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🗝️" title={t("services.secretsmanager.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.secretsmanager.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.secretsmanager.create")}</button>
      <button className="btn btn-danger" disabled={selectedNames.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedNames.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({selectedNames.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.secretsmanager.title") }, { label: t("services.secretsmanager.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedNames.size} label={t("common.selectedCount", { count: selectedNames.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-secretsmanager">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<SecretListEntry>(selectedNames, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.name), ...columns]} data={items} getRowId={(row) => row.name} onRowClick={handleRowClick} selectedId={selectedItem?.name} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.secretsmanager.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName}>
        <label>{t("services.secretsmanager.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.secretsmanager.placeholder")} className="modal-input" /></label>
        <label>{t("services.secretsmanager.descLabel")}<input value={formDesc} onChange={(e) => setFormDesc(e.target.value)} placeholder={t("services.secretsmanager.descPlaceholder")} className="modal-input" /></label>
        <label>{t("services.secretsmanager.secretValueLabel")}<textarea value={formSecretValue} onChange={(e) => setFormSecretValue(e.target.value)} placeholder={t("services.secretsmanager.secretValuePlaceholder")} className="modal-input" rows={3} /></label>
        <label>{t("services.secretsmanager.kmsKeyLabel")}<input value={formKmsKeyId} onChange={(e) => setFormKmsKeyId(e.target.value)} placeholder={t("services.secretsmanager.kmsKeyPlaceholder")} className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.secretsmanager.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedNames.size} ${t("services.secretsmanager.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedNames))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
