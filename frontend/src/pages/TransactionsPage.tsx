import { motion } from "framer-motion";
import { useEffect, useMemo, useRef, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
import { transactionsApi, type TransactionFilters } from "../api/transactions";
import { EmptyState } from "../components/ui/EmptyState";
import { AnimatedSelect } from "../components/ui/AnimatedSelect";
import { SelectionGrid } from "../components/ui/SelectionGrid";
import {
  PresenceMessage,
  Reveal,
  StaggerGroup,
  cardHover,
  fadeLeft,
  fadeRight,
  fadeUp,
} from "../components/ui/animation";
import {
  formatAccountTypeLabel,
  formatCategoryKeyLabel,
  formatCurrency,
  formatDate,
  formatTypeLabel,
} from "../lib/format";
import type { Account, Transaction, TransactionCategory, TransactionType } from "../types/api";

const transactionTypeOptions: Array<{ value: TransactionType; label: string; description: string }> = [
  { value: "income", label: "Ingreso", description: "Entradas de dinero" },
  { value: "expense", label: "Gasto", description: "Salidas de dinero" },
];

export function TransactionsPage() {
  const filtersInitializedRef = useRef(false);
  const latestRequestRef = useRef(0);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [categories, setCategories] = useState<TransactionCategory[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [filters, setFilters] = useState<TransactionFilters>({
    type: "",
    category: "",
    date_from: "",
    date_to: "",
  });

  const [title, setTitle] = useState("");
  const [amount, setAmount] = useState("");
  const [type, setType] = useState<TransactionType>("income");
  const [accountId, setAccountId] = useState("");
  const [category, setCategory] = useState("");
  const [description, setDescription] = useState("");

  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const accountMap = useMemo(
    () => new Map(accounts.map((account) => [account.id, account.name])),
    [accounts],
  );
  const categoryMap = useMemo(
    () => new Map(categories.map((transactionCategory) => [transactionCategory.key, transactionCategory.label])),
    [categories],
  );

  const activeAccounts = useMemo(
    () => accounts.filter((account) => account.is_active),
    [accounts],
  );

  const accountGridOptions = useMemo(
    () =>
      activeAccounts.map((account) => ({
        value: String(account.id),
        label: account.name,
        description: formatAccountTypeLabel(account.type),
      })),
    [activeAccounts],
  );

  const accountOptions = useMemo(
    () => [{ value: "", label: "Todas" }, ...accounts.map((account) => ({ value: String(account.id), label: account.name }))],
    [accounts],
  );

  const visibleCategories = useMemo(
    () => categories.filter((transactionCategory) => transactionCategory.allowed_types.includes(type)),
    [categories, type],
  );

  const categoryGridOptions = useMemo(
    () =>
      visibleCategories.map((transactionCategory) => ({
        value: transactionCategory.key,
        label: transactionCategory.label,
      })),
    [visibleCategories],
  );

  const filterCategoryOptions = useMemo(
    () => [{ value: "", label: "Todas" }, ...categories.map((transactionCategory) => ({
      value: transactionCategory.key,
      label: transactionCategory.label,
    }))],
    [categories],
  );

  const filterTypeOptions = useMemo(
    () => [{ value: "", label: "Todos" }, ...transactionTypeOptions.map(({ value, label }) => ({ value, label }))],
    [],
  );

  const loadInitialData = async () => {
    setLoading(true);
    setError(null);
    const requestId = latestRequestRef.current + 1;
    latestRequestRef.current = requestId;

    try {
      const [accountsData, categoriesData, transactionsData] = await Promise.all([
        accountsApi.list(),
        transactionsApi.getCategories(),
        transactionsApi.list(filters),
      ]);

      if (requestId !== latestRequestRef.current) {
        return;
      }

      setAccounts(accountsData);
      setCategories(categoriesData);
      setTransactions(transactionsData);
      filtersInitializedRef.current = true;
    } catch (err) {
      if (requestId === latestRequestRef.current) {
        setError(err instanceof ApiError ? err.message : "No se pudieron cargar las transacciones");
      }
    } finally {
      if (requestId === latestRequestRef.current) {
        setLoading(false);
      }
    }
  };

  useEffect(() => {
    void loadInitialData();
  }, []);

  useEffect(() => {
    if (!category) {
      return;
    }

    const categoryStillVisible = visibleCategories.some((transactionCategory) => transactionCategory.key === category);
    if (!categoryStillVisible) {
      setCategory("");
    }
  }, [category, visibleCategories]);

  const refreshTransactions = async (nextFilters = filters) => {
    setLoading(true);
    setError(null);
    const requestId = latestRequestRef.current + 1;
    latestRequestRef.current = requestId;

    try {
      const data = await transactionsApi.list(nextFilters);

      if (requestId !== latestRequestRef.current) {
        return;
      }

      setTransactions(data);
    } catch (err) {
      if (requestId === latestRequestRef.current) {
        setError(err instanceof ApiError ? err.message : "No se pudieron cargar las transacciones");
      }
    } finally {
      if (requestId === latestRequestRef.current) {
        setLoading(false);
      }
    }
  };

  useEffect(() => {
    if (!filtersInitializedRef.current) {
      return;
    }

    void refreshTransactions(filters);
  }, [filters]);

  const handleCreate = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

    if (!title.trim() || !amount || !type || !accountId || !category) {
      setError("Completá nombre, monto, tipo, cuenta y categoría para registrar la transacción.");
      setSubmitting(false);
      return;
    }

    try {
      await transactionsApi.create({
        title,
        amount: Number(amount),
        type,
        account_id: Number(accountId),
        category,
        description,
      });

      setTitle("");
      setAmount("");
      setDescription("");
      setSuccess("Transacción creada correctamente");
      await refreshTransactions();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo crear la transacción");
    } finally {
      setSubmitting(false);
    }
  };

  const handleDelete = async (transactionId: number) => {
    setError(null);
    setSuccess(null);

    try {
      await transactionsApi.delete(transactionId);
      setSuccess("Transacción eliminada correctamente");
      await refreshTransactions();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo eliminar la transacción");
    }
  };

  const orderedTransactions = useMemo(
    () =>
      [...transactions].sort((left, right) => Date.parse(right.created_at) - Date.parse(left.created_at)),
    [transactions],
  );

  return (
    <div className="page-stack">
      <StaggerGroup className="content-grid content-grid--wide" onView={false}>
        <motion.article className="panel" variants={fadeLeft}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Nueva transacción</span>
              <h2>Carga rápida</h2>
              <p>Solo hace falta nombre, monto, tipo, cuenta y categoría.</p>
            </div>
          </div>

          <form className="stack-form transaction-quick-form" onSubmit={handleCreate}>
            <label className="field">
              <span>Nombre</span>
              <input
                type="text"
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                placeholder="Supermercado del sábado"
                required
              />
            </label>

            <label className="field">
              <span>Monto</span>
              <input
                type="number"
                min="0.01"
                step="0.01"
                value={amount}
                onChange={(event) => setAmount(event.target.value)}
                placeholder="0,00"
                required
              />
            </label>

            <SelectionGrid
              label="Tipo"
              value={type}
              options={transactionTypeOptions}
              onChange={(nextType) => setType(nextType as TransactionType)}
              columns="compact"
            />

            <SelectionGrid
              label="Categoría"
              value={category}
              options={categoryGridOptions}
              onChange={setCategory}
              columns="wide"
              emptyMessage="No hay categorías disponibles para este tipo."
            />

            <SelectionGrid
              label="Cuenta"
              value={accountId}
              options={accountGridOptions}
              onChange={setAccountId}
              columns="regular"
              emptyMessage="Primero necesitás una cuenta activa para registrar movimientos."
              disabled={accountGridOptions.length === 0}
            />

            <label className="field">
              <span>Descripción</span>
              <textarea
                rows={3}
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                placeholder="Dato opcional para recordar el contexto"
              />
            </label>

            <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>
            <PresenceMessage className="feedback feedback--success">{success}</PresenceMessage>

            <button
              type="submit"
              className="primary-button"
              disabled={submitting || activeAccounts.length === 0}
            >
              {submitting ? "Guardando..." : "Guardar movimiento"}
            </button>
          </form>
        </motion.article>

        <motion.article className="panel panel--history panel--history-transactions" variants={fadeRight}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Historial</span>
              <h2>Movimientos y filtros</h2>
            </div>
          </div>

          <form className="stack-form stack-form--compact">
            <div className="field-grid">
              <label className="field">
                <span>Cuenta</span>
                <AnimatedSelect
                  value={filters.account_id ? String(filters.account_id) : ""}
                  options={accountOptions}
                  onChange={(value) =>
                    setFilters((current) => ({
                      ...current,
                      account_id: value ? Number(value) : undefined,
                    }))
                  }
                  ariaLabel="Filtrar por cuenta"
                />
              </label>

              <label className="field">
                <span>Tipo</span>
                <AnimatedSelect
                  value={filters.type ?? ""}
                  options={filterTypeOptions}
                  onChange={(value) =>
                    setFilters((current) => ({
                      ...current,
                      type: value as TransactionType | "",
                    }))
                  }
                  ariaLabel="Filtrar por tipo"
                />
              </label>
            </div>

            <div className="field-grid field-grid--three">
              <label className="field">
                <span>Categoría</span>
                <AnimatedSelect
                  value={filters.category ?? ""}
                  options={filterCategoryOptions}
                  onChange={(value) =>
                    setFilters((current) => ({
                      ...current,
                      category: value,
                    }))
                  }
                  ariaLabel="Filtrar por categoría"
                />
              </label>

              <label className="field">
                <span>Desde</span>
                <input
                  type="date"
                  value={filters.date_from ?? ""}
                  onChange={(event) =>
                    setFilters((current) => ({
                      ...current,
                      date_from: event.target.value,
                    }))
                  }
                />
              </label>

              <label className="field">
                <span>Hasta</span>
                <input
                  type="date"
                  value={filters.date_to ?? ""}
                  onChange={(event) =>
                    setFilters((current) => ({
                      ...current,
                      date_to: event.target.value,
                    }))
                  }
                />
              </label>
            </div>

            <div className="action-row">
              <button
                type="button"
                className="ghost-button"
                onClick={() => {
                  const cleared: TransactionFilters = {
                    type: "",
                    category: "",
                    date_from: "",
                    date_to: "",
                  };
                  setFilters(cleared);
                }}
              >
                Limpiar filtros
              </button>
            </div>
          </form>

          <div className="history-panel__content">
            {loading ? (
              <motion.p
                className="feedback"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.24, ease: [0.22, 1, 0.36, 1] }}
              >
                Cargando movimientos...
              </motion.p>
            ) : orderedTransactions.length === 0 ? (
              <EmptyState
                title="No se encontraron movimientos"
                description="Ajustá los filtros o registrá un movimiento para verlo en este historial."
              />
            ) : (
              <StaggerGroup className="stack-list scrollable-list scrollable-list--transactions" onView={false}>
                {orderedTransactions.map((transaction) => {
                  const categoryLabel =
                    transaction.category_label ??
                    categoryMap.get(transaction.category) ??
                    formatCategoryKeyLabel(transaction.category);

                  return (
                    <motion.div
                      key={transaction.id}
                      className={`list-row list-row--transaction list-row--${transaction.type}`}
                      variants={fadeUp}
                      whileHover={cardHover}
                      layout
                    >
                      <div>
                        <div className="list-row__title">
                          <strong>{transaction.title}</strong>
                          <span className={`category-badge category-badge--${transaction.type}`}>
                            {categoryLabel}
                          </span>
                        </div>
                        <p>
                          {accountMap.get(transaction.account_id) ?? `Cuenta ${transaction.account_id}`} ·{" "}
                          {formatTypeLabel(transaction.type)}
                        </p>
                        {transaction.description ? <small>{transaction.description}</small> : null}
                      </div>

                      <div className="list-row__meta">
                        <strong>{formatCurrency(transaction.amount)}</strong>
                        <span>{formatDate(transaction.created_at)}</span>
                        <button
                          type="button"
                          className="danger-link"
                          onClick={() => void handleDelete(transaction.id)}
                        >
                          Eliminar
                        </button>
                      </div>
                    </motion.div>
                  );
                })}
              </StaggerGroup>
            )}
          </div>
        </motion.article>
      </StaggerGroup>
    </div>
  );
}
