package clients

import (
	lbpb "trailbox/gen/leaderboard"
	mapspb "trailbox/gen/maps"
	notifpb "trailbox/gen/notifications"
	reviewpb "trailbox/gen/reviews"
	routespb "trailbox/gen/routes"
	userpb "trailbox/gen/users"
	workoutpb "trailbox/gen/workouts"
)

// Clients encapsulates all gRPC clients needed by the gateway.
type Clients struct {
	Users         userpb.UsersClient
	Routes        routespb.RoutesClient
	Workouts      workoutpb.WorkoutsClient
	Reviews       reviewpb.ReviewsClient
	Leaderboard   lbpb.LeaderboardClient
	Notifications notifpb.NotificationsClient
	Maps          mapspb.MapClient
}
