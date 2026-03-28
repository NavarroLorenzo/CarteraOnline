import { useState } from "react";
import { Link, useLocation, useNavigate } from "react-router-dom";
import { ApiError } from "../api/client";
import { useAuth } from "../auth/AuthContext";

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
    <div className="auth-shell">
      <section className="auth-panel auth-panel--hero">
        <span className="eyebrow">Finanzas en orden</span>
        <h1>Entrá y retomá el control de tu dinero.</h1>
        <p>
          Accedé a tus cuentas, movimientos y transferencias desde un solo lugar, con tu sesión protegida.
        </p>

        <div className="auth-hero__stats">
          <div>
            <strong>Sesión segura</strong>
            <span>Acceso protegido con token</span>
          </div>
          <div>
            <strong>Tu espacio</strong>
            <span>Cada usuario ve solo lo suyo</span>
          </div>
          <div>
            <strong>Todo conectado</strong>
            <span>Cuentas, movimientos y transferencias</span>
          </div>
        </div>
      </section>

      <section className="auth-panel auth-panel--form">
        <div className="panel-heading">
          <span className="eyebrow">Login</span>
          <h2>Bienvenido de nuevo</h2>
          <p>Podés ingresar con tu email o con tu nombre de usuario.</p>
        </div>

        <form className="stack-form" onSubmit={handleSubmit}>
          <label className="field">
            <span>Usuario o email</span>
            <input
              type="text"
              value={identifier}
              onChange={(event) => setIdentifier(event.target.value)}
              placeholder="usuario1 o usuario1@mail.com"
              required
            />
          </label>

          <label className="field">
            <span>Contraseña</span>
            <input
              type="password"
              value={password}
              onChange={(event) => setPassword(event.target.value)}
              placeholder="••••••"
              required
            />
          </label>

          {error ? <p className="feedback feedback--error">{error}</p> : null}

          <button type="submit" className="primary-button" disabled={submitting}>
            {submitting ? "Ingresando..." : "Iniciar sesión"}
          </button>
        </form>

        <p className="auth-switch">
          ¿No tenés cuenta? <Link to="/register">Crear usuario</Link>
        </p>
      </section>
    </div>
  );
}
