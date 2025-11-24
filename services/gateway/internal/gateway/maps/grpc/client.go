package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	mapspb "trailbox/gen/maps"
)

type Client struct {
	conn grpc.ClientConnInterface
	api  mapspb.MapClient
}

func Dial(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, api: mapspb.NewMapClient(conn)}, nil
}

func (c *Client) Close() error {
	if cc, ok := c.conn.(*grpc.ClientConn); ok {
		return cc.Close()
	}
	return nil
}

func (c *Client) API() mapspb.MapClient {
	return c.api
}
