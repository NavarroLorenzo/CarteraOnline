import { Navigate, Outlet } from "react-router-dom";
import { useServerWakeState } from "../api/useServerWakeState";
import { useAuth } from "./AuthContext";

export function PublicOnlyRoute() {
  const { status, isAuthenticated } = useAuth();
  const { isWaking } = useServerWakeState();

  if (status === "loading") {
    return (
      <div className="route-state">
        <div className="route-state__card">
          <span className="eyebrow">Preparando vista</span>
          <h1>Estamos preparando el acceso</h1>
          <p>
            {isWaking
              ? "El servicio se está iniciando. Te vamos a mostrar la vista apenas esté lista."
              : "Estamos recuperando el estado de tu sesión."}
          </p>
        </div>
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
