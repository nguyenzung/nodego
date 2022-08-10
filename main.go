package main

import (
	"fmt"
	"net/http"

	ev "github.com/nguyenzung/go-event-loop/eventloop"
)

func main() {

	ev.NewApp()

	ev.MakeOneTimeTask(5000, func(delay int) { fmt.Println("One time task callback", delay) })

	ev.MakeAPIHandler("/test", func(w *ev.HTTPResponseWriter, r *http.Request) {
		ev.MakeOneTimeTask(10000, func(i int) {
			w.Write([]byte("How are you"))
		})
	})

	ev.MakeAPIHandler("/counter", func(w *ev.HTTPResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			{
				w.Send("GET method in /counter")
			}
		case "POST":
			{
				w.Write([]byte("POST method in /counter"))
			}
		}
	})

	ev.RunApp()
}
