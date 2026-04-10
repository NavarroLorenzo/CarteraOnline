export type DashboardFilterPreset = "day" | "week" | "month" | "year" | "custom";

export type DashboardDateRange = {
  date_from: string;
  date_to: string;
};

export type DashboardActivePeriod = {
  preset: DashboardFilterPreset;
  customRange: DashboardDateRange;
};

const monthFormatter = new Intl.DateTimeFormat("es-AR", {
  month: "long",
  year: "numeric",
});

const shortDayFormatter = new Intl.DateTimeFormat("es-AR", {
  day: "numeric",
  month: "short",
});

const fullDayFormatter = new Intl.DateTimeFormat("es-AR", {
  day: "numeric",
  month: "short",
  year: "numeric",
});

export function getDashboardRange(
  preset: DashboardFilterPreset,
  customRange: DashboardDateRange,
  now = new Date(),
): DashboardDateRange | null {
  if (preset === "custom") {
    if (!isDashboardCustomRangeValid(customRange)) {
      return null;
    }

    return customRange;
  }

  const today = startOfDay(now);

  if (preset === "day") {
    return {
      date_from: toInputDate(today),
      date_to: toInputDate(today),
    };
  }

  if (preset === "week") {
    const weekStart = startOfWeek(today);
    return {
      date_from: toInputDate(weekStart),
      date_to: toInputDate(today),
    };
  }

  if (preset === "month") {
    const monthStart = new Date(today.getFullYear(), today.getMonth(), 1);
    return {
      date_from: toInputDate(monthStart),
      date_to: toInputDate(today),
    };
  }

  const yearStart = new Date(today.getFullYear(), 0, 1);
  return {
    date_from: toInputDate(yearStart),
    date_to: toInputDate(today),
  };
}

export function toInputDate(value: Date): string {
  const year = value.getFullYear();
  const month = String(value.getMonth() + 1).padStart(2, "0");
  const day = String(value.getDate()).padStart(2, "0");
  return `${year}-${month}-${day}`;
}

export function startOfDay(value: Date): Date {
  return new Date(value.getFullYear(), value.getMonth(), value.getDate());
}

export function startOfWeek(value: Date): Date {
  const day = value.getDay();
  const diff = day === 0 ? -6 : 1 - day;
  const result = new Date(value);
  result.setDate(value.getDate() + diff);
  return startOfDay(result);
}

export function endOfWeek(value: Date): Date {
  const result = startOfDay(value);
  result.setDate(result.getDate() + 6);
  return result;
}

export function isDashboardCustomRangeValid(range: DashboardDateRange): boolean {
  return Boolean(range.date_from && range.date_to && range.date_from <= range.date_to);
}

export function buildDashboardPeriodQuery(period: DashboardActivePeriod): {
  period: DashboardFilterPreset;
  dateFrom?: string;
  dateTo?: string;
} {
  return period.preset === "custom"
    ? {
        period: period.preset,
        dateFrom: period.customRange.date_from,
        dateTo: period.customRange.date_to,
      }
    : {
        period: period.preset,
      };
}

export function parseInputDate(value: string): Date {
  const [year, month, day] = value.split("-").map(Number);
  return new Date(year, month - 1, day);
}

export function formatDaySummaryLabel(value: string): string {
  return shortDayFormatter.format(parseInputDate(value));
}

export function formatMonthSummaryLabel(value: string): string {
  const monthDate = parseInputDate(value);
  const formatted = monthFormatter.format(monthDate);
  return formatted.charAt(0).toUpperCase() + formatted.slice(1);
}

export function formatDashboardRangeLabel(
  preset: DashboardFilterPreset,
  range: DashboardDateRange | null,
): string {
  if (!range) {
    return "Elegí un rango para ver el resumen";
  }

  if (preset === "day") {
    return fullDayFormatter.format(parseInputDate(range.date_to));
  }

  if (preset === "week") {
    return `${shortDayFormatter.format(parseInputDate(range.date_from))} - ${fullDayFormatter.format(
      parseInputDate(range.date_to),
    )}`;
  }

  if (preset === "month") {
    return formatMonthSummaryLabel(range.date_to);
  }

  if (preset === "year") {
    return parseInputDate(range.date_to).getFullYear().toString();
  }

  if (range.date_from === range.date_to) {
    return fullDayFormatter.format(parseInputDate(range.date_from));
  }

  return `${shortDayFormatter.format(parseInputDate(range.date_from))} - ${fullDayFormatter.format(
    parseInputDate(range.date_to),
  )}`;
}

export function getDashboardSummaryTitle(preset: DashboardFilterPreset): string {
  switch (preset) {
    case "day":
      return "Resumen del día";
    case "week":
      return "Resumen de la semana";
    case "year":
      return "Resumen del año";
    case "custom":
      return "Resumen del período";
    case "month":
    default:
      return "Resumen del mes";
  }
}

export function getDashboardPresetLabel(preset: DashboardFilterPreset): string {
  switch (preset) {
    case "day":
      return "Día";
    case "week":
      return "Semana";
    case "year":
      return "Año";
    case "custom":
      return "Rango personalizado";
    case "month":
    default:
      return "Mes";
  }
}

export function formatPercent(value: number): string {
  const absoluteValue = Math.abs(value);

  return new Intl.NumberFormat("es-AR", {
    maximumFractionDigits: absoluteValue >= 100 ? 0 : 1,
  }).format(absoluteValue);
}
