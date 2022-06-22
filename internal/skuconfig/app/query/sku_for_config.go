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
	GetDefaultSKUConfigForPackage(ctx context.Context, packageName string) ([]*skuconfig.SKUConfig, error)
}

type SKUForConfigCacheModel interface {
	SKUForConfig(ctx context.Context, key string, defaultKey string, randomValue int) (string, error)
	SetSKU(ctx context.Context, key string, sku string) error
	SyncConfigurations(ctx context.Context, key string, configurations []*skuconfig.SKUConfig) error
	SyncConfiguration(ctx context.Context, key string, configuration *skuconfig.SKUConfig) error
	IsCacheKeyExists(ctx context.Context, key string) (bool, error)
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

func (s skuForConfigHandler) getCacheKey(cc string, pkg string, ) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(cc), strings.ToLower(pkg))
}

func (s skuForConfigHandler) Handle(ctx context.Context, query SKUForConfig) (string, error) {

	cacheKey := s.getCacheKey(query.CountryCode, query.PackageName)
	defaultKey := s.getCacheKey("ZZ", query.PackageName)
	randomValue := generateRandomNumber(0, 100)

	// Is cache available to use or not
	if s.cacheModel != nil && query.UseCache {

		// flag for syncing cache
		syncCacheIfNotFound := true

		cachedSKU, err := s.cacheModel.SKUForConfig(ctx, cacheKey, defaultKey, randomValue)

		if err != nil {
			// If error happens here, e.g. can't reading cache etc. Don't try to sync cache
			fmt.Println("Lookup for cache error: ", err)
			syncCacheIfNotFound = false
		}

		if cachedSKU != "" {
			fmt.Println("Returning From Cache: ", cachedSKU)
			return cachedSKU, nil
		}

		if syncCacheIfNotFound {
			fmt.Println("Syncing cache from DB")

			// all configurations return configurations including the default one (cc = ZZ)
			allConfigurationsForPkgAndCountry, err := s.readModel.GetAllSKUsForConfig(ctx, query.PackageName, query.CountryCode)
			if err != nil {
				return "", nil
			}

			defaultConfigurationsForPkg, err := s.readModel.GetAllSKUsForConfig(ctx, query.PackageName, "ZZ")

			if len(allConfigurationsForPkgAndCountry) == 0 && len(defaultConfigurationsForPkg) == 0 {
				return "", nil
			}
			
			// Sync All Configurations to Cache
			err = s.cacheModel.SyncConfigurations(ctx, cacheKey, allConfigurationsForPkgAndCountry)
			// Sync Default Configurations to Cache
			err = s.cacheModel.SyncConfigurations(ctx, defaultKey, defaultConfigurationsForPkg)

			if err == nil {
				// If not fail to sync
				// Simply retrieve the cached value and return
				return s.cacheModel.SKUForConfig(ctx, cacheKey, defaultKey, randomValue)
			}

			// If we fail to sync the cache, flow will continue by reading DB
		}

		fmt.Println("Returning directly from DB")
		// Find and return from not syncing
		return s.findSKUFromDB(ctx, query.PackageName, query.CountryCode, randomValue)

	} else {
		fmt.Println("Returning directly from DB")
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
