package common

import (
	"context"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
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
	if err := processors.CacheInstance().Del(c.Key).Err(); err != nil {
		log.Error(c.Ctx, "PurgeCache | key: %v | Error while deleting from redis: %v", c.Key, err.Error())
	}
}
