package main

import (
	"context"
	"log"

	"pvz-service/internal/app"
	"pvz-service/internal/config"
)

func main() {
	cfg := config.MustLoad("./config/config.yaml")
	ctx := context.TODO()

	a, err := app.NewApp(ctx, cfg)
	if err != nil {
		log.Fatal("cannot setup server:", err.Error())
	}
	go func() {
		if err = a.Run(); err != nil {
			log.Fatal("stop server:", err.Error())
		}
	}()

	<-ctx.Done()
	log.Println("got interruption signal")
	if err := a.Stop(); err != nil {
		log.Printf("server shutdown returned an err: %v\n", err)
	}
}
