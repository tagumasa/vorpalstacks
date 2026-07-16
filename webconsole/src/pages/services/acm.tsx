/**
 * ACM service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { ACMService, type CertificateSummary, KeyAlgorithm, ValidationMethod } from "@/gen/acm_pb";
import { RequestCertificateRequestSchema } from "@/gen/acm_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { TagSection, useTags } from "@/components/shared/tag-section";

const KEY_ALGO_LABELS: Record<number, string> = {
  [KeyAlgorithm.RSA_2048]: "RSA-2048",
  [KeyAlgorithm.RSA_3072]: "RSA-3072",
  [KeyAlgorithm.RSA_4096]: "RSA-4096",
  [KeyAlgorithm.EC_PRIME256V1]: "EC-P256",
  [KeyAlgorithm.EC_SECP384R1]: "EC-P384",
};

const getColumns = (t: TFunction): ColumnDef<CertificateSummary, any>[] => [
  { accessorKey: "domainname", header: t("services.acm.domainHeader"), cell: MonoCell },
  { accessorKey: "status", header: t("services.acm.statusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 100 },
  { accessorKey: "type", header: t("services.acm.typeHeader"), size: 90 },
  { accessorKey: "keyalgorithm", header: t("services.acm.keyAlgoHeader"), cell: ({ getValue }) => { const v = getValue() as number; return KEY_ALGO_LABELS[v] ?? String(v); }, size: 100 },
  { accessorKey: "notbefore", header: t("services.acm.notBeforeHeader"), cell: DateCell },
  { accessorKey: "notafter", header: t("services.acm.notAfterHeader"), cell: DateCell },
];

type DetailTab = "detail" | "json";

export function ACMPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(ACMService);
  const { queryKey } = useListKey("acm");
  const columns = getColumns(t);

  const { selected: selectedArns, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<CertificateSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formDomain, setFormDomain] = useState("");
  const [formAltNames, setFormAltNames] = useState("");
  const [formValidation, setFormValidation] = useState<ValidationMethod>(ValidationMethod.DNS);
  const [formKeyAlgo, setFormKeyAlgo] = useState<KeyAlgorithm>(KeyAlgorithm.RSA_2048);

  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<CertificateSummary, Awaited<ReturnType<typeof client.listCertificates>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listCertificates({ nexttoken: token || undefined }),
    getItems: (r) => r.certificatesummarylist ?? [],
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "certificatearn");
  

  const createMutation = useMutation({
    mutationFn: () => client.requestCertificate(create(RequestCertificateRequestSchema, {
      domainname: formDomain,
      subjectalternativenames: formAltNames ? formAltNames.split(",").map((s) => s.trim()) : [formDomain],
      validationmethod: formValidation,
      keyalgorithm: formKeyAlgo,
    })),
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormDomain(""); setFormAltNames(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (arn: string) => client.deleteCertificate({ certificatearn: arn }),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (arns: string[]) => Promise.allSettled(arns.map((arn) => client.deleteCertificate({ certificatearn: arn }))),
    onSuccess: (_d, arns) => { invalidateList(); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && arns.includes(p.certificatearn) ? null : p)); },
  });

  const handleRowClick = (row: CertificateSummary) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.certificatearn);

  const { tags: itemTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (arn: string) => {
        const res = await client.listTagsForCertificate({ certificatearn: arn });
        return (res.tags ?? []).map((t) => ({ key: t.key, value: t.value }));
      },
      tagResource: async (arn: string, tags) => {
        await client.addTagsToCertificate({ certificatearn: arn, tags: tags.map((t) => ({ key: t.key, value: t.value })) });
      },
      untagResource: async (arn: string, tagKeys: string[]) => {
        await client.removeTagsFromCertificate({ certificatearn: arn, tags: tagKeys.map((k) => ({ key: k, value: "" })) });
      },
    },
    selectedItem?.certificatearn || undefined,
  );

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const keyAlgo = KEY_ALGO_LABELS[selectedItem.keyalgorithm] ?? String(selectedItem.keyalgorithm);
    return (
      <DetailPanel title={selectedItem.domainname} titleIcon="🔒" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <><table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.domain")}</td><td className="cell-mono">{selectedItem.domainname}</td></tr>
            <tr><td className="detail-label">{t("common.fields.status")}</td><td><span className="badge">{String(selectedItem.status)}</span></td></tr>
            <tr><td className="detail-label">{t("common.fields.type")}</td><td>{selectedItem.type || "\u2014"}</td></tr>
            <tr><td className="detail-label">{t("common.fields.keyAlgorithm")}</td><td>{keyAlgo}</td></tr>
            {selectedItem.certificatearn && <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedItem.certificatearn}</td></tr>}
            {selectedItem.notbefore && <tr><td className="detail-label">{t("common.fields.notBefore")}</td><td>{fmtDate(selectedItem.notbefore)}</td></tr>}
            {selectedItem.notafter && <tr><td className="detail-label">{t("common.fields.notAfter")}</td><td>{fmtDate(selectedItem.notafter)}</td></tr>}
          </tbody></table>
          <TagSection tags={itemTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} /></>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🔒" title={t("services.acm.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.acm.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.acm.create")}</button>
      <button className="btn btn-danger" disabled={selectedArns.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedArns.size > 0 && <span className="batch-count">({selectedArns.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.acm.title") }, { label: t("services.acm.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedArns.size} label={t("common.selectedCount", { count: selectedArns.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-acm">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<CertificateSummary>(selectedArns, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.certificatearn), ...columns]} data={items} getRowId={(row) => row.certificatearn} onRowClick={handleRowClick} selectedId={selectedItem?.certificatearn} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.acm.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formDomain}>
        <label>{t("services.acm.domainField")}<input value={formDomain} onChange={(e) => setFormDomain(e.target.value)} placeholder={t("services.acm.placeholder")} className="modal-input" /></label>
        <label>{t("services.acm.altNamesLabel")}<input value={formAltNames} onChange={(e) => setFormAltNames(e.target.value)} placeholder={t("services.acm.altNamesPlaceholder")} className="modal-input" /></label>
        <label>{t("services.acm.validationLabel")}<select value={formValidation} onChange={(e) => setFormValidation(Number(e.target.value))} className="modal-select"><option value={ValidationMethod.DNS}>DNS</option><option value={ValidationMethod.EMAIL}>Email</option></select></label>
        <label>{t("services.acm.keyAlgoLabel")}<select value={formKeyAlgo} onChange={(e) => setFormKeyAlgo(Number(e.target.value))} className="modal-select">
          <option value={KeyAlgorithm.RSA_2048}>RSA-2048</option>
          <option value={KeyAlgorithm.RSA_3072}>RSA-3072</option>
          <option value={KeyAlgorithm.RSA_4096}>RSA-4096</option>
          <option value={KeyAlgorithm.EC_PRIME256V1}>EC-P256</option>
          <option value={KeyAlgorithm.EC_SECP384R1}>EC-P384</option>
        </select></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.acm.delete")} name={selectedItem?.domainname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.certificatearn)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedArns.size} ${t("services.acm.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedArns))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
