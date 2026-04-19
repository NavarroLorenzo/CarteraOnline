import { motion } from "framer-motion";
import { useMemo, useState } from "react";
import { formatCurrency } from "../../lib/format";
import type { DashboardTrendPoint } from "../../types/api";

type IncomeExpenseAreaChartProps = {
  points: DashboardTrendPoint[];
};

type HoveredGroup = {
  label: string;
  income: number;
  expense: number;
  balance: number;
  isLoss: boolean;
  tooltipPlacement: TooltipPlacement;
  x: number;
  y: number;
};

type SeriesKey = "income" | "expense";
type TooltipPlacement = "left" | "center" | "right";

type ChartPoint = {
  x: number;
  y: number;
};

const WIDTH = 640;
const HEIGHT = 292;
const MARGIN_TOP = 22;
const MARGIN_RIGHT = 20;
const MARGIN_BOTTOM = 48;
const MARGIN_LEFT = 64;
const TICK_COUNT = 5;

export function IncomeExpenseAreaChart({ points }: IncomeExpenseAreaChartProps) {
  const [hoveredGroup, setHoveredGroup] = useState<HoveredGroup | null>(null);
  const [activeSeries, setActiveSeries] = useState<SeriesKey | null>(null);

  const chartData = useMemo(() => {
    if (points.length === 0) {
      return null;
    }

    const innerWidth = WIDTH - MARGIN_LEFT - MARGIN_RIGHT;
    const innerHeight = HEIGHT - MARGIN_TOP - MARGIN_BOTTOM;
    const rawMaxValue = Math.max(
      ...points.flatMap((point) => [point.income_total, point.expense_total]),
      1,
    );
    const maxValue = getNiceMax(rawMaxValue);
    const baselineY = MARGIN_TOP + innerHeight;
    const stepWidth = points.length > 1 ? innerWidth / (points.length - 1) : 0;

    const toY = (value: number) => baselineY - (value / maxValue) * innerHeight;
    const toX = (index: number) => MARGIN_LEFT + (points.length > 1 ? index * stepWidth : innerWidth / 2);

    const ticks = Array.from({ length: TICK_COUNT }, (_, index) => {
      const ratio = index / (TICK_COUNT - 1);
      const value = Math.round(maxValue - ratio * maxValue);
      const y = MARGIN_TOP + ratio * innerHeight;

      return { value, y };
    });

    const groups = points.map((point, index) => {
      const centerX = toX(index);
      const previousX = index > 0 ? toX(index - 1) : MARGIN_LEFT;
      const nextX = index < points.length - 1 ? toX(index + 1) : WIDTH - MARGIN_RIGHT;
      const areaX = index === 0 ? MARGIN_LEFT : previousX + (centerX - previousX) / 2;
      const areaEndX = index === points.length - 1 ? WIDTH - MARGIN_RIGHT : centerX + (nextX - centerX) / 2;
      const incomeY = toY(point.income_total);
      const expenseY = toY(point.expense_total);
      const balance = point.income_total - point.expense_total;

      return {
        point,
        balance,
        isLoss: balance < 0,
        tooltipPlacement: getTooltipPlacement(centerX),
        showLabel: shouldShowXLabel(index, points.length),
        hoverX: centerX,
        hoverY: Math.max(Math.min(incomeY, expenseY, baselineY - 34), MARGIN_TOP + 54),
        areaX,
        areaWidth: areaEndX - areaX,
      };
    });
    const incomePoints = points.map((point, index) => ({
      x: toX(index),
      y: toY(point.income_total),
    }));
    const expensePoints = points.map((point, index) => ({
      x: toX(index),
      y: toY(point.expense_total),
    }));

    return {
      innerHeight,
      baselineY,
      ticks,
      groups,
      incomePoints,
      expensePoints,
      incomeLinePath: buildSmoothPath(incomePoints),
      expenseLinePath: buildSmoothPath(expensePoints),
      incomeAreaPath: buildAreaPath(incomePoints, baselineY),
      expenseAreaPath: buildAreaPath(expensePoints, baselineY),
    };
  }, [points]);

  if (!chartData) {
    return <div className="chart-empty">No hay datos para graficar en este período.</div>;
  }

  const hoveredIncomePoint = hoveredGroup ? findPointByX(chartData.incomePoints, hoveredGroup.x) : undefined;
  const hoveredExpensePoint = hoveredGroup ? findPointByX(chartData.expensePoints, hoveredGroup.x) : undefined;

  return (
    <div
      className={`area-chart${hoveredGroup ? " area-chart--hovering" : ""}${
        activeSeries ? ` area-chart--focus-${activeSeries}` : ""
      }${
        hoveredGroup?.isLoss ? " area-chart--loss-hover" : ""
      }`}
    >
      <div className="area-chart__legend">
        <button
          type="button"
          className="area-chart__legend-item"
          aria-label="Resaltar ingresos"
          onMouseEnter={() => setActiveSeries("income")}
          onMouseLeave={() => setActiveSeries(null)}
          onFocus={() => setActiveSeries("income")}
          onBlur={() => setActiveSeries(null)}
        >
          <i className="area-chart__legend-swatch area-chart__legend-swatch--income" aria-hidden="true" />
          Ingresos
        </button>
        <button
          type="button"
          className="area-chart__legend-item"
          aria-label="Resaltar gastos"
          onMouseEnter={() => setActiveSeries("expense")}
          onMouseLeave={() => setActiveSeries(null)}
          onFocus={() => setActiveSeries("expense")}
          onBlur={() => setActiveSeries(null)}
        >
          <i className="area-chart__legend-swatch area-chart__legend-swatch--expense" aria-hidden="true" />
          Gastos
        </button>
      </div>

      <div className="area-chart__canvas">
        <svg
          className="area-chart__svg"
          viewBox={`0 0 ${WIDTH} ${HEIGHT}`}
          role="img"
          aria-label="Tendencia de ingresos y gastos"
          preserveAspectRatio="xMidYMid meet"
        >
          <defs>
            <linearGradient id="income-area-gradient" x1="0" y1={MARGIN_TOP} x2="0" y2={chartData.baselineY} gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="var(--positive)" stopOpacity="0.4" />
              <stop offset="58%" stopColor="var(--positive)" stopOpacity="0.13" />
              <stop offset="100%" stopColor="var(--positive)" stopOpacity="0" />
            </linearGradient>
            <linearGradient id="expense-area-gradient" x1="0" y1={MARGIN_TOP} x2="0" y2={chartData.baselineY} gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="var(--negative)" stopOpacity="0.38" />
              <stop offset="58%" stopColor="var(--negative)" stopOpacity="0.12" />
              <stop offset="100%" stopColor="var(--negative)" stopOpacity="0" />
            </linearGradient>
            <linearGradient id="hover-band-gradient" x1="0" y1={MARGIN_TOP} x2="0" y2={chartData.baselineY} gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#10211d" stopOpacity="0.08" />
              <stop offset="48%" stopColor="#10211d" stopOpacity="0.042" />
              <stop offset="100%" stopColor="#10211d" stopOpacity="0.01" />
            </linearGradient>
            <linearGradient id="loss-hover-band-gradient" x1="0" y1={MARGIN_TOP} x2="0" y2={chartData.baselineY} gradientUnits="userSpaceOnUse">
              <stop offset="0%" stopColor="#ca3d22" stopOpacity="0.14" />
              <stop offset="48%" stopColor="#ca3d22" stopOpacity="0.064" />
              <stop offset="100%" stopColor="#ca3d22" stopOpacity="0.014" />
            </linearGradient>
          </defs>

          {chartData.ticks.map((tick) => (
            <g key={`tick-${tick.y}`}>
              <line
                x1={MARGIN_LEFT}
                y1={tick.y}
                x2={WIDTH - MARGIN_RIGHT}
                y2={tick.y}
                className="area-chart__grid"
              />
              <text
                x={MARGIN_LEFT - 12}
                y={tick.y + 4}
                textAnchor="end"
                className="area-chart__axis-label"
              >
                {formatAxisCurrency(tick.value)}
              </text>
            </g>
          ))}

          <motion.path
            d={chartData.incomeAreaPath}
            className="area-chart__area area-chart__area--income"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.65, ease: [0.22, 1, 0.36, 1] }}
          />
          <motion.path
            d={chartData.expenseAreaPath}
            className="area-chart__area area-chart__area--expense"
            initial={{ opacity: 0 }}
            animate={{ opacity: 1 }}
            transition={{ duration: 0.65, delay: 0.08, ease: [0.22, 1, 0.36, 1] }}
          />

          <motion.path
            d={chartData.incomeLinePath}
            className="area-chart__line-glow area-chart__line-glow--income"
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 0.9, ease: [0.22, 1, 0.36, 1] }}
          />
          <motion.path
            d={chartData.expenseLinePath}
            className="area-chart__line-glow area-chart__line-glow--expense"
            initial={{ pathLength: 0 }}
            animate={{ pathLength: 1 }}
            transition={{ duration: 0.9, delay: 0.06, ease: [0.22, 1, 0.36, 1] }}
          />

          <motion.path
            d={chartData.incomeLinePath}
            className="area-chart__line area-chart__line--income"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 1 }}
            transition={{ duration: 0.9, ease: [0.22, 1, 0.36, 1] }}
          />
          <motion.path
            d={chartData.expenseLinePath}
            className="area-chart__line area-chart__line--expense"
            initial={{ pathLength: 0, opacity: 0 }}
            animate={{ pathLength: 1, opacity: 1 }}
            transition={{ duration: 0.9, delay: 0.06, ease: [0.22, 1, 0.36, 1] }}
          />

          {chartData.groups.map((group, index) => (
            <g key={`${group.point.label}-${index}`}>
              <rect
                x={group.areaX}
                y={MARGIN_TOP}
                width={group.areaWidth}
                height={chartData.innerHeight + MARGIN_BOTTOM}
                fill="transparent"
                className="area-chart__hover-zone"
                onMouseEnter={() =>
                  setHoveredGroup({
                    label: group.point.label,
                    income: group.point.income_total,
                    expense: group.point.expense_total,
                    balance: group.balance,
                    isLoss: group.isLoss,
                    tooltipPlacement: group.tooltipPlacement,
                    x: group.hoverX,
                    y: group.hoverY,
                  })
                }
                onMouseLeave={() => setHoveredGroup(null)}
              />

              {group.showLabel ? (
                <text
                  x={group.hoverX}
                  y={HEIGHT - 20}
                  textAnchor="middle"
                  className="area-chart__x-label"
                >
                  {group.point.label}
                </text>
              ) : null}
            </g>
          ))}

          {hoveredGroup ? (
            <motion.g
              pointerEvents="none"
              initial={{ opacity: 0 }}
              animate={{ opacity: 1 }}
              transition={{ duration: 0.16, ease: [0.22, 1, 0.36, 1] }}
            >
              <motion.rect
                x={hoveredGroup.x - 15}
                y={MARGIN_TOP}
                width="30"
                height={chartData.baselineY - MARGIN_TOP}
                rx="14"
                className="area-chart__hover-band"
                initial={{ opacity: 0 }}
                animate={{ opacity: 1 }}
                transition={{ duration: 0.18, ease: [0.22, 1, 0.36, 1] }}
              />
              <motion.line
                x1={hoveredGroup.x}
                y1={MARGIN_TOP}
                x2={hoveredGroup.x}
                y2={chartData.baselineY}
                className="area-chart__hover-line"
                initial={{ pathLength: 0, opacity: 0 }}
                animate={{ pathLength: 1, opacity: 1 }}
                transition={{ duration: 0.2, ease: [0.22, 1, 0.36, 1] }}
              />
              {hoveredIncomePoint ? (
                <motion.circle
                  cx={hoveredGroup.x}
                  cy={hoveredIncomePoint.y}
                  r="4.8"
                  className="area-chart__point area-chart__point--income"
                  initial={{ opacity: 0, r: 3.2 }}
                  animate={{ opacity: 1, r: 4.8 }}
                  transition={{ duration: 0.18, ease: [0.22, 1, 0.36, 1] }}
                />
              ) : null}
              {hoveredExpensePoint ? (
                <motion.circle
                  cx={hoveredGroup.x}
                  cy={hoveredExpensePoint.y}
                  r="4.8"
                  className="area-chart__point area-chart__point--expense"
                  initial={{ opacity: 0, r: 3.2 }}
                  animate={{ opacity: 1, r: 4.8 }}
                  transition={{ duration: 0.18, delay: 0.02, ease: [0.22, 1, 0.36, 1] }}
                />
              ) : null}
            </motion.g>
          ) : null}
        </svg>

        {hoveredGroup ? (
          <div
            className={`area-chart__tooltip area-chart__tooltip--${hoveredGroup.tooltipPlacement}`}
            style={{
              left: `${(hoveredGroup.x / WIDTH) * 100}%`,
              top: `${(hoveredGroup.y / HEIGHT) * 100}%`,
            }}
          >
            <strong>{hoveredGroup.label}</strong>
            <span className="area-chart__tooltip-row area-chart__tooltip-row--income">
              <i aria-hidden="true" />
              Ingresos
              <b>{formatCurrency(hoveredGroup.income)}</b>
            </span>
            <span className="area-chart__tooltip-row area-chart__tooltip-row--expense">
              <i aria-hidden="true" />
              Gastos
              <b>{formatCurrency(hoveredGroup.expense)}</b>
            </span>
            <span
              className={`area-chart__tooltip-row area-chart__tooltip-row--balance${
                hoveredGroup.isLoss ? " area-chart__tooltip-row--loss" : ""
              }`}
            >
              <i aria-hidden="true" />
              Balance
              <b>{formatCurrency(hoveredGroup.balance)}</b>
            </span>
          </div>
        ) : null}
      </div>
    </div>
  );
}

function getNiceMax(value: number): number {
  const paddedValue = value * 1.12;
  const magnitude = 10 ** Math.floor(Math.log10(paddedValue));
  const normalized = paddedValue / magnitude;
  const niceNormalized = normalized <= 1 ? 1 : normalized <= 2 ? 2 : normalized <= 5 ? 5 : 10;

  return niceNormalized * magnitude;
}

function formatAxisCurrency(value: number): string {
  const absoluteValue = Math.abs(value);

  if (absoluteValue >= 1_000_000) {
    return `$${formatCompactNumber(value / 1_000_000)}M`;
  }

  if (absoluteValue >= 1_000) {
    return `$${formatCompactNumber(value / 1_000)}k`;
  }

  return `$${Math.round(value)}`;
}

function formatCompactNumber(value: number): string {
  const roundedValue = value >= 10 ? Math.round(value).toString() : value.toFixed(1);
  return roundedValue.endsWith(".0") ? roundedValue.slice(0, -2) : roundedValue;
}

function buildSmoothPath(points: ChartPoint[]): string {
  if (points.length === 0) {
    return "";
  }

  if (points.length === 1) {
    return `M ${points[0].x} ${points[0].y}`;
  }

  const commands = [`M ${points[0].x} ${points[0].y}`];

  for (let index = 0; index < points.length - 1; index += 1) {
    const previous = points[index - 1] ?? points[index];
    const current = points[index];
    const next = points[index + 1];
    const afterNext = points[index + 2] ?? next;

    const controlOneX = current.x + (next.x - previous.x) / 6;
    const controlOneY = current.y + (next.y - previous.y) / 6;
    const controlTwoX = next.x - (afterNext.x - current.x) / 6;
    const controlTwoY = next.y - (afterNext.y - current.y) / 6;

    commands.push(`C ${controlOneX} ${controlOneY}, ${controlTwoX} ${controlTwoY}, ${next.x} ${next.y}`);
  }

  return commands.join(" ");
}

function buildAreaPath(points: ChartPoint[], baselineY: number): string {
  if (points.length === 0) {
    return "";
  }

  if (points.length === 1) {
    const point = points[0];
    const halfWidth = 12;
    return `M ${point.x - halfWidth} ${baselineY} L ${point.x - halfWidth} ${point.y} L ${point.x + halfWidth} ${point.y} L ${point.x + halfWidth} ${baselineY} Z`;
  }

  const linePath = buildSmoothPath(points);
  const firstPoint = points[0];
  const lastPoint = points[points.length - 1];

  return `${linePath} L ${lastPoint.x} ${baselineY} L ${firstPoint.x} ${baselineY} Z`;
}

function findPointByX(points: ChartPoint[], x: number): ChartPoint | undefined {
  return points.find((point) => point.x === x);
}

function getTooltipPlacement(x: number): TooltipPlacement {
  if (x < WIDTH * 0.24) {
    return "left";
  }

  if (x > WIDTH * 0.76) {
    return "right";
  }

  return "center";
}

function shouldShowXLabel(index: number, total: number): boolean {
  if (total <= 7) {
    return true;
  }

  const step = total <= 14 ? 2 : total <= 31 ? 4 : Math.ceil(total / 8);
  return index === 0 || index === total - 1 || index % step === 0;
}
