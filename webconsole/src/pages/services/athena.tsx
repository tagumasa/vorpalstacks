/**
 * Athena service page. Lists work groups with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { AthenaService, type WorkGroupSummary, WorkGroupState } from "@/gen/athena_pb";
import { CreateWorkGroupInputSchema } from "@/gen/athena_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  DateCell,
  FallbackCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Lookup map for WorkGroupState proto enum values to i18n keys. */
const WORKGROUP_STATE_I18N: Record<number, string> = {
  [WorkGroupState.DISABLED]: "services.athena.stateDisabled",
  [WorkGroupState.ENABLED]: "services.athena.stateEnabled",
};

/** Column definitions for the Athena work group table. */
const getColumns = (t: TFunction): ColumnDef<WorkGroupSummary, any>[] => [
  { accessorKey: "name", header: t("services.athena.workGroupHeader"), cell: MonoCell },
  { accessorKey: "state", header: t("services.athena.stateHeader"), cell: ({ getValue }) => { const v = getValue() as number; return <span className="badge">{WORKGROUP_STATE_I18N[v] ? t(WORKGROUP_STATE_I18N[v]) : String(v)}</span>; }, size: 90 },
  { accessorKey: "engineversion", header: t("services.athena.engineVersionHeader"), cell: ({ getValue }) => { const ev = getValue(); return ev ? String(ev) : "\u2014"; }, size: 100 },
  { accessorKey: "description", header: t("services.athena.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "creationtime", header: t("services.athena.createdHeader"), cell: DateCell },
];

/** Athena service page with list, create, and delete operations. */
export function AthenaPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<WorkGroupSummary | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDescription, setFormDescription] = useState("");

  const { client, invalidate } = useServiceClient(AthenaService);
  const { queryKey } = useListKey("athena");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listWorkGroups({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: WorkGroupSummary[] = dropEmpty(data?.workgroups ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createWorkGroup(
        create(CreateWorkGroupInputSchema, {
          name: formName,
          description: formDescription,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormDescription("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (workgroupName: string) =>
      client.deleteWorkGroup({ workgroup: workgroupName, recursivedeleteoption: true }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🔬"
      title={t("services.athena.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.athena.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.athena.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "athena-items" }}
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
        title={t("services.athena.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.athena.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.athena.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.athena.descLabel")}
          <input
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder={t("services.athena.descPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.athena.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.name)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
