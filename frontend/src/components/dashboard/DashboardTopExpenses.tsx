import { motion } from "framer-motion";
import { formatCurrency, formatDate } from "../../lib/format";
import type { DashboardTransactionItem } from "../../types/api";
import { EmptyState } from "../ui/EmptyState";
import { StaggerGroup, cardHover, fadeRight, fadeUp } from "../ui/animation";

type DashboardTopExpensesProps = {
  transactions: DashboardTransactionItem[];
  loading: boolean;
};

export function DashboardTopExpenses({ transactions, loading }: DashboardTopExpensesProps) {
  return (
    <motion.article className="panel dashboard-panel dashboard-panel--top-expenses" variants={fadeRight}>
      <div className="panel-heading">
        <div>
          <span className="eyebrow">Top gastos</span>
          <h2>Los montos más altos del período</h2>
        </div>
      </div>

      {loading ? (
        <p className="feedback">Cargando información...</p>
      ) : transactions.length === 0 ? (
        <EmptyState
          title="No hay gastos destacados"
          description="Cuando registres gastos, este ranking mostrará los importes más altos."
        />
      ) : (
        <StaggerGroup className="stack-list" onView={false}>
          {transactions.map((transaction) => (
            <motion.div
              key={transaction.id}
              className="list-row list-row--transaction list-row--expense"
              variants={fadeUp}
              whileHover={cardHover}
            >
              <div>
                <strong>{transaction.title}</strong>
                <p>
                  {transaction.category_label} · {transaction.account_name || `Cuenta ${transaction.account_id}`}
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
  );
}
