package main

import (
	"log"
	"os"
	"user-management/internal/config"
	"user-management/internal/logger"
	"user-management/internal/model"
	"user-management/internal/repository/nosql"
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

	db, err := store.NewNosqlStorage(conf.NoSQLdb, logger)
	if err != nil {
		logger.Info("New database creation failure")
	}
	logger.Info("New database created")

	u := nosql.NewUserRepo(db, logger)

	// u.Create(model.User{
	// 	Username: "Mhmd",
	// 	Password: "1234",
	// 	City:     "Tehran",
	// })

	u.UpdateByUsername(model.User{
		ID:       "64c61ffa46338762e7a4bcf0",
		Username: "Mhmd",
		Password: "12345",
		City:     "Isfahan",
	})

	return nil
}
