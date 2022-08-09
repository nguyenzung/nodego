package eventloop

type App struct {
	events chan IEvent
}

func (app *App) exec() {
	for event := range app.events {
		event.process()
	}
}

var app *App

func NewApp() {
	events := make(chan IEvent, 1<<16)
	initModules(events)
	app = &App{events}
}

func RunApp() {
	startModules()
	app.exec()
}

func initModules(events chan IEvent) {
	initTimerModule(events)
	initAPICallModule(events)
	initHTTPServerModule(events)
	initWebsocketModule(events)
}

func startModules() {
	startTimerModule()
	startAPICallModule()
	startHTTPServerModule()
	startWebsocketModule()
}
