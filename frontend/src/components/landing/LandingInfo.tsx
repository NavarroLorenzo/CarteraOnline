import { motion } from "framer-motion";
import { cardHover, fadeLeft, fadeRight, Reveal } from "./motion";

const points = [
  {
    title: "Menos fricción al registrar",
    description:
      "Cargá ingresos, gastos, transferencias y anotaciones rápidamente, sin perder contexto entre cuentas y movimientos.",
  },
  {
    title: "Mejor lectura del historial",
    description:
      "Filtrá por cuenta, categoría o fechas para responder preguntas concretas sin buscar entre datos sueltos.",
  },
  {
    title: "Más confianza en los números",
    description:
      "Cuando cada movimiento queda ordenado y trazable, el balance deja de ser una intuición y pasa a ser una referencia real.",
  },
];

export function LandingInfo() {
  return (
    <section id="como-funciona" className="landing-section landing-section--split">
      <Reveal className="landing-info__content" variants={fadeLeft}>
        <span className="eyebrow">Por qué sirve</span>
        <h2>Una plataforma pensada para transformar movimientos aislados en información útil</h2>
        <p>
          Tener visibilidad sobre tus finanzas no es solo saber cuánto dinero te queda. También es entender de dónde viene, cómo se mueve y en qué se va.
        </p>

        <div className="landing-info__list">
          {points.map((point) => (
            <motion.article
              key={point.title}
              className="landing-info__item"
              whileHover={cardHover}
            >
              <h3>{point.title}</h3>
              <p>{point.description}</p>
            </motion.article>
          ))}
        </div>
      </Reveal>

      <Reveal className="landing-story" variants={fadeRight}>
        <div className="landing-story__panel">
          <span className="landing-story__eyebrow">Flujo recomendado</span>
          <ol className="landing-story__steps">
            <li>
              <strong>1. Creá tu billetera</strong>
              <p>Separá efectivo, bancos, cuentas digitales y objetivos de ahorro.</p>
            </li>
            <li>
              <strong>2. Registrá la actividad</strong>
              <p>Sumá ingresos, gastos, movimientos internos y notas a medida que ocurren.</p>
            </li>
            <li>
              <strong>3. Revisá el panorama</strong>
              <p>Usá el dashboard para detectar hábitos y tomar decisiones con contexto.</p>
            </li>
          </ol>
        </div>
      </Reveal>
    </section>
  );
}
