import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { RDSService, type DBInstance, type DBCluster } from "@/gen/rds_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, MonoCell, BooleanBadge, DateCell, useServiceClient } from "@/components/shared/service-page";
import { DetailPanel, DetailEmpty } from "@/components/shared/inspector";
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
  const { t } = useTranslation();
  const instanceColumns = getInstanceColumns(t);
  const clusterColumns = getClusterColumns(t);

  const [tab, setTab] = useState<TabKey>("instances");
  const [selectedInstance, setSelectedInstance] = useState<DBInstance | null>(null);
  const [selectedCluster, setSelectedCluster] = useState<DBCluster | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [engineFilter, setEngineFilter] = useState<string>("all");

  const { client } = useServiceClient(RDSService);
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

  const renderInstanceDetail = () => {
    if (!selectedInstance) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedInstance.dbinstanceidentifier} titleIcon="🗄️" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Instance ID</td><td className="cell-mono">{selectedInstance.dbinstanceidentifier}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Engine</td><td>{selectedInstance.engine}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Class</td><td className="cell-mono">{selectedInstance.dbinstanceclass}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Status</td><td><span className="badge">{String(selectedInstance.dbinstancestatus)}</span></td></tr>
            <tr><td style={{ fontWeight: 600 }}>Version</td><td>{selectedInstance.engineversion || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Encrypted</td><td><BooleanBadge value={selectedInstance.storageencrypted} /></td></tr>
            {selectedInstance.dbclusteridentifier && <tr><td style={{ fontWeight: 600 }}>Cluster</td><td className="cell-mono">{selectedInstance.dbclusteridentifier}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedInstance} />}
      </DetailPanel>
    );
  };

  const renderClusterDetail = () => {
    if (!selectedCluster) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedCluster.dbclusteridentifier} titleIcon="🗄️" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Cluster ID</td><td className="cell-mono">{selectedCluster.dbclusteridentifier}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Engine</td><td>{selectedCluster.engine}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Status</td><td><span className="badge">{String(selectedCluster.status)}</span></td></tr>
            <tr><td style={{ fontWeight: 600 }}>Version</td><td>{selectedCluster.engineversion || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Encrypted</td><td><BooleanBadge value={selectedCluster.storageencrypted} /></td></tr>
            {selectedCluster.clustercreatetime && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{new Date(selectedCluster.clustercreatetime).toLocaleString()}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedCluster} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🗄️" title={t("services.rds.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange}>
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
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={instanceColumns} data={instances} getRowId={(row) => row.dbinstanceidentifier} onRowClick={(row) => { setSelectedInstance(row); setDetailTab("detail"); }} selectedId={selectedInstance?.dbinstanceidentifier} /></div>
          {renderInstanceDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "clusters" && (clusters.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-rds-cluster">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={clusterColumns} data={clusters} getRowId={(row) => row.dbclusteridentifier} onRowClick={(row) => { setSelectedCluster(row); setDetailTab("detail"); }} selectedId={selectedCluster?.dbclusteridentifier} /></div>
          {renderClusterDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}
    </ServicePageLayout>
  );
}
