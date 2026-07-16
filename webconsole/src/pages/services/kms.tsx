/**
 * KMS service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { KMSService, type KeyListEntry, KeyUsageType, KeySpec, OriginType } from "@/gen/kms_pb";
import { CreateKeyRequestSchema, ScheduleKeyDeletionRequestSchema } from "@/gen/kms_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { TagSection, useTags } from "@/components/shared/tag-section";

const getColumns = (t: TFunction): ColumnDef<KeyListEntry, any>[] => [
  { accessorKey: "keyid", header: t("services.kms.keyIdHeader"), cell: MonoCell },
  { accessorKey: "keyarn", header: t("services.kms.arnHeader"), cell: SmallMonoCell },
];

const KEY_SPECS: { i18nKey: string; value: KeySpec }[] = [
  { i18nKey: "services.kms.keySpecSymmetricDefault", value: KeySpec.SYMMETRIC_DEFAULT },
  { i18nKey: "services.kms.keySpecRsa2048", value: KeySpec.RSA_2048 },
  { i18nKey: "services.kms.keySpecRsa3072", value: KeySpec.RSA_3072 },
  { i18nKey: "services.kms.keySpecRsa4096", value: KeySpec.RSA_4096 },
  { i18nKey: "services.kms.keySpecEccNistP256", value: KeySpec.ECC_NIST_P256 },
  { i18nKey: "services.kms.keySpecEccNistP384", value: KeySpec.ECC_NIST_P384 },
  { i18nKey: "services.kms.keySpecHmac256", value: KeySpec.HMAC_256 },
];

const KEY_USAGES: { i18nKey: string; value: KeyUsageType }[] = [
  { i18nKey: "services.kms.keyUsageEncryptDecrypt", value: KeyUsageType.ENCRYPT_DECRYPT },
  { i18nKey: "services.kms.keyUsageSignVerify", value: KeyUsageType.SIGN_VERIFY },
  { i18nKey: "services.kms.keyUsageGenerateVerifyMac", value: KeyUsageType.GENERATE_VERIFY_MAC },
];

const ORIGINS: { i18nKey: string; value: OriginType }[] = [
  { i18nKey: "services.kms.originAwsKms", value: OriginType.AWS_KMS },
  { i18nKey: "services.kms.originExternal", value: OriginType.EXTERNAL },
];

type DetailTab = "detail" | "json";

export function KMSPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(KMSService);
  const { queryKey } = useListKey("kms");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<KeyListEntry | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formDesc, setFormDesc] = useState("");
  const [formKeySpec, setFormKeySpec] = useState(KeySpec.SYMMETRIC_DEFAULT);
  const [formKeyUsage, setFormKeyUsage] = useState(KeyUsageType.ENCRYPT_DECRYPT);
  const [formOrigin, setFormOrigin] = useState(OriginType.AWS_KMS);
  const [formMultiRegion, setFormMultiRegion] = useState(false);

  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<KeyListEntry, Awaited<ReturnType<typeof client.listKeys>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listKeys({ marker: token || undefined, limit: 1000 }),
    getItems: (r) => r.keys ?? [],
    getNextToken: (r) => r.nextmarker ?? "",
  });
  const items = dropEmpty(rawItems, "keyid");
  

  const createMutation = useMutation({
    mutationFn: () => client.createKey(create(CreateKeyRequestSchema, { description: formDesc, keyspec: formKeySpec, keyusage: formKeyUsage, origin: formOrigin, multiregion: formMultiRegion })),
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormDesc(""); setFormKeySpec(KeySpec.SYMMETRIC_DEFAULT); setFormKeyUsage(KeyUsageType.ENCRYPT_DECRYPT); setFormOrigin(OriginType.AWS_KMS); setFormMultiRegion(false); },
  });

  const deleteMutation = useMutation({
    mutationFn: (keyId: string) => client.scheduleKeyDeletion(create(ScheduleKeyDeletionRequestSchema, { keyid: keyId, pendingwindowindays: 7 })),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.scheduleKeyDeletion(create(ScheduleKeyDeletionRequestSchema, { keyid: id, pendingwindowindays: 7 })))),
    onSuccess: (_d, ids) => { invalidateList(); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && ids.includes(p.keyid) ? null : p)); },
  });

  const handleRowClick = (row: KeyListEntry) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.keyid);

  const { tags: itemTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (keyId: string) => {
        const res = await client.listResourceTags({ keyid: keyId });
        return (res.tags ?? []).map((t) => ({ key: t.tagkey, value: t.tagvalue }));
      },
      tagResource: async (keyId: string, tags) => {
        await client.tagResource({ keyid: keyId, tags: tags.map((t) => ({ tagkey: t.key, tagvalue: t.value })) });
      },
      untagResource: async (keyId: string, tagKeys: string[]) => {
        await client.untagResource({ keyid: keyId, tagkeys: tagKeys });
      },
    },
    selectedItem?.keyid || undefined,
  );

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedItem.keyid} titleIcon="🔑" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("services.kms.delete")}</button>}>
        {detailTab === "detail" ? (
          <><table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.keyId")}</td><td className="cell-mono">{selectedItem.keyid}</td></tr>
            <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedItem.keyarn}</td></tr>
          </tbody></table>
          <TagSection tags={itemTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} /></>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🔑" title={t("services.kms.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.kms.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.kms.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.kms.title") }, { label: t("services.kms.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-kms">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<KeyListEntry>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.keyid), ...columns]} data={items} getRowId={(row) => row.keyid} onRowClick={handleRowClick} selectedId={selectedItem?.keyid} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.kms.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={createMutation.isPending}>
        <label>{t("services.kms.descLabel")}<input value={formDesc} onChange={(e) => setFormDesc(e.target.value)} placeholder={t("services.kms.descPlaceholder")} className="modal-input" /></label>
        <label>{t("services.kms.keyUsageLabel")}<select value={formKeyUsage} onChange={(e) => setFormKeyUsage(Number(e.target.value))} className="modal-input">{KEY_USAGES.map((u) => <option key={u.value} value={u.value}>{t(u.i18nKey)}</option>)}</select></label>
        <label>{t("services.kms.keySpecLabel")}<select value={formKeySpec} onChange={(e) => setFormKeySpec(Number(e.target.value))} className="modal-input">{KEY_SPECS.map((s) => <option key={s.value} value={s.value}>{t(s.i18nKey)}</option>)}</select></label>
        <label>{t("services.kms.originLabel")}<select value={formOrigin} onChange={(e) => setFormOrigin(Number(e.target.value))} className="modal-input">{ORIGINS.map((o) => <option key={o.value} value={o.value}>{t(o.i18nKey)}</option>)}</select></label>
        <label className="checkbox-label"><input type="checkbox" checked={formMultiRegion} onChange={(e) => setFormMultiRegion(e.target.checked)} />{t("services.kms.multiregionLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.kms.delete")} name={selectedItem?.keyid} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.keyid)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.kms.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
