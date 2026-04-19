# Backend

## Seed de transacciones demo

Para cargar aproximadamente un mes de actividad financiera realista en desarrollo:

```bash
cd backend
go run ./cmd/seed
```

Por defecto el script:

- usa el primer usuario existente, o crea `demo@cenz.local` si no hay usuarios;
- usa cuentas activas existentes, o crea cuentas demo si hacen falta;
- inserta transacciones distribuidas en los ultimos 30 dias;
- modela patrones visibles: sueldo, alquiler, servicios, compras grandes, fines de semana con ocio, semanas de pagos fijos y cierre de mes mas moderado;
- borra antes los datos generados por este mismo seed para evitar duplicados.

Opciones utiles:

```bash
go run ./cmd/seed -user-id=1
go run ./cmd/seed -anchor-date=2026-04-30
go run ./cmd/seed -reset=false
```

Si `APP_ENV`, `GO_ENV` o `ENV` indican produccion, el script se bloquea. Solo se puede forzar con `SEED_ALLOW_PRODUCTION=true`.
