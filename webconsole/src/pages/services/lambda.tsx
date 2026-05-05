/**
 * Lambda service page. Lists functions with create/delete CRUD operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import {
  LambdaService,
  Runtime,
  CreateFunctionRequestSchema,
  FunctionCodeSchema,
  type FunctionConfiguration,
} from "@/gen/lambda_pb";
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
  BadgeCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Lookup map for Runtime proto enum values to human-readable labels. */
const RUNTIME_LABELS: Record<number, string> = {
  [Runtime.DOTNETCORE31]: "dotnetcore3.1",
  [Runtime.NODEJS16X]: "nodejs16.x",
  [Runtime.RUBY25]: "ruby2.5",
  [Runtime.NODEJS24X]: "nodejs24.x",
  [Runtime.PYTHON39]: "python3.9",
  [Runtime.DOTNET6]: "dotnet6",
  [Runtime.RUBY32]: "ruby3.2",
  [Runtime.GO1X]: "go1.x",
  [Runtime.PYTHON36]: "python3.6",
  [Runtime.PROVIDEDAL2023]: "provided.al2023",
  [Runtime.NODEJS22X]: "nodejs22.x",
  [Runtime.PYTHON312]: "python3.12",
  [Runtime.DOTNETCORE21]: "dotnetcore2.1",
  [Runtime.NODEJS810]: "nodejs8.10",
  [Runtime.RUBY27]: "ruby2.7",
  [Runtime.JAVA21]: "java21",
  [Runtime.NODEJS14X]: "nodejs14.x",
  [Runtime.PROVIDED]: "provided",
  [Runtime.NODEJS610]: "nodejs6.10",
  [Runtime.PYTHON310]: "python3.10",
  [Runtime.PYTHON38]: "python3.8",
  [Runtime.JAVA17]: "java17",
  [Runtime.NODEJS43]: "nodejs4.3",
  [Runtime.NODEJS]: "nodejs",
  [Runtime.NODEJS20X]: "nodejs20.x",
  [Runtime.PYTHON313]: "python3.13",
  [Runtime.DOTNETCORE20]: "dotnetcore2.0",
  [Runtime.NODEJS43EDGE]: "nodejs4.3-edge",
  [Runtime.JAVA11]: "java11",
  [Runtime.NODEJS12X]: "nodejs12.x",
  [Runtime.RUBY33]: "ruby3.3",
  [Runtime.JAVA8]: "java8",
  [Runtime.DOTNETCORE10]: "dotnetcore1.0",
  [Runtime.PYTHON37]: "python3.7",
  [Runtime.PYTHON311]: "python3.11",
  [Runtime.JAVA8AL2]: "java8.al2",
  [Runtime.RUBY34]: "ruby3.4",
  [Runtime.PROVIDEDAL2]: "provided.al2",
  [Runtime.NODEJS10X]: "nodejs10.x",
  [Runtime.JAVA25]: "java25",
  [Runtime.PYTHON27]: "python2.7",
  [Runtime.DOTNET8]: "dotnet8",
  [Runtime.NODEJS18X]: "nodejs18.x",
  [Runtime.DOTNET10]: "dotnet10",
  [Runtime.PYTHON314]: "python3.14",
};

/** Common Lambda runtime options for the create form. */
const POPULAR_RUNTIMES: { value: Runtime; i18nKey: string }[] = [
  { value: Runtime.NODEJS24X, i18nKey: "services.lambda.runtimeNodejs24" },
  { value: Runtime.NODEJS22X, i18nKey: "services.lambda.runtimeNodejs22" },
  { value: Runtime.PYTHON314, i18nKey: "services.lambda.runtimePython314" },
  { value: Runtime.PYTHON313, i18nKey: "services.lambda.runtimePython313" },
  { value: Runtime.PYTHON312, i18nKey: "services.lambda.runtimePython312" },
  { value: Runtime.PYTHON311, i18nKey: "services.lambda.runtimePython311" },
  { value: Runtime.GO1X, i18nKey: "services.lambda.runtimeGo1x" },
  { value: Runtime.JAVA21, i18nKey: "services.lambda.runtimeJava21" },
  { value: Runtime.JAVA17, i18nKey: "services.lambda.runtimeJava17" },
  { value: Runtime.DOTNET8, i18nKey: "services.lambda.runtimeDotnet8" },
  { value: Runtime.RUBY33, i18nKey: "services.lambda.runtimeRuby33" },
  { value: Runtime.PROVIDEDAL2023, i18nKey: "services.lambda.runtimeProvidedAl2023" },
];

/** Column definitions for the Lambda function table. */
const getColumns = (t: TFunction): ColumnDef<FunctionConfiguration, any>[] => [
  { accessorKey: "functionname", header: t("services.lambda.functionNameHeader"), cell: MonoCell },
  {
    accessorKey: "runtime",
    header: t("services.lambda.runtimeHeader"),
    cell: ({ getValue }) => {
      const v = getValue();
      return v !== undefined && v !== 0 ? (
        <span className="badge">{RUNTIME_LABELS[v as number] ?? String(v)}</span>
      ) : (
        "\u2014"
      );
    },
    size: 120,
  },
  { accessorKey: "memorysize", header: t("services.lambda.memoryHeader"), size: 100 },
  { accessorKey: "timeout", header: t("services.lambda.timeoutHeader"), size: 90 },
  { accessorKey: "handler", header: t("services.lambda.handlerHeader"), cell: SmallMonoCell, size: 130 },
  { accessorKey: "codesize", header: t("services.lambda.codeSizeHeader"), size: 90 },
  { accessorKey: "state", header: t("services.lambda.stateHeader"), cell: ({ getValue }) => <BadgeCell getValue={getValue} positive={["Active"]} negative={["Failed", "Inactive"]} />, size: 90 },
  { accessorKey: "lastmodified", header: t("services.lambda.lastModifiedHeader"), cell: DateCell, size: 140 },
  { accessorKey: "packagetype", header: t("services.lambda.packageTypeHeader"), size: 80 },
  { accessorKey: "functionarn", header: t("services.lambda.arnHeader"), cell: SmallMonoCell, size: 200 },
  {
    accessorKey: "description",
    header: t("services.lambda.descriptionHeader"),
    cell: FallbackCell,
  },
];

/** Lambda service page with function list, create, and delete operations. */
export function LambdaPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
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
  const [formEnvVars, setFormEnvVars] = useState("");
  const [formS3Bucket, setFormS3Bucket] = useState("");
  const [formS3Key, setFormS3Key] = useState("");

  const { client, invalidate } = useServiceClient(LambdaService);
  const { queryKey } = useListKey("lambda");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listFunctions({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: FunctionConfiguration[] = dropEmpty(data?.functions ?? [], "functionname");

  const createMutation = useMutation({
    mutationFn: () => {
      let envVars: { [key: string]: string } | undefined;
      if (formEnvVars.trim()) {
        try {
          envVars = JSON.parse(formEnvVars);
        } catch {
          throw new Error(t("common.invalidJson", { field: "Environment Variables" }));
        }
      }
      return client.createFunction(
        create(CreateFunctionRequestSchema, {
          functionname: formName,
          runtime: formRuntime,
          handler: formHandler,
          role: formRole,
          memorysize: formMemory,
          timeout: formTimeout,
          description: formDescription,
          ...(formS3Bucket || formS3Key
            ? { code: create(FunctionCodeSchema, { s3bucket: formS3Bucket, s3key: formS3Key }) }
            : {}),
          ...(envVars ? { environment: { variables: envVars } } : {}),
        }),
      );
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormRuntime(Runtime.NODEJS24X);
      setFormHandler("index.handler");
      setFormRole("");
      setFormMemory(128);
      setFormTimeout(3);
      setFormDescription("");
      setFormEnvVars("");
      setFormS3Bucket("");
      setFormS3Key("");
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

  return (
    <ServicePageLayout
      icon="⚡"
      title={t("services.lambda.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.lambda.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.lambda.create")}
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
        getRowId={(row) => row.functionname}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.functionname}
        selected={selectedItem}
        detailTitle={selectedItem?.functionname}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.lambda.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formRole}
      >
        <label>
          {t("services.lambda.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.lambda.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.lambda.runtimeLabel")}
          <select
            value={formRuntime}
            onChange={(e) => setFormRuntime(Number(e.target.value))}
            className="modal-input"
          >
            {POPULAR_RUNTIMES.map((r) => (
              <option key={r.value} value={r.value}>{t(r.i18nKey)}</option>
            ))}
          </select>
        </label>
        <label>
          {t("services.lambda.handlerLabel")}
          <input
            value={formHandler}
            onChange={(e) => setFormHandler(e.target.value)}
            placeholder={t("services.lambda.handlerPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.lambda.roleLabel")}
          <input
            value={formRole}
            onChange={(e) => setFormRole(e.target.value)}
            placeholder={t("services.lambda.rolePlaceholder")}
            className="modal-input"
          />
        </label>
        <div className="form-row">
          <label>
            {t("services.lambda.memoryLabel")}
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
            {t("services.lambda.timeoutLabel")}
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
          {t("services.lambda.descLabel")}
          <input
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder={t("services.lambda.descPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.lambda.envVarsLabel")}
          <textarea
            value={formEnvVars}
            onChange={(e) => setFormEnvVars(e.target.value)}
            placeholder='{"KEY":"value"}'
            rows={4}
            className="modal-input"
            style={{ fontFamily: "monospace", fontSize: "0.85em" }}
          />
        </label>
        <label>
          {t("services.lambda.s3BucketLabel")}
          <input
            value={formS3Bucket}
            onChange={(e) => setFormS3Bucket(e.target.value)}
            placeholder={t("services.lambda.s3BucketPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.lambda.s3KeyLabel")}
          <input
            value={formS3Key}
            onChange={(e) => setFormS3Key(e.target.value)}
            placeholder={t("services.lambda.s3KeyPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.lambda.delete")}
        name={selectedItem?.functionname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.functionname)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
