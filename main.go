package main

import (
	"fmt"

	ev "github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

var timerTask1, timerTask3 *ev.TimerTask

func timeout1(interval int) {
	fmt.Println("Timeout 1 is called: ", interval, threadutils.ThreadID())
}

func timeout2(interval int) {
	fmt.Println("Timeout 2 is called: ", interval, threadutils.ThreadID())
	ev.RemoveTimerTask(timerTask1)
	if timerTask3 == nil {
		timerTask3 = ev.MakeTimerTask(3000, timeout3)
	}
}

func timeout3(interval int) {
	fmt.Println("Timeout 3 is called: ", interval, threadutils.ThreadID())
}

func apihandle(res string) {
	fmt.Println("Response: ", res, " ThreadID", threadutils.ThreadID())
}

func exception(err error) {
	fmt.Println("API Call err", err, " ThreadID", threadutils.ThreadID())
}

func main() {
	ev.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.ThreadID())
	timerTask3 = nil
	timerTask1 = ev.MakeTimerTask(2000, timeout1)
	ev.MakeTimerTask(9000, timeout2)
	ev.MakeCallTask("http://localhost:8080", 5, apihandle, exception)
	ev.MakeCallTask("http://localhost:8081", 5, apihandle, exception)
	ev.MakeCallTask("http://localhost:8082", 11, apihandle, exception)
	ev.MakeCallTask("http://localhost:8083", 11, apihandle, exception)
	ev.MakeCallTask("http://localhost:8084", 11, apihandle, exception)
	ev.MakeCallTask("http://localhost:8085", 11, apihandle, exception)
	ev.MakeCallTask("http://localhost:8086", 10, apihandle, exception)
	ev.MakeCallTask("http://localhost:8087", 9, apihandle, exception)
	ev.RunApp()
}
