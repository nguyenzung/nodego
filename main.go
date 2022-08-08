package main

import (
	"fmt"
	"net/http"
	"time"

	ev "github.com/go-event-loop/eventloop"
	"github.com/go-event-loop/threadutils"
)

func main() {
	ev.NewApp()
	fmt.Println(" MainThread ID: ", threadutils.ThreadID())
	ev.MakeAPIHandler("/test", func(w ev.HTTPResponse, r *http.Request) {
		w.Write([]byte(" Test "))
		w.Write([]byte(fmt.Sprintln(time.Now())))
		w.Write([]byte(" ! ThreadID: "))
		w.Write([]byte(fmt.Sprintln(threadutils.ThreadID())))
		var task *ev.TimerTask
		task = ev.MakeTimerTask(10000, func(i int) {
			ev.RemoveTimerTask(task)
			w.Finish()
		})
	})

	ev.MakeAPIHandler("/counter", func(w ev.HTTPResponse, r *http.Request) {
		w.Write([]byte(" Counter "))
		w.Write([]byte(fmt.Sprintln(time.Now())))
		w.Write([]byte(" ! ThreadID: "))
		w.Write([]byte(fmt.Sprintln(threadutils.ThreadID())))
		w.Finish()
	})
	ev.RunApp()
}
