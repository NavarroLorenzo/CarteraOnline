import { Component, type ErrorInfo, type PropsWithChildren, type ReactNode } from "react";

type AppErrorBoundaryState = {
  hasError: boolean;
};

export class AppErrorBoundary extends Component<PropsWithChildren, AppErrorBoundaryState> {
  state: AppErrorBoundaryState = {
    hasError: false,
  };

  static getDerivedStateFromError(): AppErrorBoundaryState {
    return { hasError: true };
  }

  override componentDidCatch(error: Error, errorInfo: ErrorInfo) {
    console.error("App render error", error, errorInfo);
  }

  override render(): ReactNode {
    if (this.state.hasError) {
      return (
        <div className="route-state">
          <div className="route-state__card">
            <span className="eyebrow">Recuperación</span>
            <h1>No pudimos mostrar esta pantalla</h1>
            <p>
              Ocurrió un inconveniente al cargar esta vista. Podés intentarlo nuevamente para seguir usando la app.
            </p>
            <button
              type="button"
              className="primary-button"
              onClick={() => window.location.reload()}
            >
              Recargar
            </button>
          </div>
        </div>
      );
    }

    return this.props.children;
  }
}
