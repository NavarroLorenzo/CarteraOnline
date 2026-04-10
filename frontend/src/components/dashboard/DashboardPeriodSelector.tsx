import { motion } from "framer-motion";
import {
  getDashboardPresetLabel,
  type DashboardDateRange,
  type DashboardFilterPreset,
} from "../../lib/dashboard";
import { buttonHover, fadeUp } from "../ui/animation";

type DashboardPeriodSelectorProps = {
  preset: DashboardFilterPreset;
  activePreset: DashboardFilterPreset;
  customRange: DashboardDateRange;
  loading: boolean;
  customRangeValid: boolean;
  validationMessage?: string | null;
  onPresetChange: (preset: DashboardFilterPreset) => void;
  onCustomRangeChange: (nextRange: DashboardDateRange) => void;
  onApplyCustomRange: () => void;
  onClearFilters: () => void;
};

const presetOptions: DashboardFilterPreset[] = ["day", "week", "month", "year", "custom"];

export function DashboardPeriodSelector({
  preset,
  activePreset,
  customRange,
  loading,
  customRangeValid,
  validationMessage,
  onPresetChange,
  onCustomRangeChange,
  onApplyCustomRange,
  onClearFilters,
}: DashboardPeriodSelectorProps) {
  return (
    <motion.section className="panel dashboard-period-selector" variants={fadeUp}>
      <div className="dashboard-period-selector__header">
        <div>
          <span className="eyebrow">Período global</span>
          <h2>Elegí el rango principal del dashboard</h2>
          <p>La vista queda preparada para que resúmenes y analítica respondan desde un único control.</p>
        </div>

        <div className="dashboard-period-selector__status">
          <span>Activo</span>
          <strong>{getDashboardPresetLabel(activePreset)}</strong>
        </div>
      </div>

      <div className="dashboard-segmented dashboard-segmented--wide" role="tablist" aria-label="Selector de período">
        {presetOptions.map((option) => (
          <motion.button
            key={option}
            type="button"
            role="tab"
            aria-selected={preset === option}
            className={`dashboard-segmented__option ${
              preset === option ? "dashboard-segmented__option--active" : ""
            }`}
            onClick={() => onPresetChange(option)}
            whileHover={buttonHover}
            whileTap={{ scale: 0.99 }}
            disabled={loading}
          >
            {getDashboardPresetLabel(option)}
          </motion.button>
        ))}
      </div>

      {preset === "custom" ? (
        <div className="dashboard-period-selector__custom">
          <label className="field">
            <span>Fecha desde</span>
            <input
              type="date"
              value={customRange.date_from}
              onChange={(event) =>
                onCustomRangeChange({
                  ...customRange,
                  date_from: event.target.value,
                })
              }
            />
          </label>

          <label className="field">
            <span>Fecha hasta</span>
            <input
              type="date"
              value={customRange.date_to}
              onChange={(event) =>
                onCustomRangeChange({
                  ...customRange,
                  date_to: event.target.value,
                })
              }
            />
          </label>

          <div className="dashboard-period-selector__actions">
            <motion.button
              type="button"
              className="primary-button"
              onClick={onApplyCustomRange}
              whileHover={buttonHover}
              whileTap={{ scale: 0.99 }}
              disabled={loading || !customRangeValid}
            >
              Aplicar filtro
            </motion.button>

            <motion.button
              type="button"
              className="ghost-button"
              onClick={onClearFilters}
              whileHover={buttonHover}
              whileTap={{ scale: 0.99 }}
              disabled={loading}
            >
              Limpiar filtros
            </motion.button>
          </div>
        </div>
      ) : null}

      {preset === "custom" && validationMessage ? (
        <p className="feedback feedback--warning">{validationMessage}</p>
      ) : null}
    </motion.section>
  );
}
