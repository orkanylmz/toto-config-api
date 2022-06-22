package skuconfig

import (
	"context"
	"errors"
)

var KeyNotFoundError = errors.New("key not found in cache")

type CacheRepository interface {
	SKUForConfig(ctx context.Context, key string, randomValue int) (string, error)
	SetSKU(ctx context.Context, key string, sku string) error
	SyncConfigurations(ctx context.Context, key string, configurations []SKUConfig) error
	IsCacheKeyExists(ctx context.Context, key string) (bool, error)
	SyncConfiguration(ctx context.Context, key string, configuration *SKUConfig) error
}
