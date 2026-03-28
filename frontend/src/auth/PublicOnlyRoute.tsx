import { Navigate, Outlet } from "react-router-dom";
import { useAuth } from "./AuthContext";

export function PublicOnlyRoute() {
  const { status, isAuthenticated } = useAuth();

  if (status === "loading") {
    return (
      <div className="route-state">
        <div className="route-state__card">
          <span className="eyebrow">Preparando vista</span>
          <h1>Cargando cartera</h1>
          <p>Estamos recuperando tu estado de sesión.</p>
        </div>
      </div>
    );
  }

  if (isAuthenticated) {
    return <Navigate to="/dashboard" replace />;
  }

  return <Outlet />;
}
