/**
 * KMS service page. Lists keys with create/schedule-deletion operations.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { create } from "@bufbuild/protobuf";
import { KMSService, type KeyListEntry, KeyUsageType, KeySpec, OriginType } from "@/gen/kms_pb";
import { CreateKeyRequestSchema, ScheduleKeyDeletionRequestSchema } from "@/gen/kms_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Column definitions for the KMS key table. */
const getColumns = (t: TFunction): ColumnDef<KeyListEntry, any>[] => [
  { accessorKey: "keyid", header: t("services.kms.keyIdHeader"), cell: MonoCell },
  { accessorKey: "keyarn", header: t("services.kms.arnHeader"), cell: SmallMonoCell },
];

/** Key specification options for the create key form. */
const KEY_SPECS: { i18nKey: string; value: KeySpec }[] = [
  { i18nKey: "services.kms.keySpecSymmetricDefault", value: KeySpec.SYMMETRIC_DEFAULT },
  { i18nKey: "services.kms.keySpecRsa2048", value: KeySpec.RSA_2048 },
  { i18nKey: "services.kms.keySpecRsa3072", value: KeySpec.RSA_3072 },
  { i18nKey: "services.kms.keySpecRsa4096", value: KeySpec.RSA_4096 },
  { i18nKey: "services.kms.keySpecEccNistP256", value: KeySpec.ECC_NIST_P256 },
  { i18nKey: "services.kms.keySpecEccNistP384", value: KeySpec.ECC_NIST_P384 },
  { i18nKey: "services.kms.keySpecHmac256", value: KeySpec.HMAC_256 },
];

/** Key usage options for the create key form. */
const KEY_USAGES: { i18nKey: string; value: KeyUsageType }[] = [
  { i18nKey: "services.kms.keyUsageEncryptDecrypt", value: KeyUsageType.ENCRYPT_DECRYPT },
  { i18nKey: "services.kms.keyUsageSignVerify", value: KeyUsageType.SIGN_VERIFY },
  { i18nKey: "services.kms.keyUsageGenerateVerifyMac", value: KeyUsageType.GENERATE_VERIFY_MAC },
];

/** Origin options for the create key form. */
const ORIGINS: { i18nKey: string; value: OriginType }[] = [
  { i18nKey: "services.kms.originAwsKms", value: OriginType.AWS_KMS },
  { i18nKey: "services.kms.originExternal", value: OriginType.EXTERNAL },
];

/** KMS service page with list, create, and schedule-deletion operations. */
export function KMSPage() {
  const { t } = useTranslation();
  const columns = getColumns(t);
  const [selectedItem, setSelectedItem] = useState<KeyListEntry | null>(null);
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formDesc, setFormDesc] = useState("");
  const [formKeySpec, setFormKeySpec] = useState(KeySpec.SYMMETRIC_DEFAULT);
  const [formKeyUsage, setFormKeyUsage] = useState(KeyUsageType.ENCRYPT_DECRYPT);
  const [formOrigin, setFormOrigin] = useState(OriginType.AWS_KMS);
  const [formMultiRegion, setFormMultiRegion] = useState(false);

  const { client, invalidate } = useServiceClient(KMSService);
  const { queryKey } = useListKey("kms");

  const { data, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listKeys({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const items: KeyListEntry[] = dropEmpty(data?.keys ?? [], "keyid");

  const createMutation = useMutation({
    mutationFn: () =>
      client.createKey(
        create(CreateKeyRequestSchema, {
          description: formDesc,
          keyspec: formKeySpec,
          keyusage: formKeyUsage,
          origin: formOrigin,
          multiregion: formMultiRegion,
        }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormDesc("");
      setFormKeySpec(KeySpec.SYMMETRIC_DEFAULT);
      setFormKeyUsage(KeyUsageType.ENCRYPT_DECRYPT);
      setFormOrigin(OriginType.AWS_KMS);
      setFormMultiRegion(false);
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (keyId: string) =>
      client.scheduleKeyDeletion(
        create(ScheduleKeyDeletionRequestSchema, { keyid: keyId, pendingwindowindays: 7 }),
      ),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDelete(false);
      setSelectedItem(null);
    },
  });

  return (
    <ServicePageLayout
      icon="🔑"
      title={t("services.kms.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("services.kms.countLabel")}
      actions={
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.kms.create")}
          </button>
          {selectedItem && (
            <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
              {t("services.kms.delete")}
            </button>
          )}
        </>
      }
    >
      <SplitPane
        columns={columns}
        data={items}
        getRowId={(row) => row.keyid}
        onRowClick={setSelectedItem}
        selectedId={selectedItem?.keyid}
        selected={selectedItem}
        detailTitle={selectedItem?.keyid}
        onDetailClose={() => setSelectedItem(null)}
      />

      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.kms.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={createMutation.isPending}
      >
        <label>
          {t("services.kms.descLabel")}
          <input
            value={formDesc}
            onChange={(e) => setFormDesc(e.target.value)}
            placeholder={t("services.kms.descPlaceholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.kms.keyUsageLabel")}
          <select
            value={formKeyUsage}
            onChange={(e) => setFormKeyUsage(Number(e.target.value))}
            className="modal-input"
          >
            {KEY_USAGES.map((u) => (
              <option key={u.value} value={u.value}>{t(u.i18nKey)}</option>
            ))}
          </select>
        </label>
        <label>
          {t("services.kms.keySpecLabel")}
          <select
            value={formKeySpec}
            onChange={(e) => setFormKeySpec(Number(e.target.value))}
            className="modal-input"
          >
            {KEY_SPECS.map((s) => (
              <option key={s.value} value={s.value}>{t(s.i18nKey)}</option>
            ))}
          </select>
        </label>
        <label>
          {t("services.kms.originLabel")}
          <select
            value={formOrigin}
            onChange={(e) => setFormOrigin(Number(e.target.value))}
            className="modal-input"
          >
            {ORIGINS.map((o) => (
              <option key={o.value} value={o.value}>{t(o.i18nKey)}</option>
            ))}
          </select>
        </label>
        <label className="checkbox-label">
          <input
            type="checkbox"
            checked={formMultiRegion}
            onChange={(e) => setFormMultiRegion(e.target.checked)}
          />
          {t("services.kms.multiregionLabel")}
        </label>
      </ServiceCreateModal>

      <ServiceDeleteDialog
        open={showDelete && !!selectedItem}
        title={t("services.kms.delete")}
        name={selectedItem?.keyid}
        error={deleteMutation.error}
        isPending={deleteMutation.isPending}
        onConfirm={() => selectedItem && deleteMutation.mutate(selectedItem.keyid)}
        onClose={() => setShowDelete(false)}
      />
    </ServicePageLayout>
  );
}
