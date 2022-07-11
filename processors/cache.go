package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
)

type CachePurger struct {
	Ctx context.Context
	Key string
}

func NewCachePurger(ctx context.Context, cacheKey string) *CachePurger {
	return &CachePurger{
		Ctx: ctx,
		Key: cacheKey,
	}
}

func (c *CachePurger) Purge() {
	if err := CacheInstance().Del(c.Key).Err(); err != nil {
		log.Error(c.Ctx, "PurgeCache | key: %v | Error while deleting from redis: %v", c.Key, err.Error())
	}
}
