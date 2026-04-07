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

export function formatDate(value: string): string {
  return dateFormatter.format(new Date(value));
}

export function formatTypeLabel(value: string): string {
  if (value === "income") {
    return "Ingreso";
  }

  if (value === "expense") {
    return "Gasto";
  }

  return value.split("_").join(" ");
}

export function formatAccountTypeLabel(value: string): string {
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

  return value.split("_").join(" ");
}

export function formatCategoryKeyLabel(value: string): string {
  const normalized = value.trim().toLowerCase();

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
