package eventloop

type Api struct {
}

var api *Api

func initApi() {
	api = &Api{}
}
