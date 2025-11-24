MODULE      := trailbox
PROTO_DIR   := proto
GEN_ROOT    := gen

KIND_CLUSTER ?= trailbox
NAMESPACE    ?= default
SERVICES     := gateway users routes workouts reviews notifications leaderboard map

# ---------- gRPC protos ----------
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
cluster-delete:
	-kind delete cluster --name $(KIND_CLUSTER)

cluster-create:
	kind create cluster --name $(KIND_CLUSTER)

# ---------- Docker builds ----------
build-%:
	docker build \
		-f services/$*/Dockerfile \
		-t trailbox/$*:latest \
		.

build: $(SERVICES:%=build-%)

kind-load: build
	for svc in $(SERVICES); do \
		kind load docker-image --name $(KIND_CLUSTER) trailbox/$$svc:latest; \
	done

# ---------- Kubernetes ----------
k8s-postgres:
	kubectl apply -f k8s/postgres

k8s-services:
	for dir in frontend gateway users routes workouts reviews notifications leaderboard maps; do \
		kubectl apply -f k8s/$$dir; \
	done

k8s-status:
	kubectl get pods -n $(NAMESPACE)

.PHONY: proto regen clean \
	cluster-delete cluster-create \
	build build-% kind-load \
	k8s-namespace k8s-postgres k8s-services k8s-status
