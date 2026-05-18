/**
 * IAM service page — 3-panel inspector layout with tabs (users, roles, policies).
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import { IAMService, type User, type Role, type Policy, CreateUserRequestSchema, CreateRoleRequestSchema } from "@/gen/iam_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import { ServicePageLayout, ServiceCreateModal, ServiceDeleteDialog, MonoCell, SmallMonoCell, DateCell, FallbackCell, useServiceClient } from "@/components/shared/service-page";
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
  const { t } = useTranslation();
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

  const { client, invalidate } = useServiceClient(IAMService);
  const { queryKey: usersKey } = useListKey("iam-users");
  const { queryKey: rolesKey } = useListKey("iam-roles");
  const { queryKey: policiesKey } = useListKey("iam-policies");

  const usersQ = useQuery({ queryKey: usersKey, queryFn: async () => { const r = await client.listUsers({}); return dropEmpty(r.users ?? [], "username"); }, refetchInterval: REFETCH_INTERVAL });
  const rolesQ = useQuery({ queryKey: rolesKey, queryFn: async () => { const r = await client.listRoles({}); return dropEmpty(r.roles ?? [], "rolename"); }, refetchInterval: REFETCH_INTERVAL });
  const policiesQ = useQuery({ queryKey: policiesKey, queryFn: async () => { const r = await client.listPolicies({}); return dropEmpty(r.policies ?? [], "policyname"); }, refetchInterval: REFETCH_INTERVAL });

  const query = tab === "users" ? usersQ : tab === "roles" ? rolesQ : policiesQ;
  const users = usersQ.data ?? [];
  const roles = rolesQ.data ?? [];
  const policies = policiesQ.data ?? [];

  const createUserMutation = useMutation({
    mutationFn: () => client.createUser(create(CreateUserRequestSchema, { username: formUserName, path: formUserPath || undefined })),
    onSuccess: () => { invalidate(usersKey); setShowCreateUser(false); setFormUserName(""); setFormUserPath(""); },
  });
  const deleteUserMutation = useMutation({
    mutationFn: (name: string) => client.deleteUser({ username: name }),
    onSuccess: () => { invalidate(usersKey); setShowDeleteUser(false); setSelectedUser(null); userSel.clear(); },
  });
  const batchDeleteUserMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteUser({ username: n }))),
    onSuccess: (_d, names) => { invalidate(usersKey); setShowBatchDeleteUser(false); userSel.clear(); setSelectedUser((p) => (p && names.includes(p.username) ? null : p)); },
  });

  const createRoleMutation = useMutation({
    mutationFn: () => client.createRole(create(CreateRoleRequestSchema, { rolename: formRoleName, path: formRolePath || undefined, assumerolepolicydocument: formAssumePolicy })),
    onSuccess: () => { invalidate(rolesKey); setShowCreateRole(false); setFormRoleName(""); setFormRolePath(""); },
  });
  const deleteRoleMutation = useMutation({
    mutationFn: (name: string) => client.deleteRole({ rolename: name }),
    onSuccess: () => { invalidate(rolesKey); setShowDeleteRole(false); setSelectedRole(null); roleSel.clear(); },
  });
  const batchDeleteRoleMutation = useMutation({
    mutationFn: async (names: string[]) => Promise.allSettled(names.map((n) => client.deleteRole({ rolename: n }))),
    onSuccess: (_d, names) => { invalidate(rolesKey); setShowBatchDeleteRole(false); roleSel.clear(); setSelectedRole((p) => (p && names.includes(p.rolename) ? null : p)); },
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
      <DetailPanel title={selectedUser.username} titleIcon="👤" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteUser(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Username</td><td className="cell-mono">{selectedUser.username}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>ARN</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedUser.arn}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Path</td><td>{selectedUser.path || "/"}</td></tr>
            {selectedUser.createdate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{new Date(selectedUser.createdate).toLocaleString()}</td></tr>}
            {selectedUser.passwordlastused && <tr><td style={{ fontWeight: 600 }}>Last Login</td><td>{new Date(selectedUser.passwordlastused).toLocaleString()}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedUser} />}
      </DetailPanel>
    );
  };

  const renderRoleDetail = () => {
    if (!selectedRole) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedRole.rolename} titleIcon="🔑" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)} actions={<button className="btn btn-danger btn-sm" onClick={() => setShowDeleteRole(true)}>{t("common.delete")}</button>}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Role</td><td className="cell-mono">{selectedRole.rolename}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>ARN</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedRole.arn}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Description</td><td>{selectedRole.description || "\u2014"}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Path</td><td>{selectedRole.path || "/"}</td></tr>
            {selectedRole.createdate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{new Date(selectedRole.createdate).toLocaleString()}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedRole} />}
      </DetailPanel>
    );
  };

  const renderPolicyDetail = () => {
    if (!selectedPolicy) return <DetailEmpty message={t("common.noItemSelected")} />;
    return (
      <DetailPanel title={selectedPolicy.policyname} titleIcon="📜" tabs={[{ key: "detail", label: "Detail" }, { key: "json", label: t("common.rawJson") ?? "JSON" }]} activeTab={detailTab} onTabChange={(k) => setDetailTab(k as DetailTab)}>
        {detailTab === "detail" ? (
          <table className="settings-table" style={{ width: "100%" }}><tbody>
            <tr><td style={{ width: 140, fontWeight: 600 }}>Policy</td><td className="cell-mono">{selectedPolicy.policyname}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>ARN</td><td className="cell-mono" style={{ fontSize: "0.85em" }}>{selectedPolicy.arn}</td></tr>
            <tr><td style={{ fontWeight: 600 }}>Attachments</td><td>{selectedPolicy.attachmentcount ?? 0}</td></tr>
            {selectedPolicy.createdate && <tr><td style={{ fontWeight: 600 }}>Created</td><td>{new Date(selectedPolicy.createdate).toLocaleString()}</td></tr>}
            {selectedPolicy.updatedate && <tr><td style={{ fontWeight: 600 }}>Updated</td><td>{new Date(selectedPolicy.updatedate).toLocaleString()}</td></tr>}
          </tbody></table>
        ) : <JsonViewer data={selectedPolicy} />}
      </DetailPanel>
    );
  };

  return (
    <ServicePageLayout icon="👤" title={t("services.iam.title")} isLoading={query.isLoading} error={query.error} tabs={tabs} activeTab={tab} onTabChange={handleTabChange} actions={
      tab === "users" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateUser(true)}>{t("services.iam.createUser")}</button>
        <button className="btn btn-danger" disabled={userSel.selected.size === 0} onClick={() => setShowBatchDeleteUser(true)}>{t("common.deleteSelected")}{userSel.selected.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({userSel.selected.size})</span>}</button>
      </>) : tab === "roles" ? (<>
        <button className="btn btn-primary" onClick={() => setShowCreateRole(true)}>{t("services.iam.createRole")}</button>
        <button className="btn btn-danger" disabled={roleSel.selected.size === 0} onClick={() => setShowBatchDeleteRole(true)}>{t("common.deleteSelected")}{roleSel.selected.size > 0 && <span style={{ marginLeft: 4, opacity: 0.8 }}>({roleSel.selected.size})</span>}</button>
      </>) : undefined
    }>
      {tab === "users" && (users.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-users">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<User>(userSel.selected, userSel.toggle, () => userSel.toggleAll(users.map((u) => u.username)), users.map((u) => u.username), t, (row) => row.username), ...userColumns]} data={users} getRowId={(row) => row.username} onRowClick={(row) => { setSelectedUser(row); setDetailTab("detail"); }} selectedId={selectedUser?.username} /></div>
          {renderUserDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "roles" && (roles.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-roles">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={[checkboxColumn<Role>(roleSel.selected, roleSel.toggle, () => roleSel.toggleAll(roles.map((r) => r.rolename)), roles.map((r) => r.rolename), t, (row) => row.rolename), ...roleColumns]} data={roles} getRowId={(row) => row.rolename} onRowClick={(row) => { setSelectedRole(row); setDetailTab("detail"); }} selectedId={selectedRole?.rolename} /></div>
          {renderRoleDetail()}
        </Splitter>
      ) : <div className="empty-state">{t("common.noData")}</div>)}

      {tab === "policies" && (policies.length > 0 ? (
        <Splitter direction="horizontal" initialSize={240} minSize={80} maxSize={600} storageKey="vs-split-iam-policies">
          <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}><DataTable columns={policyColumns} data={policies} getRowId={(row) => row.policyname} onRowClick={(row) => { setSelectedPolicy(row); setDetailTab("detail"); }} selectedId={selectedPolicy?.policyname} /></div>
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
