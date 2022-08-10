package eventloop

type IEvent interface {
	process()
}
