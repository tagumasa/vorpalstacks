/**
 * Neptune service page — 3-panel inspector layout with CRUD operations.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { NeptuneService, type DBInstance, type DBCluster } from "@/gen/neptune_pb";
import { CreateDBInstanceMessageSchema, DeleteDBInstanceMessageSchema, CreateDBClusterMessageSchema, DeleteDBClusterMessageSchema } from "@/gen/neptune_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, BooleanBadge, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "instances" | "clusters";
type DetailTab = "detail" | "json";

const getInstanceColumns = (t: TFunction): ColumnDef<DBInstance, any>[] => [
  { accessorKey: "dbinstanceidentifier", header: t("services.neptune.instanceIdHeader"), cell: MonoCell },
  { accessorKey: "dbinstanceclass", header: t("services.neptune.instanceClassHeader"), cell: MonoCell },
  { accessorKey: "dbinstancestatus", header: t("services.neptune.statusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 90 },
  { accessorKey: "engine", header: t("services.neptune.engineHeader"), size: 80 },
  { accessorKey: "engineversion", header: t("services.neptune.engineVersionHeader"), size: 100 },
  { accessorKey: "storageencrypted", header: t("services.neptune.encryptedHeader"), cell: ({ getValue }) => <BooleanBadge value={getValue() as boolean} />, size: 80 },
];

const getClusterColumns = (t: TFunction): ColumnDef<DBCluster, any>[] => [
  { accessorKey: "dbclusteridentifier", header: t("services.neptune.clusterIdHeader"), cell: MonoCell },
  { accessorKey: "status", header: t("services.neptune.statusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 90 },
  { accessorKey: "engine", header: t("services.neptune.engineHeader"), size: 80 },
  { accessorKey: "engineversion", header: t("services.neptune.engineVersionHeader"), size: 100 },
  { accessorKey: "storageencrypted", header: t("services.neptune.encryptedHeader"), cell: ({ getValue }) => <BooleanBadge value={getValue() as boolean} />, size: 80 },
  { accessorKey: "clustercreatetime", header: t("services.neptune.createdHeader"), cell: DateCell },
];

export function NeptunePage() {
  const { t, i18n } = useTranslation();
  const instanceColumns = getInstanceColumns(t);
  const clusterColumns = getClusterColumns(t);

  const [tab, setTab] = useState<TabKey>("instances");
  const [selectedInstance, setSelectedInstance] = useState<DBInstance | null>(null);
  const [selectedCluster, setSelectedCluster] = useState<DBCluster | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  const [showCreateInst, setShowCreateInst] = useState(false);
  const [formInstId, setFormInstId] = useState("");
  const [formInstClass, setFormInstClass] = useState("db.r5.large");
  const [formInstClusterId, setFormInstClusterId] = useState("");

  const [showCreateCluster, setShowCreateCluster] = useState(false);
  const [formClusterId, setFormClusterId] = useState("");

  const [showDeleteInst, setShowDeleteInst] = useState(false);
  const [showDeleteCluster, setShowDeleteCluster] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  const instSel = useSelection<string>();
  const clusterSel = useSelection<string>();
  const selectedIds = tab === "instances" ? instSel.selected : clusterSel.selected;

  const { client, invalidate } = useServiceClient(NeptuneService);
  const { queryKey: instKey } = useListKey("neptune-instances");
  const { queryKey: clusterKey } = useListKey("neptune-clusters");

  const instQ = useQuery({ queryKey: instKey, queryFn: async () => { const r = await client.describeDBInstances({}); return dropEmpty(r.dbinstances ?? [], "dbinstanceidentifier"); }, refetchInterval: REFETCH_INTERVAL });
  const clusterQ = useQuery({ queryKey: clusterKey, queryFn: async () => { const r = await client.describeDBClusters({}); return dropEmpty(r.dbclusters ?? [], "dbclusteridentifier"); }, refetchInterval: REFETCH_INTERVAL });

  const query = tab === "instances" ? instQ : clusterQ;
  const instances = instQ.data ?? [];
  const clusters = clusterQ.data ?? [];

  const tabs = [
    { key: "instances" as TabKey, label: t("services.neptune.tabs.instances"), count: instances.length },
    { key: "clusters" as TabKey, label: t("services.neptune.tabs.clusters"), count: clusters.length },
  ];

  const handleTabChange = (k: string) => { setTab(k as TabKey); setSelectedInstance(null); setSelectedCluster(null); setDetailTab("detail"); };

  const createInstMut = useMutation({
    mutationFn: () => client.createDBInstance(create(CreateDBInstanceMessageSchema, {
      dbinstanceidentifier: formInstId, dbinstanceclass: formInstClass,
      dbclusteridentifier: formInstClusterId, engine: "neptune",
    })),
    onSuccess: () => { invalidate(instKey); setShowCreateInst(false); setFormInstId(""); setFormInstClass("db.r5.large"); setFormInstClusterId(""); },
  });

  const deleteInstMut = useMutation({
    mutationFn: (id: string) => client.deleteDBInstance(create(DeleteDBInstanceMessageSchema, { dbinstanceidentifier: id })),
    onSuccess: () => { invalidate(instKey); setShowDeleteInst(false); setSelectedInstance(null); instSel.clear(); },
  });

  const batchDeleteInstMut = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteDBInstance(create(DeleteDBInstanceMessageSchema, { dbinstanceidentifier: id })))),
    onSuccess: (_d, ids) => { invalidate(instKey); setShowBatchDelete(false); instSel.clear(); setSelectedInstance((p) => (p && ids.includes(p.dbinstanceidentifier) ? null : p)); },
  });

  const createClusterMut = useMutation({
    mutationFn: () => client.createDBCluster(create(CreateDBClusterMessageSchema, {
      dbclusteridentifier: formClusterId, engine: "neptune",
    })),
    onSuccess: () => { invalidate(clusterKey); setShowCreateCluster(false); setFormClusterId(""); },
  });

  const deleteClusterMut = useMutation({
    mutationFn: (id: string) => client.deleteDBCluster(create(DeleteDBClusterMessageSchema, { dbclusteridentifier: id })),
    onSuccess: () => { invalidate(clusterKey); setShowDeleteCluster(false); setSelectedCluster(null); clusterSel.clear(); },
  });

  const batchDeleteClusterMut = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteDBCluster(create(DeleteDBClusterMessageSchema, { dbclusteridentifier: id })))),
    onSuccess: (_d, ids) => { invalidate(clusterKey); setShowBatchDelete(false); clusterSel.clear(); setSelectedCluster((p) => (p && ids.includes(p.dbclusteridentifier) ? null : p)); },
  });

  const batchDeleteMut = tab === "instances" ? batchDeleteInstMut : batchDeleteClusterMut;
  const allInstIds = instances.map((i) => i.dbinstanceidentifier);
  const allClusterIds = clusters.map((c) => c.dbclusteridentifier);

  const renderInstanceDetail = () => {
    if (!selectedInstance) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedInstance.dbinstanceidentifier} titleIcon="🐬" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteInst(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Instance ID</td><td className="cell-mono">{selectedInstance.dbinstanceidentifier}</td></tr>
            <tr><td className="detail-label">Class</td><td className="cell-mono">{selectedInstance.dbinstanceclass}</td></tr>
            <tr><td className="detail-label">Status</td><td><span className="badge">{String(selectedInstance.dbinstancestatus)}</span></td></tr>
            <tr><td className="detail-label">Engine</td><td>{selectedInstance.engine}</td></tr>
            <tr><td className="detail-label">Version</td><td>{selectedInstance.engineversion || "\u2014"}</td></tr>
            <tr><td className="detail-label">Encrypted</td><td><BooleanBadge value={selectedInstance.storageencrypted} /></td></tr>
            {selectedInstance.endpoint?.address && <tr><td className="detail-label">Endpoint</td><td className="cell-mono">{selectedInstance.endpoint.address}:{selectedInstance.endpoint.port}</td></tr>}
            {selectedInstance.dbclusteridentifier && <tr><td className="detail-label">Cluster</td><td className="cell-mono">{selectedInstance.dbclusteridentifier}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedInstance} />}
      </DetailPanel>
    );
  };

  const renderClusterDetail = () => {
    if (!selectedCluster) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedCluster.dbclusteridentifier} titleIcon="🐬" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteCluster(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Cluster ID</td><td className="cell-mono">{selectedCluster.dbclusteridentifier}</td></tr>
            <tr><td className="detail-label">Status</td><td><span className="badge">{String(selectedCluster.status)}</span></td></tr>
            <tr><td className="detail-label">Engine</td><td>{selectedCluster.engine}</td></tr>
            <tr><td className="detail-label">Version</td><td>{selectedCluster.engineversion || "\u2014"}</td></tr>
            <tr><td className="detail-label">Encrypted</td><td><BooleanBadge value={selectedCluster.storageencrypted} /></td></tr>
            {selectedCluster.endpoint && <tr><td className="detail-label">Endpoint</td><td className="cell-mono">{selectedCluster.endpoint}</td></tr>}
            {selectedCluster.readerendpoint && <tr><td className="detail-label">Reader</td><td className="cell-mono">{selectedCluster.readerendpoint}</td></tr>}
            {selectedCluster.clustercreatetime && <tr><td className="detail-label">Created</td><td>{fmtDate(selectedCluster.clustercreatetime, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedCluster} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🐬" title={t("services.neptune.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange} actions={<>
      <button className="btn btn-primary" onClick={() => tab === "instances" ? setShowCreateInst(true) : setShowCreateCluster(true)}>{t("common.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      {tab === "instances" && (instances.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-neptune-inst">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<DBInstance>(instSel.selected, instSel.toggle, () => instSel.toggleAll(allInstIds), allInstIds, t, (row) => row.dbinstanceidentifier), ...instanceColumns]} data={instances} getRowId={(row) => row.dbinstanceidentifier} onRowClick={(row) => { setSelectedInstance(row); setDetailTab("detail"); }} selectedId={selectedInstance?.dbinstanceidentifier} /></div>
          {renderInstanceDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "clusters" && (clusters.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-neptune-cluster">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<DBCluster>(clusterSel.selected, clusterSel.toggle, () => clusterSel.toggleAll(allClusterIds), allClusterIds, t, (row) => row.dbclusteridentifier), ...clusterColumns]} data={clusters} getRowId={(row) => row.dbclusteridentifier} onRowClick={(row) => { setSelectedCluster(row); setDetailTab("detail"); }} selectedId={selectedCluster?.dbclusteridentifier} /></div>
          {renderClusterDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      <ServiceCreateModal open={showCreateInst} onClose={() => setShowCreateInst(false)} title={t("services.neptune.createInstance")} error={createInstMut.error} isPending={createInstMut.isPending} onCreate={() => createInstMut.mutate()} disabled={!formInstId}>
        <label>{t("services.neptune.instanceIdField")}<input value={formInstId} onChange={(e) => setFormInstId(e.target.value)} placeholder="my-neptune-instance" className="modal-input" /></label>
        <label>{t("services.neptune.instanceClassField")}<input value={formInstClass} onChange={(e) => setFormInstClass(e.target.value)} placeholder="db.r5.large" className="modal-input" /></label>
        <label>{t("services.neptune.clusterIdField")} {t("common.optional")}<input value={formInstClusterId} onChange={(e) => setFormInstClusterId(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
      </ServiceCreateModal>

      <ServiceCreateModal open={showCreateCluster} onClose={() => setShowCreateCluster(false)} title={t("services.neptune.createCluster")} error={createClusterMut.error} isPending={createClusterMut.isPending} onCreate={() => createClusterMut.mutate()} disabled={!formClusterId}>
        <label>{t("services.neptune.clusterIdField")}<input value={formClusterId} onChange={(e) => setFormClusterId(e.target.value)} placeholder="my-neptune-cluster" className="modal-input" /></label>
      </ServiceCreateModal>

      <ServiceDeleteDialog open={showDeleteInst && !!selectedInstance} title={t("services.neptune.deleteInstance")} name={selectedInstance?.dbinstanceidentifier} error={deleteInstMut.error} isPending={deleteInstMut.isPending} onConfirm={() => selectedInstance && deleteInstMut.mutate(selectedInstance.dbinstanceidentifier)} onClose={() => setShowDeleteInst(false)} />

      <ServiceDeleteDialog open={showDeleteCluster && !!selectedCluster} title={t("services.neptune.deleteCluster")} name={selectedCluster?.dbclusteridentifier} error={deleteClusterMut.error} isPending={deleteClusterMut.isPending} onConfirm={() => selectedCluster && deleteClusterMut.mutate(selectedCluster.dbclusteridentifier)} onClose={() => setShowDeleteCluster(false)} />

      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${tab === "instances" ? t("services.neptune.tabs.instances") : t("services.neptune.tabs.clusters")}`} error={batchDeleteMut.error} isPending={batchDeleteMut.isPending} onConfirm={() => batchDeleteMut.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
