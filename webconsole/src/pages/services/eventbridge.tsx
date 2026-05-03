import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { CloudWatchEventsService, type EventBus, type Rule, RuleState } from "@/gen/cloudwatchevents_pb";
import { CreateEventBusRequestSchema } from "@/gen/cloudwatchevents_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";

const busColumns: ColumnDef<EventBus, any>[] = [
  {
    accessorKey: "name",
    header: "Bus Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "arn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
];

const ruleColumns: ColumnDef<Rule, any>[] = [
  {
    accessorKey: "name",
    header: "Rule Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "eventbusname",
    header: "Event Bus",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "state",
    header: "State",
    cell: ({ getValue }) => {
      const s = getValue() as RuleState;
      return s === RuleState.ENABLED
        ? <span className="badge badge-green">Enabled</span>
        : s === RuleState.DISABLED
          ? <span className="badge badge-red">Disabled</span>
          : <span className="badge">{String(s)}</span>;
    },
    size: 80,
  },
];

function EventBusDetail({ item }: { item: EventBus }) {
  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Name</span><span className="cell-mono">{item.name}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
      </section>
      {item.policy && (
        <section className="detail-section">
          <h3>Resource Policy</h3>
          <JsonViewer data={(() => { try { return JSON.parse(item.policy); } catch { return item.policy; } })()} />
        </section>
      )}
      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

function RuleDetail({ item }: { item: Rule }) {
  let parsedPattern: any = null;
  if (item.eventpattern) {
    try { parsedPattern = JSON.parse(item.eventpattern); } catch { /* not JSON */ }
  }

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>General</h3>
        <div className="detail-field"><span className="detail-label">Name</span><span className="cell-mono">{item.name}</span></div>
        <div className="detail-field"><span className="detail-label">ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.arn || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Event Bus</span><span className="cell-mono">{item.eventbusname || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">State</span>
          {item.state === RuleState.ENABLED
            ? <span className="badge badge-green">Enabled</span>
            : <span className="badge badge-red">Disabled</span>}
        </div>
        <div className="detail-field"><span className="detail-label">Description</span><span>{item.description || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Managed By</span><span>{item.managedby || "\u2014"}</span></div>
        <div className="detail-field"><span className="detail-label">Role ARN</span><span className="cell-mono" style={{ fontSize: 11 }}>{item.rolearn || "\u2014"}</span></div>
      </section>

      {item.scheduleexpression && (
        <section className="detail-section">
          <h3>Schedule Expression</h3>
          <pre className="code-block" style={{ margin: 0 }}>{item.scheduleexpression}</pre>
        </section>
      )}

      {item.eventpattern && (
        <section className="detail-section">
          <h3>Event Pattern</h3>
          {parsedPattern
            ? <JsonViewer data={parsedPattern} />
            : <pre className="code-block" style={{ margin: 0 }}>{item.eventpattern}</pre>}
        </section>
      )}

      <section className="detail-section">
        <h3>Raw JSON</h3>
        <JsonViewer data={item} />
      </section>
    </div>
  );
}

type TabKey = "buses" | "rules";

export function EventBridgePage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("buses");
  const [selectedBus, setSelectedBus] = useState<EventBus | null>(null);
  const [selectedRule, setSelectedRule] = useState<Rule | null>(null);
  const [showCreateBus, setShowCreateBus] = useState(false);
  const [showDeleteBus, setShowDeleteBus] = useState(false);
  const [formBusName, setFormBusName] = useState("");

  const client = createClient(CloudWatchEventsService, transport);
  const { queryKey: busesKey } = useListKey("eventbridge-buses");
  const { queryKey: rulesKey } = useListKey("eventbridge-rules");
  const invalidate = useListInvalidator();

  const busesQ = useQuery({
    queryKey: busesKey,
    queryFn: async () => {
      const resp = await client.listEventBuses({});
      return dropEmpty(resp.eventbuses ?? [], "name");
    },
    refetchInterval: 30_000,
  });

  const rulesQ = useQuery({
    queryKey: rulesKey,
    queryFn: async () => {
      const resp = await client.listRules({});
      return dropEmpty(resp.rules ?? [], "name");
    },
    refetchInterval: 30_000,
  });

  const query = tab === "buses" ? busesQ : rulesQ;

  const createBusMutation = useMutation({
    mutationFn: () =>
      client.createEventBus(
        create(CreateEventBusRequestSchema, { name: formBusName }),
      ),
    onSuccess: () => {
      invalidate(busesKey);
      setShowCreateBus(false);
      setFormBusName("");
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

  if (query.isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📡</span>
          <h1>EventBridge</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (query.error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">📡</span>
          <h1>EventBridge</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(query.error) })}</div>
      </div>
    );
  }

  const buses = busesQ.data ?? [];
  const rules = rulesQ.data ?? [];

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">📡</span>
        <h1>EventBridge</h1>
        <div className="tab-bar">
          <button
            className={`tab-btn ${tab === "buses" ? "active" : ""}`}
            onClick={() => setTab("buses")}
          >
            Buses ({buses.length})
          </button>
          <button
            className={`tab-btn ${tab === "rules" ? "active" : ""}`}
            onClick={() => setTab("rules")}
          >
            Rules ({rules.length})
          </button>
        </div>
        <div className="page-actions">
          {tab === "buses" && (
            <>
              <button className="btn btn-primary" onClick={() => setShowCreateBus(true)}>
                Create bus
              </button>
              {selectedBus && (
                <button className="btn btn-danger" onClick={() => setShowDeleteBus(true)}>
                  {t("common.delete")}
                </button>
              )}
            </>
          )}
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          {tab === "buses" && (
            <DataTable
              columns={busColumns}
              data={buses}
              getRowId={(row) => row.name}
              onRowClick={setSelectedBus}
              selectedId={selectedBus?.name}
            />
          )}
          {tab === "rules" && (
            <DataTable
              columns={ruleColumns}
              data={rules}
              getRowId={(row) => row.name}
              onRowClick={setSelectedRule}
              selectedId={selectedRule?.name}
            />
          )}
        </div>
        {tab === "buses" && selectedBus && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedBus.name}</h2>
              <button className="detail-close" onClick={() => setSelectedBus(null)}>✕</button>
            </div>
            <EventBusDetail item={selectedBus} />
          </div>
        )}
        {tab === "rules" && selectedRule && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedRule.name}</h2>
              <button className="detail-close" onClick={() => setSelectedRule(null)}>✕</button>
            </div>
            <RuleDetail item={selectedRule} />
          </div>
        )}
      </div>

      <Modal open={showCreateBus} onClose={() => setShowCreateBus(false)}>
        <h2>Create event bus</h2>
        {createBusMutation.error && (
          <div className="modal-error">{String(createBusMutation.error)}</div>
        )}
        <label>
          Bus Name *
          <input
            value={formBusName}
            onChange={(e) => setFormBusName(e.target.value)}
            placeholder="my-custom-bus"
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCreateBus(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!formBusName || createBusMutation.isPending}
            onClick={() => createBusMutation.mutate()}
          >
            {createBusMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <ConfirmDialog
        open={showDeleteBus && !!selectedBus}
        title="Delete event bus"
        message={
          <>
            {t("confirm.confirmDelete", { name: "" })}
            <strong>{selectedBus?.name}</strong>?
          </>
        }
        error={deleteBusMutation.error}
        isPending={deleteBusMutation.isPending}
        onConfirm={() => selectedBus && deleteBusMutation.mutate(selectedBus.name)}
        onClose={() => setShowDeleteBus(false)}
      />
    </div>
  );
}
