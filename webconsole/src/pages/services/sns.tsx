/**
 * SNS service page. Lists topics with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { SNSService } from "@/gen/sns_pb";
import { CreateTopicInputSchema } from "@/gen/sns_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  SmallMonoCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Derived row shape for the SNS topic list table. */
interface TableRow {
  topicarn: string;
}

/** Column definitions for the SNS topic table. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "topicarn", header: t("services.sns.topicArnHeader"), cell: SmallMonoCell },
];

/** SNS service page with list, create, and delete operations. */
export function SNSPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formFifo, setFormFifo] = useState(false);
  const [formDisplayName, setFormDisplayName] = useState("");
  const [formTags, setFormTags] = useState("");

  const { client, invalidate } = useServiceClient(SNSService);
  const { queryKey } = useListKey("sns");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listTopics({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TableRow[] = dropEmpty(
    (data?.topics ?? []).map((topic) => ({
      topicarn: topic.topicarn,
    })),
    "topicarn",
  );

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
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormFifo(false);
      setFormDisplayName("");
      setFormTags("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (topicarn: string) =>
      client.deleteTopic({ topicarn }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

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
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "sns-items" }}
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.topicarn}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.topicarn}
        selected={selectedItem}
        detailTitle={selectedItem?.topicarn}
        onDetailClose={() => setSelectedItem(null)}
      />

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
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.sns.placeholder")}
            className="modal-input"
          />
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formFifo}
            onChange={(e) => setFormFifo(e.target.checked)}
          />
          {t("services.sns.fifoLabel")}
        </label>
        <label>
          {t("services.sns.displayNameLabel")}
          <input
            value={formDisplayName}
            onChange={(e) => setFormDisplayName(e.target.value)}
            placeholder={t("services.sns.displayNamePlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sns.tagsLabel")}
          <textarea
            value={formTags}
            onChange={(e) => setFormTags(e.target.value)}
            placeholder='{"key":"value"}'
            rows={3}
            className="modal-input"
            style={{ fontFamily: "monospace", fontSize: "0.85em" }}
          />
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
    </ServicePageLayout>
  );
}
