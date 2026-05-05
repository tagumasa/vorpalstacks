/**
 * ACM service page. Lists certificates with create/delete operations.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { ACMService, type CertificateSummary, KeyAlgorithm } from "@/gen/acm_pb";
import { RequestCertificateRequestSchema } from "@/gen/acm_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  BooleanBadge,
  BooleanCell,
  DateCell,
  fmtDate,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";

/** Column definitions for the ACM certificate table. */
const getColumns = (t: TFunction): ColumnDef<CertificateSummary, any>[] => [
  { accessorKey: "domainname", header: t("services.acm.domainHeader"), cell: MonoCell },
  {
    accessorKey: "status",
    header: t("services.acm.statusHeader"),
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "ISSUED" ? "badge-green" : v === "PENDING_VALIDATION" ? "badge-yellow" : "";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 160,
  },
  { accessorKey: "type", header: t("services.acm.typeHeader"), size: 100 },
  { accessorKey: "keyalgorithm", header: t("services.acm.keyAlgorithmHeader"), size: 120 },
  { accessorKey: "inuse", header: t("services.acm.inUseHeader"), cell: BooleanCell, size: 70 },
  { accessorKey: "notafter", header: t("services.acm.expiresHeader"), cell: DateCell },
];

/** Detail panel for an ACM certificate. */
function ACMDetail({ item }: { item: CertificateSummary }) {
  const { t } = useTranslation();
  const sans = item.subjectalternativenamesummaries ?? [];
  const keyUsages = item.keyusages ?? [];
  const extKeyUsages = item.extendedkeyusages ?? [];

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.domainLabel")}</span><span className="cell-mono">{item.domainname}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.certificatearn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.statusLabel")}</span><span className="badge">{String(item.status)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.typeLabel")}</span><span>{String(item.type)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.keyAlgorithmLabel")}</span><span>{String(item.keyalgorithm)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.inUseLabel")}</span><BooleanBadge value={item.inuse} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.managedByLabel")}</span><span>{String(item.managedby)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.renewalLabel")}</span><span>{String(item.renewaleligibility)}</span></div>
      </section>

      <section className="detail-section">
        <h3>{t("services.acm.detail.validitySection")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.createdLabel")}</span><span>{fmtDate(item.createdat)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.issuedLabel")}</span><span>{fmtDate(item.issuedat)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.notBeforeLabel")}</span><span>{fmtDate(item.notbefore)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.acm.detail.notAfterLabel")}</span><span>{fmtDate(item.notafter)}</span></div>
        {item.importedat && <div className="detail-field"><span className="detail-label">{t("services.acm.detail.importedLabel")}</span><span>{fmtDate(item.importedat)}</span></div>}
        {item.revokedat && <div className="detail-field"><span className="detail-label">{t("services.acm.detail.revokedLabel")}</span><span>{fmtDate(item.revokedat)}</span></div>}
      </section>

      {sans.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.acm.detail.sanSection")} ({sans.length})</h3>
          {sans.map((s, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono">{s}</span></div>
          ))}
        </section>
      )}

      {(keyUsages.length > 0 || extKeyUsages.length > 0) && (
        <section className="detail-section">
          <h3>{t("services.acm.detail.keyUsageSection")}</h3>
          {keyUsages.length > 0 && <div className="detail-field"><span className="detail-label">{t("services.acm.detail.usagesLabel")}</span><span>{keyUsages.join(", ")}</span></div>}
          {extKeyUsages.length > 0 && <div className="detail-field"><span className="detail-label">{t("services.acm.detail.extendedLabel")}</span><span>{extKeyUsages.join(", ")}</span></div>}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** ACM service page with certificate list and detail view. */
export function ACMPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<CertificateSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formDomain, setFormDomain] = useState("");
  const [formSAN, setFormSAN] = useState("");
  const [formKeyAlgo, setFormKeyAlgo] = useState(KeyAlgorithm.RSA_2048);

  const { client, invalidate } = useServiceClient(ACMService);
  const { queryKey } = useListKey("acm");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listCertificates({});
      return dropEmpty(resp.certificatesummarylist ?? [], "certificatearn");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: CertificateSummary[] = data ?? [];

  const createMutation = useMutation({
    mutationFn: () =>
      client.requestCertificate(
        create(RequestCertificateRequestSchema, {
          domainname: formDomain,
          subjectalternativenames: formSAN
            .split(",")
            .map((s) => s.trim())
            .filter(Boolean),
          keyalgorithm: formKeyAlgo,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormDomain("");
      setFormSAN("");
      setFormKeyAlgo(KeyAlgorithm.RSA_2048);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (arn: string) => client.deleteCertificate({ certificatearn: arn }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🔐"
      title={t("services.acm.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.acm.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.acm.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("services.acm.delete")}
            </button>
          )}
        </>
      }
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.certificatearn}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.certificatearn}
        selected={selectedItem}
        detailTitle={selectedItem?.domainname}
        onDetailClose={() => setSelectedItem(null)}
        DetailComponent={ACMDetail}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.acm.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formDomain}
      >
        <label>
          {t("services.acm.domainField")}
          <input
            value={formDomain}
            onChange={(e) => setFormDomain(e.target.value)}
            placeholder={t("services.acm.domainPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.acm.sanLabel")}
          <input
            value={formSAN}
            onChange={(e) => setFormSAN(e.target.value)}
            placeholder={t("services.acm.sanPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.acm.keyAlgorithmLabel")}
          <select
            value={formKeyAlgo}
            onChange={(e) => setFormKeyAlgo(Number(e.target.value))}
            className="modal-input"
          >
            <option value={KeyAlgorithm.RSA_2048}>RSA 2048</option>
            <option value={KeyAlgorithm.RSA_3072}>RSA 3072</option>
            <option value={KeyAlgorithm.RSA_4096}>RSA 4096</option>
            <option value={KeyAlgorithm.EC_PRIME256V1}>EC prime256v1</option>
            <option value={KeyAlgorithm.EC_SECP384R1}>EC secp384r1</option>
            <option value={KeyAlgorithm.EC_SECP521R1}>EC secp521r1</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.acm.delete")}
        name={selectedItem?.domainname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.certificatearn)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
