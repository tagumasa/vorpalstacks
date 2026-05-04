/**
 * Secrets Manager service page. Lists secrets with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { SecretsManagerService, type SecretListEntry } from "@/gen/secretsmanager_pb";
import { CreateSecretRequestSchema } from "@/gen/secretsmanager_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  FallbackCell,
  BooleanCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the Secrets Manager secret table. */
const getColumns = (t: TFunction): ColumnDef<SecretListEntry, any>[] => [
  { accessorKey: "name", header: t("services.secretsmanager.secretNameHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.secretsmanager.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "rotationenabled", header: t("services.secretsmanager.rotationHeader"), cell: BooleanCell, size: 80 },
  { accessorKey: "createddate", header: t("services.secretsmanager.createdHeader"), cell: DateCell },
  { accessorKey: "deleteddate", header: t("services.secretsmanager.deletedHeader"), cell: DateCell },
  { accessorKey: "kmskeyid", header: t("services.secretsmanager.kmsKeyHeader"), cell: SmallMonoCell },
];

/** Secrets Manager service page with list, create, and delete operations. */
export function SecretsManagerPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<SecretListEntry | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDesc, setFormDesc] = useState("");
  const [formSecretValue, setFormSecretValue] = useState("");
  const [formKmsKeyId, setFormKmsKeyId] = useState("");

  const { client, invalidate } = useServiceClient(SecretsManagerService);
  const { queryKey } = useListKey("secretsmanager");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listSecrets({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: SecretListEntry[] = dropEmpty(data?.secretlist ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () => {
      const req: Record<string, any> = {
        name: formName,
        description: formDesc,
        secretstring: formSecretValue,
      };
      if (formKmsKeyId) req.kmskeyid = formKmsKeyId;
      return client.createSecret(create(CreateSecretRequestSchema, req));
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormDesc("");
      setFormSecretValue("");
      setFormKmsKeyId("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (secretName: string) =>
      client.deleteSecret({
        secretid: secretName,
        forcedeletewithoutrecovery: true,
      }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🗝️"
      title={t("services.secretsmanager.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.secretsmanager.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.secretsmanager.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "secretsmanager-items" }}
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
        title={t("services.secretsmanager.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.secretsmanager.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.secretsmanager.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.secretsmanager.descLabel")}
          <input
            value={formDesc}
            onChange={(e) => setFormDesc(e.target.value)}
            placeholder={t("services.secretsmanager.descPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.secretsmanager.secretValueLabel")}
          <textarea
            value={formSecretValue}
            onChange={(e) => setFormSecretValue(e.target.value)}
            placeholder={t("services.secretsmanager.secretValuePlaceholder")}
            className="modal-input"
            rows={3}
          />
        </label>
        <label>
          {t("services.secretsmanager.kmsKeyLabel")}
          <input
            value={formKmsKeyId}
            onChange={(e) => setFormKmsKeyId(e.target.value)}
            placeholder={t("services.secretsmanager.kmsKeyPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.secretsmanager.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
