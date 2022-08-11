package main

import (
	"fmt"
	"net/http"

	ev "github.com/nguyenzung/go-event-loop/eventloop"
	"github.com/nguyenzung/go-event-loop/threadutils"
)

func main() {

	app := ev.NewApp()

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

	app.MakeWSHandler("/", func(message *ev.MessageEvent, session *ev.Session) {
		fmt.Println("[Message]", string(message.Data), " ||| ThreadID:", threadutils.ThreadID())
		session.WriteText([]byte("Hello from mainthread in server"))
	}, func(closeMessage *ev.CloseEvent, session *ev.Session) error {
		fmt.Println("[WS Connection is closed]", closeMessage.Code, closeMessage.Text)
		session.Close(1000, "Bye")
		return nil
	})

	app.RunApp()
}
