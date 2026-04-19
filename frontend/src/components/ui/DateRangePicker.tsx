import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";
import { formatDate } from "../../lib/format";
import { buttonHover } from "./animation";

type DateRangeValue = {
  from: string;
  to: string;
};

type DateRangePickerProps = {
  value: DateRangeValue;
  onChange: (value: DateRangeValue) => void;
  label?: string;
  disabled?: boolean;
  className?: string;
};

type ActiveField = "from" | "to";
type CalendarMode = "days" | "years";

const weekdayLabels = ["L", "M", "M", "J", "V", "S", "D"];
const monthFormatter = new Intl.DateTimeFormat("es-AR", { month: "long" });

export function DateRangePicker({
  value,
  onChange,
  label = "Fechas",
  disabled = false,
  className,
}: DateRangePickerProps) {
  const [open, setOpen] = useState(false);
  const [activeField, setActiveField] = useState<ActiveField>("from");
  const [mode, setMode] = useState<CalendarMode>("days");
  const containerRef = useRef<HTMLDivElement | null>(null);

  const selectedFrom = useMemo(() => parseInputDate(value.from), [value.from]);
  const selectedTo = useMemo(() => parseInputDate(value.to), [value.to]);
  const selectedAnchor = activeField === "to" ? selectedTo || selectedFrom : selectedFrom || selectedTo;
  const [visibleMonth, setVisibleMonth] = useState(() => startOfMonth(selectedAnchor || new Date()));

  useEffect(() => {
    if (open) {
      setVisibleMonth(startOfMonth(selectedAnchor || new Date()));
    }
  }, [open, selectedAnchor]);

  useEffect(() => {
    const handlePointerDown = (event: MouseEvent) => {
      if (!containerRef.current?.contains(event.target as Node)) {
        setOpen(false);
        setMode("days");
      }
    };

    const handleEscape = (event: KeyboardEvent) => {
      if (event.key === "Escape") {
        setOpen(false);
        setMode("days");
      }
    };

    window.addEventListener("mousedown", handlePointerDown);
    window.addEventListener("keydown", handleEscape);

    return () => {
      window.removeEventListener("mousedown", handlePointerDown);
      window.removeEventListener("keydown", handleEscape);
    };
  }, []);

  const calendarDays = useMemo(() => buildCalendarDays(visibleMonth), [visibleMonth]);
  const yearOptions = useMemo(() => buildYearOptions(visibleMonth.getFullYear()), [visibleMonth]);

  const handleTriggerClick = (field: ActiveField) => {
    if (disabled) {
      return;
    }

    setActiveField(field);
    setMode("days");
    setOpen(true);
  };

  const handleDaySelect = (date: Date) => {
    const selected = toInputDate(date);

    if (activeField === "from") {
      onChange({
        from: selected,
        to: value.to && selected > value.to ? "" : value.to,
      });
      setActiveField("to");
      return;
    }

    if (value.from && selected < value.from) {
      onChange({ from: selected, to: value.from });
    } else {
      onChange({ from: value.from, to: selected });
    }

    setOpen(false);
    setMode("days");
  };

  const handleQuickRange = (range: DateRangeValue) => {
    onChange(range);
    setVisibleMonth(startOfMonth(parseInputDate(range.to) || parseInputDate(range.from) || new Date()));
    setActiveField("from");
    setMode("days");
    setOpen(false);
  };

  const handleClear = () => {
    onChange({ from: "", to: "" });
    setActiveField("from");
    setMode("days");
    setOpen(false);
  };

  const goToMonth = (offset: number) => {
    setVisibleMonth((current) => new Date(current.getFullYear(), current.getMonth() + offset, 1));
    setMode("days");
  };

  const setVisibleYear = (year: number) => {
    setVisibleMonth((current) => new Date(year, current.getMonth(), 1));
    setMode("days");
  };

  const today = startOfDay(new Date());
  const quickRanges = getQuickRanges(today);
  const monthLabel = capitalize(monthFormatter.format(visibleMonth));

  return (
    <div
      ref={containerRef}
      className={`date-range-picker${open ? " date-range-picker--open" : ""}${disabled ? " date-range-picker--disabled" : ""}${
        className ? ` ${className}` : ""
      }`}
    >
      <div className="date-range-picker__header">
        <span>{label}</span>
        {(value.from || value.to) && !disabled ? (
          <button type="button" className="date-range-picker__clear" onClick={handleClear}>
            Limpiar fechas
          </button>
        ) : null}
      </div>

      <div className="date-range-picker__fields">
        <motion.button
          type="button"
          className={`date-range-picker__field${activeField === "from" && open ? " date-range-picker__field--active" : ""}`}
          onClick={() => handleTriggerClick("from")}
          whileHover={!disabled ? buttonHover : undefined}
          whileTap={!disabled ? { scale: 0.995 } : undefined}
          disabled={disabled}
          aria-haspopup="dialog"
          aria-expanded={open && activeField === "from"}
        >
          <small>Fecha desde</small>
          <strong>{selectedFrom ? formatDate(selectedFrom) : "Elegir fecha"}</strong>
        </motion.button>

        <motion.button
          type="button"
          className={`date-range-picker__field${activeField === "to" && open ? " date-range-picker__field--active" : ""}`}
          onClick={() => handleTriggerClick("to")}
          whileHover={!disabled ? buttonHover : undefined}
          whileTap={!disabled ? { scale: 0.995 } : undefined}
          disabled={disabled}
          aria-haspopup="dialog"
          aria-expanded={open && activeField === "to"}
        >
          <small>Fecha hasta</small>
          <strong>{selectedTo ? formatDate(selectedTo) : "Elegir fecha"}</strong>
        </motion.button>
      </div>

      <AnimatePresence>
        {open ? (
          <motion.div
            className="date-range-picker__popup"
            initial={{ opacity: 0, y: 12, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 8, scale: 0.985 }}
            transition={{ duration: 0.2, ease: [0.22, 1, 0.36, 1] }}
            role="dialog"
            aria-label="Selector de rango de fechas"
          >
            <div className="date-range-picker__quick-actions">
              {quickRanges.map((quickRange) => (
                <button
                  key={quickRange.label}
                  type="button"
                  className="date-range-picker__quick-action"
                  onClick={() => handleQuickRange(quickRange.value)}
                >
                  {quickRange.label}
                </button>
              ))}
            </div>

            <div className="date-range-picker__nav">
              <button type="button" onClick={() => goToMonth(-1)} aria-label="Mes anterior">
                ‹
              </button>
              <button
                type="button"
                className="date-range-picker__month-button"
                onClick={() => setMode((current) => (current === "years" ? "days" : "years"))}
              >
                {monthLabel} {visibleMonth.getFullYear()}
              </button>
              <button type="button" onClick={() => goToMonth(1)} aria-label="Mes siguiente">
                ›
              </button>
            </div>

            {mode === "years" ? (
              <div className="date-range-picker__years" aria-label="Seleccionar año">
                <button type="button" onClick={() => setVisibleMonth((current) => new Date(current.getFullYear() - 12, current.getMonth(), 1))}>
                  Años anteriores
                </button>
                <div className="date-range-picker__year-grid">
                  {yearOptions.map((year) => (
                    <button
                      key={year}
                      type="button"
                      className={year === visibleMonth.getFullYear() ? "date-range-picker__year--selected" : ""}
                      onClick={() => setVisibleYear(year)}
                    >
                      {year}
                    </button>
                  ))}
                </div>
                <button type="button" onClick={() => setVisibleMonth((current) => new Date(current.getFullYear() + 12, current.getMonth(), 1))}>
                  Años siguientes
                </button>
              </div>
            ) : (
              <>
                <div className="date-range-picker__weekdays" aria-hidden="true">
                  {weekdayLabels.map((weekday, index) => (
                    <span key={`${weekday}-${index}`}>{weekday}</span>
                  ))}
                </div>

                <div className="date-range-picker__days">
                  {calendarDays.map((date) => {
                    const dateValue = toInputDate(date);
                    const selected = dateValue === value.from || dateValue === value.to;
                    const rangeStart = dateValue === value.from;
                    const rangeEnd = dateValue === value.to;
                    const inRange = Boolean(value.from && value.to && dateValue > value.from && dateValue < value.to);
                    const outsideMonth = date.getMonth() !== visibleMonth.getMonth();
                    const isToday = isSameDay(date, today);

                    return (
                      <button
                        key={dateValue}
                        type="button"
                        className={[
                          "date-range-picker__day",
                          selected ? "date-range-picker__day--selected" : "",
                          rangeStart ? "date-range-picker__day--range-start" : "",
                          rangeEnd ? "date-range-picker__day--range-end" : "",
                          inRange ? "date-range-picker__day--in-range" : "",
                          outsideMonth ? "date-range-picker__day--muted" : "",
                          isToday ? "date-range-picker__day--today" : "",
                        ]
                          .filter(Boolean)
                          .join(" ")}
                        onClick={() => handleDaySelect(date)}
                        aria-pressed={selected}
                      >
                        {date.getDate()}
                      </button>
                    );
                  })}
                </div>
              </>
            )}
          </motion.div>
        ) : null}
      </AnimatePresence>
    </div>
  );
}

function buildCalendarDays(month: Date): Date[] {
  const firstDay = startOfMonth(month);
  const startOffset = (firstDay.getDay() + 6) % 7;
  const calendarStart = new Date(firstDay);
  calendarStart.setDate(firstDay.getDate() - startOffset);

  return Array.from({ length: 42 }, (_, index) => {
    const date = new Date(calendarStart);
    date.setDate(calendarStart.getDate() + index);
    return date;
  });
}

function buildYearOptions(year: number): number[] {
  const firstYear = year - 5;
  return Array.from({ length: 12 }, (_, index) => firstYear + index);
}

function getQuickRanges(today: Date): Array<{ label: string; value: DateRangeValue }> {
  const weekStart = new Date(today);
  weekStart.setDate(today.getDate() - 6);

  return [
    { label: "Hoy", value: { from: toInputDate(today), to: toInputDate(today) } },
    { label: "Últimos 7 días", value: { from: toInputDate(weekStart), to: toInputDate(today) } },
    {
      label: "Este mes",
      value: {
        from: toInputDate(new Date(today.getFullYear(), today.getMonth(), 1)),
        to: toInputDate(today),
      },
    },
  ];
}

function parseInputDate(value?: string): Date | null {
  if (!value) {
    return null;
  }

  const [year, month, day] = value.split("-").map(Number);
  const parsed = new Date(year, month - 1, day);

  if (
    Number.isNaN(parsed.getTime()) ||
    parsed.getFullYear() !== year ||
    parsed.getMonth() !== month - 1 ||
    parsed.getDate() !== day
  ) {
    return null;
  }

  return parsed;
}

function toInputDate(value: Date): string {
  const year = value.getFullYear();
  const month = String(value.getMonth() + 1).padStart(2, "0");
  const day = String(value.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

function startOfDay(value: Date): Date {
  return new Date(value.getFullYear(), value.getMonth(), value.getDate());
}

function startOfMonth(value: Date): Date {
  return new Date(value.getFullYear(), value.getMonth(), 1);
}

function isSameDay(left: Date, right: Date): boolean {
  return (
    left.getFullYear() === right.getFullYear() &&
    left.getMonth() === right.getMonth() &&
    left.getDate() === right.getDate()
  );
}

function capitalize(value: string): string {
  return value.charAt(0).toUpperCase() + value.slice(1);
}
