import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useServerWakeState } from "../api/useServerWakeState";
import { useAuth } from "./AuthContext";

export function ProtectedRoute() {
  const location = useLocation();
  const { status, isAuthenticated } = useAuth();
  const { isWaking } = useServerWakeState();

  if (status === "loading") {
    return (
      <div className="route-state">
        <div className="route-state__card">
          <span className="eyebrow">Sincronizando sesión</span>
          <h1>Estamos validando tu acceso</h1>
          <p>
            {isWaking
              ? "El servicio se está iniciando. Vamos a ingresar automáticamente cuando termine de responder."
              : "Un momento. Estamos comprobando tu acceso para mostrarte la información correcta."}
          </p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <Outlet />;
}
