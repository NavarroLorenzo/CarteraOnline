# Cenz

Cenz es una aplicación web para organizar finanzas personales. Permite administrar cuentas, registrar ingresos y gastos, transferir dinero entre cuentas propias y explorar la actividad en un dashboard con filtros y comparaciones entre períodos.

**[Probar la aplicación](https://pruebaserviciorn.site/)** · **Tecnologías:** React, TypeScript, Go, Gin y PostgreSQL.

## Qué muestra este proyecto

- Una aplicación completa: interfaz React, API REST en Go y persistencia en PostgreSQL.
- Autenticación con JWT y datos separados por usuario.
- Reglas de negocio para cuentas, movimientos y transferencias entre cuentas propias.
- Visualización de saldos, tendencias y distribución de gastos por categoría.
- Pruebas de servicios del backend y entorno local reproducible con Docker Compose.

El proyecto está en desarrollo. Las secciones de inversiones y proyecciones todavía no tienen funcionalidad; más abajo se detallan otras limitaciones y próximos pasos.

Para iniciarlo localmente en Windows: `.\scripts\start-local.cmd`. Para detenerlo: `.\scripts\stop-local.cmd`. La configuración completa está en [Desarrollo local](#desarrollo-local).

## Funcionalidades actuales

- Registro e inicio de sesión.
- Separación de datos por usuario.
- Cuentas de efectivo, banco, billetera virtual, tarjeta y ahorro.
- Saldo inicial al crear una cuenta.
- Ingresos y gastos con categorías.
- Filtros por cuenta, tipo, categoría y fechas.
- Transferencias entre cuentas con validación de saldo.
- Saldo general y por cuenta.
- Dashboard por día, semana, mes, año o rango personalizado.
- Comparación con períodos anteriores.
- Tendencias, distribución de gastos y detalle por categoría.

Los módulos de inversiones y proyecciones existen como estructura inicial, pero todavía no tienen funcionalidad.

## Arquitectura

```text
Frontend React + TypeScript + Vite
              │
              │ HTTP/JSON + JWT Bearer
              ▼
Backend Go + Gin
handlers → services → repositories
              │
              │ pgx/pgxpool
              ▼
PostgreSQL
users + accounts + transactions
```

El frontend es una SPA. Consume una API REST bajo `/api`. El backend autentica con JWT y todas las consultas privadas se filtran por `user_id`.

### Tecnologías

Backend:

- Go 1.26.
- Gin.
- PostgreSQL mediante pgx v5.
- JWT HS256.
- bcrypt.
- godotenv.
- Pruebas con `testing` de Go.

Frontend:

- React 19.
- TypeScript 5.
- Vite 7.
- React Router 7.
- Framer Motion.
- CSS propio.
- Cliente HTTP basado en `fetch`.

## Estructura del repositorio

```text
CarteraOnline/
├── README.md
├── compose.yaml                    PostgreSQL local
├── scripts/                        Inicio y apagado local
├── backend/
│   ├── cmd/api/main.go             Inicio de la API
│   ├── migrations/                 Migraciones SQL
│   └── internal/
│       ├── accounts/               Cuentas financieras
│       ├── auth/                   Usuarios, login y JWT
│       ├── config/                 Variables de entorno
│       ├── database/               PostgreSQL y esquema
│       ├── transactions/           Movimientos y dashboard
│       ├── transfers/              Transferencias
│       ├── routes/                 Endpoints
│       ├── shared/                 Utilidades comunes
│       ├── investments/            Reservado para futuro
│       └── projections/            Reservado para futuro
└── frontend/
    └── src/
        ├── api/                    Cliente y funciones HTTP
        ├── app/                    Routing React
        ├── auth/                   Estado de sesión
        ├── components/             UI, layouts y dashboard
        ├── lib/                    Storage, fechas y formatos
        ├── pages/                  Pantallas principales
        ├── types/                  Tipos de la API
        ├── main.tsx
        └── styles.css
```

## Desarrollo local

### Requisitos

- Docker Desktop iniciado.
- Go en la versión indicada en `backend/go.mod`.
- Node.js 20 o superior con npm.

### Iniciar todo

Desde la raíz del repositorio:

```powershell
.\scripts\start-local.cmd
```

Este comando:

1. inicia PostgreSQL en Docker;
2. espera a que la base esté saludable;
3. instala las dependencias del frontend si faltan;
4. inicia el backend;
5. inicia Vite;
6. verifica que ambos servicios respondan.

Servicios:

- frontend: http://localhost:5173
- API: http://localhost:8080/api/health
- PostgreSQL: `127.0.0.1:5432`
- logs: `.local/logs/`

La primera vez hay que registrar un usuario desde la aplicación.

### Detener todo

```powershell
.\scripts\stop-local.cmd
```

Esto detiene frontend, backend y PostgreSQL, pero conserva la base en un volumen de Docker.

Para borrar voluntariamente todos los datos locales:

```powershell
docker compose down -v
```

> Este último comando elimina de forma permanente el volumen local de PostgreSQL.

### Ejecución manual

Base de datos:

```powershell
docker compose up -d postgres
```

Backend:

```powershell
cd backend
go run ./cmd/api
```

Frontend:

```powershell
cd frontend
npm install
npm run dev
```

## Configuración

### Backend

El backend carga `backend/.env`. Para crear otra configuración se puede partir de `backend/.env.example`.

Variables principales:

| Variable | Uso | Predeterminado |
|---|---|---|
| `PORT` / `APP_PORT` | Puerto HTTP | `8080` |
| `DATABASE_URL` | URL completa de PostgreSQL | vacío |
| `DB_HOST` | Host de PostgreSQL | `localhost` |
| `DB_PORT` | Puerto | `5432` |
| `DB_USER` | Usuario | `postgres` |
| `DB_PASSWORD` | Contraseña | vacío |
| `DB_NAME` | Base | `postgres` |
| `DB_SSLMODE` | Modo SSL | `disable` |
| `JWT_SECRET` | Firma de tokens | inseguro para producción |
| `CORS_ALLOWED_ORIGINS` | Orígenes separados por coma | ninguno |

Si existe `DATABASE_URL`, tiene prioridad sobre las variables `DB_*`.

Nunca utilizar el `JWT_SECRET` local en producción.

### Frontend

Variables principales:

| Variable | Uso |
|---|---|
| `VITE_API_BASE_URL` | Base de la API |
| `VITE_API_PROXY_TARGET` | Proxy local de Vite |
| `VITE_PORT` | Puerto de desarrollo |

El entorno local usa `frontend/.env.development.local`, ignorado por Git y cargado solamente en desarrollo.

Las variables `VITE_*` forman parte del JavaScript del navegador. No deben contener secretos.

## Cloudflare y producción

La configuración local no modifica el build productivo. En Cloudflare Pages se debe conservar una variable equivalente a:

```env
VITE_API_BASE_URL=https://URL-DEL-BACKEND/api
```

El backend productivo debe permitir el dominio del frontend:

```env
CORS_ALLOWED_ORIGINS=https://tu-dominio.pages.dev,https://tu-dominio.com
```

El build productivo se genera con:

```powershell
cd frontend
npm run build
```

## Organización del backend

Los módulos principales siguen este patrón:

- `model.go`: estructuras de dominio, entrada y respuesta.
- `handler.go`: HTTP, validación de JSON y códigos de respuesta.
- `service.go`: reglas de negocio.
- `repository.go`: SQL y persistencia.
- `*_test.go`: pruebas.

Las dependencias se conectan manualmente en `backend/internal/routes/routes.go`.

## Modelo financiero

### Usuarios

`users` almacena email, username, hash bcrypt, fecha de creación y un campo preparado para verificación de email.

### Cuentas

`accounts` almacena nombre, tipo, estado y propietario. Los tipos válidos son:

- `cash`;
- `bank`;
- `virtual_wallet`;
- `credit_card`;
- `savings`.

Dos cuentas activas del mismo usuario no pueden tener el mismo nombre normalizado.

Las cuentas inactivas siguen visibles, pero no aceptan movimientos nuevos.

### Transacciones

`transactions` representa ingresos y gastos. Una cuenta no guarda un saldo fijo:

```text
saldo = suma de ingresos - suma de gastos
```

Reglas principales:

- título obligatorio;
- monto positivo y redondeado a dos decimales;
- tipo `income` o `expense`;
- cuenta activa y perteneciente al usuario;
- categoría válida para el tipo elegido.

Categorías actuales:

- ingresos: sueldo y ventas;
- gastos: comida, supermercado, transporte, combustible, servicios, alquiler, impuestos, farmacia, suscripciones, ocio y ropa;
- ambos: transferencia, ahorro, inversión, deudas y otros.

El saldo inicial se representa como una transacción de ingreso con categoría `otros`.

### Transferencias

Una transferencia crea atómicamente:

- un gasto en la cuenta de origen;
- un ingreso en la cuenta de destino.

Ambos movimientos comparten un `transfer_id`. El backend bloquea las cuentas, verifica que estén activas y comprueba saldo suficiente.

Las transferencias internas se excluyen de ingresos y gastos del dashboard para no duplicar actividad real.

Al borrar una transacción de transferencia se borran las dos partes.

### Eliminación de cuentas

Actualmente la eliminación es física. Se borran:

- la cuenta;
- sus transacciones;
- las contrapartes de transferencias relacionadas en otras cuentas.

Esto mantiene consistencia, pero elimina historial financiero. Es una decisión importante que debería revisarse antes de producción.

## Dashboard

Períodos disponibles:

- día;
- semana;
- mes;
- año;
- rango personalizado.

La respuesta incluye:

- ingresos, gastos y balance;
- cantidad de movimientos;
- comparación con el período anterior;
- tendencia horaria, diaria, semanal o mensual;
- distribución de gastos;
- cinco gastos principales;
- cinco movimientos recientes;
- detalle de una categoría.

## API

### Públicas

| Método | Ruta | Uso |
|---|---|---|
| `GET` | `/api/health` | Estado de la API |
| `POST` | `/api/auth/register` | Registro |
| `POST` | `/api/auth/login` | Login |

### Privadas

| Método | Ruta | Uso |
|---|---|---|
| `GET` | `/api/auth/me` | Usuario actual |
| `POST`, `GET` | `/api/accounts` | Crear/listar cuentas |
| `GET`, `PUT`, `DELETE` | `/api/accounts/:id` | Consultar/editar/borrar cuenta |
| `POST`, `GET` | `/api/transactions` | Crear/listar movimientos |
| `GET` | `/api/transactions/categories` | Categorías |
| `GET` | `/api/transactions/dashboard` | Dashboard |
| `GET` | `/api/transactions/dashboard/category` | Detalle de categoría |
| `GET` | `/api/transactions/balance` | Saldo total |
| `GET` | `/api/transactions/balance-by-account` | Saldos por cuenta |
| `GET` | `/api/transactions/summary` | Resumen filtrado |
| `GET`, `PUT`, `DELETE` | `/api/transactions/:id` | Movimiento individual |
| `POST` | `/api/transfers` | Crear transferencia |

Las rutas privadas requieren:

```http
Authorization: Bearer <token>
```

## Frontend

Rutas públicas:

- `/`;
- `/login`;
- `/register`.

Rutas privadas:

- `/dashboard`;
- `/accounts`;
- `/transactions`;
- `/transfers`.

`AuthContext` mantiene la sesión. `ProtectedRoute` protege las páginas privadas y `PublicOnlyRoute` evita mostrar login o registro a usuarios autenticados.

`src/api/client.ts` centraliza las solicitudes, agrega el JWT, limita tiempos de espera y reintenta GET cuando el servidor tarda en despertar. Las operaciones mutables no se reintentan automáticamente.

El backend permite editar transacciones normales, pero el frontend todavía no ofrece esa interfaz.

## Pruebas

Backend:

```powershell
cd backend
go test ./...
```

Frontend:

```powershell
cd frontend
npm run build
```

La suite Go cubre reglas de cuentas, filtros, categorías, dashboard, conexión de base y transferencias. El frontend todavía no tiene tests ni lint configurados.

## Pendientes importantes

1. Migrar importes de `DOUBLE PRECISION`/`float64` a decimal.
2. Preservar historial al eliminar cuentas.
3. Resolver períodos diarios según la zona horaria del usuario; actualmente se interpretan en UTC.
4. Filtrar cuentas inactivas en la pantalla de transferencias.
5. Agregar paginación al historial.
6. Llevar agregaciones grandes del dashboard a SQL.
7. Unificar la preparación automática del esquema y las migraciones versionadas.
8. Manejar globalmente respuestas `401` en el frontend.
9. Retirar `ClaimOrphanedData` cuando termine la transición legacy.
10. Incorporar tests, lint y formato automático en frontend.

Funcionalidades futuras:

- inversiones;
- proyecciones;
- presupuestos;
- metas de ahorro;
- recuperación de contraseña;
- verificación de email;
- renovación de sesión;
- importación y exportación.

El antiguo README del backend menciona `cmd/seed`, pero ese comando no existe actualmente.

## Reglas para continuar el desarrollo

- Nunca consultar datos financieros sin filtrar por `user_id`.
- Mantener las reglas en servicios y el SQL en repositorios.
- Toda transferencia debe ser atómica.
- Las dos partes de una transferencia deben permanecer sincronizadas.
- No contar transferencias internas como actividad real del dashboard.
- No aceptar movimientos en cuentas inactivas.
- Mantener sincronizados tipos y categorías entre backend y frontend.
- Agregar pruebas para cada regla financiera nueva.
- No registrar tokens, contraseñas o credenciales.
- Actualizar este README si cambia el esquema, la arquitectura o la API.

## Archivos clave

| Tema | Archivo |
|---|---|
| Rutas de la API | `backend/internal/routes/routes.go` |
| Esquema y conexión | `backend/internal/database/database.go` |
| Autenticación | `backend/internal/auth/service.go` |
| JWT | `backend/internal/auth/jwt.go` |
| Cuentas | `backend/internal/accounts/service.go` |
| Movimientos y saldos | `backend/internal/transactions/service.go` |
| Categorías | `backend/internal/transactions/categories.go` |
| Dashboard | `backend/internal/transactions/dashboard.go` |
| Transferencias | `backend/internal/transfers/service.go` y `repository.go` |
| Rutas React | `frontend/src/app/App.tsx` |
| Sesión frontend | `frontend/src/auth/AuthContext.tsx` |
| Cliente HTTP | `frontend/src/api/client.ts` |
| Diseño visual | `frontend/src/styles.css` |

Para seguir una funcionalidad de punta a punta:

```text
ruta backend → handler → service → repository → tipo frontend → módulo api → página/componente
```
