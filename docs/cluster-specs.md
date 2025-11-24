# Cluster Specification — Proyecto Final (Trailbox)

## Arquitectura y componentes
- **Gateway (public entrypoint)**: expone `/api/*` vía HTTP y traduce hacia gRPC de cada dominio. Único recurso `Service` tipo `LoadBalancer`.
- **Usuarios**: CRUD básico de usuarios y consulta de listas. gRPC accesible en `users.final-project.svc.cluster.local:50051`.
- **Rutas**: catálogo de rutas (distancia, desnivel, autor). gRPC en `routes.final-project.svc.cluster.local:50051`.
- **Entrenamientos**: historial de workouts con duración/calorías. gRPC en `workouts.final-project.svc.cluster.local:50051`.
- **Reseñas**: reseñas por ruta (rating y comentario). gRPC en `reviews.final-project.svc.cluster.local:50051`.
- **Leaderboard**: ranking de usuarios y puntajes. gRPC en `leaderboard.final-project.svc.cluster.local:50051`.
- **Notificaciones**: mensajes por usuario. gRPC en `notifications.final-project.svc.cluster.local:50051`.
- **Mapas**: almacena GeoJSON por ruta. gRPC en `maps.final-project.svc.cluster.local:50051`.
- **PostgreSQL**: base de datos compartida, expuesta solo como `ClusterIP`.
- **Frontend Svelte**: SPA estática servida por Nginx (Service `ClusterIP`), consume el gateway.

Todos los flujos de inter-servicio pasan por el gateway; no hay llamadas directas entre servicios de dominio.

## Base de datos
- **DNS de conexión**: `postgres.final-project.svc.cluster.local`, puerto `5432`.
- **Bootstrap**: ConfigMap `postgres-bootstrap` (montado en el Deployment) contiene `01-schema.sql` y `02-seed.sql`. Cada arranque ejecuta ambos scripts vía `postStart` para recrear las bases `users_db`, `routes_db`, `workouts_db`, `reviews_db`, `notifications_db`, `leaderboard_db`, `maps_db`, crear sus roles (`*_app`) y sembrar datos de `data/`.
- **Almacenamiento**: sin PVC; el pod usa `emptyDir` para `/var/lib/postgresql/data`, por lo que el contenido se repuebla en cada reinicio (ideal para demos).
- **Credenciales**: secret `trailbox-db-secret` mantiene el superuser (`DB_*`) y pares específicos por servicio (`USERS_DB_*`, `ROUTES_DB_*`, etc.) que se inyectan en cada Deployment.
- **Tablas**:
  - `users_db.users`: id (uuid), name, age, email (único), created_at.
  - `routes_db.routes`: id (uuid), path, duration, distance, user_id, created_at.
  - `workouts_db.workouts`: id (uuid), name, exercises (jsonb), duration, calories, date, user_id, route_id, created_at.
  - `reviews_db.reviews`: id (uuid), user_id, route_id, rating, comment, created_at.
  - `notifications_db.notifications`: id (uuid), user_id, message, read, created_at.
  - `leaderboard_db.leaderboard`: id (uuid), user_id, score, position, created_at.
  - `maps_db.maps`: id (uuid), route_id, geojson, created_at.

## Manifiestos Kubernetes (`k8s/`)
- `namespace/namespace.yaml`: crea el namespace `final-project`.
- `postgres/`: agrupa `secret.yaml`, `deployment.yaml`, `service.yaml` y `configmap.yaml` (SQL bootstrap) para la base de datos (sin PVC, datos efímeros).
- `users/`, `routes/`, `workouts/`, `reviews/`, `notifications/`, `maps/`, `leaderboard/`: cada carpeta contiene `deployment.yaml` y `service.yaml` (gRPC ClusterIP, probes HTTP `/health` en 8081, recursos ~150m CPU / 192Mi RAM).
- `gateway/`: `deployment.yaml` + `service.yaml` (LoadBalancer puerto 8080). Las variables apuntan a los DNS de cada servicio interno.
- `frontend/`: `deployment.yaml` + `service.yaml` (ClusterIP puerto 80; se expone vía port-forward/Ingress según el clúster).

## Exposición de servicios
- **Público**: solo el gateway (`gateway` Service tipo `LoadBalancer`, puerto 8080). Endpoint interno esperado: `http://gateway.final-project.svc.cluster.local:8080`.
- **Interno (ClusterIP)**: usuarios, rutas, workouts, reviews, leaderboard, notifications, maps y PostgreSQL.
- **Frontend**: ClusterIP (recom.: `kubectl port-forward svc/frontend -n final-project 4173:80` o publicar con un Ingress separado si el entorno lo permite).

## Capacidad y recursos
- Réplicas: 1 por Deployment (gateway, frontend y cada microservicio de dominio).
- Requests/limits aproximados: 150–300m CPU y 192–384Mi RAM para servicios gRPC; 100–200m CPU y 128–256Mi RAM para frontend; 100–250m CPU y 256–512Mi RAM para PostgreSQL.
- Con estos límites y gRPC ligero, se estima capacidad cualitativa de ~50–80 req/s por servicio (CPU-bound) y decenas de usuarios concurrentes navegando la UI/gateway en entornos de laboratorio (minikube/kind).

## Exclusiones explícitas
- **No hay HorizontalPodAutoscaler.**
- **No hay scripts de estrés/carga (Locust/K6).**
- **No hay Consul ni discovery externo**; todo el ruteo usa DNS de Kubernetes.

## Frontend
- Svelte 5 + Vite + Tailwind, servido por Nginx.
- Configurable con `VITE_API_BASE_URL` (ej. `http://gateway.final-project.svc.cluster.local:8080`).
- Rutas: `/` (home), `/users`, `/routes`, `/workouts`, `/reviews`, `/leaderboard`, `/notifications`, `/maps`; cada página invoca `/api/*` en el gateway.
