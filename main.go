package main

import (
	"fmt"

	"github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

var timerTask1, timerTask3 *eventloop.TimerTask

func timeout1(interval int) {
	fmt.Println("Timeout 1 is called: ", interval, threadutils.TheadID())
}

func timeout2(interval int) {
	fmt.Println("Timeout 2 is called: ", interval, threadutils.TheadID())
	eventloop.RemoveTimerTask(timerTask1)
	if timerTask3 == nil {
		timerTask3 = eventloop.MakeTimerTask(3000, timeout3)
	}
}

func timeout3(interval int) {
	fmt.Println("Timeout 3 is called: ", interval, threadutils.TheadID())
}

func apihandle(res string) {
	fmt.Println("Response: ", res, " ThreadID", threadutils.TheadID())
}

func exception(err error) {
	fmt.Println("API Call err", err, " ThreadID", threadutils.TheadID())
}

func main() {
	eventloop.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.TheadID())
	timerTask3 = nil
	timerTask1 = eventloop.MakeTimerTask(2000, timeout1)
	eventloop.MakeTimerTask(9000, timeout2)
	eventloop.MakeCallTask("http://localhost:8080", 5, apihandle, exception)
	eventloop.MakeCallTask("http://localhost:8081", 5, apihandle, exception)
	eventloop.MakeCallTask("http://localhost:8082", 11, apihandle, exception)
	eventloop.MakeCallTask("http://localhost:8083", 11, apihandle, exception)
	eventloop.RunApp()
}
