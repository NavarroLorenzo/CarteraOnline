import { Navigate, Outlet, useLocation } from "react-router-dom";
import { useAuth } from "./AuthContext";

export function ProtectedRoute() {
  const location = useLocation();
  const { status, isAuthenticated } = useAuth();

  if (status === "loading") {
    return (
      <div className="route-state">
        <div className="route-state__card">
          <span className="eyebrow">Sincronizando sesión</span>
          <h1>Estamos validando tu acceso</h1>
          <p>Un momento. Estamos comprobando tu token antes de abrir tus datos.</p>
        </div>
      </div>
    );
  }

  if (!isAuthenticated) {
    return <Navigate to="/login" replace state={{ from: location.pathname }} />;
  }

  return <Outlet />;
}
