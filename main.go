package main

import (
	"fmt"
	"github.com/nickrobinson/checker"
	"github.com/nickrobinson/notifer"
)

func main() {
	c := make(chan int)
	fmt.Println("Hello World")
	go checker.RunChecker(c)
	go notifier.RunNotifier(c)
	select {}
}
