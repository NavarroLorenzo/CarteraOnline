import { useState } from "react";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { useAuth } from "../../auth/AuthContext";

export function AppLayout() {
  const { user, logout } = useAuth();
  const [menuOpen, setMenuOpen] = useState(false);

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
            <span className="eyebrow">Workspace personal</span>
            <h2>Solo estás viendo tus datos</h2>
          </div>

          <button type="button" className="ghost-button" onClick={logout}>
            Cerrar sesión
          </button>
        </header>

        <div className="page-content">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
