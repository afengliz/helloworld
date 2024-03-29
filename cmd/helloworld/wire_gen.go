// Code generated by Wire. DO NOT EDIT.

//go:generate go run github.com/google/wire/cmd/wire
//+build !wireinject

package main

import (
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"helloworld/internal/biz"
	"helloworld/internal/conf"
	"helloworld/internal/data"
	"helloworld/internal/server"
	"helloworld/internal/service"
)

// Injectors from wire.go:

// initApp init kratos application.
func initApp(confServer *conf.Server, confData *conf.Data, logger log.Logger) (*kratos.App, func(), error) {
	dataData, cleanup, err := data.NewData(confData, logger)
	if err != nil {
		return nil, nil, err
	}
	greeterRepo := data.NewGreeterRepo(dataData, logger)
	greeterUsecase := biz.NewGreeterUsecase(greeterRepo, logger)
	userRepo := data.NewUserRepo(dataData, logger)
	bucketCache := biz.NewBucketCache(userRepo, logger)
	userUsecase := biz.NewUserUsecase(userRepo, logger, bucketCache)
	greeterService := service.NewGreeterService(greeterUsecase, userUsecase, logger)
	httpServer := server.NewHTTPServer(confServer, greeterService, logger)
	grpcServer := server.NewGRPCServer(confServer, greeterService, logger)
	app := newApp(logger, httpServer, grpcServer, bucketCache)
	return app, func() {
		cleanup()
	}, nil
}
