# Frontend

Frontend React + Vite + TypeScript para consumir el backend de Cenz.

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
VITE_API_PROXY_TARGET=
VITE_PORT=5173
```

Para desarrollo local, si querés usar proxy de Vite, definí también:

```env
VITE_API_PROXY_TARGET=http://localhost:8080
```

Para producción en Cloudflare Pages, configurá en el dashboard:

```env
VITE_API_BASE_URL=https://api.tudominio.com/api
```

## Desarrollo

En una terminal levantá el backend:

```bash
cd ./backend
go run ./cmd/api/main.go
```

En otra terminal levantá el frontend:

```bash
cd ./frontend
npm install
npm run dev
```

Con la configuración por defecto:

- frontend: `http://localhost:5173`
- backend: `http://localhost:8080`

Si definís `VITE_API_PROXY_TARGET`, el frontend usa proxy de Vite para `/api` y evitás problemas de CORS en desarrollo.

## Producción

Si publicás el frontend en Cloudflare Pages y el backend en otro dominio, necesitás:

- `VITE_API_BASE_URL` apuntando al backend productivo
- `CORS_ALLOWED_ORIGINS` en el backend con el dominio del frontend

Ejemplo:

```env
CORS_ALLOWED_ORIGINS=https://tu-proyecto.pages.dev,https://tudominio.com,https://www.tudominio.com
```
