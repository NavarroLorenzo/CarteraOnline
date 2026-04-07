import { motion } from "framer-motion";
import { NavLink } from "react-router-dom";
import { fadeLeft, fadeUp, staggerContainer } from "../ui/animation";

const navItems = [
  { to: "/dashboard", label: "Dashboard" },
  { to: "/accounts", label: "Cuentas" },
  { to: "/transactions", label: "Transacciones" },
  { to: "/transfers", label: "Transferencias" },
];

type SidebarProps = {
  username: string;
  email: string;
  onNavigate?: () => void;
};

export function Sidebar({ username, email, onNavigate }: SidebarProps) {
  return (
    <motion.aside
      className="sidebar"
      initial="hidden"
      animate="visible"
      variants={staggerContainer}
    >
      <motion.div className="sidebar__brand" variants={fadeLeft}>
        <div className="sidebar__brand-mark">CO</div>
        <div>
          <p className="sidebar__eyebrow">Finanzas personales</p>
          <h1>Cartera Online</h1>
        </div>
      </motion.div>

      <motion.nav className="sidebar__nav" variants={fadeUp}>
        {navItems.map((item) => (
          <motion.div key={item.to} variants={fadeUp}>
            <NavLink
              to={item.to}
              className={({ isActive }) =>
                isActive ? "sidebar__link sidebar__link--active" : "sidebar__link"
              }
              onClick={onNavigate}
            >
              {({ isActive }) => (
                <>
                  {isActive ? (
                    <motion.span
                      layoutId="sidebar-active-pill"
                      className="sidebar__link-pill"
                      transition={{ type: "spring", stiffness: 360, damping: 32 }}
                    />
                  ) : null}
                  <span className="sidebar__link-label">{item.label}</span>
                </>
              )}
            </NavLink>
          </motion.div>
        ))}
      </motion.nav>

      <motion.div className="sidebar__profile" variants={fadeUp}>
        <span className="sidebar__eyebrow">Sesión iniciada</span>
        <strong>{username}</strong>
        <p>{email}</p>
      </motion.div>
    </motion.aside>
  );
}
