import { motion } from "framer-motion";
import {
  formatDashboardRangeLabel,
  getDashboardPresetLabel,
  getDashboardSummaryTitle,
  type DashboardDateRange,
  type DashboardFilterPreset,
} from "../../lib/dashboard";
import { formatCurrency } from "../../lib/format";
import type { DashboardSummary as DashboardSummaryData } from "../../types/api";
import { StaggerGroup, cardHover, fadeUp } from "../ui/animation";

type DashboardSummaryProps = {
  preset: DashboardFilterPreset;
  range: DashboardDateRange | null;
  summary: DashboardSummaryData | null | undefined;
  loading: boolean;
};

const statItems = [
  { key: "income", label: "Ingresos", tone: "positive" },
  { key: "expense", label: "Gastos", tone: "negative" },
  { key: "balance", label: "Balance", tone: "accent" },
] as const;

export function DashboardSummary({ preset, range, summary, loading }: DashboardSummaryProps) {
  const title = getDashboardSummaryTitle(preset);
  const rangeLabel = formatDashboardRangeLabel(preset, range);
  const metricFallback = loading ? "Cargando información..." : "Sin datos";

  return (
    <section className="panel dashboard-summary">
      <div className="dashboard-summary__header">
        <div>
          <span className="eyebrow">Resumen clave</span>
          <h2>{title}</h2>
          <p>{rangeLabel}</p>
        </div>

        <div className="dashboard-summary__meta">
          <span className="dashboard-summary__meta-label">Período activo</span>
          <strong>{getDashboardPresetLabel(preset)}</strong>
          <small>{summary ? `${summary.transactions_count} movimientos analizados` : "Esperando datos del período"}</small>
        </div>
      </div>

      <StaggerGroup className="dashboard-summary__stats" onView={false}>
        {statItems.map((item) => {
          const value =
            item.key === "income"
              ? summary
                ? formatCurrency(summary.income_total)
                : metricFallback
              : item.key === "expense"
                ? summary
                  ? formatCurrency(summary.expense_total)
                  : metricFallback
                : summary
                  ? formatCurrency(summary.net_balance)
                  : metricFallback;

          return (
            <motion.article
              key={item.key}
              className={`dashboard-summary__stat dashboard-summary__stat--${item.tone}`}
              variants={fadeUp}
              whileHover={cardHover}
            >
              <span className="dashboard-summary__stat-label">{item.label}</span>
              <strong className="dashboard-summary__stat-value">{value}</strong>
            </motion.article>
          );
        })}
      </StaggerGroup>
    </section>
  );
}
