/**
 * Neptune service page. Tabbed view for DB instances and clusters with custom
 * detail panels for each resource type.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { NeptuneService, type DBInstance, type DBCluster } from "@/gen/neptune_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  MonoCell,
  SmallMonoCell,
  BooleanBadge,
  BooleanCell,
  DateCell,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";

/** Tab key type for the Neptune page. */
type TabKey = "instances" | "clusters";

/** Column definitions for the Neptune DB instance table. */
const getInstanceColumns = (t: TFunction): ColumnDef<DBInstance, any>[] => [
  { accessorKey: "dbinstanceidentifier", header: t("services.neptune.instanceIdHeader"), cell: MonoCell },
  {
    accessorKey: "dbinstancestatus",
    header: t("services.neptune.statusHeader"),
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "available" ? "badge-green" : "badge-yellow";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 100,
  },
  { accessorKey: "engine", header: t("services.neptune.engineHeader"), size: 80 },
  { accessorKey: "engineversion", header: t("services.neptune.engineVersionHeader"), size: 100 },
  { accessorKey: "dbinstanceclass", header: t("services.neptune.classHeader"), size: 120 },
  { accessorKey: "multiaz", header: t("services.neptune.multiAzHeader"), cell: BooleanCell, size: 70 },
  { accessorKey: "storageencrypted", header: t("services.neptune.encryptedHeader"), cell: BooleanCell, size: 70 },
  { accessorKey: "instancecreatetime", header: t("services.neptune.createdHeader"), cell: DateCell },
  { accessorKey: "dbinstancearn", header: t("services.neptune.arnHeader"), cell: SmallMonoCell },
];

/** Column definitions for the Neptune DB cluster table. */
const getClusterColumns = (t: TFunction): ColumnDef<DBCluster, any>[] => [
  { accessorKey: "dbclusteridentifier", header: t("services.neptune.clusterIdHeader"), cell: MonoCell },
  {
    accessorKey: "status",
    header: t("services.neptune.statusHeader"),
    cell: ({ getValue }) => {
      const v = String(getValue());
      const cls = v === "available" ? "badge-green" : "badge-yellow";
      return <span className={`badge ${cls}`}>{v}</span>;
    },
    size: 100,
  },
  { accessorKey: "engine", header: t("services.neptune.engineHeader"), size: 80 },
  { accessorKey: "engineversion", header: t("services.neptune.engineVersionHeader"), size: 100 },
  { accessorKey: "clustercreatetime", header: t("services.neptune.createdHeader"), cell: DateCell },
  { accessorKey: "endpoint", header: t("services.neptune.endpointHeader"), cell: SmallMonoCell },
  { accessorKey: "readerendpoint", header: t("services.neptune.readerEndpointHeader"), cell: SmallMonoCell },
];

/** Detail panel for a Neptune DB instance. */
function DBInstanceDetail({ item }: { item: DBInstance }) {
  const { t } = useTranslation();
  const statusCls = item.dbinstancestatus === "available" ? "badge-green" : item.dbinstancestatus === "creating" ? "badge-yellow" : "";
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.identifierLabel")}</span><span className="cell-mono">{item.dbinstanceidentifier}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.dbinstancearn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.statusLabel")}</span><span className={`badge ${statusCls}`}>{item.dbinstancestatus || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.engineLabel")}</span><span>{item.engine} {item.engineversion}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.classLabel")}</span><span>{item.dbinstanceclass || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.allocatedStorageLabel")}</span><span>{item.allocatedstorage ? `${item.allocatedstorage} GB` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.multiAzLabel")}</span><BooleanBadge value={item.multiaz} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.clusterLabel")}</span><span className="cell-mono">{item.dbclusteridentifier || "\u2014"}</span></div>
      </section>

      {item.endpoint && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.endpointSection")}</h3>
          {(() => {
            const ep = item.endpoint as { address?: string; port?: number } | undefined;
            return (
              <>
                <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.addressLabel")}</span><span className="cell-mono">{ep?.address || "\u2014"}</span></div>
                <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.portLabel")}</span><span>{ep?.port || "\u2014"}</span></div>
              </>
            );
          })()}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("services.neptune.detail.configurationSection")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.azLabel")}</span><span>{item.availabilityzone || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.publiclyAccessibleLabel")}</span><BooleanBadge value={item.publiclyaccessible} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.iamDbAuthLabel")}</span><BooleanBadge value={item.iamdatabaseauthenticationenabled} trueLabel="Enabled" falseLabel="Disabled" /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.deletionProtectionLabel")}</span><BooleanBadge value={item.deletionprotection} trueLabel="On" falseLabel="Off" /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.backupRetentionLabel")}</span><span>{item.backupretentionperiod ? `${item.backupretentionperiod} days` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.preferredBackupLabel")}</span><span>{item.preferredbackupwindow || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.preferredMaintenanceLabel")}</span><span>{item.preferredmaintenancewindow || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.portLabel")}</span><span>{item.dbinstanceport || "\u2014"}</span></div>
      </section>

      {item.kmskeyid && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.encryptionSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.kmsKeyLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.kmskeyid}</span></div>
        </section>
      )}

      {item.enabledcloudwatchlogsexports?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.cwLogsExportsSection")}</h3>
          {item.enabledcloudwatchlogsexports.map((l, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span>{l}</span></div>
          ))}
        </section>
      )}

      {item.instancecreatetime && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.timestampsSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.createdLabel")}</span><span>{new Date(item.instancecreatetime).toLocaleString()}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** Detail panel for a Neptune DB cluster. */
function DBClusterDetail({ item }: { item: DBCluster }) {
  const { t } = useTranslation();
  const statusCls = item.status === "available" ? "badge-green" : item.status === "creating" ? "badge-yellow" : "";
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.identifierLabel")}</span><span className="cell-mono">{item.dbclusteridentifier}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.statusLabel")}</span><span className={`badge ${statusCls}`}>{item.status || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.engineLabel")}</span><span>{item.engine} {item.engineversion}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.portLabel")}</span><span>{item.port || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.multiAzLabel")}</span><BooleanBadge value={item.multiaz} /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.deletionProtectionLabel")}</span><BooleanBadge value={item.deletionprotection} trueLabel="On" falseLabel="Off" /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.iamDbAuthLabel")}</span><BooleanBadge value={item.iamdatabaseauthenticationenabled} trueLabel="Enabled" falseLabel="Disabled" /></div>
        <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.storageEncryptedLabel")}</span><BooleanBadge value={item.storageencrypted} /></div>
      </section>

      {item.endpoint && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.writerEndpointSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.addressLabel")}</span><span className="cell-mono">{item.endpoint}</span></div>
        </section>
      )}

      {item.readerendpoint && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.readerEndpointSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.addressLabel")}</span><span className="cell-mono">{item.readerendpoint}</span></div>
        </section>
      )}

      {item.backupretentionperiod > 0 && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.backupSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.retentionLabel")}</span><span>{item.backupretentionperiod} days</span></div>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.preferredWindowLabel")}</span><span>{item.preferredbackupwindow || "\u2014"}</span></div>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.preferredMaintenanceLabel")}</span><span>{item.preferredmaintenancewindow || "\u2014"}</span></div>
        </section>
      )}

      {item.enabledcloudwatchlogsexports?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.cwLogsExportsSection")}</h3>
          {item.enabledcloudwatchlogsexports.map((l, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span>{l}</span></div>
          ))}
        </section>
      )}

      {item.clustercreatetime && (
        <section className="detail-section">
          <h3>{t("services.neptune.detail.timestampsSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.neptune.detail.createdLabel")}</span><span>{new Date(item.clustercreatetime).toLocaleString()}</span></div>
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** Neptune service page with tabbed instances/clusters view and detail panels. */
export function NeptunePage() {
  const { t } = useTranslation();
  const instanceColumns = getInstanceColumns(t);
  const clusterColumns = getClusterColumns(t);
  const [tab, setTab] = useState<TabKey>("instances");
  const [selectedInstance, setSelectedInstance] = useState<DBInstance | null>(null);
  const [selectedCluster, setSelectedCluster] = useState<DBCluster | null>(null);

  const { client } = useServiceClient(NeptuneService);
  const { queryKey: instKey } = useListKey("neptune-instances");
  const { queryKey: clKey } = useListKey("neptune-clusters");

  const instQ = useQuery({
    queryKey: instKey,
    queryFn: async () => {
      const resp = await client.describeDBInstances({});
      return dropEmpty(resp.dbinstances ?? [], "dbinstanceidentifier");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const clusterQ = useQuery({
    queryKey: clKey,
    queryFn: async () => {
      const resp = await client.describeDBClusters({});
      return dropEmpty(resp.dbclusters ?? [], "dbclusteridentifier");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const query = tab === "instances" ? instQ : clusterQ;
  const instances: DBInstance[] = instQ.data ?? [];
  const clusters: DBCluster[] = clusterQ.data ?? [];

  const tabs = [
    { key: "instances" as TabKey, label: t("services.neptune.tabs.instances"), count: instances.length },
    { key: "clusters" as TabKey, label: t("services.neptune.tabs.clusters"), count: clusters.length },
  ];

  return (
    <ServicePageLayout
      icon="🔮"
      title={t("services.neptune.title")}
      isLoading={query.isLoading}
      error={query.error}
      tabs={tabs}
      activeTab={tab}
      onTabChange={(k) => { setTab(k as TabKey); setSelectedInstance(null); setSelectedCluster(null); }}
    >
      {tab === "instances" && (
        <SplitPane
          columns={instanceColumns}
          data={instances}
          getRowId={(row) => row.dbinstanceidentifier}
          onRowClick={setSelectedInstance}
          selectedId={selectedInstance?.dbinstanceidentifier}
          selected={selectedInstance}
          detailTitle={selectedInstance?.dbinstanceidentifier}
          onDetailClose={() => setSelectedInstance(null)}
          DetailComponent={DBInstanceDetail}
        />
      )}
      {tab === "clusters" && (
        <SplitPane
          columns={clusterColumns}
          data={clusters}
          getRowId={(row) => row.dbclusteridentifier}
          onRowClick={setSelectedCluster}
          selectedId={selectedCluster?.dbclusteridentifier}
          selected={selectedCluster}
          detailTitle={selectedCluster?.dbclusteridentifier}
          onDetailClose={() => setSelectedCluster(null)}
          DetailComponent={DBClusterDetail}
        />
      )}
    </ServicePageLayout>
  );
}
