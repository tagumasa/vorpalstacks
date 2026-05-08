/**
 * DynamoDB service page. Lists tables with create/delete operations, a detail
 * panel showing table metadata (DescribeTable), item browser (Scan), and item
 * CRUD (PutItem, DeleteItem).
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  DynamoDBService,
  type AttributeValue,
  type TableDescription,
  type KeySchemaElement,
  type AttributeDefinition,
  KeyType,
  ScalarAttributeType,
  BillingMode,
} from "@/gen/dynamodb_pb";
import {
  CreateTableInputSchema,
  DeleteTableInputSchema,
  ScanInputSchema,
  PutItemInputSchema,
  DeleteItemInputSchema,
  DescribeTableInputSchema,
} from "@/gen/dynamodb_pb";
import { useListKey, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  useServiceClient,
} from "@/components/shared/service-page";
import { JsonViewer } from "@/components/shared/json-viewer";
import { DataTable } from "@/components/shared/data-table";
import { Modal } from "@/components/shared/modal";

/** Derived row shape for the DynamoDB table list. */
interface TableRow {
  name: string;
}

/** Column definitions for the DynamoDB table list. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.dynamodb.tableNameHeader"), cell: MonoCell },
];

/** Render an AttributeValue as a short string. */
function fmtAttr(v: AttributeValue): string {
  if (v.s) return v.s;
  if (v.n) return v.n;
  if (v.bool !== undefined) return String(v.bool);
  if (v.null) return "null";
  if (v.ss?.length) return `[${v.ss.join(", ")}]`;
  if (v.ns?.length) return `[${v.ns.join(", ")}]`;
  if (v.l?.length) return `List(${v.l.length})`;
  if (v.m && Object.keys(v.m).length) return `Map(${Object.keys(v.m).length})`;
  if (v.b !== undefined && v.b.length > 0) return `B(${v.b.length})`;
  if (v.bs?.length) return `BS(${v.bs.length})`;
  return "—";
}

/** Extract PK and optionally SK from an item using the key schema. */
function extractKey(
  item: Record<string, AttributeValue>,
  keySchema: KeySchemaElement[],
): Record<string, AttributeValue> {
  const key: Record<string, AttributeValue> = {};
  for (const ks of keySchema) {
    const v = item[ks.attributename];
    if (v) key[ks.attributename] = v;
  }
  return key;
}

/** Key schema summary row. */
function KeySchemaDisplay({ elements, t }: { elements: KeySchemaElement[]; t: TFunction }) {
  return (
    <div className="detail-field">
      <span className="detail-label">{t("services.dynamodb.detail.keySchema")}</span>
      <span>
        {elements.map((e, i) => (
          <span key={i}>
            {i > 0 && ", "}
            <span className="cell-mono">{e.attributename}</span>{" "}
            <span className="badge">{e.keytype === KeyType.HASH ? t("services.dynamodb.detail.hashKey") : t("services.dynamodb.detail.rangeKey")}</span>
          </span>
        ))}
      </span>
    </div>
  );
}

/** Attribute definitions display. */
function AttrDefsDisplay({ defs, t }: { defs: AttributeDefinition[]; t: TFunction }) {
  return (
    <div className="detail-field">
      <span className="detail-label">{t("services.dynamodb.detail.attributes")}</span>
      <span>
        {defs.map((d, i) => (
          <span key={i}>
            {i > 0 && ", "}
            <span className="cell-mono">{d.attributename}</span>{" "}
            <span className="badge">{d.attributetype}</span>
          </span>
        ))}
      </span>
    </div>
  );
}

/** Item row shape for the items sub-table. */
interface ItemRow {
  keyJson: string;
  pk: string;
  sk: string;
  item: Record<string, AttributeValue>;
}

/** DynamoDB table detail panel with DescribeTable metadata + Scan items. */
function DynamoDBDetail({ item: tableRow }: { item: TableRow }) {
  const { t } = useTranslation();
  const { client } = useServiceClient(DynamoDBService);
  const queryClient = useQueryClient();

  const [showAddItem, setShowAddItem] = useState(false);
  const [itemJson, setItemJson] = useState("");
  const [jsonError, setJsonError] = useState("");
  const [selectedItem, setSelectedItem] = useState<ItemRow | null>(null);
  const [showDeleteItem, setShowDeleteItem] = useState(false);

  const { data: descData } = useQuery({
    queryKey: ["dynamodb", "describe", tableRow.name],
    queryFn: () => client.describeTable(create(DescribeTableInputSchema, { tablename: tableRow.name })),
    refetchInterval: REFETCH_INTERVAL,
  });

  const table: TableDescription | undefined = descData?.table;
  const keySchema = table?.keyschema ?? [];
  const attrDefs = table?.attributedefinitions ?? [];
  const pkName = keySchema.find((k) => k.keytype === KeyType.HASH)?.attributename ?? "";
  const skName = keySchema.find((k) => k.keytype === KeyType.RANGE)?.attributename ?? "";

  const { data: scanData } = useQuery({
    queryKey: ["dynamodb", "items", tableRow.name],
    queryFn: () => client.scan(create(ScanInputSchema, { tablename: tableRow.name, limit: 100 })),
    refetchInterval: REFETCH_INTERVAL,
    enabled: !!keySchema.length,
  });

  const items: ItemRow[] = (scanData?.items ?? []).map((entry) => {
    const item = entry.value as Record<string, AttributeValue>;
    const key = extractKey(item, keySchema);
    return {
      keyJson: JSON.stringify(key),
      pk: item[pkName] ? fmtAttr(item[pkName]) : "—",
      sk: skName && item[skName] ? fmtAttr(item[skName]) : "",
      item,
    };
  });

  const putItemMutation = useMutation({
    mutationFn: (parsed: Record<string, AttributeValue>) =>
      client.putItem(create(PutItemInputSchema, { tablename: tableRow.name, item: parsed })),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dynamodb", "items", tableRow.name] });
      setShowAddItem(false);
      setItemJson("");
      setJsonError("");
    },
    onError: (err: Error) => setJsonError(String(err)),
  });

  const deleteItemMutation = useMutation({
    mutationFn: (key: Record<string, AttributeValue>) =>
      client.deleteItem(create(DeleteItemInputSchema, { tablename: tableRow.name, key })),
    onSuccess: () => {
      queryClient.invalidateQueries({ queryKey: ["dynamodb", "items", tableRow.name] });
      setShowDeleteItem(false);
      setSelectedItem(null);
    },
  });

  const handleAddItem = () => {
    try {
      const parsed = JSON.parse(itemJson);
      setJsonError("");
      putItemMutation.mutate(parsed);
    } catch {
      setJsonError(t("services.dynamodb.detail.invalidJson"));
    }
  };

  const itemColumns: ColumnDef<ItemRow, any>[] = [
    { accessorKey: "pk", header: pkName || "PK", cell: MonoCell },
    ...(skName ? [{ accessorKey: "sk" as const, header: skName, cell: MonoCell }] : []),
  ];

  return (
    <div className="detail-body">
      <section className="detail-section">
        <h3>{t("common.general")}</h3>
        {keySchema.length > 0 && <KeySchemaDisplay elements={keySchema} t={t} />}
        {attrDefs.length > 0 && <AttrDefsDisplay defs={attrDefs} t={t} />}
        <div className="detail-field">
          <span className="detail-label">{t("services.dynamodb.detail.itemCount")}</span>
          <span>{table?.itemcount ?? 0}</span>
        </div>
        <div className="detail-field">
          <span className="detail-label">{t("services.dynamodb.detail.tableSize")}</span>
          <span>{table?.tablesizebytes ? Number(table.tablesizebytes).toLocaleString() : "0"}</span>
        </div>
        <div className="detail-field">
          <span className="detail-label">{t("services.dynamodb.detail.status")}</span>
          <span>{table?.tablestatus || "—"}</span>
        </div>
        <div className="detail-field">
          <span className="detail-label">{t("services.dynamodb.detail.creationDate")}</span>
          <span>{table?.creationdatetime || "—"}</span>
        </div>
      </section>

      <section className="detail-section">
        <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
          <h3>{t("services.dynamodb.detail.items")} ({items.length})</h3>
          <button className="btn btn-primary btn-sm" onClick={() => setShowAddItem(true)}>
            {t("services.dynamodb.detail.addJson")}
          </button>
        </div>
        {items.length > 0 ? (
          <DataTable
            columns={itemColumns}
            data={items}
            getRowId={(row) => row.keyJson}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.keyJson}
          />
        ) : (
          <div className="empty-state">{t("services.dynamodb.detail.noItems")}</div>
        )}
        {selectedItem && (
          <div style={{ marginTop: 8 }}>
            <div style={{ display: "flex", justifyContent: "space-between", alignItems: "center" }}>
              <h4>{t("services.dynamodb.detail.value")}</h4>
              <button className="btn btn-danger btn-sm" onClick={() => setShowDeleteItem(true)}>
                {t("services.dynamodb.detail.deleteItem")}
              </button>
            </div>
            <JsonViewer data={selectedItem.item} />
          </div>
        )}
      </section>

      <Modal open={showAddItem} onClose={() => { setShowAddItem(false); setJsonError(""); }}>
        <h2>{t("services.dynamodb.detail.addJson")}</h2>
        {jsonError && <div className="modal-error">{jsonError}</div>}
        {putItemMutation.error && <div className="modal-error">{String(putItemMutation.error)}</div>}
        <textarea
          value={itemJson}
          onChange={(e) => setItemJson(e.target.value)}
          placeholder={t("services.dynamodb.detail.jsonPlaceholder")}
          className="modal-input"
          rows={8}
          style={{ fontFamily: "monospace", fontSize: 12 }}
        />
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => { setShowAddItem(false); setJsonError(""); }}>
            {t("common.cancel")}
          </button>
          <button className="btn btn-primary" disabled={!itemJson || putItemMutation.isPending} onClick={handleAddItem}>
            {putItemMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <Modal open={showDeleteItem && !!selectedItem} onClose={() => setShowDeleteItem(false)}>
        <h2>{t("services.dynamodb.detail.deleteItem")}</h2>
        <p>{t("services.dynamodb.detail.confirmDeleteItem")}</p>
        {selectedItem && <JsonViewer data={selectedItem.item} />}
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowDeleteItem(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-danger"
            disabled={deleteItemMutation.isPending}
            onClick={() => {
              if (selectedItem) {
                deleteItemMutation.mutate(extractKey(selectedItem.item, keySchema));
              }
            }}
          >
            {deleteItemMutation.isPending ? t("common.deleting") : t("common.delete")}
          </button>
        </div>
      </Modal>
    </div>
  );
}

/** DynamoDB service page with list, create, delete, detail, and item CRUD. */
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
        DetailComponent={DynamoDBDetail}
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
