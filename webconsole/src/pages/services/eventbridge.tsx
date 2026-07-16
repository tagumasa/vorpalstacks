/**
 * EventBridge service page — 3-panel inspector layout with tabs.
 * Tabbed view for event buses, rules, and schedules.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { CloudWatchEventsService, type EventBus, type Rule, RuleState } from "@/gen/cloudwatchevents_pb";
import { CreateEventBusRequestSchema } from "@/gen/cloudwatchevents_pb";
import { SchedulerService, type ScheduleSummary } from "@/gen/scheduler_pb";
import { CreateScheduleInputSchema, FlexibleTimeWindowSchema, TargetSchema } from "@/gen/scheduler_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const getBusColumns = (t: TFunction): ColumnDef<EventBus, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.busNameHeader"), cell: MonoCell },
  { accessorKey: "arn", header: t("services.eventbridge.arnHeader"), cell: SmallMonoCell },
];

const getRuleColumns = (t: TFunction): ColumnDef<Rule, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.ruleNameHeader"), cell: MonoCell },
  { accessorKey: "eventbusname", header: t("services.eventbridge.eventBusHeader"), cell: MonoCell },
  { accessorKey: "state", header: t("services.eventbridge.stateHeader"), cell: ({ getValue }) => { const s = getValue() as RuleState; return s === RuleState.ENABLED ? <span className="badge badge-green">{t("services.eventbridge.stateEnabled")}</span> : s === RuleState.DISABLED ? <span className="badge badge-red">{t("services.eventbridge.stateDisabled")}</span> : <span className="badge">{String(s)}</span>; }, size: 80 },
  { accessorKey: "scheduleexpression", header: t("services.eventbridge.scheduleExpressionHeader"), cell: ({ getValue }) => (getValue() as string) || "\u2014" },
  { accessorKey: "description", header: t("services.eventbridge.ruleDescriptionHeader"), cell: ({ getValue }) => (getValue() as string) || "\u2014" },
];

const getScheduleColumns = (t: TFunction): ColumnDef<ScheduleSummary, any>[] => [
  { accessorKey: "name", header: t("services.eventbridge.scheduleNameHeader"), cell: MonoCell },
  { accessorKey: "groupname", header: t("services.eventbridge.groupHeader"), cell: ({ getValue }) => (getValue() as string) || t("services.eventbridge.defaultGroup") },
  { accessorKey: "state", header: t("services.eventbridge.stateHeader"), cell: ({ getValue }) => { const s = getValue() as string; return s ? <span className={`badge ${s === "ENABLED" ? "badge-green" : "badge-red"}`}>{s}</span> : "\u2014"; }, size: 90 },
  { accessorKey: "target", header: t("services.eventbridge.targetHeader"), cell: ({ getValue }) => { const tgt = getValue(); return tgt ? String(tgt) : "\u2014"; }, size: 120 },
  { accessorKey: "creationdate", header: t("services.eventbridge.createdHeader"), cell: DateCell },
  { accessorKey: "lastmodificationdate", header: t("services.eventbridge.lastModifiedHeader"), cell: DateCell },
];

type TabKey = "buses" | "rules" | "schedules";
type DetailTab = "detail" | "json";

export function EventBridgePage() {
  const { t, i18n } = useTranslation();
  const busColumns = getBusColumns(t);
  const ruleColumns = getRuleColumns(t);
  const scheduleColumns = getScheduleColumns(t);

  const [tab, setTab] = useState<TabKey>("buses");

  /* Selection state — one per tab */
  const busSelection = useSelection<string>();
  const scheduleSelection = useSelection<string>();

  /* Detail state — one selected item per tab */
  const [selectedBus, setSelectedBus] = useState<EventBus | null>(null);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
  const [selectedSchedule, setSelectedSchedule] = useState<ScheduleSummary | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  /* Modal state */
  const [showCreateBus, setShowCreateBus] = useState(false);
  const [showDeleteBus, setShowDeleteBus] = useState(false);
  const [showBatchDeleteBus, setShowBatchDeleteBus] = useState(false);
  const [showCreateSchedule, setShowCreateSchedule] = useState(false);
  const [showDeleteSchedule, setShowDeleteSchedule] = useState(false);
  const [showBatchDeleteSchedule, setShowBatchDeleteSchedule] = useState(false);

  /* Form state */
  const [formBusName, setFormBusName] = useState("");
  const [formEventSource, setFormEventSource] = useState("");
  const [formScheduleName, setFormScheduleName] = useState("");
  const [formScheduleExpression, setFormScheduleExpression] = useState("rate(5 minutes)");
  const [formScheduleGroup, setFormScheduleGroup] = useState("");
  const [formScheduleDesc, setFormScheduleDesc] = useState("");
  const [formTargetArn, setFormTargetArn] = useState("");
  const [formRoleArn, setFormRoleArn] = useState("");

  const { client } = useServiceClient(CloudWatchEventsService);
  const { client: schedulerClient } = useServiceClient(SchedulerService);
  const { queryKey: busesKey } = useListKey("eventbridge-buses");
  const { queryKey: rulesKey } = useListKey("eventbridge-rules");
  const { queryKey: schedulesKey } = useListKey("eventbridge-schedules");

  const busesList = usePaginatedList<EventBus, Awaited<ReturnType<typeof client.listEventBuses>>>({
    queryKeyBase: busesKey, fetchPage: (token) => client.listEventBuses({ nexttoken: token || undefined }), getItems: (r) => r.eventbuses ?? [], getNextToken: (r) => r.nexttoken ?? "",
  });
  const rulesList = usePaginatedList<Rule, Awaited<ReturnType<typeof client.listRules>>>({
    queryKeyBase: rulesKey, fetchPage: (token) => client.listRules({ nexttoken: token || undefined }), getItems: (r) => r.rules ?? [], getNextToken: (r) => r.nexttoken ?? "",
  });
  const schedulesList = usePaginatedList<ScheduleSummary, Awaited<ReturnType<typeof schedulerClient.listSchedules>>>({
    queryKeyBase: schedulesKey, fetchPage: (token) => schedulerClient.listSchedules({ nexttoken: token || undefined }), getItems: (r) => r.schedules ?? [], getNextToken: (r) => r.nexttoken ?? "",
  });

  const buses = dropEmpty(busesList.items, "name");
  const rules = dropEmpty(rulesList.items, "name");
  const schedules = dropEmpty(schedulesList.items, "name");

  const query = tab === "buses" ? busesList : tab === "rules" ? rulesList : schedulesList;

  const createBusMutation = useMutation({
    mutationFn: () => client.createEventBus(create(CreateEventBusRequestSchema, { name: formBusName, ...(formEventSource ? { eventsourcename: formEventSource } : {}) })),
    onSuccess: () => { busesList.invalidate(); setShowCreateBus(false); setFormBusName(""); setFormEventSource(""); },
  });
  const deleteBusMutation = useMutation({
    mutationFn: (name: string) => client.deleteEventBus({ name }),
    onSuccess: () => { busesList.invalidate(); setShowDeleteBus(false); setSelectedBus(null); busSelection.clear(); },
  });
  const batchDeleteBusMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteEventBus({ name: n }))),
    onSuccess: (_d, names) => { busesList.invalidate(); setShowBatchDeleteBus(false); busSelection.clear(); setSelectedBus((p) => (p && names.includes(p.name) ? null : p)); },
  });

  const createScheduleMutation = useMutation({
    mutationFn: () => schedulerClient.createSchedule(create(CreateScheduleInputSchema, { name: formScheduleName, scheduleexpression: formScheduleExpression, groupname: formScheduleGroup, description: formScheduleDesc, flexibletimewindow: create(FlexibleTimeWindowSchema, { mode: "OFF" }), target: create(TargetSchema, { arn: formTargetArn, rolearn: formRoleArn }) })),
    onSuccess: () => { schedulesList.invalidate(); setShowCreateSchedule(false); setFormScheduleName(""); setFormScheduleExpression("rate(5 minutes)"); setFormScheduleGroup(""); setFormScheduleDesc(""); setFormTargetArn(""); setFormRoleArn(""); },
  });
  const deleteScheduleMutation = useMutation({
    mutationFn: (item: ScheduleSummary) => schedulerClient.deleteSchedule({ name: item.name, groupname: item.groupname }),
    onSuccess: () => { schedulesList.invalidate(); setShowDeleteSchedule(false); setSelectedSchedule(null); scheduleSelection.clear(); },
  });
  const batchDeleteScheduleMutation = useMutation({
    mutationFn: async (items: ScheduleSummary[]) => Promise.allSettled(items.map((s) => schedulerClient.deleteSchedule({ name: s.name, groupname: s.groupname }))),
    onSuccess: (_d, deletedItems) => { schedulesList.invalidate(); setShowBatchDeleteSchedule(false); scheduleSelection.clear(); const deletedNames = new Set(deletedItems.map((s) => s.name)); setSelectedSchedule((p) => (p && deletedNames.has(p.name) ? null : p)); },
  });

  const tabs = [
    { key: "buses" as TabKey, label: t("services.eventbridge.tabs.buses"), count: buses.length },
    { key: "rules" as TabKey, label: t("services.eventbridge.tabs.rules"), count: rules.length },
    { key: "schedules" as TabKey, label: t("services.eventbridge.tabs.schedules"), count: schedules.length },
  ];

  const handleTabChange = (k: string) => { setTab(k as TabKey); setSelectedBus(null); setSelectedRule(null); setSelectedSchedule(null); setDetailTab("detail"); busSelection.clear(); scheduleSelection.clear(); };

  /* Detail panel renderers */
  const renderBusDetail = () => {
    if (!selectedBus) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedBus.name} titleIcon="📡" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteBus(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.name")}</td><td className="cell-mono">{selectedBus.name}</td></tr>
            <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedBus.arn || "\u2014"}</td></tr>
            {selectedBus.policy && <tr><td className="detail-label">{t("common.fields.policy")}</td><td><JsonViewer data={(() => { try { return JSON.parse(selectedBus.policy); } catch { return selectedBus.policy; } })()} /></td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedBus} />}
      </DetailPanel>
    );
  };

  const renderRuleDetail = () => {
    if (!selectedRule) return <DetailEmpty message={t("common.noItemSelected")} />;
    let parsedPattern: Record<string, unknown> | null = null;
    if (selectedRule.eventpattern) { try { parsedPattern = JSON.parse(selectedRule.eventpattern); } catch { /* not JSON */ } }
    return (
      <DetailPanel title={selectedRule.name} titleIcon="📋" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.name")}</td><td className="cell-mono">{selectedRule.name}</td></tr>
            <tr><td className="detail-label">{t("common.fields.eventBus")}</td><td className="cell-mono">{selectedRule.eventbusname || "\u2014"}</td></tr>
            <tr><td className="detail-label">{t("common.fields.state")}</td><td>{selectedRule.state === RuleState.ENABLED ? <span className="badge badge-green">{t("services.eventbridge.stateEnabled")}</span> : <span className="badge badge-red">{t("services.eventbridge.stateDisabled")}</span>}</td></tr>
            <tr><td className="detail-label">{t("common.fields.description")}</td><td>{selectedRule.description || "\u2014"}</td></tr>
            {selectedRule.scheduleexpression && <tr><td className="detail-label">{t("common.fields.schedule")}</td><td className="cell-mono">{selectedRule.scheduleexpression}</td></tr>}
            {selectedRule.rolearn && <tr><td className="detail-label">{t("common.fields.roleArn")}</td><td className="cell-mono cell-long">{selectedRule.rolearn}</td></tr>}
            {selectedRule.eventpattern && <tr><td className="detail-label">{t("common.fields.eventPattern")}</td><td>{parsedPattern ? <JsonViewer data={parsedPattern} /> : <pre className="code-block" style={{ margin: 0 }}>{selectedRule.eventpattern}</pre>}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedRule} />}
      </DetailPanel>
    );
  };

  const renderScheduleDetail = () => {
    if (!selectedSchedule) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedSchedule.name} titleIcon="⏰" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteSchedule(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.name")}</td><td className="cell-mono">{selectedSchedule.name}</td></tr>
            <tr><td className="detail-label">{t("common.fields.group")}</td><td>{selectedSchedule.groupname || t("services.eventbridge.defaultGroup")}</td></tr>
            <tr><td className="detail-label">{t("common.fields.state")}</td><td><span className="badge">{String(selectedSchedule.state || "\u2014")}</span></td></tr>
            <tr><td className="detail-label">{t("common.fields.target")}</td><td>{selectedSchedule.target ? String(selectedSchedule.target) : "\u2014"}</td></tr>
            {selectedSchedule.creationdate && <tr><td className="detail-label">{t("common.fields.created")}</td><td>{fmtDate(selectedSchedule.creationdate, i18n.language)}</td></tr>}
            {selectedSchedule.lastmodificationdate && <tr><td className="detail-label">{t("common.fields.modified")}</td><td>{fmtDate(selectedSchedule.lastmodificationdate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedSchedule} />}
      </DetailPanel>
    );
  };

  const selectedScheduleItems = schedules.filter((s) => scheduleSelection.selected.has(s.name));

  return (
    <ServicePageLayout icon="📡" title={t("services.eventbridge.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange} actions={
      tab === "buses" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateBus(true)}>{t("services.eventbridge.create")}</button>
        <button className="btn btn-danger" disabled={busSelection.selected.size === 0} onClick={() => setShowBatchDeleteBus(true)}>{t("common.deleteSelected")}{busSelection.selected.size > 0 && <span className="batch-count">({busSelection.selected.size})</span>}</button>
      </>) : tab === "schedules" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateSchedule(true)}>{t("services.scheduler.create")}</button>
        <button className="btn btn-danger" disabled={scheduleSelection.selected.size === 0} onClick={() => setShowBatchDeleteSchedule(true)}>{t("common.deleteSelected")}{scheduleSelection.selected.size > 0 && <span className="batch-count">({scheduleSelection.selected.size})</span>}</button>
      </>) : undefined
    }>
      {/* Buses tab */}
      {tab === "buses" && (
        buses.length > 0 ? (
          <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-eb-buses">
            <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<EventBus>(busSelection.selected, busSelection.toggle, () => busSelection.toggleAll(buses.map((b) => b.name)), buses.map((b) => b.name), t, (row) => row.name), ...busColumns]} data={buses} getRowId={(row) => row.name} onRowClick={(row) => { setSelectedBus(row); setDetailTab("detail"); }} selectedId={selectedBus?.name} hasMore={busesList.hasMore} onLoadMore={busesList.loadMore} loadingMore={busesList.isFetchingMore} /></div>
            {renderBusDetail()}
          </Splitter>
        ) : <div className="empty-state">{t("common.noData")}</div>
      )}

      {/* Rules tab */}
      {tab === "rules" && (
        rules.length > 0 ? (
          <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-eb-rules">
            <div className="flex-fill-scroll"><DataTable columns={ruleColumns} data={rules} getRowId={(row) => row.name} onRowClick={(row) => { setSelectedRule(row); setDetailTab("detail"); }} selectedId={selectedRule?.name} hasMore={rulesList.hasMore} onLoadMore={rulesList.loadMore} loadingMore={rulesList.isFetchingMore} /></div>
            {renderRuleDetail()}
          </Splitter>
        ) : <div className="empty-state">{t("common.noData")}</div>
      )}

      {/* Schedules tab */}
      {tab === "schedules" && (
        schedules.length > 0 ? (
          <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-eb-schedules">
            <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<ScheduleSummary>(scheduleSelection.selected, scheduleSelection.toggle, () => scheduleSelection.toggleAll(schedules.map((s) => s.name)), schedules.map((s) => s.name), t, (row) => row.name), ...scheduleColumns]} data={schedules} getRowId={(row) => row.name} onRowClick={(row) => { setSelectedSchedule(row); setDetailTab("detail"); }} selectedId={selectedSchedule?.name} hasMore={schedulesList.hasMore} onLoadMore={schedulesList.loadMore} loadingMore={schedulesList.isFetchingMore} /></div>
            {renderScheduleDetail()}
          </Splitter>
        ) : <div className="empty-state">{t("common.noData")}</div>
      )}

      {/* Bus modals */}
      <ServiceCreateModal open={showCreateBus} onClose={() => setShowCreateBus(false)} title={t("services.eventbridge.create")} error={createBusMutation.error} isPending={createBusMutation.isPending} onCreate={() => createBusMutation.mutate()} disabled={!formBusName}>
        <label>{t("services.eventbridge.nameField")}<input value={formBusName} onChange={(e) => setFormBusName(e.target.value)} placeholder={t("services.eventbridge.placeholder")} className="modal-input" /></label>
        <label>{t("services.eventbridge.eventSourceLabel")}<input value={formEventSource} onChange={(e) => setFormEventSource(e.target.value)} placeholder={t("services.eventbridge.eventSourcePlaceholder")} className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDeleteBus && !!selectedBus} title={t("services.eventbridge.delete")} name={selectedBus?.name} error={deleteBusMutation.error} isPending={deleteBusMutation.isPending} onConfirm={() => selectedBus && deleteBusMutation.mutate(selectedBus.name)} onClose={() => setShowDeleteBus(false)} />
      <ServiceDeleteDialog open={showBatchDeleteBus} title={t("common.deleteSelected")} name={`${busSelection.selected.size} ${t("services.eventbridge.tabs.buses")}`} error={batchDeleteBusMutation.error} isPending={batchDeleteBusMutation.isPending} onConfirm={() => batchDeleteBusMutation.mutate(Array.from(busSelection.selected))} onClose={() => setShowBatchDeleteBus(false)} />

      {/* Schedule modals */}
      <ServiceCreateModal open={showCreateSchedule} onClose={() => setShowCreateSchedule(false)} title={t("services.scheduler.create")} error={createScheduleMutation.error} isPending={createScheduleMutation.isPending} onCreate={() => createScheduleMutation.mutate()} disabled={!formScheduleName || !formScheduleExpression || !formTargetArn || !formRoleArn}>
        <label>{t("services.scheduler.nameField")}<input value={formScheduleName} onChange={(e) => setFormScheduleName(e.target.value)} placeholder={t("services.scheduler.placeholder")} className="modal-input" /></label>
        <label>{t("services.scheduler.scheduleLabel")}<input value={formScheduleExpression} onChange={(e) => setFormScheduleExpression(e.target.value)} placeholder={t("services.scheduler.schedulePlaceholder")} className="modal-input" /></label>
        <label>{t("services.scheduler.groupLabel")}<input value={formScheduleGroup} onChange={(e) => setFormScheduleGroup(e.target.value)} placeholder={t("services.scheduler.groupPlaceholder")} className="modal-input" /></label>
        <label>{t("services.scheduler.descriptionLabel")}<input value={formScheduleDesc} onChange={(e) => setFormScheduleDesc(e.target.value)} placeholder={t("services.scheduler.descriptionPlaceholder")} className="modal-input" /></label>
        <label>{t("services.scheduler.targetLabel")}<input value={formTargetArn} onChange={(e) => setFormTargetArn(e.target.value)} placeholder={t("services.eventbridge.targetArnPlaceholder")} className="modal-input" /></label>
        <label>{t("services.sfn.roleArnLabel")}<input value={formRoleArn} onChange={(e) => setFormRoleArn(e.target.value)} placeholder={t("services.eventbridge.roleArnPlaceholder")} className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDeleteSchedule && !!selectedSchedule} title={t("services.scheduler.delete")} name={selectedSchedule?.name} error={deleteScheduleMutation.error} isPending={deleteScheduleMutation.isPending} onConfirm={() => selectedSchedule && deleteScheduleMutation.mutate(selectedSchedule)} onClose={() => setShowDeleteSchedule(false)} />
      <ServiceDeleteDialog open={showBatchDeleteSchedule} title={t("common.deleteSelected")} name={`${scheduleSelection.selected.size} ${t("services.eventbridge.tabs.schedules")}`} error={batchDeleteScheduleMutation.error} isPending={batchDeleteScheduleMutation.isPending} onConfirm={() => batchDeleteScheduleMutation.mutate(selectedScheduleItems)} onClose={() => setShowBatchDeleteSchedule(false)} />
    </ServicePageLayout>
  );
}
