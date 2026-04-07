import { motion } from "framer-motion";
import { useMemo } from "react";
import type { DashboardTrendPoint } from "../../types/api";

type TrendLineChartProps = {
  points: DashboardTrendPoint[];
  showIncome?: boolean;
  showExpense?: boolean;
};

const WIDTH = 640;
const HEIGHT = 260;
const PADDING_X = 26;
const PADDING_Y = 24;

export function TrendLineChart({
  points,
  showIncome = true,
  showExpense = true,
}: TrendLineChartProps) {
  const chartData = useMemo(() => {
    const innerWidth = WIDTH - PADDING_X * 2;
    const innerHeight = HEIGHT - PADDING_Y * 2;
    const maxValue = Math.max(
      ...points.flatMap((point) => [
        showIncome ? point.income_total : 0,
        showExpense ? point.expense_total : 0,
      ]),
      1,
    );

    const toY = (value: number) => HEIGHT - PADDING_Y - (value / maxValue) * innerHeight;
    const toX = (index: number) =>
      points.length === 1 ? WIDTH / 2 : PADDING_X + (index / (points.length - 1)) * innerWidth;

    const incomePoints = points.map((point, index) => ({
      x: toX(index),
      y: toY(point.income_total),
      label: point.label,
      value: point.income_total,
    }));

    const expensePoints = points.map((point, index) => ({
      x: toX(index),
      y: toY(point.expense_total),
      label: point.label,
      value: point.expense_total,
    }));

    const gridLines = Array.from({ length: 4 }, (_, index) => {
      const ratio = index / 3;
      return {
        y: PADDING_Y + ratio * innerHeight,
        value: Math.round(maxValue - ratio * maxValue),
      };
    });

    return {
      maxValue,
      incomePath: createLinePath(incomePoints),
      expensePath: createLinePath(expensePoints),
      incomePoints,
      expensePoints,
      gridLines,
    };
  }, [points, showExpense, showIncome]);

  if (points.length === 0) {
    return <div className="chart-empty">Todavía no hay datos para graficar en este período.</div>;
  }

  return (
    <div className="trend-chart">
      <div className="trend-chart__legend">
        {showIncome ? (
          <span className="trend-chart__legend-item">
            <i className="trend-chart__dot trend-chart__dot--income" />
            Ingresos
          </span>
        ) : null}
        {showExpense ? (
          <span className="trend-chart__legend-item">
            <i className="trend-chart__dot trend-chart__dot--expense" />
            Gastos
          </span>
        ) : null}
      </div>

      <svg className="trend-chart__svg" viewBox={`0 0 ${WIDTH} ${HEIGHT}`} role="img">
        {chartData.gridLines.map((line) => (
          <g key={line.y}>
            <line
              x1={PADDING_X}
              y1={line.y}
              x2={WIDTH - PADDING_X}
              y2={line.y}
              className="trend-chart__grid"
            />
          </g>
        ))}

        {showIncome ? (
          <motion.path
            d={chartData.incomePath}
            className="trend-chart__line trend-chart__line--income"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 1 }}
            transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1] }}
          />
        ) : null}

        {showExpense ? (
          <motion.path
            d={chartData.expensePath}
            className="trend-chart__line trend-chart__line--expense"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 1 }}
            transition={{ duration: 0.8, ease: [0.22, 1, 0.36, 1], delay: 0.1 }}
          />
        ) : null}

        {showIncome
          ? chartData.incomePoints.map((point) => (
              <motion.circle
                key={`income-${point.label}`}
                cx={point.x}
                cy={point.y}
                r="4"
                className="trend-chart__point trend-chart__point--income"
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ duration: 0.28, delay: 0.2 }}
              />
            ))
          : null}

        {showExpense
          ? chartData.expensePoints.map((point) => (
              <motion.circle
                key={`expense-${point.label}`}
                cx={point.x}
                cy={point.y}
                r="4"
                className="trend-chart__point trend-chart__point--expense"
                initial={{ scale: 0, opacity: 0 }}
                animate={{ scale: 1, opacity: 1 }}
                transition={{ duration: 0.28, delay: 0.25 }}
              />
            ))
          : null}
      </svg>

      <div className="trend-chart__labels">
        {points.map((point) => (
          <span key={point.label}>{point.label}</span>
        ))}
      </div>
    </div>
  );
}

type Point = {
  x: number;
  y: number;
};

function createLinePath(points: Point[]): string {
  if (points.length === 0) {
    return "";
  }

  if (points.length === 1) {
    return `M ${points[0].x} ${points[0].y}`;
  }

  return points.reduce((path, point, index) => {
    if (index === 0) {
      return `M ${point.x} ${point.y}`;
    }

    const previous = points[index - 1];
    const controlX = (previous.x + point.x) / 2;
    return `${path} C ${controlX} ${previous.y}, ${controlX} ${point.y}, ${point.x} ${point.y}`;
  }, "");
}
