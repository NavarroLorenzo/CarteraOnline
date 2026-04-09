import { useState } from "react";
import { motion } from "framer-motion";
import { Outlet } from "react-router-dom";
import { Sidebar } from "./Sidebar";
import { useServerWakeState } from "../../api/useServerWakeState";
import { useAuth } from "../../auth/AuthContext";
import { PresenceMessage, buttonHover } from "../ui/animation";

export function AppLayout() {
  const { user, logout } = useAuth();
  const { isWaking } = useServerWakeState();
  const [menuOpen, setMenuOpen] = useState(false);

  if (!user) {
    return (
      <div className="route-state">
        <div className="route-state__card">
          <span className="eyebrow">Sincronizando sesión</span>
          <h1>Estamos preparando tu espacio</h1>
          <p>
            {isWaking
              ? "El servicio se está iniciando. Tu información aparecerá automáticamente en unos segundos."
              : "Estamos cargando tu información para abrir el panel."}
          </p>
        </div>
      </div>
    );
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
            Menú
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

        <PresenceMessage className="feedback feedback--warning">
          {isWaking ? "Iniciando servicio, por favor esperá unos segundos..." : null}
        </PresenceMessage>

        <div className="page-content">
          <Outlet />
        </div>
      </main>
    </div>
  );
}
