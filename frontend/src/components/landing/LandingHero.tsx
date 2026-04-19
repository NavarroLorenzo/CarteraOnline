import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { HeroVisual } from "./HeroVisual";
import { buttonHover, delayedTransition, fadeLeft, fadeRight, StaggerGroup, fadeUp } from "./motion";

const highlights = [
  "Tus cuentas separadas en un solo lugar",
  "Transferencias internas con trazabilidad",
  "Historial ordenado por fecha y categoría",
];

export function LandingHero() {
  return (
    <section className="landing-hero">
      <StaggerGroup className="landing-hero__content">
        <motion.div className="landing-hero__copy" variants={fadeLeft}>
          <motion.span className="eyebrow" variants={fadeUp}>
            Finanzas personales claras
          </motion.span>
          <motion.h1 variants={fadeUp}>
            Ordená tu dinero, entendé tus movimientos y tomá mejores decisiones.
          </motion.h1>
          <motion.p variants={fadeUp}>
            Cenz te ayuda a registrar cuentas, ingresos, gastos y transferencias con una experiencia simple, prolija y pensada para uso diario.
          </motion.p>

          <motion.div className="landing-hero__actions" variants={fadeUp}>
            <motion.div whileHover={buttonHover} whileTap={{ scale: 0.99 }}>
              <Link to="/login" className="primary-button">
                Ingresar
              </Link>
            </motion.div>
            <Link to="/register" className="ghost-button landing-hero__secondary">
              Comenzar ahora
            </Link>
          </motion.div>

          <motion.ul className="landing-hero__highlights" variants={fadeUp}>
            {highlights.map((item, index) => (
              <motion.li
                key={item}
                variants={fadeUp}
                transition={delayedTransition(0.12 + index * 0.08)}
              >
                {item}
              </motion.li>
            ))}
          </motion.ul>
        </motion.div>

        <motion.div className="landing-hero__visual" variants={fadeRight}>
          <HeroVisual />
        </motion.div>
      </StaggerGroup>
    </section>
  );
}
