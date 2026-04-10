import { motion } from "framer-motion";
import { CategoryPieChart } from "./CategoryPieChart";
import type { DashboardCategorySummary } from "../../types/api";
import { EmptyState } from "../ui/EmptyState";
import { StaggerGroup, fadeLeft } from "../ui/animation";

type DashboardCategoryDistributionProps = {
  categories: DashboardCategorySummary[];
  activeCategoryKey: string | null;
  loading: boolean;
  onCategorySelect: (categoryKey: string) => void;
};

export function DashboardCategoryDistribution({
  categories,
  activeCategoryKey,
  loading,
  onCategorySelect,
}: DashboardCategoryDistributionProps) {
  return (
    <motion.article className="panel dashboard-panel dashboard-panel--distribution" variants={fadeLeft}>
      <div className="panel-heading">
        <div>
          <span className="eyebrow">Categorías</span>
          <h2>Distribución de gastos</h2>
          <p>Seleccioná una categoría para ver su detalle.</p>
        </div>
      </div>

      {loading ? (
        <p className="feedback">Cargando información...</p>
      ) : categories.length === 0 ? (
        <EmptyState
          title="No hay gastos para analizar"
          description="Cuando registres gastos en este período, vas a ver la distribución por categoría."
        />
      ) : (
        <StaggerGroup onView={false}>
          <CategoryPieChart
            categories={categories}
            activeCategoryKey={activeCategoryKey}
            onCategorySelect={(categoryKey) => void onCategorySelect(categoryKey)}
          />
        </StaggerGroup>
      )}
    </motion.article>
  );
}
