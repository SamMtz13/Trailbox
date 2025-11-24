# Trailbox Cluster Specification

## 1. Cluster Components

### Microservices and responsibilities

| Service | Responsibility | External Interfaces |
| --- | --- | --- |
| `gateway` | Single HTTP entrypoint. Terminates REST calls from browsers, translates them into gRPC requests for the backend. Handles CORS and request logging. | HTTP `/api/*` |
| `users` | CRUD for athlete profiles (name, email, age). Performs migrations on start and exposes gRPC + HTTP health endpoints. | gRPC `Users` |
| `routes` | Stores running/cycling routes (path, distance, duration). | gRPC `Routes` |
| `workouts` | Tracks completed workouts (user, route, calories, exercises JSON). | gRPC `Workouts` |
| `reviews` | Manages user-generated reviews per route (rating + comment). | gRPC `Reviews` |
| `leaderboard` | Aggregates scores per user and exposes a ranking API. | gRPC `Leaderboard` |
| `notifications` | Persists simple inbox messages per user; supports list/send APIs. | gRPC `Notifications` |
| `map` | Stores GeoJSON blobs for each route so they can be rendered on maps. | gRPC `Map` |
| `consul` (support) | Keeps backward-compatible service registration for the Go services. | HTTP :8500, DNS :8600 |
| `frontend` | Svelte 5 SPA that consumes the gateway. Built with Vite + Tailwind. | HTTP (static assets) |

### Public entrypoint

The `gateway` Deployment + LoadBalancer is the only public API entrypoint. All browsers talk to `http(s)://<gateway-lb>/api/...`. The Svelte UI only calls the gateway.

### gRPC dependencies

- `gateway` dials every business service via gRPC (`users`, `routes`, `workouts`, `reviews`, `leaderboard`, `notifications`, `map`).
- No other service calls another one directly. (`routes` still ships an optional HTTP peer client but it is dormant inside Kubernetes.)

### Database

A single PostgreSQL 16 instance (`postgres` Deployment) holds all state inside the `trailbox` database. Tables are created by each service on boot:

| Table | Owner | Notes |
| --- | --- | --- |
| `users` | users svc | UUID PK, name, email, age. |
| `routes` | routes svc | UUID PK, name/path, distance, duration, user UUID. |
| `workouts` | workouts svc | UUID PK, exercises JSONB, calories, duration, FK to users/routes. |
| `reviews` | reviews svc | UUID PK, user/route UUIDs, rating, comment. |
| `leaderboard` | leaderboard svc | UUID PK, user UUID, score, computed rank. |
| `notifications` | notifications svc | UUID PK, user UUID, message, read flag. |
| `maps` | map svc | UUID PK, route UUID, geojson text. |

All services share `DB_HOST=postgres.final-project.svc.cluster.local` plus credentials injected via `trailbox-db-secret`.

## 2. Kubernetes Resources

### Namespace & shared resources

- Namespace: `final-project`.
- Secrets: `trailbox-db-secret` (contains DB credentials for Postgres and the Go services).
- PersistentVolumeClaim: `trailbox-postgres-data` (5 GiB, `ReadWriteOnce`).
- Support Deployments: `postgres`, `consul`.

### Deployments & Services

| Component | Deployment | Service | Type | Replicas | Notes |
| --- | --- | --- | --- | --- | --- |
| postgres | `postgres` | `postgres` | ClusterIP | 1 | Mounted PVC, liveness/readiness via `pg_isready`. |
| consul | `consul` | `consul` | ClusterIP (8500 + 8600/UDP) | 1 | Keeps legacy discovery happy. |
| gateway | `gateway` | `gateway` | **LoadBalancer** | 2 | Exposes port 80 → pod 8080, CORS enabled. |
| frontend | `frontend` | `frontend` | LoadBalancer | 1 | Serves Svelte build via Nginx. |
| users | `users` | `users-service` | ClusterIP | 2 | gRPC 50051, HTTP health 8081. |
| routes | `routes` | `routes-service` | ClusterIP | 2 | Same pattern as users. |
| workouts | `workouts` | `workouts-service` | ClusterIP | 2 | Stores JSONB workouts. |
| reviews | `reviews` | `reviews-service` | ClusterIP | 2 | Handles review CRUD. |
| leaderboard | `leaderboard` | `leaderboard-service` | ClusterIP | 2 | Ranking service. |
| notifications | `notifications` | `notifications-service` | ClusterIP | 2 | Inbox service. |
| map | `maps` | `map-service` | ClusterIP | 1 | GeoJSON persistence. |

All pods declare explicit CPU/memory requests and limits, HTTP readiness/liveness probes (on `/health`, port 8081 for gRPC services, `/health` on 8080 for the gateway). No Horizontal Pod Autoscalers are defined (per professor request); scaling is manual via the `replicas` field or `kubectl scale`.

### Storage & volumes

- PostgreSQL uses the `trailbox-postgres-data` PVC.
- All other services are stateless and rely purely on Postgres, so they do not mount volumes.

## 3. Capacity & Sizing

Assumptions: local kind/minikube cluster with 4 vCPU / 8 GiB RAM nodes.

- Aggregate requests (approximate):
  - CPU: `postgres 200m + gateway 200m + 8×150m + frontend 50m + consul 100m ≈ 1.95 vCPU`.
  - Memory: `postgres 256Mi + gateway 128Mi + 8×256Mi + frontend 64Mi + consul 128Mi ≈ 2.6 GiB`.
- Limits leave enough headroom (roughly 4–5 GiB cluster usage), so a single worker node fits everything comfortably.
- With two replicas per business service and synchronous Postgres writes, the theoretical concurrency is bounded by gRPC handlers and DB connections:
  - Each Go service opens up to 50–100 DB connections; with 8 services the Postgres pool stays below default connection limits.
  - Gateway is stateless and can multiplex hundreds of concurrent HTTP requests (the Go server is configured with sensible read/write timeouts).

Without HPA or stress tests, manual scaling guidelines:

1. Increase `replicas` on the CPU-bound service when p95 latency grows.
2. Bump `postgres` limits first if connection saturation shows up.

## 4. Frontend Integration

- The UI lives under `frontend/` and is a Vite + Svelte 5 SPA styled with Tailwind and `@tailwindcss/forms`.
- All API calls are centralized in `src/lib/api.ts` and point to `import.meta.env.VITE_API_BASE_URL`. Provide `VITE_API_BASE_URL=http://localhost:8080` for local dev and the gateway LoadBalancer hostname for Kubernetes builds.
- Building locally:
  ```bash
  cd frontend
  npm install
  npm run dev          # hot reload @ http://localhost:5173
  npm run build        # outputs dist/ for Docker/Nginx
  ```
- Docker: `frontend/Dockerfile` performs a Node build and serves the static assets with Nginx. Set `--build-arg VITE_API_BASE_URL="https://<gateway-host>"` before pushing.
- Deployment: `k8s/frontend.yaml` creates a `frontend` Deployment + LoadBalancer Service. Point your browser to that LB URL; it will call the gateway LB for every API request.

No load-testing tooling (Locust/K6) or HPAs are part of the manifests anymore, as requested by the professor.
