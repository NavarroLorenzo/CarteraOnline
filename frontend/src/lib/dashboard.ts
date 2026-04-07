export type DashboardFilterPreset = "today" | "last_7_days" | "this_month" | "custom";
export type DashboardMode = "simple" | "advanced";

export type DashboardDateRange = {
  date_from: string;
  date_to: string;
};

const monthFormatter = new Intl.DateTimeFormat("es-AR", {
  month: "long",
  year: "numeric",
});

const dayFormatter = new Intl.DateTimeFormat("es-AR", {
  day: "numeric",
  month: "short",
});

export function getDashboardRange(
  preset: DashboardFilterPreset,
  customRange: DashboardDateRange,
  now = new Date(),
): DashboardDateRange | null {
  if (preset === "custom") {
    if (!customRange.date_from || !customRange.date_to) {
      return null;
    }

    return customRange;
  }

  const today = startOfDay(now);

  if (preset === "today") {
    return {
      date_from: toInputDate(today),
      date_to: toInputDate(today),
    };
  }

  if (preset === "last_7_days") {
    const dateFrom = new Date(today);
    dateFrom.setDate(today.getDate() - 6);

    return {
      date_from: toInputDate(dateFrom),
      date_to: toInputDate(today),
    };
  }

  const monthStart = new Date(today.getFullYear(), today.getMonth(), 1);
  return {
    date_from: toInputDate(monthStart),
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

export function parseInputDate(value: string): Date {
  const [year, month, day] = value.split("-").map(Number);
  return new Date(year, month - 1, day);
}

export function formatDaySummaryLabel(value: string): string {
  return dayFormatter.format(parseInputDate(value));
}

export function formatMonthSummaryLabel(value: string): string {
  const monthDate = parseInputDate(value);
  const formatted = monthFormatter.format(monthDate);
  return formatted.charAt(0).toUpperCase() + formatted.slice(1);
}

export function formatPercent(value: number): string {
  const absoluteValue = Math.abs(value);

  return new Intl.NumberFormat("es-AR", {
    maximumFractionDigits: absoluteValue >= 100 ? 0 : 1,
  }).format(absoluteValue);
}
