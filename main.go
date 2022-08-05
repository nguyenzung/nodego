package main

import (
	"fmt"

	"github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

var timerTask1 *eventloop.TimerTask

func timeout1(interval int) {
	fmt.Println("Timeout 1 is called: ", interval, threadutils.TheadID())
}

func timeout2(interval int) {
	fmt.Println("Timeout 2 is called: ", interval, threadutils.TheadID())
	eventloop.RemoveTimerTask(timerTask1)
}

func main() {
	eventloop.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.TheadID())
	timerTask1 = eventloop.NewTimerTask(2000, timeout1)
	eventloop.NewTimerTask(9000, timeout2)
	eventloop.RunApp()
}
