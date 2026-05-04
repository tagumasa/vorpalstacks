/**
 * EventBridge service page. Tabbed view for event buses and rules with
 * create/delete operations for buses and custom detail panels.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudWatchEventsService, type EventBus, type Rule, RuleState } from "@/gen/cloudwatchevents_pb";
import { CreateEventBusRequestSchema } from "@/gen/cloudwatchevents_pb";
import { SchedulerService, type ScheduleSummary } from "@/gen/scheduler_pb";
import { CreateScheduleInputSchema, FlexibleTimeWindowSchema, TargetSchema } from "@/gen/scheduler_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";

/** Column definitions for the EventBridge event bus table. */
const getBusColumns = (t: TFunction): ColumnDef<EventBus, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.busNameHeader"), cell: MonoCell },
  { accessorKey: "arn", header: t("services.eventbridge.arnHeader"), cell: SmallMonoCell },
];

/** Column definitions for the EventBridge rule table. */
const getRuleColumns = (t: TFunction): ColumnDef<Rule, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.ruleNameHeader"), cell: MonoCell },
  { accessorKey: "eventbusname", header: t("services.eventbridge.eventBusHeader"), cell: MonoCell },
  {
    accessorKey: "state",
    header: t("services.eventbridge.stateHeader"),
    cell: ({ getValue }) => {
      const s = getValue() as RuleState;
      return s === RuleState.ENABLED
        ? <span className="badge badge-green">{t("services.eventbridge.stateEnabled")}</span>
        : s === RuleState.DISABLED
          ? <span className="badge badge-red">{t("services.eventbridge.stateDisabled")}</span>
          : <span className="badge">{String(s)}</span>;
    },
    size: 80,
  },
  { accessorKey: "scheduleexpression", header: t("services.eventbridge.scheduleExpressionHeader"), cell: ({ getValue }) => (getValue() as string) || "\u2014" },
  { accessorKey: "description", header: t("services.eventbridge.ruleDescriptionHeader"), cell: ({ getValue }) => (getValue() as string) || "\u2014" },
];

/** Column definitions for the EventBridge Scheduler schedule table. */
const getScheduleColumns = (t: TFunction): ColumnDef<ScheduleSummary, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.scheduleNameHeader"), cell: MonoCell },
  { accessorKey: "groupname", header: t("services.eventbridge.groupHeader"), cell: ({ getValue }) => (getValue() as string) || t("services.eventbridge.defaultGroup") },
  { accessorKey: "state", header: t("services.eventbridge.stateHeader"), cell: ({ getValue }) => {
    const s = getValue() as string;
    return s ? <span className={`badge ${s === "ENABLED" ? "badge-green" : "badge-red"}`}>{s}</span> : "\u2014";
  }, size: 90 },
  { accessorKey: "target", header: t("services.eventbridge.targetHeader"), cell: ({ getValue }) => { const tgt = getValue(); return tgt ? String(tgt) : "\u2014"; }, size: 120 },
  { accessorKey: "creationdate", header: t("services.eventbridge.createdHeader"), cell: DateCell },
  { accessorKey: "lastmodificationdate", header: t("services.eventbridge.lastModifiedHeader"), cell: DateCell },
];

/** Detail panel for an EventBridge event bus. */
function EventBusDetail({ item }: { item: EventBus }) {
  const { t } = useTranslation();
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.nameLabel")}</span><span className="cell-mono">{item.name}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
      </section>
      {item.policy && (
        <section className="detail-section">
          <h3>{t("services.eventbridge.detail.resourcePolicySection")}</h3>
          <JsonViewer data={(() => { try { return JSON.parse(item.policy); } catch { return item.policy; } })()} />
        </section>
      )}
      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** Detail panel for an EventBridge rule. */
function RuleDetail({ item }: { item: Rule }) {
  const { t } = useTranslation();
  let parsedPattern: Record<string, unknown> | null = null;
  if (item.eventpattern) {
    try { parsedPattern = JSON.parse(item.eventpattern); } catch { /* not JSON */ }
  }

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.nameLabel")}</span><span className="cell-mono">{item.name}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.arnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.eventBusLabel")}</span><span className="cell-mono">{item.eventbusname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.stateLabel")}</span>
          {item.state === RuleState.ENABLED
            ? <span className="badge badge-green">{t("services.eventbridge.stateEnabled")}</span>
            : <span className="badge badge-red">{t("services.eventbridge.stateDisabled")}</span>}
        </div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.descriptionLabel")}</span><span>{item.description || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.managedByLabel")}</span><span>{item.managedby || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">{t("services.eventbridge.detail.roleArnLabel")}</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.rolearn || "\u2014"}</span></div>
      </section>

      {item.scheduleexpression && (
        <section className="detail-section">
          <h3>{t("services.eventbridge.detail.scheduleExpressionSection")}</h3>
          <pre className="code-block" style={{ margin: 0 }}>{item.scheduleexpression}</pre>
        </section>
      )}

      {item.eventpattern && (
        <section className="detail-section">
          <h3>{t("services.eventbridge.detail.eventPatternSection")}</h3>
          {parsedPattern
            ? <JsonViewer data={parsedPattern} />
            : <pre className="code-block" style={{ margin: 0 }}>{item.eventpattern}</pre>}
        </section>
      )}

      <section className="detail-section">
        <h3>{t("common.rawJson")}</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

/** Tab key type for the EventBridge page. */
type TabKey = "buses" | "rules" | "schedules";

/** EventBridge service page with tabbed buses/rules/schedules view and CRUD for buses and schedules. */
export function EventBridgePage() {
  const { t } = useTranslation();
  const busColumns = getBusColumns(t);
  const ruleColumns = getRuleColumns(t);
  const scheduleColumns = getScheduleColumns(t);
  const [tab, setTab] = useState<TabKey>("buses");
  const [selectedBus, setSelectedBus] = useState<EventBus | null>(null);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
  const [selectedSchedule, setSelectedSchedule] = useState<ScheduleSummary | null>(null);
  const [showCreateBus, setShowCreateBus] = useState(false);
  const [showDeleteBus, setShowDeleteBus] = useState(false);
  const [showCreateSchedule, setShowCreateSchedule] = useState(false);
  const [showDeleteSchedule, setShowDeleteSchedule] = useState(false);
  const [formBusName, setFormBusName] = useState("");
  const [formEventSource, setFormEventSource] = useState("");
  const [formScheduleName, setFormScheduleName] = useState("");
  const [formScheduleExpression, setFormScheduleExpression] = useState("rate(5 minutes)");
  const [formScheduleGroup, setFormScheduleGroup] = useState("");
  const [formScheduleDesc, setFormScheduleDesc] = useState("");
  const [formTargetArn, setFormTargetArn] = useState("");
  const [formRoleArn, setFormRoleArn] = useState("");

  const { client, invalidate } = useServiceClient(CloudWatchEventsService);
  const { client: schedulerClient } = useServiceClient(SchedulerService);
  const { queryKey: busesKey } = useListKey("eventbridge-buses");
  const { queryKey: rulesKey } = useListKey("eventbridge-rules");
  const { queryKey: schedulesKey } = useListKey("eventbridge-schedules");

  const busesQ = useQuery({
    queryKey: busesKey,
    queryFn: async () => {
      const resp = await client.listEventBuses({});
      return dropEmpty(resp.eventbuses ?? [], "name");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const rulesQ = useQuery({
    queryKey: rulesKey,
    queryFn: async () => {
      const resp = await client.listRules({});
      return dropEmpty(resp.rules ?? [], "name");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const schedulesQ = useQuery({
    queryKey: schedulesKey,
    queryFn: async () => {
      const resp = await schedulerClient.listSchedules({});
      return dropEmpty(resp.schedules ?? [], "name");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const query = tab === "buses" ? busesQ : tab === "rules" ? rulesQ : schedulesQ;
  const buses = busesQ.data ?? [];
  const rules = rulesQ.data ?? [];
  const schedules = schedulesQ.data ?? [];

  const createBusMutation = useMutation({
    mutationFn: () =>
      client.createEventBus(
        create(CreateEventBusRequestSchema, {
          name: formBusName,
          ...(formEventSource ? { eventsourcename: formEventSource } : {}),
        }),
      ),
    onSuccess: () => {
      invalidate(busesKey);
      setShowCreateBus(false);
      setFormBusName("");
      setFormEventSource("");
    },
  });

  const deleteBusMutation = useMutation({
    mutationFn: (busName: string) =>
      client.deleteEventBus({ name: busName }),
    onSuccess: () => {
      invalidate(busesKey);
      setShowDeleteBus(false);
      setSelectedBus(null);
    },
  });

  const createScheduleMutation = useMutation({
    mutationFn: () =>
      schedulerClient.createSchedule(
        create(CreateScheduleInputSchema, {
          name: formScheduleName,
          scheduleexpression: formScheduleExpression,
          groupname: formScheduleGroup,
          description: formScheduleDesc,
          flexibletimewindow: create(FlexibleTimeWindowSchema, { mode: "OFF" }),
          target: create(TargetSchema, { arn: formTargetArn, rolearn: formRoleArn }),
        }),
      ),
    onSuccess: () => {
      invalidate(schedulesKey);
      setShowCreateSchedule(false);
      setFormScheduleName("");
      setFormScheduleExpression("rate(5 minutes)");
      setFormScheduleGroup("");
      setFormScheduleDesc("");
      setFormTargetArn("");
      setFormRoleArn("");
    },
  });

  const deleteScheduleMutation = useMutation({
    mutationFn: (item: ScheduleSummary) =>
      schedulerClient.deleteSchedule({ name: item.name, groupname: item.groupname }),
    onSuccess: () => {
      invalidate(schedulesKey);
      setShowDeleteSchedule(false);
      setSelectedSchedule(null);
    },
  });

  const tabs = [
    { key: "buses" as TabKey, label: t("services.eventbridge.tabs.buses"), count: buses.length },
    { key: "rules" as TabKey, label: t("services.eventbridge.tabs.rules"), count: rules.length },
    { key: "schedules" as TabKey, label: t("services.eventbridge.tabs.schedules"), count: schedules.length },
  ];

  return (
    <ServicePageLayout
      icon="📡"
      title={t("services.eventbridge.title")}
      isLoading={query.isLoading}
      error={query.error}
      tabs={tabs}
      activeTab={tab}
      onTabChange={(k) => { setTab(k as TabKey); setSelectedBus(null); setSelectedRule(null); setSelectedSchedule(null); }}
      actions={
        tab === "buses" ? (
          <>
            <button className="btn btn-primary" onClick={() => setShowCreateBus(true)}>
              {t("services.eventbridge.create")}
            </button>
            {selectedBus && (
              <button className="btn btn-danger" onClick={() => setShowDeleteBus(true)}>
                {t("common.delete")}
              </button>
            )}
          </>
        ) : tab === "schedules" ? (
          <>
            <button className="btn btn-primary" onClick={() => setShowCreateSchedule(true)}>
              {t("services.scheduler.create")}
            </button>
            {selectedSchedule && (
              <button className="btn btn-danger" onClick={() => setShowDeleteSchedule(true)}>
                {t("common.delete")}
              </button>
            )}
          </>
        ) : undefined
      }
    >
      {tab === "buses" && (
        <SplitPane
          columns={busColumns}
          data={buses}
          getRowId={(row) => row.name}
          onRowClick={setSelectedBus}
          selectedId={selectedBus?.name}
          selected={selectedBus}
          detailTitle={selectedBus?.name}
          onDetailClose={() => setSelectedBus(null)}
          DetailComponent={EventBusDetail}
        />
      )}
      {tab === "rules" && (
        <SplitPane
          columns={ruleColumns}
          data={rules}
          getRowId={(row) => row.name}
          onRowClick={setSelectedRule}
          selectedId={selectedRule?.name}
          selected={selectedRule}
          detailTitle={selectedRule?.name}
          onDetailClose={() => setSelectedRule(null)}
          DetailComponent={RuleDetail}
        />
      )}
      {tab === "schedules" && (
        <SplitPane
          columns={scheduleColumns}
          data={schedules}
          getRowId={(row) => row.name}
          onRowClick={setSelectedSchedule}
          selectedId={selectedSchedule?.name}
          selected={selectedSchedule}
          detailTitle={selectedSchedule?.name}
          onDetailClose={() => setSelectedSchedule(null)}
        />
      )}

      <ServiceCreateModal
        open={showCreateBus}
        onClose={() => setShowCreateBus(false)}
        title={t("services.eventbridge.create")}
        error={createBusMutation.error}
        isPending={createBusMutation.isPending}
        onCreate={() => createBusMutation.mutate()}
        disabled={!formBusName}
      >
        <label>
          {t("services.eventbridge.nameField")}
          <input
            value={formBusName}
            onChange={(e) => setFormBusName(e.target.value)}
            placeholder={t("services.eventbridge.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.eventbridge.eventSourceLabel")}
          <input
            value={formEventSource}
            onChange={(e) => setFormEventSource(e.target.value)}
            placeholder={t("services.eventbridge.eventSourcePlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDeleteBus && !!selectedBus}
        title={t("services.eventbridge.delete")}
        name={selectedBus?.name}
        error={deleteBusMutation.error}
        isPending={deleteBusMutation.isPending}
        onConfirm={() => selectedBus && deleteBusMutation.mutate(selectedBus.name)}
        onClose={() => setShowDeleteBus(false)}
      />

      <ServiceCreateModal
        open={showCreateSchedule}
        onClose={() => setShowCreateSchedule(false)}
        title={t("services.scheduler.create")}
        error={createScheduleMutation.error}
        isPending={createScheduleMutation.isPending}
        onCreate={() => createScheduleMutation.mutate()}
        disabled={!formScheduleName || !formScheduleExpression || !formTargetArn || !formRoleArn}
      >
        <label>
          {t("services.scheduler.nameField")}
          <input
            value={formScheduleName}
            onChange={(e) => setFormScheduleName(e.target.value)}
            placeholder={t("services.scheduler.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.scheduler.scheduleLabel")}
          <input
            value={formScheduleExpression}
            onChange={(e) => setFormScheduleExpression(e.target.value)}
            placeholder={t("services.scheduler.schedulePlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.scheduler.groupLabel")}
          <input
            value={formScheduleGroup}
            onChange={(e) => setFormScheduleGroup(e.target.value)}
            placeholder={t("services.scheduler.groupPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.scheduler.descriptionLabel")}
          <input
            value={formScheduleDesc}
            onChange={(e) => setFormScheduleDesc(e.target.value)}
            placeholder={t("services.scheduler.descriptionPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.scheduler.targetLabel")}
          <input
            value={formTargetArn}
            onChange={(e) => setFormTargetArn(e.target.value)}
            placeholder={t("services.eventbridge.targetArnPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sfn.roleArnLabel")}
          <input
            value={formRoleArn}
            onChange={(e) => setFormRoleArn(e.target.value)}
            placeholder={t("services.eventbridge.roleArnPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDeleteSchedule && !!selectedSchedule}
        title={t("services.scheduler.delete")}
        name={selectedSchedule?.name}
        error={deleteScheduleMutation.error}
        isPending={deleteScheduleMutation.isPending}
        onConfirm={() => selectedSchedule && deleteScheduleMutation.mutate(selectedSchedule)}
        onClose={() => setShowDeleteSchedule(false)}
      />
    </ServicePageLayout>
  );
}
