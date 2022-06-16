package cache

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"

	lru "github.com/hashicorp/golang-lru"
	"go.uber.org/zap"

	"l0/order/repository"
)

type Cache struct {
	cache      *lru.Cache
	size       int32
	orderStore repository.Querier
	logger     *zap.Logger
}

func NewCache(size int, store repository.Querier, logger *zap.Logger) (*Cache, error) {
	cache, err := lru.New(size)
	if err != nil {
		return nil, err
	}

	c := &Cache{
		cache:      cache,
		size:       int32(size),
		orderStore: store,
		logger:     logger,
	}

	return c, nil
}

func (c *Cache) Store(key string, value json.RawMessage) {
	c.cache.ContainsOrAdd(key, value)
}

func (c *Cache) Get(ctx context.Context, key string) (json.RawMessage, bool, error) {
	v, ok := c.cache.Get(key)
	if ok {
		c.logger.Info("got cache hit")
	}
	if !ok {
		o, err := c.orderStore.GetOrderByID(ctx, key)
		if errors.Is(err, sql.ErrNoRows) {
			return json.RawMessage{}, false, nil
		}
		if err != nil {
			return json.RawMessage{}, false, err
		}
		c.Store(o.OrderUid, o.Data)
		v = o.Data
	}

	order, _ := v.(json.RawMessage)

	return order, true, nil
}

func (c *Cache) Recover(ctx context.Context) error {
	params := repository.ListOrdersParams{
		ID:    0,
		Limit: c.size,
	}

	orders, err := c.orderStore.ListOrders(ctx, params)
	if err != nil {
		return err
	}

	for _, v := range orders {
		c.Store(v.OrderUid, v.Data)
	}

	return nil
}
