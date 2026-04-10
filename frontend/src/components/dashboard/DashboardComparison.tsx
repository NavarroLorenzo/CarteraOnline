import { motion } from "framer-motion";
import { formatPercent } from "../../lib/dashboard";
import { formatCurrency } from "../../lib/format";
import type { DashboardComparison as DashboardComparisonData } from "../../types/api";
import { EmptyState } from "../ui/EmptyState";
import { fadeRight, fadeUp } from "../ui/animation";

type DashboardComparisonProps = {
  comparison: DashboardComparisonData | null | undefined;
  loading: boolean;
};

export function DashboardComparison({ comparison, loading }: DashboardComparisonProps) {
  return (
    <motion.article className="panel dashboard-panel dashboard-panel--comparison" variants={fadeRight}>
      <div className="panel-heading">
        <div>
          <span className="eyebrow">Comparación</span>
          <h2>{comparison?.title ?? "Comparación del período"}</h2>
          <p>Compara el período activo con su equivalente anterior para detectar variaciones con el mismo criterio.</p>
        </div>
      </div>

      {loading ? (
        <p className="feedback">Cargando información...</p>
      ) : comparison ? (
        <div className="comparison-grid">
          <ComparisonBlock
            label="Ingresos"
            previousLabel={comparison.previous_label}
            currentValue={comparison.current.income_total}
            previousValue={comparison.previous.income_total}
            change={comparison.income_change_pct}
            variant="income"
          />
          <ComparisonBlock
            label="Gastos"
            previousLabel={comparison.previous_label}
            currentValue={comparison.current.expense_total}
            previousValue={comparison.previous.expense_total}
            change={comparison.expense_change_pct}
            variant="expense"
          />
        </div>
      ) : (
        <EmptyState
          title="Comparación no disponible"
          description="No pudimos construir una comparación clara para este período seleccionado."
        />
      )}
    </motion.article>
  );
}

type ComparisonBlockProps = {
  label: string;
  previousLabel: string;
  currentValue: number;
  previousValue: number;
  change: number;
  variant: "income" | "expense";
};

function ComparisonBlock({ label, previousLabel, currentValue, previousValue, change, variant }: ComparisonBlockProps) {
  const direction = change > 0 ? "up" : change < 0 ? "down" : "flat";
  const isPositive =
    variant === "expense" ? direction === "down" || direction === "flat" : direction === "up" || direction === "flat";
  const symbol = direction === "up" ? "↑" : direction === "down" ? "↓" : "•";

  return (
    <motion.div className="comparison-card" variants={fadeUp}>
      <span className="comparison-card__label">{label}</span>
      <strong className="comparison-card__value">{formatCurrency(currentValue)}</strong>
      <span
        className={`comparison-card__trend ${
          isPositive ? "comparison-card__trend--positive" : "comparison-card__trend--negative"
        }`}
      >
        {symbol} {formatPercent(change)}%
      </span>
      <small>
        {previousLabel}: {formatCurrency(previousValue)}
      </small>
    </motion.div>
  );
}
