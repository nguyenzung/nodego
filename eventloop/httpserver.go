package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HTTPResponseWriter struct {
	http.ResponseWriter
	flagChannel chan struct{}
}

func (response *HTTPResponseWriter) Send(data string) {
	response.ResponseWriter.Write([]byte(data))
	response.ResponseWriter = nil
	response.flagChannel <- struct{}{}
}

func (response *HTTPResponseWriter) Write(data []byte) {
	response.ResponseWriter.Write(data)
	response.ResponseWriter = nil
	response.flagChannel <- struct{}{}
}

func (response *HTTPResponseWriter) wait() {
	<-response.flagChannel
}

type HTTPServeEvent struct {
	w       *HTTPResponseWriter
	r       *http.Request
	handler func(*HTTPResponseWriter, *http.Request)
}

func (event *HTTPServeEvent) process() {
	w := event.w
	r := event.r
	event.handler(w, r)
}

type HTTPServerModule struct {
	BaseModule
	server     *http.Server
	locker     *sync.Mutex
	numHandler int
}

func (httpModule *HTTPServerModule) makeAPIHandler(path string, handler func(*HTTPResponseWriter, *http.Request)) {
	httpModule.numHandler++
	httpModule.server.Handler.(*http.ServeMux).HandleFunc(path, func(rw http.ResponseWriter, r *http.Request) {
		flagChannel := make(chan struct{})
		w := &HTTPResponseWriter{rw, flagChannel}
		serveEvent := HTTPServeEvent{w, r, handler}
		httpModule.events <- &serveEvent
		w.wait()
		close(flagChannel)
	})
}

func (httpModule *HTTPServerModule) exec() {
	if httpModule.numHandler > 0 {
		log.Fatal("[LOG]", httpModule.server.ListenAndServe())
	}
}

func makeHTTPServerModule(events chan IEvent) *HTTPServerModule {
	server := makeServer(fmt.Sprintf("%s:%d", HTTP_IP, HTTP_PORT))
	return &HTTPServerModule{BaseModule{events}, server, &sync.Mutex{}, 0}
}
