package eventloop

type HTTPServerModule struct {
	events chan IResult
}

func (http *HTTPServerModule) exec() {

}

var httpModule *HTTPServerModule

func initHTTPServerModule(events chan IResult) {
	httpModule = &HTTPServerModule{events}
	go httpModule.exec()
}
