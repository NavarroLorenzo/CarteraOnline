import { useReducedMotion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";

const transactions = [
  { id: 1, label: "Sueldo acreditado", amount: "+$1.240.000", type: "income", badge: "IN" },
  { id: 2, label: "Supermercado", amount: "-$82.400", type: "expense", badge: "GS" },
  { id: 3, label: "Nota: alquiler mayo", amount: "Anotado", type: "note", badge: "NT" },
  { id: 4, label: "Ahorro mensual", amount: "+$180.000", type: "income", badge: "AH" },
  { id: 5, label: "Servicios", amount: "-$46.900", type: "expense", badge: "SV" },
  { id: 6, label: "Transferencia a billetera", amount: "$75.000", type: "note", badge: "TR" },
] as const;

const chartPoints = [38, 55, 42, 70, 60, 82, 74, 91, 85, 95, 88, 100];
const floatingMarks = ["$", "+", "%", "$", "C", "$"] as const;

type SparklineProps = {
  points: number[];
  progress: number;
};

type BudgetRingProps = {
  percent: number;
  tick: number;
};

type SavingsBarProps = {
  percent: number;
  tick: number;
};

type FlowLineProps = {
  x1: number;
  y1: number;
  x2: number;
  y2: number;
  tick: number;
  speed?: number;
  delay?: number;
};

function lerp(a: number, b: number, t: number) {
  return a + (b - a) * t;
}

function useTick(fps = 30) {
  const prefersReducedMotion = useReducedMotion();
  const [tick, setTick] = useState(0);
  const start = useRef(Date.now());

  useEffect(() => {
    if (prefersReducedMotion) {
      setTick(0);
      return undefined;
    }

    const interval = window.setInterval(() => {
      setTick((Date.now() - start.current) / 1000);
    }, 1000 / fps);

    return () => window.clearInterval(interval);
  }, [fps, prefersReducedMotion]);

  return tick;
}

function Sparkline({ points, progress }: SparklineProps) {
  const width = 210;
  const height = 54;

  const { pathD, areaD, lastPoint } = useMemo(() => {
    const visibleCount = Math.max(1, Math.round(progress * (points.length - 1)));
    const xs = points.map((_, index) => (index / (points.length - 1)) * width);
    const ys = points.map((value) => height - (value / 100) * height * 0.82 - 4);
    const visibleX = xs.slice(0, visibleCount + 1);
    const visibleY = ys.slice(0, visibleCount + 1);
    const path = visibleX.map((x, index) => `${index === 0 ? "M" : "L"}${x},${visibleY[index]}`).join(" ");
    const lastX = visibleX[visibleX.length - 1] ?? 0;

    return {
      pathD: path,
      areaD: `${path} L${lastX},${height} L0,${height} Z`,
      lastPoint: { x: lastX, y: visibleY[visibleY.length - 1] ?? height },
    };
  }, [points, progress]);

  return (
    <svg className="finance-hero-chart" viewBox={`0 0 ${width} ${height}`} aria-hidden="true">
      <defs>
        <linearGradient id="financeHeroChartFill" x1="0" y1="0" x2="0" y2="1">
          <stop offset="0%" stopColor="var(--positive)" stopOpacity="0.26" />
          <stop offset="100%" stopColor="var(--positive)" stopOpacity="0" />
        </linearGradient>
      </defs>
      <path d={areaD} fill="url(#financeHeroChartFill)" />
      <path className="finance-hero-chart__line" d={pathD} />
      <circle className="finance-hero-chart__point" cx={lastPoint.x} cy={lastPoint.y} r="4.5" />
    </svg>
  );
}

function BudgetRing({ percent, tick }: BudgetRingProps) {
  const radius = 27;
  const center = 32;
  const circumference = 2 * Math.PI * radius;
  const visibleStroke = (percent / 100) * circumference;
  const pulse = 1 + Math.sin(tick * 1.8) * 0.025;

  return (
    <svg
      className="finance-hero-ring"
      viewBox="0 0 64 64"
      style={{ transform: `scale(${pulse})` }}
      aria-hidden="true"
    >
      <circle className="finance-hero-ring__track" cx={center} cy={center} r={radius} />
      <circle
        className="finance-hero-ring__value"
        cx={center}
        cy={center}
        r={radius}
        strokeDasharray={`${visibleStroke} ${circumference - visibleStroke}`}
        strokeDashoffset={circumference / 4}
      />
      <circle className="finance-hero-ring__center" cx={center} cy={center} r={radius - 11} />
    </svg>
  );
}

function SavingsBar({ percent, tick }: SavingsBarProps) {
  const glow = 0.5 + Math.sin(tick * 1.4) * 0.5;

  return (
    <div className="finance-hero-savings__track" aria-hidden="true">
      <span
        className="finance-hero-savings__value"
        style={{
          width: `${percent}%`,
          boxShadow: `0 0 ${8 * glow}px rgba(25, 134, 92, 0.38)`,
        }}
      />
    </div>
  );
}

function FlowLine({ x1, y1, x2, y2, tick, speed = 1, delay = 0 }: FlowLineProps) {
  const mx = (x1 + x2) / 2 + (y2 - y1) * 0.22;
  const my = (y1 + y2) / 2 - (x2 - x1) * 0.12;
  const offset = -((tick * speed * 38 + delay * 28) % 55);

  return (
    <path
      className="finance-hero-flow__line"
      d={`M${x1},${y1} Q${mx},${my} ${x2},${y2}`}
      strokeDashoffset={offset}
    />
  );
}

export function HeroVisual() {
  const tick = useTick(30);
  const loop = 6;
  const phase = (tick % loop) / loop;
  const chartProgress = Math.min(phase * 2, 1);
  const balancePulse = 1 + Math.sin(tick * 1.2) * 0.006;
  const transactionIndex = Math.floor(phase * transactions.length) % transactions.length;
  const budgetPercent = Math.round(lerp(52, 78, Math.sin(tick * 0.6) * 0.5 + 0.5));

  const walletY = Math.sin(tick * 0.7) * 7;
  const walletRotate = Math.sin(tick * 0.5) * 1.8;
  const chartY = Math.sin(tick * 0.9 + 1) * 5;
  const chartRotate = Math.sin(tick * 0.6 + 2) * 1.4;
  const budgetY = Math.sin(tick * 0.8 + 3) * 6;
  const budgetRotate = Math.sin(tick * 0.55 + 1) * 2;
  const savingsY = Math.sin(tick * 0.65 + 2) * 5;
  const shimmerLeft = -60 + ((tick * 35) % 220);

  return (
    <div className="finance-hero-visual" aria-label="Vista previa animada de Cenz">
      <svg className="finance-hero-flow" viewBox="0 0 480 480" preserveAspectRatio="none" aria-hidden="true">
        <defs>
          <linearGradient id="financeHeroFlow" x1="0" y1="0" x2="1" y2="1">
            <stop offset="0%" stopColor="var(--positive)" stopOpacity="0" />
            <stop offset="50%" stopColor="var(--positive)" stopOpacity="0.42" />
            <stop offset="100%" stopColor="var(--accent)" stopOpacity="0" />
          </linearGradient>
        </defs>
        <FlowLine x1={200} y1={145} x2={355} y2={200} tick={tick} speed={0.8} />
        <FlowLine x1={355} y1={255} x2={275} y2={310} tick={tick} speed={1.1} delay={0.5} />
        <FlowLine x1={100} y1={160} x2={55} y2={265} tick={tick} speed={0.9} delay={1} />
        <FlowLine x1={200} y1={165} x2={130} y2={375} tick={tick} speed={0.7} delay={1.5} />
      </svg>

      <div
        className="finance-hero-wallet"
        style={{ transform: `translateY(${walletY}px) rotate(${walletRotate}deg) scale(${balancePulse})` }}
      >
        <div className="finance-hero-wallet__inner">
          <div className="finance-hero-wallet__top">
            <div>
              <span>Saldo total</span>
              <strong>
                $2.480.000<span>,50</span>
              </strong>
            </div>
            <i aria-hidden="true" />
          </div>
          <div className="finance-hero-wallet__bottom">
            <span>Billetera Cenz</span>
            <strong>+8,4% este mes</strong>
          </div>
        </div>
        <span className="finance-hero-wallet__shine" style={{ left: `${shimmerLeft}%` }} aria-hidden="true" />
      </div>

      <article
        className="finance-hero-card finance-hero-card--transactions"
        style={{ transform: `translateY(${chartY * 0.4}px)` }}
      >
        <header>
          <span>Movimientos</span>
          <strong>En vivo</strong>
        </header>
        <div className="finance-hero-transactions">
          {transactions.map((transaction, index) => {
            const offset = (index - transactionIndex + transactions.length) % transactions.length;
            const opacity = offset === 0 ? 1 : offset === 1 ? 0.72 : offset === 2 ? 0.42 : 0;

            return (
              <div
                className="finance-hero-transaction"
                key={transaction.id}
                style={{
                  opacity,
                  transform: `translateY(${offset * 46}px)`,
                }}
              >
                <span>{transaction.badge}</span>
                <p>{transaction.label}</p>
                <strong data-type={transaction.type}>{transaction.amount}</strong>
              </div>
            );
          })}
        </div>
      </article>

      <article
        className="finance-hero-card finance-hero-card--chart"
        style={{ transform: `translateY(${chartY}px) rotate(${chartRotate}deg)` }}
      >
        <header>
          <span>Tendencia mensual</span>
          <strong>+12,3%</strong>
        </header>
        <Sparkline points={chartPoints} progress={chartProgress} />
        <footer>
          {["Ene", "Abr", "Jul", "Oct", "Dic"].map((month) => (
            <span key={month}>{month}</span>
          ))}
        </footer>
      </article>

      <article
        className="finance-hero-card finance-hero-card--budget"
        style={{ transform: `translateY(${budgetY}px) rotate(${budgetRotate}deg)` }}
      >
        <BudgetRing percent={budgetPercent} tick={tick} />
        <div>
          <span>Presupuesto</span>
          <strong>{budgetPercent}%</strong>
          <p>de $620.000 usados</p>
        </div>
      </article>

      <article className="finance-hero-card finance-hero-card--savings" style={{ transform: `translateY(${savingsY}px)` }}>
        <header>
          <span>Meta de ahorro</span>
          <strong>62%</strong>
        </header>
        <p>
          $310.000 <span>/ $500.000</span>
        </p>
        <SavingsBar percent={62} tick={tick} />
        <small>Viaje y emergencias</small>
      </article>

      {floatingMarks.map((mark, index) => {
        const speed = 0.22 + index * 0.09;
        const startX = [25, 105, 195, 295, 370, 445][index];
        const yOffset = (tick * speed * 48 + index * 75) % 440;
        const xWave = Math.sin(tick * speed + index) * 12;
        const opacity = yOffset < 70 ? yOffset / 70 : yOffset > 360 ? (440 - yOffset) / 80 : 0.1 + (index % 3) * 0.05;

        return (
          <span
            aria-hidden="true"
            className="finance-hero-mark"
            key={`${mark}-${index}`}
            style={{
              left: startX + xWave,
              top: yOffset,
              opacity,
            }}
          >
            {mark}
          </span>
        );
      })}
    </div>
  );
}
