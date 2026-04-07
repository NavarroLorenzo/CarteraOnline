import { useState } from "react";
import { motion } from "framer-motion";
import { Outlet, useLocation } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { useAuth } from "../../auth/AuthContext";
import { RouteTransition, buttonHover } from "../ui/animation";

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
        <motion.header
          className="topbar"
          initial={{ opacity: 0, y: -14 }}
          animate={{ opacity: 1, y: 0 }}
          transition={{ duration: 0.42, ease: [0.22, 1, 0.36, 1] }}
        >
          <motion.button
            type="button"
            className="topbar__menu-button"
            onClick={() => setMenuOpen((current) => !current)}
            whileHover={buttonHover}
            whileTap={{ scale: 0.98 }}
          >
            Menu
          </motion.button>

          <motion.div
            className="topbar__copy"
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.46, delay: 0.04, ease: [0.22, 1, 0.36, 1] }}
          >
            <span className="eyebrow">Espacio personal</span>
            <h2>Solo estás viendo tus propios datos</h2>
          </motion.div>

          <motion.button
            type="button"
            className="ghost-button"
            onClick={logout}
            whileHover={buttonHover}
            whileTap={{ scale: 0.98 }}
            initial={{ opacity: 0, y: 12 }}
            animate={{ opacity: 1, y: 0 }}
            transition={{ duration: 0.46, delay: 0.08, ease: [0.22, 1, 0.36, 1] }}
          >
            Cerrar sesión
          </motion.button>
        </motion.header>

        <div className="page-content">
          <RouteTransition routeKey={location.pathname}>
            <Outlet />
          </RouteTransition>
        </div>
      </main>
    </div>
  );
}
