import { motion } from "framer-motion";
import type { DashboardTrendPoint } from "../../types/api";
import { IncomeExpenseAreaChart } from "./IncomeExpenseAreaChart";
import { EmptyState } from "../ui/EmptyState";
import { fadeLeft } from "../ui/animation";

type DashboardIncomeVsExpensesProps = {
  points: DashboardTrendPoint[];
  loading: boolean;
};

export function DashboardIncomeVsExpenses({ points, loading }: DashboardIncomeVsExpensesProps) {
  return (
    <motion.article className="panel dashboard-panel dashboard-panel--trend" variants={fadeLeft}>
      <div className="panel-heading">
        <div>
          <span className="eyebrow">Tendencia</span>
          <h2>Ingresos vs gastos</h2>
          <p>La serie mantiene el foco en el período activo del dashboard.</p>
        </div>
      </div>

      {loading ? (
        <p className="feedback">Cargando información...</p>
      ) : points.length > 0 ? (
        <IncomeExpenseAreaChart points={points} />
      ) : (
        <EmptyState
          title="No hay tendencia disponible"
          description="Cuando haya ingresos o gastos en este período, vas a ver la evolución en este panel."
        />
      )}
    </motion.article>
  );
}
