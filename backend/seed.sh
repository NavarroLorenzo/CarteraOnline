#!/bin/bash

AUTH_PAYLOAD='{"email":"demo@cartera.local","username":"demo","password":"123456"}'

REGISTER_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/register \
-H "Content-Type: application/json" \
-d "$AUTH_PAYLOAD")

TOKEN=$(printf '%s' "$REGISTER_RESPONSE" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')

if [ -z "$TOKEN" ]; then
  LOGIN_RESPONSE=$(curl -s -X POST http://localhost:8080/api/auth/login \
-H "Content-Type: application/json" \
-d '{"identifier":"demo","password":"123456"}')
  TOKEN=$(printf '%s' "$LOGIN_RESPONSE" | sed -n 's/.*"token":"\([^"]*\)".*/\1/p')
fi

if [ -z "$TOKEN" ]; then
  echo "No se pudo obtener un token JWT"
  echo "$REGISTER_RESPONSE"
  exit 1
fi

echo "Creando cuentas..."

curl -s -X POST http://localhost:8080/api/accounts \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $TOKEN" \
-d '{"name":"Efectivo","type":"cash","initial_amount":50000}'

curl -s -X POST http://localhost:8080/api/accounts \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $TOKEN" \
-d '{"name":"Mercado Pago","type":"virtual_wallet","initial_amount":20000}'

echo "Creando transacciones..."

curl -s -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $TOKEN" \
-d '{"title":"Cobro cliente","amount":15000,"type":"income","account_id":1,"category":"ventas"}'

curl -s -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-H "Authorization: Bearer $TOKEN" \
-d '{"title":"Compra insumos","amount":10000,"type":"expense","account_id":2,"category":"compras"}'

echo "Listo 🚀"
