import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";
import { buttonHover } from "./animation";

export type SelectOption = {
  value: string;
  label: string;
};

type AnimatedSelectProps = {
  value: string;
  options: SelectOption[];
  onChange: (value: string) => void;
  placeholder?: string;
  disabled?: boolean;
  ariaLabel?: string;
};

export function AnimatedSelect({
  value,
  options,
  onChange,
  placeholder = "Seleccionar",
  disabled = false,
  ariaLabel,
}: AnimatedSelectProps) {
  const [open, setOpen] = useState(false);
  const containerRef = useRef<HTMLDivElement | null>(null);

  const selectedOption = useMemo(
    () => options.find((option) => option.value === value),
    [options, value],
  );

  useEffect(() => {
    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setOpen(false);
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
      }
    };

    window.addEventListener("mousedown", handlePointerDown);
    window.addEventListener("keydown", handleEscape);

    return () => {
      window.removeEventListener("mousedown", handlePointerDown);
      window.removeEventListener("keydown", handleEscape);
    };
  }, []);

  return (
    <div
      ref={containerRef}
      className={`animated-select${open ? " animated-select--open" : ""}${disabled ? " animated-select--disabled" : ""}`}
    >
      <motion.button
        type="button"
        className="animated-select__trigger"
        onClick={() => {
          if (!disabled) {
            setOpen((current) => !current);
          }
        }}
        whileHover={!disabled ? buttonHover : undefined}
        whileTap={!disabled ? { scale: 0.995 } : undefined}
        disabled={disabled}
        aria-haspopup="listbox"
        aria-expanded={open}
        aria-label={ariaLabel}
      >
        <span className="animated-select__value">{selectedOption?.label || placeholder}</span>
        <motion.span
          className="animated-select__caret"
          animate={{ rotate: open ? 180 : 0 }}
          transition={{ duration: 0.2, ease: "easeInOut" }}
          aria-hidden="true"
        >
          ▾
        </motion.span>
      </motion.button>

      <AnimatePresence>
        {open ? (
          <motion.div
            className="animated-select__menu"
            initial={{ opacity: 0, y: 10, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 6, scale: 0.985 }}
            transition={{ duration: 0.2, ease: [0.22, 1, 0.36, 1] }}
            role="listbox"
          >
            {options.map((option) => {
              const selected = option.value === value;

              return (
                <motion.button
                  key={option.value}
                  type="button"
                  className={`animated-select__option${selected ? " animated-select__option--selected" : ""}`}
                  onClick={() => {
                    onChange(option.value);
                    setOpen(false);
                  }}
                  whileHover={{ x: 4 }}
                  transition={{ duration: 0.16 }}
                  role="option"
                  aria-selected={selected}
                >
                  <span>{option.label}</span>
                  {selected ? <small aria-hidden="true">✓</small> : null}
                </motion.button>
              );
            })}
          </motion.div>
        ) : null}
      </AnimatePresence>
    </div>
  );
}
