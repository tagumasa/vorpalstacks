/**
 * Lambda service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { LambdaService, Runtime, CreateFunctionRequestSchema, FunctionCodeSchema, type FunctionConfiguration } from "@/gen/lambda_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, BadgeCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const RUNTIME_LABELS: Record<number, string> = {
  [Runtime.DOTNETCORE31]: "dotnetcore3.1", [Runtime.NODEJS16X]: "nodejs16.x", [Runtime.RUBY25]: "ruby2.5", [Runtime.NODEJS24X]: "nodejs24.x",
  [Runtime.PYTHON39]: "python3.9", [Runtime.DOTNET6]: "dotnet6", [Runtime.RUBY32]: "ruby3.2", [Runtime.GO1X]: "go1.x",
  [Runtime.PYTHON36]: "python3.6", [Runtime.PROVIDEDAL2023]: "provided.al2023", [Runtime.NODEJS22X]: "nodejs22.x", [Runtime.PYTHON312]: "python3.12",
  [Runtime.DOTNETCORE21]: "dotnetcore2.1", [Runtime.NODEJS810]: "nodejs8.10", [Runtime.RUBY27]: "ruby2.7", [Runtime.JAVA21]: "java21",
  [Runtime.NODEJS14X]: "nodejs14.x", [Runtime.PROVIDED]: "provided", [Runtime.NODEJS610]: "nodejs6.10", [Runtime.PYTHON310]: "python3.10",
  [Runtime.PYTHON38]: "python3.8", [Runtime.JAVA17]: "java17", [Runtime.NODEJS43]: "nodejs4.3", [Runtime.NODEJS]: "nodejs",
  [Runtime.NODEJS20X]: "nodejs20.x", [Runtime.PYTHON313]: "python3.13", [Runtime.DOTNETCORE20]: "dotnetcore2.0", [Runtime.NODEJS43EDGE]: "nodejs4.3-edge",
  [Runtime.JAVA11]: "java11", [Runtime.NODEJS12X]: "nodejs12.x", [Runtime.RUBY33]: "ruby3.3", [Runtime.JAVA8]: "java8",
  [Runtime.DOTNETCORE10]: "dotnetcore1.0", [Runtime.PYTHON37]: "python3.7", [Runtime.PYTHON311]: "python3.11", [Runtime.JAVA8AL2]: "java8.al2",
  [Runtime.RUBY34]: "ruby3.4", [Runtime.PROVIDEDAL2]: "provided.al2", [Runtime.NODEJS10X]: "nodejs10.x", [Runtime.JAVA25]: "java25",
  [Runtime.PYTHON27]: "python2.7", [Runtime.DOTNET8]: "dotnet8", [Runtime.NODEJS18X]: "nodejs18.x", [Runtime.DOTNET10]: "dotnet10",
  [Runtime.PYTHON314]: "python3.14",
};

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

const getColumns = (t: TFunction): ColumnDef<FunctionConfiguration, any>[] => [
  { accessorKey: "functionname", header: t("services.lambda.functionNameHeader"), cell: MonoCell },
  { accessorKey: "runtime", header: t("services.lambda.runtimeHeader"), cell: ({ getValue }) => { const v = getValue(); return v !== undefined && v !== 0 ? <span className="badge">{RUNTIME_LABELS[v as number] ?? String(v)}</span> : "\u2014"; }, size: 120 },
  { accessorKey: "memorysize", header: t("services.lambda.memoryHeader"), size: 100 },
  { accessorKey: "timeout", header: t("services.lambda.timeoutHeader"), size: 90 },
  { accessorKey: "handler", header: t("services.lambda.handlerHeader"), cell: SmallMonoCell, size: 130 },
  { accessorKey: "state", header: t("services.lambda.stateHeader"), cell: ({ getValue }) => <BadgeCell getValue={getValue} positive={["Active"]} negative={["Failed", "Inactive"]} />, size: 90 },
  { accessorKey: "lastmodified", header: t("services.lambda.lastModifiedHeader"), cell: DateCell, size: 140 },
];

type DetailTab = "detail" | "json";

export function LambdaPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(LambdaService);
  const { queryKey } = useListKey("lambda");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<FunctionConfiguration | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
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

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listFunctions({}), refetchInterval: REFETCH_INTERVAL });
  const items: FunctionConfiguration[] = dropEmpty(data?.functions ?? [], "functionname");

  const createMutation = useMutation({
    mutationFn: () => {
      const envVars: Record<string, string> = {};
      if (formEnvVars) { for (const line of formEnvVars.split("\n")) { const [k, ...v] = line.split("="); if (k) envVars[k.trim()] = v.join("=").trim(); } }
      return client.createFunction(create(CreateFunctionRequestSchema, {
        functionname: formName, runtime: formRuntime, handler: formHandler, role: formRole,
        code: create(FunctionCodeSchema, { s3bucket: formS3Bucket || undefined, s3key: formS3Key || undefined }),
        memorysize: formMemory, timeout: formTimeout, description: formDescription,
        environment: Object.keys(envVars).length > 0 ? { variables: envVars } : undefined,
      }));
    },
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormHandler("index.handler"); setFormRole(""); setFormMemory(128); setFormTimeout(3); setFormDescription(""); setFormEnvVars(""); setFormS3Bucket(""); setFormS3Key(""); },
  });

  const deleteMutation = useMutation({
    mutationFn: (name: string) => client.deleteFunction({ functionname: name }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteFunction({ functionname: n }))),
    onSuccess: (_d, names) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && names.includes(p.functionname) ? null : p)); },
  });

  const handleRowClick = (row: FunctionConfiguration) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.functionname);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const rt = RUNTIME_LABELS[selectedItem.runtime] ?? String(selectedItem.runtime);
    return (
      <DetailPanel title={selectedItem.functionname} titleIcon="⚡" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">Function</td><td className="cell-mono">{selectedItem.functionname}</td></tr>
            <tr><td className="detail-label">Runtime</td><td><span className="badge">{rt}</span></td></tr>
            <tr><td className="detail-label">Handler</td><td className="cell-mono">{selectedItem.handler || "\u2014"}</td></tr>
            <tr><td className="detail-label">Memory</td><td>{selectedItem.memorysize} MB</td></tr>
            <tr><td className="detail-label">Timeout</td><td>{selectedItem.timeout}s</td></tr>
            <tr><td className="detail-label">State</td><td><span className="badge">{selectedItem.state || "\u2014"}</span></td></tr>
            <tr><td className="detail-label">Code Size</td><td>{selectedItem.codesize ? `${selectedItem.codesize} bytes` : "\u2014"}</td></tr>
            {selectedItem.description && <tr><td className="detail-label">Description</td><td>{selectedItem.description}</td></tr>}
            {selectedItem.functionarn && <tr><td className="detail-label">ARN</td><td className="cell-mono cell-long">{selectedItem.functionarn}</td></tr>}
            {selectedItem.lastmodified && <tr><td className="detail-label">Modified</td><td>{fmtDate(selectedItem.lastmodified, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="⚡" title={t("services.lambda.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.lambda.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.lambda.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span className="batch-count">({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.lambda.title") }, { label: t("services.lambda.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-lambda">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<FunctionConfiguration>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.functionname), ...columns]} data={items} getRowId={(row) => row.functionname} onRowClick={handleRowClick} selectedId={selectedItem?.functionname} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.lambda.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName || !formRole}>
        <label>{t("services.lambda.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.lambda.placeholder")} className="modal-input" /></label>
        <label>{t("services.lambda.runtimeLabel")}<select value={formRuntime} onChange={(e) => setFormRuntime(Number(e.target.value))} className="modal-input">{POPULAR_RUNTIMES.map((r) => <option key={r.value} value={r.value}>{t(r.i18nKey)}</option>)}</select></label>
        <label>{t("services.lambda.handlerLabel")}<input value={formHandler} onChange={(e) => setFormHandler(e.target.value)} className="modal-input" /></label>
        <label>{t("services.lambda.roleLabel")}<input value={formRole} onChange={(e) => setFormRole(e.target.value)} placeholder={t("services.lambda.rolePlaceholder")} className="modal-input" /></label>
        <label>{t("services.lambda.memoryLabel")}<input type="number" min={128} max={10240} value={formMemory} onChange={(e) => setFormMemory(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.lambda.timeoutLabel")}<input type="number" min={1} max={900} value={formTimeout} onChange={(e) => setFormTimeout(Number(e.target.value))} className="modal-input" /></label>
        <label>{t("services.lambda.descriptionLabel")}<input value={formDescription} onChange={(e) => setFormDescription(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label>{t("services.lambda.s3BucketLabel")}<input value={formS3Bucket} onChange={(e) => setFormS3Bucket(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label>{t("services.lambda.s3KeyLabel")}<input value={formS3Key} onChange={(e) => setFormS3Key(e.target.value)} placeholder={t("common.optional")} className="modal-input" /></label>
        <label>{t("services.lambda.envVarsLabel")}<textarea value={formEnvVars} onChange={(e) => setFormEnvVars(e.target.value)} placeholder="KEY=value" rows={3} className="modal-input" style={{ fontFamily: "monospace", fontSize: "0.85em" }} /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.lambda.delete")} name={selectedItem?.functionname} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.functionname)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.lambda.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
