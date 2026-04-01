import { motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
import { transactionsApi, type TransactionFilters } from "../api/transactions";
import { EmptyState } from "../components/ui/EmptyState";
import { AnimatedSelect } from "../components/ui/AnimatedSelect";
import {
  PresenceMessage,
  Reveal,
  StaggerGroup,
  cardHover,
  fadeLeft,
  fadeRight,
  fadeUp,
} from "../components/ui/animation";
import { formatCurrency, formatDate, formatTypeLabel } from "../lib/format";
import type { Account, Transaction, TransactionType } from "../types/api";

const transactionTypeOptions = [
  { value: "income", label: "Ingreso" },
  { value: "expense", label: "Gasto" },
];

export function TransactionsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
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
  const accountOptions = useMemo(
    () => accounts.map((account) => ({ value: String(account.id), label: account.name })),
    [accounts],
  );
  const filterAccountOptions = useMemo(
    () => [{ value: "", label: "Todas" }, ...accountOptions],
    [accountOptions],
  );
  const filterTypeOptions = useMemo(
    () => [{ value: "", label: "Todos" }, ...transactionTypeOptions],
    [],
  );

  const loadInitialData = async () => {
    setLoading(true);
    setError(null);

    try {
      const [accountsData, transactionsData] = await Promise.all([
        accountsApi.list(),
        transactionsApi.list(filters),
      ]);

      setAccounts(accountsData);
      setTransactions(transactionsData);

      if (!accountId && accountsData[0]) {
        setAccountId(String(accountsData[0].id));
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudieron cargar las transacciones");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadInitialData();
  }, []);

  const refreshTransactions = async (nextFilters = filters) => {
    setLoading(true);
    setError(null);

    try {
      const data = await transactionsApi.list(nextFilters);
      setTransactions(data);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudieron cargar las transacciones");
    } finally {
      setLoading(false);
    }
  };

  const handleCreate = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

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
      setCategory("");
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

  const handleFilterSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    await refreshTransactions(filters);
  };

  const orderedTransactions = useMemo(
    () =>
      [...transactions].sort((left, right) => Date.parse(right.created_at) - Date.parse(left.created_at)),
    [transactions],
  );

  return (
    <div className="page-stack">
      <Reveal>
        <section className="page-header">
          <div>
            <span className="eyebrow">Transacciones</span>
            <h1>Registrá y consultá tus movimientos</h1>
            <p>Guardá ingresos y gastos, filtralos por cuenta o fecha y revisá tu historial cuando lo necesites.</p>
          </div>
        </section>
      </Reveal>

      <StaggerGroup className="content-grid content-grid--wide">
        <motion.article className="panel" variants={fadeLeft}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Nueva transacción</span>
              <h2>Cargar movimiento</h2>
            </div>
          </div>

          <form className="stack-form" onSubmit={handleCreate}>
            <label className="field">
              <span>Título</span>
              <input
                type="text"
                value={title}
                onChange={(event) => setTitle(event.target.value)}
                placeholder="Sueldo de marzo"
                required
              />
            </label>

            <div className="field-grid">
              <label className="field">
                <span>Monto</span>
                <input
                  type="number"
                  min="0.01"
                  step="0.01"
                  value={amount}
                  onChange={(event) => setAmount(event.target.value)}
                  required
                />
              </label>

              <label className="field">
                <span>Tipo</span>
                <AnimatedSelect
                  value={type}
                  options={transactionTypeOptions}
                  onChange={(value) => setType(value as TransactionType)}
                  ariaLabel="Tipo de transacción"
                />
              </label>
            </div>

            <label className="field">
              <span>Cuenta</span>
              <AnimatedSelect
                value={accountId}
                options={accountOptions}
                onChange={setAccountId}
                placeholder={accounts.length === 0 ? "Primero creá una cuenta" : "Seleccioná una cuenta"}
                disabled={accounts.length === 0}
                ariaLabel="Cuenta"
              />
            </label>

            <label className="field">
              <span>Categoría</span>
              <input
                type="text"
                value={category}
                onChange={(event) => setCategory(event.target.value)}
                placeholder="supermercado"
              />
            </label>

            <label className="field">
              <span>Descripción</span>
              <textarea
                rows={4}
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                placeholder="Observaciones opcionales"
              />
            </label>

            <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>
            <PresenceMessage className="feedback feedback--success">{success}</PresenceMessage>

            <button
              type="submit"
              className="primary-button"
              disabled={submitting || accounts.length === 0}
            >
              {submitting ? "Guardando..." : "Guardar movimiento"}
            </button>
          </form>
        </motion.article>

        <motion.article className="panel" variants={fadeRight}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Historial</span>
              <h2>Movimientos y filtros</h2>
            </div>
          </div>

          <form className="stack-form stack-form--compact" onSubmit={handleFilterSubmit}>
            <div className="field-grid">
              <label className="field">
                <span>Cuenta</span>
                <AnimatedSelect
                  value={filters.account_id ? String(filters.account_id) : ""}
                  options={filterAccountOptions}
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

            <div className="field-grid">
              <label className="field">
                <span>Categoría</span>
                <input
                  type="text"
                  value={filters.category ?? ""}
                  onChange={(event) =>
                    setFilters((current) => ({
                      ...current,
                      category: event.target.value,
                    }))
                  }
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
              <button type="submit" className="ghost-button">
                Aplicar filtros
              </button>
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
                  void refreshTransactions(cleared);
                }}
              >
                Limpiar
              </button>
            </div>
          </form>

          {loading ? (
            <p className="feedback">Cargando movimientos...</p>
          ) : orderedTransactions.length === 0 ? (
            <EmptyState
              title="No hay resultados"
              description="Probá ajustando los filtros o registrando un movimiento nuevo."
            />
          ) : (
            <StaggerGroup className="stack-list">
              {orderedTransactions.map((transaction) => (
                <motion.div
                  key={transaction.id}
                  className="list-row list-row--transaction"
                  variants={fadeUp}
                  whileHover={cardHover}
                  layout
                >
                  <div>
                    <strong>{transaction.title}</strong>
                    <p>
                      {accountMap.get(transaction.account_id) ?? `Cuenta ${transaction.account_id}`} ·{" "}
                      {formatTypeLabel(transaction.type)} · {transaction.category || "sin categoría"}
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
              ))}
            </StaggerGroup>
          )}
        </motion.article>
      </StaggerGroup>
    </div>
  );
}
