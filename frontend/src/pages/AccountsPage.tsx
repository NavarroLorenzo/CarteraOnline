import { useEffect, useState } from "react";
import { accountsApi } from "../api/accounts";
import { ApiError } from "../api/client";
import { EmptyState } from "../components/ui/EmptyState";
import { formatDate } from "../lib/format";
import type { Account } from "../types/api";

const accountTypes = [
  { value: "cash", label: "Efectivo" },
  { value: "bank", label: "Banco" },
  { value: "virtual_wallet", label: "Billetera virtual" },
  { value: "savings", label: "Ahorro" },
];

export function AccountsPage() {
  const [accounts, setAccounts] = useState<Account[]>([]);
  const [name, setName] = useState("");
  const [type, setType] = useState(accountTypes[0].value);
  const [initialAmount, setInitialAmount] = useState("0");
  const [editingAccountId, setEditingAccountId] = useState<number | null>(null);
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
    setError(null);
    setSuccess(`Editando la cuenta "${account.name}"`);
  };

  const handleDelete = async (account: Account) => {
    const confirmed = window.confirm(
      `Vas a eliminar la cuenta "${account.name}". Esto también va a borrar sus transacciones asociadas. ¿Querés continuar?`,
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
      <section className="page-header">
        <div>
          <span className="eyebrow">Cuentas</span>
          <h1>Administrá tus cuentas personales</h1>
          <p>Podés crear, editar nombre/tipo y borrar cuentas con confirmación previa.</p>
        </div>
      </section>

      <section className="content-grid">
        <article className="panel">
          <div className="panel-heading">
            <div>
              <span className="eyebrow">{editingAccountId !== null ? "Editar cuenta" : "Nueva cuenta"}</span>
              <h2>{editingAccountId !== null ? "Actualizar datos" : "Formulario"}</h2>
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
              <select value={type} onChange={(event) => setType(event.target.value)}>
                {accountTypes.map((option) => (
                  <option key={option.value} value={option.value}>
                    {option.label}
                  </option>
                ))}
              </select>
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

            {editingAccountId !== null ? (
              <p className="feedback feedback--warning">
                Estás editando una cuenta existente. Desde acá solo se actualizan nombre y tipo.
              </p>
            ) : null}

            {error ? <p className="feedback feedback--error">{error}</p> : null}
            {success ? <p className="feedback feedback--success">{success}</p> : null}

            <button type="submit" className="primary-button" disabled={submitting}>
              {submitting
                ? "Guardando..."
                : editingAccountId !== null
                  ? "Guardar cambios"
                  : "Crear cuenta"}
            </button>
          </form>
        </article>

        <article className="panel">
          <div className="panel-heading">
            <div>
              <span className="eyebrow">Listado</span>
              <h2>Tus cuentas</h2>
            </div>
            <button type="button" className="ghost-button" onClick={() => void loadAccounts()}>
              Recargar
            </button>
          </div>

          {loading ? (
            <p className="feedback">Cargando cuentas...</p>
          ) : accounts.length === 0 ? (
            <EmptyState
              title="No tenés cuentas cargadas"
              description="Cuando crees la primera, la vas a ver listada acá."
            />
          ) : (
            <div className="stack-list">
              {accounts.map((account) => (
                <div key={account.id} className="list-row">
                  <div>
                    <strong>{account.name}</strong>
                    <p>{account.type}</p>
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
                </div>
              ))}
            </div>
          )}
        </article>
      </section>
    </div>
  );
}
