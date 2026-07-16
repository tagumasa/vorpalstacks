/**
 * Reusable tag display and editor for service detail panels.
 * Each service wires its own RPC calls; this component handles
 * rendering and the add/remove interaction.
 */
import { useState } from "react";
import { useTranslation } from "react-i18next";
import { useQuery, useMutation, useQueryClient } from "@tanstack/react-query";
import { Modal } from "./modal";

export interface TagPair {
  key: string;
  value: string;
}

export interface TagSectionConfig {
  /** Query key prefix for tag list queries. */
  queryKeyBase: unknown[];
  /** Fetch tags for the given resource ARN. Returns [] if no tags. */
  fetchTags: (resourceArn: string) => Promise<TagPair[]>;
  /** Add tags to a resource. */
  tagResource: (resourceArn: string, tags: TagPair[]) => Promise<unknown>;
  /** Remove tag keys from a resource. */
  untagResource: (resourceArn: string, tagKeys: string[]) => Promise<unknown>;
}

/**
 * Hook for managing tags on a selected resource.
 * Automatically fetches tags when resourceArn changes and provides
 * add/remove mutations with cache invalidation.
 */
export function useTags(
  cfg: TagSectionConfig,
  resourceArn: string | undefined,
) {
  const qc = useQueryClient();
  const queryKey = [...cfg.queryKeyBase, resourceArn ?? ""];

  const { data: tags, isLoading } = useQuery({
    queryKey,
    queryFn: () => cfg.fetchTags(resourceArn!),
    enabled: !!resourceArn,
  });

  const tagMut = useMutation({
    mutationFn: (newTags: TagPair[]) => cfg.tagResource(resourceArn!, newTags),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
  });

  const untagMut = useMutation({
    mutationFn: (keys: string[]) => cfg.untagResource(resourceArn!, keys),
    onSuccess: () => qc.invalidateQueries({ queryKey }),
  });

  return {
    tags: tags ?? [],
    isLoading,
    addTags: async (newTags: TagPair[]) => { await tagMut.mutateAsync(newTags); },
    removeTag: (key: string) => untagMut.mutate([key]),
    isPending: tagMut.isPending || untagMut.isPending,
    isError: tagMut.isError || untagMut.isError,
  };
}

/** Renders a single tag as a badge with optional remove button. */
function TagBadge({
  tagKey,
  tagValue,
  onRemove,
}: {
  tagKey: string;
  tagValue: string;
  onRemove?: () => void;
}) {
  return (
    <span className="tag-badge">
      <span className="tag-key">{tagKey}</span>
      <span className="tag-sep">:</span>
      <span className="tag-value">{tagValue}</span>
      {onRemove && (
        <button
          className="tag-remove"
          onClick={onRemove}
          title="Remove tag"
        >
          {"\u00d7"}
        </button>
      )}
    </span>
  );
}

/**
 * Tag section for service detail panels. Shows tags and provides
 * an inline editor modal for adding/removing tags.
 */
export function TagSection({
  tags,
  isLoading,
  onAddTags,
  onRemoveTag,
  isPending,
  isError,
  readOnly,
}: {
  tags: TagPair[];
  isLoading?: boolean;
  onAddTags: (tags: TagPair[]) => Promise<void>;
  onRemoveTag: (key: string) => void;
  isPending?: boolean;
  isError?: boolean;
  readOnly?: boolean;
}) {
  const { t } = useTranslation();
  const [showEditor, setShowEditor] = useState(false);
  const [newKey, setNewKey] = useState("");
  const [newValue, setNewValue] = useState("");
  const [pendingTags, setPendingTags] = useState<TagPair[]>([]);
  const [saveError, setSaveError] = useState(false);

  const handleAddPending = () => {
    const k = newKey.trim();
    if (!k) return;
    setPendingTags([...pendingTags.filter((t) => t.key !== k), { key: k, value: newValue.trim() }]);
    setNewKey("");
    setNewValue("");
  };

  const handleRemovePending = (key: string) => {
    setPendingTags(pendingTags.filter((t) => t.key !== key));
  };

  const handleSave = async () => {
    setSaveError(false);
    try {
      await onAddTags(pendingTags);
      setPendingTags([]);
      setShowEditor(false);
    } catch {
      setSaveError(true);
    }
  };

  const handleClose = () => {
    setPendingTags([]);
    setNewKey("");
    setNewValue("");
    setShowEditor(false);
  };

  if (isLoading) {
    return (
      <section className="detail-section">
        <h3>{t("common.fields.tags")}</h3>
        <div className="loading-state">{t("common.loading")}</div>
      </section>
    );
  }

  const allTags = [...tags, ...pendingTags.filter((pt) => !tags.some((t) => t.key === pt.key))];

  return (
    <section className="detail-section">
      <h3>
        {t("common.fields.tags")}
        {!readOnly && (
          <button
            className="btn btn-secondary btn-sm"
            style={{ marginLeft: 8 }}
            onClick={() => {
              setPendingTags([]);
              setShowEditor(true);
            }}
          >
            {t("common.edit")}
          </button>
        )}
      </h3>
      {allTags.length === 0 ? (
        <span className="text-muted">{"\u2014"}</span>
      ) : (
        <div className="tag-list">
          {allTags.map((tag) => (
            <TagBadge
              key={tag.key}
              tagKey={tag.key}
              tagValue={tag.value || "\u2014"}
              onRemove={
                showEditor
                  ? undefined
                  : () => onRemoveTag(tag.key)
              }
            />
          ))}
        </div>
      )}

      <Modal open={showEditor} onClose={handleClose}>
        <h2>{t("common.editTags")}</h2>
        <div className="tag-editor-existing">
          {tags.map((tag) => (
            <TagBadge
              key={tag.key}
              tagKey={tag.key}
              tagValue={tag.value || "\u2014"}
              onRemove={() => {
                onRemoveTag(tag.key);
              }}
            />
          ))}
        </div>
        <div className="tag-editor-add" style={{ display: "flex", gap: 8, marginTop: 12 }}>
          <input
            className="modal-input"
            placeholder={t("common.fields.key")}
            value={newKey}
            onChange={(e) => setNewKey(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") handleAddPending();
            }}
            style={{ flex: "0 0 120px" }}
          />
          <input
            className="modal-input"
            placeholder={t("common.fields.value")}
            value={newValue}
            onChange={(e) => setNewValue(e.target.value)}
            onKeyDown={(e) => {
              if (e.key === "Enter") handleAddPending();
            }}
            style={{ flex: 1 }}
          />
          <button
            className="btn btn-secondary btn-sm"
            onClick={handleAddPending}
            disabled={!newKey.trim()}
          >
            {t("common.add")}
          </button>
        </div>
        {pendingTags.length > 0 && (
          <div className="tag-editor-pending" style={{ marginTop: 8 }}>
            {pendingTags.map((tag) => (
              <TagBadge
                key={tag.key}
                tagKey={tag.key}
                tagValue={tag.value || "\u2014"}
                onRemove={() => handleRemovePending(tag.key)}
              />
            ))}
          </div>
        )}
        {(saveError || isError) && (
          <div className="modal-error" style={{ marginTop: 8, color: "var(--color-danger, #e53e3e)" }}>
            {t("common.saveFailed")}
          </div>
        )}
        <div className="modal-actions">
          <button className="btn btn-secondary" onClick={handleClose}>
            {t("common.cancel")}
          </button>
          <button
            className="btn btn-primary"
            disabled={isPending || pendingTags.length === 0}
            onClick={handleSave}
          >
            {isPending ? t("common.saving") : t("common.save")}
          </button>
        </div>
      </Modal>
    </section>
  );
}
