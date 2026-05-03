/**
 * IAM service page. Lists users, roles, and policies via tabbed views.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { IAMService, type User, type Role, type Policy } from "@/gen/iam_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "users" | "roles" | "policies";

const userColumns: ColumnDef<User, any>[] = [
  {
    accessorKey: "username",
    header: "User Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "arn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "createdate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

const roleColumns: ColumnDef<Role, any>[] = [
  {
    accessorKey: "rolename",
    header: "Role Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "arn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "createdate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

const policyColumns: ColumnDef<Policy, any>[] = [
  {
    accessorKey: "policyname",
    header: "Policy Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "arn",
    header: "ARN",
    cell: ({ getValue }) => (
      <span className="cell-mono" style={{ fontSize: "0.85em" }}>{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "attachmentcount",
    header: "Attachments",
    size: 100,
  },
];

export function IAMPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("users");
  const [selectedUser, setSelectedUser] = useState<User | null>(null);
  const [selectedRole, setSelectedRole] = useState<Role | null>(null);
  const [selectedPolicy, setSelectedPolicy] = useState<Policy | null>(null);

  const client = createClient(IAMService, transport);
  const { queryKey: usersKey } = useListKey("iam-users");
  const { queryKey: rolesKey } = useListKey("iam-roles");
  const { queryKey: policiesKey } = useListKey("iam-policies");

  const usersQ = useQuery({
    queryKey: usersKey,
    queryFn: () => client.listUsers({}),
    refetchInterval: 30_000,
  });

  const rolesQ = useQuery({
    queryKey: rolesKey,
    queryFn: () => client.listRoles({}),
    refetchInterval: 30_000,
  });

  const policiesQ = useQuery({
    queryKey: policiesKey,
    queryFn: () => client.listPolicies({}),
    refetchInterval: 30_000,
  });

  const query = tab === "users" ? usersQ : tab === "roles" ? rolesQ : policiesQ;

  if (query.isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔒</span>
          <h1>IAM</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (query.error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">🔒</span>
          <h1>IAM</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(query.error) })}</div>
      </div>
    );
  }

  const users = dropEmpty(usersQ.data?.users ?? [], "username");
  const roles = dropEmpty(rolesQ.data?.roles ?? [], "rolename");
  const policies = dropEmpty(policiesQ.data?.policies ?? [], "policyname");

  const tabs: { key: TabKey; label: string; count: number }[] = [
    { key: "users", label: "Users", count: users.length },
    { key: "roles", label: "Roles", count: roles.length },
    { key: "policies", label: "Policies", count: policies.length },
  ];

  const selected = tab === "users" ? selectedUser : tab === "roles" ? selectedRole : selectedPolicy;
  const setSelected = tab === "users" ? setSelectedUser : tab === "roles" ? setSelectedRole : setSelectedPolicy;

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">🔒</span>
        <h1>IAM</h1>
        <div className="tab-bar">
          {tabs.map((t) => (
            <button
              key={t.key}
              className={`tab-btn ${tab === t.key ? "active" : ""}`}
              onClick={() => setTab(t.key)}
            >
              {t.label} ({t.count})
            </button>
          ))}
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          {tab === "users" && (
            <DataTable
              columns={userColumns}
              data={users}
              getRowId={(row) => row.username}
              onRowClick={setSelectedUser as any}
              selectedId={selectedUser?.username}
            />
          )}
          {tab === "roles" && (
            <DataTable
              columns={roleColumns}
              data={roles}
              getRowId={(row) => row.rolename}
              onRowClick={setSelectedRole as any}
              selectedId={selectedRole?.rolename}
            />
          )}
          {tab === "policies" && (
            <DataTable
              columns={policyColumns}
              data={policies}
              getRowId={(row) => row.policyname}
              onRowClick={setSelectedPolicy as any}
              selectedId={selectedPolicy?.policyname}
            />
          )}
        </div>
        {selected && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{(selected as any).username || (selected as any).rolename || (selected as any).policyname}</h2>
              <button className="detail-close" onClick={() => setSelected(null as any)}>
                ✕
              </button>
            </div>
            <JsonViewer data={selected} />
          </div>
        )}
      </div>
    </div>
  );
}
