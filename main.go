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
	ev.MakeAPIHandler("/test", func(w http.ResponseWriter, r *http.Request, terminate chan struct{}) {
		w.Write([]byte(fmt.Sprintln(time.Now())))
		w.Write([]byte(" ! ThreadID: "))
		w.Write([]byte(fmt.Sprintln(threadutils.ThreadID())))
		time.Sleep(time.Second * 10)
		terminate <- struct{}{}
	})
	ev.RunApp()
}
