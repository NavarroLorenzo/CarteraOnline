import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { buttonHover, fadeUp, Reveal } from "./motion";

export function LandingCTA() {
  return (
    <Reveal id="cta-final" className="landing-cta" variants={fadeUp}>
      <div>
        <span className="eyebrow">Listo para empezar</span>
        <h2>Entrá a tu espacio y empezá a ver tus finanzas con más claridad.</h2>
        <p>
          La app ya está preparada para que cargues tus cuentas, registres movimientos y sigas cada transferencia desde el primer día.
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
