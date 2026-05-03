/**
 * Reusable dialog components built on the shared Modal primitive.
 * ConfirmDialog handles destructive confirmations (delete, etc).
 */
import { Modal } from "./modal";

interface ConfirmDialogProps {
  open: boolean;
  title: string;
  message: React.ReactNode;
  confirmLabel?: string;
  isPending?: boolean;
  error?: Error | null;
  onConfirm: () => void;
  onClose: () => void;
}

export function ConfirmDialog({
  open,
  title,
  message,
  confirmLabel = "Delete",
  isPending = false,
  error,
  onConfirm,
  onClose,
}: ConfirmDialogProps) {
  return (
    <Modal open={open} onClose={onClose}>
      <h2>{title}</h2>
      <p>{message}</p>
      {error && <div className="modal-error">{String(error)}</div>}
      <div className="modal-actions">
        <button className="btn btn-secondary" onClick={onClose}>
          Cancel
        </button>
        <button
          className="btn btn-danger"
          disabled={isPending}
          onClick={onConfirm}
        >
          {isPending ? `${confirmLabel}ing...` : confirmLabel}
        </button>
      </div>
    </Modal>
  );
}
