package main

import (
	"log"
	"os"
	"user-management/internal/config"
	"user-management/internal/logger"
	"user-management/internal/repository/mongo"
	"user-management/internal/store"
)

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
}

func run() error {
	logger := logger.InitLogger()
	logger.Info("Logger initialized")

	conf, err := config.Load()
	if err != nil {
		logger.Info("Config load failure")
		return err
	}
	logger.Info("Config loaded")

	db, err := store.NewMongoStorage(conf.MongoDtabaseConfig, logger)
	if err != nil {
		logger.Info("New database creation failure")
	}
	logger.Info("New database created")

	_ = mongo.NewUserRepo(db, logger)

	return err
}
