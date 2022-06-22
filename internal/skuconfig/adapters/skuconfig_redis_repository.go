package adapters

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/pkg/errors"
	"os"
	"strconv"
	"toto-config-api/internal/skuconfig/domain"
)

type RedisSKUConfigRepository struct {
	redisClient *redis.Client
}

func (r RedisSKUConfigRepository) SyncConfigurations(ctx context.Context, key string, configurations []*skuconfig.SKUConfig) error {

	sortedSetValues := make([]*redis.Z, 0)

	for _, conf := range configurations {
		sortedSetValues = append(sortedSetValues, &redis.Z{
			Score:  float64(conf.PercentileMax()),
			Member: conf.SKU(),
		})
	}

	err := r.redisClient.ZAdd(ctx, key, sortedSetValues...).Err()
	if err != nil {
		return err
	}

	return nil
}

func (r RedisSKUConfigRepository) SKUForConfig(ctx context.Context, key string, randomValue int) (string, error) {
	isKeyExists, err := r.redisClient.Exists(ctx, key).Result()

	if err != nil {
		return "", errors.Wrap(err, "redisClient.Exists")
	}

	if isKeyExists == 0 {
		return "", skuconfig.KeyNotFoundError
	}

	res, err := r.redisClient.ZRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:   strconv.Itoa(randomValue),
		Max:   "+inf",
		Count: 1,
	}).Result()

	if err != nil {
		return "", errors.Wrap(err, "redisRepo.SKUForConfig.redisClient.Get")
	}

	if len(res) == 0 {
		return "", nil
	}

	return res[0], nil
}

func (r RedisSKUConfigRepository) SetSKU(ctx context.Context, key string, sku string) error {
	fmt.Println("Setting SKU to Redis")
	return nil
}

func NewRedisSKUConfigRepository(redisClient *redis.Client) *RedisSKUConfigRepository {
	return &RedisSKUConfigRepository{redisClient: redisClient}
}

func NewRedisClient(ctx context.Context) (*redis.Client, error) {
	poolSize, err := strconv.Atoi(os.Getenv("REDIS_POOL_SIZE"))
	if err != nil {
		poolSize = 12000
	}

	opts := redis.Options{
		Password: os.Getenv("REDIS_PASSWORD"),
		PoolSize: poolSize,
		Addr:     os.Getenv("REDIS_HOST"),
		DB:       0,
	}
	client := redis.NewClient(&opts)

	return client, ping(client, ctx)
}

func ping(client *redis.Client, ctx context.Context) error {
	pong, err := client.Ping(ctx).Result()
	if err != nil {
		return err
	}
	fmt.Println("REDIS PING RES: ", pong, err)
	return nil
}
