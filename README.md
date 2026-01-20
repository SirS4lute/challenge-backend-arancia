# challenge-backend-arancia

ToDo microservice in Go (Gin) with BoltDB persistence, Docker multi-stage image, and Kubernetes manifests.

## Running locally

```bash
make run
```

Environment variables:

- `PORT` (default `8080`)
- `DB_PATH` (default `todo.db`)
- `LOG_LEVEL` (`debug|info|warn|error`, default `info`)
- `GIN_MODE` (`debug|release|test`, default `release`)

## API

- `GET /healthz`
- `GET /readyz`
- `GET /todos`
- `POST /todos`
- `PUT /todos/:id`
- `DELETE /todos/:id`

## Docker

Build:

```bash
docker build -t todo-api:local .
```

Run (persisting DB on a local folder):

```bash
mkdir -p .data
chmod 777 .data
docker run --rm -p 18080:8080 -e DB_PATH=/data/todo.db -v "$PWD/.data:/data" todo-api:local
```

## Kubernetes

Manifests are in `k8s/`.

Apply:

```bash
kubectl apply -f k8s/
```

Port-forward:

```bash
kubectl port-forward service/todo-api 8080:8080
```

Using Minikube:

```bash
minikube image build -t todo-api:local .
kubectl rollout restart deployment/todo-api
kubectl get pods -w
kubectl port-forward service/todo-api 8080:8080
```

Note: BoltDB is a local file. With `replicas: 2`, each Pod will have its own DB file (acceptable for the challenge; documented here).

## Demo

With the service running (local/Docker/K8s port-forward), run:

```bash
BASE_URL=http://localhost:18080 ./scripts/demo.sh
```

