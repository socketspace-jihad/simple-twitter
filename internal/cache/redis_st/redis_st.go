package redis_st

import (
	"context"
	"log"
	"simple_twitter/internal/cache"
	"time"

	"github.com/redis/go-redis/v9"
)

type Redis struct {
	Addr   string
	User   string
	Passwd string
	*redis.Client
	ctx context.Context
}

type RedisConfig func(*Redis)

func WithAddr(addr string) RedisConfig {
	return func(r *Redis) {
		r.Addr = addr
	}
}

func WithPassword(passwd string) RedisConfig {
	return func(r *Redis) {
		r.Passwd = passwd
	}
}

func WithUser(user string) RedisConfig {
	return func(r *Redis) {
		r.User = user
	}
}

func NewRedis() cache.Cache {
	rd := &Redis{
		Addr:   "127.0.0.1:6379",
		User:   "default",
		Passwd: "",
		ctx:    context.Background(),
	}
	return rd
}

func (r *Redis) Connect() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     r.Addr,
		Password: r.Passwd,
		DB:       0,
	})
	r.Client = rdb
	pong, err := r.Client.Ping(r.ctx).Result()
	if err != nil {
		return err
	}
	log.Println(pong)
	return nil
}

func (r *Redis) Disconnect() error {
	return r.Client.Close()
}

func (r *Redis) GetString(key string) string {
	res := r.Client.Get(r.ctx, key)
	return res.String()
}

func (r *Redis) GetBytes(key string) []byte {
	res := r.Client.Get(r.ctx, key)
	bytes, err := res.Bytes()
	if err != nil {
		log.Println(err)
		return nil
	}
	return bytes
}

func (r *Redis) GetInt(key string) int {
	res := r.Client.Get(r.ctx, key)
	val, err := res.Int()
	if err != nil {
		log.Println(err)
		return 0
	}
	return val
}

func (r *Redis) Set(key string, value any) {
	r.Client.Set(r.ctx, key, value, 5*time.Minute)
}

func (r *Redis) Delete(key string) {
	r.Client.Del(r.ctx, key)
}

func init() {
	cache.RegisterCache("redis", NewRedis())
}
