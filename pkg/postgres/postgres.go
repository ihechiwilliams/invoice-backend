package postgres

import (
	"fmt"
	"time"

	gormtrace "gopkg.in/DataDog/dd-trace-go.v1/contrib/gorm.io/gorm.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

const (
	defaultConnectionLifetime = 5 * time.Minute
	defaultMaxPoolSize        = 25
)

type Config struct {
	Name            string
	Password        string
	PrimaryHost     string
	ReadReplicaHost string
	User            string
	Port            string
}

func InitDB(serviceName string, config *Config) (*gorm.DB, error) {
	fullServiceName := fmt.Sprintf("%s.postgres", serviceName)

	primaryConfig := buildPostgresConfig(config.PrimaryHost, config.User, config.Password, config.Name, config.Port)
	readReplicaConfig := buildPostgresConfig(
		config.ReadReplicaHost,
		config.User,
		config.Password,
		config.Name,
		config.Port,
	)

	gormDB, err := gormtrace.Open(
		postgres.Open(primaryConfig.DSN),
		&gorm.Config{Logger: logger.Default.LogMode(logger.Silent)},
		gormtrace.WithServiceName(fullServiceName),
	)
	if err != nil {
		return nil, err
	}

	err = gormDB.Use(
		dbresolver.Register(
			dbresolver.Config{
				Replicas: []gorm.Dialector{postgres.Open(readReplicaConfig.DSN)},
				Policy:   dbresolver.RandomPolicy{},
			},
		),
	)
	if err != nil {
		return nil, err
	}

	sqlDB, err := gormDB.DB()
	if err != nil {
		return nil, err
	}

	sqlDB.SetMaxOpenConns(defaultMaxPoolSize)
	sqlDB.SetMaxIdleConns(defaultMaxPoolSize)
	sqlDB.SetConnMaxLifetime(defaultConnectionLifetime)

	return gormDB, nil
}

func buildPostgresConfig(host, user, password, name, port string) postgres.Config {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=disable TimeZone=UTC",
		host,
		user,
		password,
		name,
		port,
	)

	return postgres.Config{DSN: dsn}
}
