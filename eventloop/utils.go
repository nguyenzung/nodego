package eventloop

import "net/http"

func makeServer(addr string) *http.Server {
	serverMux := http.NewServeMux()
	server := &http.Server{Addr: addr, Handler: serverMux}

	return server
}
