package main

import (
	"errors"
	"fmt"
	"log"
	"net/http"

	ev "github.com/nguyenzung/nodego/eventloop"
	"github.com/nguyenzung/nodego/runtimeutils"
)

func sum(arr ...int) (int, error) {
	fmt.Println("Sum function Thread ID", runtimeutils.ThreadID())
	if len(arr) == 0 {
		return -1, errors.New("should have params")
	}
	sum := 0
	for _, value := range arr {
		sum += value
	}
	return sum, nil
}

func main() {
	fmt.Println(" Init ")
	ev.InitApp()

	ev.MakeAPIHandler("/test", func(hw *ev.HTTPResponseWriter, r *http.Request) {
		// Wait 2 seconds before response data to client
		ev.MakeOneTimeTask(2000, func(i int) {
			hw.SendText("123456")
		})
	})

	ev.MakeOneTimeTask(1000, func(i int) {
		ev.MakeCallTask("http://localhost:9090/test", 12, func(s string) { fmt.Println("Response from server", s) }, func(err error) { fmt.Println(err) })
	})

	task := ev.MakeTask(sum,
		func(result int) {
			fmt.Println("Sum =", result, " ||| ThreadID", runtimeutils.ThreadID())
		},
		func(err error) {
			fmt.Println("Error", err, " ||| ThreadID", runtimeutils.ThreadID())
		})
	fmt.Println(task.Exec(1, 2, 3, 5, 6), runtimeutils.ThreadID())
	counter := 0
	ev.MakeWSHandler("/",
		func(s *ev.Session) {
			counter++
			// if counter%100 == 99 {
			log.Println("Total connections", counter)
			// }
		},
		func(me *ev.MessageEvent, s *ev.Session) {
			log.Println("Message:", string(me.Data))
			if string(me.Data) == "q" {
				s.CloseSession(1000, "BYE")
			}
		},
		func(ce *ev.CloseEvent, s *ev.Session) error {
			counter--
			log.Println("Close ", counter)
			return nil
		},
	)
	ev.ExecApp()
}
