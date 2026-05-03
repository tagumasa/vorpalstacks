import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { CloudWatchService, type MetricAlarm, StateValue } from "@/gen/cloudwatch_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

const columns: ColumnDef<MetricAlarm, any>[] = [
  {
    accessorKey: "alarmname",
    header: "Alarm Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "statevalue",
    header: "State",
    cell: ({ getValue }) => {
      const v = getValue() as StateValue;
      if (v === StateValue.OK) return <span className="badge badge-green">OK</span>;
      if (v === StateValue.ALARM) return <span className="badge badge-red">ALARM</span>;
      return <span className="badge badge-yellow">INSUFFICIENT</span>;
    },
    size: 90,
  },
  {
    accessorKey: "namespace",
    header: "Namespace",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      return v || "\u2014";
    },
  },
  {
    accessorKey: "metricname",
    header: "Metric",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      return v || "\u2014";
    },
  },
];

function CloudWatchDetail({ item }: { item: MetricAlarm }) {
  const stateCls = item.statevalue === StateValue.OK
    ? "badge-green"
    : item.statevalue === StateValue.ALARM
      ? "badge-red"
      : "badge-yellow";

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Name</span><span className="cell-mono">{item.alarmname}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.alarmarn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">State</span><span className={`badge ${stateCls}`}>{String(item.statevalue)}</span></div>
        <div className="detail-field"><span className="detail-label">Description</span><span>{item.alarmdescription || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Actions Enabled</span>{item.actionsenabled ? <span className="badge badge-green">Yes</span> : "No"}</div>
      </section>

      <section className="detail-section">
        <h3>Metric</h3>
        <div className="detail-field"><span className="detail-label">Namespace</span><span className="cell-mono">{item.namespace || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Metric</span><span className="cell-mono">{item.metricname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Statistic</span><span>{String(item.statistic) || "\u2014"}</span></div>
        {item.extendedstatistic && <div className="detail-field"><span className="detail-label">Extended Stat</span><span>{item.extendedstatistic}</span></div>}
        <div className="detail-field"><span className="detail-label">Period</span><span>{item.period ? `${item.period}s` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Evaluation Periods</span><span>{item.evaluationperiods || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Datapoints to Alarm</span><span>{item.datapointstoalarm || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Comparison</span><span>{String(item.comparisonoperator)}</span></div>
        <div className="detail-field"><span className="detail-label">Threshold</span><span>{item.threshold}</span></div>
        {item.treatmissingdata && <div className="detail-field"><span className="detail-label">Treat Missing</span><span>{item.treatmissingdata}</span></div>}
      </section>

      {item.dimensions?.length > 0 && (
        <section className="detail-section">
          <h3>Dimensions ({item.dimensions.length})</h3>
          {item.dimensions.map((d, i) => (
            <div key={i} className="detail-field">
              <span className="detail-label">{(d as any).name || `dim ${i}`}</span>
              <span className="cell-mono">{(d as any).value || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      <section className="detail-section">
        <h3>State</h3>
        {item.statereason && <div className="detail-field"><span className="detail-label">Reason</span><span>{item.statereason}</span></div>}
        {item.stateupdatedtimestamp && (
          <div className="detail-field"><span className="detail-label">Updated</span><span>{new Date(item.stateupdatedtimestamp).toLocaleString()}</span></div>
        )}
        {item.statetransitionedtimestamp && (
          <div className="detail-field"><span className="detail-label">Transitioned</span><span>{new Date(item.statetransitionedtimestamp).toLocaleString()}</span></div>
        )}
      </section>

      {item.alarmconfigurationupdatedtimestamp && (
        <section className="detail-section">
          <h3>Configuration</h3>
          <div className="detail-field"><span className="detail-label">Last Updated</span><span>{new Date(item.alarmconfigurationupdatedtimestamp).toLocaleString()}</span></div>
        </section>
      )}

      {item.alarmactions?.length > 0 && (
        <section className="detail-section">
          <h3>Alarm Actions ({item.alarmactions.length})</h3>
          {item.alarmactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      {item.okactions?.length > 0 && (
        <section className="detail-section">
          <h3>OK Actions ({item.okactions.length})</h3>
          {item.okactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      {item.insufficientdataactions?.length > 0 && (
        <section className="detail-section">
          <h3>Insufficient Data Actions ({item.insufficientdataactions.length})</h3>
          {item.insufficientdataactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

export function CloudWatchPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<MetricAlarm | null>(null);

  const client = createClient(CloudWatchService, transport);
  const { queryKey } = useListKey("cloudwatch");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.describeAlarms({}),
    refetchInterval: 30_000,
  });

  const items: MetricAlarm[] = dropEmpty(data?.metricalarms ?? [], "alarmname");

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📊</span>
          <h1>CloudWatch</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📊</span>
          <h1>CloudWatch</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📊</span>
        <h1>CloudWatch</h1>
        <span className="resource-count">{items.length} alarms</span>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.alarmname}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.alarmname}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.alarmname}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <CloudWatchDetail item={selectedItem} />
          </div>
        )}
      </div>
    </div>
  );
}
