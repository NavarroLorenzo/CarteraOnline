import { motion } from "framer-motion";
import { useMemo, useState } from "react";
import { formatCurrency } from "../../lib/format";
import type { DashboardCategorySummary } from "../../types/api";

type CategoryPieChartProps = {
  categories: DashboardCategorySummary[];
  activeCategoryKey?: string | null;
  onCategorySelect?: (categoryKey: string) => void;
  showList?: boolean;
};

type HoveredSlice = {
  category: DashboardCategorySummary;
  x: number;
  y: number;
};

const SIZE = 280;
const CENTER = SIZE / 2;
const BASE_RADIUS = 94;

const sliceColors = [
  "#f06543",
  "#19865c",
  "#0f766e",
  "#d4a72c",
  "#7f5539",
  "#4062bb",
  "#2d6a4f",
  "#c44536",
  "#3d405b",
  "#a68a64",
  "#ba3f1d",
  "#118ab2",
  "#ef476f",
  "#8f5a3c",
  "#4f772d",
  "#495057",
  "#bc6c25",
  "#6c757d",
];

export function CategoryPieChart({
  categories,
  activeCategoryKey,
  onCategorySelect,
  showList = true,
}: CategoryPieChartProps) {
  const [hoveredSlice, setHoveredSlice] = useState<HoveredSlice | null>(null);

  const chartData = useMemo(() => {
    const totalAmount = categories.reduce((accumulator, category) => accumulator + category.amount, 0);
    if (totalAmount <= 0) {
      return null;
    }

    let currentAngle = -Math.PI / 2;

    return categories.map((category, index) => {
      const ratio = category.amount / totalAmount;
      const angle = ratio * Math.PI * 2;
      const startAngle = currentAngle;
      const endAngle = currentAngle + angle;
      const midAngle = startAngle + angle / 2;
      const isActive = activeCategoryKey === category.key || hoveredSlice?.category.key === category.key;
      const radius = isActive ? BASE_RADIUS + 8 : BASE_RADIUS;
      const offset = isActive ? 9 : 0;
      const centerX = CENTER + Math.cos(midAngle) * offset;
      const centerY = CENTER + Math.sin(midAngle) * offset;
      const centroidX = centerX + Math.cos(midAngle) * (radius * 0.55);
      const centroidY = centerY + Math.sin(midAngle) * (radius * 0.55);

      currentAngle = endAngle;

      return {
        category,
        color: sliceColors[index % sliceColors.length],
        path: describeSlice(centerX, centerY, radius, startAngle, endAngle),
        centroidX,
        centroidY,
        isWhole: ratio >= 0.999,
        centerX,
        centerY,
        radius,
      };
    });
  }, [activeCategoryKey, categories, hoveredSlice?.category.key]);

  if (!chartData) {
    return <div className="chart-empty">No hay gastos categorizados para este período.</div>;
  }

  return (
    <div className={`pie-chart ${showList ? "" : "pie-chart--compact"}`}>
      <div className="pie-chart__canvas">
        <svg className="pie-chart__svg" viewBox={`0 0 ${SIZE} ${SIZE}`} role="img">
          {chartData.map((slice) =>
            slice.isWhole ? (
              <motion.circle
                key={slice.category.key}
                cx={slice.centerX}
                cy={slice.centerY}
                r={slice.radius}
                fill={slice.color}
                className="pie-chart__slice"
                initial={{ scale: 0.92, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
                onMouseEnter={() =>
                  setHoveredSlice({
                    category: slice.category,
                    x: slice.centroidX,
                    y: slice.centroidY,
                  })
                }
                onMouseLeave={() => setHoveredSlice(null)}
                onClick={() => onCategorySelect?.(slice.category.key)}
              />
            ) : (
              <motion.path
                key={slice.category.key}
                d={slice.path}
                fill={slice.color}
                className="pie-chart__slice"
                initial={{ scale: 0.92, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ duration: 0.45, ease: [0.22, 1, 0.36, 1] }}
                onMouseEnter={() =>
                  setHoveredSlice({
                    category: slice.category,
                    x: slice.centroidX,
                    y: slice.centroidY,
                  })
                }
                onMouseLeave={() => setHoveredSlice(null)}
                onClick={() => onCategorySelect?.(slice.category.key)}
              />
            ),
          )}
        </svg>

        {hoveredSlice ? (
          <div
            className="pie-chart__tooltip"
            style={{
              left: `${(hoveredSlice.x / SIZE) * 100}%`,
              top: `${(hoveredSlice.y / SIZE) * 100}%`,
            }}
          >
            <strong>{hoveredSlice.category.label}</strong>
            <span>{formatCurrency(hoveredSlice.category.amount)}</span>
            <span>{hoveredSlice.category.percentage.toFixed(1)}%</span>
          </div>
        ) : null}
      </div>

      {showList ? (
        <div className="pie-chart__list-scroll scrollable-list" aria-label="Lista de categorías">
          <div className="pie-chart__list">
            {chartData.map((slice) => (
              <button
                key={slice.category.key}
                type="button"
                className={`pie-chart__list-item ${
                  activeCategoryKey === slice.category.key ? "pie-chart__list-item--active" : ""
                }`}
                onClick={() => onCategorySelect?.(slice.category.key)}
              >
                <div className="pie-chart__list-meta">
                  <i className="pie-chart__swatch" style={{ backgroundColor: slice.color }} />
                  <div>
                    <strong>{slice.category.label}</strong>
                    <span>{slice.category.transactions_count} movimientos</span>
                  </div>
                </div>

                <div className="pie-chart__list-values">
                  <strong>{formatCurrency(slice.category.amount)}</strong>
                  <span>{slice.category.percentage.toFixed(1)}%</span>
                </div>
              </button>
            ))}
          </div>
        </div>
      ) : null}
    </div>
  );
}

function describeSlice(
  centerX: number,
  centerY: number,
  radius: number,
  startAngle: number,
  endAngle: number,
): string {
  const startX = centerX + Math.cos(startAngle) * radius;
  const startY = centerY + Math.sin(startAngle) * radius;
  const endX = centerX + Math.cos(endAngle) * radius;
  const endY = centerY + Math.sin(endAngle) * radius;
  const largeArcFlag = endAngle - startAngle > Math.PI ? 1 : 0;

  return [
    `M ${centerX} ${centerY}`,
    `L ${startX} ${startY}`,
    `A ${radius} ${radius} 0 ${largeArcFlag} 1 ${endX} ${endY}`,
    "Z",
  ].join(" ");
}
