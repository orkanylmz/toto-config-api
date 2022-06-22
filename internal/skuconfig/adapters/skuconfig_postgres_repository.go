package adapters

import (
	"context"
	"fmt"
	_ "github.com/GoogleCloudPlatform/cloudsql-proxy/proxy/dialers/postgres"
	"github.com/pkg/errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"os"
	"strconv"
	"toto-config-api/internal/skuconfig/domain"
)

type PostgresSKUConfigRepository struct {
	db *gorm.DB
}

func (p PostgresSKUConfigRepository) GetDefaultSKUConfigForPackage(ctx context.Context, packageName string) (*skuconfig.SKUConfig, error) {
	var foundConf *SKUConfigModel
	query := "package = ? AND country_code = ZZ"
	err := p.db.Debug().WithContext(ctx).Where(query, packageName).Last(&foundConf).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Can't find a record with country code, so lets find with a default country code (ZZ)
			return nil, nil
		}

		return nil, err
	}

	newSkuConfig, err := skuconfig.UnmarshalSKUConfigFromDatabase(foundConf.ID, foundConf.Package, foundConf.CountryCode, foundConf.PercentileMin, foundConf.PercentileMax, foundConf.SKU)
	if err != nil {
		fmt.Println("error while converting db object to domain entity")
	}
	return newSkuConfig, nil
}

func (p PostgresSKUConfigRepository) GetAllSKUsForConfig(ctx context.Context, packageName string, countryCode string) ([]*skuconfig.SKUConfig, error) {

	var foundSKUs []SKUConfigModel

	err := p.db.WithContext(ctx).Where("package = ? AND (country_code = ? OR country_code = ZZ)", packageName, countryCode).Find(&foundSKUs).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return make([]*skuconfig.SKUConfig, 0), nil
		}

		return nil, err
	}

	result := make([]*skuconfig.SKUConfig, 0)

	for _, s := range foundSKUs {
		newSkuConfig, err := skuconfig.UnmarshalSKUConfigFromDatabase(s.ID, s.Package, s.CountryCode, s.PercentileMin, s.PercentileMax, s.SKU)
		if err != nil {
			fmt.Println("error while converting db object to domain entity")
		}
		result = append(result, newSkuConfig)
	}

	return result, nil
}

func (p PostgresSKUConfigRepository) GetSKUForConfig(ctx context.Context, packageName string, countryCode string, randomValue int) (string, error) {
	var foundConf *SKUConfigModel
	query := "package = ? AND country_code = ? AND (percentile_min < ? AND percentile_max >= ?)"
	err := p.db.Debug().WithContext(ctx).Where(query, packageName, countryCode, randomValue, randomValue).Last(&foundConf).Error

	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			// Can't find a record with country code, so lets find with a default country code (ZZ)
			err := p.db.Debug().WithContext(ctx).Where(query, packageName, "ZZ", randomValue, randomValue).Last(&foundConf).Error
			if err != nil {
				if errors.Is(err, gorm.ErrRecordNotFound) {
					return "", nil
				}
				return "", err
			}

			return foundConf.SKU, nil
		}

		return "", err
	}

	return foundConf.SKU, nil
}

func NewPostgresSKUConfigRepository(db *gorm.DB) *PostgresSKUConfigRepository {
	if db == nil {
		panic("missing db in postgres repository")
	}

	return &PostgresSKUConfigRepository{db: db}
}

func getConfig() (string, error) {

	host := os.Getenv("POSTGRES_HOST")
	db := os.Getenv("POSTGRES_DB")
	port, err := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if err != nil {
		port = 5432
	}

	password := os.Getenv("POSTGRES_PASSWORD")
	user := os.Getenv("POSTGRES_USER")

	return fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, db), nil
}

func NewPostgresConnection(ctx context.Context) (*gorm.DB, error) {

	dbConnString := os.Getenv("DB_CONN_STRING")

	if dbConnString != "" {
		dbDriver := os.Getenv("DB_DRIVER")

		if dbDriver == "cloudsqlpostgres" {
			db, err := gorm.Open(postgres.New(postgres.Config{
				DriverName: "cloudsqlpostgres",
				DSN: fmt.Sprintf("host=%s user=%s dbname=%s password=%s sslmode=disable",
					os.Getenv("DB_CONN_STRING"),
					os.Getenv("POSTGRES_USER"),
					os.Getenv("POSTGRES_DB"),
					os.Getenv("POSTGRES_PASSWORD"),
				),
			}))

			return db, err
		}
	}

	conf, err := getConfig()

	if err != nil {
		return nil, err
	}

	db, err := gorm.Open(postgres.Open(conf), &gorm.Config{})

	if err != nil {
		return nil, errors.Wrap(err, "Could not connect to DB")
	}

	return db, nil
}
