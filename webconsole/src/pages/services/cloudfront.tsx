/**
 * CloudFront service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudFrontService, type DistributionSummary, ViewerProtocolPolicy } from "@/gen/cloudfront_pb";
import { CreateDistributionRequestSchema, DistributionConfigSchema, OriginsSchema, OriginSchema, DefaultCacheBehaviorSchema } from "@/gen/cloudfront_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, FallbackCell, BooleanCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { TagSection, useTags } from "@/components/shared/tag-section";

const getColumns = (t: TFunction): ColumnDef<DistributionSummary, any>[] => [
  { accessorKey: "id", header: t("services.cloudfront.distributionIdHeader"), cell: MonoCell },
  { accessorKey: "domainname", header: t("services.cloudfront.domainHeader"), cell: MonoCell },
  { accessorKey: "status", header: t("services.cloudfront.statusHeader"), size: 100 },
  { accessorKey: "enabled", header: t("services.cloudfront.enabledHeader"), cell: BooleanCell, size: 80 },
  { accessorKey: "lastmodifiedtime", header: t("services.cloudfront.lastModifiedHeader"), cell: DateCell },
  { accessorKey: "comment", header: t("services.cloudfront.commentHeader"), cell: FallbackCell },
];

type DetailTab = "detail" | "json";

export function CloudFrontPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(CloudFrontService);
  const { queryKey } = useListKey("cloudfront");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<DistributionSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formOriginDomain, setFormOriginDomain] = useState("");
  const [formOriginId, setFormOriginId] = useState("");
  const [formComment, setFormComment] = useState("");
  const [formEnabled, setFormEnabled] = useState(true);

  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<DistributionSummary, Awaited<ReturnType<typeof client.listDistributions>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listDistributions({ marker: token || undefined }),
    getItems: (r) => r.distributionlist?.items ?? [],
    getNextToken: (r) => r.distributionlist?.nextmarker ?? "",
  });
  const items = dropEmpty(rawItems, "id");

  const createMutation = useMutation({
    mutationFn: () => client.createDistribution(create(CreateDistributionRequestSchema, {
      distributionconfig: create(DistributionConfigSchema, {
        enabled: formEnabled, comment: formComment,
        origins: create(OriginsSchema, { items: [create(OriginSchema, { id: formOriginId || "1", domainname: formOriginDomain })] }),
        defaultcachebehavior: create(DefaultCacheBehaviorSchema, { targetoriginid: formOriginId || "1", viewerprotocolpolicy: ViewerProtocolPolicy.ALLOW_ALL }),
      }),
    })),
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormOriginDomain(""); setFormOriginId(""); setFormComment(""); setFormEnabled(true); },
  });

  const deleteMutation = useMutation({
    mutationFn: (dist: DistributionSummary) => client.deleteDistribution({ id: dist.id, ifmatch: dist.etag }),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (dists: DistributionSummary[]) => Promise.allSettled(dists.map((d) => client.deleteDistribution({ id: d.id, ifmatch: d.etag }))),
    onSuccess: (_d, dists) => { invalidateList(); setShowBatchDelete(false); clearSelection(); const deletedIds = new Set(dists.map((d) => d.id)); setSelectedItem((p) => (p && deletedIds.has(p.id) ? null : p)); },
  });

  const handleRowClick = (row: DistributionSummary) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.id);
  const selectedDists = items.filter((i) => selectedIds.has(i.id));

  const { tags: itemTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (arn: string) => {
        const res = await client.listTagsForResource({ resource: arn });
        return (res.tags?.items ?? []).map((t) => ({ key: t.key, value: t.value }));
      },
      tagResource: async (arn: string, tags) => {
        await client.tagResource({ resource: arn, tags: { items: tags.map((t) => ({ key: t.key, value: t.value })) } });
      },
      untagResource: async (arn: string, tagKeys: string[]) => {
        await client.untagResource({ resource: arn, tagkeys: { items: tagKeys } });
      },
    },
    selectedItem?.arn || undefined,
  );

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const origins = selectedItem.origins?.items ?? [];
    const aliases = selectedItem.aliases?.items ?? [];
    return (
      <DetailPanel title={selectedItem.id} titleIcon="☁️" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <div className="detail-body">
            <section className="detail-section">
              <h3>{t("common.general")}</h3>
              <div className="detail-field"><span className="detail-label">{t("common.fields.id")}</span><span className="cell-mono">{selectedItem.id}</span></div>
              <div className="detail-field"><span className="detail-label">{t("common.fields.domain")}</span><span className="cell-mono">{selectedItem.domainname || "\u2014"}</span></div>
              <div className="detail-field"><span className="detail-label">{t("common.fields.status")}</span><span>{selectedItem.status || "\u2014"}</span></div>
              <div className="detail-field"><span className="detail-label">{t("common.fields.enabled")}</span><span>{selectedItem.enabled ? t("common.yes") : t("common.no")}</span></div>
              {selectedItem.comment && <div className="detail-field"><span className="detail-label">{t("common.fields.comment")}</span><span>{selectedItem.comment}</span></div>}
              {selectedItem.lastmodifiedtime && <div className="detail-field"><span className="detail-label">{t("common.fields.modified")}</span><span>{fmtDate(selectedItem.lastmodifiedtime)}</span></div>}
            </section>
            {aliases.length > 0 && <section className="detail-section"><h3>{t("common.fields.aliases")} ({aliases.length})</h3>{aliases.map((a, i) => <div key={i} className="detail-field"><span className="cell-mono">{a}</span></div>)}</section>}
            {origins.length > 0 && <section className="detail-section"><h3>{t("common.fields.origins")} ({origins.length})</h3>{origins.map((o) => <div key={o.id} className="detail-field"><span className="detail-label">{o.id}</span><span className="cell-mono">{o.domainname || "\u2014"}</span></div>)}</section>}
            <TagSection tags={itemTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} />
          </div>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="☁️" title={t("services.cloudfront.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cloudfront.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cloudfront.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cloudfront.title") }, { label: t("services.cloudfront.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cf">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<DistributionSummary>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...columns]} data={items} getRowId={(row) => row.id} onRowClick={handleRowClick} selectedId={selectedItem?.id} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cloudfront.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formOriginDomain}>
        <label>{t("services.cloudfront.originDomainLabel")}<input value={formOriginDomain} onChange={(e) => setFormOriginDomain(e.target.value)} placeholder={t("services.cloudfront.originDomainPlaceholder")} className="modal-input" /></label>
        <label>{t("services.cloudfront.originIdLabel")}<input value={formOriginId} onChange={(e) => setFormOriginId(e.target.value)} placeholder={t("services.cloudfront.originIdPlaceholder")} className="modal-input" /></label>
        <label>{t("services.cloudfront.commentLabel")}<input value={formComment} onChange={(e) => setFormComment(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label className="checkbox-label"><input type="checkbox" checked={formEnabled} onChange={(e) => setFormEnabled(e.target.checked)} />{t("services.cloudfront.enabledLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cloudfront.delete")} name={selectedItem?.id} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.cloudfront.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(selectedDists)} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
