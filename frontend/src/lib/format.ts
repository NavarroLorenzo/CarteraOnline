const currencyFormatter = new Intl.NumberFormat("es-AR", {
  style: "currency",
  currency: "ARS",
  maximumFractionDigits: 2,
});

const dateFormatter = new Intl.DateTimeFormat("es-AR", {
  dateStyle: "medium",
  timeStyle: "short",
});

export function formatCurrency(value: number): string {
  return currencyFormatter.format(value);
}

export function formatDate(value?: string | Date | null): string {
  if (!value) {
    return "Fecha no disponible";
  }

  const parsed = value instanceof Date ? value : new Date(value);

  if (Number.isNaN(parsed.getTime())) {
    return "Fecha no disponible";
  }

  return dateFormatter.format(parsed);
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
