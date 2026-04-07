import { motion } from "framer-motion";
import { useEffect, useMemo, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
import { transactionsApi } from "../api/transactions";
import { MetricCard } from "../components/ui/MetricCard";
import { EmptyState } from "../components/ui/EmptyState";
import {
  PresenceMessage,
  Reveal,
  StaggerGroup,
  cardHover,
  fadeLeft,
  fadeRight,
  fadeUp,
} from "../components/ui/animation";
import { formatAccountTypeLabel, formatCurrency, formatDate, formatTypeLabel } from "../lib/format";
import type { Account, AccountBalance, Transaction, TransactionSummary } from "../types/api";

export function DashboardPage() {
  const [summary, setSummary] = useState<TransactionSummary | null>(null);
  const [balances, setBalances] = useState<AccountBalance[]>([]);
  const [totalBalance, setTotalBalance] = useState(0);
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [transactions, setTransactions] = useState<Transaction[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    let cancelled = false;

    const loadDashboard = async () => {
      setLoading(true);
      setError(null);

      try {
        const [summaryData, balanceData, accountsData, transactionsData] = await Promise.all([
          transactionsApi.getSummary(),
          transactionsApi.getBalanceByAccount(),
          accountsApi.list(),
          transactionsApi.list(),
        ]);

        if (cancelled) {
          return;
        }

        setSummary(summaryData);
        setBalances(balanceData.accounts);
        setTotalBalance(balanceData.total);
        setAccounts(accountsData);
        setTransactions(transactionsData);
      } catch (err) {
        if (!cancelled) {
          setError(err instanceof ApiError ? err.message : "No se pudo cargar el dashboard");
        }
      } finally {
        if (!cancelled) {
          setLoading(false);
        }
      }
    };

    void loadDashboard();

    return () => {
      cancelled = true;
    };
  }, []);

  const recentTransactions = useMemo(
    () =>
      [...transactions]
        .sort((left, right) => Date.parse(right.created_at) - Date.parse(left.created_at))
        .slice(0, 6),
    [transactions],
  );

  return (
    <div className="page-stack">
      <Reveal onView={false}>
        <section className="hero-card">
          <div>
            <span className="eyebrow">Dashboard</span>
            <h1>Tu panorama financiero de hoy</h1>
            <p>
              Revisá el balance total, el resultado de tus movimientos y la actividad reciente de cada cuenta.
            </p>
          </div>
        </section>
      </Reveal>

      <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>

      <StaggerGroup className="metrics-grid" onView={false}>
        <MetricCard
          label="Balance total"
          value={loading ? "Cargando..." : formatCurrency(totalBalance)}
          tone="accent"
        />
        <MetricCard
          label="Ingresos"
          value={summary ? formatCurrency(summary.income_total) : "Cargando..."}
          tone="positive"
        />
        <MetricCard
          label="Gastos"
          value={summary ? formatCurrency(summary.expense_total) : "Cargando..."}
          tone="negative"
        />
        <MetricCard
          label="Movimientos"
          value={summary ? String(summary.transactions_count) : "Cargando..."}
        />
      </StaggerGroup>

      <StaggerGroup className="dashboard-grid" onView={false}>
        <motion.article className="panel" variants={fadeLeft}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Cuentas</span>
              <h2>Balance por cuenta</h2>
            </div>
          </div>

          {balances.length === 0 ? (
            <EmptyState
              title="Todavía no hay cuentas"
              description="Creá tu primera cuenta desde la sección Cuentas para empezar."
            />
          ) : (
            <StaggerGroup className="stack-list" onView={false}>
              {balances.map((account) => (
                <motion.div
                  key={account.id}
                  className="list-row"
                  variants={fadeUp}
                  whileHover={cardHover}
                >
                  <div>
                    <strong>{account.name}</strong>
                    <p>Cuenta #{account.id}</p>
                  </div>
                  <strong>{formatCurrency(account.balance)}</strong>
                </motion.div>
              ))}
            </StaggerGroup>
          )}
        </motion.article>

        <motion.article className="panel" variants={fadeRight}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Actividad</span>
              <h2>Últimos movimientos</h2>
            </div>
          </div>

          {recentTransactions.length === 0 ? (
            <EmptyState
              title="No hay transacciones cargadas"
              description="Cuando registres ingresos o gastos, van a aparecer acá."
            />
          ) : (
            <StaggerGroup className="stack-list" onView={false}>
              {recentTransactions.map((transaction) => (
                <motion.div
                  key={transaction.id}
                  className="list-row list-row--transaction"
                  variants={fadeUp}
                  whileHover={cardHover}
                >
                  <div>
                    <strong>{transaction.title}</strong>
                    <p>
                      {formatTypeLabel(transaction.type)} · {transaction.category || "sin categoría"}
                    </p>
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

      <Reveal onView={false}>
        <section className="panel">
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Vista general</span>
              <h2>Cuentas disponibles</h2>
            </div>
          </div>

          {accounts.length === 0 ? (
            <EmptyState
              title="No hay cuentas para mostrar"
              description="Creá una cuenta nueva para empezar a ver tu cartera resumida acá."
            />
          ) : (
            <StaggerGroup className="pill-list" onView={false}>
              {accounts.map((account) => (
                <motion.span
                  key={account.id}
                  className="pill"
                  variants={fadeUp}
                  whileHover={cardHover}
                >
                  {account.name} · {formatAccountTypeLabel(account.type)}
                </motion.span>
              ))}
            </StaggerGroup>
          )}
        </section>
      </Reveal>
    </div>
  );
}
