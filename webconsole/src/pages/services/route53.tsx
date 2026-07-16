/**
 * Route 53 service page — 3-panel inspector layout.
 * Supports hosted zone CRUD + drill-down record set management.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  Route53Service,
  type HostedZone,
  type ResourceRecordSet,
  RRType,
  ChangeAction,
  TagResourceType,
} from "@/gen/route53_pb";
import {
  CreateHostedZoneRequestSchema,
  HostedZoneConfigSchema,
  ChangeResourceRecordSetsRequestSchema,
  ChangeBatchSchema,
  ChangeSchema,
  ResourceRecordSetSchema,
  ResourceRecordSchema,
} from "@/gen/route53_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { TagSection, useTags } from "@/components/shared/tag-section";

const getZoneColumns = (t: TFunction): ColumnDef<HostedZone, any>[] => [
  { accessorKey: "name", header: t("services.route53.zoneNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.route53.zoneIdHeader"), cell: MonoCell, size: 120 },
  { accessorKey: "resourcerecordsetcount", header: t("services.route53.recordCountHeader"), size: 80 },
  { accessorKey: "config.comment", header: t("services.route53.commentHeader") },
];

const RR_TYPE_LABELS: Record<number, string> = {
  [RRType.A]: "A", [RRType.AAAA]: "AAAA", [RRType.CNAME]: "CNAME", [RRType.MX]: "MX",
  [RRType.TXT]: "TXT", [RRType.NS]: "NS", [RRType.SOA]: "SOA", [RRType.SRV]: "SRV",
  [RRType.PTR]: "PTR", [RRType.CAA]: "CAA", [RRType.DS]: "DS", [RRType.SSHFP]: "SSHFP",
  [RRType.TLSA]: "TLSA", [RRType.NAPTR]: "NAPTR", [RRType.HTTPS]: "HTTPS", [RRType.SPF]: "SPF",
  [RRType.SVCB]: "SVCB",
};

const COMMON_RR_TYPES = [RRType.A, RRType.AAAA, RRType.CNAME, RRType.MX, RRType.TXT, RRType.NS, RRType.SRV, RRType.CAA];

function getRecordValues(rrs: ResourceRecordSet): string {
  return (rrs.resourcerecords ?? []).map((r) => r.value).join(", ") || "\u2014";
}

type DetailTab = "detail" | "json";
type View = { type: "zones" } | { type: "records"; zone: HostedZone };

export function Route53Page() {
  const { t } = useTranslation();
  const { client } = useServiceClient(Route53Service);
  const { queryKey } = useListKey("route53");
  const zoneColumns = getZoneColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedZone, setSelectedZone] = useState<HostedZone | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formComment, setFormComment] = useState("");
  const [formPrivate, setFormPrivate] = useState(false);
  const [view, setView] = useState<View>({ type: "zones" });

  // ── Record set view state ────────────────────────────────────
  const [selectedRecord, setSelectedRecord] = useState<ResourceRecordSet | null>(null);
  const [recordDetailTab, setRecordDetailTab] = useState<DetailTab>("detail");
  const [showCreateRecord, setShowCreateRecord] = useState(false);
  const [showDeleteRecord, setShowDeleteRecord] = useState(false);
  const [formRecName, setFormRecName] = useState("");
  const [formRecType, setFormRecType] = useState(RRType.A);
  const [formRecTtl, setFormRecTtl] = useState("300");
  const [formRecValue, setFormRecValue] = useState("");

  // ── Hosted zone list ─────────────────────────────────────────
  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<HostedZone, Awaited<ReturnType<typeof client.listHostedZones>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listHostedZones({ marker: token || undefined }),
    getItems: (r) => r.hostedzones ?? [],
    getNextToken: (r) => r.nextmarker ?? "",
  });
  const items = dropEmpty(rawItems, "id");

  // ── Zone tags ────────────────────────────────────────────────
  const { tags: zoneTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (resourceId: string) => {
        const res = await client.listTagsForResource({ resourceid: resourceId, resourcetype: TagResourceType.HOSTEDZONE });
        return (res.resourcetagset?.tags ?? []).map((t) => ({ key: t.key, value: t.value }));
      },
      tagResource: async () => { throw new Error("Route53 tags are read-only in this console"); },
      untagResource: async () => { throw new Error("Route53 tags are read-only in this console"); },
    },
    selectedZone?.id || undefined,
  );

  // ── Record sets (for drill-down view) ────────────────────────
  const recordsQuery = useQuery({
    queryKey: ["route53", "records", view.type === "records" ? view.zone.id : ""],
    queryFn: () => client.listResourceRecordSets({
      hostedzoneid: view.type === "records" ? view.zone.id : "",
    }),
    enabled: view.type === "records",
  });
  const recordSets = dropEmpty(recordsQuery.data?.resourcerecordsets ?? [], "name");

  // ── Mutations ────────────────────────────────────────────────
  const createZoneMutation = useMutation({
    mutationFn: () => client.createHostedZone(create(CreateHostedZoneRequestSchema, { callerreference: `vs-${Date.now()}`, name: formName, hostedzoneconfig: create(HostedZoneConfigSchema, { comment: formComment, privatezone: formPrivate }) })),
    onSuccess: () => { invalidateList(); setShowCreate(false); setFormName(""); setFormComment(""); setFormPrivate(false); },
  });

  const deleteZoneMutation = useMutation({
    mutationFn: (zoneId: string) => client.deleteHostedZone({ id: zoneId }),
    onSuccess: () => { invalidateList(); setShowDelete(false); setSelectedZone(null); clearSelection(); },
  });

  const batchDeleteZoneMutation = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteHostedZone({ id }))),
    onSuccess: (_d, ids) => { invalidateList(); setShowBatchDelete(false); clearSelection(); setSelectedZone((p) => (p && ids.includes(p.id) ? null : p)); },
  });

  const createRecordMutation = useMutation({
    mutationFn: () => {
      const zoneId = view.type === "records" ? view.zone.id : "";
      return client.changeResourceRecordSets(create(ChangeResourceRecordSetsRequestSchema, {
        hostedzoneid: zoneId,
        changebatch: create(ChangeBatchSchema, {
          changes: [create(ChangeSchema, {
            action: ChangeAction.CREATE,
            resourcerecordset: create(ResourceRecordSetSchema, {
              name: formRecName,
              type: formRecType,
              ttl: BigInt(formRecTtl || "300"),
              resourcerecords: formRecValue.split("\n").filter(Boolean).map((v) => create(ResourceRecordSchema, { value: v.trim() })),
            }),
          })],
        }),
      }));
    },
    onSuccess: () => { recordsQuery.refetch(); setShowCreateRecord(false); setFormRecName(""); setFormRecType(RRType.A); setFormRecTtl("300"); setFormRecValue(""); },
  });

  const deleteRecordMutation = useMutation({
    mutationFn: (rrs: ResourceRecordSet) => {
      const zoneId = view.type === "records" ? view.zone.id : "";
      return client.changeResourceRecordSets(create(ChangeResourceRecordSetsRequestSchema, {
        hostedzoneid: zoneId,
        changebatch: create(ChangeBatchSchema, {
          changes: [create(ChangeSchema, {
            action: ChangeAction.DELETE,
            resourcerecordset: create(ResourceRecordSetSchema, {
              name: rrs.name,
              type: rrs.type,
              ttl: rrs.ttl,
              resourcerecords: rrs.resourcerecords ?? [],
            }),
          })],
        }),
      }));
    },
    onSuccess: () => { recordsQuery.refetch(); setShowDeleteRecord(false); setSelectedRecord(null); },
  });

  // ── Handlers ─────────────────────────────────────────────────
  const allIds = items.map((i) => i.id);

  const handleZoneRowClick = (row: HostedZone) => {
    setSelectedZone(row);
    setDetailTab("detail");
  };

  const handleZoneDoubleClick = (row: HostedZone) => {
    setView({ type: "records", zone: row });
    setSelectedRecord(null);
  };

  const handleRecordRowClick = (row: ResourceRecordSet) => {
    setSelectedRecord(row);
    setRecordDetailTab("detail");
  };

  // ── Zone detail panel ────────────────────────────────────────
  const renderZoneDetailPanel = () => {
    if (!selectedZone) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedZone.name} titleIcon="🌐" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<>
        <button className="btn btn-primary btn-sm" onClick={() => handleZoneDoubleClick(selectedZone)}>{t("services.route53.viewRecords")}</button>
        <button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>
      </>}>
        {detailTab === "detail" ? (
          <><table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.name")}</td><td className="cell-mono">{selectedZone.name}</td></tr>
            <tr><td className="detail-label">{t("common.fields.zoneId")}</td><td className="cell-mono">{selectedZone.id}</td></tr>
            <tr><td className="detail-label">{t("common.fields.records")}</td><td>{String(selectedZone.resourcerecordsetcount)}</td></tr>
            {selectedZone.config?.comment && <tr><td className="detail-label">{t("common.fields.comment")}</td><td>{selectedZone.config.comment}</td></tr>}
            {selectedZone.config && <tr><td className="detail-label">{t("common.fields.private")}</td><td>{selectedZone.config.privatezone ? t("common.yes") : t("common.no")}</td></tr>}
          </tbody></table>
          <TagSection tags={zoneTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} readOnly /></>
        ) : <JsonViewer data={selectedZone} />}
      </DetailPanel>
    );
  };

  // ── Record set columns ───────────────────────────────────────
  const recordColumns: ColumnDef<ResourceRecordSet, any>[] = [
    { accessorKey: "name", header: t("services.route53.recNameHeader"), cell: MonoCell },
    { accessorKey: "type", header: t("services.route53.recTypeHeader"), cell: ({ getValue }) => <span className="badge">{RR_TYPE_LABELS[getValue() as number] ?? String(getValue())}</span>, size: 80 },
    { accessorKey: "ttl", header: t("services.route53.recTtlHeader"), cell: ({ getValue }) => <span className="cell-mono">{String(getValue())}</span>, size: 80 },
    { id: "values", header: t("services.route53.recValueHeader"), cell: ({ row }) => <span className="cell-mono cell-long">{getRecordValues(row.original)}</span> },
  ];

  // ── Record set detail panel ──────────────────────────────────
  const renderRecordDetailPanel = () => {
    if (!selectedRecord) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedRecord.name} titleIcon=" DNS" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={recordDetailTab} onTabChange={(k) => setRecordDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteRecord(true)}>{t("common.delete")}</button>}>
        {recordDetailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label">{t("common.fields.name")}</td><td className="cell-mono">{selectedRecord.name}</td></tr>
            <tr><td className="detail-label">{t("services.route53.recTypeHeader")}</td><td><span className="badge">{RR_TYPE_LABELS[selectedRecord.type] ?? String(selectedRecord.type)}</span></td></tr>
            <tr><td className="detail-label">{t("services.route53.recTtlHeader")}</td><td className="cell-mono">{String(selectedRecord.ttl)}</td></tr>
            <tr><td className="detail-label">{t("services.route53.recValueHeader")}</td><td className="cell-mono">{getRecordValues(selectedRecord)}</td></tr>
          </tbody></table>
        ) : <JsonViewer data={selectedRecord} />}
      </DetailPanel>
    );
  };

  // ── Render ───────────────────────────────────────────────────
  if (view.type === "records") {
    const zone = view.zone;
    return (
      <ServicePageLayout icon="🌐" title={`${t("services.route53.title")} — ${zone.name}`} isLoading={recordsQuery.isLoading} error={recordsQuery.error} actions={<>
        <button className="btn btn-secondary" onClick={() => { setView({ type: "zones" }); setSelectedRecord(null); }}>{t("services.route53.backToZones")}</button>
        <button className="btn btn-primary" onClick={() => setShowCreateRecord(true)}>{t("services.route53.createRecord")}</button>
      </>}>
        <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.route53.title"), onClick: () => { setView({ type: "zones" }); setSelectedRecord(null); } }, { label: zone.name }]} /></div>
        {recordSets.length > 0 ? (
          <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-route53-records">
            <div className="flex-fill-scroll"><DataTable columns={recordColumns} data={recordSets} getRowId={(row) => `${row.name}-${row.type}`} onRowClick={handleRecordRowClick} selectedId={selectedRecord ? `${selectedRecord.name}-${selectedRecord.type}` : undefined} /></div>
            {renderRecordDetailPanel()}
          </Splitter>
        ) : <div className="empty-state">{t("common.noData")}</div>}

        <ServiceCreateModal open={showCreateRecord} onClose={() => setShowCreateRecord(false)} title={t("services.route53.createRecord")} error={createRecordMutation.error} isPending={createRecordMutation.isPending} onCreate={() => createRecordMutation.mutate()} disabled={!formRecName || !formRecValue}>
          <label>{t("services.route53.recNameField")}<input value={formRecName} onChange={(e) => setFormRecName(e.target.value)} placeholder="www.example.com." className="modal-input" /></label>
          <label>{t("services.route53.recTypeField")}<select value={formRecType} onChange={(e) => setFormRecType(Number(e.target.value))} className="modal-input">{COMMON_RR_TYPES.map((rt) => <option key={rt} value={rt}>{RR_TYPE_LABELS[rt]}</option>)}</select></label>
          <label>{t("services.route53.recTtlField")}<input value={formRecTtl} onChange={(e) => setFormRecTtl(e.target.value)} type="number" className="modal-input" /></label>
          <label>{t("services.route53.recValueField")}<textarea value={formRecValue} onChange={(e) => setFormRecValue(e.target.value)} placeholder={"192.0.2.1\n192.0.2.2"} className="modal-input" rows={3} style={{ fontFamily: "monospace", fontSize: "0.85em" }} /></label>
        </ServiceCreateModal>
        <ServiceDeleteDialog open={showDeleteRecord && !!selectedRecord} title={t("services.route53.deleteRecord")} name={selectedRecord?.name} error={deleteRecordMutation.error} isPending={deleteRecordMutation.isPending} onConfirm={() => selectedRecord && deleteRecordMutation.mutate(selectedRecord)} onClose={() => setShowDeleteRecord(false)} />
      </ServicePageLayout>
    );
  }

  return (
    <ServicePageLayout icon="🌐" title={t("services.route53.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.route53.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.route53.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.route53.title") }, { label: t("services.route53.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-route53">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<HostedZone>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.id), ...zoneColumns]} data={items} getRowId={(row) => row.id} onRowClick={handleZoneRowClick} selectedId={selectedZone?.id} hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore} /></div>
          {renderZoneDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.route53.create")} error={createZoneMutation.error} isPending={createZoneMutation.isPending} onCreate={() => createZoneMutation.mutate()} disabled={!formName}>
        <label>{t("services.route53.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.route53.placeholder")} className="modal-input" /></label>
        <label>{t("services.route53.commentLabel")}<input value={formComment} onChange={(e) => setFormComment(e.target.value)} placeholder={t("services.route53.commentPlaceholder")} className="modal-input" /></label>
        <label className="checkbox-label"><input type="checkbox" checked={formPrivate} onChange={(e) => setFormPrivate(e.target.checked)} />{t("services.route53.privateZoneLabel")}</label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedZone} title={t("services.route53.delete")} name={selectedZone?.name} error={deleteZoneMutation.error} isPending={deleteZoneMutation.isPending} onConfirm={() => selectedZone && deleteZoneMutation.mutate(selectedZone.id)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.route53.countLabel")}`} error={batchDeleteZoneMutation.error} isPending={batchDeleteZoneMutation.isPending} onConfirm={() => batchDeleteZoneMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
