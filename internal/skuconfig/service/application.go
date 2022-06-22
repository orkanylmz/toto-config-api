package service

import (
	"context"
	"github.com/sirupsen/logrus"
	"toto-config-api/internal/common/metrics"
	"toto-config-api/internal/skuconfig/adapters"
	"toto-config-api/internal/skuconfig/app"
	"toto-config-api/internal/skuconfig/app/query"
)

// NewApplication creates a new application with the required Commands & Queries,
// Repositories, Loggers and Metrics Client
func NewApplication(ctx context.Context) app.Application {

	logger := logrus.NewEntry(logrus.StandardLogger())
	metricsClient := metrics.NoOp{}

	redisClient := adapters.NewRedisClient()

	skuconfigRedisRepository := adapters.NewRedisSKUConfigRepository(redisClient)

	conn, err := adapters.NewPostgresConnection(ctx)

	if err != nil {
		panic(err)
	}

	_ = conn.AutoMigrate(adapters.SKUConfigModel{})

	// Here we are creating a repository for read the data for our Queries and Commands
	// Each Query or Command can use a different DB technology as long as they implement
	// our repository interfaces
	// Each implementation can  be found in adapters layer!

	skuconfigPostgreRepository := adapters.NewPostgresSKUConfigRepository(conn)

	return app.Application{
		Queries:  app.Queries{SKUForConfig: query.NewGetSKUForConfigHandler(skuconfigPostgreRepository, skuconfigRedisRepository, logger, metricsClient)},
		Commands: app.Commands{},
	}

}
