/**
 * IAM service page — 3-panel inspector layout with tabs (users, roles, policies).
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { IAMService, type User, type Role, type Policy, CreateUserRequestSchema, CreateRoleRequestSchema } from "@/gen/iam_pb";
import { useListKey, dropEmpty, usePaginatedList } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, FallbackCell, fmtDate, useServiceClient } from "@/components/shared/service-page";
import { checkboxColumn, DetailPanel, DetailEmpty, useSelection } from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "users" | "roles" | "policies";
type DetailTab = "detail" | "json";

const getUserColumns = (t: TFunction): ColumnDef<User, any>[] => [
  { accessorKey: "username", header: t("services.iam.userNameHeader"), cell: MonoCell },
  { accessorKey: "path", header: t("services.iam.pathHeader"), cell: FallbackCell, size: 90 },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "createdate", header: t("services.iam.createdHeader"), cell: DateCell },
];

const getRoleColumns = (t: TFunction): ColumnDef<Role, any>[] => [
  { accessorKey: "rolename", header: t("services.iam.roleNameHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.iam.roleDescriptionHeader"), cell: FallbackCell },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "createdate", header: t("services.iam.createdHeader"), cell: DateCell },
];

const getPolicyColumns = (t: TFunction): ColumnDef<Policy, any>[] => [
  { accessorKey: "policyname", header: t("services.iam.policyNameHeader"), cell: MonoCell },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "attachmentcount", header: t("services.iam.attachmentsHeader"), size: 100 },
  { accessorKey: "createdate", header: t("services.iam.policyCreatedHeader"), cell: DateCell },
];

export function IAMPage() {
  const { t, i18n } = useTranslation();
  const userColumns = getUserColumns(t);
  const roleColumns = getRoleColumns(t);
  const policyColumns = getPolicyColumns(t);

  const [tab, setTab] = useState<TabKey>("users");
  const userSel = useSelection<string>();
  const roleSel = useSelection<string>();

  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedPolicy, setSelectedPolicy] = useState<Policy | null>(null);
  const [detailTab, setDetailTab] = useState<DetailTab>("detail");

  const [showCreateUser, setShowCreateUser] = useState(false);
  const [showDeleteUser, setShowDeleteUser] = useState(false);
  const [showBatchDeleteUser, setShowBatchDeleteUser] = useState(false);
  const [showCreateRole, setShowCreateRole] = useState(false);
  const [showDeleteRole, setShowDeleteRole] = useState(false);
  const [showBatchDeleteRole, setShowBatchDeleteRole] = useState(false);

  const [formUserName, setFormUserName] = useState("");
  const [formUserPath, setFormUserPath] = useState("");
  const [formRoleName, setFormRoleName] = useState("");
  const [formRolePath, setFormRolePath] = useState("");
  const [formAssumePolicy, setFormAssumePolicy] = useState(JSON.stringify({ Version: "2012-10-17", Statement: [{ Effect: "Allow", Principal: { Service: "lambda.amazonaws.com" }, Action: "sts:AssumeRole" }] }, null, 2));

  const { client } = useServiceClient(IAMService);
  const { queryKey: usersKey } = useListKey("iam-users");
  const { queryKey: rolesKey } = useListKey("iam-roles");
  const { queryKey: policiesKey } = useListKey("iam-policies");

  const usersList = usePaginatedList<User, Awaited<ReturnType<typeof client.listUsers>>>({
    queryKeyBase: usersKey, fetchPage: (token) => client.listUsers({ marker: token || undefined }), getItems: (r) => r.users ?? [], getNextToken: (r) => r.marker ?? "",
  });
  const rolesList = usePaginatedList<Role, Awaited<ReturnType<typeof client.listRoles>>>({
    queryKeyBase: rolesKey, fetchPage: (token) => client.listRoles({ marker: token || undefined }), getItems: (r) => r.roles ?? [], getNextToken: (r) => r.marker ?? "",
  });
  const policiesList = usePaginatedList<Policy, Awaited<ReturnType<typeof client.listPolicies>>>({
    queryKeyBase: policiesKey, fetchPage: (token) => client.listPolicies({ marker: token || undefined }), getItems: (r) => r.policies ?? [], getNextToken: (r) => r.marker ?? "",
  });

  const users = dropEmpty(usersList.items, "username");
  const roles = dropEmpty(rolesList.items, "rolename");
  const policies = dropEmpty(policiesList.items, "policyname");

  const query = tab === "users" ? usersList : tab === "roles" ? rolesList : policiesList;

  const createUserMutation = useMutation({
    mutationFn: () => client.createUser(create(CreateUserRequestSchema, { username: formUserName, path: formUserPath || undefined })),
    onSuccess: () => { usersList.invalidate(); setShowCreateUser(false); setFormUserName(""); setFormUserPath(""); },
  });
  const deleteUserMutation = useMutation({
    mutationFn: (name: string) => client.deleteUser({ username: name }),
    onSuccess: () => { usersList.invalidate(); setShowDeleteUser(false); setSelectedUser(null); userSel.clear(); },
  });
  const batchDeleteUserMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteUser({ username: n }))),
    onSuccess: (_d, names) => { usersList.invalidate(); setShowBatchDeleteUser(false); userSel.clear(); setSelectedUser((p) => (p && names.includes(p.username) ? null : p)); },
  });

  const createRoleMutation = useMutation({
    mutationFn: () => client.createRole(create(CreateRoleRequestSchema, { rolename: formRoleName, path: formRolePath || undefined, assumerolepolicydocument: formAssumePolicy })),
    onSuccess: () => { rolesList.invalidate(); setShowCreateRole(false); setFormRoleName(""); setFormRolePath(""); },
  });
  const deleteRoleMutation = useMutation({
    mutationFn: (name: string) => client.deleteRole({ rolename: name }),
    onSuccess: () => { rolesList.invalidate(); setShowDeleteRole(false); setSelectedRole(null); roleSel.clear(); },
  });
  const batchDeleteRoleMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteRole({ rolename: n }))),
    onSuccess: (_d, names) => { rolesList.invalidate(); setShowBatchDeleteRole(false); roleSel.clear(); setSelectedRole((p) => (p && names.includes(p.rolename) ? null : p)); },
  });

  const tabs = [
    { key: "users" as TabKey, label: t("services.iam.tabs.users"), count: users.length },
    { key: "roles" as TabKey, label: t("services.iam.tabs.roles"), count: roles.length },
    { key: "policies" as TabKey, label: t("services.iam.tabs.policies"), count: policies.length },
  ];

  const handleTabChange = (k: string) => { setTab(k as TabKey); setSelectedUser(null); setSelectedRole(null); setSelectedPolicy(null); setDetailTab("detail"); userSel.clear(); roleSel.clear(); };

  const renderUserDetail = () => {
    if (!selectedUser) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedUser.username} titleIcon="👤" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteUser(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.username")}</td><td className="cell-mono">{selectedUser.username}</td></tr>
            <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedUser.arn}</td></tr>
            <tr><td className="detail-label">{t("common.fields.path")}</td><td>{selectedUser.path || "/"}</td></tr>
            {selectedUser.createdate && <tr><td className="detail-label">{t("common.fields.created")}</td><td>{fmtDate(selectedUser.createdate, i18n.language)}</td></tr>}
            {selectedUser.passwordlastused && <tr><td className="detail-label">{t("common.fields.lastLogin")}</td><td>{fmtDate(selectedUser.passwordlastused, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedUser} />}
      </DetailPanel>
    );
  };

  const renderRoleDetail = () => {
    if (!selectedRole) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedRole.rolename} titleIcon="🔑" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteRole(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.role")}</td><td className="cell-mono">{selectedRole.rolename}</td></tr>
            <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedRole.arn}</td></tr>
            <tr><td className="detail-label">{t("common.fields.description")}</td><td>{selectedRole.description || "\u2014"}</td></tr>
            <tr><td className="detail-label">{t("common.fields.path")}</td><td>{selectedRole.path || "/"}</td></tr>
            {selectedRole.createdate && <tr><td className="detail-label">{t("common.fields.created")}</td><td>{fmtDate(selectedRole.createdate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedRole} />}
      </DetailPanel>
    );
  };

  const renderPolicyDetail = () => {
    if (!selectedPolicy) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedPolicy.policyname} titleIcon="📜" tabs={[{ key: "detail", label: t("common.tabDetail") }, { key: "json", label: t("common.rawJson") }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)}>
        {detailTab === "detail" ? (
          <table className="settings-table"><tbody>
            <tr><td className="detail-label-fixed">{t("common.fields.policy")}</td><td className="cell-mono">{selectedPolicy.policyname}</td></tr>
            <tr><td className="detail-label">{t("common.fields.arn")}</td><td className="cell-mono cell-long">{selectedPolicy.arn}</td></tr>
            <tr><td className="detail-label">{t("common.fields.attachments")}</td><td>{selectedPolicy.attachmentcount ?? 0}</td></tr>
            {selectedPolicy.createdate && <tr><td className="detail-label">{t("common.fields.created")}</td><td>{fmtDate(selectedPolicy.createdate, i18n.language)}</td></tr>}
            {selectedPolicy.updatedate && <tr><td className="detail-label">{t("common.fields.updated")}</td><td>{fmtDate(selectedPolicy.updatedate, i18n.language)}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedPolicy} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="👤" title={t("services.iam.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange} actions={
      tab === "users" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateUser(true)}>{t("services.iam.createUser")}</button>
        <button className="btn btn-danger" disabled={userSel.selected.size === 0} onClick={() => setShowBatchDeleteUser(true)}>{t("common.deleteSelected")}{userSel.selected.size > 0 && <span className="batch-count">({userSel.selected.size})</span>}</button>
      </>) : tab === "roles" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateRole(true)}>{t("services.iam.createRole")}</button>
        <button className="btn btn-danger" disabled={roleSel.selected.size === 0} onClick={() => setShowBatchDeleteRole(true)}>{t("common.deleteSelected")}{roleSel.selected.size > 0 && <span className="batch-count">({roleSel.selected.size})</span>}</button>
      </>) : undefined
    }>
      {tab === "users" && (users.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-users">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<User>(userSel.selected, userSel.toggle, () => userSel.toggleAll(users.map((u) => u.username)), users.map((u) => u.username), t, (row) => row.username), ...userColumns]} data={users} getRowId={(row) => row.username} onRowClick={(row) => { setSelectedUser(row); setDetailTab("detail"); }} selectedId={selectedUser?.username} hasMore={usersList.hasMore} onLoadMore={usersList.loadMore} loadingMore={usersList.isFetchingMore} /></div>
          {renderUserDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "roles" && (roles.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-roles">
          <div className="flex-fill-scroll"><DataTable columns={[checkboxColumn<Role>(roleSel.selected, roleSel.toggle, () => roleSel.toggleAll(roles.map((r) => r.rolename)), roles.map((r) => r.rolename), t, (row) => row.rolename), ...roleColumns]} data={roles} getRowId={(row) => row.rolename} onRowClick={(row) => { setSelectedRole(row); setDetailTab("detail"); }} selectedId={selectedRole?.rolename} hasMore={rolesList.hasMore} onLoadMore={rolesList.loadMore} loadingMore={rolesList.isFetchingMore} /></div>
          {renderRoleDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "policies" && (policies.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-policies">
          <div className="flex-fill-scroll"><DataTable columns={policyColumns} data={policies} getRowId={(row) => row.policyname} onRowClick={(row) => { setSelectedPolicy(row); setDetailTab("detail"); }} selectedId={selectedPolicy?.policyname} hasMore={policiesList.hasMore} onLoadMore={policiesList.loadMore} loadingMore={policiesList.isFetchingMore} /></div>
          {renderPolicyDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {/* User modals */}
      <ServiceCreateModal open={showCreateUser} onClose={() => setShowCreateUser(false)} title={t("services.iam.createUser")} error={createUserMutation.error} isPending={createUserMutation.isPending} onCreate={() => createUserMutation.mutate()} disabled={!formUserName}>
        <label>{t("services.iam.userNameLabel")}<input value={formUserName} onChange={(e) => setFormUserName(e.target.value)} placeholder={t("services.iam.userNamePlaceholder")} className="modal-input" /></label>
        <label>{t("services.iam.pathLabel")}<input value={formUserPath} onChange={(e) => setFormUserPath(e.target.value)} placeholder="/" className="modal-input" /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDeleteUser && !!selectedUser} title={t("services.iam.deleteUser")} name={selectedUser?.username} error={deleteUserMutation.error} isPending={deleteUserMutation.isPending} onConfirm={() => selectedUser && deleteUserMutation.mutate(selectedUser.username)} onClose={() => setShowDeleteUser(false)} />
      <ServiceDeleteDialog open={showBatchDeleteUser} title={t("common.deleteSelected")} name={`${userSel.selected.size} ${t("services.iam.tabs.users")}`} error={batchDeleteUserMutation.error} isPending={batchDeleteUserMutation.isPending} onConfirm={() => batchDeleteUserMutation.mutate(Array.from(userSel.selected))} onClose={() => setShowBatchDeleteUser(false)} />

      {/* Role modals */}
      <ServiceCreateModal open={showCreateRole} onClose={() => setShowCreateRole(false)} title={t("services.iam.createRole")} error={createRoleMutation.error} isPending={createRoleMutation.isPending} onCreate={() => createRoleMutation.mutate()} disabled={!formRoleName}>
        <label>{t("services.iam.roleNameLabel")}<input value={formRoleName} onChange={(e) => setFormRoleName(e.target.value)} placeholder={t("services.iam.roleNamePlaceholder")} className="modal-input" /></label>
        <label>{t("services.iam.pathLabel")}<input value={formRolePath} onChange={(e) => setFormRolePath(e.target.value)} placeholder="/" className="modal-input" /></label>
        <label>{t("services.iam.assumePolicyLabel")}<textarea value={formAssumePolicy} onChange={(e) => setFormAssumePolicy(e.target.value)} rows={8} className="modal-input" style={{ fontFamily: "monospace", fontSize: "0.85em" }} /></label>
      </ServiceCreateModal>
      <ServiceDeleteDialog open={showDeleteRole && !!selectedRole} title={t("services.iam.deleteRole")} name={selectedRole?.rolename} error={deleteRoleMutation.error} isPending={deleteRoleMutation.isPending} onConfirm={() => selectedRole && deleteRoleMutation.mutate(selectedRole.rolename)} onClose={() => setShowDeleteRole(false)} />
      <ServiceDeleteDialog open={showBatchDeleteRole} title={t("common.deleteSelected")} name={`${roleSel.selected.size} ${t("services.iam.tabs.roles")}`} error={batchDeleteRoleMutation.error} isPending={batchDeleteRoleMutation.isPending} onConfirm={() => batchDeleteRoleMutation.mutate(Array.from(roleSel.selected))} onClose={() => setShowBatchDeleteRole(false)} />
    </ServicePageLayout>
  );
}
