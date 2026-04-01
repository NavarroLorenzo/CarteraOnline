import { motion } from "framer-motion";
import { cardHover, fadeUp } from "./animation";

type MetricCardProps = {
  label: string;
  value: string;
  tone?: "default" | "positive" | "negative" | "accent";
};

export function MetricCard({ label, value, tone = "default" }: MetricCardProps) {
  return (
    <motion.article
      className={`metric-card metric-card--${tone}`}
      variants={fadeUp}
      whileHover={cardHover}
    >
      <span>{label}</span>
      <strong>{value}</strong>
    </motion.article>
  );
}
