import { motion, useScroll, useTransform } from "framer-motion";
import { useRef } from "react";
import { cardHover } from "./motion";

export function HeroVisual() {
  const visualRef = useRef<HTMLDivElement | null>(null);
  const { scrollYProgress } = useScroll({
    target: visualRef,
    offset: ["start start", "end start"],
  });

  const topCardY = useTransform(scrollYProgress, [0, 1], [0, -42]);
  const bottomCardY = useTransform(scrollYProgress, [0, 1], [0, 30]);
  const ringY = useTransform(scrollYProgress, [0, 1], [0, -54]);
  const ringRotate = useTransform(scrollYProgress, [0, 1], [0, 8]);

  return (
    <div ref={visualRef} className="hero-visual">
      <motion.div
        className="hero-visual__orb hero-visual__orb--primary"
        style={{ y: ringY, rotate: ringRotate }}
      />
      <motion.div className="hero-visual__orb hero-visual__orb--soft" />

      <motion.article
        className="hero-visual__card hero-visual__card--summary"
        style={{ y: topCardY }}
        whileHover={cardHover}
      >
        <span className="hero-visual__label">Saldo total</span>
        <strong>$ 2.480.000</strong>
        <p>Una vista clara de todas tus cuentas en el mismo panel.</p>

        <div className="hero-visual__meter">
          <span style={{ width: "78%" }} />
        </div>

        <div className="hero-visual__chips">
          <span>Banco</span>
          <span>Billetera</span>
          <span>Efectivo</span>
        </div>
      </motion.article>

      <motion.article
        className="hero-visual__card hero-visual__card--flow"
        style={{ y: bottomCardY }}
        whileHover={cardHover}
      >
        <header>
          <span className="hero-visual__label">Movimientos del mes</span>
          <strong>+12.4%</strong>
        </header>

        <div className="hero-visual__bars" aria-hidden="true">
          <span style={{ height: "28%" }} />
          <span style={{ height: "46%" }} />
          <span style={{ height: "38%" }} />
          <span style={{ height: "61%" }} />
          <span style={{ height: "82%" }} />
          <span style={{ height: "70%" }} />
          <span style={{ height: "94%" }} />
        </div>

        <div className="hero-visual__transfer">
          <div>
            <small>Origen</small>
            <strong>Efectivo</strong>
          </div>
          <span aria-hidden="true">→</span>
          <div>
            <small>Destino</small>
            <strong>Mercado Pago</strong>
          </div>
        </div>
      </motion.article>
    </div>
  );
}
