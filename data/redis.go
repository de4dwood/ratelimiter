package data

import (
	"context"
	"fmt"
	"time"

	"github.com/go-redis/redis/v8"
)

type DataRedis struct {
	client    *redis.Client
	ttl       time.Duration
	keyPrefix string
}

func NewDataRedis(s *redis.Options, ttl time.Duration) Data {
	return &DataRedis{
		ttl:       ttl,
		keyPrefix: "ratelimiter",
		client:    redis.NewClient(s),
	}
}
func (d *DataRedis) Check(ctx context.Context, user, URL, window string, t time.Time) (int64, error) {
	key := d.keyPrefix + ":user:" + user + ":" + URL
	windowt, err := time.ParseDuration(window)
	if err != nil {
		return -1, err
	}
	max := fmt.Sprintf("%f", float64(t.UnixMicro()))
	min := fmt.Sprintf("%f", float64(t.Add(-windowt).UnixMicro()))

	count, err := d.client.ZCount(ctx, key, min, max).Result()
	if err != nil {
		return -1, err
	}
	return count, nil
}
func (d *DataRedis) Request(ctx context.Context, user, URL string, t time.Time) error {
	key := d.keyPrefix + ":user:" + user + ":" + URL
	nowUnix := t.UnixMicro()
	expire := fmt.Sprint(t.Add(-d.ttl).UnixMicro())
	z := redis.Z{
		Score:  float64(nowUnix),
		Member: nowUnix,
	}
	pipe := d.client.TxPipeline()
	pipe.ZRemRangeByScore(ctx, key, "", expire).Result()
	pipe.ZAdd(ctx, key, &z).Result()
	pipe.Expire(ctx, key, d.ttl).Result()
	_, err := pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (d *DataRedis) AddConfig(ctx context.Context, URL, id string, c map[string]interface{}) error {
	key := d.keyPrefix + ":config:" + id + ":" + URL
	result, err := d.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	if len(result) != 0 {
		return fmt.Errorf("config allready exists")
	}
	_, err = d.client.HSet(ctx, key, c).Result()
	if err != nil {
		return err
	}
	return nil
}
func (d *DataRedis) UpdateConfig(ctx context.Context, URL, id string, c map[string]interface{}) error {
	key := d.keyPrefix + ":config:" + id + ":" + URL
	result, err := d.client.Keys(ctx, key).Result()
	if err != nil {
		return err
	}
	if len(result) == 0 {
		return fmt.Errorf("config not exists")
	}
	pipe := d.client.Pipeline()
	pipe.Del(ctx, key).Result()
	pipe.HSet(ctx, key, c)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return err
	}
	return nil

}
func (d *DataRedis) GetConfig(ctx context.Context, URL, id string) (map[string]string, error) {
	key := d.keyPrefix + ":config:" + id + ":" + URL
	result, err := d.client.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, err
	}
	return result, nil

}

func (d *DataRedis) GetConfigs(ctx context.Context, URL string) ([]map[string]string, error) {
	keyPattern := d.keyPrefix + ":config:*:" + URL
	var configs []map[string]string
	keys, err := d.client.Keys(ctx, keyPattern).Result()
	if err != nil {
		return nil, err
	}
	for _, key := range keys {
		result, err := d.client.HGetAll(ctx, key).Result()
		if err != nil {
			return configs, err
		}
		configs = append(configs, result)
	}
	return configs, nil
}
func (d *DataRedis) GetAllConfigs(ctx context.Context) ([]map[string]string, error) {
	keyPattern := d.keyPrefix + ":config:*:*"
	var configs []map[string]string
	keys, err := d.client.Keys(ctx, keyPattern).Result()
	if err != nil {
		return configs, err
	}
	for _, key := range keys {
		result, err := d.client.HGetAll(ctx, key).Result()
		if err != nil {
			return configs, err
		}
		configs = append(configs, result)
	}
	return configs, nil
}
func (d *DataRedis) DeleteConfig(ctx context.Context, URL, id string) error {
	key := d.keyPrefix + ":config:" + id + ":" + URL
	_, err := d.client.Del(ctx, key).Result()
	if err != nil {
		return err
	}
	return nil
}
