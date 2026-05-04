/**
 * AppSync service page. Lists GraphQL APIs with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { AppSyncService, type GraphqlApi, AuthenticationType } from "@/gen/appsync_pb";
import { CreateGraphqlApiRequestSchema, DeleteGraphqlApiRequestSchema } from "@/gen/appsync_pb";
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

/** Column definitions for the AppSync GraphQL API table. */
const getColumns = (t: TFunction): ColumnDef<GraphqlApi, any>[] => [
  { accessorKey: "name", header: t("services.appsync.apiNameHeader"), cell: MonoCell },
  { accessorKey: "apiid", header: t("services.appsync.apiIdHeader"), cell: MonoCell, size: 100 },
  { accessorKey: "authenticationtype", header: t("services.appsync.authTypeHeader"), cell: FallbackCell, size: 120 },
  { accessorKey: "apitype", header: t("services.appsync.apiTypeHeader"), cell: FallbackCell, size: 90 },
  { accessorKey: "arn", header: t("services.appsync.arnHeader"), cell: SmallMonoCell },
];

/** Authentication type options for the create API form. */
const AUTH_TYPES = [
  { label: "API_KEY", value: AuthenticationType.API_KEY },
  { label: "AWS_IAM", value: AuthenticationType.AWS_IAM },
  { label: "AMAZON_COGNITO_USER_POOLS", value: AuthenticationType.AMAZON_COGNITO_USER_POOLS },
  { label: "OPENID_CONNECT", value: AuthenticationType.OPENID_CONNECT },
  { label: "AWS_LAMBDA", value: AuthenticationType.AWS_LAMBDA },
];

/** AppSync service page with list, create, and delete operations. */
export function AppSyncPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<GraphqlApi | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formAuthType, setFormAuthType] = useState(AuthenticationType.API_KEY);

  const { client, invalidate } = useServiceClient(AppSyncService);
  const { queryKey } = useListKey("appsync");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listGraphqlApis({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: GraphqlApi[] = dropEmpty(data?.graphqlapis ?? [], "name");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createGraphqlApi(
        create(CreateGraphqlApiRequestSchema, {
          name: formName,
          authenticationtype: formAuthType,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (apiId: string) =>
      client.deleteGraphqlApi(
        create(DeleteGraphqlApiRequestSchema, { apiid: apiId }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🔗"
      title={t("services.appsync.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.appsync.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.appsync.create")}
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
        getRowId={(row) => row.apiid}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.apiid}
        selected={selectedItem}
        detailTitle={selectedItem?.name}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.appsync.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.appsync.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.appsync.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.appsync.authTypeLabel")}
          <select
            value={formAuthType}
            onChange={(e) => setFormAuthType(Number(e.target.value))}
            className="modal-input"
          >
            {AUTH_TYPES.map((a) => (
              <option key={a.value} value={a.value}>{a.label}</option>
            ))}
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.appsync.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.apiid)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
