import { AnimatePresence, motion } from "framer-motion";
import type { DashboardCategoryDetail } from "../../types/api";
import { formatCurrency, formatDate, formatTypeLabel } from "../../lib/format";
import { TrendLineChart } from "./TrendLineChart";

type CategoryDetailModalProps = {
  detail: DashboardCategoryDetail | null;
  loading: boolean;
  onClose: () => void;
};

export function CategoryDetailModal({ detail, loading, onClose }: CategoryDetailModalProps) {
  return (
    <AnimatePresence>
      {detail || loading ? (
        <motion.div
          className="dashboard-modal"
          initial={{ opacity: 0 }}
          animate={{ opacity: 1 }}
          exit={{ opacity: 0 }}
          onClick={onClose}
        >
          <motion.div
            className="dashboard-modal__panel"
            initial={{ opacity: 0, y: 32, scale: 0.98 }}
            animate={{ opacity: 1, y: 0, scale: 1 }}
            exit={{ opacity: 0, y: 18, scale: 0.98 }}
            transition={{ duration: 0.28, ease: [0.22, 1, 0.36, 1] }}
            onClick={(event) => event.stopPropagation()}
          >
            <div className="dashboard-modal__header">
              <div>
                <span className="eyebrow">Categoría</span>
                <h2>{detail ? detail.category.label : "Cargando detalle..."}</h2>
                <p>
                  {detail
                    ? `${formatCurrency(detail.category.amount)} · ${detail.category.transactions_count} movimientos`
                    : "Estamos preparando el desglose de tus gastos."}
                </p>
              </div>

              <button type="button" className="ghost-button ghost-button--small" onClick={onClose}>
                Cerrar
              </button>
            </div>

            {loading || !detail ? (
              <p className="feedback">Cargando detalle de la categoría...</p>
            ) : (
              <div className="dashboard-modal__content">
                <div className="dashboard-modal__chart">
                  <TrendLineChart points={detail.trend} showIncome={false} />
                </div>

                <div className="stack-list">
                  {detail.transactions.length === 0 ? (
                    <div className="chart-empty">No hay movimientos para esta categoría en el período activo.</div>
                  ) : (
                    detail.transactions.map((transaction) => (
                      <div
                        key={transaction.id}
                        className={`list-row list-row--transaction list-row--${transaction.type}`}
                      >
                        <div>
                          <strong>{transaction.title}</strong>
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
                      </div>
                    ))
                  )}
                </div>
              </div>
            )}
          </motion.div>
        </motion.div>
      ) : null}
    </AnimatePresence>
  );
}
