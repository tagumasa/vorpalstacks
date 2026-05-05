/**
 * Route 53 service page. Lists hosted zones with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { Route53Service, type HostedZone } from "@/gen/route53_pb";
import { CreateHostedZoneRequestSchema, HostedZoneConfigSchema } from "@/gen/route53_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  FallbackCell,
  BooleanCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the Route 53 hosted zone table. */
const getColumns = (t: TFunction): ColumnDef<HostedZone, any>[] => [
  { accessorKey: "name", header: t("services.route53.zoneNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.route53.zoneIdHeader"), cell: MonoCell },
  { accessorKey: "resourcerecordsetcount", header: t("services.route53.recordsHeader"), size: 80 },
  { accessorKey: "config.comment", header: t("services.route53.commentHeader"), cell: FallbackCell },
  { accessorKey: "config.privatezone", header: t("services.route53.privateHeader"), cell: BooleanCell, size: 70 },
];

/** Route 53 service page with list, create, and delete operations. */
export function Route53Page() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<HostedZone | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formComment, setFormComment] = useState("");

  const { client, invalidate } = useServiceClient(Route53Service);
  const { queryKey } = useListKey("route53");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const resp = await client.listHostedZones({});
      return dropEmpty(resp.hostedzones ?? [], "id");
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: HostedZone[] = data ?? [];

  const createMutation = useMutation({
    mutationFn: () =>
      client.createHostedZone(
        create(CreateHostedZoneRequestSchema, {
          name: formName,
          callerreference: `${Date.now()}-${Math.random().toString(36).slice(2)}`,
          hostedzoneconfig: formComment
            ? create(HostedZoneConfigSchema, { comment: formComment })
            : undefined,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormComment("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (zoneId: string) =>
      client.deleteHostedZone({ id: zoneId }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🌐"
      title={t("services.route53.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.route53.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.route53.create")}
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
        getRowId={(row) => row.id}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.id}
        selected={selectedItem}
        detailTitle={selectedItem?.name}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.route53.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.route53.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.route53.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.route53.commentLabel")}
          <input
            value={formComment}
            onChange={(e) => setFormComment(e.target.value)}
            placeholder={t("services.route53.commentPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.route53.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
