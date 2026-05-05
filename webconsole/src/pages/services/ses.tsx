/**
 * SESv2 service page. Lists email identities with create/delete operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { SESv2Service, type IdentityInfo, IdentityType, VerificationStatus } from "@/gen/sesv2_pb";
import { CreateEmailIdentityRequestSchema } from "@/gen/sesv2_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  BooleanBadge,
  useServiceClient,
} from "@/components/shared/service-page";

/** Lookup map for IdentityType proto enum values to i18n keys. */
const IDENTITY_TYPE_I18N: Record<number, string> = {
  [IdentityType.MANAGED_DOMAIN]: "services.ses.typeManagedDomain",
  [IdentityType.DOMAIN]: "services.ses.typeDomain",
  [IdentityType.EMAIL_ADDRESS]: "services.ses.typeEmailAddress",
};

/** Lookup map for VerificationStatus proto enum values to i18n keys. */
const VERIFICATION_STATUS_I18N: Record<number, string> = {
  [VerificationStatus.PENDING]: "services.ses.statusPending",
  [VerificationStatus.SUCCESS]: "services.ses.statusSuccess",
  [VerificationStatus.TEMPORARY_FAILURE]: "services.ses.statusTemporaryFailure",
  [VerificationStatus.FAILED]: "services.ses.statusFailed",
  [VerificationStatus.NOT_STARTED]: "services.ses.statusNotStarted",
};

/** Column definitions for the SES email identity table. */
const getColumns = (t: TFunction): ColumnDef<IdentityInfo, any>[] => [
  { accessorKey: "identityname", header: t("services.ses.identityHeader"), cell: MonoCell },
  {
    accessorKey: "identitytype",
    header: t("services.ses.typeHeader"),
    cell: ({ getValue }) => {
      const v = getValue() as number;
      return <span className="badge">{IDENTITY_TYPE_I18N[v] ? t(IDENTITY_TYPE_I18N[v]) : String(v)}</span>;
    },
    size: 90,
  },
  {
    accessorKey: "verificationstatus",
    header: t("services.ses.verificationHeader"),
    cell: ({ getValue }) => {
      const v = getValue() as number;
      return <span className="badge">{VERIFICATION_STATUS_I18N[v] ? t(VERIFICATION_STATUS_I18N[v]) : String(v)}</span>;
    },
    size: 110,
  },
  {
    accessorKey: "sendingenabled",
    header: t("services.ses.sendingHeader"),
    cell: ({ getValue }) => <BooleanBadge value={getValue() as boolean} />,
    size: 90,
  },
];

/** SESv2 service page with list, create, and delete operations. */
export function SESPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<IdentityInfo | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formIdentity, setFormIdentity] = useState("");
  const [formConfigSetName, setFormConfigSetName] = useState("");

  const { client, invalidate } = useServiceClient(SESv2Service);
  const { queryKey } = useListKey("ses");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listEmailIdentities({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: IdentityInfo[] = dropEmpty(data?.emailidentities ?? [], "identityname");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createEmailIdentity(
        create(CreateEmailIdentityRequestSchema, {
          emailidentity: formIdentity,
          ...(formConfigSetName ? { configurationsetname: formConfigSetName } : {}),
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormIdentity("");
      setFormConfigSetName("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (identity: string) =>
      client.deleteEmailIdentity({ emailidentity: identity }),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="✉️"
      title={t("services.ses.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.ses.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.ses.create")}
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
        getRowId={(row) => row.identityname}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.identityname}
        selected={selectedItem}
        detailTitle={selectedItem?.identityname}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.ses.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formIdentity}
      >
        <label>
          {t("services.ses.nameField")}
          <input
            value={formIdentity}
            onChange={(e) => setFormIdentity(e.target.value)}
            placeholder={t("services.ses.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.ses.configurationSetLabel")}
          <input
            value={formConfigSetName}
            onChange={(e) => setFormConfigSetName(e.target.value)}
            placeholder={t("services.ses.configurationSetPlaceholder")}
            className="modal-input"
          />
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.ses.delete")}
        name={selectedItem?.identityname}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.identityname)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
