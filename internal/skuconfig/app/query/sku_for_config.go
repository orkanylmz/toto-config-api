package query

import (
	"context"
	"errors"
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
		fmt.Println("USING CACHE")
		cachedSKU, err := s.cacheModel.SKUForConfig(ctx, key, randomValue)
		fmt.Println("CACHED SKU: ", cachedSKU)
		if err != nil {
			fmt.Println("SKUForConfig Err: ", err)
			if !errors.Is(err, skuconfig.KeyNotFoundError) {
				fmt.Println("SKUForConfig Err: ", "NOT FOUND")
				return "", err
			}

			// Key not found in cache, maybe in db?. Check for package name and country
			// And set a cache from the configurations
			fmt.Println("GetAllSKUsForConfig: ", "GETTING")
			allConfigurationsForPkgAndCountry, err := s.readModel.GetAllSKUsForConfig(ctx, query.PackageName, query.CountryCode)
			if err != nil {
				fmt.Println("Can't get all configuration to sync db to cache")
			}

			if len(allConfigurationsForPkgAndCountry) == 0 {
				return "", errors.New("not found")
			}

			fmt.Println("SyncConfigurations")
			_ = s.cacheModel.SyncConfigurations(ctx, key, allConfigurationsForPkgAndCountry)

			// Simply retrieve the cached value and return
			return s.cacheModel.SKUForConfig(ctx, key, randomValue)
		}

		// If we find in cache, return
		if cachedSKU != "" {
			return cachedSKU, nil
		}

		sku, err := s.readModel.GetSKUForConfig(ctx, query.PackageName, query.CountryCode, randomValue)

		if err != nil {
			return "", err
		}

		if err = s.cacheModel.SetSKU(ctx, key, sku); err != nil {
			return "", err
		}
		return sku, nil
	} else {
		sku, err := s.readModel.GetSKUForConfig(ctx, query.PackageName, query.CountryCode, randomValue)

		if err != nil {
			return "", err
		}

		return sku, nil
	}

}
