package main

import (
	"watchlist/business"
	serviceConstants "watchlist/commons/constants"
	"watchlist/router"

	"context"
	"fmt"
	genericConstants "stock_broker_application/src/constants"
	loggerMiddleware "stock_broker_application/src/middleware/logger"
	"stock_broker_application/src/utils/configs"
	"stock_broker_application/src/utils/flags"
	"stock_broker_application/src/utils/logger"
	"stock_broker_application/src/utils/postgres"

	"go.uber.org/zap"
)

// @title Stock Broker Application
// @description   api for Stock Broker using gin and gorm
// @version 2.0
// @host localhost:8081
// @securityDefinitions.apiKey JWT
// @in header
// @name Authorization
// @BasePath /
func main() {
	ctx := context.Background()
	initConfigs(ctx)
	initLogger(ctx)
	log := logger.GetLoggerWithoutContext()
	err := postgres.InitPostgresDBConfig(ctx)
	if err != nil {
		log.With(zap.Error(err)).Fatal(genericConstants.PostgresDBInitializationError)
	}
	defer postgres.ClosePostgres(ctx)
	err = configs.InitApplicationConfigs(ctx)
	if err != nil {
		log.With(zap.Error(err)).Fatal(genericConstants.ConfigBindingFailedError)
		panic(fmt.Errorf(genericConstants.ConfigBindingFailedError))
	}

	// Start Kafka producer in a separate goroutine
	go business.StartProducer()

	startRouter(ctx)

}

func initConfigs(ctx context.Context) {
	configs.Init([]string{flags.BaseConfigPath(), flags.MockConfigPath()})
}

func initLogger(ctx context.Context) {
	LoggerConfig, err := configs.Get(genericConstants.LoggerConfig)
	if err != nil {
		panic(err)
	}
	logger.SetupLogging(LoggerConfig.GetString(genericConstants.LogLevelKey))
}

func startRouter(ctx context.Context) {
	log := logger.GetLoggerWithoutContext()
	router := router.GetRouter(loggerMiddleware.Logger())
	log.Info(fmt.Sprintf(genericConstants.RunningServerPort, serviceConstants.PortDefaultValue))
	err := router.Run(fmt.Sprintf(":%d", serviceConstants.PortDefaultValue))
	if err != nil {
		log.With(zap.Error(err)).Fatal(genericConstants.ExternalServiceError)
	}
}
