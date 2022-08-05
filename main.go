package main

import (
	"fmt"

	"github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

func timeout1(interval int) {
	fmt.Println("Timeout 1 is called: ", interval, threadutils.TheadID())
}

func timeout2(interval int) {
	fmt.Println("Timeout 2 is called: ", interval, threadutils.TheadID())
}

func main() {
	eventloop.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.TheadID())
	eventloop.NewTimerTask(3000, timeout1)
	eventloop.NewTimerTask(5000, timeout2)
	eventloop.RunApp()
}
