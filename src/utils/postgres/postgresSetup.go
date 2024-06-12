package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"os"
	"stock_broker_application/src/constants"
	genericConstants "stock_broker_application/src/constants"
	"stock_broker_application/src/models"
	"stock_broker_application/src/utils/configs"
	"stock_broker_application/src/utils/logger"
	"time"

	"github.com/uptrace/opentelemetry-go-extra/otelgorm"
	gormLogger "gorm.io/gorm/logger"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/plugin/prometheus"
)

var DB *gorm.DB

type PostGresConClient struct {
	GormDb *gorm.DB
	SqlDb  *sql.DB
}

var postgresClient *PostGresConClient

// InitPostgresDBConfig initializes the Postgres database configuration and establishes a connection to the database.
func InitPostgresDBConfig(ctx context.Context) error {

	log := logger.GetLoggerWithoutContext()
	PostgresConfig, err := configs.Get(genericConstants.PostgresConfig)
	if err != nil {
		return fmt.Errorf(genericConstants.GetPostgresConfigError, err)
	}

	var postgresConfig models.PostgresConfig

	postgresConfig.Host = getEnv(genericConstants.PostgresHostEnv, PostgresConfig.GetString(genericConstants.PostgresHostKey))
	postgresConfig.Port = getEnv(genericConstants.PostgresPortEnv, PostgresConfig.GetString(genericConstants.PostgresPortKey))
	postgresConfig.User = getEnv(genericConstants.PostgresUserEnv, PostgresConfig.GetString(genericConstants.PostgresUserKey))
	postgresConfig.Password = getEnv(genericConstants.PostgresPasswordEnv, PostgresConfig.GetString(genericConstants.PostgresPasswordKey))
	postgresConfig.DBName = getEnv(genericConstants.PostgresDBNameEnv, PostgresConfig.GetString(genericConstants.PostgresDBNameKey))
	postgresConfig.SSLMode = PostgresConfig.GetString(genericConstants.PostgresSSLModeKey)
	postgresConfig.TimeZone = PostgresConfig.GetString(genericConstants.PostgresTimeZoneKey)
	postgresConfig.IsMockConnection = PostgresConfig.GetBool(genericConstants.PostgresIsMockConnectionKey)

	if postgresConfig.IsMockConnection {
		// Connection Mock for postgres
		err = ConnectMockDatabase(ctx, postgresConfig, log)
		if err != nil {
			return fmt.Errorf(genericConstants.PostgresMockConnectionError, err)
		}
	} else {
		// Connection for postgres
		err = ConnectPostgresDatabase(ctx, postgresConfig, log)
		if err != nil {
			return fmt.Errorf(genericConstants.PostgresConnectionError, err)
		}
	}
	return nil
}

// ConnectPostgresDatabase establishes a connection to the Postgres database using the provided configuration.
func ConnectPostgresDatabase(ctx context.Context, postgresConfig models.PostgresConfig, log logger.Logger) error {
	// Construct the DSN string
	dsn := fmt.Sprintf(genericConstants.DNSString, postgresConfig.Host, postgresConfig.Port, postgresConfig.User, postgresConfig.Password, postgresConfig.DBName, postgresConfig.SSLMode, postgresConfig.TimeZone)

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.Silent),
	})
	if err != nil {
		return err
	}

	// Adds the Prometheus plugin to the Postgres database connection
	db.Use(prometheus.New(prometheus.Config{
		DBName:          postgresConfig.DBName,
		RefreshInterval: 15,
		StartServer:     false,
		MetricsCollector: []prometheus.MetricsCollector{
			&prometheus.Postgres{VariableNames: []string{"Threads_running"}},
		},
	}))

	if err := db.Use(otelgorm.NewPlugin()); err != nil {
		return err
	}

	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	err = sqlDB.PingContext(ctx)
	if err != nil {
		return err
	}
	log.Info(genericConstants.PostgresConnectionSuccessful)

	SetPostgresClient(db, sqlDB)
	DB = db

	return nil
}

// ConnectMockDatabase establishes a connection to a mock Postgres database using the provided configuration.
func ConnectMockDatabase(ctx context.Context, postgresConfig models.PostgresConfig, log logger.Logger) error {
	mockDb, _, err := sqlmock.New()
	if err != nil {
		return err
	}
	dialConnection := postgres.New(postgres.Config{
		Conn:       mockDb,
		DriverName: constants.PostgresDriverName,
	})
	db, err := gorm.Open(dialConnection, &gorm.Config{})
	if err != nil {
		return err
	}
	SetPostgresClient(db, mockDb)
	log.Info(genericConstants.PostgresMockConnectionSuccessful)
	return nil
}

func SetPostgresClient(db *gorm.DB, sqlDB *sql.DB) {
	postgresClient = &PostGresConClient{GormDb: db, SqlDb: sqlDB}
}

func ClosePostgres(ctx context.Context) {
	log := logger.GetLoggerWithoutContext()
	if postgresClient != nil {
		err := postgresClient.SqlDb.Close()
		if err != nil {
			log.Error(genericConstants.ClosePostgresClientError, zap.Error(err))
		}
	}
}

func GetPostGresClient() *PostGresConClient {
	return postgresClient
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}
