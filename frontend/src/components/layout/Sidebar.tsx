import { NavLink } from "react-router-dom";

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
    <aside className="sidebar">
      <div className="sidebar__brand">
        <div className="sidebar__brand-mark">CO</div>
        <div>
          <p className="sidebar__eyebrow">Finanzas personales</p>
          <h1>Cartera Online</h1>
        </div>
      </div>

      <nav className="sidebar__nav">
        {navItems.map((item) => (
          <NavLink
            key={item.to}
            to={item.to}
            className={({ isActive }) =>
              isActive ? "sidebar__link sidebar__link--active" : "sidebar__link"
            }
            onClick={onNavigate}
          >
            {item.label}
          </NavLink>
        ))}
      </nav>

      <div className="sidebar__profile">
        <span className="sidebar__eyebrow">Sesión iniciada</span>
        <strong>{username}</strong>
        <p>{email}</p>
      </div>
    </aside>
  );
}
