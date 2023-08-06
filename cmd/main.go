package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"user-management/internal/config"
	"user-management/internal/controller"
	"user-management/internal/jwtpkg"
	"user-management/internal/logger"
	authapiv1 "user-management/internal/proto/usermgt/authapi/v1"
	userapiv1 "user-management/internal/proto/usermgt/userapi/v1"
	"user-management/internal/repository/mongo"
	service "user-management/internal/service/impl"
	"user-management/internal/store"

	"google.golang.org/grpc"
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

	cl, db, err := store.NewMongoStorage(conf.MongoDtabaseConfig, logger)
	if err != nil {
		logger.Info("New database creation failure")
	}
	logger.Info("New database created")

	r := mongo.NewUserRepo(db, cl, logger)
	j := jwtpkg.NewJwtHandler(conf.RsaPair)

	as := service.NewAuthService(r, j)
	us := service.NewUserService(r, as)

	uss := controller.NewUserServiceServer(us, as)
	ass := controller.NewAuthServiceServer(us, as)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	userapiv1.RegisterUserServiceServer(server, uss)
	authapiv1.RegisterAuthServiceServer(server, ass)

	err = server.Serve(lis)

	return err
}
