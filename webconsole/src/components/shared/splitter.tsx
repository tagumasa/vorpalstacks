/**
 * Draggable splitter component for resizing two panels.
 * Used by SplitPane in service-page.tsx to allow users to resize
 * the table/detail panel divider.
 */
import { useState, useRef, useCallback, useEffect, type ReactNode } from "react";

interface SplitterProps {
  initialSize?: number;
  minSize?: number;
  maxSize?: number;
  direction?: "horizontal" | "vertical";
  storageKey?: string;
  children: [ReactNode, ReactNode];
}

export function Splitter({
  initialSize = 360,
  minSize = 200,
  maxSize = 800,
  direction = "vertical",
  storageKey,
  children,
}: SplitterProps) {
  const savedSize = storageKey
    ? parseInt(localStorage.getItem(storageKey) || "", 10)
    : 0;
  const [size, setSize] = useState(savedSize || initialSize);
  const dragging = useRef(false);
  const startPos = useRef(0);
  const startSize = useRef(0);
  const sizeRef = useRef(size);
  sizeRef.current = size;

  const onMouseDown = useCallback(
    (e: React.MouseEvent) => {
      e.preventDefault();
      dragging.current = true;
      startPos.current = direction === "vertical" ? e.clientX : e.clientY;
      startSize.current = sizeRef.current;
      document.body.style.cursor =
        direction === "vertical" ? "col-resize" : "row-resize";
      document.body.style.userSelect = "none";
    },
    [direction],
  );

  useEffect(() => {
    const onMouseMove = (e: MouseEvent) => {
      if (!dragging.current) return;
      const delta =
        direction === "vertical"
          ? e.clientX - startPos.current
          : e.clientY - startPos.current;
      const next = Math.min(maxSize, Math.max(minSize, startSize.current + delta));
      setSize(next);
    };

    const onMouseUp = () => {
      if (!dragging.current) return;
      dragging.current = false;
      document.body.style.cursor = "";
      document.body.style.userSelect = "";
      if (storageKey) {
        localStorage.setItem(storageKey, String(sizeRef.current));
      }
    };

    window.addEventListener("mousemove", onMouseMove);
    window.addEventListener("mouseup", onMouseUp);
    return () => {
      window.removeEventListener("mousemove", onMouseMove);
      window.removeEventListener("mouseup", onMouseUp);
    };
  }, [direction, minSize, maxSize, storageKey]);

  const isVertical = direction === "vertical";

  return (
    <div
      style={{
        display: "flex",
        flexDirection: isVertical ? "row" : "column",
        flex: 1,
        minHeight: 0,
      }}
    >
      <div style={{ flex: 1, minWidth: 0, minHeight: 0, overflow: "hidden" }}>
        {children[0]}
      </div>
      <div
        onMouseDown={onMouseDown}
        style={{
          flexShrink: 0,
          width: isVertical ? 4 : "auto",
          height: isVertical ? "auto" : 4,
          cursor: isVertical ? "col-resize" : "row-resize",
          background: "var(--border-dim)",
          position: "relative",
          display: "flex",
          alignItems: "center",
          justifyContent: "center",
        }}
      >
        <div
          style={{
            width: isVertical ? 2 : 16,
            height: isVertical ? 16 : 2,
            borderRadius: 1,
            background: "var(--text-muted)",
            opacity: 0.5,
          }}
        />
      </div>
      <div
        style={{
          flexShrink: 0,
          width: isVertical ? size : "auto",
          height: isVertical ? "auto" : size,
          overflow: "auto",
        }}
      >
        {children[1]}
      </div>
    </div>
  );
}
