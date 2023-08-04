package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
	"user-management/internal/config"
	"user-management/internal/controller"
	"user-management/internal/jwtpkg"
	"user-management/internal/logger"
	"user-management/internal/model"
	"user-management/internal/proto/protoconnect"
	"user-management/internal/repository/mongo"
	service "user-management/internal/service/impl"
	"user-management/internal/store"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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
	j := jwtpkg.NewJwtHandler(conf.RsaPair)
	as := service.NewAuthService(r, j)
	us := service.NewUserService(r, as)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	fmt.Println(as.Login(ctx, model.Username("su"), model.Password("Admin@123")))

	mux := http.NewServeMux()
	path, handler := protoconnect.NewUserServiceHandler(&controller.UserServiceServer{
		UserSrv: us,
		AuthSrv: as,
	})

	mux.Handle(path, handler)
	err = http.ListenAndServe(
		conf.HttpURI,
		h2c.NewHandler(mux, &http2.Server{}),
	)

	// uss := controller.NewUserServiceServer(us, as)

	// lis, err := net.Listen("tcp", fmt.Sprintf(":%s", conf.Port))
	// if err != nil {
	// 	return err
	// }

	// server := grpc.NewServer()
	// protobuf.RegisterUserServiceServer(server, uss)

	// err = server.Serve(lis)

	return err
}
