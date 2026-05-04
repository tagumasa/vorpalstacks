/**
 * Step Functions service page. Lists state machines with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { SFNService, type StateMachineListItem, StateMachineType } from "@/gen/sfn_pb";
import { CreateStateMachineInputSchema } from "@/gen/sfn_pb";
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

/** Lookup map for StateMachineType proto enum values to human-readable labels. */
const SFN_TYPE_LABELS: Record<number, string> = {
  [StateMachineType.EXPRESS]: "EXPRESS",
  [StateMachineType.STANDARD]: "STANDARD",
};

/** Column definitions for the Step Functions state machine table. */
const getColumns = (t: TFunction): ColumnDef<StateMachineListItem, any>[] => [
  { accessorKey: "name", header: t("services.sfn.nameHeader"), cell: MonoCell },
  { accessorKey: "statemachinearn", header: t("services.sfn.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "type", header: t("services.sfn.typeHeader"), cell: ({ getValue }) => <span className="badge">{SFN_TYPE_LABELS[getValue() as number] ?? String(getValue())}</span>, size: 100 },
  { accessorKey: "creationdate", header: t("services.sfn.createdHeader"), cell: DateCell },
];

/** Step Functions service page with list, create, and delete operations. */
export function SFNPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<StateMachineListItem | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDefinition, setFormDefinition] = useState(
    JSON.stringify(
      { Comment: "A Hello World example", StartAt: "HelloWorld", States: { HelloWorld: { Type: "Pass", Result: "Hello World!", End: true } } },
      null,
      2,
    ),
  );
  const [formRoleArn, setFormRoleArn] = useState("");
  const [formType, setFormType] = useState(StateMachineType.STANDARD);

  const { client, invalidate } = useServiceClient(SFNService);
  const { queryKey } = useListKey("sfn");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listStateMachines({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: StateMachineListItem[] = dropEmpty(data?.statemachines ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createStateMachine(
        create(CreateStateMachineInputSchema, {
          name: formName,
          definition: formDefinition,
          rolearn: formRoleArn,
          type: formType,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormRoleArn("");
      setFormType(StateMachineType.STANDARD);
      setFormDefinition(
        JSON.stringify(
          { Comment: "A Hello World example", StartAt: "HelloWorld", States: { HelloWorld: { Type: "Pass", Result: "Hello World!", End: true } } },
          null,
          2,
        ),
      );
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (statemachinearn: string) =>
      client.deleteStateMachine({ statemachinearn }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🔀"
      title={t("services.sfn.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.sfn.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.sfn.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "sfn-items" }}
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
        title={t("services.sfn.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName || !formRoleArn}
      >
        <label>
          {t("services.sfn.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.sfn.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sfn.roleArnLabel")}
          <input
            value={formRoleArn}
            onChange={(e) => setFormRoleArn(e.target.value)}
            placeholder={t("services.sfn.roleArnPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.sfn.typeLabel")}
          <select
            value={formType}
            onChange={(e) => setFormType(Number(e.target.value))}
            className="modal-input"
          >
            <option value={StateMachineType.STANDARD}>{t("services.sfn.typeStandard")}</option>
            <option value={StateMachineType.EXPRESS}>{t("services.sfn.typeExpress")}</option>
          </select>
        </label>
        <label>
          {t("services.sfn.defLabel")}
          <textarea
            value={formDefinition}
            onChange={(e) => setFormDefinition(e.target.value)}
            rows={10}
            className="modal-input"
            style={{ fontFamily: "monospace", fontSize: "0.85em" }}
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.sfn.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.statemachinearn)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
