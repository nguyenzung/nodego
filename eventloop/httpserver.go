package eventloop

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-event-loop/threadutils"
)

type HTTPServeEvent struct {
	w       http.ResponseWriter
	r       *http.Request
	s       chan struct{}
	handler func(http.ResponseWriter, *http.Request, chan struct{})
}

func (event *HTTPServeEvent) process() {
	fmt.Println("Process a request ", threadutils.ThreadID())
	w := event.w
	r := event.r
	s := event.s
	event.handler(w, r, s)
}

type HTTPServerModule struct {
	events    chan IEvent
	apiMapper map[string]func(http.ResponseWriter, *http.Request, chan struct{})
}

func (httpModule *HTTPServerModule) exec() {
	log.Fatal(http.ListenAndServe("0.0.0.0:8080", nil))
}

var httpModule *HTTPServerModule

func MakeAPIHandler(path string, handler func(http.ResponseWriter, *http.Request, chan struct{})) {
	httpModule.apiMapper[path] = handler
	http.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("A request comes ", httpModule.apiMapper)
		handler, status := httpModule.apiMapper[path]
		fmt.Println("Check", &handler, status)
		if status {
			fmt.Println("Found handler")
			flag := make(chan struct{})
			serveEvent := HTTPServeEvent{w, r, flag, handler}
			httpModule.events <- &serveEvent
			fmt.Println("Finish a request:", <-flag)
		}
	})
}

func initHTTPServerModule(events chan IEvent) {
	httpModule = &HTTPServerModule{events, make(map[string]func(http.ResponseWriter, *http.Request, chan struct{}))}
}

func startHTTPServerModule() {
	go httpModule.exec()
}
