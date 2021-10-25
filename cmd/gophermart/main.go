package main

import (
	"github.com/im-tollu/go-musthave-diploma-tpl/config"
	"log"
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)

	conf, errConf := config.Load()
	if errConf != nil {
		log.Fatalf("Cannot load config: %s", errConf.Error())
	}

	log.Printf("Starting with config %#v", conf)
}
