package store

import (
	"context"
	"time"
	"user-management/internal/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewMongoStorage(conf config.MongoDtabaseConfig, logger *zap.Logger) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.URI))
	if err != nil {
		logger.Info("Mongo storage connect failure")
		return nil, err
	}
	logger.Info("Mongo storage connected")

	return client.Database(conf.Name), nil
}
