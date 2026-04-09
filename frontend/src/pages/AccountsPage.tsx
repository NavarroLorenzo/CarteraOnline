import { AnimatePresence, motion } from "framer-motion";
import { useEffect, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
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
import { formatAccountTypeLabel, formatDate } from "../lib/format";
import type { Account } from "../types/api";

const accountTypes = [
  { value: "cash", label: "Efectivo" },
  { value: "bank", label: "Banco" },
  { value: "virtual_wallet", label: "Billetera virtual" },
  { value: "credit_card", label: "Tarjeta de crédito" },
  { value: "savings", label: "Ahorro" },
];

const accountStatusOptions = [
  { value: "active", label: "Activa" },
  { value: "inactive", label: "Inactiva" },
];

export function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [name, setName] = useState("");
  const [type, setType] = useState(accountTypes[0].value);
  const [initialAmount, setInitialAmount] = useState("0");
  const [editingAccountId, setEditingAccountId] = useState<number | null>(null);
  const [isActive, setIsActive] = useState(true);
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState<string | null>(null);
  const [success, setSuccess] = useState<string | null>(null);

  const loadAccounts = async () => {
    setLoading(true);
    setError(null);

    try {
      const response = await accountsApi.list();
      setAccounts(response);
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudieron cargar las cuentas");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    void loadAccounts();
  }, []);

  const resetForm = () => {
    setName("");
    setType(accountTypes[0].value);
    setInitialAmount("0");
    setIsActive(true);
    setEditingAccountId(null);
  };

  const handleSubmit = async (event: React.FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    setSubmitting(true);
    setError(null);
    setSuccess(null);

    try {
      if (editingAccountId !== null) {
        const confirmed = window.confirm(
          "Vas a modificar esta cuenta. El saldo y las transacciones no cambian, solo el nombre y el tipo. ¿Querés continuar?",
        );
        if (!confirmed) {
          setSubmitting(false);
          return;
        }

        await accountsApi.update(editingAccountId, {
          name,
          type,
          is_active: isActive,
        });
        setSuccess("Cuenta actualizada correctamente");
      } else {
        await accountsApi.create({
          name,
          type,
          initial_amount: Number(initialAmount || 0),
        });
        setSuccess("Cuenta creada correctamente");
      }

      resetForm();
      await loadAccounts();
    } catch (err) {
      setError(
        err instanceof ApiError
          ? err.message
          : editingAccountId !== null
            ? "No se pudo actualizar la cuenta"
            : "No se pudo crear la cuenta",
      );
    } finally {
      setSubmitting(false);
    }
  };

  const startEditing = (account: Account) => {
    setEditingAccountId(account.id);
    setName(account.name);
    setType(account.type);
    setInitialAmount("0");
    setIsActive(account.is_active);
    setError(null);
    setSuccess(`Editando la cuenta "${account.name}"`);
  };

  const handleDelete = async (account: Account) => {
    const confirmed = window.confirm(
      `Vas a eliminar la cuenta "${account.name}". También se van a borrar las transacciones de esa cuenta, pero no las de otras cuentas relacionadas. ¿Querés continuar?`,
    );

    if (!confirmed) {
      return;
    }

    setError(null);
    setSuccess(null);

    try {
      const response = await accountsApi.delete(account.id);
      if (editingAccountId === account.id) {
        resetForm();
      }
      setSuccess(response.warning ? `${response.message}. ${response.warning}` : response.message);
      await loadAccounts();
    } catch (err) {
      setError(err instanceof ApiError ? err.message : "No se pudo eliminar la cuenta");
    }
  };

  return (
    <div className="page-stack">
      <Reveal onView={false}>
        <section className="page-header">
          <div>
            <span className="eyebrow">Cuentas</span>
            <h1>Administrá tus cuentas personales</h1>
            <p>Organizá tu dinero entre efectivo, bancos, billeteras virtuales, tarjetas o ahorro.</p>
          </div>
        </section>
      </Reveal>

      <StaggerGroup className="content-grid" onView={false}>
        <motion.article className="panel" variants={fadeLeft}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">{editingAccountId !== null ? "Editar cuenta" : "Nueva cuenta"}</span>
              <h2>{editingAccountId !== null ? "Actualizar cuenta" : "Datos de la cuenta"}</h2>
            </div>
            {editingAccountId !== null ? (
              <button type="button" className="ghost-button" onClick={resetForm}>
                Cancelar edición
              </button>
            ) : null}
          </div>

          <form className="stack-form" onSubmit={handleSubmit}>
            <label className="field">
              <span>Nombre</span>
              <input
                type="text"
                value={name}
                onChange={(event) => setName(event.target.value)}
                placeholder="Caja diaria"
                required
              />
            </label>

            <label className="field">
              <span>Tipo</span>
              <AnimatedSelect
                value={type}
                options={accountTypes}
                onChange={setType}
                ariaLabel="Tipo de cuenta"
              />
            </label>

            <label className="field">
              <span>Saldo inicial</span>
              <input
                type="number"
                min="0"
                step="0.01"
                value={initialAmount}
                onChange={(event) => setInitialAmount(event.target.value)}
                disabled={editingAccountId !== null}
              />
            </label>

            <AnimatePresence initial={false}>
              {editingAccountId !== null ? (
                <motion.div
                  className="stack-form__cluster"
                  initial={{ opacity: 0, height: 0, y: 10 }}
                  animate={{ opacity: 1, height: "auto", y: 0 }}
                  exit={{ opacity: 0, height: 0, y: -8 }}
                  transition={{ duration: 0.24, ease: [0.22, 1, 0.36, 1] }}
                >
                  <p className="feedback feedback--warning">
                    Podés cambiar el nombre, el tipo y el estado. El saldo se mantiene según los movimientos ya registrados.
                  </p>

                  <label className="field">
                    <span>Estado</span>
                    <AnimatedSelect
                      value={isActive ? "active" : "inactive"}
                      options={accountStatusOptions}
                      onChange={(value) => setIsActive(value === "active")}
                      ariaLabel="Estado de la cuenta"
                    />
                  </label>
                </motion.div>
              ) : null}
            </AnimatePresence>

            <PresenceMessage className="feedback feedback--error">{error}</PresenceMessage>
            <PresenceMessage className="feedback feedback--success">{success}</PresenceMessage>

            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting
                ? "Guardando..."
                : editingAccountId !== null
                  ? "Guardar cambios"
                  : "Crear cuenta"}
            </button>
          </form>
        </motion.article>

        <motion.article className="panel" variants={fadeRight}>
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Resumen</span>
              <h2>Cuentas registradas</h2>
            </div>
            <button type="button" className="ghost-button" onClick={() => void loadAccounts()}>
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
              Cargando cuentas...
            </motion.p>
          ) : accounts.length === 0 ? (
            <EmptyState
              title="No hay cuentas registradas"
              description="Creá una cuenta para empezar a separar tu dinero por origen o uso."
            />
          ) : (
            <StaggerGroup className="stack-list" onView={false}>
              {accounts.map((account) => (
                <motion.div
                  key={account.id}
                  className="list-row"
                  variants={fadeUp}
                  whileHover={cardHover}
                  layout
                >
                  <div>
                    <strong>{account.name}</strong>
                    <p>
                      {formatAccountTypeLabel(account.type)} · {account.is_active ? "activa" : "inactiva"}
                    </p>
                  </div>
                  <div className="list-row__meta">
                    <span>{formatDate(account.created_at)}</span>
                    <div className="action-row action-row--end">
                      <button
                        type="button"
                        className="ghost-button ghost-button--small"
                        onClick={() => startEditing(account)}
                      >
                        Editar
                      </button>
                      <button
                        type="button"
                        className="danger-button danger-button--small"
                        onClick={() => void handleDelete(account)}
                      >
                        Borrar
                      </button>
                    </div>
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
