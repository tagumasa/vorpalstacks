/**
 * CloudWatch service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudWatchService, type MetricAlarm, StateValue, ComparisonOperator, Statistic } from "@/gen/cloudwatch_pb";
import { PutMetricAlarmInputSchema } from "@/gen/cloudwatch_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, FallbackCell, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const STATE_VALUE_I18N: Record<number, string> = {
  [StateValue.ALARM]: "services.cloudwatch.stateAlarm",
  [StateValue.OK]: "services.cloudwatch.stateOk",
  [StateValue.INSUFFICIENT_DATA]: "services.cloudwatch.stateInsufficientData",
};

const STATISTIC_I18N: Record<number, string> = {
  [Statistic.SUM]: "services.cloudwatch.statSum", [Statistic.SAMPLECOUNT]: "services.cloudwatch.statSampleCount",
  [Statistic.AVERAGE]: "services.cloudwatch.statAverage", [Statistic.MAXIMUM]: "services.cloudwatch.statMaximum",
  [Statistic.MINIMUM]: "services.cloudwatch.statMinimum",
};

const COMPARISON_I18N: Record<number, string> = {
  [ComparisonOperator.GREATERTHANOREQUALTOTHRESHOLD]: "services.cloudwatch.compGte",
  [ComparisonOperator.GREATERTHANTHRESHOLD]: "services.cloudwatch.compGt",
  [ComparisonOperator.LESSTHANTHRESHOLD]: "services.cloudwatch.compLt",
  [ComparisonOperator.LESSTHANOREQUALTOTHRESHOLD]: "services.cloudwatch.compLte",
};

const getColumns = (t: TFunction): ColumnDef<MetricAlarm, any>[] => [
  { accessorKey: "alarmname", header: t("services.cloudwatch.alarmNameHeader"), cell: MonoCell },
  { accessorKey: "statevalue", header: t("services.cloudwatch.stateHeader"), cell: ({ getValue }) => { const v = getValue() as number; const key = STATE_VALUE_I18N[v]; return key ? <span className={`badge ${v === StateValue.ALARM ? "badge-red" : v === StateValue.OK ? "badge-green" : ""}`}>{t(key) ?? String(v)}</span> : String(v); }, size: 100 },
  { accessorKey: "metricname", header: t("services.cloudwatch.metricHeader"), cell: FallbackCell },
  { accessorKey: "namespace", header: t("services.cloudwatch.namespaceHeader"), cell: FallbackCell },
  { accessorKey: "statistic", header: t("services.cloudwatch.statisticHeader"), cell: ({ getValue }) => { const v = getValue() as number; return STATISTIC_I18N[v] ? t(STATISTIC_I18N[v]) ?? String(v) : String(v); }, size: 100 },
  { accessorKey: "period", header: t("services.cloudwatch.periodHeader"), size: 80 },
  { accessorKey: "evaluationperiods", header: t("services.cloudwatch.evalPeriodsHeader"), size: 80 },
  { accessorKey: "threshold", header: t("services.cloudwatch.thresholdHeader"), size: 80 },
  { accessorKey: "comparisonoperator", header: t("services.cloudwatch.comparisonHeader"), cell: ({ getValue }) => { const v = getValue() as number; return COMPARISON_I18N[v] ? t(COMPARISON_I18N[v]) ?? String(v) : String(v); }, size: 120 },
  { accessorKey: "alarmdescription", header: t("services.cloudwatch.descriptionHeader"), cell: FallbackCell },
];

type DetailTab = "detail" | "json";

export function CloudWatchPage() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(CloudWatchService);
  const { queryKey } = useListKey("cloudwatch");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<MetricAlarm | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formMetric, setFormMetric] = useState("");
  const [formNamespace, setFormNamespace] = useState("");
  const [formStat, setFormStat] = useState(Statistic.AVERAGE);
  const [formPeriod, setFormPeriod] = useState(60);
  const [formEvalPeriods, setFormEvalPeriods] = useState(1);
  const [formThreshold, setFormThreshold] = useState(0);
  const [formComp, setFormComp] = useState(ComparisonOperator.GREATERTHANTHRESHOLD);
  const [formDesc, setFormDesc] = useState("");

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.describeAlarms({}), refetchInterval: REFETCH_INTERVAL });
  const items: MetricAlarm[] = dropEmpty(data?.metricalarms ?? [], "alarmname");

  const createMutation = useMutation({
    mutationFn: () => client.putMetricAlarm(create(PutMetricAlarmInputSchema, { alarmname: formName, metricname: formMetric, namespace: formNamespace, statistic: formStat, period: formPeriod, evaluationperiods: formEvalPeriods, threshold: formThreshold, comparisonoperator: formComp, alarmdescription: formDesc })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormMetric(""); setFormNamespace(""); setFormDesc(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteAlarms({ alarmnames: [name] }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteAlarms({ alarmnames: [n] }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.alarmname) ? null : p)); },
  });

  const handleRowClick = (row: MetricAlarm) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.alarmname);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const stateKey = STATE_VALUE_I18N[selectedItem.statevalue];
    const statKey = STATISTIC_I18N[selectedItem.statistic];
    const compKey = COMPARISON_I18N[selectedItem.comparisonoperator];
    return (
      <DetailPanel title={selectedItem.alarmname} titleIcon="📊" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Alarm</td><td className="cell-mono">{selectedItem.alarmname}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>State</td><td><span className="badge">{stateKey ? t(stateKey) ?? String(selectedItem.statevalue) : String(selectedItem.statevalue)}</span></td></tr>
            <tr><td style={{ fontWeight: 600 }}>Metric</td><td>{selectedItem.metricname || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Namespace</td><td>{selectedItem.namespace || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Statistic</td><td>{statKey ? t(statKey) ?? String(selectedItem.statistic) : String(selectedItem.statistic)}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Period</td><td>{selectedItem.period}s</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Eval Periods</td><td>{selectedItem.evaluationperiods}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Threshold</td><td>{selectedItem.threshold}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Comparison</td><td>{compKey ? t(compKey) ?? String(selectedItem.comparisonoperator) : String(selectedItem.comparisonoperator)}</td></tr>
            {selectedItem.alarmdescription && <tr><td style={{ fontWeight: 600 }}>Description</td><td>{selectedItem.alarmdescription}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="📊" title={t("services.cloudwatch.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.cloudwatch.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.cloudwatch.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.cloudwatch.title") }, { label: t("services.cloudwatch.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-cw">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<MetricAlarm>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.alarmname), ...columns]} data={items} getRowId={(row) => row.alarmname} onRowClick={handleRowClick} selectedId={selectedItem?.alarmname} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.cloudwatch.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName || !formMetric}>
        <label>{t("services.cloudwatch.alarmNameLabel")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.cloudwatch.alarmNamePlaceholder")} className="modal-input" /></label>
        <label>{t("services.cloudwatch.metricNameLabel")}<input value={formMetric} onChange={(e) => setFormMetric(e.target.value)} className="modal-input" /></label>
        <label>{t("services.cloudwatch.namespaceLabel")}<input value={formNamespace} onChange={(e) => setFormNamespace(e.target.value)} className="modal-input" /></label>
        <label>{t("services.cloudwatch.statisticLabel")}<select value={formStat} onChange={(e) => setFormStat(Number(e.target.value))} className="modal-input">{Object.entries(STATISTIC_I18N).map(([v, k]) => <option key={v} value={v}>{t(k) ?? k}</option>)}</select></label>
        <label>{t("services.cloudwatch.periodLabel")}<input type="number" min={1} value={formPeriod} onChange={(e) => setFormPeriod(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.cloudwatch.evalPeriodsLabel")}<input type="number" min={1} value={formEvalPeriods} onChange={(e) => setFormEvalPeriods(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.cloudwatch.thresholdLabel")}<input type="number" step="any" value={formThreshold} onChange={(e) => setFormThreshold(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.cloudwatch.comparisonLabel")}<select value={formComp} onChange={(e) => setFormComp(Number(e.target.value))} className="modal-input">{Object.entries(COMPARISON_I18N).map(([v, k]) => <option key={v} value={v}>{t(k) ?? k}</option>)}</select></label>
        <label>{t("services.cloudwatch.descriptionLabel")}<input value={formDesc} onChange={(e) => setFormDesc(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.cloudwatch.delete")} name={selectedItem?.alarmname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.alarmname)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.cloudwatch.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
