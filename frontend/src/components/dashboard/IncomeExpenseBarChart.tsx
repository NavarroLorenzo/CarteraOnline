import { motion } from "framer-motion";
import { useMemo, useState } from "react";
import { formatCurrency } from "../../lib/format";
import type { DashboardTrendPoint } from "../../types/api";

type IncomeExpenseBarChartProps = {
  points: DashboardTrendPoint[];
};

type HoveredGroup = {
  label: string;
  income: number;
  expense: number;
  x: number;
  y: number;
};

const WIDTH = 760;
const HEIGHT = 340;
const MARGIN_TOP = 24;
const MARGIN_RIGHT = 24;
const MARGIN_BOTTOM = 56;
const MARGIN_LEFT = 88;
const TICK_COUNT = 5;

export function IncomeExpenseBarChart({ points }: IncomeExpenseBarChartProps) {
  const [hoveredGroup, setHoveredGroup] = useState<HoveredGroup | null>(null);

  const chartData = useMemo(() => {
    if (points.length === 0) {
      return null;
    }

    const innerWidth = WIDTH - MARGIN_LEFT - MARGIN_RIGHT;
    const innerHeight = HEIGHT - MARGIN_TOP - MARGIN_BOTTOM;
    const maxValue = Math.max(
      ...points.flatMap((point) => [point.income_total, point.expense_total]),
      1,
    );
    const groupWidth = innerWidth / points.length;
    const barWidth = Math.min(28, Math.max(12, groupWidth * 0.28));
    const barGap = Math.max(6, groupWidth * 0.08);
    const baselineY = MARGIN_TOP + innerHeight;

    const toHeight = (value: number) => (value / maxValue) * innerHeight;

    const ticks = Array.from({ length: TICK_COUNT }, (_, index) => {
      const ratio = index / (TICK_COUNT - 1);
      const value = Math.round(maxValue - ratio * maxValue);
      const y = MARGIN_TOP + ratio * innerHeight;

      return { value, y };
    });

    const groups = points.map((point, index) => {
      const groupStartX = MARGIN_LEFT + index * groupWidth;
      const centerX = groupStartX + groupWidth / 2;
      const incomeHeight = toHeight(point.income_total);
      const expenseHeight = toHeight(point.expense_total);
      const incomeY = baselineY - incomeHeight;
      const expenseY = baselineY - expenseHeight;

      return {
        point,
        hoverX: centerX,
        hoverY: Math.min(incomeY, expenseY, baselineY - 28),
        areaX: groupStartX,
        areaWidth: groupWidth,
        income: {
          x: centerX - barWidth - barGap / 2,
          y: incomeY,
          width: barWidth,
          height: incomeHeight,
        },
        expense: {
          x: centerX + barGap / 2,
          y: expenseY,
          width: barWidth,
          height: expenseHeight,
        },
      };
    });

    return {
      innerHeight,
      baselineY,
      ticks,
      groups,
    };
  }, [points]);

  if (!chartData) {
    return <div className="chart-empty">Todavía no hay datos para graficar en este período.</div>;
  }

  return (
    <div className="bar-chart">
      <div className="bar-chart__legend">
        <span className="bar-chart__legend-item">
          <i className="bar-chart__legend-swatch bar-chart__legend-swatch--income" />
          Ingresos
        </span>
        <span className="bar-chart__legend-item">
          <i className="bar-chart__legend-swatch bar-chart__legend-swatch--expense" />
          Gastos
        </span>
      </div>

      <div className="bar-chart__canvas">
        <svg className="bar-chart__svg" viewBox={`0 0 ${WIDTH} ${HEIGHT}`} role="img">
          {chartData.ticks.map((tick) => (
            <g key={`tick-${tick.y}`}>
              <line
                x1={MARGIN_LEFT}
                y1={tick.y}
                x2={WIDTH - MARGIN_RIGHT}
                y2={tick.y}
                className="bar-chart__grid"
              />
              <text
                x={MARGIN_LEFT - 12}
                y={tick.y + 4}
                textAnchor="end"
                className="bar-chart__axis-label"
              >
                {formatCurrency(tick.value)}
              </text>
            </g>
          ))}

          {chartData.groups.map((group) => (
            <g key={group.point.label}>
              <motion.rect
                x={group.income.x}
                y={group.income.y}
                width={group.income.width}
                height={group.income.height}
                rx="10"
                className="bar-chart__bar bar-chart__bar--income"
                initial={{ height: 0, y: chartData.baselineY }}
                animate={{ height: group.income.height, y: group.income.y }}
                transition={{ duration: 0.5, ease: [0.22, 1, 0.36, 1] }}
              />
              <motion.rect
                x={group.expense.x}
                y={group.expense.y}
                width={group.expense.width}
                height={group.expense.height}
                rx="10"
                className="bar-chart__bar bar-chart__bar--expense"
                initial={{ height: 0, y: chartData.baselineY }}
                animate={{ height: group.expense.height, y: group.expense.y }}
                transition={{ duration: 0.5, delay: 0.06, ease: [0.22, 1, 0.36, 1] }}
              />

              <rect
                x={group.areaX}
                y={MARGIN_TOP}
                width={group.areaWidth}
                height={chartData.innerHeight + MARGIN_BOTTOM}
                fill="transparent"
                onMouseEnter={() =>
                  setHoveredGroup({
                    label: group.point.label,
                    income: group.point.income_total,
                    expense: group.point.expense_total,
                    x: group.hoverX,
                    y: group.hoverY,
                  })
                }
                onMouseLeave={() => setHoveredGroup(null)}
              />

              <text
                x={group.hoverX}
                y={HEIGHT - 20}
                textAnchor="middle"
                className="bar-chart__x-label"
              >
                {group.point.label}
              </text>
            </g>
          ))}
        </svg>

        {hoveredGroup ? (
          <div
            className="bar-chart__tooltip"
            style={{
              left: `${(hoveredGroup.x / WIDTH) * 100}%`,
              top: `${(hoveredGroup.y / HEIGHT) * 100}%`,
            }}
          >
            <strong>{hoveredGroup.label}</strong>
            <span>Ingresos: {formatCurrency(hoveredGroup.income)}</span>
            <span>Gastos: {formatCurrency(hoveredGroup.expense)}</span>
          </div>
        ) : null}
      </div>
    </div>
  );
}
