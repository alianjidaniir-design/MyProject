package dataSources

import (
	"context"
	"time"
)

type RedisDS interface {
	Logout(ctx context.Context, req string, time time.Time) error
}
