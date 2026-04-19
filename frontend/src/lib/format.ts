const currencyFormatter = new Intl.NumberFormat("es-AR", {
  style: "currency",
  currency: "ARS",
  maximumFractionDigits: 2,
});

const shortDateFormatter = new Intl.DateTimeFormat("es-AR", {
  day: "numeric",
  month: "short",
  year: "numeric",
});

const timeFormatter = new Intl.DateTimeFormat("es-AR", {
  hour: "2-digit",
  minute: "2-digit",
  hour12: false,
});

export function formatCurrency(value: number): string {
  return currencyFormatter.format(value);
}

export function formatDate(value?: string | Date | null): string {
  const parsed = parseDateValue(value);

  if (!parsed) {
    return "Fecha no disponible";
  }

  if (isSameDay(parsed, new Date())) {
    return "Hoy";
  }

  const yesterday = new Date();
  yesterday.setDate(yesterday.getDate() - 1);

  if (isSameDay(parsed, yesterday)) {
    return "Ayer";
  }

  return formatShortDate(parsed);
}

export function formatDateTime(value?: string | Date | null): string {
  const parsed = parseDateValue(value);

  if (!parsed) {
    return "Fecha no disponible";
  }

  if (!hasExplicitTime(value)) {
    return formatDate(parsed);
  }

  return `${formatDate(parsed)}, ${timeFormatter.format(parsed)} hs`;
}

export function formatTypeLabel(value?: string | null): string {
  if (value === "income") {
    return "Ingreso";
  }

  if (value === "expense") {
    return "Gasto";
  }

  if (!value) {
    return "Sin tipo";
  }

  return value.split("_").join(" ");
}

export function formatAccountTypeLabel(value?: string | null): string {
  if (value === "cash") {
    return "Efectivo";
  }

  if (value === "bank") {
    return "Banco";
  }

  if (value === "virtual_wallet") {
    return "Billetera virtual";
  }

  if (value === "credit_card") {
    return "Tarjeta de crédito";
  }

  if (value === "savings") {
    return "Ahorro";
  }

  if (!value) {
    return "Sin tipo";
  }

  return value.split("_").join(" ");
}

export function formatCategoryKeyLabel(value?: string | null): string {
  const normalized = value?.trim().toLowerCase() ?? "";

  if (normalized === "inversion") {
    return "Inversión";
  }

  if (normalized === "suscripciones") {
    return "Suscripciones";
  }

  if (!normalized) {
    return "Otros";
  }

  return normalized.charAt(0).toUpperCase() + normalized.slice(1);
}

function parseDateValue(value?: string | Date | null): Date | null {
  if (!value) {
    return null;
  }

  if (value instanceof Date) {
    return Number.isNaN(value.getTime()) ? null : value;
  }

  if (/^\d{4}-\d{2}-\d{2}$/.test(value)) {
    const [year, month, day] = value.split("-").map(Number);
    const parsedInputDate = new Date(year, month - 1, day);

    if (
      Number.isNaN(parsedInputDate.getTime()) ||
      parsedInputDate.getFullYear() !== year ||
      parsedInputDate.getMonth() !== month - 1 ||
      parsedInputDate.getDate() !== day
    ) {
      return null;
    }

    return parsedInputDate;
  }

  const parsed = new Date(value);
  return Number.isNaN(parsed.getTime()) ? null : parsed;
}

function isSameDay(left: Date, right: Date): boolean {
  return (
    left.getFullYear() === right.getFullYear() &&
    left.getMonth() === right.getMonth() &&
    left.getDate() === right.getDate()
  );
}

function hasExplicitTime(value?: string | Date | null): boolean {
  if (!value || value instanceof Date) {
    return Boolean(value);
  }

  return /(?:T|\s)\d{1,2}:\d{2}/.test(value);
}

function formatShortDate(value: Date): string {
  const parts = shortDateFormatter.formatToParts(value);
  const day = parts.find((part) => part.type === "day")?.value;
  const month = parts.find((part) => part.type === "month")?.value.replace(/\./g, "");
  const year = parts.find((part) => part.type === "year")?.value;

  if (!day || !month || !year) {
    return shortDateFormatter.format(value).replace(/\./g, "").replace(/\sde\s/g, " ");
  }

  return `${day} ${month} ${year}`;
}
