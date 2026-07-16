/**
 * Error Boundary component that catches rendering errors in child components.
 * Displays a recoverable error message instead of crashing the entire app.
 */
import { Component, type ErrorInfo, type ReactNode } from "react";
import i18n from "@/i18n";

interface Props {
  children: ReactNode;
}

interface State {
  hasError: boolean;
  error: Error | null;
}

/**
 * Catches JavaScript errors anywhere in its child component tree,
 * logs those errors, and displays a fallback UI instead of crashing.
 */
export class ErrorBoundary extends Component<Props, State> {
  constructor(props: Props) {
    super(props);
    this.state = { hasError: false, error: null };
  }

  static getDerivedStateFromError(error: Error): State {
    return { hasError: true, error };
  }

  componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("[ErrorBoundary]", error, errorInfo);
  }

  render() {
    if (!this.state.hasError) return this.props.children;

    return (
      <div className="content-area" style={{ padding: 24 }}>
        <h1 style={{ color: "var(--accent-red)", fontSize: 14, marginBottom: 12 }}>
          {i18n.t("common.errorTitle")}
        </h1>
        <pre
          style={{
            fontSize: 11,
            fontFamily: "var(--pixel-font)",
            color: "var(--text-muted)",
            whiteSpace: "pre-wrap",
            marginBottom: 16,
          }}
        >
          {this.state.error?.message ?? i18n.t("common.errorUnknown")}
        </pre>
        <button
          className="btn btn-primary"
          onClick={() => this.setState({ hasError: false, error: null })}
        >
          {i18n.t("common.errorRetry")}
        </button>
      </div>
    );
  }
}
