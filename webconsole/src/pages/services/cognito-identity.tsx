/**
 * Cognito Identity service page. Lists identity pools with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import {
  CognitoIdentityService,
  type IdentityPoolShortDescription,
} from "@/gen/cognitoidentity_pb";
import { CreateIdentityPoolInputSchema } from "@/gen/cognitoidentity_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the Cognito Identity pool table. */
const getColumns = (t: TFunction): ColumnDef<IdentityPoolShortDescription, any>[] => [
  { accessorKey: "identitypoolname", header: t("services.cognitoidentity.poolNameHeader"), cell: MonoCell },
  { accessorKey: "identitypoolid", header: t("services.cognitoidentity.poolIdHeader"), cell: MonoCell },
];

/** Cognito Identity service page with list, create, and delete operations. */
export function CognitoIdentityPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<IdentityPoolShortDescription | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formAllowUnauth, setFormAllowUnauth] = useState(true);
  const [formProviderName, setFormProviderName] = useState("");
  const [formProviderClientId, setFormProviderClientId] = useState("");
  const [formTags, setFormTags] = useState("");

  const { client, invalidate } = useServiceClient(CognitoIdentityService);
  const { queryKey } = useListKey("cognito-identity");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listIdentityPools({ maxresults: 60 }),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: IdentityPoolShortDescription[] = dropEmpty(data?.identitypools ?? [], "identitypoolid");

  const createMutation = useMutation({
    mutationFn: () => {
      let identitypooltags: Record<string, string> = {};
      if (formTags.trim()) {
        try {
          identitypooltags = JSON.parse(formTags);
        } catch {
          throw new Error(t("common.invalidJson", { field: "Tags" }));
        }
      }
      return client.createIdentityPool(
        create(CreateIdentityPoolInputSchema, {
          identitypoolname: formName,
          allowunauthenticatedidentities: formAllowUnauth,
          ...(formProviderName ? {
            cognitoidentityproviders: [{ providername: formProviderName, clientid: formProviderClientId }],
          } : {}),
          ...(Object.keys(identitypooltags).length > 0 ? { identitypooltags } : {}),
        }),
      );
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormAllowUnauth(true);
      setFormProviderName("");
      setFormProviderClientId("");
      setFormTags("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (poolId: string) =>
      client.deleteIdentityPool({ identitypoolid: poolId }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="👤"
      title={t("services.cognitoidentity.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cognitoidentity.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cognitoidentity.create")}
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
        getRowId={(row) => row.identitypoolid}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.identitypoolid}
        selected={selectedItem}
        detailTitle={selectedItem?.identitypoolname}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.cognitoidentity.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.cognitoidentity.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.cognitoidentity.placeholder")}
            className="modal-input"
          />
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formAllowUnauth}
            onChange={(e) => setFormAllowUnauth(e.target.checked)}
          />
          {t("services.cognitoidentity.allowUnauthenticatedLabel")}
        </label>
        <label>
          {t("services.cognitoidentity.providerNameLabel")}
          <input
            value={formProviderName}
            onChange={(e) => setFormProviderName(e.target.value)}
            placeholder={t("services.cognitoidentity.providerNamePlaceholder")}
            className="modal-input"
          />
        </label>
        {formProviderName && (
          <label>
            {t("services.cognitoidentity.providerClientIdLabel")}
            <input
              value={formProviderClientId}
              onChange={(e) => setFormProviderClientId(e.target.value)}
              placeholder={t("services.cognitoidentity.providerClientIdPlaceholder")}
              className="modal-input"
            />
          </label>
        )}
        <label>
          {t("services.cognitoidentity.tagsLabel")}
          <textarea
            value={formTags}
            onChange={(e) => setFormTags(e.target.value)}
            placeholder='{"key":"value"}'
            rows={3}
            className="modal-input"
            style={{ fontFamily: "monospace", fontSize: "0.85em" }}
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cognitoidentity.delete")}
        name={selectedItem?.identitypoolname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.identitypoolid)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
