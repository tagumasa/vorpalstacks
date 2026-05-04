/**
 * Cognito IDP service page. Lists user pools with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { CognitoIdentityProviderService, VerifiedAttributeType } from "@/gen/cognitoidentityprovider_pb";
import { CreateUserPoolRequestSchema, PasswordPolicyTypeSchema, UserPoolPolicyTypeSchema } from "@/gen/cognitoidentityprovider_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  BadgeCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Derived row shape for the Cognito user pool list table. */
interface TableRow {
  name: string;
  id: string;
  status: string;
  creationdate: string;
  lastmodifieddate: string;
}

/** Column definitions for the Cognito user pool table. */
const getColumns = (t: TFunction): ColumnDef<TableRow, any>[] => [
  { accessorKey: "name", header: t("services.cognito.poolNameHeader"), cell: MonoCell },
  { accessorKey: "id", header: t("services.cognito.poolIdHeader"), cell: SmallMonoCell },
  { accessorKey: "status", header: t("services.cognito.statusHeader"), cell: ({ getValue }) => <BadgeCell getValue={getValue} positive={["Enabled"]} negative={["Disabled"]} />, size: 90 },
  { accessorKey: "creationdate", header: t("services.cognito.creationDateHeader"), cell: DateCell },
  { accessorKey: "lastmodifieddate", header: t("services.cognito.creationDateHeader"), cell: DateCell },
];

/** Cognito IDP service page with list, create, and delete operations. */
export function CognitoPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<TableRow | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formAutoVerifyEmail, setFormAutoVerifyEmail] = useState(false);
  const [formAutoVerifyPhone, setFormAutoVerifyPhone] = useState(false);
  const [formPasswordMinLength, setFormPasswordMinLength] = useState(8);
  const [formTempPasswordDays, setFormTempPasswordDays] = useState(7);

  const { client, invalidate } = useServiceClient(CognitoIdentityProviderService);
  const { queryKey } = useListKey("cognito");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listUserPools({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: TableRow[] = dropEmpty(
    (data?.userpools ?? []).map((pool) => ({
      name: pool.name,
      id: pool.id,
      status: String(pool.status ?? ""),
      creationdate: pool.creationdate ?? "",
      lastmodifieddate: pool.lastmodifieddate ?? "",
    })),
    "name",
  );

  const createMutation = useMutation({
    mutationFn: () => {
      const autoVerified: VerifiedAttributeType[] = [];
      if (formAutoVerifyEmail) autoVerified.push(VerifiedAttributeType.EMAIL);
      if (formAutoVerifyPhone) autoVerified.push(VerifiedAttributeType.PHONE_NUMBER);
      const passwordpolicy = create(PasswordPolicyTypeSchema, {
        minimumlength: formPasswordMinLength,
        temporarypasswordvaliditydays: formTempPasswordDays,
      });
      const policies = create(UserPoolPolicyTypeSchema, { passwordpolicy });
      return client.createUserPool(
        create(CreateUserPoolRequestSchema, {
          poolname: formName,
          ...(autoVerified.length > 0 ? { autoverifiedattributes: autoVerified } : {}),
          policies,
        }),
      );
    },
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormName("");
      setFormAutoVerifyEmail(false);
      setFormAutoVerifyPhone(false);
      setFormPasswordMinLength(8);
      setFormTempPasswordDays(7);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (userpoolid: string) =>
      client.deleteUserPool({ userpoolid }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="👤"
      title={t("services.cognito.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.cognito.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.cognito.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("common.delete")}
            </button>
          )}
        </>
      }
      exportData={{ rows: items as unknown as Record<string, unknown>[], columns, filenamePrefix: "cognito-items" }}
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
        title={t("services.cognito.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formName}
      >
        <label>
          {t("services.cognito.nameField")}
          <input
            value={formName}
            onChange={(e) => setFormName(e.target.value)}
            placeholder={t("services.cognito.placeholder")}
            className="modal-input"
          />
        </label>
        <div style={{ marginTop: 8 }}>
          <span style={{ fontSize: 13, fontWeight: 500, display: "block", marginBottom: 4 }}>
            {t("services.cognito.autoVerifiedAttributesLabel")}
          </span>
          <label style={{ display: "flex", alignItems: "center", gap: 6, marginBottom: 4 }}>
            <input
              type="checkbox"
              checked={formAutoVerifyEmail}
              onChange={(e) => setFormAutoVerifyEmail(e.target.checked)}
            />
            {t("services.cognito.autoVerifiedEmailLabel")}
          </label>
          <label style={{ display: "flex", alignItems: "center", gap: 6 }}>
            <input
              type="checkbox"
              checked={formAutoVerifyPhone}
              onChange={(e) => setFormAutoVerifyPhone(e.target.checked)}
            />
            {t("services.cognito.autoVerifiedPhoneLabel")}
          </label>
        </div>
        <div className="form-row">
          <label>
            {t("services.cognito.passwordMinLengthLabel")}
            <input
              type="number"
              min={6}
              max={99}
              value={formPasswordMinLength}
              onChange={(e) => setFormPasswordMinLength(Number(e.target.value))}
              className="modal-input"
            />
          </label>
          <label>
            {t("services.cognito.tempPasswordDaysLabel")}
            <input
              type="number"
              min={1}
              max={365}
              value={formTempPasswordDays}
              onChange={(e) => setFormTempPasswordDays(Number(e.target.value))}
              className="modal-input"
            />
          </label>
        </div>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.cognito.delete")}
        name={selectedItem?.name}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.id)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
