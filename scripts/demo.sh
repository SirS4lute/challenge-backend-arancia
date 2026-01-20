#!/usr/bin/env bash
set -euo pipefail

BASE_URL="${BASE_URL:-http://localhost:8080}"

echo "Using BASE_URL=$BASE_URL"

echo
echo "== healthz =="
curl -sS "$BASE_URL/healthz" | jq .

echo
echo "== readyz =="
curl -sS "$BASE_URL/readyz" | jq .

echo
echo "== create todo =="
CREATE_RES="$(curl -sS -X POST "$BASE_URL/todos" -H 'Content-Type: application/json' -d '{"title":"buy milk"}')"
echo "$CREATE_RES" | jq .

ID="$(echo "$CREATE_RES" | jq -r '.id')"
echo "Created id=$ID"

echo
echo "== list todos =="
curl -sS "$BASE_URL/todos" | jq .

echo
echo "== update todo =="
curl -sS -X PUT "$BASE_URL/todos/$ID" -H 'Content-Type: application/json' -d '{"title":"buy milk (done)","completed":true}' | jq .

echo
echo "== delete todo =="
curl -sS -X DELETE "$BASE_URL/todos/$ID" -i | head -n 1

echo
echo "== list todos (empty) =="
curl -sS "$BASE_URL/todos" | jq .

