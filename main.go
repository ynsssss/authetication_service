package main

import (
	"authentication_service/config"
	"authentication_service/database"
	"authentication_service/repositories"
	"authentication_service/server"
	"authentication_service/services"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	cfg := config.MustGet("conf.yaml")

	client := database.MustGetConnection(cfg.Database.Url)

	tokensRepo := repositories.NewRefreshTokenRepository(
		client.Database("authentication_service").Collection("refresh_tokens"),
	)

	authService := services.NewAuthService(tokensRepo, cfg.Secret)

	handler := server.NewAuthHandler(authService)
	srv := server.NewHttpServer(cfg.Server.Port, handler)
	go srv.Start()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down")

	ctxServer, cancelServer := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancelServer()
	if err := srv.Shutdown(ctxServer); err != nil {
		log.Fatal("Server forced to shutdown: ", err)
	}

	ctxMongo, cancelMongo := context.WithDeadline(
		context.Background(),
		time.Now().Add(time.Second*10),
	)
	err := client.Disconnect(ctxMongo)
	if err != nil {
		log.Println(err)
	}

	defer cancelMongo()
}
