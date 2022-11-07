package cash

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/sanyarise/hezzl/internal/pb"
)

type CashStore interface {
	CheckCash(key string) bool
	CreateCash(ctx context.Context, res chan *pb.User, key string) error
	GetCash(key string) ([]*pb.User, error)
}

type RedisClient struct {
	*redis.Client
	TTL time.Duration
}

type results struct {
	Responses []*pb.User
}

// NewRedisClient initialize redis client
func NewRedisClient(host, port string, ttl time.Duration) (*RedisClient, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", host, port),
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	err := client.Ping(context.Background()).Err()
	if err != nil {
		return nil, fmt.Errorf("try to ping to redis: %w", err)
	}
	c := &RedisClient{
		Client: client,
		TTL:    ttl,
	}
	return c, nil
}

// Close close redis client
func (c *RedisClient) Close() error {
	return c.Client.Close()
}

// CheckCash checks for data in the cache
func (c *RedisClient) CheckCash(key string) bool {
	item, err := c.GetCash(key)
	if err != nil {
		log.Printf("redis: get record %q: %v", key, err)
		return false
	}

	if item != nil {
		log.Printf("key %q in cash found success", key)
		return true
	}
	log.Printf("redis: get record %q not exist", key)
	return false
}

// CreateCash add data in the cash
func (c *RedisClient) CreateCash(ctx context.Context, res chan *pb.User, key string) error {
	in := results{}
	for resUser := range res {
		in.Responses = append(in.Responses, resUser)
	}

	data, err := json.Marshal(in)
	if err != nil {
		return fmt.Errorf("marshal unknown user: %w", err)
	}

	ttl := 1 * time.Minute

	err = c.Set(ctx, key, data, ttl).Err()
	if err != nil {
		return fmt.Errorf("redis: set key %q: %w", key, err)
	}
	return nil
}

// GetCash retrieves data from the cache
func (c *RedisClient) GetCash(key string) ([]*pb.User, error) {
	res := results{}
	data, err := c.Get(context.Background(), key).Bytes()
	if err == redis.Nil {
		// we got empty result, it's not an error
		return nil, nil
	} else if err != nil {
		return nil, err
	}
	json.Unmarshal(data, &res)
	return res.Responses, nil
}
