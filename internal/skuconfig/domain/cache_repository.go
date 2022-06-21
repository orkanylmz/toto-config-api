package skuconfig

import (
	"context"
	"errors"
)

var KeyNotFoundError = errors.New("divide by zero")

type CacheRepository interface {
	GetSKUConfig(ctx context.Context, key string) (string, error)
	SetSKU(ctx context.Context, key string, sku string) error
	SyncConfigurations(ctx context.Context, key string, configurations []SKUConfig) error
}
