package model

import (
	lbpb "trailbox/gen/leaderboard"
	mapspb "trailbox/gen/maps"
	notifpb "trailbox/gen/notifications"
	reviewpb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userpb "trailbox/gen/users"
	workoutpb "trailbox/gen/workouts"
)

// UserProfile representa la vista agregada de un usuario y su actividad.
type UserProfile struct {
	User           *userpb.User               `json:"user,omitempty"`
	Workouts       []*workoutpb.Workout       `json:"workouts,omitempty"`
	Routes         []*routespb.Route          `json:"routes,omitempty"`
	Reviews        []*reviewpb.Review         `json:"reviews,omitempty"`
	Notifications  []*notifpb.Notification    `json:"notifications,omitempty"`
	Leaderboard    *lbpb.LeaderboardEntry     `json:"leaderboard_entry,omitempty"`
	RouteMaps      []*mapspb.GetRouteResponse `json:"maps,omitempty"`
	AggregatedFrom []string                   `json:"aggregated_from,omitempty"`
}
