package handler

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"

	commonpb "trailbox/gen/common"
	leaderboardpb "trailbox/gen/leaderboard"
	mappb "trailbox/gen/maps"
	notificationspb "trailbox/gen/notifications"
	reviewspb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userspb "trailbox/gen/users"
	workoutspb "trailbox/gen/workouts"
)

const (
	requestTimeout = 5 * time.Second
	maxBodyBytes   = 1 << 20 // 1 MiB
)

// API wires HTTP handlers to the gRPC clients exposed by the domain services.
type API struct {
	users         userspb.UsersClient
	routes        routespb.RoutesClient
	workouts      workoutspb.WorkoutsClient
	reviews       reviewspb.ReviewsClient
	leaderboard   leaderboardpb.LeaderboardClient
	notifications notificationspb.NotificationsClient
	maps          mappb.MapClient
}

// NewAPI creates a new API aggregator.
func NewAPI(
	users userspb.UsersClient,
	routes routespb.RoutesClient,
	workouts workoutspb.WorkoutsClient,
	reviews reviewspb.ReviewsClient,
	leaderboard leaderboardpb.LeaderboardClient,
	notifications notificationspb.NotificationsClient,
	maps mappb.MapClient,
) *API {
	return &API{
		users:         users,
		routes:        routes,
		workouts:      workouts,
		reviews:       reviews,
		leaderboard:   leaderboard,
		notifications: notifications,
		maps:          maps,
	}
}

// Register attaches all API routes to the provided mux.
func (a *API) Register(mux *http.ServeMux) {
	mux.HandleFunc("/api/users", a.handleUsers)
	mux.HandleFunc("/api/users/", a.handleUserByID)

	mux.HandleFunc("/api/routes", a.handleRoutes)
	mux.HandleFunc("/api/routes/", a.handleRouteByID)

	mux.HandleFunc("/api/workouts", a.handleWorkouts)
	mux.HandleFunc("/api/workouts/", a.handleWorkoutByID)

	mux.HandleFunc("/api/reviews", a.handleReviews)
	mux.HandleFunc("/api/leaderboard", a.handleLeaderboard)
	mux.HandleFunc("/api/notifications", a.handleNotifications)

	mux.HandleFunc("/api/maps", a.handleMapWrite)
	mux.HandleFunc("/api/maps/", a.handleMapByRoute)
}

func (a *API) handleUsers(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.users.ListUsers(ctx, &userspb.ListUsersRequest{})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to list users")
		return
	}
	respondJSON(w, http.StatusOK, resp.GetUsers())
}

func (a *API) handleUserByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, ok := extractID(r.URL.Path, "/api/users/")
	if !ok {
		respondError(w, http.StatusBadRequest, "missing user id")
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.users.GetUser(ctx, &commonpb.UserId{Id: id})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to fetch user")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (a *API) handleRoutes(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.routes.ListRoutes(ctx, &routespb.ListRoutesRequest{})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to list routes")
		return
	}
	respondJSON(w, http.StatusOK, resp.GetRoutes())
}

func (a *API) handleRouteByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, ok := extractID(r.URL.Path, "/api/routes/")
	if !ok {
		respondError(w, http.StatusBadRequest, "missing route id")
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.routes.GetRoute(ctx, &commonpb.RouteId{Id: id})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to fetch route")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (a *API) handleWorkouts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.workouts.ListWorkouts(ctx, &workoutspb.ListWorkoutsRequest{})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to list workouts")
		return
	}
	respondJSON(w, http.StatusOK, resp.GetWorkouts())
}

func (a *API) handleWorkoutByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	id, ok := extractID(r.URL.Path, "/api/workouts/")
	if !ok {
		respondError(w, http.StatusBadRequest, "missing workout id")
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.workouts.GetWorkout(ctx, &commonpb.UserId{Id: id})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to fetch workout")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (a *API) handleReviews(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		routeID := strings.TrimSpace(r.URL.Query().Get("route_id"))
		if routeID == "" {
			respondError(w, http.StatusBadRequest, "route_id query param is required")
			return
		}

		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.reviews.GetReviews(ctx, &reviewspb.ReviewListRequest{RouteId: routeID})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to fetch reviews")
			return
		}
		respondJSON(w, http.StatusOK, resp.GetReviews())
	case http.MethodPost:
		var payload struct {
			UserID  string `json:"userId"`
			RouteID string `json:"routeId"`
			Rating  int32  `json:"rating"`
			Comment string `json:"comment"`
		}
		if err := decodeJSON(r, &payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if strings.TrimSpace(payload.UserID) == "" ||
			strings.TrimSpace(payload.RouteID) == "" ||
			payload.Rating <= 0 {
			respondError(w, http.StatusBadRequest, "userId, routeId and rating are required")
			return
		}

		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.reviews.CreateReview(ctx, &reviewspb.CreateReviewRequest{
			UserId:  payload.UserID,
			RouteId: payload.RouteID,
			Rating:  payload.Rating,
			Comment: payload.Comment,
		})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to create review")
			return
		}
		respondJSON(w, http.StatusCreated, resp)
	default:
		methodNotAllowed(w)
	}
}

func (a *API) handleLeaderboard(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		limit := parseLimit(r.URL.Query().Get("limit"), 10)
		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.leaderboard.GetTop(ctx, &leaderboardpb.GetTopRequest{Limit: limit})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to fetch leaderboard")
			return
		}
		respondJSON(w, http.StatusOK, resp.GetEntries())
	case http.MethodPost:
		var payload struct {
			UserID string `json:"userId"`
			Score  int32  `json:"score"`
		}
		if err := decodeJSON(r, &payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if strings.TrimSpace(payload.UserID) == "" {
			respondError(w, http.StatusBadRequest, "userId is required")
			return
		}

		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.leaderboard.Upsert(ctx, &leaderboardpb.UpsertRequest{
			UserId: payload.UserID,
			Score:  payload.Score,
		})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to update leaderboard")
			return
		}
		respondJSON(w, http.StatusOK, resp)
	default:
		methodNotAllowed(w)
	}
}

func (a *API) handleNotifications(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		userID := strings.TrimSpace(r.URL.Query().Get("user_id"))
		if userID == "" {
			respondError(w, http.StatusBadRequest, "user_id query param is required")
			return
		}

		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.notifications.GetNotifications(ctx, &notificationspb.UserIdRequest{UserId: userID})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to fetch notifications")
			return
		}
		respondJSON(w, http.StatusOK, resp.GetNotifications())
	case http.MethodPost:
		var payload struct {
			UserID  string `json:"userId"`
			Message string `json:"message"`
		}
		if err := decodeJSON(r, &payload); err != nil {
			respondError(w, http.StatusBadRequest, "invalid request body")
			return
		}
		if strings.TrimSpace(payload.UserID) == "" || strings.TrimSpace(payload.Message) == "" {
			respondError(w, http.StatusBadRequest, "userId and message are required")
			return
		}

		ctx, cancel := contextWithTimeout(r)
		defer cancel()

		resp, err := a.notifications.SendNotification(ctx, &notificationspb.SendNotificationRequest{
			UserId:  payload.UserID,
			Message: payload.Message,
		})
		if err != nil {
			respondError(w, http.StatusBadGateway, "failed to send notification")
			return
		}
		respondJSON(w, http.StatusCreated, resp)
	default:
		methodNotAllowed(w)
	}
}

func (a *API) handleMapWrite(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		methodNotAllowed(w)
		return
	}
	var payload struct {
		RouteID string `json:"routeId"`
		GeoJSON string `json:"geoJson"`
	}
	if err := decodeJSON(r, &payload); err != nil {
		respondError(w, http.StatusBadRequest, "invalid request body")
		return
	}
	if strings.TrimSpace(payload.RouteID) == "" || strings.TrimSpace(payload.GeoJSON) == "" {
		respondError(w, http.StatusBadRequest, "routeId and geoJson are required")
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.maps.SetRoute(ctx, &mappb.SetRouteRequest{
		RouteId: payload.RouteID,
		GeoJson: payload.GeoJSON,
	})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to store map")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func (a *API) handleMapByRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		methodNotAllowed(w)
		return
	}
	routeID, ok := extractID(r.URL.Path, "/api/maps/")
	if !ok {
		respondError(w, http.StatusBadRequest, "missing route id")
		return
	}

	ctx, cancel := contextWithTimeout(r)
	defer cancel()

	resp, err := a.maps.GetRoute(ctx, &mappb.GetRouteRequest{RouteId: routeID})
	if err != nil {
		respondError(w, http.StatusBadGateway, "failed to fetch map")
		return
	}
	respondJSON(w, http.StatusOK, resp)
}

func contextWithTimeout(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), requestTimeout)
}

func respondJSON(w http.ResponseWriter, status int, payload interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func respondError(w http.ResponseWriter, status int, msg string) {
	respondJSON(w, status, map[string]string{"error": msg})
}

func decodeJSON(r *http.Request, dst interface{}) error {
	defer r.Body.Close()
	reader := io.LimitReader(r.Body, maxBodyBytes)
	dec := json.NewDecoder(reader)
	dec.DisallowUnknownFields()
	return dec.Decode(dst)
}

func methodNotAllowed(w http.ResponseWriter) {
	respondError(w, http.StatusMethodNotAllowed, "method not allowed")
}

func extractID(path, prefix string) (string, bool) {
	if !strings.HasPrefix(path, prefix) {
		return "", false
	}
	id := strings.Trim(strings.TrimPrefix(path, prefix), "/")
	if id == "" {
		return "", false
	}
	return id, true
}

func parseLimit(value string, def int32) int32 {
	if value == "" {
		return def
	}
	n, err := strconv.Atoi(value)
	if err != nil || n <= 0 {
		return def
	}
	return int32(n)
}
