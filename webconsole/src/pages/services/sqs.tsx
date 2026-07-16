/**
 * SQS service page — 3-panel inspector layout.
 *
 * Panel 1 (toolbar): Breadcrumb navigation
 * Panel 2 (table):   Queue list with checkbox multi-select
 * Panel 3 (detail):  Queue detail (URL + attributes)
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { SQSService, CreateQueueRequestSchema, SendMessageRequestSchema } from "@/gen/sqs_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  useServiceClient,
} from "@/components/shared/service-page";
import {
  checkboxColumn,
  Breadcrumb,
  SelectionBadge,
  DetailPanel,
  DetailEmpty,
  useSelection,
} from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { TagSection, useTags } from "@/components/shared/tag-section";

// ─── Row Type ───────────────────────────────────────────────────

interface TableRow {
  url: string;
  name: string;
}

/** Extracts the queue name from a full SQS URL. */
function queueNameFromUrl(url: string): string {
  const idx = url.lastIndexOf("/");
  return idx >= 0 ? url.slice(idx + 1) : url;
}

// ─── Column Definitions ─────────────────────────────────────────

const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.sqs.queueNameHeader"), cell: MonoCell },
  { accessorKey: "url", header: t("services.sqs.urlHeader"), cell: SmallMonoCell },
];

// ─── Detail panel tab ───────────────────────────────────────────

type DetailTab = "detail" | "json" | "messages";

// ─── SQS Page ───────────────────────────────────────────────────

export function SQSPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(SQSService);
  const { queryKey } = useListKey("sqs");
  const columns = getColumns(t);

  // ── Selection state ──────────────────────────────────────────
  const { selected: selectedUrls, toggle: toggleUrl, toggleAll: toggleAllUrls, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  // ── Create form state ────────────────────────────────────────
  const [formName, setFormName] = useState("");
  const [formVisibilityTimeout, setFormVisibilityTimeout] = useState("30");
  const [formRetentionPeriod, setFormRetentionPeriod] = useState("345600");
  const [formDelaySeconds, setFormDelaySeconds] = useState("0");

  // ── Data ─────────────────────────────────────────────────────
  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<TableRow, Awaited<ReturnType<typeof client.listQueues>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listQueues({ nexttoken: token || undefined }),
    getItems: (r) => (r.queueurls ?? []).map((url) => ({ url, name: queueNameFromUrl(url) })),
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "url");

  // ── Mutations ────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: () =>
      client.createQueue(create(CreateQueueRequestSchema, {
        queuename: formName,
        attributes: {
          VisibilityTimeout: formVisibilityTimeout,
          MessageRetentionPeriod: formRetentionPeriod,
          DelaySeconds: formDelaySeconds,
        },
      })),
    onSuccess: () => {
      invalidateList();
      setShowCreate(false);
      setFormName("");
      setFormVisibilityTimeout("30");
      setFormRetentionPeriod("345600");
      setFormDelaySeconds("0");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (queueUrl: string) => client.deleteQueue({ queueurl: queueUrl }),
    onSuccess: () => {
      invalidateList();
      setShowDelete(false);
      setSelectedItem(null);
      clearSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (urls: string[]) => {
      const results = await Promise.allSettled(
        urls.map((url) => client.deleteQueue({ queueurl: url }))
      );
      return results;
    },
    onSuccess: (_data, urls) => {
      invalidateList();
      setShowBatchDelete(false);
      clearSelection();
      setSelectedItem((prev) => (prev && urls.includes(prev.url) ? null : prev));
    },
  });

  // ── Handlers ─────────────────────────────────────────────────

  const handleRowClick = (row: TableRow) => {
    setSelectedItem(row);
    setDetailTab("detail");
  };

  const allIds = items.map((i) => i.url);

  const { tags: itemTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (queueUrl: string) => {
        const res = await client.listQueueTags({ queueurl: queueUrl });
        return Object.entries(res.tags ?? {}).map(([key, value]) => ({ key, value }));
      },
      tagResource: async (queueUrl: string, tags) => {
        const m: Record<string, string> = {};
        for (const t of tags) m[t.key] = t.value;
        await client.tagQueue({ queueurl: queueUrl, tags: m });
      },
      untagResource: async (queueUrl: string, tagKeys: string[]) => {
        await client.untagQueue({ queueurl: queueUrl, tagkeys: tagKeys });
      },
    },
    selectedItem?.url || undefined,
  );

  // ── Message send/receive ─────────────────────────────────────
  const [showSend, setShowSend] = useState(false);
  const [formMsgBody, setFormMsgBody] = useState("");
  const messagesQueryKey = ["sqs", "messages", selectedItem?.url ?? ""];

  const { data: messagesData, isLoading: msgsLoading, refetch: refetchMessages } = useQuery({
    queryKey: messagesQueryKey,
    queryFn: () => client.receiveMessage({ queueurl: selectedItem!.url, maxnumberofmessages: 10, visibilitytimeout: 30, waittimeseconds: 0 }),
    enabled: false,
  });
  const receivedMessages = messagesData?.messages ?? [];

  const sendMsgMutation = useMutation({
    mutationFn: () => client.sendMessage(create(SendMessageRequestSchema, { queueurl: selectedItem!.url, messagebody: formMsgBody })),
    onSuccess: () => { setShowSend(false); setFormMsgBody(""); },
  });

  const deleteMsgMutation = useMutation({
    mutationFn: (receiptHandle: string) => client.deleteMessage({ queueurl: selectedItem!.url, receipthandle: receiptHandle }),
    onSuccess: () => refetchMessages(),
  });

  // ── Detail Panel ─────────────────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedItem) {
      return <DetailEmpty message={t("common.noItemSelected")} />;
    }

    const detailTabs = [
      { key: "detail", label: t("common.tabDetail") },
      { key: "messages", label: t("common.tabMessages") },
      { key: "json", label: t("common.rawJson") },
    ];

    return (
      <DetailPanel
        title={selectedItem.name}
        titleIcon="📬"
        tabs={detailTabs}
        activeTab={detailTab}
        onTabChange={(k) => setDetailTab(k as DetailTab)}
        actions={
          <button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>
            {t("common.delete")}
          </button>
        }
      >
        {detailTab === "detail" ? (
          <><table className="settings-table">
            <tbody>
              <tr>
                <td className="detail-label-sm">Name</td>
                <td className="cell-mono">{selectedItem.name}</td>
              </tr>
              <tr>
                <td className="detail-label">{t("common.fields.url")}</td>
                <td className="cell-mono cell-long">{selectedItem.url}</td>
              </tr>
            </tbody>
          </table>
          <TagSection tags={itemTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} /></>
        ) : detailTab === "messages" ? (
          <div className="flex-fill-scroll" style={{ padding: 8 }}>
            <div style={{ marginBottom: 8, display: "flex", gap: 8 }}>
              <button className="btn btn-primary btn-sm" onClick={() => setShowSend(true)}>{t("services.sqs.sendMessage")}</button>
              <button className="btn btn-secondary btn-sm" onClick={() => refetchMessages()} disabled={msgsLoading}>{msgsLoading ? t("common.loading") : t("services.sqs.pollMessages")}</button>
            </div>
            {receivedMessages.length > 0 ? (
              <div className="message-list">
                {receivedMessages.map((msg, i) => (
                  <div key={i} className="message-item">
                    <div className="message-header">
                      <span className="cell-mono" style={{ fontSize: "0.8em" }}>MD5: {msg.md5ofbody || "\u2014"}</span>
                      <button className="btn btn-danger btn-sm" style={{ marginLeft: "auto" }} onClick={() => msg.receipthandle && deleteMsgMutation.mutate(msg.receipthandle)}>{t("common.delete")}</button>
                    </div>
                    <pre className="message-body">{msg.body || "\u2014"}</pre>
                  </div>
                ))}
              </div>
            ) : (
              <div className="text-muted" style={{ padding: 16 }}>{msgsLoading ? t("common.loading") : t("services.sqs.noMessages")}</div>
            )}
          </div>
        ) : (
          <JsonViewer data={selectedItem} />
        )}
      </DetailPanel>
    );
  };

  // ── Render ───────────────────────────────────────────────────

  return (
    <ServicePageLayout
      icon="📬"
      title={t("services.sqs.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.sqs.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.sqs.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedUrls.size === 0}
            onClick={() => setShowBatchDelete(true)}
          >
            {t("common.deleteSelected")}
            {selectedUrls.size > 0 && <span className="batch-count">({selectedUrls.size})</span>}
          </button>
        </>
      }
    >
      <div className="inspector-toolbar">
        <Breadcrumb parts={[
          { label: t("services.sqs.title") },
          { label: t("services.sqs.countLabel") },
        ]} />
        <div className="toolbar-selection-info">
          <SelectionBadge count={selectedUrls.size} label={t("common.selectedCount", { count: selectedUrls.size })} />
        </div>
      </div>

      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-sqs">
          <div className="flex-fill-scroll">
            <DataTable
              columns={[
                checkboxColumn<TableRow>(selectedUrls, toggleUrl, () => toggleAllUrls(allIds), allIds, t, (row) => row.url),
                ...columns,
              ]}
              data={items}
              getRowId={(row) => row.url}
              onRowClick={handleRowClick}
              selectedId={selectedItem?.url}
              hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore}
            />
          </div>
          {renderDetailPanel()}
        </Splitter>
      ) : (
        <div className="empty-state">{t("common.noData")}</div>
      )}

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.sqs.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.sqs.nameField")}
          <input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.sqs.placeholder")} className="modal-input" />
        </label>
        <label>
          {t("services.sqs.visibilityTimeoutLabel")}
          <input type="number" min={0} max={43200} value={formVisibilityTimeout} onChange={(e) => setFormVisibilityTimeout(e.target.value)} className="modal-input" />
        </label>
        <label>
          {t("services.sqs.messageRetentionLabel")}
          <input type="number" min={60} max={1209600} value={formRetentionPeriod} onChange={(e) => setFormRetentionPeriod(e.target.value)} className="modal-input" />
        </label>
        <label>
          {t("services.sqs.delaySecondsLabel")}
          <input type="number" min={0} max={900} value={formDelaySeconds} onChange={(e) => setFormDelaySeconds(e.target.value)} className="modal-input" />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.sqs.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.url)}
        onClose={() => setShowDelete(false)}
      />

      <ServiceDeleteDialog
        open={showBatchDelete}
        title={t("common.deleteSelected")}
        name={`${selectedUrls.size} ${t("services.sqs.countLabel")}`}
        error={batchDeleteMutation.error}
        isPending={batchDeleteMutation.isPending}
        onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedUrls))}
        onClose={() => setShowBatchDelete(false)}
      />

      <ServiceCreateModal
        open={showSend}
        onClose={() => setShowSend(false)}
        title={t("services.sqs.sendMessage")}
        error={sendMsgMutation.error}
        isPending={sendMsgMutation.isPending}
        onCreate={() => sendMsgMutation.mutate()}
        disabled={!formMsgBody}
      >
        <label>
          {t("services.sqs.messageBody")}
          <textarea value={formMsgBody} onChange={(e) => setFormMsgBody(e.target.value)} className="modal-input" rows={6} style={{ fontFamily: "monospace", fontSize: "0.85em" }} />
        </label>
      </ServiceCreateModal>
    </ServicePageLayout>
  );
}
