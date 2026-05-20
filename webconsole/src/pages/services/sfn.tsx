/**
 * Step Functions service page — 3-panel inspector layout.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { SFNService, type StateMachineListItem, StateMachineType } from "@/gen/sfn_pb";
import { CreateStateMachineInputSchema } from "@/gen/sfn_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, Breadcrumb, SelectionBadge, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

const SFN_TYPE_LABELS: Record<number, string> = {
  [StateMachineType.EXPRESS]: "EXPRESS",
  [StateMachineType.STANDARD]: "STANDARD",
};

const DEFAULT_DEFINITION = JSON.stringify({ Comment: "A Hello World example", StartAt: "HelloWorld", States: { HelloWorld: { Type: "Pass", Result: "Hello World!", End: true } } }, null, 2);

const getColumns = (t: TFunction): ColumnDef<StateMachineListItem, any>[] => [
  { accessorKey: "name", header: t("services.sfn.nameHeader"), cell: MonoCell },
  { accessorKey: "statemachinearn", header: t("services.sfn.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "type", header: t("services.sfn.typeHeader"), cell: ({ getValue }) => <span className="badge">{SFN_TYPE_LABELS[getValue() as number] ?? String(getValue())}</span>, size: 100 },
  { accessorKey: "creationdate", header: t("services.sfn.createdHeader"), cell: DateCell },
];

type DetailTab = "detail" | "json";

export function SFNPage() {
  const { t, i18n } = useTranslation();
  const { client, invalidate } = useServiceClient(SFNService);
  const { queryKey } = useListKey("sfn");
  const columns = getColumns(t);

  const { selected: selectedIds, toggle, toggleAll: toggleAll_, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<StateMachineListItem | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDefinition, setFormDefinition] = useState(DEFAULT_DEFINITION);
  const [formRoleArn, setFormRoleArn] = useState("");
  const [formType, setFormType] = useState(StateMachineType.STANDARD);

  const { data, isLoading, error } = useQuery({ queryKey, queryFn: () => client.listStateMachines({}), refetchInterval: REFETCH_INTERVAL });
  const items: StateMachineListItem[] = dropEmpty(data?.statemachines ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () => client.createStateMachine(create(CreateStateMachineInputSchema, { name: formName, definition: formDefinition, rolearn: formRoleArn, type: formType })),
    onSuccess: () => { invalidate(queryKey); setShowCreate(false); setFormName(""); setFormRoleArn(""); setFormType(StateMachineType.STANDARD); setFormDefinition(DEFAULT_DEFINITION); },
  });

  const deleteMutation = useMutation({
    mutationFn: (arn: string) => client.deleteStateMachine({ statemachinearn: arn }),
    onSuccess: () => { invalidate(queryKey); setShowDelete(false); setSelectedItem(null); clearSelection(); },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (arns: string[]) => Promise.allSettled(arns.map((arn) => client.deleteStateMachine({ statemachinearn: arn }))),
    onSuccess: (_d, arns) => { invalidate(queryKey); setShowBatchDelete(false); clearSelection(); setSelectedItem((p) => (p && arns.includes(p.statemachinearn) ? null : p)); },
  });

  const handleRowClick = (row: StateMachineListItem) => { setSelectedItem(row); setDetailTab("detail"); };
  const allIds = items.map((i) => i.statemachinearn);

  const renderDetailPanel = () => {
    if (!selectedItem) return <DetailEmpty message={t("common.noItemSelected")} />;
    const typeLabel = SFN_TYPE_LABELS[selectedItem.type] ?? String(selectedItem.type);
    return (
      <DetailPanel title={selectedItem.name} titleIcon="🔀" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Name</td><td className="cell-mono">{selectedItem.name}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Type</td><td><span className="badge">{typeLabel}</span></td></tr>
            <tr><td style={{ fontWeight: 600 }}>ARN</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedItem.statemachinearn}</td></tr>
            {selectedItem.creationdate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{fmtDate(selectedItem.creationdate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedItem} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="🔀" title={t("services.sfn.title")} isLoading={isLoading} error={error} count={items.length} countLabel={t("services.sfn.countLabel")} actions={<>
      <button className="btn btn-primary" onClick={() => setShowCreate(true)}>{t("services.sfn.create")}</button>
      <button className="btn btn-danger" disabled={selectedIds.size === 0} onClick={() => setShowBatchDelete(true)}>{t("common.deleteSelected")}{selectedIds.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({selectedIds.size})</span>}</button>
    </>}>
      <div className="inspector-toolbar"><Breadcrumb parts={[{ label: t("services.sfn.title") }, { label: t("services.sfn.countLabel") }]} /><div className="toolbar-selection-info"><SelectionBadge count={selectedIds.size} label={t("common.selectedCount", { count: selectedIds.size })} /></div></div>
      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-sfn">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<StateMachineListItem>(selectedIds, toggle, () => toggleAll_(allIds), allIds, t, (row) => row.statemachinearn), ...columns]} data={items} getRowId={(row) => row.statemachinearn} onRowClick={handleRowClick} selectedId={selectedItem?.statemachinearn} /></div>
          {renderDetailPanel()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>}

      <ServiceCreateModal open={showCreate} onClose={() => setShowCreate(false)} title={t("services.sfn.create")} error={createMutation.error} isPending={createMutation.isPending} onCreate={() => createMutation.mutate()} disabled={!formName || !formRoleArn}>
        <label>{t("services.sfn.nameField")}<input value={formName} onChange={(e) => setFormName(e.target.value)} placeholder={t("services.sfn.placeholder")} className="modal-input" /></label>
        <label>{t("services.sfn.roleArnLabel")}<input value={formRoleArn} onChange={(e) => setFormRoleArn(e.target.value)} placeholder={t("services.sfn.roleArnPlaceholder")} className="modal-input" /></label>
        <label>{t("services.sfn.typeLabel")}<select value={formType} onChange={(e) => setFormType(Number(e.target.value))} className="modal-input"><option value={StateMachineType.STANDARD}>{t("services.sfn.typeStandard")}</option><option value={StateMachineType.EXPRESS}>{t("services.sfn.typeExpress")}</option></select></label>
        <label>{t("services.sfn.defLabel")}<textarea value={formDefinition} onChange={(e) => setFormDefinition(e.target.value)} rows={10} className="modal-input" style={{ fontFamily: "monospace", fontSize: "0.85em" }} /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDelete && !!selectedItem} title={t("services.sfn.delete")} name={selectedItem?.name} error={deleteMutation.error} isPending={deleteMutation.isPending} onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.statemachinearn)} onClose={() => setShowDelete(false)} />
      <ServiceDeleteDialog open={showBatchDelete} title={t("common.deleteSelected")} name={`${selectedIds.size} ${t("services.sfn.countLabel")}`} error={batchDeleteMutation.error} isPending={batchDeleteMutation.isPending} onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedIds))} onClose={() => setShowBatchDelete(false)} />
    </ServicePageLayout>
  );
}
