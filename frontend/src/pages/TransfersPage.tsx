import { motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
import { transactionsApi } from "../api/transactions";
import { transfersApi } from "../api/transfers";
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
import { formatCurrency, formatDate } from "../lib/format";
import type { Account, Transaction } from "../types/api";

export function TransfersPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [transferRows, setTransferRows] = useState<Transaction[]>([]);
  const [fromAccountId, setFromAccountId] = useState("");
  const [toAccountId, setToAccountId] = useState("");
  const [amount, setAmount] = useState("");
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

  const loadPage = async () => {
    setLoading(true);
    setError(null);

    try {
      const [accountsData, transfersData] = await Promise.all([
        accountsApi.list(),
        transactionsApi.list({ category: "transfer" }),
      ]);

      setAccounts(accountsData);
      setTransferRows(transfersData);

      if (accountsData[0]) {
        setFromAccountId((current) => current || String(accountsData[0].id));
      }

      if (accountsData[1]) {
        setToAccountId((current) => current || String(accountsData[1].id));
      }
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudieron cargar las transferencias");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadPage();
  }, []);

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      await transfersApi.create({
        from_account_id: Number(fromAccountId),
        to_account_id: Number(toAccountId),
        amount: Number(amount),
        description,
      });

      setAmount("");
      setDescription("");
      setSuccess("Transferencia registrada correctamente");
      await loadPage();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo registrar la transferencia");
    } finally {
      setSubmitting(false);
    }
  };

  const orderedTransfers = useMemo(
    () =>
      [...transferRows].sort((left, right) => Date.parse(right.created_at) - Date.parse(left.created_at)),
    [transferRows],
  );

  return (
    <div className="page-stack">
      <Reveal>
        <section className="page-header">
          <div>
            <span className="eyebrow">Transferencias</span>
            <h1>Mover saldo entre tus cuentas</h1>
            <p>Pasá dinero entre tus cuentas y seguí el registro de cada movimiento generado por la transferencia.</p>
          </div>
        </section>
      </Reveal>

      <StaggerGroup className="content-grid">
        <motion.article className="panel" variants={fadeLeft}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Nueva transferencia</span>
              <h2>Enviar entre cuentas</h2>
            </div>
          </div>

          <form className="stack-form" onSubmit={handleSubmit}>
            <label className="field">
              <span>Cuenta origen</span>
              <AnimatedSelect
                value={fromAccountId}
                options={accountOptions}
                onChange={setFromAccountId}
                placeholder={accounts.length < 2 ? "Necesitás al menos dos cuentas" : "Elegí una cuenta"}
                disabled={accounts.length < 2}
                ariaLabel="Cuenta origen"
              />
            </label>

            <label className="field">
              <span>Cuenta destino</span>
              <AnimatedSelect
                value={toAccountId}
                options={accountOptions}
                onChange={setToAccountId}
                placeholder={accounts.length < 2 ? "Necesitás al menos dos cuentas" : "Elegí una cuenta"}
                disabled={accounts.length < 2}
                ariaLabel="Cuenta destino"
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
                required
              />
            </label>

            <label className="field">
              <span>Descripción</span>
              <textarea
                rows={4}
                value={description}
                onChange={(event) => setDescription(event.target.value)}
                placeholder="Ejemplo: dinero para gastos del mes"
              />
            </label>

            {accounts.length < 2 ? (
              <motion.p
                className="feedback feedback--warning"
                initial={{ opacity: 0, y: 10 }}
                animate={{ opacity: 1, y: 0 }}
                transition={{ duration: 0.24, ease: [0.22, 1, 0.36, 1] }}
              >
                Necesitás al menos dos cuentas para poder transferir.
              </motion.p>
            ) : null}
            <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>
            <PresenceMessage className="feedback feedback--success">{success}</PresenceMessage>

            <button
              type="submit"
              className="primary-button"
              disabled={submitting || accounts.length < 2}
            >
              {submitting ? "Enviando..." : "Registrar transferencia"}
            </button>
          </form>
        </motion.article>

        <motion.article className="panel" variants={fadeRight}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Historial</span>
              <h2>Transferencias recientes</h2>
            </div>
            <button type="button" className="ghost-button" onClick={() => void loadPage()}>
              Recargar
            </button>
          </div>

          {loading ? (
            <motion.p
              className="feedback"
              initial={{ opacity: 0, y: 10 }}
              animate={{ opacity: 1, y: 0 }}
              transition={{ duration: 0.24, ease: [0.22, 1, 0.36, 1] }}
            >
              Cargando transferencias...
            </motion.p>
          ) : orderedTransfers.length === 0 ? (
            <EmptyState
              title="No hay transferencias registradas"
              description="Cuando hagas la primera, vas a verla reflejada en las cuentas involucradas."
            />
          ) : (
            <StaggerGroup className="stack-list">
              {orderedTransfers.map((transaction) => (
                <motion.div
                  key={transaction.id}
                  className="list-row list-row--transaction"
                  variants={fadeUp}
                  whileHover={cardHover}
                  layout
                >
                  <div>
                    <strong>{transaction.title}</strong>
                    <p>{accountMap.get(transaction.account_id) ?? `Cuenta ${transaction.account_id}`}</p>
                    {transaction.transfer_id ? <small>ID: {transaction.transfer_id}</small> : null}
                  </div>

                  <div className="list-row__meta">
                    <strong>{formatCurrency(transaction.amount)}</strong>
                    <span>{formatDate(transaction.created_at)}</span>
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
