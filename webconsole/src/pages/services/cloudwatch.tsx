/**
 * CloudWatch service page. Lists metric alarms with create/delete operations and
 * a custom detail panel showing metric configuration, dimensions, and actions.
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
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  FallbackCell,
  BooleanBadge,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";

/** Lookup map for StateValue proto enum values to i18n keys. */
const STATE_VALUE_I18N: Record<number, string> = {
  [StateValue.ALARM]: "services.cloudwatch.stateAlarm",
  [StateValue.OK]: "services.cloudwatch.stateOk",
  [StateValue.INSUFFICIENT_DATA]: "services.cloudwatch.stateInsufficientData",
};

/** Lookup map for Statistic proto enum values to i18n keys. */
const STATISTIC_I18N: Record<number, string> = {
  [Statistic.SUM]: "services.cloudwatch.statSum",
  [Statistic.SAMPLECOUNT]: "services.cloudwatch.statSampleCount",
  [Statistic.AVERAGE]: "services.cloudwatch.statAverage",
  [Statistic.MAXIMUM]: "services.cloudwatch.statMaximum",
  [Statistic.MINIMUM]: "services.cloudwatch.statMinimum",
};

/** Lookup map for ComparisonOperator proto enum values to i18n keys. */
const COMPARISON_I18N: Record<number, string> = {
  [ComparisonOperator.GREATERTHANOREQUALTOTHRESHOLD]: "services.cloudwatch.compGte",
  [ComparisonOperator.GREATERTHANTHRESHOLD]: "services.cloudwatch.compGt",
  [ComparisonOperator.LESSTHANTHRESHOLD]: "services.cloudwatch.compLt",
  [ComparisonOperator.LESSTHANOREQUALTOTHRESHOLD]: "services.cloudwatch.compLte",
  [ComparisonOperator.LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD]: "services.cloudwatch.compLtOrGt",
  [ComparisonOperator.LESSTHANLOWERTHRESHOLD]: "services.cloudwatch.compLtLower",
  [ComparisonOperator.GREATERTHANUPPERTHRESHOLD]: "services.cloudwatch.compGtUpper",
};

/** Shared state badge renderer used in both column and detail panel. */
const renderStateBadge = (state: number, t: TFunction) => {
  const cls = state === StateValue.OK ? "badge-green" : state === StateValue.ALARM ? "badge-red" : "badge-yellow";
  return <span className={`badge ${cls}`}>{STATE_VALUE_I18N[state] ? t(STATE_VALUE_I18N[state]) : String(state)}</span>;
};

/** Column definitions for the CloudWatch alarm table. */
const getColumns = (t: TFunction): ColumnDef<MetricAlarm, any>[] => [
  { accessorKey: "alarmname", header: t("services.cloudwatch.alarmNameHeader"), cell: MonoCell },
  {
    accessorKey: "statevalue",
    header: t("services.cloudwatch.stateHeader"),
    cell: ({ getValue }) => renderStateBadge(getValue() as StateValue, t),
    size: 90,
  },
  { accessorKey: "namespace", header: t("services.cloudwatch.namespaceHeader"), cell: FallbackCell },
  { accessorKey: "metricname", header: t("services.cloudwatch.metricHeader"), cell: FallbackCell },
  {
    accessorKey: "statistic",
    header: t("services.cloudwatch.statisticHeader"),
    cell: ({ getValue }) => { const k = STATISTIC_I18N[getValue() as number]; return k ? t(k) : "\u2014"; },
    size: 80,
  },
  { accessorKey: "threshold", header: t("services.cloudwatch.thresholdHeader"), size: 80 },
  { accessorKey: "period", header: t("services.cloudwatch.periodHeader"), cell: ({ getValue }) => { const v = getValue() as number; return v ? `${v}s` : "\u2014"; }, size: 70 },
  { accessorKey: "evaluationperiods", header: t("services.cloudwatch.evaluationPeriodsHeader"), size: 70 },
];

/** Detail panel for a CloudWatch metric alarm. */
function CloudWatchDetail({ item }: { item: MetricAlarm }) {
  const { t } = useTranslation();

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.nameLabel")}</span><span className="cell-mono">{item.alarmname}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.alarmarn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.stateLabel")}</span>{renderStateBadge(item.statevalue, t)}</div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.descriptionLabel")}</span><span>{item.alarmdescription || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.actionsEnabledLabel")}</span><BooleanBadge value={!!item.actionsenabled} /></div>
      </section>

      <section className="detail-section">
        <h3>{t("services.cloudwatch.detail.metricSection")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.namespaceLabel")}</span><span className="cell-mono">{item.namespace || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.metricLabel")}</span><span className="cell-mono">{item.metricname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.statisticLabel")}</span><span>{(() => { const k = STATISTIC_I18N[item.statistic]; return k ? t(k) : String(item.statistic); })()}</span></div>
        {item.extendedstatistic && <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.extendedStatLabel")}</span><span>{item.extendedstatistic}</span></div>}
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.periodLabel")}</span><span>{item.period ? `${item.period}s` : "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.evaluationPeriodsLabel")}</span><span>{item.evaluationperiods || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.datapointsToAlarmLabel")}</span><span>{item.datapointstoalarm || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.comparisonLabel")}</span><span>{(() => { const k = COMPARISON_I18N[item.comparisonoperator]; return k ? t(k) : String(item.comparisonoperator); })()}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.thresholdLabel")}</span><span>{item.threshold}</span></div>
        {item.treatmissingdata && <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.treatMissingLabel")}</span><span>{item.treatmissingdata}</span></div>}
      </section>

      {item.dimensions?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudwatch.detail.dimensionsSection")} ({item.dimensions.length})</h3>
          {item.dimensions.map((d, i) => (
            <div key={i} className="detail-field">
              <span className="detail-label">{(d as { name?: string; value?: string }).name || `dim ${i}`}</span>
              <span className="cell-mono">{(d as { name?: string; value?: string }).value || "\u2014"}</span>
            </div>
          ))}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("services.cloudwatch.detail.stateSection")}</h3>
        {item.statereason && <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.reasonLabel")}</span><span>{item.statereason}</span></div>}
        {item.stateupdatedtimestamp && (
          <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.updatedLabel")}</span><span>{new Date(item.stateupdatedtimestamp).toLocaleString()}</span></div>
        )}
        {item.statetransitionedtimestamp && (
          <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.transitionedLabel")}</span><span>{new Date(item.statetransitionedtimestamp).toLocaleString()}</span></div>
        )}
      </section>

      {item.alarmconfigurationupdatedtimestamp && (
        <section className="detail-section">
          <h3>{t("services.cloudwatch.detail.configurationSection")}</h3>
          <div className="detail-field"><span className="detail-label">{t("services.cloudwatch.detail.lastUpdatedLabel")}</span><span>{new Date(item.alarmconfigurationupdatedtimestamp).toLocaleString()}</span></div>
        </section>
      )}

      {item.alarmactions?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudwatch.detail.alarmActionsSection")} ({item.alarmactions.length})</h3>
          {item.alarmactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      {item.okactions?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudwatch.detail.okActionsSection")} ({item.okactions.length})</h3>
          {item.okactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      {item.insufficientdataactions?.length > 0 && (
        <section className="detail-section">
          <h3>{t("services.cloudwatch.detail.insufficientDataActionsSection")} ({item.insufficientdataactions.length})</h3>
          {item.insufficientdataactions.map((a, i) => (
            <div key={i} className="detail-field"><span className="detail-label">{i + 1}</span><span className="cell-mono" style={{ fontSize: 11 }}>{a}</span></div>
          ))}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** CloudWatch service page with alarm list, create, delete, and detail operations. */
export function CloudWatchPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<MetricAlarm | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formNamespace, setFormNamespace] = useState("AWS/Test");
  const [formMetricName, setFormMetricName] = useState("TestMetric");
  const [formStatistic, setFormStatistic] = useState<Statistic>(Statistic.AVERAGE);
  const [formPeriod, setFormPeriod] = useState(60);
  const [formEvalPeriods, setFormEvalPeriods] = useState(1);
  const [formThreshold, setFormThreshold] = useState(1);
  const [formCompOperator, setFormCompOperator] = useState<ComparisonOperator>(ComparisonOperator.GREATERTHANOREQUALTOTHRESHOLD);

  const { client, invalidate } = useServiceClient(CloudWatchService);
  const { queryKey } = useListKey("cloudwatch");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.describeAlarms({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: MetricAlarm[] = dropEmpty(data?.metricalarms ?? [], "alarmname");

  const createMutation = useMutation({
    mutationFn: () =>
      client.putMetricAlarm(
        create(PutMetricAlarmInputSchema, {
          alarmname: formName,
          namespace: formNamespace,
          metricname: formMetricName,
          statistic: formStatistic,
          period: formPeriod,
          evaluationperiods: formEvalPeriods,
          threshold: formThreshold,
          comparisonoperator: formCompOperator,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormNamespace("AWS/Test");
      setFormMetricName("TestMetric");
      setFormStatistic(Statistic.AVERAGE);
      setFormPeriod(60);
      setFormEvalPeriods(1);
      setFormThreshold(1);
      setFormCompOperator(ComparisonOperator.GREATERTHANOREQUALTOTHRESHOLD);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (alarmName: string) =>
      client.deleteAlarms({ alarmnames: [alarmName] }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="📊"
      title={t("services.cloudwatch.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cloudwatch.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cloudwatch.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.alarmname}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.alarmname}
        selected={selectedItem}
        detailTitle={selectedItem?.alarmname}
        onDetailClose={() => setSelectedItem(null)}
        DetailComponent={CloudWatchDetail}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.cloudwatch.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formNamespace || !formMetricName}
      >
        <label>
          {t("services.cloudwatch.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.cloudwatch.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.namespaceLabel")}
          <input
            value={formNamespace}
            onChange={(e) => setFormNamespace(e.target.value)}
            placeholder={t("services.cloudwatch.namespacePlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.metricNameLabel")}
          <input
            value={formMetricName}
            onChange={(e) => setFormMetricName(e.target.value)}
            placeholder={t("services.cloudwatch.metricNamePlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.statisticLabel")}
          <select value={formStatistic} onChange={(e) => setFormStatistic(Number(e.target.value) as Statistic)} className="modal-input">
            <option value={Statistic.AVERAGE}>{t("services.cloudwatch.statAverage")}</option>
            <option value={Statistic.SUM}>{t("services.cloudwatch.statSum")}</option>
            <option value={Statistic.MAXIMUM}>{t("services.cloudwatch.statMaximum")}</option>
            <option value={Statistic.MINIMUM}>{t("services.cloudwatch.statMinimum")}</option>
            <option value={Statistic.SAMPLECOUNT}>{t("services.cloudwatch.statSampleCount")}</option>
          </select>
        </label>
        <label>
          {t("services.cloudwatch.periodLabel")}
          <input
            type="number"
            value={formPeriod}
            onChange={(e) => setFormPeriod(Number(e.target.value))}
            min={1}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.evaluationPeriodsLabel")}
          <input
            type="number"
            value={formEvalPeriods}
            onChange={(e) => setFormEvalPeriods(Number(e.target.value))}
            min={1}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.thresholdLabel")}
          <input
            type="number"
            value={formThreshold}
            onChange={(e) => setFormThreshold(Number(e.target.value))}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatch.comparisonOperatorLabel")}
          <select value={formCompOperator} onChange={(e) => setFormCompOperator(Number(e.target.value) as ComparisonOperator)} className="modal-input">
            <option value={ComparisonOperator.GREATERTHANOREQUALTOTHRESHOLD}>{t("services.cloudwatch.compGte")}</option>
            <option value={ComparisonOperator.GREATERTHANTHRESHOLD}>{t("services.cloudwatch.compGt")}</option>
            <option value={ComparisonOperator.LESSTHANTHRESHOLD}>{t("services.cloudwatch.compLt")}</option>
            <option value={ComparisonOperator.LESSTHANOREQUALTOTHRESHOLD}>{t("services.cloudwatch.compLte")}</option>
            <option value={ComparisonOperator.LESSTHANLOWERORGREATERTHANUPPERTHRESHOLD}>{t("services.cloudwatch.compLtOrGt")}</option>
            <option value={ComparisonOperator.LESSTHANLOWERTHRESHOLD}>{t("services.cloudwatch.compLtLower")}</option>
            <option value={ComparisonOperator.GREATERTHANUPPERTHRESHOLD}>{t("services.cloudwatch.compGtUpper")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cloudwatch.delete")}
        name={selectedItem?.alarmname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.alarmname)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
