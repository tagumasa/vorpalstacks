/**
 * Kinesis service page. Lists streams with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { KinesisService, type StreamSummary, StreamMode } from "@/gen/kinesis_pb";
import { CreateStreamInputSchema, StreamModeDetailsSchema } from "@/gen/kinesis_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  BadgeCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the Kinesis stream table. */
const getColumns = (t: TFunction): ColumnDef<StreamSummary, any>[] => [
  { accessorKey: "streamname", header: t("services.kinesis.streamNameHeader"), cell: MonoCell },
  { accessorKey: "streamstatus", header: t("services.kinesis.streamStatusHeader"), cell: ({ getValue }) => <BadgeCell getValue={getValue} positive={["ACTIVE"]} negative={["CREATING", "DELETING"]} />, size: 90 },
  { accessorKey: "streamarn", header: t("services.kinesis.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "streamcreationtimestamp", header: t("services.kinesis.streamCreatedHeader"), cell: DateCell },
];

/** Kinesis service page with list, create, and delete operations. */
export function KinesisPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<StreamSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formShardCount, setFormShardCount] = useState("1");
  const [formStreamMode, setFormStreamMode] = useState(StreamMode.PROVISIONED);

  const { client, invalidate } = useServiceClient(KinesisService);
  const { queryKey } = useListKey("kinesis");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listStreams({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: StreamSummary[] = dropEmpty(data?.streamsummaries ?? [], "streamname");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createStream(
        create(CreateStreamInputSchema, {
          streamname: formName,
          shardcount: parseInt(formShardCount, 10) || 1,
          streammodedetails: create(StreamModeDetailsSchema, {
            streammode: formStreamMode,
          }),
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormShardCount("1");
      setFormStreamMode(StreamMode.PROVISIONED);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (streamName: string) =>
      client.deleteStream({ streamname: streamName }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🌊"
      title={t("services.kinesis.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.kinesis.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.kinesis.create")}
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
        getRowId={(row) => row.streamname}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.streamname}
        selected={selectedItem}
        detailTitle={selectedItem?.streamname}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.kinesis.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.kinesis.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.kinesis.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.kinesis.shardCountLabel")}
          <input
            type="number"
            min="1"
            value={formShardCount}
            onChange={(e) => setFormShardCount(e.target.value)}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.kinesis.modeLabel")}
          <select
            value={formStreamMode}
            onChange={(e) => setFormStreamMode(Number(e.target.value))}
            className="modal-input"
          >
            <option value={StreamMode.PROVISIONED}>{t("services.kinesis.modeProvisioned")}</option>
            <option value={StreamMode.ON_DEMAND}>{t("services.kinesis.modeOnDemand")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.kinesis.delete")}
        name={selectedItem?.streamname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.streamname)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
