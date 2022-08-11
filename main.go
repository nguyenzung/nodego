package main

import (
	"fmt"
	"net/http"

	ev "github.com/nguyenzung/go-event-loop/eventloop"
)

func main() {

	app := ev.NewApp()
	// app.MakeTimerTask(2000, func(delay int) { fmt.Println(" Timer task callback", delay) })
	// app.MakeOneTimeTask(5000, func(delay int) { fmt.Println("One time task callback", delay) })
	app.MakeCallTask("http://localhost:8080", 10, func(s string) {
		fmt.Println("Got data:", s)
	}, func(e error) {
		fmt.Println("[CALL ERROR]", e)
	})
	app.MakeCallTask("http://localhost:8080", 11, func(s string) {
		fmt.Println("Got data:", s)
	}, func(e error) {
		fmt.Println("[CALL ERROR]", e)
	})

	app.MakeAPIHandler("/test", func(w *ev.HTTPResponseWriter, r *http.Request) {
		app.MakeOneTimeTask(10000, func(i int) {
			w.Write([]byte("How are you"))
		})
	})

	app.MakeAPIHandler("/counter", func(w *ev.HTTPResponseWriter, r *http.Request) {
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

	app.RunApp()
}
