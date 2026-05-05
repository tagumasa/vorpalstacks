/**
 * CloudWatch Logs service page. Lists log groups with create/delete operations.
 *
 * BUG FIX: `LogGroupSummary` only has `loggrouparn`, `loggroupclass`,
 * `loggroupname`. Removed `creationtime` and `storedbytes` columns that
 * referenced non-existent fields.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { CloudWatchLogsService, type LogGroupSummary } from "@/gen/cloudwatchlogs_pb";
import { CreateLogGroupRequestSchema, PutRetentionPolicyRequestSchema } from "@/gen/cloudwatchlogs_pb";
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

/** Column definitions for the CloudWatch Logs log group table. */
const getColumns = (t: TFunction): ColumnDef<LogGroupSummary, any>[] => [
  { accessorKey: "loggroupname", header: t("services.cloudwatchlogs.logGroupHeader"), cell: MonoCell },
  { accessorKey: "loggroupclass", header: t("services.cloudwatchlogs.classHeader"), size: 100 },
  { accessorKey: "loggrouparn", header: t("services.cloudwatchlogs.arnHeader"), cell: SmallMonoCell },
];

/** CloudWatch Logs service page with list, create, and delete operations. */
export function CloudWatchLogsPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<LogGroupSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");
  const [formRetentionDays, setFormRetentionDays] = useState("");

  const { client, invalidate } = useServiceClient(CloudWatchLogsService);
  const { queryKey } = useListKey("cloudwatchlogs");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listLogGroups({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: LogGroupSummary[] = dropEmpty(data?.loggroups ?? [], "loggroupname");

  const createMutation = useMutation({
    mutationFn: async () => {
      const req: Record<string, any> = { loggroupname: formName };
      if (formKmsKeyId) req.kmskeyid = formKmsKeyId;
      await client.createLogGroup(create(CreateLogGroupRequestSchema, req));
      if (formRetentionDays) {
        await client.putRetentionPolicy(
          create(PutRetentionPolicyRequestSchema, {
            loggroupname: formName,
            retentionindays: Number(formRetentionDays),
          }),
        );
      }
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormKmsKeyId("");
      setFormRetentionDays("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (logGroupName: string) =>
      client.deleteLogGroup({ loggroupname: logGroupName }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="📜"
      title={t("services.cloudwatchlogs.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cloudwatchlogs.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cloudwatchlogs.create")}
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
        getRowId={(row) => row.loggroupname}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.loggroupname}
        selected={selectedItem}
        detailTitle={selectedItem?.loggroupname}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.cloudwatchlogs.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.cloudwatchlogs.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.cloudwatchlogs.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatchlogs.kmsKeyLabel")}
          <input
            value={formKmsKeyId}
            onChange={(e) => setFormKmsKeyId(e.target.value)}
            placeholder={t("services.cloudwatchlogs.kmsKeyPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudwatchlogs.retentionLabel")}
          <select
            value={formRetentionDays}
            onChange={(e) => setFormRetentionDays(e.target.value)}
            className="modal-input"
          >
            <option value="">{t("services.cloudwatchlogs.retentionNone")}</option>
            <option value="1">1</option>
            <option value="3">3</option>
            <option value="5">5</option>
            <option value="7">7</option>
            <option value="14">14</option>
            <option value="30">30</option>
            <option value="60">60</option>
            <option value="90">90</option>
            <option value="120">120</option>
            <option value="150">150</option>
            <option value="180">180</option>
            <option value="365">365</option>
            <option value="400">400</option>
            <option value="545">545</option>
            <option value="731">731</option>
            <option value="1827">1827</option>
            <option value="3653">3653</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cloudwatchlogs.delete")}
        name={selectedItem?.loggroupname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.loggroupname)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
