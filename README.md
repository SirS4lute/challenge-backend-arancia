# challenge-backend-arancia

ToDo microservice in Go (Gin) with BoltDB persistence, Docker multi-stage image, and Kubernetes manifests.

## Running

This project is intended to be run via **Docker** or **Kubernetes (Minikube)**.

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

### Minikube

Start Minikube:

```bash
minikube start
```

Build the image inside Minikube:

```bash
minikube image build -t todo-api:local .
```

Apply:

```bash
kubectl apply -f k8s/
```

Important note about persistence/replicas:

- BoltDB is a **local file** and the manifests use a **single PVC** (`ReadWriteOnce`).
- On Minikube (default storage classes), this means the workload should run with **1 replica**.

If needed, scale it down:

```bash
kubectl scale deployment/todo-api --replicas=1
```

Port-forward:

```bash
kubectl port-forward service/todo-api 8080:8080
```

Then access the API at `http://localhost:8080`.

## Demo

With the service running (local/Docker/K8s port-forward), run:

```bash
BASE_URL=http://localhost:18080 ./scripts/demo.sh   # Docker example
# BASE_URL=http://localhost:8080 ./scripts/demo.sh  # Kubernetes port-forward example
```


