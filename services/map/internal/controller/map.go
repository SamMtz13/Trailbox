package controller

import "context"

type LatLng struct {
	Lat float64 `json:"lat"`
	Lng float64 `json:"lng"`
}

type RouteMap struct {
	RouteID string   `json:"route_id"`
	Coords  []LatLng `json:"coords"`
}

type Repository interface {
	SetRoute(ctx context.Context, rm RouteMap) error
	GetRoute(ctx context.Context, routeID string) (RouteMap, error)
}

type Controller struct{ repo Repository }

func NewController(r Repository) *Controller { return &Controller{repo: r} }

func (c *Controller) SetRoute(ctx context.Context, rm RouteMap) error {
	if rm.RouteID == "" {
		return ErrBadRoute
	}
	return c.repo.SetRoute(ctx, rm)
}

func (c *Controller) GetRoute(ctx context.Context, id string) (RouteMap, error) {
	return c.repo.GetRoute(ctx, id)
}

var ErrBadRoute = &BadReq{"route_id required"}

type BadReq struct{ Msg string }

func (e *BadReq) Error() string { return e.Msg }
