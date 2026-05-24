package redis

import (
	"MyProject/models/student/dataSources"
	"MyProject/pkg/timeLoc"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

type redisDS struct {
	client *redis.Client
}

func NewRedisDS(addr string) (dataSources.RedisDS, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:         addr,
		DialTimeout:  5 * time.Second,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	})
	ctx := context.Background()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}
	return &redisDS{client: rdb}, nil

}

func (red *redisDS) Logout(ctx context.Context, req string, tim time.Time) error {
	if req == "" || time.Now().In(timeLoc.MyLocation()).After(tim) {
		log.Println("jti یا exp نامعتبر است، توکن access به لیست سیاه اضافه نشد.")
		return nil
	}
	blKey := "bl:access:" + req // کلید blacklist برای access token
	ttl := time.Until(tim)
	if ttl > 0 {
		err := red.client.Set(ctx, blKey, tim.String(), ttl).Err()
		if err != nil {
			log.Printf("خطا در اضافه کردن توکن access به لیست سیاه Redis: %v\n", err)
			return fmt.Errorf("failed to set access token in blacklist: %w", err)
		}
		log.Printf("Access token '%s' blacklisted until %s\n", req, time.Now().In(timeLoc.MyLocation()))
	} else {
		log.Printf("Access Token already expired\n")
	}
	return nil
}
