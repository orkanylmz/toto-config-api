package query

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"math/rand"
	"strings"
	"time"
	"toto-config-api/internal/common/decorator"
	"toto-config-api/internal/skuconfig/domain"
)

type SKUForConfig struct {
	PackageName string
	CountryCode string
	UseCache    bool
}

type SKUForConfigHandler decorator.QueryHandler[SKUForConfig, string]

type SKUForConfigReadModel interface {
	GetAllSKUsForConfig(ctx context.Context, packageName string, countryCode string) ([]*skuconfig.SKUConfig, error)
	GetSKUForConfig(ctx context.Context, packageName string, countryCode string, randomValue int) (string, error)
}

type SKUForConfigCacheModel interface {
	SKUForConfig(ctx context.Context, key string, randomValue int) (string, error)
	SetSKU(ctx context.Context, key string, sku string) error
	SyncConfigurations(ctx context.Context, key string, configurations []*skuconfig.SKUConfig) error
}

type skuForConfigHandler struct {
	readModel  SKUForConfigReadModel
	cacheModel SKUForConfigCacheModel
}

func NewGetSKUForConfigHandler(
	readModel SKUForConfigReadModel,
	cacheModel SKUForConfigCacheModel,
	logger *logrus.Entry,
	metricsClient decorator.MetricsClient,
) SKUForConfigHandler {
	return decorator.ApplyQueryDecorators[SKUForConfig, string](
		skuForConfigHandler{readModel: readModel, cacheModel: cacheModel},
		logger,
		metricsClient,
	)
}

func generateRandomNumber(min int, max int) int {
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))

	return min + r.Intn(max-min)
}

func (s skuForConfigHandler) Handle(ctx context.Context, query SKUForConfig) (string, error) {

	key := fmt.Sprintf("%s_%s", strings.ToLower(query.CountryCode), strings.ToLower(query.PackageName))
	randomValue := generateRandomNumber(0, 100)

	if s.cacheModel != nil && query.UseCache {

		syncCacheIfNotFound := true

		cachedSKU, err := s.cacheModel.SKUForConfig(ctx, key, randomValue)

		if err != nil {
			fmt.Println("Lookup for cache error: ", err)
			syncCacheIfNotFound = false
		}

		if cachedSKU != "" {
			fmt.Println("Returning From Cache: ", cachedSKU)
			return cachedSKU, nil
		}

		if syncCacheIfNotFound {
			fmt.Println("Syncing cache from DB")
			allConfigurationsForPkgAndCountry, err := s.readModel.GetAllSKUsForConfig(ctx, query.PackageName, query.CountryCode)
			if err != nil {
				return "", nil
			}

			_ = s.cacheModel.SyncConfigurations(ctx, key, allConfigurationsForPkgAndCountry)

			// Simply retrieve the cached value and return
			return s.cacheModel.SKUForConfig(ctx, key, randomValue)
		}

		fmt.Println("Returning directly from cache since cache reading error")
		// Find and return from not syncing
		return s.findSKUFromDB(ctx, query.PackageName, query.CountryCode, randomValue)

	} else {
		fmt.Println("Returning directly from cache, no-opt in cache")
		return s.findSKUFromDB(ctx, query.PackageName, query.CountryCode, randomValue)
	}

}

func (s skuForConfigHandler) findSKUFromDB(ctx context.Context, pkg string, cc string, val int) (string, error) {
	sku, err := s.readModel.GetSKUForConfig(ctx, pkg, cc, val)

	if err != nil {
		return "", err
	}

	return sku, nil
}