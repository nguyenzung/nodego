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
		fmt.Println("[Message]", string(message.Data), message.MessageType, message.Err, " ||| [ThreadID]:", threadutils.ThreadID())
		if string(message.Data) == "Bye" {
			session.WriteClose(1000, "")
		} else {
			session.WriteText([]byte("Hello from mainthread in server"))
		}

	}, func(closeMessage *ev.CloseEvent, session *ev.Session) error {
		fmt.Println("[WS Connection is closed]", closeMessage.Code, closeMessage.Text)
		return nil
	})
	app.RunApp()
}
