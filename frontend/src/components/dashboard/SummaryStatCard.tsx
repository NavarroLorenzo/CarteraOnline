import { motion } from "framer-motion";
import { cardHover, fadeUp } from "../ui/animation";

type SummaryStatCardProps = {
  label: string;
  value: string;
  tone?: "default" | "positive" | "negative" | "accent";
  subtitle?: string;
};

export function SummaryStatCard({
  label,
  value,
  tone = "default",
  subtitle,
}: SummaryStatCardProps) {
  return (
    <motion.article
      className={`summary-stat-card summary-stat-card--${tone}`}
      variants={fadeUp}
      whileHover={cardHover}
    >
      <span className="summary-stat-card__label">{label}</span>
      <strong className="summary-stat-card__value">{value}</strong>
      {subtitle ? <small className="summary-stat-card__subtitle">{subtitle}</small> : null}
    </motion.article>
  );
}
