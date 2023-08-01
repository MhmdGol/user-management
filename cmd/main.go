package main

import (
	"log"
	"os"
	"user-management/internal/config"
	"user-management/internal/logger"
	"user-management/internal/repository/mongo"
	service "user-management/internal/service/impl"
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

	r := mongo.NewUserRepo(db, logger)
	as := service.NewAuthService(r, conf.RsaPair)
	_ = service.NewUserService(r, as)

	// uss := controller.NewUserServiceServer(s)

	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	return err
	// }

	// server := grpc.NewServer()
	// protobuf.RegisterUserServiceServer(server, uss)

	// err = server.Serve(lis)
	return err
}
