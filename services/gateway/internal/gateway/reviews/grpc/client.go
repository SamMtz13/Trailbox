package grpc

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	reviewpb "trailbox/gen/reviews"
)

type Client struct {
	conn grpc.ClientConnInterface
	api  reviewpb.ReviewsClient
}

func Dial(addr string) (*Client, error) {
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, api: reviewpb.NewReviewsClient(conn)}, nil
}

func (c *Client) Close() error {
	if cc, ok := c.conn.(*grpc.ClientConn); ok {
		return cc.Close()
	}
	return nil
}

func (c *Client) API() reviewpb.ReviewsClient {
	return c.api
}
