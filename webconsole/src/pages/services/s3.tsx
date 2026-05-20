/**
 * S3 service page — 3-panel inspector layout matching the PoC mockup.
 *
 * Panel 1 (toolbar): Breadcrumb navigation  S3 / Buckets  or  S3 / bucket / prefix/
 * Panel 2 (table):   Bucket list OR object list with checkbox multi-select
 * Panel 3 (detail):  Object metadata JSON/RAW/HEADERS tabs (drag-split bottom panel)
 */
import { useState, useCallback, type ChangeEvent } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import type { ColumnDef } from "@tanstack/react-table";
import { create } from "@bufbuild/protobuf";
import {
  S3Service,
  CreateBucketRequestSchema,
  DeleteBucketRequestSchema,
  ListObjectsV2RequestSchema,
  GetObjectRequestSchema,
  PutObjectRequestSchema,
  DeleteObjectRequestSchema,
  DeleteObjectsRequestSchema,
  CopyObjectRequestSchema,
  ObjectOwnership,
  type Bucket,
  type Object$,
  type CommonPrefix,
} from "@/gen/s3_pb";
import { useListKey, dropEmpty, REFETCH_INTERVAL } from "@/lib/use-service-list";
import {
  ServicePageLayout,
  ServiceCreateModal,
  ServiceDeleteDialog,
  MonoCell,
  SmallMonoCell,
  DateCell,
  useServiceClient,
} from "@/components/shared/service-page";
import {
  checkboxColumn,
  useSelection,
} from "@/components/shared/inspector";
import { DataTable } from "@/components/shared/data-table";
import { Modal } from "@/components/shared/modal";
import { Splitter } from "@/components/shared/splitter";
import { JsonViewer } from "@/components/shared/json-viewer";
import { formatBytes } from "@/lib/format";

// ─── Constants ──────────────────────────────────────────────────

const MAX_UPLOAD_BYTES = 100 * 1024 * 1024;
const MAX_TEXT_PREVIEW = 256 * 1024;
const MAX_IMAGE_PREVIEW = 5 * 1024 * 1024;
const LIST_PAGE_SIZE = 100;

const TEXT_EXTENSIONS = new Set([
  ".txt", ".json", ".csv", ".md", ".xml", ".log", ".yaml", ".yml",
  ".toml", ".ini", ".cfg", ".conf", ".html", ".css", ".js", ".ts",
  ".go", ".py", ".rs", ".java", ".sh", ".bash", ".zsh", ".env",
]);

const IMAGE_EXTENSIONS = new Set([
  ".png", ".jpg", ".jpeg", ".gif", ".svg", ".webp", ".ico", ".bmp",
]);

const MIME_MAP: Record<string, string> = {
  ".json": "application/json", ".csv": "text/csv", ".xml": "application/xml",
  ".html": "text/html", ".css": "text/css", ".js": "application/javascript",
  ".ts": "application/typescript", ".md": "text/markdown", ".yaml": "application/x-yaml",
  ".yml": "application/x-yaml", ".txt": "text/plain", ".png": "image/png",
  ".jpg": "image/jpeg", ".jpeg": "image/jpeg", ".gif": "image/gif",
  ".svg": "image/svg+xml", ".webp": "image/webp", ".pdf": "application/pdf",
};

// ─── Helpers ────────────────────────────────────────────────────

function getExtension(key: string): string {
  const lastDot = key.lastIndexOf(".");
  if (lastDot === -1) return "";
  return key.slice(lastDot).toLowerCase();
}

function inferContentType(filename: string): string {
  return MIME_MAP[getExtension(filename)] ?? "application/octet-stream";
}

function isTextFile(key: string): boolean {
  return TEXT_EXTENSIONS.has(getExtension(key));
}

function isImageFile(key: string): boolean {
  return IMAGE_EXTENSIONS.has(getExtension(key));
}

/** Format bytes for display in the object table, falling back to the shared formatBytes. */
function formatSize(bytes: bigint | number): string {
  return formatBytes(bytes);
}

function fileNameFromKey(key: string): string {
  const idx = key.lastIndexOf("/");
  return idx === -1 ? key : key.slice(idx + 1);
}

function stripTrailingSlash(s: string): string {
  return s.endsWith("/") ? s.slice(0, -1) : s;
}

// ─── Unified row type for object table (folders + files) ────────

interface ObjectRow {
  id: string;
  key: string;
  displayKey: string;
  size: string;
  lastModified: string;
  storageClass: string;
  isFolder: boolean;
  rawSize: bigint;
}

function toObjectRows(
  prefixes: CommonPrefix[],
  objects: Object$[],
  currentPrefix: string,
): ObjectRow[] {
  const folders: ObjectRow[] = prefixes.map((p) => ({
    id: `folder:${p.prefix}`,
    key: p.prefix,
    displayKey: stripTrailingSlash(p.prefix.slice(currentPrefix.length)),
    size: "—",
    lastModified: "—",
    storageClass: "—",
    isFolder: true,
    rawSize: 0n,
  }));
  const files: ObjectRow[] = objects
    .filter((o) => o.key !== currentPrefix)
    .map((o) => ({
      id: `file:${o.key}`,
      key: o.key,
      displayKey: o.key.slice(currentPrefix.length),
      size: formatSize(o.size),
      lastModified: o.lastmodified,
      storageClass: String(o.storageclass ?? "STANDARD"),
      isFolder: false,
      rawSize: o.size,
    }));
  return [...folders, ...files];
}



// ─── View state ─────────────────────────────────────────────────

type ViewState =
  | { type: "buckets" }
  | { type: "objects"; bucket: Bucket };

// ─── Object detail tabs ─────────────────────────────────────────

type DetailTab = "json" | "raw" | "headers";

// ─── S3 Page ────────────────────────────────────────────────────

export function S3Page() {
  const { t } = useTranslation();
  const { client, invalidate } = useServiceClient(S3Service);
  const queryClient = useQueryClient();
  const { queryKey } = useListKey("s3");

  // ── View state ───────────────────────────────────────────────
  const [view, setView] = useState<ViewState>({ type: "buckets" });

  // ── Bucket checkbox multi-select ─────────────────────────────
  const {
    selected: selectedBucketNames,
    toggle: toggleBucket,
    toggleAll: toggleAllBuckets,
    clear: clearBucketSelection,
  } = useSelection<string>();

  // ── Object browsing state ────────────────────────────────────
  const [prefix, setPrefix] = useState("");
  const [continuationToken, setContinuationToken] = useState("");
  const [accumulatedObjects, setAccumulatedObjects] = useState<Object$[]>([]);
  const [accumulatedPrefixes, setAccumulatedPrefixes] = useState<CommonPrefix[]>([]);

  // ── Object selection (single click for detail) ───────────────
  const [selectedObject, setSelectedObject] = useState<ObjectRow | null>(null);

  // ── Object checkbox multi-select ─────────────────────────────
  const {
    selected: selectedObjectIds,
    toggle: toggleObject,
    toggleAll: toggleAllObjects,
    clear: clearObjectSelection,
  } = useSelection<string>();

  // ── Detail panel tab ─────────────────────────────────────────
  const [detailTab, setDetailTab] = useState<DetailTab>("json");

  // ── Modals ───────────────────────────────────────────────────
  const [showCreate, setShowCreate] = useState(false);
  const [formBucket, setFormBucket] = useState("");
  const [formOwnership, setFormOwnership] = useState<string>("BUCKETOWNERENFORCED");
  const [showDeleteBuckets, setShowDeleteBuckets] = useState(false);

  const [showUpload, setShowUpload] = useState(false);
  const [uploadFile, setUploadFile] = useState<File | null>(null);
  const [uploadKey, setUploadKey] = useState("");
  const [uploadError, setUploadError] = useState("");
  const [showDeleteObj, setShowDeleteObj] = useState(false);
  const [showCopy, setShowCopy] = useState(false);
  const [copyDestBucket, setCopyDestBucket] = useState("");
  const [copyDestKey, setCopyDestKey] = useState("");

  // ── Bucket list query ────────────────────────────────────────
  const { data: bucketData, isLoading, error } = useQuery({
    queryKey,
    queryFn: () => client.listBuckets({}),
    refetchInterval: REFETCH_INTERVAL,
  });

  const buckets: Bucket[] = dropEmpty(bucketData?.buckets ?? [], "name");

  // ── Object list query (only in objects view) ─────────────────
  const activeBucket = view.type === "objects" ? view.bucket.name : "";
  const objectQueryKey = ["s3", "objects", activeBucket, prefix];

  const { data: listData, isFetching: listFetching } = useQuery({
    queryKey: [...objectQueryKey, continuationToken],
    queryFn: () =>
      client.listObjectsV2(
        create(ListObjectsV2RequestSchema, {
          bucket: activeBucket,
          prefix,
          delimiter: "/",
          maxkeys: LIST_PAGE_SIZE,
          continuationtoken: continuationToken || undefined,
        }),
      ),
    enabled: view.type === "objects",
    refetchInterval: REFETCH_INTERVAL,
  });

  const allPrefixes = listData?.commonprefixes ?? [];
  const allObjects = listData?.contents ?? [];
  const nextToken = listData?.nextcontinuationtoken ?? "";

  const objectRows = view.type === "objects"
    ? toObjectRows(
        continuationToken ? accumulatedPrefixes : allPrefixes,
        continuationToken ? [...accumulatedObjects, ...allObjects] : allObjects,
        prefix,
      )
    : [];

  // ── Navigation helpers ───────────────────────────────────────

  const navigateToBucket = useCallback((bucket: Bucket) => {
    setView({ type: "objects", bucket });
    setPrefix("");
    setContinuationToken("");
    setAccumulatedObjects([]);
    setAccumulatedPrefixes([]);
    setSelectedObject(null);
    clearObjectSelection();
    setDetailTab("json");
  }, []);

  const navigateToBuckets = useCallback(() => {
    setView({ type: "buckets" });
    clearBucketSelection();
    setSelectedObject(null);
    clearObjectSelection();
  }, []);

  const handleBreadcrumb = useCallback((index: number) => {
    if (index === -1) {
      setPrefix("");
    } else {
      const parts = prefix.split("/").filter(Boolean);
      setPrefix(parts.slice(0, index + 1).join("/") + "/");
    }
    setContinuationToken("");
    setAccumulatedObjects([]);
    setAccumulatedPrefixes([]);
    setSelectedObject(null);
    clearObjectSelection();
  }, [prefix]);

  const handleFolderClick = useCallback((folderPrefix: string) => {
    setPrefix(folderPrefix);
    setContinuationToken("");
    setAccumulatedObjects([]);
    setAccumulatedPrefixes([]);
    setSelectedObject(null);
    clearObjectSelection();
  }, []);

  const handleLoadMore = () => {
    setAccumulatedPrefixes(allPrefixes);
    setAccumulatedObjects(allObjects);
    setContinuationToken(nextToken);
  };

  const invalidateObjects = useCallback(() => {
    queryClient.invalidateQueries({ queryKey: objectQueryKey });
    setContinuationToken("");
    setAccumulatedObjects([]);
    setAccumulatedPrefixes([]);
  }, [queryClient, objectQueryKey]);

  // ── Derived toggle helpers for checkbox columns ──────────────

  const allBucketIds = buckets.map((b) => b.name);
  const allObjectIds = objectRows.map((r) => r.id);

  // ── Preview query ────────────────────────────────────────────

  const selectedFile = selectedObject && !selectedObject.isFolder ? selectedObject : null;
  const canPreview = selectedFile !== null
    && (isTextFile(selectedFile.key) || isImageFile(selectedFile.key))
    && selectedFile.rawSize <= (isTextFile(selectedFile.key) ? MAX_TEXT_PREVIEW : MAX_IMAGE_PREVIEW);

  const { data: previewData, isFetching: previewFetching, error: previewError } = useQuery({
    queryKey: ["s3", "preview", activeBucket, selectedFile?.key],
    queryFn: () =>
      client.getObject(create(GetObjectRequestSchema, { bucket: activeBucket, key: selectedFile!.key })),
    enabled: canPreview,
  });

  // ── Bucket mutations ─────────────────────────────────────────

  const createMutation = useMutation({
    mutationFn: () =>
      client.createBucket(create(CreateBucketRequestSchema, {
        bucket: formBucket,
        objectownership: ObjectOwnership[formOwnership as keyof typeof ObjectOwnership] ?? ObjectOwnership.BUCKETOWNERENFORCED,
      })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowCreate(false);
      setFormBucket("");
      setFormOwnership("BUCKETOWNERENFORCED");
    },
  });

  const deleteBucketMutation = useMutation({
    mutationFn: (bucket: string) =>
      client.deleteBucket(create(DeleteBucketRequestSchema, { bucket })),
    onSuccess: () => {
      invalidate(queryKey);
      setShowDeleteBuckets(false);
      clearBucketSelection();
    },
  });

  // ── Object mutations ─────────────────────────────────────────

  const uploadMutation = useMutation({
    mutationFn: async () => {
      if (!uploadFile) throw new Error("No file selected");
      if (uploadFile.size > MAX_UPLOAD_BYTES) {
        throw new Error(`File size exceeds ${MAX_UPLOAD_BYTES / 1024 / 1024} MB limit`);
      }
      const buffer = await uploadFile.arrayBuffer();
      return client.putObject(
        create(PutObjectRequestSchema, {
          bucket: activeBucket,
          key: prefix + uploadKey,
          body: new Uint8Array(buffer),
          contenttype: inferContentType(uploadKey),
        }),
      );
    },
    onSuccess: () => {
      invalidateObjects();
      setShowUpload(false);
      setUploadFile(null);
      setUploadKey("");
      setUploadError("");
    },
    onError: (err: Error) => setUploadError(err.message),
  });

  const downloadMutation = useMutation({
    mutationFn: (key: string) =>
      client.getObject(create(GetObjectRequestSchema, { bucket: activeBucket, key })),
    onSuccess: (data, key) => {
      const blob = new Blob([data.body as Uint8Array<ArrayBuffer>]);
      const url = URL.createObjectURL(blob);
      const a = document.createElement("a");
      a.href = url;
      a.download = fileNameFromKey(key);
      document.body.appendChild(a);
      a.click();
      document.body.removeChild(a);
      URL.revokeObjectURL(url);
    },
  });

  const deleteObjMutation = useMutation({
    mutationFn: (key: string) =>
      client.deleteObject(create(DeleteObjectRequestSchema, { bucket: activeBucket, key })),
    onSuccess: () => {
      invalidateObjects();
      setShowDeleteObj(false);
      setSelectedObject(null);
      clearObjectSelection();
    },
  });

  const batchDeleteMutation = useMutation({
    mutationFn: (keys: string[]) =>
      client.deleteObjects(
        create(DeleteObjectsRequestSchema, {
          bucket: activeBucket,
          delete: {
            objects: keys.map((k) => ({ key: k })),
            quiet: false,
          },
        }),
      ),
    onSuccess: () => {
      invalidateObjects();
      setSelectedObject(null);
      clearObjectSelection();
    },
  });

  const copyMutation = useMutation({
    mutationFn: () =>
      client.copyObject(
        create(CopyObjectRequestSchema, {
          bucket: copyDestBucket,
          key: copyDestKey,
          copysource: `${activeBucket}/${selectedObject?.key}`,
        }),
      ),
    onSuccess: () => {
      invalidateObjects();
      setShowCopy(false);
      setCopyDestBucket(activeBucket);
      setCopyDestKey("");
    },
  });

  // ── Column definitions ───────────────────────────────────────

  const bucketColumns: ColumnDef<Bucket, any>[] = [
    checkboxColumn<Bucket>(selectedBucketNames, toggleBucket, () => toggleAllBuckets(allBucketIds), allBucketIds, t, (row) => row.name),
    { accessorKey: "name", header: t("services.s3.bucketNameHeader"), cell: MonoCell },
    { accessorKey: "bucketregion", header: t("services.s3.regionHeader"), size: 100 },
    { accessorKey: "bucketarn", header: t("services.s3.arnHeader"), cell: SmallMonoCell },
    { accessorKey: "creationdate", header: t("services.s3.createdHeader"), cell: DateCell },
  ];

  const objectColumns: ColumnDef<ObjectRow, any>[] = [
    checkboxColumn<ObjectRow>(selectedObjectIds, toggleObject, () => toggleAllObjects(allObjectIds), allObjectIds, t, (row) => row.id),
    {
      accessorKey: "displayKey",
      header: t("services.s3.keyHeader"),
      cell: ({ getValue, row }) => {
        const v = getValue() as string;
        return row.original.isFolder ? (
          <span className="cell-mono" style={{ cursor: "pointer" }}>📁 {v}</span>
        ) : (
          <span className="cell-mono">{v}</span>
        );
      },
    },
    { accessorKey: "size", header: t("services.s3.sizeHeader"), size: 100 },
    { accessorKey: "lastModified", header: t("services.s3.lastModifiedHeader"), cell: DateCell, size: 160 },
    { accessorKey: "storageClass", header: t("services.s3.storageClassHeader"), size: 120 },
  ];

  // ── Breadcrumb rendering ─────────────────────────────────────

  const breadcrumbParts = prefix.split("/").filter(Boolean);

  const renderBreadcrumb = () => {
    if (view.type === "buckets") {
      return (
        <div className="toolbar-breadcrumb">
          S3 / <span className="breadcrumb-current">{t("services.s3.backToBuckets")}</span>
        </div>
      );
    }
    return (
      <div className="toolbar-breadcrumb">
        <button className="breadcrumb-link" onClick={navigateToBuckets}>
          {t("services.s3.backToBuckets")}
        </button>
        <span className="breadcrumb-sep">/</span>
        <span className="breadcrumb-current">{view.bucket.name}</span>
        {breadcrumbParts.length > 0 && (
          <>
            <span className="breadcrumb-sep">/</span>
            <button className="breadcrumb-link" onClick={() => handleBreadcrumb(-1)}>
              {t("services.s3.root")}
            </button>
          </>
        )}
        {breadcrumbParts.map((part, i) => (
          <span key={i}>
            <span className="breadcrumb-sep">/</span>
            {i < breadcrumbParts.length - 1 ? (
              <button className="breadcrumb-link" onClick={() => handleBreadcrumb(i)}>
                {part}
              </button>
            ) : (
              <span className="breadcrumb-current">{part}</span>
            )}
          </span>
        ))}
      </div>
    );
  };

  // ── Detail panel rendering ───────────────────────────────────

  const renderDetailPanel = () => {
    if (!selectedFile) {
      return (
        <div className="detail-panel-empty">
          {t("services.s3.noObjectSelected")}
        </div>
      );
    }

    // Build JSON data for the selected object
    const jsonObj: Record<string, unknown> = {
      Key: selectedFile.key,
      Size: Number(selectedFile.rawSize),
      StorageClass: selectedFile.storageClass,
      LastModified: selectedFile.lastModified,
    };

    const renderTabContent = () => {
      switch (detailTab) {
        case "json":
          return <JsonViewer data={jsonObj} />;
        case "raw": {
          if (!canPreview) {
            return <div className="detail-panel-empty">{t("services.s3.previewNotAvailable")}</div>;
          }
          if (previewError) {
            return <div className="detail-panel-empty" style={{ color: "var(--accent-red)" }}>{String(previewError)}</div>;
          }
          if (!previewData) return <div className="detail-panel-empty">{previewFetching ? t("common.loading") : "No data"}</div>;
          const bodyBytes = previewData.body as Uint8Array;
          const text = new TextDecoder().decode(bodyBytes);
          return (
            <pre className="code-block" style={{ padding: 8, fontSize: 12, maxHeight: "100%", overflow: "auto" }}>
              {text}
            </pre>
          );
        }
        case "headers": {
          const headers: Record<string, unknown> = {
            ContentType: inferContentType(selectedFile.key),
            ContentLength: Number(selectedFile.rawSize),
          };
          return <JsonViewer data={headers} />;
        }
      }
    };

    return (
      <div className="detail-panel-content">
        <div className="detail-header">
          <span className="detail-title">📄 {selectedFile.displayKey}</span>
          <div className="detail-tabs">
            <button
              className={`detail-tab ${detailTab === "json" ? "active" : ""}`}
              onClick={() => setDetailTab("json")}
            >
              {t("services.s3.tabJson")}
            </button>
            <button
              className={`detail-tab ${detailTab === "raw" ? "active" : ""}`}
              onClick={() => setDetailTab("raw")}
            >
              {t("services.s3.tabRaw")}
            </button>
            <button
              className={`detail-tab ${detailTab === "headers" ? "active" : ""}`}
              onClick={() => setDetailTab("headers")}
            >
              {t("services.s3.tabHeaders")}
            </button>
          </div>
          <div className="detail-actions">
            <button
              className="btn btn-primary btn-sm"
              disabled={downloadMutation.isPending}
              onClick={() => downloadMutation.mutate(selectedFile.key)}
            >
              {t("services.s3.download")}
            </button>
            <button
              className="btn btn-secondary btn-sm"
              onClick={() => { setShowCopy(true); setCopyDestBucket(activeBucket); setCopyDestKey(selectedFile.key); }}
            >
              {t("services.s3.copyObject")}
            </button>
            <button
              className="btn btn-danger btn-sm"
              onClick={() => setShowDeleteObj(true)}
            >
              {t("common.delete")}
            </button>
          </div>
        </div>
        <div className="detail-body">
          {renderTabContent()}
        </div>
      </div>
    );
  };

  // ── Toolbar actions ──────────────────────────────────────────

  const renderActions = () => {
    if (view.type === "buckets") {
      return (
        <>
          <button className="btn btn-primary" onClick={() => setShowCreate(true)}>
            {t("services.s3.create")}
          </button>
          <button
            className="btn btn-danger"
            disabled={selectedBucketNames.size === 0}
            onClick={() => setShowDeleteBuckets(true)}
          >
            {t("common.delete")}
            {selectedBucketNames.size > 0 && (
              <span style={{ marginLeft: 4, opacity: 0.8 }}>
                ({selectedBucketNames.size})
              </span>
            )}
          </button>
        </>
      );
    }
    return (
      <>
        <button className="btn btn-primary btn-sm" onClick={() => setShowUpload(true)}>
          {t("services.s3.upload")}
        </button>
        <button
          className="btn btn-danger btn-sm"
          disabled={selectedObjectIds.size === 0}
          onClick={() => {
            const keys = objectRows
              .filter((r) => selectedObjectIds.has(r.id) && !r.isFolder)
              .map((r) => r.key);
            if (keys.length > 0) batchDeleteMutation.mutate(keys);
          }}
        >
          {t("services.s3.deleteSelected")}
          {selectedObjectIds.size > 0 && (
            <span style={{ marginLeft: 4, opacity: 0.8 }}>
              ({selectedObjectIds.size})
            </span>
          )}
        </button>
      </>
    );
  };

  // ── Row click handler ────────────────────────────────────────

  const handleBucketRowClick = (row: Bucket) => {
    navigateToBucket(row);
  };

  const handleObjectRowClick = (row: ObjectRow) => {
    if (row.isFolder) {
      handleFolderClick(row.key);
    } else {
      setSelectedObject(row);
    }
  };

  // ── Render ───────────────────────────────────────────────────

  return (
    <ServicePageLayout
      icon="📦"
      title={t("services.s3.title")}
      isLoading={isLoading && view.type === "buckets"}
      error={error}
      count={view.type === "buckets" ? buckets.length : undefined}
      countLabel={t("services.s3.countLabel")}
      actions={renderActions()}
    >
      {/* Breadcrumb toolbar */}
      <div className="inspector-toolbar">
        {renderBreadcrumb()}
        <div className="toolbar-selection-info">
          {view.type === "buckets" && selectedBucketNames.size > 0 && (
            <span className="selection-count">
              {t("services.s3.selectedCount", { count: selectedBucketNames.size })}
            </span>
          )}
          {view.type === "objects" && selectedObjectIds.size > 0 && (
            <span className="selection-count">
              {t("services.s3.selectedCount", { count: selectedObjectIds.size })}
            </span>
          )}
        </div>
      </div>

      {/* Main content area */}
      {view.type === "buckets" ? (
        /* Bucket list — full width, no bottom detail */
        <DataTable
          columns={bucketColumns}
          data={buckets}
          getRowId={(row) => row.name}
          onRowClick={handleBucketRowClick}
        />
      ) : (
        /* Object browser — table + bottom detail panel (drag split) */
        objectRows.length > 0 ? (
          <Splitter
            direction="horizontal"
            initialSize={240}
            minSize={80}
            maxSize={600}
            storageKey="vs-split-s3-detail"
          >
            <div style={{ flex: 1, minHeight: 0, overflow: "auto" }}>
              <DataTable
                columns={objectColumns}
                data={objectRows}
                getRowId={(row) => row.id}
                onRowClick={handleObjectRowClick}
                selectedId={selectedObject?.id}
              />
              {nextToken && (
                <div style={{ textAlign: "center", padding: 8 }}>
                  <button
                    className="btn btn-secondary btn-sm"
                    onClick={handleLoadMore}
                    disabled={listFetching}
                  >
                    {listFetching ? t("common.loading") : t("services.s3.loadMore")}
                  </button>
                </div>
              )}
            </div>
            {renderDetailPanel()}
          </Splitter>
        ) : (
          <div className="empty-state">{t("services.s3.emptyBucket")}</div>
        )
      )}

      {/* ── Modals ──────────────────────────────────────────── */}

      {/* Create Bucket */}
      <ServiceCreateModal
        open={showCreate}
        onClose={() => setShowCreate(false)}
        title={t("services.s3.create")}
        error={createMutation.error}
        isPending={createMutation.isPending}
        onCreate={() => createMutation.mutate()}
        disabled={!formBucket}
      >
        <label>
          {t("services.s3.nameField")}
          <input
            value={formBucket}
            onChange={(e) => setFormBucket(e.target.value)}
            placeholder={t("services.s3.placeholder")}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.s3.objectOwnershipLabel")}
          <select
            value={formOwnership}
            onChange={(e) => setFormOwnership(e.target.value)}
            className="modal-input"
          >
            <option value="BUCKETOWNERENFORCED">{t("services.s3.ownershipBucketOwnerEnforced")}</option>
            <option value="BUCKETOWNERPREFERRED">{t("services.s3.ownershipBucketOwnerPreferred")}</option>
            <option value="OBJECTWRITER">{t("services.s3.ownershipObjectWriter")}</option>
          </select>
        </label>
      </ServiceCreateModal>

      {/* Delete Bucket(s) */}
      <ServiceDeleteDialog
        open={showDeleteBuckets}
        title={t("services.s3.delete")}
        name={Array.from(selectedBucketNames).join(", ")}
        error={deleteBucketMutation.error}
        isPending={deleteBucketMutation.isPending}
        onConfirm={() => {
          for (const name of selectedBucketNames) {
            deleteBucketMutation.mutate(name);
          }
        }}
        onClose={() => setShowDeleteBuckets(false)}
      />

      {/* Upload Object */}
      <Modal open={showUpload} onClose={() => { setShowUpload(false); setUploadError(""); }}>
        <h2>{t("services.s3.uploadTitle")}</h2>
        {uploadError && <div className="modal-error">{uploadError}</div>}
        {uploadMutation.error && <div className="modal-error">{String(uploadMutation.error)}</div>}
        <label>
          {t("services.s3.selectFile")}
          <input type="file" onChange={(e: ChangeEvent<HTMLInputElement>) => {
            const f = e.target.files?.[0];
            if (f) {
              setUploadFile(f);
              setUploadKey(f.name);
              setUploadError("");
            }
          }} className="modal-input" />
        </label>
        {uploadFile && (
          <label>
            {t("services.s3.keyName")}
            <input
              value={uploadKey}
              onChange={(e) => setUploadKey(e.target.value)}
              className="modal-input"
            />
          </label>
        )}
        <div style={{ fontSize: "0.85em", color: "var(--text-muted)" }}>
          {t("services.s3.uploadPath")}: {prefix}[{uploadKey || "..."}]
        </div>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => { setShowUpload(false); setUploadError(""); }}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!uploadFile || !uploadKey || uploadMutation.isPending}
            onClick={() => uploadMutation.mutate()}
          >
            {uploadMutation.isPending ? t("services.s3.uploading") : t("services.s3.upload")}
          </button>
        </div>
      </Modal>

      {/* Delete Single Object */}
      <Modal open={showDeleteObj && !!selectedObject} onClose={() => setShowDeleteObj(false)}>
        <h2>{t("services.s3.deleteObject")}</h2>
        <p>{t("services.s3.deleteConfirm")}</p>
        <span className="cell-mono">{selectedObject?.key}</span>
        {deleteObjMutation.error && <div className="modal-error">{String(deleteObjMutation.error)}</div>}
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowDeleteObj(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-danger"
            disabled={deleteObjMutation.isPending}
            onClick={() => selectedObject && deleteObjMutation.mutate(selectedObject.key)}
          >
            {deleteObjMutation.isPending ? t("common.deleting") : t("common.delete")}
          </button>
        </div>
      </Modal>

      {/* Copy Object */}
      <Modal open={showCopy} onClose={() => setShowCopy(false)}>
        <h2>{t("services.s3.copyObject")}</h2>
        {copyMutation.error && <div className="modal-error">{String(copyMutation.error)}</div>}
        <div className="detail-field">
          <span className="detail-label">{t("services.s3.sourceBucket")}</span>
          <span className="cell-mono">{activeBucket}</span>
        </div>
        <div className="detail-field">
          <span className="detail-label">{t("services.s3.sourceKey")}</span>
          <span className="cell-mono">{selectedObject?.key}</span>
        </div>
        <label>
          {t("services.s3.destBucket")}
          <input
            value={copyDestBucket}
            onChange={(e) => setCopyDestBucket(e.target.value)}
            className="modal-input"
          />
        </label>
        <label>
          {t("services.s3.destKey")}
          <input
            value={copyDestKey}
            onChange={(e) => setCopyDestKey(e.target.value)}
            className="modal-input"
          />
        </label>
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={() => setShowCopy(false)}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={!copyDestBucket || !copyDestKey || copyMutation.isPending}
            onClick={() => copyMutation.mutate()}
          >
            {copyMutation.isPending ? t("services.s3.copying") : t("services.s3.copyObject")}
          </button>
        </div>
      </Modal>
    </ServicePageLayout>
  );
}
