
MODULE = trailbox
PROTO_DIR = proto
GEN_ROOT = gen

proto:
	@if not exist $(GEN_ROOT) mkdir $(GEN_ROOT)
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

clean:
	@if exist $(GEN_ROOT) rmdir /S /Q $(GEN_ROOT)

regen: clean proto

.PHONY: proto clean regen
