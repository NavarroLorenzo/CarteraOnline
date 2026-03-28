import { useState } from "react";
import { Link, useNavigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";

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
    <div className="auth-shell">
      <section className="auth-panel auth-panel--hero auth-panel--hero-alt">
        <span className="eyebrow">Onboarding</span>
        <h1>Creá un acceso propio antes de cargar cuentas y movimientos.</h1>
        <p>
          Más adelante podés sumar verificación de email o login con Gmail sin rehacer la base del
          frontend.
        </p>
      </section>

      <section className="auth-panel auth-panel--form">
        <div className="panel-heading">
          <span className="eyebrow">Registro</span>
          <h2>Nuevo usuario</h2>
          <p>Guardamos el token y te redirigimos a tu panel al terminar.</p>
        </div>

        <form className="stack-form" onSubmit={handleSubmit}>
          <label className="field">
            <span>Email</span>
            <input
              type="email"
              value={email}
              onChange={(event) => setEmail(event.target.value)}
              placeholder="nombre@mail.com"
              required
            />
          </label>

          <label className="field">
            <span>Usuario</span>
            <input
              type="text"
              value={username}
              onChange={(event) => setUsername(event.target.value)}
              placeholder="tuusuario"
              required
            />
          </label>

          <label className="field">
            <span>Contraseña</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="mínimo 6 caracteres"
              required
            />
          </label>

          {error ? <p className="feedback feedback--error">{error}</p> : null}

          <button type="submit" className="primary-button" disabled={submitting}>
            {submitting ? "Creando cuenta..." : "Registrarme"}
          </button>
        </form>

        <p className="auth-switch">
          ¿Ya tenés cuenta? <Link to="/login">Iniciar sesión</Link>
        </p>
      </section>
    </div>
  );
}
