import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { NeptuneService, type DBInstance, type DBCluster } from "@/gen/neptune_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "instances" | "clusters";

const instanceColumns: ColumnDef<DBInstance, any>[] = [
  {
    accessorKey: "dbinstanceidentifier",
    header: "Instance ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "dbinstancestatus",
    header: "Status",
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "available" ? "badge-green" : "badge-yellow";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 100,
  },
  {
    accessorKey: "engine",
    header: "Engine",
    size: 100,
  },
  {
    accessorKey: "dbinstanceclass",
    header: "Class",
    size: 120,
  },
];

const clusterColumns: ColumnDef<DBCluster, any>[] = [
  {
    accessorKey: "dbclusteridentifier",
    header: "Cluster ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "status",
    header: "Status",
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "available" ? "badge-green" : "badge-yellow";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 100,
  },
  {
    accessorKey: "engine",
    header: "Engine",
    size: 100,
  },
];

function DBInstanceDetail({ item }: { item: DBInstance }) {
  const statusCls = item.dbinstancestatus === "available" ? "badge-green" : item.dbinstancestatus === "creating" ? "badge-yellow" : "";
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Identifier</span><span className="cell-mono">{item.dbinstanceidentifier}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.dbinstancearn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Status</span><span className={`badge ${statusCls}`}>{item.dbinstancestatus || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Engine</span><span>{item.engine} {item.engineversion}</span></div>
        <div className="detail-field"><span className="detail-label">Class</span><span>{item.dbinstanceclass || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Allocated Storage</span><span>{item.allocatedstorage ? `${item.allocatedstorage} GB` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Multi-AZ</span>{item.multiaz ? <span className="badge badge-green">Yes</span> : "No"}</div>
        <div className="detail-field"><span className="detail-label">Cluster</span><span className="cell-mono">{item.dbclusteridentifier || "\u2014"}</span></div>
      </section>

      {item.endpoint && (
        <section className="detail-section">
          <h3>Endpoint</h3>
          <div className="detail-field"><span className="detail-label">Address</span><span className="cell-mono">{(item.endpoint as any).address || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">Port</span><span>{(item.endpoint as any).port || "\u2014"}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>Configuration</h3>
        <div className="detail-field"><span className="detail-label">AZ</span><span>{item.availabilityzone || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Publicly Accessible</span>{item.publiclyaccessible ? "Yes" : "No"}</div>
        <div className="detail-field"><span className="detail-label">IAM DB Auth</span>{item.iamdatabaseauthenticationenabled ? <span className="badge badge-green">Enabled</span> : "Disabled"}</div>
        <div className="detail-field"><span className="detail-label">Deletion Protection</span>{item.deletionprotection ? <span className="badge badge-green">On</span> : "Off"}</div>
        <div className="detail-field"><span className="detail-label">Backup Retention</span><span>{item.backupretentionperiod ? `${item.backupretentionperiod} days` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Preferred Backup</span><span>{item.preferredbackupwindow || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Preferred Maintenance</span><span>{item.preferredmaintenancewindow || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Port</span><span>{item.dbinstanceport || "\u2014"}</span></div>
      </section>

      {item.kmskeyid && (
        <section className="detail-section">
          <h3>Encryption</h3>
          <div className="detail-field"><span className="detail-label">KMS Key</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.kmskeyid}</span></div>
        </section>
      )}

      {item.enabledcloudwatchlogsexports?.length > 0 && (
        <section className="detail-section">
          <h3>CloudWatch Logs Exports</h3>
          {item.enabledcloudwatchlogsexports.map((l, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span>{l}</span></div>
          ))}
        </section>
      )}

      {item.instancecreatetime && (
        <section className="detail-section">
          <h3>Timestamps</h3>
          <div className="detail-field"><span className="detail-label">Created</span><span>{new Date(item.instancecreatetime).toLocaleString()}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

function DBClusterDetail({ item }: { item: DBCluster }) {
  const statusCls = item.status === "available" ? "badge-green" : item.status === "creating" ? "badge-yellow" : "";
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Identifier</span><span className="cell-mono">{item.dbclusteridentifier}</span></div>
        <div className="detail-field"><span className="detail-label">Status</span><span className={`badge ${statusCls}`}>{item.status || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Engine</span><span>{item.engine} {item.engineversion}</span></div>
        <div className="detail-field"><span className="detail-label">Port</span><span>{item.port || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Multi-AZ</span>{item.multiaz ? <span className="badge badge-green">Yes</span> : "No"}</div>
        <div className="detail-field"><span className="detail-label">Deletion Protection</span>{item.deletionprotection ? <span className="badge badge-green">On</span> : "Off"}</div>
        <div className="detail-field"><span className="detail-label">IAM DB Auth</span>{item.iamdatabaseauthenticationenabled ? <span className="badge badge-green">Enabled</span> : "Disabled"}</div>
        <div className="detail-field"><span className="detail-label">Storage Encrypted</span>{item.storageencrypted ? <span className="badge badge-green">Yes</span> : "No"}</div>
      </section>

      {item.endpoint && (
        <section className="detail-section">
          <h3>Writer Endpoint</h3>
          <div className="detail-field"><span className="detail-label">Address</span><span className="cell-mono">{item.endpoint}</span></div>
        </section>
      )}

      {item.readerendpoint && (
        <section className="detail-section">
          <h3>Reader Endpoint</h3>
          <div className="detail-field"><span className="detail-label">Address</span><span className="cell-mono">{item.readerendpoint}</span></div>
        </section>
      )}

      {item.backupretentionperiod > 0 && (
        <section className="detail-section">
          <h3>Backup</h3>
          <div className="detail-field"><span className="detail-label">Retention</span><span>{item.backupretentionperiod} days</span></div>
          <div className="detail-field"><span className="detail-label">Preferred Window</span><span>{item.preferredbackupwindow || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">Preferred Maintenance</span><span>{item.preferredmaintenancewindow || "\u2014"}</span></div>
        </section>
      )}

      {item.enabledcloudwatchlogsexports?.length > 0 && (
        <section className="detail-section">
          <h3>CloudWatch Logs Exports</h3>
          {item.enabledcloudwatchlogsexports.map((l, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span>{l}</span></div>
          ))}
        </section>
      )}

      {item.clustercreatetime && (
        <section className="detail-section">
          <h3>Timestamps</h3>
          <div className="detail-field"><span className="detail-label">Created</span><span>{new Date(item.clustercreatetime).toLocaleString()}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

export function NeptunePage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("instances");
  const [selectedInstance, setSelectedInstance] = useState<DBInstance | null>(null);
  const [selectedCluster, setSelectedCluster] = useState<DBCluster | null>(null);

  const client = createClient(NeptuneService, transport);
  const { queryKey: instKey } = useListKey("neptune-instances");
  const { queryKey: clKey } = useListKey("neptune-clusters");

  const instQ = useQuery({
    queryKey: instKey,
    queryFn: async () => {
      const resp = await client.describeDBInstances({});
      return dropEmpty(resp.dbinstances ?? [], "dbinstanceidentifier");
    },
    refetchInterval: 30_000,
  });

  const clusterQ = useQuery({
    queryKey: clKey,
    queryFn: async () => {
      const resp = await client.describeDBClusters({});
      return dropEmpty(resp.dbclusters ?? [], "dbclusteridentifier");
    },
    refetchInterval: 30_000,
  });

  const query = tab === "instances" ? instQ : clusterQ;

  if (query.isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔮</span>
          <h1>Neptune</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (query.error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔮</span>
          <h1>Neptune</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(query.error) })}</div>
      </div>
    );
  }

  const instances: DBInstance[] = instQ.data ?? [];
  const clusters: DBCluster[] = clusterQ.data ?? [];

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔮</span>
        <h1>Neptune</h1>
        <div className="tab-bar">
          <button
            className={`tab-btn ${tab === "instances" ? "active" : ""}`}
            onClick={() => setTab("instances")}
          >
            Instances ({instances.length})
          </button>
          <button
            className={`tab-btn ${tab === "clusters" ? "active" : ""}`}
            onClick={() => setTab("clusters")}
          >
            Clusters ({clusters.length})
          </button>
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          {tab === "instances" && (
            <DataTable
              columns={instanceColumns}
              data={instances}
              getRowId={(row) => row.dbinstanceidentifier}
              onRowClick={setSelectedInstance}
              selectedId={selectedInstance?.dbinstanceidentifier}
            />
          )}
          {tab === "clusters" && (
            <DataTable
              columns={clusterColumns}
              data={clusters}
              getRowId={(row) => row.dbclusteridentifier}
              onRowClick={setSelectedCluster}
              selectedId={selectedCluster?.dbclusteridentifier}
            />
          )}
        </div>
        {tab === "instances" && selectedInstance && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedInstance.dbinstanceidentifier}</h2>
              <button className="detail-close" onClick={() => setSelectedInstance(null)}>✕</button>
            </div>
            <DBInstanceDetail item={selectedInstance} />
          </div>
        )}
        {tab === "clusters" && selectedCluster && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedCluster.dbclusteridentifier}</h2>
              <button className="detail-close" onClick={() => setSelectedCluster(null)}>✕</button>
            </div>
            <DBClusterDetail item={selectedCluster} />
          </div>
        )}
      </div>
    </div>
  );
}
