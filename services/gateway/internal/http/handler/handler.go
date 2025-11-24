package handler

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"
	"time"

	commonpb "trailbox/gen/common"
	lbpb "trailbox/gen/leaderboard"
	mapspb "trailbox/gen/maps"
	notifpb "trailbox/gen/notifications"
	reviewpb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userpb "trailbox/gen/users"
	workoutpb "trailbox/gen/workouts"

	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/proto"

	aggcontroller "trailbox/services/gateway/internal/aggregator/controller"
	"trailbox/services/gateway/internal/clients"
)

const requestTimeout = 5 * time.Second

type Handler struct {
	clients    clients.Clients
	aggregator *aggcontroller.Controller
}

func New(cl clients.Clients, agg *aggcontroller.Controller) *Handler {
	return &Handler{
		clients:    cl,
		aggregator: agg,
	}
}

func (h *Handler) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", h.handleUsers)
	mux.HandleFunc("/api/users/", h.handleUserByID)

	mux.HandleFunc("/api/routes", h.handleRoutes)
	mux.HandleFunc("/api/routes/", h.handleRouteByID)

	mux.HandleFunc("/api/workouts", h.handleWorkouts)
	mux.HandleFunc("/api/workouts/", h.handleWorkoutByID)

	mux.HandleFunc("/api/reviews", h.handleReviews)
	mux.HandleFunc("/api/leaderboard", h.handleLeaderboard)

	mux.HandleFunc("/api/notifications", h.handleNotifications)
	mux.HandleFunc("/api/notifications/", h.handleNotificationsByUser)

	mux.HandleFunc("/api/maps", h.handleMaps)
	mux.HandleFunc("/api/maps/", h.handleMapByRoute)
	mux.HandleFunc("/api/aggregate/users/", h.handleAggregateUserByID)
}

func (h *Handler) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Users.ListUsers(ctx, &userpb.ListUsersRequest{})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Users.GetUser(ctx, &commonpb.UserId{Id: id})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Routes.ListRoutes(ctx, &routespb.ListRoutesRequest{})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleRouteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/routes/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Routes.GetRoute(ctx, &commonpb.RouteId{Id: id})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Workouts.ListWorkouts(ctx, &workoutpb.ListWorkoutsRequest{})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/workouts/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Workouts.GetWorkout(ctx, &commonpb.UserId{Id: id})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		routeID := r.URL.Query().Get("routeId")
		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		resp, err := h.clients.Reviews.GetReviews(ctx, &reviewpb.ReviewListRequest{RouteId: routeID})
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeProto(w, http.StatusOK, resp)
	case http.MethodPost:
		var req reviewpb.CreateReviewRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		resp, err := h.clients.Reviews.CreateReview(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeProto(w, http.StatusCreated, resp)
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit := int32(10)
		if q := r.URL.Query().Get("limit"); q != "" {
			if v, err := strconv.Atoi(q); err == nil && v > 0 {
				limit = int32(v)
			}
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		resp, err := h.clients.Leaderboard.GetTop(ctx, &lbpb.GetTopRequest{Limit: limit})
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeProto(w, http.StatusOK, resp)
	case http.MethodPost:
		var req lbpb.UpsertRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			writeError(w, http.StatusBadRequest, err)
			return
		}

		ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
		defer cancel()

		resp, err := h.clients.Leaderboard.Upsert(ctx, &req)
		if err != nil {
			writeError(w, http.StatusBadGateway, err)
			return
		}
		writeProto(w, http.StatusOK, resp)
	default:
		methodNotAllowed(w)
	}
}

func (h *Handler) handleNotifications(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var req notifpb.SendNotificationRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Notifications.SendNotification(ctx, &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusCreated, resp)
}

func (h *Handler) handleNotificationsByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	userID := strings.TrimPrefix(r.URL.Path, "/api/notifications/")
	if userID == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Notifications.GetNotifications(ctx, &notifpb.UserIdRequest{UserId: userID})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleMaps(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}

	var req mapspb.SetRouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Maps.SetRoute(ctx, &req)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusCreated, resp)
}

func (h *Handler) handleMapByRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}

	routeID := strings.TrimPrefix(r.URL.Path, "/api/maps/")
	if routeID == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	resp, err := h.clients.Maps.GetRoute(ctx, &mapspb.GetRouteRequest{RouteId: routeID})
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeProto(w, http.StatusOK, resp)
}

func (h *Handler) handleAggregateUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	if h.aggregator == nil {
		writeError(w, http.StatusServiceUnavailable, errors.New("aggregator not configured"))
		return
	}

	id := strings.TrimPrefix(r.URL.Path, "/api/aggregate/users/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	ctx, cancel := context.WithTimeout(r.Context(), requestTimeout)
	defer cancel()

	profile, err := h.aggregator.GetUserProfile(ctx, id)
	if err != nil {
		writeError(w, http.StatusBadGateway, err)
		return
	}
	writeJSON(w, http.StatusOK, profile)
}

func writeProto(w http.ResponseWriter, status int, msg proto.Message) {
	out, err := protojson.MarshalOptions{
		UseProtoNames:   true,
		EmitUnpopulated: true,
	}.Marshal(msg)
	if err != nil {
		writeError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_, _ = w.Write(out)
}

func writeError(w http.ResponseWriter, status int, err error) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]string{"error": err.Error()})
}

func methodNotAllowed(w http.ResponseWriter) {
	http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
}

func writeJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}
