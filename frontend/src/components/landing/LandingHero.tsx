import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { HeroVisual } from "./HeroVisual";
import { buttonHover, delayedTransition, fadeLeft, fadeRight, StaggerGroup, fadeUp } from "./motion";
import { scrollToSection } from "./scroll";

const highlights = [
  "Billetera virtual y cuentas en una vista",
  "Gastos, ingresos y ahorro con contexto",
  "Notas financieras para no perder detalles",
];

export function LandingHero() {
  return (
    <section className="landing-hero">
      <StaggerGroup className="landing-hero__content">
        <motion.div className="landing-hero__copy" variants={fadeLeft}>
          <motion.span className="eyebrow" variants={fadeUp}>
            Cenz wallet personal
          </motion.span>
          <motion.h1 variants={fadeUp}>
            Tu billetera virtual para ordenar gastos, ahorro y notas.
          </motion.h1>
          <motion.p variants={fadeUp}>
            Cenz reúne tus cuentas, movimientos y anotaciones financieras para que sepas cuánto tenés, en qué se va tu dinero y qué querés mejorar cada mes.
          </motion.p>

          <motion.div className="landing-hero__actions" variants={fadeUp}>
            <motion.div whileHover={buttonHover} whileTap={{ scale: 0.99 }}>
              <Link to="/register" className="primary-button">
                Comenzar
              </Link>
            </motion.div>
            <Link to="/login" className="ghost-button landing-hero__secondary">
              Ingresar
            </Link>
            <button
              className="ghost-button landing-hero__secondary landing-hero__learn"
              type="button"
              onClick={() => scrollToSection("como-funciona")}
            >
              Ver cómo funciona
            </button>
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
