/**
 * Cognito service page. Lists user pools and identity pools.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { createClient } from "@connectrpc/connect";
import { CognitoIdentityProviderService, type UserPoolDescriptionType } from "@/gen/cognitoidentityprovider_pb";
import { CognitoIdentityService, type IdentityPoolShortDescription } from "@/gen/cognitoidentity_pb";
import { transport } from "@/lib/transport";
import { useListKey, dropEmpty } from "@/lib/use-service-list";
import { DataTable } from "@/components/shared/data-table";
import { JsonViewer } from "@/components/shared/json-viewer";

type TabKey = "userpools" | "identitypools";

const userPoolColumns: ColumnDef<UserPoolDescriptionType, any>[] = [
  {
    accessorKey: "name",
    header: "Pool Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "id",
    header: "Pool ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "status",
    header: "Status",
    size: 90,
  },
  {
    accessorKey: "creationdate",
    header: "Created",
    cell: ({ getValue }) => {
      const v = getValue() as string;
      if (!v) return "\u2014";
      try { return new Date(v).toLocaleString(); } catch { return v; }
    },
  },
];

const identityPoolColumns: ColumnDef<IdentityPoolShortDescription, any>[] = [
  {
    accessorKey: "identitypoolname",
    header: "Pool Name",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
  {
    accessorKey: "identitypoolid",
    header: "Pool ID",
    cell: ({ getValue }) => (
      <span className="cell-mono">{getValue() as string}</span>
    ),
  },
];

export function CognitoPage() {
  const { t } = useTranslation();
  const [tab, setTab] = useState<TabKey>("userpools");
  const [selectedUP, setSelectedUP] = useState<UserPoolDescriptionType | null>(null);
  const [selectedIP, setSelectedIP] = useState<IdentityPoolShortDescription | null>(null);

  const idpClient = createClient(CognitoIdentityProviderService, transport);
  const identityClient = createClient(CognitoIdentityService, transport);
  const { queryKey: upKey } = useListKey("cognito-userpools");
  const { queryKey: ipKey } = useListKey("cognito-identitypools");

  const upQ = useQuery({
    queryKey: upKey,
    queryFn: () => idpClient.listUserPools({}),
    refetchInterval: 30_000,
  });

  const ipQ = useQuery({
    queryKey: ipKey,
    queryFn: () => identityClient.listIdentityPools({}),
    refetchInterval: 30_000,
  });

  const query = tab === "userpools" ? upQ : ipQ;

  if (query.isLoading) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">👤</span>
          <h1>Cognito</h1>
        </div>
        <div className="loading-state">{t("common.loading")}</div>
      </div>
    );
  }

  if (query.error) {
    return (
      <div className="content-area">
        <div className="page-header">
          <span className="page-icon">👤</span>
          <h1>Cognito</h1>
        </div>
        <div className="error-state">{t("common.failedToLoad", { error: String(query.error) })}</div>
      </div>
    );
  }

  const userPools = dropEmpty(upQ.data?.userpools ?? [], "id");
  const identityPools = dropEmpty(ipQ.data?.identitypools ?? [], "identitypoolid");
  const selected = tab === "userpools" ? selectedUP : selectedIP;

  return (
    <div className="content-area">
      <div className="page-header">
        <span className="page-icon">👤</span>
        <h1>Cognito</h1>
        <div className="tab-bar">
          <button
            className={`tab-btn ${tab === "userpools" ? "active" : ""}`}
            onClick={() => setTab("userpools")}
          >
            User Pools ({userPools.length})
          </button>
          <button
            className={`tab-btn ${tab === "identitypools" ? "active" : ""}`}
            onClick={() => setTab("identitypools")}
          >
            Identity Pools ({identityPools.length})
          </button>
        </div>
      </div>

      <div className="split-pane">
        <div className="split-table">
          {tab === "userpools" && (
            <DataTable
              columns={userPoolColumns}
              data={userPools}
              getRowId={(row) => row.id}
              onRowClick={setSelectedUP}
              selectedId={selectedUP?.id}
            />
          )}
          {tab === "identitypools" && (
            <DataTable
              columns={identityPoolColumns}
              data={identityPools}
              getRowId={(row) => row.identitypoolid}
              onRowClick={setSelectedIP}
              selectedId={selectedIP?.identitypoolid}
            />
          )}
        </div>
        {selected && (
          <div className="split-detail">
            <div className="detail-header">
              <h2>{(selected as any).name || (selected as any).identitypoolname}</h2>
              <button className="detail-close" onClick={() => { setSelectedUP(null); setSelectedIP(null); }}>
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
