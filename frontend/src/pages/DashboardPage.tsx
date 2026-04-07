import { motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { ApiError } from "../api/client";
import { transactionsApi } from "../api/transactions";
import { CategoryDetailModal } from "../components/dashboard/CategoryDetailModal";
import { CategoryDistributionChart } from "../components/dashboard/CategoryDistributionChart";
import { DashboardFilterBar } from "../components/dashboard/DashboardFilterBar";
import { SummaryStatCard } from "../components/dashboard/SummaryStatCard";
import { TrendLineChart } from "../components/dashboard/TrendLineChart";
import { EmptyState } from "../components/ui/EmptyState";
import { MetricCard } from "../components/ui/MetricCard";
import {
  PresenceMessage,
  Reveal,
  StaggerGroup,
  cardHover,
  fadeLeft,
  fadeRight,
  fadeUp,
} from "../components/ui/animation";
import {
  formatDaySummaryLabel,
  formatMonthSummaryLabel,
  formatPercent,
  getDashboardRange,
  type DashboardDateRange,
  type DashboardFilterPreset,
  type DashboardMode,
  toInputDate,
} from "../lib/dashboard";
import { formatCurrency, formatDate, formatTypeLabel } from "../lib/format";
import type {
  DashboardAnalytics,
  DashboardCategoryDetail,
  DashboardTransactionItem,
} from "../types/api";

export function DashboardPage() {
  const todayString = useMemo(() => toInputDate(new Date()), []);
  const [preset, setPreset] = useState<DashboardFilterPreset>("today");
  const [mode, setMode] = useState<DashboardMode>("advanced");
  const [customRange, setCustomRange] = useState<DashboardDateRange>({
    date_from: todayString,
    date_to: todayString,
  });
  const [dashboard, setDashboard] = useState<DashboardAnalytics | null>(null);
  const [selectedCategoryKey, setSelectedCategoryKey] = useState<string | null>(null);
  const [categoryDetail, setCategoryDetail] = useState<DashboardCategoryDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [detailLoading, setDetailLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const activeRange = useMemo(() => getDashboardRange(preset, customRange), [customRange, preset]);

  useEffect(() => {
    if (!activeRange) {
      setLoading(false);
      return;
    }

    let cancelled = false;

    const loadDashboard = async () => {
      setLoading(true);
      setError(null);
      setSelectedCategoryKey(null);
      setCategoryDetail(null);

      try {
        const data = await transactionsApi.getDashboard(activeRange);

        if (!cancelled) {
          setDashboard(data);
        }
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof ApiError ? err.message : "No se pudo cargar el dashboard");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void loadDashboard();

    return () => {
      cancelled = true;
    };
  }, [activeRange]);

  const handleCategorySelect = async (categoryKey: string) => {
    if (!activeRange) {
      return;
    }

    setSelectedCategoryKey(categoryKey);
    setDetailLoading(true);
    setError(null);

    try {
      const detail = await transactionsApi.getDashboardCategoryDetail(categoryKey, activeRange);
      setCategoryDetail(detail);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo cargar el detalle de la categoría");
    } finally {
      setDetailLoading(false);
    }
  };

  const periodSummary = dashboard?.period_summary;
  const daySummary = dashboard?.day_summary;
  const monthSummary = dashboard?.month_summary;
  const recentTransactions = dashboard?.recent_transactions ?? [];
  const topExpenses = dashboard?.top_expenses ?? [];
  const categories = dashboard?.expense_categories ?? [];
  const comparison = dashboard?.comparison;

  return (
    <div className="page-stack">
      <Reveal onView={false}>
        <section className="hero-card dashboard-hero">
          <div className="dashboard-hero__copy">
            <span className="eyebrow">Dashboard</span>
            <h1>Tu analítica financiera en una sola vista</h1>
            <p>
              Seguís tu balance, entendés hábitos de gasto y explorás tendencias sin mover la lógica sensible
              al frontend.
            </p>
          </div>

          <DashboardFilterBar
            preset={preset}
            mode={mode}
            customRange={customRange}
            loading={loading}
            onPresetChange={setPreset}
            onModeChange={setMode}
            onCustomRangeChange={setCustomRange}
          />

          {preset === "custom" && !activeRange ? (
            <p className="feedback feedback--warning">Elegí fecha de inicio y fin para cargar el rango personalizado.</p>
          ) : null}
        </section>
      </Reveal>

      <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>

      <Reveal onView={false}>
        <section className="panel dashboard-summary-panel">
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Resumen clave</span>
              <h2>Lectura rápida del día y del mes</h2>
            </div>
          </div>

          <div className="dashboard-summary-groups">
            <div className="dashboard-summary-group">
              <div className="dashboard-summary-group__header">
                <h3>Resumen del día</h3>
                <span>{formatDaySummaryLabel(activeRange?.date_to ?? todayString)}</span>
              </div>

              <StaggerGroup className="dashboard-summary-grid" onView={false}>
                <SummaryStatCard
                  label="Ingresos del día"
                  value={daySummary ? formatCurrency(daySummary.income_total) : "Cargando..."}
                  tone="positive"
                />
                <SummaryStatCard
                  label="Gastos del día"
                  value={daySummary ? formatCurrency(daySummary.expense_total) : "Cargando..."}
                  tone="negative"
                />
                <SummaryStatCard
                  label="Balance del día"
                  value={daySummary ? formatCurrency(daySummary.net_balance) : "Cargando..."}
                  tone="accent"
                />
              </StaggerGroup>
            </div>

            <div className="dashboard-summary-group">
              <div className="dashboard-summary-group__header">
                <h3>Resumen del mes</h3>
                <span>{formatMonthSummaryLabel(activeRange?.date_to ?? todayString)}</span>
              </div>

              <StaggerGroup className="dashboard-summary-grid" onView={false}>
                <SummaryStatCard
                  label="Ingresos del mes"
                  value={monthSummary ? formatCurrency(monthSummary.income_total) : "Cargando..."}
                  tone="positive"
                />
                <SummaryStatCard
                  label="Gastos del mes"
                  value={monthSummary ? formatCurrency(monthSummary.expense_total) : "Cargando..."}
                  tone="negative"
                />
                <SummaryStatCard
                  label="Ahorro del mes"
                  value={monthSummary ? formatCurrency(monthSummary.net_balance) : "Cargando..."}
                  tone="accent"
                />
              </StaggerGroup>
            </div>
          </div>
        </section>
      </Reveal>

      <StaggerGroup className="metrics-grid" onView={false}>
        <MetricCard
          label="Balance del período"
          value={periodSummary ? formatCurrency(periodSummary.net_balance) : "Cargando..."}
          tone="accent"
        />
        <MetricCard
          label="Ingresos analizados"
          value={periodSummary ? formatCurrency(periodSummary.income_total) : "Cargando..."}
          tone="positive"
        />
        <MetricCard
          label="Gastos analizados"
          value={periodSummary ? formatCurrency(periodSummary.expense_total) : "Cargando..."}
          tone="negative"
        />
        <MetricCard
          label="Movimientos recientes"
          value={loading ? "Cargando..." : String(recentTransactions.length)}
        />
      </StaggerGroup>

      {mode === "advanced" ? (
        <StaggerGroup className="dashboard-analytics-grid" onView={false}>
          <motion.article className="panel dashboard-panel dashboard-panel--wide" variants={fadeLeft}>
            <div className="panel-heading">
              <div>
                <span className="eyebrow">Tendencia</span>
                <h2>Ingresos vs gastos</h2>
                <p>La serie se adapta automáticamente al rango activo.</p>
              </div>
            </div>

            {dashboard && dashboard.trend.length > 0 ? (
              <TrendLineChart points={dashboard.trend} />
            ) : (
              <EmptyState
                title="Sin tendencia disponible"
                description="Cuando tengas ingresos o gastos en este período, vas a ver la evolución acá."
              />
            )}
          </motion.article>

          <motion.article className="panel dashboard-panel" variants={fadeRight}>
            <div className="panel-heading">
              <div>
                <span className="eyebrow">Comparación</span>
                <h2>Mes actual vs anterior</h2>
              </div>
            </div>

            {comparison ? (
              <div className="comparison-grid">
                <ComparisonBlock
                  label="Ingresos"
                  currentValue={comparison.current_month.income_total}
                  previousValue={comparison.previous_month.income_total}
                  change={comparison.income_change_pct}
                  variant="income"
                />
                <ComparisonBlock
                  label="Gastos"
                  currentValue={comparison.current_month.expense_total}
                  previousValue={comparison.previous_month.expense_total}
                  change={comparison.expense_change_pct}
                  variant="expense"
                />
              </div>
            ) : (
              <EmptyState
                title="Todavía no hay comparación"
                description="Hace falta al menos un mes con movimientos para calcular la variación."
              />
            )}
          </motion.article>

          <motion.article className="panel dashboard-panel" variants={fadeLeft}>
            <div className="panel-heading">
              <div>
                <span className="eyebrow">Categorías</span>
                <h2>Distribución de gastos</h2>
                <p>Seleccioná una categoría para abrir su detalle.</p>
              </div>
            </div>

            <StaggerGroup onView={false}>
              <CategoryDistributionChart
                categories={categories}
                activeCategoryKey={selectedCategoryKey}
                onCategorySelect={(categoryKey) => void handleCategorySelect(categoryKey)}
              />
            </StaggerGroup>
          </motion.article>

          <motion.article className="panel dashboard-panel" variants={fadeRight}>
            <div className="panel-heading">
              <div>
                <span className="eyebrow">Top gastos</span>
                <h2>Los montos más altos del período</h2>
              </div>
            </div>

            {topExpenses.length === 0 ? (
              <EmptyState
                title="No hay gastos para rankear"
                description="Cuando registres gastos, acá vas a ver cuáles pesan más."
              />
            ) : (
              <StaggerGroup className="stack-list" onView={false}>
                {topExpenses.map((transaction) => (
                  <TopExpenseRow key={transaction.id} transaction={transaction} />
                ))}
              </StaggerGroup>
            )}
          </motion.article>
        </StaggerGroup>
      ) : null}

      <Reveal onView={false}>
        <section className="panel dashboard-panel">
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Actividad</span>
              <h2>Últimos movimientos</h2>
              <p>Vista rápida de lo más reciente dentro del filtro actual.</p>
            </div>
          </div>

          {recentTransactions.length === 0 && !loading ? (
            <EmptyState
              title="No hay movimientos para mostrar"
              description="Probá otro rango o registrá una transacción para empezar a ver actividad."
            />
          ) : (
            <StaggerGroup className="stack-list" onView={false}>
              {recentTransactions.map((transaction) => (
                <motion.div
                  key={transaction.id}
                  className={`list-row list-row--transaction list-row--${transaction.type}`}
                  variants={fadeUp}
                  whileHover={cardHover}
                >
                  <div>
                    <div className="list-row__title">
                      <strong>{transaction.title}</strong>
                      <span className={`category-badge category-badge--${transaction.type}`}>
                        {transaction.category_label}
                      </span>
                    </div>
                    <p>
                      {transaction.account_name || `Cuenta ${transaction.account_id}`} ·{" "}
                      {formatTypeLabel(transaction.type)}
                    </p>
                    {transaction.description ? <small>{transaction.description}</small> : null}
                  </div>

                  <div className="list-row__meta">
                    <strong>{formatCurrency(transaction.amount)}</strong>
                    <span>{formatDate(transaction.created_at)}</span>
                  </div>
                </motion.div>
              ))}
            </StaggerGroup>
          )}
        </section>
      </Reveal>

      <CategoryDetailModal
        detail={categoryDetail}
        loading={detailLoading}
        onClose={() => {
          setSelectedCategoryKey(null);
          setCategoryDetail(null);
          setDetailLoading(false);
        }}
      />
    </div>
  );
}

type ComparisonBlockProps = {
  label: string;
  currentValue: number;
  previousValue: number;
  change: number;
  variant: "income" | "expense";
};

function ComparisonBlock({
  label,
  currentValue,
  previousValue,
  change,
  variant,
}: ComparisonBlockProps) {
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
      <small>Mes anterior: {formatCurrency(previousValue)}</small>
    </motion.div>
  );
}

type TopExpenseRowProps = {
  transaction: DashboardTransactionItem;
};

function TopExpenseRow({ transaction }: TopExpenseRowProps) {
  return (
    <motion.div
      className="list-row list-row--transaction list-row--expense"
      variants={fadeUp}
      whileHover={cardHover}
    >
      <div>
        <strong>{transaction.title}</strong>
        <p>
          {transaction.category_label} · {transaction.account_name || `Cuenta ${transaction.account_id}`}
        </p>
      </div>

      <div className="list-row__meta">
        <strong>{formatCurrency(transaction.amount)}</strong>
        <span>{formatDate(transaction.created_at)}</span>
      </div>
    </motion.div>
  );
}
