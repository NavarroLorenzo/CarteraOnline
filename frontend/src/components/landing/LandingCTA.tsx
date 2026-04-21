import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { buttonHover, fadeUp, Reveal } from "./motion";

export function LandingCTA() {
  return (
    <Reveal id="cta-final" className="landing-cta" variants={fadeUp}>
      <div>
        <span className="eyebrow">Listo para empezar</span>
        <h2>Entrá a Cenz y empezá a ver tus finanzas con más claridad.</h2>
        <p>
          Cenz ya está preparado para que cargues tu billetera, registres movimientos, sigas ahorros y guardes notas financieras desde el primer día.
        </p>
      </div>

      <motion.div whileHover={buttonHover} whileTap={{ scale: 0.99 }}>
        <Link to="/login" className="primary-button landing-cta__button">
          Ir al login
        </Link>
      </motion.div>
    </Reveal>
  );
}
