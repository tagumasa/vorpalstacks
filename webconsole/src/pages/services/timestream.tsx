/**
 * Timestream service page. Lists databases with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { TimestreamWriteService } from "@/gen/timestreamwrite_pb";
import { CreateDatabaseRequestSchema } from "@/gen/timestreamwrite_pb";
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

/** Derived row shape for the Timestream database list table. */
interface TableRow {
  databasename: string;
  arn: string;
  tablecount: number;
  creationtime: string;
  lastupdatedtime: string;
}

/** Column definitions for the Timestream database table. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "databasename", header: t("services.timestream.databaseNameHeader"), cell: MonoCell },
  { accessorKey: "tablecount", header: t("services.timestream.tableCountHeader"), size: 80 },
  { accessorKey: "creationtime", header: t("services.timestream.createdHeader"), cell: DateCell },
  { accessorKey: "lastupdatedtime", header: t("services.timestream.lastUpdatedHeader"), cell: DateCell },
  { accessorKey: "arn", header: t("services.timestream.arnHeader"), cell: SmallMonoCell },
];

/** Timestream service page with list, create, and delete operations. */
export function TimestreamPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");
  const [formMagneticRetentionDays, setFormMagneticRetentionDays] = useState(365);
  const [formMemoryRetentionHours, setFormMemoryRetentionHours] = useState(24);

  const { client, invalidate } = useServiceClient(TimestreamWriteService);
  const { queryKey } = useListKey("timestream");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listDatabases({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TableRow[] = dropEmpty(
    (data?.databases ?? []).map((db) => ({
      databasename: db.databasename,
      arn: db.arn,
      tablecount: Number(db.tablecount ?? 0),
      creationtime: db.creationtime ?? "",
      lastupdatedtime: db.lastupdatedtime ?? "",
    })),
    "databasename",
  );

  const createMutation = useMutation({
    mutationFn: () => {
      const req: Record<string, any> = {
        databasename: formName,
        magneticretentiondays: formMagneticRetentionDays,
        memoryretentionhours: formMemoryRetentionHours,
      };
      if (formKmsKeyId) req.kmskeyid = formKmsKeyId;
      return client.createDatabase(create(CreateDatabaseRequestSchema, req));
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormKmsKeyId("");
      setFormMagneticRetentionDays(365);
      setFormMemoryRetentionHours(24);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (databasename: string) =>
      client.deleteDatabase({ databasename }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="⏱"
      title={t("services.timestream.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.timestream.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.timestream.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "timestream-items" }}
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.databasename}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.databasename}
        selected={selectedItem}
        detailTitle={selectedItem?.databasename}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.timestream.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.timestream.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.timestream.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.timestream.kmsKeyLabel")}
          <input
            value={formKmsKeyId}
            onChange={(e) => setFormKmsKeyId(e.target.value)}
            placeholder={t("services.timestream.kmsKeyPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.timestream.magneticRetentionLabel")}
          <input
            type="number"
            min={1}
            max={73000}
            value={formMagneticRetentionDays}
            onChange={(e) => setFormMagneticRetentionDays(Number(e.target.value))}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.timestream.memoryRetentionLabel")}
          <input
            type="number"
            min={1}
            max={8766}
            value={formMemoryRetentionHours}
            onChange={(e) => setFormMemoryRetentionHours(Number(e.target.value))}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.timestream.delete")}
        name={selectedItem?.databasename}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.databasename)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
