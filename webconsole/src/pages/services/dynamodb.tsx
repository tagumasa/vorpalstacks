/**
 * DynamoDB service page. Lists tables with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { DynamoDBService } from "@/gen/dynamodb_pb";
import { CreateTableInputSchema, DeleteTableInputSchema, KeyType, ScalarAttributeType, BillingMode } from "@/gen/dynamodb_pb";
import { useListKey, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Derived row shape for the DynamoDB table list. */
interface TableRow {
  name: string;
}

/** Column definitions for the DynamoDB table list. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.dynamodb.tableNameHeader"), cell: MonoCell },
];

/** DynamoDB service page with list, create, and delete operations. */
export function DynamoDBPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formPkName, setFormPkName] = useState("pk");
  const [formPkType, setFormPkType] = useState<ScalarAttributeType>(ScalarAttributeType.S);
  const [formSkName, setFormSkName] = useState("");
  const [formSkType, setFormSkType] = useState<ScalarAttributeType>(ScalarAttributeType.S);
  const [formBillingMode, setFormBillingMode] = useState<BillingMode>(BillingMode.PAY_PER_REQUEST);

  const { client, invalidate } = useServiceClient(DynamoDBService);
  const { queryKey } = useListKey("dynamodb");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listTables({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TableRow[] = (data?.tablenames ?? []).map((name) => ({ name }));

  const createMutation = useMutation({
    mutationFn: () =>
      client.createTable(
        create(CreateTableInputSchema, {
          tablename: formName,
          keyschema: [
            { attributename: formPkName, keytype: KeyType.HASH },
            ...(formSkName ? [{ attributename: formSkName, keytype: KeyType.RANGE }] : []),
          ],
          attributedefinitions: [
            { attributename: formPkName, attributetype: formPkType },
            ...(formSkName ? [{ attributename: formSkName, attributetype: formSkType }] : []),
          ],
          billingmode: formBillingMode,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormPkName("pk");
      setFormPkType(ScalarAttributeType.S);
      setFormSkName("");
      setFormSkType(ScalarAttributeType.S);
      setFormBillingMode(BillingMode.PAY_PER_REQUEST);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (tablename: string) =>
      client.deleteTable(create(DeleteTableInputSchema, { tablename })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🗃️"
      title={t("services.dynamodb.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.dynamodb.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.dynamodb.create")}
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
        title={t("services.dynamodb.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formPkName}
      >
        <label>
          {t("services.dynamodb.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.dynamodb.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.dynamodb.partitionKeyLabel")}
          <input
            value={formPkName}
            onChange={(e) => setFormPkName(e.target.value)}
            placeholder={t("services.dynamodb.pkPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.dynamodb.partitionKeyTypeLabel")}
          <select value={formPkType} onChange={(e) => setFormPkType(Number(e.target.value) as ScalarAttributeType)} className="modal-input">
            <option value={ScalarAttributeType.S}>{t("services.dynamodb.attrTypeString")}</option>
            <option value={ScalarAttributeType.N}>{t("services.dynamodb.attrTypeNumber")}</option>
            <option value={ScalarAttributeType.B}>{t("services.dynamodb.attrTypeBinary")}</option>
          </select>
        </label>
        <label>
          {t("services.dynamodb.sortKeyLabel")}
            <input
              value={formSkName}
              onChange={(e) => setFormSkName(e.target.value)}
              placeholder={t("common.optional")}
            className="modal-input"
          />
        </label>
        {formSkName && (
          <label>
            {t("services.dynamodb.sortKeyTypeLabel")}
            <select value={formSkType} onChange={(e) => setFormSkType(Number(e.target.value) as ScalarAttributeType)} className="modal-input">
              <option value={ScalarAttributeType.S}>{t("services.dynamodb.attrTypeString")}</option>
              <option value={ScalarAttributeType.N}>{t("services.dynamodb.attrTypeNumber")}</option>
              <option value={ScalarAttributeType.B}>{t("services.dynamodb.attrTypeBinary")}</option>
            </select>
          </label>
        )}
        <label>
          {t("services.dynamodb.billingModeLabel")}
          <select value={formBillingMode} onChange={(e) => setFormBillingMode(Number(e.target.value) as BillingMode)} className="modal-input">
            <option value={BillingMode.PAY_PER_REQUEST}>{t("services.dynamodb.billingPayPerRequest")}</option>
            <option value={BillingMode.PROVISIONED}>{t("services.dynamodb.billingProvisioned")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.dynamodb.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
