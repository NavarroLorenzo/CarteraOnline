import { motion } from "framer-motion";
import type { DashboardCategorySummary } from "../../types/api";
import { formatCurrency } from "../../lib/format";
import { cardHover, fadeUp } from "../ui/animation";

type CategoryDistributionChartProps = {
  categories: DashboardCategorySummary[];
  activeCategoryKey?: string | null;
  onCategorySelect: (categoryKey: string) => void;
};

const categoryColors = ["#f06543", "#19865c", "#0f766e", "#d4a72c", "#7f5539"];

export function CategoryDistributionChart({
  categories,
  activeCategoryKey,
  onCategorySelect,
}: CategoryDistributionChartProps) {
  if (categories.length === 0) {
    return <div className="chart-empty">Todavía no hay gastos categorizados para este rango.</div>;
  }

  return (
    <div className="category-chart">
      <div className="category-chart__bar" aria-hidden="true">
        {categories.map((category, index) => (
          <span
            key={category.key}
            style={{
              width: `${Math.max(category.percentage, 6)}%`,
              backgroundColor: categoryColors[index % categoryColors.length],
            }}
          />
        ))}
      </div>

      <div className="category-chart__list">
        {categories.map((category, index) => (
          <motion.button
            key={category.key}
            type="button"
            className={`category-chart__item ${
              activeCategoryKey === category.key ? "category-chart__item--active" : ""
            }`}
            variants={fadeUp}
            whileHover={cardHover}
            onClick={() => onCategorySelect(category.key)}
          >
            <div className="category-chart__meta">
              <i
                className="category-chart__swatch"
                style={{ backgroundColor: categoryColors[index % categoryColors.length] }}
              />
              <div>
                <strong>{category.label}</strong>
                <span>{category.transactions_count} movimientos</span>
              </div>
            </div>

            <div className="category-chart__values">
              <strong>{category.percentage.toFixed(1)}%</strong>
              <span>{formatCurrency(category.amount)}</span>
            </div>
          </motion.button>
        ))}
      </div>
    </div>
  );
}
