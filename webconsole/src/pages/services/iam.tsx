/**
 * IAM service page. Tabbed view for users, roles, and policies with create/delete
 * operations for users and roles.
 */
import { useState } from "react";
import type { TFunction } from "i18next";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  IAMService,
  type User,
  type Role,
  type Policy,
  CreateUserRequestSchema,
  CreateRoleRequestSchema,
} from "@/gen/iam_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  SplitPane,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  FallbackCell,
  useServiceClient,
} from "@/components/shared/service-page";

/** Tab key type for the IAM page. */
type TabKey = "users" | "roles" | "policies";

/** Column definitions for the IAM user table. */
const getUserColumns = (t: TFunction): ColumnDef<User, any>[] => [
  { accessorKey: "username", header: t("services.iam.userNameHeader"), cell: MonoCell },
  { accessorKey: "path", header: t("services.iam.pathHeader"), cell: FallbackCell, size: 90 },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "createdate", header: t("services.iam.createdHeader"), cell: DateCell },
  { accessorKey: "passwordlastused", header: t("services.iam.passwordLastUsedHeader"), cell: DateCell },
];

/** Column definitions for the IAM role table. */
const getRoleColumns = (t: TFunction): ColumnDef<Role, any>[] => [
  { accessorKey: "rolename", header: t("services.iam.roleNameHeader"), cell: MonoCell },
  { accessorKey: "description", header: t("services.iam.roleDescriptionHeader"), cell: FallbackCell },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "createdate", header: t("services.iam.createdHeader"), cell: DateCell },
];

/** Column definitions for the IAM policy table. */
const getPolicyColumns = (t: TFunction): ColumnDef<Policy, any>[] => [
  { accessorKey: "policyname", header: t("services.iam.policyNameHeader"), cell: MonoCell },
  { accessorKey: "arn", header: t("services.iam.arnHeader"), cell: SmallMonoCell },
  { accessorKey: "attachmentcount", header: t("services.iam.attachmentsHeader"), size: 100 },
  { accessorKey: "createdate", header: t("services.iam.policyCreatedHeader"), cell: DateCell },
  { accessorKey: "updatedate", header: t("services.iam.policyUpdatedHeader"), cell: DateCell },
];

/** IAM service page with tabbed users/roles/policies view and CRUD operations. */
export function IAMPage() {
  const { t } = useTranslation();
  const userColumns = getUserColumns(t);
  const roleColumns = getRoleColumns(t);
  const policyColumns = getPolicyColumns(t);
  const [tab, setTab] = useState<TabKey>("users");
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedPolicy, setSelectedPolicy] = useState<Policy | null>(null);
  const DEFAULT_TRUST_POLICY = JSON.stringify({
    Version: "2012-10-17",
    Statement: [
      { Effect: "Allow", Principal: { Service: "lambda.amazonaws.com" }, Action: "sts:AssumeRole" },
    ],
  }, null, 2);

  const [showCreate, setShowCreate] = useState(false);
  const [showDelete, setShowDelete] = useState(false);
  const [formName, setFormName] = useState("");
  const [formPath, setFormPath] = useState("/");
  const [formTrustPolicy, setFormTrustPolicy] = useState(DEFAULT_TRUST_POLICY);
  const [formMaxSessionDuration, setFormMaxSessionDuration] = useState(3600);
  const [formDescription, setFormDescription] = useState("");

  const { client, invalidate } = useServiceClient(IAMService);
  const { queryKey: usersKey } = useListKey("iam-users");
  const { queryKey: rolesKey } = useListKey("iam-roles");
  const { queryKey: policiesKey } = useListKey("iam-policies");

  const usersQ = useQuery({
    queryKey: usersKey,
    queryFn: () => client.listUsers({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const rolesQ = useQuery({
    queryKey: rolesKey,
    queryFn: () => client.listRoles({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const policiesQ = useQuery({
    queryKey: policiesKey,
    queryFn: () => client.listPolicies({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const query = tab === "users" ? usersQ : tab === "roles" ? rolesQ : policiesQ;

  const createUserMutation = useMutation({
    mutationFn: () =>
      client.createUser(create(CreateUserRequestSchema, {
        username: formName,
        path: formPath || "/",
      })),
    onSuccess: () => {
      invalidate(usersKey);
      setShowCreate(false);
      setFormName("");
      setFormPath("/");
    },
  });

  const createRoleMutation = useMutation({
    mutationFn: () =>
      client.createRole(
        create(CreateRoleRequestSchema, {
          rolename: formName,
          assumerolepolicydocument: formTrustPolicy,
          maxsessionduration: formMaxSessionDuration,
          ...(formDescription ? { description: formDescription } : {}),
        }),
      ),
    onSuccess: () => {
      invalidate(rolesKey);
      setShowCreate(false);
      setFormName("");
      setFormTrustPolicy(DEFAULT_TRUST_POLICY);
      setFormMaxSessionDuration(3600);
      setFormDescription("");
    },
  });

  const deleteUserMutation = useMutation({
    mutationFn: (username: string) => client.deleteUser({ username }),
    onSuccess: () => {
      invalidate(usersKey);
      setShowDelete(false);
      setSelectedUser(null);
    },
  });

  const deleteRoleMutation = useMutation({
    mutationFn: (rolename: string) => client.deleteRole({ rolename }),
    onSuccess: () => {
      invalidate(rolesKey);
      setShowDelete(false);
      setSelectedRole(null);
    },
  });

  const createMutation = tab === "users" ? createUserMutation : createRoleMutation;
  const deleteMutation = tab === "users" ? deleteUserMutation : deleteRoleMutation;

  const users = dropEmpty(usersQ.data?.users ?? [], "username");
  const roles = dropEmpty(rolesQ.data?.roles ?? [], "rolename");
  const policies = dropEmpty(policiesQ.data?.policies ?? [], "policyname");

  const tabs = [
    { key: "users" as TabKey, label: t("services.iam.tabs.users"), count: users.length },
    { key: "roles" as TabKey, label: t("services.iam.tabs.roles"), count: roles.length },
    { key: "policies" as TabKey, label: t("services.iam.tabs.policies"), count: policies.length },
  ];

  const selected = tab === "users" ? selectedUser : tab === "roles" ? selectedRole : selectedPolicy;
  const hasCrud = tab === "users" || tab === "roles";

  return (
    <ServicePageLayout
      icon="🔒"
      title={t("services.iam.title")}
      isLoading={query.isLoading}
      error={query.error}
      tabs={tabs}
      activeTab={tab}
      onTabChange={(k) => { setTab(k as TabKey); setFormName(""); setFormPath("/"); setFormTrustPolicy(DEFAULT_TRUST_POLICY); setFormMaxSessionDuration(3600); setFormDescription(""); setShowCreate(false); setShowDelete(false); setSelectedUser(null); setSelectedRole(null); setSelectedPolicy(null); }}
      actions={
        hasCrud ? (
          <>
            <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
              {tab === "users" ? t("services.iam.createUser") : t("services.iam.createRole")}
            </button>
            {selected && (
              <button className="btn btn-danger" onClick={() => setShowDelete(true)}>
                {t("common.delete")}
              </button>
            )}
          </>
        ) : undefined
      }
    >
      {tab === "users" && (
        <SplitPane
          columns={userColumns}
          data={users}
          getRowId={(row) => row.username}
          onRowClick={(row) => setSelectedUser(row)}
          selectedId={selectedUser?.username}
          selected={selectedUser}
          detailTitle={selectedUser?.username}
          onDetailClose={() => setSelectedUser(null)}
        />
      )}
      {tab === "roles" && (
        <SplitPane
          columns={roleColumns}
          data={roles}
          getRowId={(row) => row.rolename}
          onRowClick={(row) => setSelectedRole(row)}
          selectedId={selectedRole?.rolename}
          selected={selectedRole}
          detailTitle={selectedRole?.rolename}
          onDetailClose={() => setSelectedRole(null)}
        />
      )}
      {tab === "policies" && (
        <SplitPane
          columns={policyColumns}
          data={policies}
          getRowId={(row) => row.policyname}
          onRowClick={(row) => setSelectedPolicy(row)}
          selectedId={selectedPolicy?.policyname}
          selected={selectedPolicy}
          detailTitle={selectedPolicy?.policyname}
          onDetailClose={() => setSelectedPolicy(null)}
        />
      )}

      {hasCrud && (
        <ServiceCreateModal
          open={showCreate}
          onClose={() => setShowCreate(false)}
          title={tab === "users" ? t("services.iam.createUser") : t("services.iam.createRole")}
          error={createMutation.error}
          isPending={createMutation.isPending}
          onCreate={() => createMutation.mutate()}
          disabled={!formName}
        >
          <label>
            {tab === "users" ? t("services.iam.userNameField") : t("services.iam.roleNameField")}
            <input
              value={formName}
              onChange={(e) => setFormName(e.target.value)}
              placeholder={tab === "users" ? t("services.iam.userNamePlaceholder") : t("services.iam.roleNamePlaceholder")}
              className="modal-input"
            />
          </label>
          {tab === "users" && (
            <label>
              {t("services.iam.pathLabel")}
              <input
                value={formPath}
                onChange={(e) => setFormPath(e.target.value)}
                placeholder="/"
                className="modal-input"
              />
            </label>
          )}
          {tab === "roles" && (
            <>
              <label>
                {t("services.iam.trustPolicyLabel")}
                <textarea
                  value={formTrustPolicy}
                  onChange={(e) => setFormTrustPolicy(e.target.value)}
                  rows={8}
                  className="modal-input"
                  style={{ fontFamily: "monospace", fontSize: 12 }}
                />
              </label>
              <label>
                {t("services.iam.maxSessionDurationLabel")}
                <input
                  type="number"
                  value={formMaxSessionDuration}
                  onChange={(e) => setFormMaxSessionDuration(Number(e.target.value))}
                  min={900}
                  max={43200}
                  className="modal-input"
                />
              </label>
              <label>
                {t("services.iam.descriptionLabel")}
                <input
                  value={formDescription}
                  onChange={(e) => setFormDescription(e.target.value)}
                  placeholder={t("services.iam.descriptionPlaceholder")}
                  className="modal-input"
                />
              </label>
            </>
          )}
        </ServiceCreateModal>
      )}

      {hasCrud && (
        <ServiceDeleteDialog
          open={showDelete && !!selected}
          title={tab === "users" ? t("services.iam.deleteUser") : t("services.iam.deleteRole")}
          name={tab === "users" ? selectedUser?.username : selectedRole?.rolename}
          error={deleteMutation.error}
          isPending={deleteMutation.isPending}
          onConfirm={() => {
            if (tab === "users" && selectedUser) deleteUserMutation.mutate(selectedUser.username);
            if (tab === "roles" && selectedRole) deleteRoleMutation.mutate(selectedRole.rolename);
          }}
          onClose={() => setShowDelete(false)}
        />
      )}
    </ServicePageLayout>
  );
}
