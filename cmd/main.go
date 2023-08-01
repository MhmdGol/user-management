package main

import (
	"log"
	"os"
	"user-management/internal/config"
	"user-management/internal/logger"
	"user-management/internal/model"
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
	us := service.NewUserService(r, as)
	// uss := controller.NewUserServiceServer(s)
	// t := "eyJhbGciOiJSUzI1NiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2OTA5MzE4NDAsInJvbGUiOiJhZG1pbiJ9.C8mMuK9KKrP-G7XtST_5BcuLCegt8DL8GpIxeT8M9azJqFHUrx53hJy23uJvem4pHEq5RrlkcWuFTkSWk775WtWyNIDhkxiU2kjajv10SYBMv1PfMQoPVmIEjcVZ6VHtpVpvFLdcmzFwFP4aM68q086tFn3DN-PkTST8avXtHqQ"

	// us.Create(model.User{
	// 	Username: "Mhmd",
	// 	Password: "1234",
	// 	Role:     "staff",
	// 	City:     "Tehran",
	// }, model.JwtToken{Token: t})

	t2, _ := as.Login(model.LoginRequest{
		Username: "Mhmd",
		Password: "1234",
	})

	// us.DeleteByID(model.ID("64c88ca7a82bc4c248bafa04"), model.JwtToken{Token: t})
	us.UpdateByID(model.User{
		ID:      model.ID("64c88d09a82bc4c248bafa05"),
		Role:    "user",
		City:    "Tehran",
		Version: 1,
	}, t2)
	//
	//
	//
	// lis, err := net.Listen("tcp", ":50051")
	// if err != nil {
	// 	return err
	// }

	// server := grpc.NewServer()
	// protobuf.RegisterUserServiceServer(server, uss)

	// err = server.Serve(lis)
	return err
}
