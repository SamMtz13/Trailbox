package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	routespb "trailbox/gen/routes"
)

type Client struct {
	conn grpc.ClientConnInterface
	api  routespb.RoutesClient
}

func Dial(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, api: routespb.NewRoutesClient(conn)}, nil
}

func (c *Client) Close() error {
	if cc, ok := c.conn.(*grpc.ClientConn); ok {
		return cc.Close()
	}
	return nil
}

func (c *Client) API() routespb.RoutesClient {
	return c.api
}
