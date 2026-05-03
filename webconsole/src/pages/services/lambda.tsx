import { useTranslation } from "react-i18next";
import { useState } from "react";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { create } from "@bufbuild/protobuf";
import { LambdaService, Runtime } from "@/gen/lambda_pb";
import { CreateFunctionRequestSchema } from "@/gen/lambda_pb";
import { transport } from "@/lib/transport";
import { useListKey, useListInvalidator, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";
import { Modal } from "@/components/shared/modal";
import { ConfirmDialog } from "@/components/shared/confirm-dialog";
import type { FunctionConfiguration } from "@/gen/lambda_pb";

const POPULAR_RUNTIMES: { value: Runtime; label: string }[] = [
  { value: Runtime.NODEJS24X, label: "Node.js 24" },
  { value: Runtime.NODEJS22X, label: "Node.js 22" },
  { value: Runtime.PYTHON314, label: "Python 3.14" },
  { value: Runtime.PYTHON313, label: "Python 3.13" },
  { value: Runtime.PYTHON312, label: "Python 3.12" },
  { value: Runtime.PYTHON311, label: "Python 3.11" },
  { value: Runtime.GO1X, label: "Go 1.x" },
  { value: Runtime.JAVA21, label: "Java 21" },
  { value: Runtime.JAVA17, label: "Java 17" },
  { value: Runtime.DOTNET8, label: ".NET 8" },
  { value: Runtime.RUBY33, label: "Ruby 3.3" },
  { value: Runtime.PROVIDEDAL2023, label: "Custom (AL2023)" },
];

const columns: ColumnDef<FunctionConfiguration, any>[] = [
  {
    accessorKey: "functionname",
    header: "Function Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "runtime",
    header: "Runtime",
    cell: ({ getValue }) => {
      const v = getValue();
      return v !== undefined && v !== 0 ? <span className="badge">{String(v)}</span> : "\u2014";
    },
    size: 120,
  },
  {
    accessorKey: "memorysize",
    header: "Memory (MB)",
    size: 100,
  },
  {
    accessorKey: "timeout",
    header: "Timeout (s)",
    size: 90,
  },
  {
    accessorKey: "description",
    header: "Description",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      return v || "\u2014";
    },
  },
];

export function LambdaPage() {
  const { t } = useTranslation();
  const [selectedItem, setSelectedItem] = useState<FunctionConfiguration | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);

  const [formName, setFormName] = useState("");
  const [formRuntime, setFormRuntime] = useState(Runtime.NODEJS24X);
  const [formHandler, setFormHandler] = useState("index.handler");
  const [formRole, setFormRole] = useState("");
  const [formMemory, setFormMemory] = useState(128);
  const [formTimeout, setFormTimeout] = useState(3);
  const [formDescription, setFormDescription] = useState("");

  const client = createClient(LambdaService, transport);
  const { queryKey } = useListKey("lambda");
  const invalidate = useListInvalidator();

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listFunctions({}),
    refetchInterval: 30_000,
  });

  const items: FunctionConfiguration[] = dropEmpty(data?.functions ?? [], "functionname");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createFunction(
        create(CreateFunctionRequestSchema, {
          functionname: formName,
          runtime: formRuntime,
          handler: formHandler,
          role: formRole,
          memorysize: formMemory,
          timeout: formTimeout,
          description: formDescription,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormHandler("index.handler");
      setFormRole("");
      setFormMemory(128);
      setFormTimeout(3);
      setFormDescription("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (functionName: string) =>
      client.deleteFunction({ functionname: functionName }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  if (isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">⚡</span>
          <h1>Lambda</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">⚡</span>
          <h1>Lambda</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(error) })}</div>
      </div>
    );
  }

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">⚡</span>
        <h1>Lambda</h1>
        <span className="resource-count">{items.length} functions</span>
        <div className="page-actions">
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            Create function
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          <DataTable
            columns={columns}
            data={items}
            getRowId={(row) => row.functionname}
            onRowClick={setSelectedItem}
            selectedId={selectedItem?.functionname}
          />
        </div>
        {selectedItem && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{selectedItem.functionname}</h2>
              <button className="detail-close" onClick={() => setSelectedItem(null)}>
                ✕
              </button>
            </div>
            <JsonViewer data={selectedItem} />
          </div>
        )}
      </div>

      <Modal open={showCreate} onClose={() => setShowCreate(false)}>
        <h2>Create function</h2>
        {createMutation.error && (
          <div className="modal-error">{String(createMutation.error)}</div>
        )}
        <label>
          Function Name *
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder="my-function"
            className="modal-input"
          />
        </label>
        <label>
          Runtime *
          <select
            value={formRuntime}
            onChange={(e) => setFormRuntime(Number(e.target.value))}
            className="modal-input"
          >
            {POPULAR_RUNTIMES.map((r) => (
              <option key={r.value} value={r.value}>{r.label}</option>
            ))}
          </select>
        </label>
        <label>
          Handler *
          <input
            value={formHandler}
            onChange={(e) => setFormHandler(e.target.value)}
            placeholder="index.handler"
            className="modal-input"
          />
        </label>
        <label>
          Role ARN *
          <input
            value={formRole}
            onChange={(e) => setFormRole(e.target.value)}
            placeholder="arn:aws:iam::123456789012:role/lambda-role"
            className="modal-input"
          />
        </label>
        <div className="form-row">
          <label>
            Memory (MB)
            <input
              type="number"
              value={formMemory}
              onChange={(e) => setFormMemory(Number(e.target.value))}
              min={128}
              max={10240}
              step={64}
              className="modal-input"
            />
          </label>
          <label>
            Timeout (s)
            <input
              type="number"
              value={formTimeout}
              onChange={(e) => setFormTimeout(Number(e.target.value))}
              min={1}
              max={900}
              className="modal-input"
            />
          </label>
        </div>
        <label>
          Description
          <input
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder="Function description"
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCreate(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!formName || !formRole || createMutation.isPending}
            onClick={() => createMutation.mutate()}
          >
            {createMutation.isPending ? t("common.creating") : t("common.create")}
          </button>
        </div>
      </Modal>

      <ConfirmDialog
        open={showDelete && !!selectedItem}
        title="Delete function"
        message={
          <>
            {t("confirm.confirmDelete", { name: "" })}
            <strong>{selectedItem?.functionname}</strong>?
          </>
        }
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.functionname)}
        onClose={() => setShowDelete(false)}
      />
    </div>
  );
}
