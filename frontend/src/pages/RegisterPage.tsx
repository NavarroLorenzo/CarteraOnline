import { motion } from "framer-motion";
import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
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

const registerStats = [
  {
    title: "Alta simple",
    description: "Creá tu usuario y entrá directo a tu panel.",
  },
  {
    title: "Base ordenada",
    description: "Cuentas, movimientos y transferencias desde el primer día.",
  },
  {
    title: "Seguimiento claro",
    description: "Cada cambio queda visible en tu historial personal.",
  },
];

export function RegisterPage() {
  const navigate = useNavigate();
  const { register } = useAuth();

  const [email, setEmail] = useState("");
  const [username, setUsername] = useState("");
  const [password, setPassword] = useState("");
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);

    try {
      await register({
        email,
        username,
        password,
      });
      navigate("/dashboard", { replace: true });
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo crear la cuenta");
    } finally {
      setSubmitting(false);
    }
  };

  return (
    <StaggerGroup className="auth-shell" onView={false}>
      <motion.section className="auth-panel auth-panel--hero auth-panel--hero-alt" variants={fadeLeft}>
        <motion.span className="eyebrow" variants={fadeUp}>
          Cenz
        </motion.span>
        <motion.h1 variants={fadeUp}>Creá tu usuario en Cenz y empezá a ordenar tus finanzas.</motion.h1>
        <motion.p variants={fadeUp}>
          Después vas a poder cargar cuentas, registrar ingresos y gastos, y mover saldo entre tus cuentas.
        </motion.p>

        <motion.div className="auth-hero__stats auth-hero__stats--compact" variants={fadeUp}>
          {registerStats.map((item) => (
            <motion.div key={item.title} variants={scaleIn} whileHover={cardHover}>
              <strong>{item.title}</strong>
              <span>{item.description}</span>
            </motion.div>
          ))}
        </motion.div>
      </motion.section>

      <motion.section className="auth-panel auth-panel--form auth-panel--form-register" variants={fadeRight}>
        <motion.div className="auth-form__header" variants={fadeUp}>
          <span className="eyebrow">Registro</span>
          <h2>Nuevo usuario</h2>
          <p>Completá tus datos para entrar directo a tu panel personal.</p>
        </motion.div>

        <motion.form className="stack-form" onSubmit={handleSubmit} variants={fadeUp}>
          <motion.label className="field" variants={fadeUp}>
            <span>Email</span>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="nombre@mail.com"
              required
            />
          </motion.label>

          <motion.label className="field" variants={fadeUp}>
            <span>Usuario</span>
            <input
              type="text"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              placeholder="tuusuario"
              required
            />
          </motion.label>

          <motion.label className="field" variants={fadeUp}>
            <span>Contraseña</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="mínimo 8 caracteres"
              required
            />
          </motion.label>

          <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>

          <motion.div variants={fadeUp} whileHover={buttonHover} whileTap={{ scale: 0.995 }}>
            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting ? "Creando cuenta..." : "Registrarme"}
            </button>
          </motion.div>
        </motion.form>

        <motion.p className="auth-switch" variants={fadeUp}>
          ¿Ya tenés cuenta? <Link to="/login">Iniciar sesión</Link>
        </motion.p>
      </motion.section>
    </StaggerGroup>
  );
}
