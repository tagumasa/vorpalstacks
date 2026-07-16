/**
 * SESv2 service page — 3-panel inspector layout.
 *
 * Panel 1 (toolbar): Breadcrumb navigation
 * Panel 2 (table):   Email identity list with checkbox multi-select
 * Panel 3 (detail):  Identity detail (type, verification status, sending enabled)
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { SESv2Service, type IdentityInfo, IdentityType, VerificationStatus } from "@/gen/sesv2_pb";
import { CreateEmailIdentityRequestSchema } from "@/gen/sesv2_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  BooleanBadge,
  useServiceClient,
} from "@/components/shared/service-page";
import {
  checkboxColumn,
  Breadcrumb,
  SelectionBadge,
  DetailPanel,
  DetailEmpty,
  useSelection,
} from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

// ─── Helpers ────────────────────────────────────────────────────

const IDENTITY_TYPE_I18N: Record<number, string> = {
  [IdentityType.MANAGED_DOMAIN]: "services.ses.typeManagedDomain",
  [IdentityType.DOMAIN]: "services.ses.typeDomain",
  [IdentityType.EMAIL_ADDRESS]: "services.ses.typeEmailAddress",
};

const VERIFICATION_STATUS_I18N: Record<number, string> = {
  [VerificationStatus.PENDING]: "services.ses.statusPending",
  [VerificationStatus.SUCCESS]: "services.ses.statusSuccess",
  [VerificationStatus.TEMPORARY_FAILURE]: "services.ses.statusTemporaryFailure",
  [VerificationStatus.FAILED]: "services.ses.statusFailed",
  [VerificationStatus.NOT_STARTED]: "services.ses.statusNotStarted",
};

// ─── Column Definitions ─────────────────────────────────────────

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

// ─── Detail panel tab ───────────────────────────────────────────

type DetailTab = "detail" | "json";

// ─── SESv2 Page ─────────────────────────────────────────────────

export function SESPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(SESv2Service);
  const { queryKey } = useListKey("ses");
  const columns = getColumns(t);

  // ── Selection state ──────────────────────────────────────────
  const { selected: selectedNames, toggle: toggleName, toggleAll: toggleAllNames, clear: clearSelection } = useSelection<string>();
  const [selectedItem, setSelectedItem] = useState<IdentityInfo | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [showBatchDelete, setShowBatchDelete] = useState(false);

  // ── Create form state ────────────────────────────────────────
  const [formIdentity, setFormIdentity] = useState("");
  const [formConfigSetName, setFormConfigSetName] = useState("");

  // ── Data ─────────────────────────────────────────────────────
  const { items: rawItems, hasMore, loadMore, isFetchingMore, isLoading, error, invalidate: invalidateList } = usePaginatedList<IdentityInfo, Awaited<ReturnType<typeof client.listEmailIdentities>>>({
    queryKeyBase: queryKey,
    fetchPage: (token) => client.listEmailIdentities({ nexttoken: token || undefined }),
    getItems: (r) => r.emailidentities ?? [],
    getNextToken: (r) => r.nexttoken ?? "",
  });
  const items = dropEmpty(rawItems, "identityname");

  // ── Mutations ────────────────────────────────────────────────
  const createMutation = useMutation({
    mutationFn: () =>
      client.createEmailIdentity(
        create(CreateEmailIdentityRequestSchema, {
          emailidentity: formIdentity,
          ...(formConfigSetName ? { configurationsetname: formConfigSetName } : {}),
        }),
      ),
    onSuccess: () => {
      invalidateList();
      setShowCreate(false);
      setFormIdentity("");
      setFormConfigSetName("");
    },
  });

  const deleteMutation = useMutation({
    mutationFn: (identity: string) =>
      client.deleteEmailIdentity({ emailidentity: identity }),
    onSuccess: () => {
      invalidateList();
      setShowDelete(false);
      setSelectedItem(null);
      clearSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: async (identities: string[]) => {
      const results = await Promise.allSettled(
        identities.map((id) => client.deleteEmailIdentity({ emailidentity: id }))
      );
      return results;
    },
    onSuccess: (_data, identities) => {
      invalidateList();
      setShowBatchDelete(false);
      clearSelection();
      setSelectedItem((prev) => (prev && identities.includes(prev.identityname) ? null : prev));
    },
  });

  // ── Handlers ─────────────────────────────────────────────────

  const handleRowClick = (row: IdentityInfo) => {
    setSelectedItem(row);
    setDetailTab("detail");
  };

  const allIds = items.map((i) => i.identityname);

  // ── Detail Panel ─────────────────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedItem) {
      return <DetailEmpty message={t("common.noItemSelected")} />;
    }

    const typeLabel = IDENTITY_TYPE_I18N[selectedItem.identitytype] ? t(IDENTITY_TYPE_I18N[selectedItem.identitytype]!) : String(selectedItem.identitytype);
    const statusLabel = VERIFICATION_STATUS_I18N[selectedItem.verificationstatus] ? t(VERIFICATION_STATUS_I18N[selectedItem.verificationstatus]!) : String(selectedItem.verificationstatus);

    const detailTabs = [
      { key: "detail", label: t("common.tabDetail") },
      { key: "json", label: t("common.rawJson") },
    ];

    return (
      <DetailPanel
        title={selectedItem.identityname}
        titleIcon="✉️"
        tabs={detailTabs}
        activeTab={detailTab}
        onTabChange={(k) => setDetailTab(k as DetailTab)}
        actions={
          <button className="btn btn-danger btn-sm" onClick={() => setShowDelete(true)}>
            {t("common.delete")}
          </button>
        }
      >
        {detailTab === "detail" ? (
          <table className="settings-table">
            <tbody>
              <tr><td className="detail-label-fixed">{t("common.fields.identity")}</td><td className="cell-mono">{selectedItem.identityname}</td></tr>
              <tr><td className="detail-label">{t("common.fields.type")}</td><td><span className="badge">{typeLabel}</span></td></tr>
              <tr><td className="detail-label">{t("common.fields.verification")}</td><td><span className="badge">{statusLabel}</span></td></tr>
              <tr><td className="detail-label">{t("common.fields.sendingEnabled")}</td><td><BooleanBadge value={selectedItem.sendingenabled} /></td></tr>
            </tbody>
          </table>
        ) : (
          <JsonViewer data={selectedItem} />
        )}
      </DetailPanel>
    );
  };

  // ── Render ───────────────────────────────────────────────────

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
          <button
            className="btn btn-danger"
            disabled={selectedNames.size === 0}
            onClick={() => setShowBatchDelete(true)}
          >
            {t("common.deleteSelected")}
            {selectedNames.size > 0 && <span className="batch-count">({selectedNames.size})</span>}
          </button>
        </>
      }
    >
      <div className="inspector-toolbar">
        <Breadcrumb parts={[
          { label: t("services.ses.title") },
          { label: t("services.ses.countLabel") },
        ]} />
        <div className="toolbar-selection-info">
          <SelectionBadge count={selectedNames.size} label={t("common.selectedCount", { count: selectedNames.size })} />
        </div>
      </div>

      {items.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-ses">
          <div className="flex-fill-scroll">
            <DataTable
              columns={[
                checkboxColumn<IdentityInfo>(selectedNames, toggleName, () => toggleAllNames(allIds), allIds, t, (row) => row.identityname),
                ...columns,
              ]}
              data={items}
              getRowId={(row) => row.identityname}
              onRowClick={handleRowClick}
              selectedId={selectedItem?.identityname}
              hasMore={hasMore} onLoadMore={loadMore} loadingMore={isFetchingMore}
            />
          </div>
          {renderDetailPanel()}
        </Splitter>
      ) : (
        <div className="empty-state">{t("common.noData")}</div>
      )}

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
          <input value={formIdentity} onChange={(e) => setFormIdentity(e.target.value)} placeholder={t("services.ses.placeholder")} className="modal-input" />
        </label>
        <label>
          {t("services.ses.configurationSetLabel")}
          <input value={formConfigSetName} onChange={(e) => setFormConfigSetName(e.target.value)} placeholder={t("services.ses.configurationSetPlaceholder")} className="modal-input" />
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

      <ServiceDeleteDialog
        open={showBatchDelete}
        title={t("common.deleteSelected")}
        name={`${selectedNames.size} ${t("services.ses.countLabel")}`}
        error={batchDeleteMutation.error}
        isPending={batchDeleteMutation.isPending}
        onConfirm={() => batchDeleteMutation.mutate(Array.from(selectedNames))}
        onClose={() => setShowBatchDelete(false)}
      />
    </ServicePageLayout>
  );
}
