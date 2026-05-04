/**
 * SSM Parameter Store service page. Lists parameters with create/delete CRUD.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import {
  SSMService,
  PutParameterRequestSchema,
  DeleteParameterRequestSchema,
  type ParameterMetadata,
  ParameterType,
  ParameterTier,
} from "@/gen/ssm_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  DateCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the SSM parameter table. */
const getColumns = (t: TFunction): ColumnDef<ParameterMetadata, any>[] => [
  { accessorKey: "name", header: t("services.ssm.nameHeader"), cell: MonoCell },
  {
    accessorKey: "type",
    header: t("services.ssm.typeHeader"),
    cell: ({ getValue }) => {
      const val = getValue() as ParameterType;
      const labels: Record<number, string> = {
        [ParameterType.SECURE_STRING]: t("services.ssm.paramTypeSecureString"),
        [ParameterType.STRING_LIST]: t("services.ssm.paramTypeStringList"),
        [ParameterType.STRING]: t("services.ssm.paramTypeString"),
      };
      return <span className="badge">{labels[val] ?? String(val)}</span>;
    },
  },
  { accessorKey: "version", header: t("services.ssm.versionHeader"), size: 50 },
  {
    accessorKey: "tier",
    header: t("services.ssm.tierHeader"),
    cell: ({ getValue }) => {
      const val = getValue() as ParameterTier;
      const labels: Record<number, string> = {
        [ParameterTier.STANDARD]: t("services.ssm.tierStandard"),
        [ParameterTier.ADVANCED]: t("services.ssm.tierAdvanced"),
        [ParameterTier.INTELLIGENT_TIERING]: t("services.ssm.tierIntelligent"),
      };
      return <span>{labels[val] ?? String(val)}</span>;
    },
    size: 90,
  },
  { accessorKey: "description", header: t("services.ssm.descriptionHeader"), cell: ({ getValue }) => (getValue() as string) || "\u2014" },
  { accessorKey: "datatype", header: t("services.ssm.dataTypeHeader"), cell: ({ getValue }) => (getValue() as string) || "text" },
  { accessorKey: "lastmodifieddate", header: t("services.ssm.lastModifiedHeader"), cell: DateCell },
];

/** Available parameter types for the create form. */
const PARAM_TYPES = [
  { value: ParameterType.STRING, i18nKey: "services.ssm.paramTypeString" },
  { value: ParameterType.SECURE_STRING, i18nKey: "services.ssm.paramTypeSecureString" },
  { value: ParameterType.STRING_LIST, i18nKey: "services.ssm.paramTypeStringList" },
];

/** SSM Parameter Store page with list, create, and delete operations. */
export function SSMPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<ParameterMetadata | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formValue, setFormValue] = useState("");
  const [formType, setFormType] = useState(ParameterType.STRING);
  const [formDesc, setFormDesc] = useState("");

  const { client, invalidate } = useServiceClient(SSMService);
  const { queryKey } = useListKey("ssm");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.describeParameters({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: ParameterMetadata[] = dropEmpty(data?.parameters ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.putParameter(
        create(PutParameterRequestSchema, {
          name: formName,
          value: formValue,
          type: formType,
          description: formDesc,
          tier: ParameterTier.STANDARD,
          overwrite: false,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormValue("");
      setFormType(ParameterType.STRING);
      setFormDesc("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) =>
      client.deleteParameter(create(DeleteParameterRequestSchema, { name })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="📋"
      title={t("services.ssm.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.ssm.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.ssm.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "ssm-items" }}
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
        title={t("services.ssm.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formValue}
      >
        <label>
          {t("services.ssm.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.ssm.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.ssm.valueLabel")}
          <textarea
            value={formValue}
            onChange={(e) => setFormValue(e.target.value)}
            placeholder={t("services.ssm.valuePlaceholder")}
            className="modal-textarea"
            rows={3}
          />
        </label>
        <label>
          {t("services.ssm.typeLabel")}
          <select
            value={formType}
            onChange={(e) => setFormType(Number(e.target.value))}
            className="modal-select"
          >
            {PARAM_TYPES.map((pt) => (
              <option key={pt.value} value={pt.value}>
                {t(pt.i18nKey)}
              </option>
            ))}
          </select>
        </label>
        <label>
          {t("services.ssm.descLabel")}
          <input
            value={formDesc}
            onChange={(e) => setFormDesc(e.target.value)}
            placeholder={t("services.ssm.descPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.ssm.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
