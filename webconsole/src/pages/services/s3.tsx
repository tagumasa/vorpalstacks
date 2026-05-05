/**
 * S3 service page. Lists buckets with create/delete CRUD operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import {
  S3Service,
  CreateBucketRequestSchema,
  DeleteBucketRequestSchema,
  ObjectOwnership,
  type Bucket,
} from "@/gen/s3_pb";
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

/** Column definitions for the S3 bucket table. */
const getColumns = (t: TFunction): ColumnDef<Bucket, any>[] => [
  { accessorKey: "name", header: t("services.s3.bucketNameHeader"), cell: MonoCell },
  { accessorKey: "bucketregion", header: t("services.s3.regionHeader"), size: 100 },
  { accessorKey: "bucketarn", header: t("services.s3.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "creationdate", header: t("services.s3.createdHeader"), cell: DateCell },
];

/** S3 service page with bucket list, create, and delete operations. */
export function S3Page() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<Bucket | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formBucket, setFormBucket] = useState("");
  const [formOwnership, setFormOwnership] = useState<string>("BUCKETOWNERENFORCED");

  const { client, invalidate } = useServiceClient(S3Service);
  const { queryKey } = useListKey("s3");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listBuckets({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: Bucket[] = dropEmpty(data?.buckets ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createBucket(create(CreateBucketRequestSchema, {
        bucket: formBucket,
        objectownership: ObjectOwnership[formOwnership as keyof typeof ObjectOwnership] ?? ObjectOwnership.BUCKETOWNERENFORCED,
      })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormBucket("");
      setFormOwnership("BUCKETOWNERENFORCED");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (bucket: string) =>
      client.deleteBucket(create(DeleteBucketRequestSchema, { bucket })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="📦"
      title={t("services.s3.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.s3.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.s3.create")}
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
        title={t("services.s3.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formBucket}
      >
        <label>
          {t("services.s3.nameField")}
          <input
            value={formBucket}
            onChange={(e) => setFormBucket(e.target.value)}
            placeholder={t("services.s3.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.s3.objectOwnershipLabel")}
          <select
            value={formOwnership}
            onChange={(e) => setFormOwnership(e.target.value)}
            className="modal-input"
          >
            <option value="BUCKETOWNERENFORCED">{t("services.s3.ownershipBucketOwnerEnforced")}</option>
            <option value="BUCKETOWNERPREFERRED">{t("services.s3.ownershipBucketOwnerPreferred")}</option>
            <option value="OBJECTWRITER">{t("services.s3.ownershipObjectWriter")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.s3.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
