package controller

import (
	"context"
	"errors"
	"fmt"
	"time"

	commonpb "trailbox/gen/common"
	lbpb "trailbox/gen/leaderboard"
	mapspb "trailbox/gen/maps"
	notifpb "trailbox/gen/notifications"
	reviewpb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	workoutpb "trailbox/gen/workouts"

	"trailbox/services/gateway/internal/aggregator/model"
	"trailbox/services/gateway/internal/clients"
)

const requestTimeout = 5 * time.Second

type Controller struct {
	clients clients.Clients
}

func New(cl clients.Clients) *Controller {
	return &Controller{clients: cl}
}

func (c *Controller) GetUserProfile(ctx context.Context, userID string) (*model.UserProfile, error) {
	if userID == "" {
		return nil, errors.New("user id is required")
	}

	ctxUser, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	user, err := c.clients.Users.GetUser(ctxUser, &commonpb.UserId{Id: userID})
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	profile := &model.UserProfile{
		User:           user,
		Workouts:       []*workoutpb.Workout{},
		Routes:         []*routespb.Route{},
		Reviews:        []*reviewpb.Review{},
		Notifications:  []*notifpb.Notification{},
		RouteMaps:      []*mapspb.GetRouteResponse{},
		AggregatedFrom: []string{},
	}

	workouts, routeIDs, err := c.fetchWorkouts(ctx, userID)
	if err != nil {
		return nil, err
	}
	profile.Workouts = workouts
	if len(workouts) > 0 {
		profile.AggregatedFrom = append(profile.AggregatedFrom, "workouts")
	}

	if routes, err := c.fetchRoutes(ctx, routeIDs); err == nil {
		profile.Routes = routes
		if len(routes) > 0 {
			profile.AggregatedFrom = append(profile.AggregatedFrom, "routes")
		}
	}

	if reviews, err := c.fetchReviews(ctx, userID, routeIDs); err == nil {
		profile.Reviews = reviews
		if len(reviews) > 0 {
			profile.AggregatedFrom = append(profile.AggregatedFrom, "reviews")
		}
	}

	if notifs, err := c.fetchNotifications(ctx, userID); err == nil {
		profile.Notifications = notifs
		if len(notifs) > 0 {
			profile.AggregatedFrom = append(profile.AggregatedFrom, "notifications")
		}
	}

	if mapsData, err := c.fetchMaps(ctx, routeIDs); err == nil {
		profile.RouteMaps = mapsData
		if len(mapsData) > 0 {
			profile.AggregatedFrom = append(profile.AggregatedFrom, "maps")
		}
	}

	if entry, err := c.fetchLeaderboardEntry(ctx, userID); err == nil && entry != nil {
		profile.Leaderboard = entry
		profile.AggregatedFrom = append(profile.AggregatedFrom, "leaderboard")
	}

	return profile, nil
}

func (c *Controller) fetchWorkouts(ctx context.Context, userID string) ([]*workoutpb.Workout, []string, error) {
	ctxList, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()

	resp, err := c.clients.Workouts.ListWorkouts(ctxList, &workoutpb.ListWorkoutsRequest{})
	if err != nil {
		return nil, nil, fmt.Errorf("list workouts: %w", err)
	}
	var workouts []*workoutpb.Workout
	routeIDs := make(map[string]struct{})
	for _, w := range resp.Workouts {
		if w.GetUserId() == userID {
			workouts = append(workouts, w)
			if w.GetRouteId() != "" {
				routeIDs[w.GetRouteId()] = struct{}{}
			}
		}
	}
	var ids []string
	for id := range routeIDs {
		ids = append(ids, id)
	}
	return workouts, ids, nil
}

func (c *Controller) fetchRoutes(ctx context.Context, ids []string) ([]*routespb.Route, error) {
	var routes []*routespb.Route
	for _, id := range ids {
		ctxRoute, cancel := context.WithTimeout(ctx, requestTimeout)
		route, err := c.clients.Routes.GetRoute(ctxRoute, &commonpb.RouteId{Id: id})
		cancel()
		if err == nil {
			routes = append(routes, route)
		}
	}
	return routes, nil
}

func (c *Controller) fetchReviews(ctx context.Context, userID string, routeIDs []string) ([]*reviewpb.Review, error) {
	var reviews []*reviewpb.Review
	for _, routeID := range routeIDs {
		ctxReviews, cancel := context.WithTimeout(ctx, requestTimeout)
		resp, err := c.clients.Reviews.GetReviews(ctxReviews, &reviewpb.ReviewListRequest{RouteId: routeID})
		cancel()
		if err != nil {
			continue
		}
		for _, r := range resp.GetReviews() {
			if r.GetUserId() == userID {
				reviews = append(reviews, r)
			}
		}
	}
	return reviews, nil
}

func (c *Controller) fetchNotifications(ctx context.Context, userID string) ([]*notifpb.Notification, error) {
	ctxNotif, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	resp, err := c.clients.Notifications.GetNotifications(ctxNotif, &notifpb.UserIdRequest{UserId: userID})
	if err != nil {
		return nil, err
	}
	return resp.GetNotifications(), nil
}

func (c *Controller) fetchMaps(ctx context.Context, routeIDs []string) ([]*mapspb.GetRouteResponse, error) {
	var maps []*mapspb.GetRouteResponse
	for _, routeID := range routeIDs {
		ctxMap, cancel := context.WithTimeout(ctx, requestTimeout)
		resp, err := c.clients.Maps.GetRoute(ctxMap, &mapspb.GetRouteRequest{RouteId: routeID})
		cancel()
		if err == nil {
			maps = append(maps, resp)
		}
	}
	return maps, nil
}

func (c *Controller) fetchLeaderboardEntry(ctx context.Context, userID string) (*lbpb.LeaderboardEntry, error) {
	ctxTop, cancel := context.WithTimeout(ctx, requestTimeout)
	defer cancel()
	resp, err := c.clients.Leaderboard.GetTop(ctxTop, &lbpb.GetTopRequest{Limit: 100})
	if err != nil {
		return nil, err
	}
	for _, entry := range resp.GetEntries() {
		if entry.GetUserId() == userID {
			return entry, nil
		}
	}
	return nil, nil
}
