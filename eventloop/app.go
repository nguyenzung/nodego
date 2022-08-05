package eventloop

type App struct {
	events chan IResult
}

func (app *App) exec() {
	for event := range app.events {
		event.process()
	}
}

var app *App

func NewApp() {
	events := make(chan IResult, 1<<16)
	initModules(events)
	app = &App{events}
}

func RunApp() {
	app.exec()
}

func initModules(events chan IResult) {
	initTimerModule(events)
	initAPICallModule(events)
}
