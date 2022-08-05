package main

import (
	"fmt"

	ev "github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

var timerTask1, timerTask3 *ev.TimerTask

func timeout1(interval int) {
	fmt.Println("Timeout 1 is called: ", interval, threadutils.TheadID())
}

func timeout2(interval int) {
	fmt.Println("Timeout 2 is called: ", interval, threadutils.TheadID())
	ev.RemoveTimerTask(timerTask1)
	if timerTask3 == nil {
		timerTask3 = ev.MakeTimerTask(3000, timeout3)
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
	ev.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.TheadID())
	timerTask3 = nil
	timerTask1 = ev.MakeTimerTask(2000, timeout1)
	ev.MakeTimerTask(9000, timeout2)
	ev.MakeCallTask("http://localhost:8080", 5, apihandle, exception)
	ev.MakeCallTask("http://localhost:8081", 5, apihandle, exception)
	ev.MakeCallTask("http://localhost:8082", 11, apihandle, exception)
	ev.MakeCallTask("http://localhost:8083", 11, apihandle, exception)
	ev.RunApp()
}
