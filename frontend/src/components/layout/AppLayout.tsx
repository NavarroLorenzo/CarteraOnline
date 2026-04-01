import { useState } from "react";
import { Outlet, useLocation } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { useAuth } from "../../auth/AuthContext";
import { RouteTransition } from "../ui/animation";

export function AppLayout() {
  const { user, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);
  const location = useLocation();

  if (!user) {
    return null;
  }

  return (
    <div className={`app-shell ${menuOpen ? "app-shell--menu-open" : ""}`}>
      <Sidebar
        username={user.username}
        email={user.email}
        onNavigate={() => setMenuOpen(false)}
      />

      <div className="app-shell__backdrop" onClick={() => setMenuOpen(false)} />

      <main className="app-main">
        <header className="topbar">
          <button
            type="button"
            className="topbar__menu-button"
            onClick={() => setMenuOpen((current) => !current)}
          >
            Menu
          </button>

          <div className="topbar__copy">
            <span className="eyebrow">Espacio personal</span>
            <h2>Solo estás viendo tus propios datos</h2>
          </div>

          <button type="button" className="ghost-button" onClick={logout}>
            Cerrar sesión
          </button>
        </header>

        <div className="page-content">
          <RouteTransition routeKey={location.pathname}>
            <Outlet />
          </RouteTransition>
        </div>
      </main>
    </div>
  );
}
