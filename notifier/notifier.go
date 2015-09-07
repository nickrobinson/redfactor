package main

import (
	"gopkg.in/gcfg.v1"
	"log"
	"fmt"
)

type Config struct {
        Smtp struct {
                Host string
                Port int
        }
}

type Notifier struct {
	
}

func Notify() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "test.gcfg")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Smtp.Host)
	fmt.Println(cfg.Smtp.Port)
}
