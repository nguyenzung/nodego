package eventloop

import (
	"fmt"
	"log"
	"net/http"
	"sync"
)

type HTTPResponseWriter struct {
	http.ResponseWriter
	responseChannel chan []byte
}

func (response *HTTPResponseWriter) SendText(data string) {
	response.WriteByteArray([]byte(data))
}

func (response *HTTPResponseWriter) WriteByteArray(data []byte) {
	response.responseChannel <- data
}

func (response *HTTPResponseWriter) wait() {
	response.ResponseWriter.Write(<-response.responseChannel)
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
		responseChannel := make(chan []byte)
		w := &HTTPResponseWriter{rw, responseChannel}
		serveEvent := HTTPServeEvent{w, r, handler}
		httpModule.events <- &serveEvent
		w.wait()
		close(responseChannel)
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
