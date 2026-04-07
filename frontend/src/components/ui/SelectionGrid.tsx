import { motion } from "framer-motion";
import { fadeUp } from "./animation";

type SelectionOption = {
  value: string;
  label: string;
  description?: string;
};

type SelectionGridProps = {
  label: string;
  value: string;
  options: SelectionOption[];
  onChange: (value: string) => void;
  columns?: "compact" | "regular" | "wide";
  emptyMessage?: string;
  disabled?: boolean;
};

export function SelectionGrid({
  label,
  value,
  options,
  onChange,
  columns = "regular",
  emptyMessage,
  disabled = false,
}: SelectionGridProps) {
  return (
    <div className="selection-grid">
      <div className="selection-grid__header">
        <span>{label}</span>
      </div>

      {options.length === 0 ? (
        <div className="selection-grid__empty">{emptyMessage ?? "No hay opciones disponibles."}</div>
      ) : (
        <div className={`selection-grid__options selection-grid__options--${columns}`}>
          {options.map((option) => (
            <motion.button
              key={option.value}
              type="button"
              className={`selection-grid__option ${
                value === option.value ? "selection-grid__option--active" : ""
              }`}
              variants={fadeUp}
              whileTap={{ scale: 0.99 }}
              onClick={() => onChange(option.value)}
              disabled={disabled}
              aria-pressed={value === option.value}
            >
              <strong>{option.label}</strong>
              {option.description ? <span>{option.description}</span> : null}
            </motion.button>
          ))}
        </div>
      )}
    </div>
  );
}
