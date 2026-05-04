/**
 * SQS service page. Lists queues with create/delete CRUD operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { SQSService, CreateQueueRequestSchema } from "@/gen/sqs_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Derived row shape for the queue list table. */
interface TableRow {
  url: string;
  name: string;
}

/** Extracts the queue name from a full SQS URL. */
function queueNameFromUrl(url: string): string {
  const idx = url.lastIndexOf("/");
  return idx >= 0 ? url.slice(idx + 1) : url;
}

/** Column definitions for the SQS queue table. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.sqs.queueNameHeader"), cell: MonoCell },
  { accessorKey: "url", header: t("services.sqs.urlHeader"), cell: SmallMonoCell },
];

/** SQS service page with list, create, and delete operations. */
export function SQSPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formVisibilityTimeout, setFormVisibilityTimeout] = useState("30");
  const [formRetentionPeriod, setFormRetentionPeriod] = useState("345600");
  const [formDelaySeconds, setFormDelaySeconds] = useState("0");

  const { client, invalidate } = useServiceClient(SQSService);
  const { queryKey } = useListKey("sqs");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listQueues({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TableRow[] = dropEmpty(
    (data?.queueurls ?? []).map((url) => ({ url, name: queueNameFromUrl(url) })),
    "url",
  );

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
      invalidate(queryKey);
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
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

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
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "sqs-items" }}
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.url}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.url}
        selected={selectedItem}
        detailTitle={selectedItem?.name}
        onDetailClose={() => setSelectedItem(null)}
      />

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
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.sqs.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sqs.visibilityTimeoutLabel")}
          <input
            type="number"
            min={0}
            max={43200}
            value={formVisibilityTimeout}
            onChange={(e) => setFormVisibilityTimeout(e.target.value)}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sqs.messageRetentionLabel")}
          <input
            type="number"
            min={60}
            max={1209600}
            value={formRetentionPeriod}
            onChange={(e) => setFormRetentionPeriod(e.target.value)}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sqs.delaySecondsLabel")}
          <input
            type="number"
            min={0}
            max={900}
            value={formDelaySeconds}
            onChange={(e) => setFormDelaySeconds(e.target.value)}
            className="modal-input"
          />
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
    </ServicePageLayout>
  );
}
