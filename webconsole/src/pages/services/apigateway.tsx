/**
 * API Gateway service page. Lists REST APIs with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { APIGatewayService, type RestApi, EndpointType } from "@/gen/apigateway_pb";
import { CreateRestApiRequestSchema, EndpointConfigurationSchema } from "@/gen/apigateway_pb";
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

/** Column definitions for the API Gateway REST API table. */
const getColumns = (t: TFunction): ColumnDef<RestApi, any>[] => [
  { accessorKey: "name", header: t("services.apigateway.apiNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.apigateway.apiIdHeader"), cell: MonoCell, size: 100 },
  { accessorKey: "apistatus", header: t("services.apigateway.apiStatusHeader"), cell: FallbackCell, size: 90 },
  { accessorKey: "description", header: t("services.apigateway.descriptionHeader"), cell: FallbackCell },
  { accessorKey: "createddate", header: t("services.apigateway.createdHeader"), cell: DateCell },
];

/** API Gateway service page with list, create, and delete operations. */
export function APIGatewayPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<RestApi | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formDescription, setFormDescription] = useState("");
  const [formVersion, setFormVersion] = useState("1.0");
  const [formEndpointType, setFormEndpointType] = useState<string>("REGIONAL");

  const { client, invalidate } = useServiceClient(APIGatewayService);
  const { queryKey } = useListKey("apigateway");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.getRestApis({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: RestApi[] = dropEmpty(data?.items ?? [], "id");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createRestApi(
        create(CreateRestApiRequestSchema, {
          name: formName,
          description: formDescription,
          version: formVersion,
          endpointconfiguration: create(EndpointConfigurationSchema, {
            types: [EndpointType[formEndpointType as keyof typeof EndpointType] ?? EndpointType.REGIONAL],
          }),
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormDescription("");
      setFormVersion("1.0");
      setFormEndpointType("REGIONAL");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (apiId: string) =>
      client.deleteRestApi({ restapiid: apiId }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🌉"
      title={t("services.apigateway.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.apigateway.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.apigateway.create")}
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
        title={t("services.apigateway.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.apigateway.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.apigateway.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.apigateway.descriptionLabel")}
          <input
            value={formDescription}
            onChange={(e) => setFormDescription(e.target.value)}
            placeholder={t("services.apigateway.descPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.apigateway.versionLabel")}
          <input
            value={formVersion}
            onChange={(e) => setFormVersion(e.target.value)}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.apigateway.endpointTypeLabel")}
          <select
            value={formEndpointType}
            onChange={(e) => setFormEndpointType(e.target.value)}
            className="modal-input"
          >
            <option value="REGIONAL">{t("services.apigateway.endpointRegional")}</option>
            <option value="EDGE">{t("services.apigateway.endpointEdge")}</option>
            <option value="PRIVATE">{t("services.apigateway.endpointPrivate")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.apigateway.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
