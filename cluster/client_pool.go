package cluster

import (
	"context"
	"errors"
	pool "github.com/jolestar/go-commons-pool/v2"
	"godis/resp/client"
)

type connectionFactory struct {
	Peer string
}

func (c *connectionFactory) MakeObject(ctx context.Context) (*pool.PooledObject, error) {
	cc, err := client.NewClient(c.Peer)
	if err != nil {
		return nil, err
	}
	cc.Start()
	return pool.NewPooledObject(cc), nil
}

func (c *connectionFactory) DestroyObject(ctx context.Context, object *pool.PooledObject) error {
	cc, ok := object.Object.(*client.Client)
	if !ok {
		return errors.New("type mismatch")
	}
	cc.Close()
	return nil
}

func (c *connectionFactory) ValidateObject(ctx context.Context, object *pool.PooledObject) bool {
	//TODO implement me
	panic("implement me")
}

func (c *connectionFactory) ActivateObject(ctx context.Context, object *pool.PooledObject) error {
	//TODO implement me
	panic("implement me")
}

func (c *connectionFactory) PassivateObject(ctx context.Context, object *pool.PooledObject) error {
	//TODO implement me
	panic("implement me")
}
