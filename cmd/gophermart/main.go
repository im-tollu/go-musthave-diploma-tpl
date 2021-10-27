package main

import (
	"context"
	"github.com/im-tollu/go-musthave-diploma-tpl/api"
	"github.com/im-tollu/go-musthave-diploma-tpl/config"
	"log"
	"os"
	"os/signal"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	conf, errConf := config.Load()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	log.Printf("Starting with config %#v", conf)

	server := api.NewServer(conf.RunAddress)

	awaitTermination()

	if errShutdown := server.Shutdown(context.Background()); errShutdown != nil {
		log.Fatalf("Could not gracefully stop the server: %s", errShutdown.Error())
	}
}

func awaitTermination() {
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, os.Interrupt)
	<-sigint
}
