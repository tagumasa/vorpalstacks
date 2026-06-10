/**
 * DynamoDB service page — 3-panel inspector layout matching S3 pattern.
 *
 * Panel 1 (toolbar): Breadcrumb navigation  DynamoDB / Tables  or  DynamoDB / tableName
 * Panel 2 (table):   Table list OR item list with checkbox multi-select
 * Panel 3 (detail):  Item detail with JSON + Attributes tabs (drag-split bottom panel)
 */
import { useState, useCallback, useMemo } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  DynamoDBService,
  type AttributeValue,
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
import {
  type AttrRow, type AttrType,
  fmtAttr, extractKey, buildAV, rowToAV, avToRow, rowToPlain, plainToRow,
  jsonToAV, avToPlain, typePlaceholder,
  TYPE_LABEL_KEYS, ATTR_TYPES,
} from "@/lib/dynamodb-attrs";
import { useListKey, usePaginatedList } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  useServiceClient,
} from "@/components/shared/service-page";
import {
  checkboxColumn,
  Breadcrumb,
  SelectionBadge,
  useSelection,
} from "@/components/shared/inspector";
import { JsonViewer } from "@/components/shared/json-viewer";
import { DataTable } from "@/components/shared/data-table";
import { Modal } from "@/components/shared/modal";
import { Splitter } from "@/components/shared/splitter";

// ─── Row types ──────────────────────────────────────────────────

/** Derived row shape for the DynamoDB table list. */
interface TableRow {
  name: string;
}

/** Item row shape for the items sub-table. */
interface ItemRow {
  keyJson: string;
  pk: string;
  sk: string;
  item: Record<string, AttributeValue>;
}

// ─── View state ──────────────────────────────────────────────────

type ViewState =
  | { type: "tables" }
  | { type: "items"; tableName: string };

// ─── Detail panel tab ────────────────────────────────────────────

type DetailTab = "json" | "edit";

// ─── PutItem modal tab ───────────────────────────────────────────

type PutItemTab = "structured" | "json";

// ─── Table list columns ─────────────────────────────────────────

const tableColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.dynamodb.tableNameHeader"), cell: MonoCell },
];

// ─── DynamoDB Page ──────────────────────────────────────────────

export function DynamoDBPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(DynamoDBService);
  const queryClient = useQueryClient();
  const { queryKey } = useListKey("dynamodb");

  // ── View state ───────────────────────────────────────────────
  const [view, setView] = useState<ViewState>({ type: "tables" });

  // ── Table list selection ──────────────────────────────────────
  const {
    selected: selectedTableNames,
    toggle: toggleTable,
    toggleAll: toggleAllTables,
    clear: clearTableSelection,
  } = useSelection<string>();

  // ── Item selection (single click for detail) ──────────────────
  const [selectedItem, setSelectedItem] = useState<ItemRow | null>(null);

  // ── Item checkbox multi-select ────────────────────────────────
  const {
    selected: selectedItemKeys,
    toggle: toggleItem,
    toggleAll: toggleAllItems,
    clear: clearItemSelection,
  } = useSelection<string>();

  // ── Detail panel tab ──────────────────────────────────────────
  const [detailTab, setDetailTab] = useState<DetailTab>("edit");

  // ── Modals ────────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDeleteTables, setShowDeleteTables] = useState(false);
  const [showDeleteItems, setShowDeleteItems] = useState(false);
  const [showAddItem, setShowAddItem] = useState(false);
  const [showDeleteItem, setShowDeleteItem] = useState(false);
  const [isCreating, setIsCreating] = useState(false);

  // ── Create table form ─────────────────────────────────────────
  const [formName, setFormName] = useState("");
  const [formPkName, setFormPkName] = useState("pk");
  const [formPkType, setFormPkType] = useState<ScalarAttributeType>(ScalarAttributeType.S);
  const [formSkName, setFormSkName] = useState("");
  const [formSkType, setFormSkType] = useState<ScalarAttributeType>(ScalarAttributeType.S);
  const [formBillingMode, setFormBillingMode] = useState<BillingMode>(BillingMode.PAY_PER_REQUEST);

  // ── PutItem modal state ───────────────────────────────────────
  const [putItemTab, setPutItemTab] = useState<PutItemTab>("structured");
  const [itemJson, setItemJson] = useState("");
  const [jsonError, setJsonError] = useState("");
  const [attrRows, setAttrRows] = useState<AttrRow[]>([]);

  // ── Scan pagination ───────────────────────────────────────────
  const [lastEvaluatedKey, setLastEvaluatedKey] = useState<Record<string, AttributeValue> | undefined>(undefined);
  const [accumulatedItems, setAccumulatedItems] = useState<ItemRow[]>([]);

  // ── Batch delete result ───────────────────────────────────────
  const [batchResult, setBatchResult] = useState<string | null>(null);

  // ── Table list query ──────────────────────────────────────────
  const tableList = usePaginatedList<string, Awaited<ReturnType<typeof client.listTables>>>({
    queryKeyBase: queryKey, fetchPage: (token) => client.listTables({ exclusivestarttablename: token || undefined }), getItems: (r) => r.tablenames ?? [], getNextToken: (r) => r.lastevaluatedtablename ?? "",
  });

  const tables: TableRow[] = tableList.items.map((name) => ({ name }));

  // ── DescribeTable query (items view) ──────────────────────────
  const activeTableName = view.type === "items" ? view.tableName : "";

  const { data: descData } = useQuery({
    queryKey: ["dynamodb", "describe", activeTableName],
    queryFn: () => client.describeTable(create(DescribeTableInputSchema, { tablename: activeTableName })),
    enabled: view.type === "items",
  });

  const keySchema = descData?.table?.keyschema ?? [];
  const attrDefs = descData?.table?.attributedefinitions ?? [];
  const pkName = keySchema.find((k) => k.keytype === KeyType.HASH)?.attributename ?? "";
  const skName = keySchema.find((k) => k.keytype === KeyType.RANGE)?.attributename ?? "";

  const scalarToAttrType = (sat: ScalarAttributeType): AttrType =>
    sat === ScalarAttributeType.N ? "N" : sat === ScalarAttributeType.B ? "B" : "S";
  const pkType = scalarToAttrType(attrDefs.find((a) => a.attributename === pkName)?.attributetype ?? ScalarAttributeType.S);
  const skType = scalarToAttrType(attrDefs.find((a) => a.attributename === skName)?.attributetype ?? ScalarAttributeType.S);

  // ── Scan query (items view) ───────────────────────────────────
  const itemsQueryKey = ["dynamodb", "items", activeTableName];

  const { data: scanData, isFetching: scanFetching } = useQuery({
    queryKey: [...itemsQueryKey, JSON.stringify(lastEvaluatedKey)],
    queryFn: () =>
      client.scan(create(ScanInputSchema, {
        tablename: activeTableName,
        limit: 100,
        exclusivestartkey: lastEvaluatedKey,
      })),
    enabled: view.type === "items" && !!keySchema.length,
  });

  const currentItems: ItemRow[] = (scanData?.items ?? []).map((entry) => {
    const item = entry.value as Record<string, AttributeValue>;
    const key = extractKey(item, keySchema);
    return {
      keyJson: JSON.stringify(key),
      pk: item[pkName] ? fmtAttr(item[pkName]) : "—",
      sk: skName && item[skName] ? fmtAttr(item[skName]) : "",
      item,
    };
  });

  const allItems = lastEvaluatedKey ? [...accumulatedItems, ...currentItems] : currentItems;
  const nextEvalKey = scanData?.lastevaluatedkey && Object.keys(scanData.lastevaluatedkey).length > 0
    ? scanData.lastevaluatedkey
    : undefined;

  // ── Item columns (dynamic — includes all attributes) ──────────
  const extraAttrs = useMemo(() => {
    const s = new Set<string>();
    for (const row of allItems) {
      for (const k of Object.keys(row.item)) {
        if (k !== pkName && k !== skName) s.add(k);
      }
    }
    return Array.from(s).sort();
  }, [allItems, pkName, skName]);

  const itemColumns: ColumnDef<ItemRow, any>[] = useMemo(() => [
    { accessorKey: "pk", header: pkName || "PK", cell: MonoCell, size: 120 },
    ...(skName ? [{ accessorKey: "sk" as const, header: skName, cell: MonoCell, size: 120 }] : []),
    ...extraAttrs.map(attr => ({
      accessorKey: attr,
      header: attr,
      cell: ({ row }: { row: { original: ItemRow } }) => {
        const av = row.original.item[attr];
        return av ? <span className="cell-mono">{fmtAttr(av)}</span> : <span className="attr-mute-inline">—</span>;
      },
      size: 120,
    })),
  ], [pkName, skName, extraAttrs]);

  // ── Navigation helpers ────────────────────────────────────────

  const navigateToTable = useCallback((name: string) => {
    setView({ type: "items", tableName: name });
    setSelectedItem(null);
    clearItemSelection();
    setDetailTab("edit");
    setLastEvaluatedKey(undefined);
    setAccumulatedItems([]);
    setBatchResult(null);
    setJsonError("");
  }, []);

  const navigateToTables = useCallback(() => {
    setView({ type: "tables" });
    clearTableSelection();
    setSelectedItem(null);
    clearItemSelection();
    setBatchResult(null);
  }, []);

  const invalidateItems = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: itemsQueryKey });
    setLastEvaluatedKey(undefined);
    setAccumulatedItems([]);
  }, [queryClient, itemsQueryKey]);

  // ── Checkbox toggle helpers ───────────────────────────────────

  /** Derived IDs for checkbox column toggle-all helpers. */
  const allTableIds = tables.map((tbl) => tbl.name);
  const allItemIds = allItems.map((i) => i.keyJson);

  // ── Table mutations ───────────────────────────────────────────

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
      tableList.invalidate();
      setShowCreate(false);
      setFormName("");
      setFormPkName("pk");
      setFormPkType(ScalarAttributeType.S);
      setFormSkName("");
      setFormSkType(ScalarAttributeType.S);
      setFormBillingMode(BillingMode.PAY_PER_REQUEST);
    },
  });

  const deleteTableMutation = useMutation({
    mutationFn: (tablename: string) =>
      client.deleteTable(create(DeleteTableInputSchema, { tablename })),
    onSuccess: () => {
      tableList.invalidate();
    },
  });

  // ── Item mutations ────────────────────────────────────────────

  const putItemMutation = useMutation({
    mutationFn: (parsed: Record<string, AttributeValue>) =>
      client.putItem(create(PutItemInputSchema, { tablename: activeTableName, item: parsed })),
    onSuccess: () => {
      invalidateItems();
      if (showAddItem) {
        setShowAddItem(false);
        setItemJson("");
        setAttrRows([]);
      }
      if (isCreating) {
        setIsCreating(false);
        setAttrRows([]);
      }
      setJsonError("");
    },
    onError: (err: Error) => setJsonError(String(err)),
  });

  const deleteItemMutation = useMutation({
    mutationFn: (key: Record<string, AttributeValue>) =>
      client.deleteItem(create(DeleteItemInputSchema, { tablename: activeTableName, key })),
    onSuccess: () => {
      invalidateItems();
      setShowDeleteItem(false);
      setSelectedItem(null);
      if (selectedItem) toggleItem(selectedItem.keyJson);
    },
  });

  // ── Batch delete items (Promise.allSettled) ───────────────────

  const batchDeleteItems = async () => {
    const keysToDelete = allItems
      .filter((i) => selectedItemKeys.has(i.keyJson))
      .map((i) => extractKey(i.item, keySchema));

    const results = await Promise.allSettled(
      keysToDelete.map((key) =>
        client.deleteItem(create(DeleteItemInputSchema, { tablename: activeTableName, key }))
      ),
    );

    const succeeded = results.filter((r) => r.status === "fulfilled").length;
    const failed = results.filter((r) => r.status === "rejected").length;

    setBatchResult(
      failed === 0
        ? t("services.dynamodb.batchDeleteResult", { count: succeeded })
        : t("services.dynamodb.batchDeleteResult", { count: succeeded }) +
          ` (${failed} ${t("services.dynamodb.failed")})`,
    );

    invalidateItems();
    clearItemSelection();
    setSelectedItem(null);
    setShowDeleteItems(false);

    // Clear result after 5s
    setTimeout(() => setBatchResult(null), 5000);
  };

  // ── PutItem handlers ──────────────────────────────────────────

  const handleAddItemJson = () => {
    try {
      const parsed = JSON.parse(itemJson);
      if (typeof parsed !== "object" || Array.isArray(parsed)) {
        setJsonError(t("services.dynamodb.detail.invalidJson"));
        return;
      }
      const avMap: Record<string, AttributeValue> = {};
      for (const [k, v] of Object.entries(parsed as Record<string, unknown>)) {
        avMap[k] = jsonToAV(v);
      }
      setJsonError("");
      putItemMutation.mutate(avMap);
    } catch {
      setJsonError(t("services.dynamodb.detail.invalidJson"));
    }
  };

  const handleAddItemStructured = () => {
    const item: Record<string, AttributeValue> = {};
    for (const row of attrRows) {
      if (!row.name.trim()) continue;
      if (row.type === "NULL") {
        item[row.name] = buildAV("NULL", "");
        continue;
      }
      if (row.type === "L" || row.type === "M") {
        if (!(row.children?.length)) continue;
        try {
          item[row.name] = rowToAV(row);
        } catch (e) {
          setJsonError(`${row.name}: ${e instanceof Error ? e.message : String(e)}`);
          return;
        }
        continue;
      }
      if (!row.value.trim()) continue;
      try {
        item[row.name] = rowToAV(row);
      } catch (e) {
        setJsonError(`${row.name}: ${e instanceof Error ? e.message : String(e)}`);
        return;
      }
    }
    if (Object.keys(item).length === 0) {
      setJsonError(t("services.dynamodb.detail.invalidJson"));
      return;
    }
    setJsonError("");
    putItemMutation.mutate(item);
  };

  const openAddItemModal = () => {
    const keyRows: AttrRow[] = [];
    if (pkName) keyRows.push({ name: pkName, type: pkType, value: "" });
    if (skName) keyRows.push({ name: skName, type: skType, value: "" });
    setAttrRows(keyRows);
    setItemJson("");
    setJsonError("");
    setSelectedItem(null);
    setIsCreating(true);
    setDetailTab("edit");
  };

  const switchToStructured = () => {
    try {
      const parsed = JSON.parse(itemJson);
      if (typeof parsed === "object" && !Array.isArray(parsed)) {
        const rows: AttrRow[] = Object.entries(parsed as Record<string, unknown>).map(([name, v]) => plainToRow(name, v));
        setAttrRows(rows);
      }
    } catch { /* keep existing attrRows */ }
    setPutItemTab("structured");
  };

  const switchToJson = () => {
    const plainObj: Record<string, unknown> = {};
    for (const row of attrRows) {
      if (!row.name.trim()) continue;
      plainObj[row.name] = rowToPlain(row);
    }
    if (Object.keys(plainObj).length > 0) setItemJson(JSON.stringify(plainObj, null, 2));
    setPutItemTab("json");
  };

  const addAttrRow = () => {
    setAttrRows((prev) => [...prev, { name: "", type: "S" as AttrType, value: "" }]);
  };

  const removeAttrRow = (index: number) => {
    setAttrRows((prev) => prev.filter((_, i) => i !== index));
  };

  const updateAttrRow = (index: number, field: keyof AttrRow, value: string | AttrType) => {
    setAttrRows((prev) =>
      prev.map((row, i) => {
        if (i !== index) return row;
        const updated = { ...row, [field]: value };
        if (field === "type") {
          if (value === "L") updated.children = [newChild("L")];
          else if (value === "M") updated.children = [newChild("M")];
          else updated.children = undefined;
        }
        return updated;
      }),
    );
  };

  const newChild = (parentType: "L" | "M"): AttrRow =>
    parentType === "L" ? { name: "", type: "S", value: "" } : { name: "", type: "S", value: "" };

  const addChildRow = (parentIndex: number) => {
    setAttrRows((prev) =>
      prev.map((row, i) => {
        if (i !== parentIndex) return row;
        const parentType = row.type as "L" | "M";
        return { ...row, children: [...(row.children ?? []), newChild(parentType)] };
      }),
    );
  };

  const removeChildRow = (parentIndex: number, childIndex: number) => {
    setAttrRows((prev) =>
      prev.map((row, i) => {
        if (i !== parentIndex) return row;
        return { ...row, children: (row.children ?? []).filter((_, ci) => ci !== childIndex) };
      }),
    );
  };

  const updateChildRow = (parentIndex: number, childIndex: number, field: keyof AttrRow, value: string | AttrType) => {
    setAttrRows((prev) =>
      prev.map((row, i) => {
        if (i !== parentIndex) return row;
        const children = (row.children ?? []).map((child, ci) => {
          if (ci !== childIndex) return child;
          const updated = { ...child, [field]: value };
          if (field === "type") {
            if (value === "L") updated.children = [newChild("L")];
            else if (value === "M") updated.children = [newChild("M")];
            else updated.children = undefined;
          }
          return updated;
        });
        return { ...row, children };
      }),
    );
  };

  // ── Breadcrumb ────────────────────────────────────────────────

  const breadcrumbParts = view.type === "items"
    ? [
        { label: t("services.dynamodb.title") },
        { label: t("services.dynamodb.tablesLabel"), onClick: navigateToTables },
        { label: view.tableName },
      ]
    : [
        { label: t("services.dynamodb.title") },
        { label: t("services.dynamodb.tablesLabel") },
      ];

  // ── Detail panel ──────────────────────────────────────────────

  const renderDetailPanel = () => {
    const isEditing = isCreating || selectedItem;

    if (!isEditing) {
      return (
        <div className="detail-panel-empty">
          {t("services.dynamodb.noItemSelected")}
        </div>
      );
    }

    const cellInputStyle = (mono?: boolean): React.CSSProperties => ({
      flex: 1, height: 28, fontSize: 11, padding: "2px 6px",
      border: "1px solid var(--border-dim)", background: "var(--bg-tertiary)",
      color: "var(--text-primary)", borderRadius: 3,
      ...(mono ? { fontFamily: "monospace" } : {}),
    });

    const isKeyAttr = (row: AttrRow) =>
      row.name !== "" && (row.name === pkName || row.name === skName);

    const delBtnStyle: React.CSSProperties = {
      padding: "2px 8px", fontSize: 12, height: 28, lineHeight: 1, cursor: "pointer",
      border: "1px solid #dc2626", background: "transparent", color: "#dc2626",
      borderRadius: 3, flexShrink: 0,
    };

    const renderChildRow = (parentIdx: number, child: AttrRow, ci: number, parentType: "L" | "M") => (
      <div key={ci} className="attr-row attr-row--child">
        {parentType === "M" ? (
          <input
            value={child.name}
            onChange={(e) => updateChildRow(parentIdx, ci, "name", e.target.value)}
            placeholder="key"
            className="attr-input attr-input--narrow"
          />
        ) : (
          <span className="attr-index">{ci}</span>
        )}
        <select
          value={child.type}
          onChange={(e) => updateChildRow(parentIdx, ci, "type", e.target.value as AttrType)}
          className="attr-input attr-input--value"
        >
          {ATTR_TYPES.map((at) => (
            <option key={at} value={at}>{t(`services.dynamodb.${TYPE_LABEL_KEYS[at]}`)}</option>
          ))}
        </select>
        {child.type === "BOOL" ? (
          <input type="checkbox" checked={child.value === "true"} onChange={(e) => updateChildRow(parentIdx, ci, "value", e.target.checked ? "true" : "false")} className="attr-checkbox" />
        ) : child.type === "NULL" ? (
          <span className="attr-null">null</span>
        ) : (child.type === "L" || child.type === "M") ? (
          <span className="attr-count">{child.type === "L" ? `[${child.children?.length ?? 0} items]` : `{${child.children?.length ?? 0} entries}`}</span>
        ) : (
          <input value={child.value} onChange={(e) => updateChildRow(parentIdx, ci, "value", e.target.value)} placeholder={typePlaceholder(child.type)} style={cellInputStyle(false)} />
        )}
        <button onClick={() => removeChildRow(parentIdx, ci)} title={t("common.delete")} style={delBtnStyle}>✕</button>
      </div>
    );

    const renderAttrRow = (row: AttrRow, i: number) => {
      const isKey = isKeyAttr(row);
      const isListOrMap = row.type === "L" || row.type === "M";
      return (
        <div key={i} style={{ marginBottom: isListOrMap ? 6 : 0 }}>
          <div className="attr-row">
            {isKey ? (
              <span className="cell-mono attr-name" title={t("services.dynamodb.keyAttribute")}>
                {row.name}
              </span>
            ) : (
              <input
                value={row.name}
                onChange={(e) => updateAttrRow(i, "name", e.target.value)}
                placeholder={t("services.dynamodb.attributeName")}
                className="attr-input"
              />
            )}
            <select
              value={row.type}
              onChange={(e) => updateAttrRow(i, "type", e.target.value as AttrType)}
              disabled={isKey}
              className="attr-input attr-input--wide"
            >
              {ATTR_TYPES.map((at) => (
                <option key={at} value={at}>{t(`services.dynamodb.${TYPE_LABEL_KEYS[at]}`)}</option>
              ))}
            </select>
            {isListOrMap ? (
              <span className="attr-count">
                {row.type === "L" ? `[${row.children?.length ?? 0} items]` : `{${row.children?.length ?? 0} entries}`}
              </span>
            ) : row.type === "BOOL" ? (
              <input type="checkbox" checked={row.value === "true"} onChange={(e) => updateAttrRow(i, "value", e.target.checked ? "true" : "false")} className="attr-checkbox" />
            ) : row.type === "NULL" ? (
              <span className="attr-null">null</span>
            ) : (
              <input value={row.value} onChange={(e) => updateAttrRow(i, "value", e.target.value)} placeholder={typePlaceholder(row.type)} style={cellInputStyle(false)} />
            )}
            {!isKey && <button onClick={() => removeAttrRow(i)} title={t("common.delete")} style={delBtnStyle}>✕</button>}
          </div>
          {isListOrMap && (
            <div className="attr-nested-border">
              {(row.children ?? []).map((child, ci) => renderChildRow(i, child, ci, row.type as "L" | "M"))}
              <button onClick={() => addChildRow(i)} className="attr-add-btn">
                + {row.type === "L" ? t("services.dynamodb.addItem") || "item" : t("services.dynamodb.addAttribute") || "entry"}
              </button>
            </div>
          )}
        </div>
      );
    };

    const renderEditForm = () => {
      const keyRows: number[] = [];
      const nonKeyRows: number[] = [];
      attrRows.forEach((row, i) => {
        if (isKeyAttr(row)) keyRows.push(i);
        else nonKeyRows.push(i);
      });

      return (
        <div className="attr-section">
          {keyRows.length > 0 && (
            <div className="attr-section-head">
              {keyRows.map((idx) => attrRows[idx] ? renderAttrRow(attrRows[idx], idx) : null)}
            </div>
          )}
          {nonKeyRows.map((idx) => attrRows[idx] ? renderAttrRow(attrRows[idx], idx) : null)}
          <button onClick={addAttrRow} className="attr-section-add">
            + {t("services.dynamodb.addAttribute")}
          </button>
        </div>
      );
    };

    return (
      <div className="detail-panel-content">
        <div className="detail-header">
          <span className="detail-title">
            {isCreating
              ? t("services.dynamodb.addItem")
              : `🔑 ${selectedItem!.pk}${selectedItem!.sk ? ` / ${selectedItem!.sk}` : ""}`}
          </span>
          <div className="detail-tabs">
            <button
              className={`detail-tab ${detailTab === "edit" ? "active" : ""}`}
              onClick={() => {
                if (detailTab === "json") {
                  try {
                    const parsed = JSON.parse(itemJson);
                    if (typeof parsed === "object" && !Array.isArray(parsed)) {
                      setAttrRows(Object.entries(parsed as Record<string, unknown>).map(([name, v]) => plainToRow(name, v)));
                    }
                  } catch { /* keep existing attrRows */ }
                }
                setDetailTab("edit");
              }}
            >
              {t("services.dynamodb.structuredForm")}
            </button>
            <button
              className={`detail-tab ${detailTab === "json" ? "active" : ""}`}
              onClick={() => {
                if (detailTab === "edit") {
                  const plainObj: Record<string, unknown> = {};
                  for (const row of attrRows) {
                    if (!row.name.trim()) continue;
                    plainObj[row.name] = rowToPlain(row);
                  }
                  if (Object.keys(plainObj).length > 0) setItemJson(JSON.stringify(plainObj, null, 2));
                }
                setDetailTab("json");
              }}
            >
              {t("services.dynamodb.tabJson")}
            </button>
          </div>
          <div className="detail-actions">
            {detailTab === "edit" && (
              <button
                className="btn btn-primary btn-sm"
                disabled={attrRows.length === 0 || attrRows.some((r) => !r.name.trim() || (r.type !== "NULL" && r.type !== "L" && r.type !== "M" && !r.value.trim())) || putItemMutation.isPending}
                onClick={handleAddItemStructured}
              >
                {putItemMutation.isPending ? t("services.dynamodb.saving") : isCreating ? t("services.dynamodb.createItem") : t("services.dynamodb.saveChanges")}
              </button>
            )}
            {detailTab === "json" && (
              <button
                className="btn btn-primary btn-sm"
                disabled={!itemJson || putItemMutation.isPending}
                onClick={handleAddItemJson}
              >
                {putItemMutation.isPending ? t("services.dynamodb.saving") : isCreating ? t("services.dynamodb.createItem") : t("services.dynamodb.saveChanges")}
              </button>
            )}
            {!isCreating && (
              <button
                className="btn btn-danger btn-sm"
                onClick={() => setShowDeleteItem(true)}
              >
                {t("services.dynamodb.detail.deleteItem")}
              </button>
            )}
            {isCreating && (
              <button
                className="btn btn-ghost btn-sm"
                onClick={() => { setIsCreating(false); setAttrRows([]); setJsonError(""); }}
              >
                {t("common.cancel")}
              </button>
            )}
          </div>
        </div>
        {jsonError && <div className="modal-error json-error-inline">{jsonError}</div>}
        <div className="detail-body">
          {detailTab === "edit"
            ? renderEditForm()
            : (
              <textarea
                value={itemJson}
                onChange={(e) => setItemJson(e.target.value)}
                className="json-editor"
              />
            )}
        </div>
      </div>
    );
  };

  // ── Toolbar actions ───────────────────────────────────────────

  const renderActions = () => {
    if (view.type === "tables") {
      return (
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.dynamodb.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedTableNames.size === 0}
            onClick={() => setShowDeleteTables(true)}
          >
            {t("common.delete")}
            {selectedTableNames.size > 0 && (
              <span className="badge-count">({selectedTableNames.size})</span>
            )}
          </button>
        </>
      );
    }
    return (
      <>
        <button className="btn btn-primary btn-sm" onClick={openAddItemModal}>
          {t("services.dynamodb.addItem")}
        </button>
        <button
          className="btn btn-danger btn-sm"
          disabled={selectedItemKeys.size === 0}
          onClick={() => setShowDeleteItems(true)}
        >
          {t("services.dynamodb.deleteSelected")}
          {selectedItemKeys.size > 0 && (
            <span className="badge-count">({selectedItemKeys.size})</span>
          )}
        </button>
      </>
    );
  };

  // ── Row click handlers ────────────────────────────────────────

  const handleTableRowClick = (row: TableRow) => {
    navigateToTable(row.name);
  };

  const handleItemRowClick = (row: ItemRow) => {
    setSelectedItem(row);
    setIsCreating(false);
    setShowAddItem(false);
    const raw: AttrRow[] = Object.entries(row.item).map(([name, av]) => avToRow(name, av));
    const sorted = raw.sort((a, b) => {
      if (a.name === pkName) return -1;
      if (b.name === pkName) return 1;
      if (a.name === skName) return -1;
      if (b.name === skName) return 1;
      return a.name.localeCompare(b.name);
    });
    setAttrRows(sorted);
    const plainObj: Record<string, unknown> = {};
    for (const [k, av] of Object.entries(row.item)) plainObj[k] = avToPlain(av);
    setItemJson(JSON.stringify(plainObj, null, 2));
    setJsonError("");
    setDetailTab("edit");
  };

  const handleLoadMore = () => {
    setAccumulatedItems(allItems);
    setLastEvaluatedKey(nextEvalKey);
  };

  // ── Render ────────────────────────────────────────────────────

  return (
    <ServicePageLayout
      icon="🗃️"
      title={t("services.dynamodb.title")}
      isLoading={tableList.isLoading && view.type === "tables"}
      error={tableList.error}
      count={view.type === "tables" ? tables.length : undefined}
      countLabel={t("services.dynamodb.countLabel")}
      actions={renderActions()}
    >
      {/* Breadcrumb toolbar */}
      <div className="inspector-toolbar">
         <Breadcrumb parts={breadcrumbParts} />
         <div className="toolbar-selection-info">
           <SelectionBadge count={view.type === "tables" ? selectedTableNames.size : selectedItemKeys.size} label={t("services.dynamodb.selectedCount", { count: view.type === "tables" ? selectedTableNames.size : selectedItemKeys.size })} />
           {batchResult && (
             <span className="selection-count selection-count--green">
               {batchResult}
             </span>
           )}
         </div>
      </div>

      {/* Main content area */}
      {view.type === "tables" ? (
        /* Table list — full width with checkbox */
        <DataTable
          columns={[
            checkboxColumn<TableRow>(selectedTableNames, toggleTable, () => toggleAllTables(allTableIds), allTableIds, t, (row) => row.name),
            ...tableColumns(t),
          ]}
          data={tables}
          getRowId={(row) => row.name}
          onRowClick={handleTableRowClick}
          hasMore={tableList.hasMore}
          onLoadMore={tableList.loadMore}
          loadingMore={tableList.isFetchingMore}
        />
      ) : (
        /* Items view — always show Splitter */
        <Splitter
          direction="horizontal"
          initialSize={240}
          minSize={80}
          maxSize={600}
          storageKey="vs-split-dynamodb-detail"
        >
          <div className="scroll-panel">
            {allItems.length > 0 ? (
              <DataTable
                columns={[
                  checkboxColumn<ItemRow>(selectedItemKeys, toggleItem, () => toggleAllItems(allItemIds), allItemIds, t, (row) => row.keyJson),
                  ...itemColumns,
                ]}
                data={allItems}
                getRowId={(row) => row.keyJson}
                onRowClick={handleItemRowClick}
                selectedId={selectedItem?.keyJson}
              />
            ) : (
              <div className="empty-state">
                {t("services.dynamodb.detail.noItems")}
                <button
                  className="btn btn-primary btn-sm attr-sk"
                  onClick={openAddItemModal}
                >
                  {t("services.dynamodb.addItem")}
                </button>
              </div>
            )}
            {nextEvalKey && (
              <div className="load-more">
                <button
                  className="btn btn-secondary btn-sm"
                  onClick={handleLoadMore}
                  disabled={scanFetching}
                >
                  {scanFetching ? t("common.loading") : t("common.loadMore")}
                </button>
              </div>
            )}
          </div>
          {renderDetailPanel()}
        </Splitter>
      )}

      {/* ── Modals ──────────────────────────────────────────── */}

      {/* Create Table */}
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

      {/* Delete Table(s) */}
      <ServiceDeleteDialog
        open={showDeleteTables}
        title={t("services.dynamodb.delete")}
        name={Array.from(selectedTableNames).join(", ")}
        error={deleteTableMutation.error}
        isPending={deleteTableMutation.isPending}
        onConfirm={async () => {
          const results = await Promise.allSettled(
            Array.from(selectedTableNames).map(name =>
              client.deleteTable(create(DeleteTableInputSchema, { tablename: name }))
            ),
          );
          const succeeded = results.filter(r => r.status === "fulfilled").length;
          const failed = results.filter(r => r.status === "rejected").length;
          setBatchResult(failed === 0
            ? t("services.dynamodb.batchDeleteResult", { count: succeeded })
            : `${t("services.dynamodb.batchDeleteResult", { count: succeeded })} (${failed} ${t("services.dynamodb.failed")})`);
          tableList.invalidate();
          clearTableSelection();
          setShowDeleteTables(false);
          setTimeout(() => setBatchResult(null), 5000);
        }}
        onClose={() => setShowDeleteTables(false)}
      />

      {/* Batch Delete Items */}
      <Modal open={showDeleteItems} onClose={() => setShowDeleteItems(false)}>
        <h2>{t("services.dynamodb.batchDelete")}</h2>
        <p>{t("services.dynamodb.batchDeleteConfirm", { count: selectedItemKeys.size })}</p>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowDeleteItems(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-danger"
            onClick={batchDeleteItems}
          >
            {t("common.delete")}
          </button>
        </div>
      </Modal>

      {/* Add Item — Structured + JSON tabs */}
      <Modal open={showAddItem} onClose={() => { setShowAddItem(false); setJsonError(""); }}>
        <h2>{t("services.dynamodb.addItem")}</h2>

        {/* Tab switcher */}
        <div className="attr-row-head--flex attr-tab-bar">
          <button
            className={`detail-tab ${putItemTab === "structured" ? "active" : ""}`}
            onClick={switchToStructured}
          >
            {t("services.dynamodb.structuredForm")}
          </button>
          <button
            className={`detail-tab ${putItemTab === "json" ? "active" : ""}`}
            onClick={switchToJson}
          >
            {t("services.dynamodb.jsonInput")}
          </button>
        </div>

        {jsonError && <div className="modal-error">{jsonError}</div>}
        {putItemMutation.error && <div className="modal-error">{String(putItemMutation.error)}</div>}

        {putItemTab === "structured" ? (
          <div>
            {/* PK/SK info */}
            <div className="footer-note attr-info-bar">
              <span>PK: <strong className="cell-mono">{pkName}</strong></span>
              {skName && <span className="attr-sk">SK: <strong className="cell-mono">{skName}</strong></span>}
            </div>

            {/* Attribute rows */}
            {attrRows.map((row, i) => {
              const isKey = row.name !== "" && (row.name === pkName || row.name === skName);
              const isListOrMap = row.type === "L" || row.type === "M";
              return (
              <div key={i} style={{ marginBottom: 4 }}>
                <div className="attr-row-head">
                  <input
                    value={row.name}
                    onChange={(e) => updateAttrRow(i, "name", e.target.value)}
                    placeholder={t("services.dynamodb.attributeName")}
                    className="modal-input attr-row-input"
                    readOnly={isKey}
                    title={isKey ? t("services.dynamodb.keyAttribute") : undefined}
                  />
                  <select
                    value={row.type}
                    onChange={(e) => updateAttrRow(i, "type", e.target.value as AttrType)}
                    className="modal-input attr-row-enum"
                  >
                    {ATTR_TYPES.map((at) => (
                      <option key={at} value={at}>{t(`services.dynamodb.${TYPE_LABEL_KEYS[at]}`)}</option>
                    ))}
                  </select>
                  {!isKey && <button className="btn btn-secondary btn-sm attr-remove-btn" onClick={() => removeAttrRow(i)}>✕</button>}
                </div>
                {isListOrMap ? (
                  <div className="attr-nested-border">
                    {(row.children ?? []).map((child, ci) => (
                      <div key={ci} className="attr-child-row">
                        {row.type === "M" ? (
                          <input value={child.name} onChange={(e) => updateChildRow(i, ci, "name", e.target.value)} placeholder="key" className="modal-input attr-enum-input" />
                        ) : (
                          <span className="attr-enum">{ci}</span>
                        )}
                        <select value={child.type} onChange={(e) => updateChildRow(i, ci, "type", e.target.value as AttrType)} className="modal-input attr-enum-input">
                          {ATTR_TYPES.map((at) => (<option key={at} value={at}>{t(`services.dynamodb.${TYPE_LABEL_KEYS[at]}`)}</option>))}
                        </select>
                        {child.type === "BOOL" ? (
                          <input type="checkbox" checked={child.value === "true"} onChange={(e) => updateChildRow(i, ci, "value", e.target.checked ? "true" : "false")} />
                        ) : child.type === "NULL" ? (
                          <span className="attr-null">null</span>
                        ) : (child.type === "L" || child.type === "M") ? (
                          <span className="attr-count">{child.type === "L" ? `[${child.children?.length ?? 0}]` : `{${child.children?.length ?? 0}}`}</span>
                        ) : (
                          <input value={child.value} onChange={(e) => updateChildRow(i, ci, "value", e.target.value)} placeholder={typePlaceholder(child.type)} className="modal-input attr-row-input" />
                        )}
                        <button className="btn btn-secondary btn-sm attr-remove-btn" onClick={() => removeChildRow(i, ci)}>✕</button>
                      </div>
                    ))}
                    <button className="btn btn-secondary btn-sm attr-add-btn" onClick={() => addChildRow(i)}>
                      + {row.type === "L" ? "item" : "entry"}
                    </button>
                  </div>
                ) : row.type === "BOOL" ? (
                  <label className="attr-label">
                    <input
                      type="checkbox"
                      checked={row.value === "true"}
                      onChange={(e) => updateAttrRow(i, "value", e.target.checked ? "true" : "false")}
                    />
                    {row.value === "true" ? "true" : "false"}
                  </label>
                ) : row.type !== "NULL" && (
                  <input
                    value={row.value}
                    onChange={(e) => updateAttrRow(i, "value", e.target.value)}
                    placeholder={typePlaceholder(row.type)}
                    className="modal-input attr-input-full"
                  />
                )}
              </div>
              );
            })}

            <button className="btn btn-secondary btn-sm attr-add-sm" onClick={addAttrRow}>
              {t("services.dynamodb.addAttribute")}
            </button>

            <div className="modal-actions modal-actions--spaced">
              <button className="btn btn-secondary" onClick={() => { setShowAddItem(false); setJsonError(""); }}>
                {t("common.cancel")}
              </button>
              <button
                className="btn btn-primary"
                disabled={attrRows.length === 0 || attrRows.some((r) => !r.name.trim() || (r.type !== "NULL" && r.type !== "L" && r.type !== "M" && !r.value.trim())) || putItemMutation.isPending}
                onClick={handleAddItemStructured}
              >
                {putItemMutation.isPending ? t("common.creating") : t("common.create")}
              </button>
            </div>
          </div>
        ) : (
          <div>
            <textarea
              value={itemJson}
              onChange={(e) => setItemJson(e.target.value)}
              placeholder={t("services.dynamodb.jsonPlaceholderPlain")}
              className="modal-input attr-monospace"
              rows={8}
            />
            <div className="modal-actions">
              <button className="btn btn-secondary" onClick={() => { setShowAddItem(false); setJsonError(""); }}>
                {t("common.cancel")}
              </button>
              <button
                className="btn btn-primary"
                disabled={!itemJson || putItemMutation.isPending}
                onClick={handleAddItemJson}
              >
                {putItemMutation.isPending ? t("common.creating") : t("common.create")}
              </button>
            </div>
          </div>
        )}
      </Modal>

      {/* Delete Single Item */}
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
    </ServicePageLayout>
  );
}
