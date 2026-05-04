/**
 * CloudFront service page. Lists distributions with create/delete operations and a
 * custom detail panel showing origins, aliases, and cache behaviours.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudFrontService, type DistributionSummary, ViewerProtocolPolicy } from "@/gen/cloudfront_pb";
import { CreateDistributionRequestSchema, DistributionConfigSchema, OriginsSchema, OriginSchema, DefaultCacheBehaviorSchema } from "@/gen/cloudfront_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  FallbackCell,
  BooleanBadge,
  BooleanCell,
  DateCell,
  fmtDate,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";

/** Column definitions for the CloudFront distribution table. */
const getColumns = (t: TFunction): ColumnDef<DistributionSummary, any>[] => [
  { accessorKey: "id", header: t("services.cloudfront.distributionIdHeader"), cell: MonoCell },
  { accessorKey: "domainname", header: t("services.cloudfront.domainHeader"), cell: MonoCell },
  { accessorKey: "status", header: t("services.cloudfront.statusHeader"), size: 100 },
  {
    accessorKey: "enabled",
    header: t("services.cloudfront.enabledHeader"),
    cell: BooleanCell,
    size: 80,
  },
  { accessorKey: "lastmodifiedtime", header: t("services.cloudfront.lastModifiedHeader"), cell: DateCell },
  { accessorKey: "comment", header: t("services.cloudfront.commentHeader"), cell: FallbackCell },
];

/** Detail panel for a CloudFront distribution. */
function CloudFrontDetail({ item }: { item: DistributionSummary }) {
  const { t } = useTranslation();
  const origins = item.origins?.items ?? [];
  const aliases = item.aliases?.items ?? [];
  const cacheBehaviors = item.cachebehaviors?.items ?? [];
  const defBehavior = item.defaultcachebehavior;

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.idLabel")}</span><span className="cell-mono">{item.id}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.domainLabel")}</span><span className="cell-mono">{item.domainname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.statusLabel")}</span><span>{item.status || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.enabledLabel")}</span><BooleanBadge value={item.enabled} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.ipv6Label")}</span><BooleanBadge value={item.isipv6enabled} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.stagingLabel")}</span><BooleanBadge value={item.staging} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.commentLabel")}</span><span>{item.comment || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.etagLabel")}</span><span className="cell-mono">{item.etag || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.lastModifiedLabel")}</span><span>{fmtDate(item.lastmodifiedtime)}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.webAclIdLabel")}</span><span className="cell-mono">{item.webaclid || "\u2014"}</span></div>
      </section>

      {aliases.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudfront.detail.aliasesSection")} ({aliases.length})</h3>
          {aliases.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono">{a}</span></div>
          ))}
        </section>
      )}

      {origins.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudfront.detail.originsSection")} ({origins.length})</h3>
          {origins.map((o) => (
            <div key={o.id} className="detail-field">
              <span className="detail-label">{o.id}</span>
              <span className="cell-mono">{o.domainname || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      {defBehavior && (
        <section className="detail-section">
          <h3>{t("services.cloudfront.detail.defaultCacheBehaviourSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.targetOriginLabel")}</span><span className="cell-mono">{defBehavior.targetoriginid || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.viewerProtocolLabel")}</span><span>{String(defBehavior.viewerprotocolpolicy) || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.allowedMethodsLabel")}</span><span>{defBehavior.allowedmethods?.items?.join(", ") || "\u2014"}</span></div>
        </section>
      )}

      {cacheBehaviors.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudfront.detail.cacheBehavioursSection")} ({cacheBehaviors.length})</h3>
          {cacheBehaviors.map((cb) => (
            <div key={cb.pathpattern} className="detail-field">
              <span className="detail-label">{cb.pathpattern}</span>
              <span className="cell-mono">{cb.targetoriginid || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      {item.viewercertificate && (
        <section className="detail-section">
          <h3>{t("services.cloudfront.detail.viewerCertificateSection")}</h3>
          {(() => {
            const cert = item.viewercertificate;
            return (
              <>
                <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.certificateSourceLabel")}</span><span>{String(cert?.certificatesource ?? "") || "\u2014"}</span></div>
                <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.minimumProtocolLabel")}</span><span>{String(cert?.minimumprotocolversion ?? "") || "\u2014"}</span></div>
                <div className="detail-field"><span className="detail-label">{t("services.cloudfront.detail.sslSupportLabel")}</span><span>{String(cert?.sslsupportmethod ?? "") || "\u2014"}</span></div>
              </>
            );
          })()}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** CloudFront service page with distribution list, detail, and delete operations. */
export function CloudFrontPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<DistributionSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formOriginDomain, setFormOriginDomain] = useState("");
  const [formOriginId, setFormOriginId] = useState("");
  const [formComment, setFormComment] = useState("");
  const [formEnabled, setFormEnabled] = useState(true);

  const { client, invalidate } = useServiceClient(CloudFrontService);
  const { queryKey } = useListKey("cloudfront");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listDistributions({});
      return dropEmpty(resp.distributionlist?.items ?? [], "id");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: DistributionSummary[] = data ?? [];

  const createMutation = useMutation({
    mutationFn: () =>
      client.createDistribution(
        create(CreateDistributionRequestSchema, {
          distributionconfig: create(DistributionConfigSchema, {
            enabled: formEnabled,
            comment: formComment,
            origins: create(OriginsSchema, {
              quantity: 1,
              items: [
                create(OriginSchema, {
                  id: formOriginId || "default-origin",
                  domainname: formOriginDomain,
                }),
              ],
            }),
            defaultcachebehavior: create(DefaultCacheBehaviorSchema, {
              targetoriginid: formOriginId || "default-origin",
              viewerprotocolpolicy: ViewerProtocolPolicy.ALLOW_ALL,
            }),
          }),
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormOriginDomain("");
      setFormOriginId("");
      setFormComment("");
      setFormEnabled(true);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (dist: DistributionSummary) =>
      client.deleteDistribution({
        id: dist.id,
        ifmatch: dist.etag,
      }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="☁️"
      title={t("services.cloudfront.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cloudfront.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cloudfront.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("services.cloudfront.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "cloudfront-items" }}
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.id}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.id}
        selected={selectedItem}
        detailTitle={selectedItem?.id}
        onDetailClose={() => setSelectedItem(null)}
        DetailComponent={CloudFrontDetail}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.cloudfront.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formOriginDomain}
      >
        <label>
          {t("services.cloudfront.originDomainLabel")}
          <input
            value={formOriginDomain}
            onChange={(e) => setFormOriginDomain(e.target.value)}
            placeholder={t("services.cloudfront.originDomainPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudfront.originIdLabel")}
          <input
            value={formOriginId}
            onChange={(e) => setFormOriginId(e.target.value)}
            placeholder={t("services.cloudfront.originIdPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudfront.commentLabel")}
          <input
            value={formComment}
            onChange={(e) => setFormComment(e.target.value)}
            placeholder={t("services.cloudfront.commentPlaceholder")}
            className="modal-input"
          />
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formEnabled}
            onChange={(e) => setFormEnabled(e.target.checked)}
          />
          {t("services.cloudfront.enabledLabel")}
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cloudfront.delete")}
        name={selectedItem?.id}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
