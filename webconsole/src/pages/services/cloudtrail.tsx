/**
 * CloudTrail service page. Lists trails with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { CloudTrailService, type TrailInfo } from "@/gen/cloudtrail_pb";
import { CreateTrailRequestSchema } from "@/gen/cloudtrail_pb";
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

/** Column definitions for the CloudTrail trail table. */
const getColumns = (t: TFunction): ColumnDef<TrailInfo, any>[] => [
  { accessorKey: "name", header: t("services.cloudtrail.trailNameHeader"), cell: MonoCell },
  { accessorKey: "trailarn", header: t("services.cloudtrail.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "homeregion", header: t("services.cloudtrail.homeRegionHeader"), size: 120 },
];

/** CloudTrail service page with list, create, and delete operations. */
export function CloudTrailPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TrailInfo | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formS3Bucket, setFormS3Bucket] = useState("");
  const [formMultiRegion, setFormMultiRegion] = useState(true);
  const [formGlobalEvents, setFormGlobalEvents] = useState(true);

  const { client, invalidate } = useServiceClient(CloudTrailService);
  const { queryKey } = useListKey("cloudtrail");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listTrails({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TrailInfo[] = dropEmpty(data?.trails ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createTrail(
        create(CreateTrailRequestSchema, {
          name: formName,
          s3bucketname: formS3Bucket,
          ismultiregiontrail: formMultiRegion,
          includeglobalserviceevents: formGlobalEvents,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormS3Bucket("");
      setFormMultiRegion(true);
      setFormGlobalEvents(true);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteTrail({ name }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="📋"
      title={t("services.cloudtrail.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cloudtrail.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cloudtrail.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "cloudtrail-items" }}
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.name}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.name}
        selected={selectedItem}
        detailTitle={selectedItem?.name}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.cloudtrail.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formS3Bucket}
      >
        <label>
          {t("services.cloudtrail.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.cloudtrail.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.cloudtrail.s3BucketLabel")}
          <input
            value={formS3Bucket}
            onChange={(e) => setFormS3Bucket(e.target.value)}
            placeholder={t("services.cloudtrail.s3BucketPlaceholder")}
            className="modal-input"
          />
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formMultiRegion}
            onChange={(e) => setFormMultiRegion(e.target.checked)}
          />
          {t("services.cloudtrail.multiRegionLabel")}
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formGlobalEvents}
            onChange={(e) => setFormGlobalEvents(e.target.checked)}
          />
          {t("services.cloudtrail.globalServiceEventsLabel")}
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cloudtrail.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
