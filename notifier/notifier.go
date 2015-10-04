package notifier

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

func StartNotifier(c chan int) {
	for {
		rowId := <-c
		fmt.Println("Row Id: ", rowId)
	}
}

func Notify() {
	var cfg Config
	err := gcfg.ReadFileInto(&cfg, "../redfactor.cfg")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(cfg.Smtp.Host)
	fmt.Println(cfg.Smtp.Port)
}
