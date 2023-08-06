package store

import (
	"context"
	"time"
	"user-management/internal/config"
	"user-management/internal/dbmigrate"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func NewMongoStorage(conf config.MongoDtabaseConfig, logger *zap.Logger) (*mongo.Client, *mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(conf.URI))
	if err != nil {
		logger.Info("Mongo storage connect failure")
		return nil, nil, err
	}
	logger.Info("Mongo storage connected")

	client.Database(conf.DbName).CreateCollection(ctx, conf.CollectionName, dbmigrate.UsersSchema())

	return client, client.Database(conf.DbName), nil
}
