/**
 * RDS service page — 3-panel inspector layout with CRUD operations.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { RDSService, type DBInstance, type DBCluster } from "@/gen/rds_pb";
import { CreateDBInstanceMessageSchema, DeleteDBInstanceMessageSchema, CreateDBClusterMessageSchema, DeleteDBClusterMessageSchema } from "@/gen/rds_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, BooleanBadge, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "instances" | "clusters";
type DetailTab = "detail" | "json";

const getInstanceColumns = (t: TFunction): ColumnDef<DBInstance, any>[] => [
  { accessorKey: "dbinstanceidentifier", header: t("services.rds.instanceIdHeader"), cell: MonoCell },
  { accessorKey: "engine", header: t("services.rds.engineHeader"), size: 90 },
  { accessorKey: "dbinstanceclass", header: t("services.rds.instanceClassHeader"), cell: MonoCell },
  { accessorKey: "dbinstancestatus", header: t("services.rds.statusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 90 },
  { accessorKey: "engineversion", header: t("services.rds.engineVersionHeader"), size: 100 },
  { accessorKey: "storageencrypted", header: t("services.rds.encryptedHeader"), cell: ({ getValue }) => <BooleanBadge value={getValue() as boolean} />, size: 80 },
];

const getClusterColumns = (t: TFunction): ColumnDef<DBCluster, any>[] => [
  { accessorKey: "dbclusteridentifier", header: t("services.rds.clusterIdHeader"), cell: MonoCell },
  { accessorKey: "engine", header: t("services.rds.engineHeader"), size: 90 },
  { accessorKey: "status", header: t("services.rds.statusHeader"), cell: ({ getValue }) => <span className="badge">{String(getValue())}</span>, size: 90 },
  { accessorKey: "engineversion", header: t("services.rds.engineVersionHeader"), size: 100 },
  { accessorKey: "storageencrypted", header: t("services.rds.encryptedHeader"), cell: ({ getValue }) => <BooleanBadge value={getValue() as boolean} />, size: 80 },
  { accessorKey: "clustercreatetime", header: t("services.rds.createdHeader"), cell: DateCell },
];

export function RDSPage() {
  const { t, i18n } = useTranslation();
  const instanceColumns = getInstanceColumns(t);
  const clusterColumns = getClusterColumns(t);

  const [tab, setTab] = useState<TabKey>("instances");
  const [selectedInstance, setSelectedInstance] = useState<DBInstance | null>(null);
  const [selectedCluster, setSelectedCluster] = useState<DBCluster | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [engineFilter, setEngineFilter] = useState<string>("all");

  // Instance create form state
  const [showCreateInst, setShowCreateInst] = useState(false);
  const [formInstId, setFormInstId] = useState("");
  const [formInstEngine, setFormInstEngine] = useState("mysql");
  const [formInstClass, setFormInstClass] = useState("db.t3.micro");
  const [formInstStorage, setFormInstStorage] = useState(20);
  const [formInstUser, setFormInstUser] = useState("admin");
  const [formInstPass, setFormInstPass] = useState("");

  // Cluster create form state
  const [showCreateCluster, setShowCreateCluster] = useState(false);
  const [formClusterId, setFormClusterId] = useState("");
  const [formClusterEngine, setFormClusterEngine] = useState("aurora-mysql");
  const [formClusterUser, setFormClusterUser] = useState("admin");
  const [formClusterPass, setFormClusterPass] = useState("");

  // Delete dialog state
  const [showDeleteInst, setShowDeleteInst] = useState(false);
  const [showDeleteCluster, setShowDeleteCluster] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  const instSel = useSelection<string>();
  const clusterSel = useSelection<string>();
  const selectedIds = tab === "instances" ? instSel.selected : clusterSel.selected;

  const { client, invalidate } = useServiceClient(RDSService);
  const { queryKey: instKey } = useListKey("rds-instances");
  const { queryKey: clusterKey } = useListKey("rds-clusters");

  const instQ = useQuery({ queryKey: instKey, queryFn: async () => { const r = await client.describeDBInstances({}); return dropEmpty(r.dbinstances ?? [], "dbinstanceidentifier"); }, refetchInterval: REFETCH_INTERVAL });
  const clusterQ = useQuery({ queryKey: clusterKey, queryFn: async () => { const r = await client.describeDBClusters({}); return dropEmpty(r.dbclusters ?? [], "dbclusteridentifier"); }, refetchInterval: REFETCH_INTERVAL });

  const query = tab === "instances" ? instQ : clusterQ;
  const allInstances = instQ.data ?? [];
  const allClusters = clusterQ.data ?? [];

  const instances = engineFilter === "all" ? allInstances : allInstances.filter((i: DBInstance) => i.engine === engineFilter);
  const clusters = engineFilter === "all" ? allClusters : allClusters.filter((c: DBCluster) => c.engine === engineFilter);

  const tabs = [
    { key: "instances" as TabKey, label: t("services.rds.tabs.instances"), count: instances.length },
    { key: "clusters" as TabKey, label: t("services.rds.tabs.clusters"), count: clusters.length },
  ];

  const handleTabChange = (k: string) => { setTab(k as TabKey); setSelectedInstance(null); setSelectedCluster(null); setDetailTab("detail"); };

  // Instance CRUD mutations
  const createInstMut = useMutation({
    mutationFn: () => client.createDBInstance(create(CreateDBInstanceMessageSchema, {
      dbinstanceidentifier: formInstId, engine: formInstEngine, dbinstanceclass: formInstClass,
      allocatedstorage: formInstStorage, masterusername: formInstUser, masteruserpassword: formInstPass,
    })),
    onSuccess: () => { invalidate(instKey); setShowCreateInst(false); resetInstForm(); },
  });

  const deleteInstMut = useMutation({
    mutationFn: (id: string) => client.deleteDBInstance(create(DeleteDBInstanceMessageSchema, { dbinstanceidentifier: id })),
    onSuccess: () => { invalidate(instKey); setShowDeleteInst(false); setSelectedInstance(null); instSel.clear(); },
  });

  const batchDeleteInstMut = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteDBInstance(create(DeleteDBInstanceMessageSchema, { dbinstanceidentifier: id })))),
    onSuccess: (_d, ids) => { invalidate(instKey); setShowBatchDelete(false); instSel.clear(); setSelectedInstance((p) => (p && ids.includes(p.dbinstanceidentifier) ? null : p)); },
  });

  // Cluster CRUD mutations
  const createClusterMut = useMutation({
    mutationFn: () => client.createDBCluster(create(CreateDBClusterMessageSchema, {
      dbclusteridentifier: formClusterId, engine: formClusterEngine, masterusername: formClusterUser, masteruserpassword: formClusterPass,
    })),
    onSuccess: () => { invalidate(clusterKey); setShowCreateCluster(false); resetClusterForm(); },
  });

  const deleteClusterMut = useMutation({
    mutationFn: (id: string) => client.deleteDBCluster(create(DeleteDBClusterMessageSchema, { dbclusteridentifier: id })),
    onSuccess: () => { invalidate(clusterKey); setShowDeleteCluster(false); setSelectedCluster(null); clusterSel.clear(); },
  });

  const batchDeleteClusterMut = useMutation({
    mutationFn: async (ids: string[]) => Promise.allSettled(ids.map((id) => client.deleteDBCluster(create(DeleteDBClusterMessageSchema, { dbclusteridentifier: id })))),
    onSuccess: (_d, ids) => { invalidate(clusterKey); setShowBatchDelete(false); clusterSel.clear(); setSelectedCluster((p) => (p && ids.includes(p.dbclusteridentifier) ? null : p)); },
  });

  const resetInstForm = () => { setFormInstId(""); setFormInstEngine("mysql"); setFormInstClass("db.t3.micro"); setFormInstStorage(20); setFormInstUser("admin"); setFormInstPass(""); };
  const resetClusterForm = () => { setFormClusterId(""); setFormClusterEngine("aurora-mysql"); setFormClusterUser("admin"); setFormClusterPass(""); };

  const allInstIds = instances.map((i) => i.dbinstanceidentifier);
  const allClusterIds = clusters.map((c) => c.dbclusteridentifier);

  const renderInstanceDetail = () => {
    if (!selectedInstance) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedInstance.dbinstanceidentifier} titleIcon="🗄️" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteInst(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Instance ID</td><td className="cell-mono">{selectedInstance.dbinstanceidentifier}</td></tr>
            <tr><td className="detail-label">Engine</td><td>{selectedInstance.engine}</td></tr>
            <tr><td className="detail-label">Class</td><td className="cell-mono">{selectedInstance.dbinstanceclass}</td></tr>
            <tr><td className="detail-label">Status</td><td><span className="badge">{String(selectedInstance.dbinstancestatus)}</span></td></tr>
            <tr><td className="detail-label">Version</td><td>{selectedInstance.engineversion || "\u2014"}</td></tr>
            <tr><td className="detail-label">Encrypted</td><td><BooleanBadge value={selectedInstance.storageencrypted} /></td></tr>
            {selectedInstance.dbclusteridentifier && <tr><td className="detail-label">Cluster</td><td className="cell-mono">{selectedInstance.dbclusteridentifier}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedInstance} />}
      </DetailPanel>
    );
  };

  const renderClusterDetail = () => {
    if (!selectedCluster) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedCluster.dbclusteridentifier} titleIcon="🗄️" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteCluster(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Cluster ID</td><td className="cell-mono">{selectedCluster.dbclusteridentifier}</td></tr>
            <tr><td className="detail-label">Engine</td><td>{selectedCluster.engine}</td></tr>
            <tr><td className="detail-label">Status</td><td><span className="badge">{String(selectedCluster.status)}</span></td></tr>
            <tr><td className="detail-label">Version</td><td>{selectedCluster.engineversion || "\u2014"}</td></tr>
            <tr><td className="detail-label">Encrypted</td><td><BooleanBadge value={selectedCluster.storageencrypted} /></td></tr>
            {selectedCluster.clustercreatetime && <tr><td className="detail-label">Created</td><td>{fmtDate(selectedCluster.clustercreatetime, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedCluster} />}
      </DetailPanel>
    );
  };

  const batchDeleteMut = tab === "instances" ? batchDeleteInstMut : batchDeleteClusterMut;

  return (
    <ServicePageLayout icon="🗄️" title={t("services.rds.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange} actions={<>
      <button className="btn btn-primary" onClick={() => tab === "instances" ? setShowCreateInst(true) : setShowCreateCluster(true)}>{t("common.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div style={{ display: "flex", gap: 6, alignItems: "center", padding: "4px 12px", borderBottom: "1px solid var(--border)" }}>
        <span style={{ fontSize: 12, color: "var(--text-muted)" }}>Engine:</span>
        {(["all", "neptune", "mysql"] as const).map((eng) => (
          <button key={eng} className={`btn btn-sm ${engineFilter === eng ? "btn-primary" : ""}`} onClick={() => setEngineFilter(eng)} style={{ padding: "2px 8px", fontSize: 11 }}>
            {eng === "all" ? "All" : eng.charAt(0).toUpperCase() + eng.slice(1)}
          </button>
        ))}
      </div>
      {tab === "instances" && (instances.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-rds-inst">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<DBInstance>(instSel.selected, instSel.toggle, () => instSel.toggleAll(allInstIds), allInstIds, t, (row) => row.dbinstanceidentifier), ...instanceColumns]} data={instances} getRowId={(row) => row.dbinstanceidentifier} onRowClick={(row) => { setSelectedInstance(row); setDetailTab("detail"); }} selectedId={selectedInstance?.dbinstanceidentifier} /></div>
          {renderInstanceDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "clusters" && (clusters.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-rds-cluster">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<DBCluster>(clusterSel.selected, clusterSel.toggle, () => clusterSel.toggleAll(allClusterIds), allClusterIds, t, (row) => row.dbclusteridentifier), ...clusterColumns]} data={clusters} getRowId={(row) => row.dbclusteridentifier} onRowClick={(row) => { setSelectedCluster(row); setDetailTab("detail"); }} selectedId={selectedCluster?.dbclusteridentifier} /></div>
          {renderClusterDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {/* Create Instance Modal */}
      <ServiceCreateModal open={showCreateInst} onClose={() => setShowCreateInst(false)} title={t("services.rds.createInstance")} error={createInstMut.error} isPending={createInstMut.isPending} onCreate={() => createInstMut.mutate()} disabled={!formInstId}>
        <label>{t("services.rds.instanceIdField")}<input value={formInstId} onChange={(e) => setFormInstId(e.target.value)} placeholder="my-db-instance" className="modal-input" /></label>
        <label>{t("services.rds.engineField")}<select value={formInstEngine} onChange={(e) => setFormInstEngine(e.target.value)} className="modal-input"><option value="mysql">mysql</option><option value="postgres">postgres</option></select></label>
        <label>{t("services.rds.instanceClassField")}<input value={formInstClass} onChange={(e) => setFormInstClass(e.target.value)} placeholder="db.t3.micro" className="modal-input" /></label>
        <label>{t("services.rds.storageField")}<input type="number" value={formInstStorage} onChange={(e) => setFormInstStorage(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.rds.usernameField")}<input value={formInstUser} onChange={(e) => setFormInstUser(e.target.value)} placeholder="admin" className="modal-input" /></label>
        <label>{t("services.rds.passwordField")}<input type="password" value={formInstPass} onChange={(e) => setFormInstPass(e.target.value)} className="modal-input" /></label>
      </ServiceCreateModal>

      {/* Create Cluster Modal */}
      <ServiceCreateModal open={showCreateCluster} onClose={() => setShowCreateCluster(false)} title={t("services.rds.createCluster")} error={createClusterMut.error} isPending={createClusterMut.isPending} onCreate={() => createClusterMut.mutate()} disabled={!formClusterId}>
        <label>{t("services.rds.clusterIdField")}<input value={formClusterId} onChange={(e) => setFormClusterId(e.target.value)} placeholder="my-db-cluster" className="modal-input" /></label>
        <label>{t("services.rds.engineField")}<select value={formClusterEngine} onChange={(e) => setFormClusterEngine(e.target.value)} className="modal-input"><option value="aurora-mysql">aurora-mysql</option><option value="aurora-postgresql">aurora-postgresql</option></select></label>
        <label>{t("services.rds.usernameField")}<input value={formClusterUser} onChange={(e) => setFormClusterUser(e.target.value)} placeholder="admin" className="modal-input" /></label>
        <label>{t("services.rds.passwordField")}<input type="password" value={formClusterPass} onChange={(e) => setFormClusterPass(e.target.value)} className="modal-input" /></label>
      </ServiceCreateModal>

      {/* Delete Instance Dialog */}
      <ServiceDeleteDialog open={showDeleteInst && !!selectedInstance} title={t("services.rds.deleteInstance")} name={selectedInstance?.dbinstanceidentifier} error={deleteInstMut.error} isPending={deleteInstMut.isPending} onConfirm={() => selectedInstance && deleteInstMut.mutate(selectedInstance.dbinstanceidentifier)} onClose={() => setShowDeleteInst(false)} />

      {/* Delete Cluster Dialog */}
      <ServiceDeleteDialog open={showDeleteCluster && !!selectedCluster} title={t("services.rds.deleteCluster")} name={selectedCluster?.dbclusteridentifier} error={deleteClusterMut.error} isPending={deleteClusterMut.isPending} onConfirm={() => selectedCluster && deleteClusterMut.mutate(selectedCluster.dbclusteridentifier)} onClose={() => setShowDeleteCluster(false)} />

      {/* Batch Delete Dialog */}
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${tab === "instances" ? t("services.rds.tabs.instances") : t("services.rds.tabs.clusters")}`} error={batchDeleteMut.error} isPending={batchDeleteMut.isPending} onConfirm={() => batchDeleteMut.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
