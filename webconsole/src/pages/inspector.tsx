/**
 * Pebble Inspector page — read-only key-value browser for config store.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery } from "@tanstack/react-query";
import { AdminConfigService } from "@/gen/admin_config_pb";
import { ServicePageLayout, MonoCell, useServiceClient } from "@/components/shared/service-page";
import { DataTable } from "@/components/shared/data-table";
import type { ColumnDef } from "@tanstack/react-table";
import type { TFunction } from "i18next";
import { JsonViewer } from "@/components/shared/json-viewer";
import { REFETCH_INTERVAL } from "@/lib/use-service-list";

interface ConfigRow {
  key: string;
  value: string;
  type: string;
  description: string;
}

const getColumns = (t: TFunction): ColumnDef<ConfigRow, any>[] => [
  { accessorKey: "key", header: t("inspector.key"), cell: MonoCell },
  { accessorKey: "type", header: t("inspector.typeHeader"), size: 100 },
  { accessorKey: "value", header: t("inspector.value"), size: 300 },
];

export function InspectorPage() {
  const { t } = useTranslation();
  const { client } = useServiceClient(AdminConfigService);
  const [selectedItem, setSelectedItem] = useState<ConfigRow | null>(null);
  const [prefixFilter, setPrefixFilter] = useState("");

  const { data, isLoading, error } = useQuery({
    queryKey: ["inspector-config"],
    queryFn: () => client.listConfig({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const allRows: ConfigRow[] = (data?.entries ?? []).map((e) => ({
    key: e.key,
    value: e.value,
    type: e.type,
    description: e.description,
  }));

  const items = prefixFilter
    ? allRows.filter((r) => r.key.startsWith(prefixFilter))
    : allRows;

  const columns = getColumns(t);

  return (
    <ServicePageLayout
      icon="🔍"
      title={t("inspector.title")}
      isLoading={isLoading}
      error={error}
      count={items.length}
      countLabel={t("inspector.countLabel")}
    >
      <div className="inspector-toolbar">
        <input
          type="text"
          className="data-table-filter"
          value={prefixFilter}
          onChange={(e) => setPrefixFilter(e.target.value)}
          placeholder={t("inspector.filter")}
          aria-label={t("inspector.filter")}
        />
        {prefixFilter && (
          <span className="data-table-filter-count">
            {items.length}/{allRows.length}
          </span>
        )}
      </div>
      <div className="flex-fill-scroll">
        <DataTable
          columns={columns}
          data={items}
          getRowId={(row) => row.key}
          onRowClick={(row) => setSelectedItem(row)}
          selectedId={selectedItem?.key}
          pageSize={100}
          exportName="inspector"
        />
      </div>
      {selectedItem && (
        <div className="detail-panel" style={{ borderTop: "1px solid var(--border-color)", padding: 12 }}>
          <h3 style={{ margin: "0 0 8px", fontSize: 13 }}>{selectedItem.key}</h3>
          {selectedItem.description && (
            <div style={{ fontSize: 11, color: "var(--text-muted)", marginBottom: 8 }}>
              {selectedItem.description}
            </div>
          )}
          <JsonViewer data={selectedItem} />
        </div>
      )}
    </ServicePageLayout>
  );
}
