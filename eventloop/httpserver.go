package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HTTPResponse struct {
	http.ResponseWriter
	flagChannel chan struct{}
}

func (response HTTPResponse) Finish() {
	response.flagChannel <- struct{}{}
}

func (response HTTPResponse) wait() {
	<-response.flagChannel

}

type HTTPServeEvent struct {
	w       HTTPResponse
	r       *http.Request
	handler func(HTTPResponse, *http.Request)
}

func (event *HTTPServeEvent) process() {
	w := event.w
	r := event.r
	event.handler(w, r)
}

type HTTPServerModule struct {
	BaseModule
	locker    *sync.Mutex
	apiMapper map[string]func(HTTPResponse, *http.Request)
}

func (httpModule *HTTPServerModule) exec() {
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%d", HTTP_IP, HTTP_PORT), nil))
}

var httpModule *HTTPServerModule

func MakeAPIHandler(path string, handler func(HTTPResponse, *http.Request)) {
	httpModule.apiMapper[path] = handler
	http.HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		handler, status := httpModule.apiMapper[path]
		if status {
			flagChannel := make(chan struct{})
			w := HTTPResponse{rw, flagChannel}
			serveEvent := HTTPServeEvent{w, r, handler}
			httpModule.events <- &serveEvent
			w.wait()
		}
	})
}

func initHTTPServerModule(events chan IEvent) {
	httpModule = &HTTPServerModule{BaseModule{events}, &sync.Mutex{}, make(map[string]func(HTTPResponse, *http.Request))}
}

func startHTTPServerModule() {
	go httpModule.exec()
}
