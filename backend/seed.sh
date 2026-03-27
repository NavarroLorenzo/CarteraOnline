#!/bin/bash

echo "Creando cuentas..."

curl -s -X POST http://localhost:8080/api/accounts \
-H "Content-Type: application/json" \
-d '{"name":"Efectivo","type":"cash","initial_amount":50000}'

curl -s -X POST http://localhost:8080/api/accounts \
-H "Content-Type: application/json" \
-d '{"name":"Mercado Pago","type":"virtual_wallet","initial_amount":20000}'

echo "Creando transacciones..."

curl -s -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"title":"Cobro cliente","amount":15000,"type":"income","account_id":1,"category":"ventas"}'

curl -s -X POST http://localhost:8080/api/transactions \
-H "Content-Type: application/json" \
-d '{"title":"Compra insumos","amount":10000,"type":"expense","account_id":2,"category":"compras"}'

echo "Listo 🚀"
