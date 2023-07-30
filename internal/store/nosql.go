package store

import (
	"context"
	"fmt"
	"time"
	"user-management/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewNosqlStorage(conf config.NoSQLDatabaseConfig, logger *zap.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	uri := fmt.Sprintf("mongodb://%s:%s", conf.Host, conf.Port)
	client, err := mongo.NewClient(options.Client().ApplyURI(uri))
	if err != nil {
		logger.Info("Nosql client make failure")
		return nil, err
	}
	logger.Info("Nosql client made")

	err = client.Connect(ctx)
	if err != nil {
		logger.Info("Nosql storage connect failure")
		return nil, err
	}
	logger.Info("Nosql storage connected")

	return client.Database(conf.Name), nil
}
