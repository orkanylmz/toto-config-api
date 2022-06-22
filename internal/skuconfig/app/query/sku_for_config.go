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

func (s skuForConfigHandler) getCacheKey(cc string, pkg string) string {
	return fmt.Sprintf("%s_%s", strings.ToLower(cc), strings.ToLower(pkg))
}

func (s skuForConfigHandler) Handle(ctx context.Context, query SKUForConfig) (string, error) {

	cacheKey := s.getCacheKey(query.CountryCode, query.PackageName)
	defaultKey := s.getCacheKey("ZZ", query.PackageName)

	randomValue := generateRandomNumber(0, 100)

	// Is cache available to use or not
	if s.cacheModel != nil && query.UseCache {

		// If country code is explicitly passed
		if query.CountryCode == "ZZ" {
			cacheKey = s.getCacheKey("ZZ", query.PackageName)
		}

		// Try to get from cache with cacheKey
		cachedSKU, err := s.cacheModel.SKUForConfig(ctx, cacheKey, randomValue)

		if err != nil {
			// If error happens here, e.g. can't reading cache etc. Don't try to sync cache
			fmt.Println("Lookup for cache error: ", err)

		}

		if cachedSKU != "" {
			fmt.Println("Returning From Cache: ", cachedSKU)
			return cachedSKU, nil
		}

		// If not found in cache and not explicitly set the country code
		if query.CountryCode != "ZZ" {
			// Could not find with cacheKey
			// Try to fetch from db
			allConfigurationsForPkgAndCountry, err := s.readModel.GetAllSKUsForConfig(ctx, query.PackageName, query.CountryCode)
			if err != nil {
				return "", nil
			}

			// If Configurations for for country
			if len(allConfigurationsForPkgAndCountry) > 0 {
				err = s.cacheModel.SyncConfigurations(ctx, cacheKey, allConfigurationsForPkgAndCountry)
				if err == nil {
					return s.cacheModel.SKUForConfig(ctx, cacheKey, randomValue)
				}
				// Sync Err
				return searchInConfigSlice(allConfigurationsForPkgAndCountry, randomValue), nil
			}

			// Check For Default Configuration
			if len(allConfigurationsForPkgAndCountry) == 0 {
				return s.syncAndServeDefaultConfigurationsFromCache(ctx, query.PackageName, defaultKey, randomValue)
			}

		} else {
			return s.syncAndServeDefaultConfigurationsFromCache(ctx, query.PackageName, defaultKey, randomValue)
		}

	} else {
		fmt.Println("Returning directly from DB")
		return s.findSKUFromDB(ctx, query.PackageName, query.CountryCode, randomValue)
	}

	return "", nil
}

func (s skuForConfigHandler) syncAndServeDefaultConfigurationsFromCache(ctx context.Context, packageName string, key string, randomValue int) (string, error) {
	defaultConfigurationsForPkg, err := s.readModel.GetAllSKUsForConfig(ctx, packageName, "ZZ")
	if err != nil {
		return "", nil
	}

	if len(defaultConfigurationsForPkg) == 0 {
		return "", nil
	}

	err = s.cacheModel.SyncConfigurations(ctx, key, defaultConfigurationsForPkg)
	if err == nil {
		return s.cacheModel.SKUForConfig(ctx, key, randomValue)
	}

	// Error for syncing cache, serve it directly from defaultConfigurationsForPkg
	return searchInConfigSlice(defaultConfigurationsForPkg, randomValue), nil
}

func searchInConfigSlice(configSlice []*skuconfig.SKUConfig, val int) string {
	for _, c := range configSlice {
		if int(c.PercentileMin()) < val && val <= int(c.PercentileMax()) {
			return c.SKU()
		}
	}

	return ""
}

func (s skuForConfigHandler) findSKUFromDB(ctx context.Context, pkg string, cc string, val int) (string, error) {
	sku, err := s.readModel.GetSKUForConfig(ctx, pkg, cc, val)

	if err != nil {
		return "", err
	}

	return sku, nil
}
