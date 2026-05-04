/**
 * WAFv2 service page. Lists WebACLs with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { WAFV2Service, type WebACLSummary, Scope } from "@/gen/wafv2_pb";
import { CreateWebACLRequestSchema } from "@/gen/wafv2_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  FallbackCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Extended type carrying the scope retrieved from the list response. */
type WebACLItem = WebACLSummary & { _scope: Scope };

/** Column definitions for the WAFv2 WebACL table. */
const getColumns = (t: TFunction): ColumnDef<WebACLItem, any>[] => [
  { accessorKey: "name", header: t("services.wafv2.webaclNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.wafv2.idHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.wafv2.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "arn", header: t("services.wafv2.arnHeader"), cell: SmallMonoCell },
];

/** WAFv2 service page with list, create, and delete operations. */
export function WAFv2Page() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<WebACLItem | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formScope, setFormScope] = useState<Scope>(Scope.REGIONAL);
  const [formDescription, setFormDescription] = useState("");
  const [formDefaultAction, setFormDefaultAction] = useState<"allow" | "block">("allow");

  const { client, invalidate } = useServiceClient(WAFV2Service);
  const { queryKey } = useListKey("wafv2");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: async () => {
      const [regional, cloudfront] = await Promise.all([
        client.listWebACLs({ scope: Scope.REGIONAL }),
        client.listWebACLs({ scope: Scope.CLOUDFRONT }),
      ]);
      return [
        ...(regional.webacls ?? []).map((acl) => ({ ...acl, _scope: Scope.REGIONAL })),
        ...(cloudfront.webacls ?? []).map((acl) => ({ ...acl, _scope: Scope.CLOUDFRONT })),
      ];
    },
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: WebACLItem[] = dropEmpty(data ?? [], "id");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createWebACL(
        create(CreateWebACLRequestSchema, {
          name: formName,
          scope: formScope,
          description: formDescription,
          defaultaction: formDefaultAction === "allow" ? { allow: {} } : { block: {} },
          visibilityconfig: {
            sampledrequestsenabled: false,
            cloudwatchmetricsenabled: false,
            metricname: formName,
          },
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormScope(Scope.REGIONAL);
      setFormDescription("");
      setFormDefaultAction("allow");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (acl: WebACLItem) =>
      client.deleteWebACL({
        id: acl.id,
        name: acl.name,
        scope: acl._scope,
        locktoken: acl.locktoken,
      }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🛡️"
      title={t("services.wafv2.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.wafv2.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.wafv2.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "wafv2-items" }}
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
        title={t("services.wafv2.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.wafv2.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.wafv2.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.wafv2.scopeLabel")}
          <select value={formScope} onChange={(e) => setFormScope(Number(e.target.value) as Scope)} className="modal-input">
            <option value={Scope.REGIONAL}>{t("services.wafv2.scopeRegional")}</option>
            <option value={Scope.CLOUDFRONT}>{t("services.wafv2.scopeCloudfront")}</option>
          </select>
        </label>
        <label>
          {t("services.wafv2.descriptionLabel")}
          <input
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder={t("common.optional")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.wafv2.defaultActionLabel")}
          <select value={formDefaultAction} onChange={(e) => setFormDefaultAction(e.target.value as "allow" | "block")} className="modal-input">
            <option value="allow">{t("services.wafv2.actionAllow")}</option>
            <option value="block">{t("services.wafv2.actionBlock")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.wafv2.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
