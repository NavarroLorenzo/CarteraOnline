import { motion } from "framer-motion";
import { Link } from "react-router-dom";
import { formatCurrency, formatDate, formatTypeLabel } from "../../lib/format";
import type { DashboardTransactionItem } from "../../types/api";
import { EmptyState } from "../ui/EmptyState";
import { Reveal, StaggerGroup, cardHover, fadeUp } from "../ui/animation";

type DashboardRecentTransactionsProps = {
  transactions: DashboardTransactionItem[];
  loading: boolean;
};

export function DashboardRecentTransactions({ transactions, loading }: DashboardRecentTransactionsProps) {
  return (
    <Reveal onView={false}>
      <section className="panel dashboard-panel dashboard-panel--recent">
        <div className="panel-heading">
          <div>
            <span className="eyebrow">Actividad reciente</span>
            <h2>Últimos movimientos</h2>
            <p>Este bloque se mantiene como una lectura rápida e independiente de la actividad más reciente.</p>
          </div>
          <Link to="/transactions" className="ghost-button">
            Ver todas las transacciones
          </Link>
        </div>

        {loading ? (
          <p className="feedback">Cargando información...</p>
        ) : transactions.length === 0 ? (
          <EmptyState
            title="No hay movimientos para mostrar"
            description="Cuando registres transacciones, esta vista rápida mostrará la actividad más reciente."
          />
        ) : (
          <StaggerGroup className="stack-list" onView={false}>
            {transactions.map((transaction) => (
              <motion.div
                key={transaction.id}
                className={`list-row list-row--transaction list-row--${transaction.type}`}
                variants={fadeUp}
                whileHover={cardHover}
              >
                <div>
                  <div className="list-row__title">
                    <strong>{transaction.title}</strong>
                    <span className={`category-badge category-badge--${transaction.type}`}>
                      {transaction.category_label}
                    </span>
                  </div>
                  <p>
                    {transaction.account_name || `Cuenta ${transaction.account_id}`} ·{" "}
                    {formatTypeLabel(transaction.type)}
                  </p>
                  {transaction.description ? <small>{transaction.description}</small> : null}
                </div>

                <div className="list-row__meta">
                  <strong>{formatCurrency(transaction.amount)}</strong>
                  <span>{formatDate(transaction.created_at)}</span>
                </div>
              </motion.div>
            ))}
          </StaggerGroup>
        )}
      </section>
    </Reveal>
  );
}
