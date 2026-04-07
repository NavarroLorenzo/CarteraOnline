import { motion } from "framer-motion";
import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";
import {
  PresenceMessage,
  StaggerGroup,
  buttonHover,
  cardHover,
  fadeLeft,
  fadeRight,
  fadeUp,
  scaleIn,
} from "../components/ui/animation";

const loginStats = [
  {
    title: "Sesion segura",
    description: "Acceso protegido con token",
  },
  {
    title: "Tu espacio",
    description: "Cada usuario ve solo lo suyo",
  },
  {
    title: "Todo conectado",
    description: "Cuentas, movimientos y transferencias",
  },
];

export function LoginPage() {
  const navigate = useNavigate();
  const location = useLocation();
  const { login } = useAuth();

  const [identifier, setIdentifier] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const from = (location.state as { from?: string } | null)?.from || "/dashboard";

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      await login({
        identifier,
        password,
      });
      navigate(from, { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo iniciar sesión");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <StaggerGroup className="auth-shell" onView={false}>
      <motion.section className="auth-panel auth-panel--hero" variants={fadeLeft}>
        <motion.span className="eyebrow" variants={fadeUp}>
          Finanzas en orden
        </motion.span>
        <motion.h1 variants={fadeUp}>Entrá y retomá el control de tu dinero.</motion.h1>
        <motion.p variants={fadeUp}>
          Accedé a tus cuentas, movimientos y transferencias desde un solo lugar, con tu sesión protegida.
        </motion.p>

        <motion.div className="auth-hero__stats" variants={fadeUp}>
          {loginStats.map((item) => (
            <motion.div key={item.title} variants={scaleIn} whileHover={cardHover}>
              <strong>{item.title}</strong>
              <span>{item.description}</span>
            </motion.div>
          ))}
        </motion.div>
      </motion.section>

      <motion.section className="auth-panel auth-panel--form auth-panel--form-login" variants={fadeRight}>
        <motion.div className="auth-panel__actions" variants={fadeUp}>
          <motion.div whileHover={buttonHover} whileTap={{ scale: 0.99 }}>
            <Link to="/" className="ghost-button auth-back-button">
              Volver al inicio
            </Link>
          </motion.div>
        </motion.div>

        <motion.div className="auth-form__header" variants={fadeUp}>
          <span className="eyebrow">Login</span>
          <h2>Bienvenido de nuevo</h2>
          <p>Podés ingresar con tu email o con tu nombre de usuario.</p>
        </motion.div>

        <motion.form className="stack-form" onSubmit={handleSubmit} variants={fadeUp}>
          <motion.label className="field" variants={fadeUp}>
            <span>Usuario o email</span>
            <input
              type="text"
              value={identifier}
              onChange={(event) => setIdentifier(event.target.value)}
              placeholder="usuario1 o usuario1@mail.com"
              required
            />
          </motion.label>

          <motion.label className="field" variants={fadeUp}>
            <span>Contraseña</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="••••••"
              required
            />
          </motion.label>

          <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>

          <motion.div variants={fadeUp} whileHover={buttonHover} whileTap={{ scale: 0.995 }}>
            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting ? "Ingresando..." : "Iniciar sesión"}
            </button>
          </motion.div>
        </motion.form>

        <motion.p className="auth-switch" variants={fadeUp}>
          ¿No tenés cuenta? <Link to="/register">Crear usuario</Link>
        </motion.p>
      </motion.section>
    </StaggerGroup>
  );
}
