import { motion } from "framer-motion";
import type { DashboardDateRange, DashboardFilterPreset, DashboardMode } from "../../lib/dashboard";
import { buttonHover, fadeUp } from "../ui/animation";

type DashboardFilterBarProps = {
  preset: DashboardFilterPreset;
  mode: DashboardMode;
  customRange: DashboardDateRange;
  loading: boolean;
  onPresetChange: (preset: DashboardFilterPreset) => void;
  onModeChange: (mode: DashboardMode) => void;
  onCustomRangeChange: (nextRange: DashboardDateRange) => void;
};

const presetOptions: Array<{ value: DashboardFilterPreset; label: string }> = [
  { value: "today", label: "Hoy" },
  { value: "last_7_days", label: "Últimos 7 días" },
  { value: "this_month", label: "Este mes" },
  { value: "custom", label: "Rango personalizado" },
];

const modeOptions: Array<{ value: DashboardMode; label: string }> = [
  { value: "simple", label: "Modo simple" },
  { value: "advanced", label: "Modo avanzado" },
];

export function DashboardFilterBar({
  preset,
  mode,
  customRange,
  loading,
  onPresetChange,
  onModeChange,
  onCustomRangeChange,
}: DashboardFilterBarProps) {
  return (
    <motion.div className="dashboard-toolbar" variants={fadeUp}>
      <div className="dashboard-toolbar__section">
        <span className="dashboard-toolbar__label">Período</span>
        <div className="dashboard-segmented">
          {presetOptions.map((option) => (
            <motion.button
              key={option.value}
              type="button"
              className={`dashboard-segmented__option ${
                preset === option.value ? "dashboard-segmented__option--active" : ""
              }`}
              onClick={() => onPresetChange(option.value)}
              whileHover={buttonHover}
              whileTap={{ scale: 0.99 }}
              disabled={loading}
            >
              {option.label}
            </motion.button>
          ))}
        </div>

        {preset === "custom" ? (
          <div className="dashboard-toolbar__dates">
            <label className="field">
              <span>Desde</span>
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
              <span>Hasta</span>
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
          </div>
        ) : null}
      </div>

      <div className="dashboard-toolbar__section dashboard-toolbar__section--mode">
        <span className="dashboard-toolbar__label">Vista</span>
        <div className="dashboard-segmented dashboard-segmented--compact">
          {modeOptions.map((option) => (
            <motion.button
              key={option.value}
              type="button"
              className={`dashboard-segmented__option ${
                mode === option.value ? "dashboard-segmented__option--active" : ""
              }`}
              onClick={() => onModeChange(option.value)}
              whileHover={buttonHover}
              whileTap={{ scale: 0.99 }}
            >
              {option.label}
            </motion.button>
          ))}
        </div>
      </div>
    </motion.div>
  );
}
