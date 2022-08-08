package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HTTPResponse struct {
	http.ResponseWriter
	locker *sync.Mutex
}

func (response HTTPResponse) start() {
	response.locker.Lock()
}

func (response HTTPResponse) Finish() {
	response.locker.Unlock()
	// response.flagChannel <- struct{}{}
}

func (response HTTPResponse) wait() {
	response.locker.Lock()
	response.locker.Unlock()
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
	events    chan IEvent
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
			w := HTTPResponse{rw, httpModule.locker}
			w.start()
			serveEvent := HTTPServeEvent{w, r, handler}
			httpModule.events <- &serveEvent
			w.wait()
		}
	})
}

func initHTTPServerModule(events chan IEvent) {
	httpModule = &HTTPServerModule{events, &sync.Mutex{}, make(map[string]func(HTTPResponse, *http.Request))}
}

func startHTTPServerModule() {
	go httpModule.exec()
}
