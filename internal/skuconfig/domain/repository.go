package skuconfig

import (
	"context"
	"fmt"
)

// This is our Repository Abstraction that can be used to implement with different db solutions / mechanisms

type NotFoundError struct {
	packageName string
}

func (e NotFoundError) Error() string {
	return fmt.Sprintf("sku for package:'%s' not found", e.packageName)

}

type Repository interface {
	GetAllSKUsForConfig(ctx context.Context, packageName string, countryCode string) ([]*SKUConfig, error)
	GetSKUForConfig(ctx context.Context, packageName string, countryCode string, randomValue int) (string, error)
}
