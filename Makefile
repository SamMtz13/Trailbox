# ---------- Variables ----------
# No tocar
MODULE      := trailbox
PROTO_DIR   := proto
GEN_ROOT    := gen

KIND_CLUSTER ?= trailbox
NAMESPACE    ?= default
SERVICES     := gateway users routes workouts reviews notifications leaderboard map

# ---------- gRPC protos ----------
# Solo tocar si se agregan nuevos protos
proto:
	mkdir -p $(GEN_ROOT)
	protoc -I $(PROTO_DIR) \
		--go_out=. --go_opt=module=$(MODULE) \
		--go-grpc_out=. --go-grpc_opt=module=$(MODULE) \
		$(PROTO_DIR)/common.proto \
		$(PROTO_DIR)/users.proto \
		$(PROTO_DIR)/routes.proto \
		$(PROTO_DIR)/workouts.proto \
		$(PROTO_DIR)/reviews.proto \
		$(PROTO_DIR)/leaderboard.proto \
		$(PROTO_DIR)/maps.proto \
		$(PROTO_DIR)/notifications.proto

regen: clean proto
clean:
	rm -rf $(GEN_ROOT)

# ---------- Kind cluster ----------
# Borra el cluster existente
cluster-delete:
	-kind delete cluster --name $(KIND_CLUSTER)

# Crea un nuevo cluster
cluster-create:
	kind create cluster --name $(KIND_CLUSTER)

# ---------- Docker builds ----------
# Construye la imagen de un servicio especifico
build-%:
	docker build \
		-f services/$*/Dockerfile \
		-t trailbox/$*:latest \
		.

# Construye las imagenes de todos los servicios
build: $(SERVICES:%=build-%)

# Carga las imagenes de los servicios al cluster kind
kind-load: build
	for svc in $(SERVICES); do \
		kind load docker-image --name $(KIND_CLUSTER) trailbox/$$svc:latest; \
	done

# ---------- Kubernetes ----------
# Aplica los manifiestos de postgres
k8s-postgres:
	kubectl apply -f k8s/postgres

# Aplica los manifiestos de los servicios
k8s-services:
	for dir in frontend gateway users routes workouts reviews notifications leaderboard maps; do \
		kubectl apply -f k8s/$$dir; \
	done

# Muestra el estado de los pods en el namespace especificado
k8s-status:
	kubectl get pods -n $(NAMESPACE)

# Establece los nombres de los objetivos phony (no tocar)
.PHONY: proto regen clean \
	cluster-delete cluster-create \
	build build-% kind-load \
	k8s-namespace k8s-postgres k8s-services k8s-status
