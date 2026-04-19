import { motion } from "framer-motion";
import { cardHover, fadeUp, Reveal, StaggerGroup } from "./motion";

const features = [
  {
    title: "Cuentas organizadas",
    description:
      "Separá efectivo, bancos, billeteras virtuales y ahorro con estados claros para cada cuenta.",
  },
  {
    title: "Movimientos simples",
    description:
      "Registrá ingresos y gastos con categoría, fecha y detalle para mantener tu historial limpio.",
  },
  {
    title: "Transferencias consistentes",
    description:
      "Mové saldo entre cuentas y mantené trazabilidad entre ambos movimientos relacionados.",
  },
  {
    title: "Panel listo para actuar",
    description:
      "Visualizá balance total, actividad reciente y evolución de tus cuentas en segundos.",
  },
];

export function LandingFeatures() {
  return (
    <section id="funcionalidades" className="landing-section">
      <Reveal className="landing-section__intro">
        <span className="eyebrow">Qué podés hacer</span>
        <h2>Cenz te da una base sólida para llevar tus finanzas sin planillas dispersas</h2>
        <p>
          Todo está pensado para que puedas entender tu dinero con claridad y tomar decisiones más rápido.
        </p>
      </Reveal>

      <StaggerGroup className="landing-feature-grid">
        {features.map((feature) => (
          <motion.article
            key={feature.title}
            className="landing-feature-card"
            variants={fadeUp}
            whileHover={cardHover}
          >
            <span className="landing-feature-card__index" aria-hidden="true" />
            <h3>{feature.title}</h3>
            <p>{feature.description}</p>
          </motion.article>
        ))}
      </StaggerGroup>
    </section>
  );
}
