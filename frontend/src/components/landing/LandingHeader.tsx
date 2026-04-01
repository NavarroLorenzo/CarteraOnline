import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { buttonHover, delayedTransition, fadeUp } from "./motion";

export function LandingHeader() {
  return (
    <motion.header
      className="landing-nav"
      initial="hidden"
      animate="visible"
      variants={fadeUp}
      transition={delayedTransition(0.05)}
    >
      <Link to="/" className="landing-brand" aria-label="Cartera Online">
        <span className="landing-brand__mark">CO</span>
        <span className="landing-brand__copy">
          <strong>Cartera Online</strong>
          <small>Finanzas personales sin fricción</small>
        </span>
      </Link>

      <nav className="landing-nav__links" aria-label="Navegación principal">
        <a href="#funcionalidades">Funcionalidades</a>
        <a href="#como-funciona">Cómo funciona</a>
        <a href="#cta-final">Empezar</a>
      </nav>

      <div className="landing-nav__actions">
        <Link to="/login" className="ghost-button landing-nav__ghost">
          Ingresar
        </Link>
        <motion.div whileHover={buttonHover} whileTap={{ scale: 0.99 }}>
          <Link to="/register" className="primary-button landing-nav__cta">
            Crear cuenta
          </Link>
        </motion.div>
      </div>
    </motion.header>
  );
}
