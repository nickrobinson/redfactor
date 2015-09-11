package main

import (
	"fmt"
	"./checker"
	"./notifier"
)

var checkClient checker.Checker

func main() {
	c := make(chan int)
	fmt.Println("Hello World")

	checkClient.NewChecker("./redfactor.db", "192.168.1.206", 8086, "redfactor", "cpu")

	go checkClient.StartChecker(c)
	go notifier.RunNotifier(c)
	select {}
}
