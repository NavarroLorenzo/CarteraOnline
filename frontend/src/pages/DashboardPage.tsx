import { useEffect, useMemo, useState } from "react";
import { ApiError } from "../api/client";
import { transactionsApi } from "../api/transactions";
import { CategoryDetailModal } from "../components/dashboard/CategoryDetailModal";
import { DashboardCategoryDistribution } from "../components/dashboard/DashboardCategoryDistribution";
import { DashboardComparison } from "../components/dashboard/DashboardComparison";
import { DashboardIncomeVsExpenses } from "../components/dashboard/DashboardIncomeVsExpenses";
import { DashboardPeriodSelector } from "../components/dashboard/DashboardPeriodSelector";
import { DashboardRecentTransactions } from "../components/dashboard/DashboardRecentTransactions";
import { DashboardSummary } from "../components/dashboard/DashboardSummary";
import { DashboardTopExpenses } from "../components/dashboard/DashboardTopExpenses";
import { MetricCard } from "../components/ui/MetricCard";
import { PresenceMessage, Reveal, StaggerGroup } from "../components/ui/animation";
import {
  buildDashboardPeriodQuery,
  getDashboardRange,
  isDashboardCustomRangeValid,
  type DashboardActivePeriod,
  type DashboardDateRange,
  type DashboardFilterPreset,
  toInputDate,
} from "../lib/dashboard";
import { formatCurrency } from "../lib/format";
import type { DashboardAnalytics, DashboardCategoryDetail } from "../types/api";

export function DashboardPage() {
  const todayString = useMemo(() => toInputDate(new Date()), []);
  const [activePeriod, setActivePeriod] = useState<DashboardActivePeriod>({
    preset: "month",
    customRange: {
      date_from: todayString,
      date_to: todayString,
    },
  });
  const [selectorPreset, setSelectorPreset] = useState<DashboardFilterPreset>("month");
  const [customRangeDraft, setCustomRangeDraft] = useState<DashboardDateRange>({
    date_from: todayString,
    date_to: todayString,
  });
  const [dashboard, setDashboard] = useState<DashboardAnalytics | null>(null);
  const [selectedCategoryKey, setSelectedCategoryKey] = useState<string | null>(null);
  const [categoryDetail, setCategoryDetail] = useState<DashboardCategoryDetail | null>(null);
  const [loading, setLoading] = useState(true);
  const [detailLoading, setDetailLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const activeRange = useMemo(
    () => getDashboardRange(activePeriod.preset, activePeriod.customRange),
    [activePeriod.customRange, activePeriod.preset],
  );
  const customRangeValid = useMemo(() => isDashboardCustomRangeValid(customRangeDraft), [customRangeDraft]);
  const dashboardQuery = useMemo(() => buildDashboardPeriodQuery(activePeriod), [activePeriod]);

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
        const data = await transactionsApi.getDashboard(dashboardQuery);

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
  }, [activeRange, dashboardQuery]);

  const handleCategorySelect = async (categoryKey: string) => {
    if (!activeRange) {
      return;
    }

    setSelectedCategoryKey(categoryKey);
    setDetailLoading(true);
    setError(null);

    try {
      const detail = await transactionsApi.getDashboardCategoryDetail(categoryKey, dashboardQuery);
      setCategoryDetail(detail);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo cargar el detalle de la categoría");
    } finally {
      setDetailLoading(false);
    }
  };

  const handlePresetChange = (nextPreset: DashboardFilterPreset) => {
    setSelectorPreset(nextPreset);
    setError(null);

    if (nextPreset === "custom") {
      if (activeRange) {
        setCustomRangeDraft(activeRange);
      }
      return;
    }

    setActivePeriod((current) => ({
      ...current,
      preset: nextPreset,
    }));
  };

  const handleApplyCustomRange = () => {
    if (!customRangeValid) {
      return;
    }

    setActivePeriod({
      preset: "custom",
      customRange: customRangeDraft,
    });
    setSelectorPreset("custom");
    setSelectedCategoryKey(null);
    setCategoryDetail(null);
  };

  const handleClearFilters = () => {
    const resetRange = {
      date_from: todayString,
      date_to: todayString,
    };

    setSelectorPreset("month");
    setCustomRangeDraft(resetRange);
    setActivePeriod({
      preset: "month",
      customRange: resetRange,
    });
    setSelectedCategoryKey(null);
    setCategoryDetail(null);
    setError(null);
  };

  const periodSummary = dashboard?.period_summary;
  const recentTransactions = dashboard?.recent_transactions ?? [];
  const topExpenses = dashboard?.top_expenses ?? [];
  const categories = dashboard?.expense_categories ?? [];
  const comparison = dashboard?.comparison;
  const trendPoints = dashboard?.trend ?? [];
  const metricFallback = loading ? "Cargando información..." : "Sin datos";
  const customRangeValidationMessage =
    selectorPreset === "custom" && !customRangeValid
      ? "Elegí una fecha desde y una fecha hasta válidas para aplicar el rango personalizado."
      : null;

  return (
    <div className="page-stack">
      <Reveal onView={false}>
      </Reveal>

      <Reveal onView={false}>
        <DashboardPeriodSelector
          preset={selectorPreset}
          activePreset={activePeriod.preset}
          customRange={customRangeDraft}
          loading={loading}
          customRangeValid={customRangeValid}
          validationMessage={customRangeValidationMessage}
          onPresetChange={handlePresetChange}
          onCustomRangeChange={setCustomRangeDraft}
          onApplyCustomRange={handleApplyCustomRange}
          onClearFilters={handleClearFilters}
        />
      </Reveal>

      <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>

      <Reveal onView={false}>
        <DashboardSummary preset={activePeriod.preset} range={activeRange} summary={periodSummary} loading={loading} />
      </Reveal>

      <StaggerGroup className="metrics-grid" onView={false}>
        <MetricCard
          label="Balance del período"
          value={periodSummary ? formatCurrency(periodSummary.net_balance) : metricFallback}
          tone="accent"
        />
        <MetricCard
          label="Ingresos analizados"
          value={periodSummary ? formatCurrency(periodSummary.income_total) : metricFallback}
          tone="positive"
        />
        <MetricCard
          label="Gastos analizados"
          value={periodSummary ? formatCurrency(periodSummary.expense_total) : metricFallback}
          tone="negative"
        />
        <MetricCard
          label="Movimientos del período"
          value={periodSummary ? String(periodSummary.transactions_count) : metricFallback}
        />
      </StaggerGroup>

      <StaggerGroup className="dashboard-analytics-grid" onView={false}>
        <DashboardComparison comparison={comparison} loading={loading} />
        <DashboardCategoryDistribution
          categories={categories}
          activeCategoryKey={selectedCategoryKey}
          loading={loading}
          onCategorySelect={(categoryKey) => void handleCategorySelect(categoryKey)}
        />
        <DashboardTopExpenses transactions={topExpenses} loading={loading} />
        <DashboardIncomeVsExpenses points={trendPoints} loading={loading} />
      </StaggerGroup>

      <DashboardRecentTransactions transactions={recentTransactions} loading={loading} />

      <CategoryDetailModal
        detail={categoryDetail}
        categories={categories}
        activeCategoryKey={selectedCategoryKey}
        loading={detailLoading}
        onCategorySelect={(categoryKey) => void handleCategorySelect(categoryKey)}
        onClose={() => {
          setSelectedCategoryKey(null);
          setCategoryDetail(null);
          setDetailLoading(false);
        }}
      />
    </div>
  );
}
