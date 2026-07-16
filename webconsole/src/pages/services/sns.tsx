/**
 * SNS service page — 3-panel inspector layout.
 *
 * Panel 1 (toolbar): Breadcrumb navigation
 * Panel 2 (table):   Topic list with checkbox multi-select
 * Panel 3 (detail):  Topic detail (ARN + JSON view)
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation, useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { SNSService, type Subscription } from "@/gen/sns_pb";
import { CreateTopicInputSchema, SubscribeInputSchema, PublishInputSchema } from "@/gen/sns_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
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
  topicarn: string;
}

// ─── Column Definitions ─────────────────────────────────────────

const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "topicarn", header: t("services.sns.topicArnHeader"), cell: SmallMonoCell },
];

// ─── Detail panel tab ───────────────────────────────────────────

type DetailTab = "detail" | "subscriptions" | "json";

// ─── SNS Page ───────────────────────────────────────────────────

export function SNSPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(SNSService);
  const { queryKey } = useListKey("sns");
  const columns = getColumns(t);

  // ── Selection state ──────────────────────────────────────────
  const { selected: selectedArns, toggle: toggleArn, toggleAll: toggleAllArns, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  // ── Create form state ────────────────────────────────────────
  const [formName, setFormName] = useState("");
  const [formFifo, setFormFifo] = useState(false);
  const [formDisplayName, setFormDisplayName] = useState("");
  const [formTags, setFormTags] = useState("");

  // ── Data ─────────────────────────────────────────────────────
  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<TableRow, Awaited<ReturnType<typeof client.listTopics>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listTopics({ nexttoken: token || undefined }),
    getItems: (r) => (r.topics ?? []).map((topic) => ({ topicarn: topic.topicarn })),
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "topicarn");

  // ── Mutations ────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: () => {
      const attributes: Record<string, string> = {};
      if (formFifo) attributes.FifoTopic = "true";
      if (formDisplayName) attributes.DisplayName = formDisplayName;
      let tags: { key: string; value: string }[] = [];
      if (formTags.trim()) {
        try {
          const parsed = JSON.parse(formTags);
          tags = Object.entries(parsed).map(([key, value]) => ({ key, value: String(value) }));
        } catch {
          throw new Error(t("common.invalidJson", { field: "Tags" }));
        }
      }
      return client.createTopic(
        create(CreateTopicInputSchema, {
          name: formName,
          ...(Object.keys(attributes).length > 0 ? { attributes } : {}),
          ...(tags.length > 0 ? { tags } : {}),
        }),
      );
    },
    onSuccess: () => {
      invalidateList();
      setShowCreate(false);
      setFormName("");
      setFormFifo(false);
      setFormDisplayName("");
      setFormTags("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (topicarn: string) => client.deleteTopic({ topicarn }),
    onSuccess: () => {
      invalidateList();
      setShowDelete(false);
      setSelectedItem(null);
      clearSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (arns: string[]) => {
      const results = await Promise.allSettled(
        arns.map((arn) => client.deleteTopic({ topicarn: arn }))
      );
      return results;
    },
    onSuccess: (_data, arns) => {
      invalidateList();
      setShowBatchDelete(false);
      clearSelection();
      setSelectedItem((prev) => (prev && arns.includes(prev.topicarn) ? null : prev));
    },
  });

  // ── Handlers ─────────────────────────────────────────────────

  const handleRowClick = (row: TableRow) => {
    setSelectedItem(row);
    setDetailTab("detail");
  };

  const allIds = items.map((i) => i.topicarn);

  const { tags: itemTags, isLoading: tagsLoading, addTags, removeTag, isPending: tagsPending } = useTags(
    {
      queryKeyBase: [...queryKey, "tags"],
      fetchTags: async (arn: string) => {
        const res = await client.listTagsForResource({ resourcearn: arn });
        return (res.tags ?? []).map((t) => ({ key: t.key, value: t.value }));
      },
      tagResource: async (arn: string, tags) => {
        await client.tagResource({ resourcearn: arn, tags: tags.map((t) => ({ key: t.key, value: t.value })) });
      },
      untagResource: async (arn: string, tagKeys: string[]) => {
        await client.untagResource({ resourcearn: arn, tagkeys: tagKeys });
      },
    },
    selectedItem?.topicarn || undefined,
  );

  // ── Subscriptions ────────────────────────────────────────────
  const [showSubscribe, setShowSubscribe] = useState(false);
  const [showPublish, setShowPublish] = useState(false);
  const [formProtocol, setFormProtocol] = useState("sqs");
  const [formEndpoint, setFormEndpoint] = useState("");
  const [formMsg, setFormMsg] = useState("");

  const { data: subsData, refetch: refetchSubs } = useQuery({
    queryKey: ["sns", "subs", selectedItem?.topicarn ?? ""],
    queryFn: () => client.listSubscriptionsByTopic({ topicarn: selectedItem!.topicarn }),
    enabled: false,
  });
  const subscriptions = subsData?.subscriptions ?? [];

  const subscribeMutation = useMutation({
    mutationFn: () => client.subscribe(create(SubscribeInputSchema, {
      topicarn: selectedItem!.topicarn, protocol: formProtocol, endpoint: formEndpoint, returnsubscriptionarn: true,
    })),
    onSuccess: () => { setShowSubscribe(false); setFormEndpoint(""); refetchSubs(); },
  });

  const unsubscribeMutation = useMutation({
    mutationFn: (subArn: string) => client.unsubscribe({ subscriptionarn: subArn }),
    onSuccess: () => refetchSubs(),
  });

  const publishMutation = useMutation({
    mutationFn: () => client.publish(create(PublishInputSchema, {
      topicarn: selectedItem!.topicarn, message: formMsg,
    })),
    onSuccess: () => { setShowPublish(false); setFormMsg(""); },
  });

  /** Extract topic name from ARN for display. */
  const topicName = (arn: string) => {
    const parts = arn.split(":");
    return parts[5] ?? arn;
  };

  // ── Detail Panel ─────────────────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedItem) {
      return <DetailEmpty message={t("common.noItemSelected")} />;
    }

    const detailTabs = [
      { key: "detail", label: t("common.tabDetail") },
      { key: "subscriptions", label: t("services.sns.subscriptions") },
      { key: "json", label: t("common.rawJson") },
    ];

    return (
      <DetailPanel
        title={topicName(selectedItem.topicarn)}
        titleIcon="📢"
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
                <td className="detail-label-sm">ARN</td>
                <td className="cell-mono cell-long">{selectedItem.topicarn}</td>
              </tr>
              <tr>
                <td className="detail-label">{t("common.fields.name")}</td>
                <td className="cell-mono">{topicName(selectedItem.topicarn)}</td>
              </tr>
            </tbody>
          </table>
          <TagSection tags={itemTags} isLoading={tagsLoading} onAddTags={addTags} onRemoveTag={removeTag} isPending={tagsPending} /></>
        ) : detailTab === "subscriptions" ? (
          <div className="flex-fill-scroll" style={{ padding: 8 }}>
            <div style={{ marginBottom: 8, display: "flex", gap: 8 }}>
              <button className="btn btn-primary btn-sm" onClick={() => setShowSubscribe(true)}>{t("services.sns.subscribe")}</button>
              <button className="btn btn-secondary btn-sm" onClick={() => refetchSubs()}>{t("services.sns.refreshSubs")}</button>
              <button className="btn btn-primary btn-sm" onClick={() => setShowPublish(true)}>{t("services.sns.publish")}</button>
            </div>
            {subscriptions.length > 0 ? (
              <table className="settings-table"><tbody>
                {subscriptions.map((sub: Subscription, i) => (
                  <tr key={i}>
                    <td className="cell-mono">{sub.protocol}</td>
                    <td className="cell-mono cell-long">{sub.endpoint || "\u2014"}</td>
                    <td><button className="btn btn-danger btn-sm" onClick={() => unsubscribeMutation.mutate(sub.subscriptionarn)}>{t("services.sns.unsubscribe")}</button></td>
                  </tr>
                ))}
              </tbody></table>
            ) : (
              <div className="text-muted" style={{ padding: 16 }}>{t("services.sns.noSubscriptions")}</div>
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
      icon="📢"
      title={t("services.sns.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.sns.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.sns.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedArns.size === 0}
            onClick={() => setShowBatchDelete(true)}
          >
            {t("common.deleteSelected")}
            {selectedArns.size > 0 && <span className="batch-count">({selectedArns.size})</span>}
          </button>
        </>
      }
    >
      <div className="inspector-toolbar">
        <Breadcrumb parts={[
          { label: t("services.sns.title") },
          { label: t("services.sns.countLabel") },
        ]} />
        <div className="toolbar-selection-info">
          <SelectionBadge count={selectedArns.size} label={t("common.selectedCount", { count: selectedArns.size })} />
        </div>
      </div>

      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-sns">
          <div className="flex-fill-scroll">
            <DataTable
              columns={[
                checkboxColumn<TableRow>(selectedArns, toggleArn, () => toggleAllArns(allIds), allIds, t, (row) => row.topicarn),
                ...columns,
              ]}
              data={items}
              getRowId={(row) => row.topicarn}
              onRowClick={handleRowClick}
              selectedId={selectedItem?.topicarn}
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
        title={t("services.sns.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.sns.nameField")}
          <input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.sns.placeholder")} className="modal-input" />
        </label>
        <label className="checkbox-label">
          <input type="checkbox" checked={formFifo} onChange={(e) => setFormFifo(e.target.checked)} />
          {t("services.sns.fifoLabel")}
        </label>
        <label>
          {t("services.sns.displayNameLabel")}
          <input value={formDisplayName} onChange={(e) => setFormDisplayName(e.target.value)} placeholder={t("services.sns.displayNamePlaceholder")} className="modal-input" />
        </label>
        <label>
          {t("services.sns.tagsLabel")}
          <textarea value={formTags} onChange={(e) => setFormTags(e.target.value)} placeholder='{"key":"value"}' rows={3} className="modal-input" style={{ fontFamily: "monospace", fontSize: "0.85em" }} />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.sns.delete")}
        name={selectedItem?.topicarn}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.topicarn)}
        onClose={() => setShowDelete(false)}
      />

      <ServiceDeleteDialog
        open={showBatchDelete}
        title={t("common.deleteSelected")}
        name={`${selectedArns.size} ${t("services.sns.countLabel")}`}
        error={batchDeleteMutation.error}
        isPending={batchDeleteMutation.isPending}
        onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedArns))}
        onClose={() => setShowBatchDelete(false)}
      />

      <ServiceCreateModal open={showSubscribe} onClose={() => setShowSubscribe(false)} title={t("services.sns.subscribe")} error={subscribeMutation.error} isPending={subscribeMutation.isPending} onCreate={() => subscribeMutation.mutate()} disabled={!formEndpoint}>
        <label>{t("services.sns.protocolLabel")}<select value={formProtocol} onChange={(e) => setFormProtocol(e.target.value)} className="modal-input"><option value="sqs">SQS</option><option value="lambda">Lambda</option><option value="https">HTTPS</option><option value="http">HTTP</option><option value="email">Email</option><option value="sms">SMS</option><option value="firehose">Firehose</option></select></label>
        <label>{t("services.sns.endpointLabel")}<input value={formEndpoint} onChange={(e) => setFormEndpoint(e.target.value)} placeholder="arn:aws:sqs:..." className="modal-input" /></label>
      </ServiceCreateModal>

      <ServiceCreateModal open={showPublish} onClose={() => setShowPublish(false)} title={t("services.sns.publish")} error={publishMutation.error} isPending={publishMutation.isPending} onCreate={() => publishMutation.mutate()} disabled={!formMsg}>
        <label>{t("services.sns.messageLabel")}<textarea value={formMsg} onChange={(e) => setFormMsg(e.target.value)} className="modal-input" rows={6} style={{ fontFamily: "monospace", fontSize: "0.85em" }} /></label>
      </ServiceCreateModal>
    </ServicePageLayout>
  );
}
