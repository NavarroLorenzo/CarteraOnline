# Frontend

Frontend React + Vite + TypeScript para consumir el backend de `CarteraOnline`.

## Requisitos

- Node.js 20+
- npm 10+

## Variables

Copiá `.env.example` a `.env` y ajustá si hace falta:

```bash
cp .env.example .env
```

Valores por defecto:

```env
VITE_API_BASE_URL=/api
VITE_API_PROXY_TARGET=http://localhost:8080
VITE_PORT=5173
```

## Desarrollo

En una terminal levantá el backend:

```bash
cd /home/lorenzonavarro/Documents/Pruebas/CarteraOnline/backend
go run ./cmd/api/main.go
```

En otra terminal levantá el frontend:

```bash
cd /home/lorenzonavarro/Documents/Pruebas/CarteraOnline/frontend
npm install
npm run dev
```

Con la configuración por defecto:

- frontend: `http://localhost:5173`
- backend: `http://localhost:8080`

El frontend usa proxy de Vite para `/api`, así evitás problemas de CORS en desarrollo.
