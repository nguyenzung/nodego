package eventloop

import (
	"net/http"
)

type App struct {
	api   *APICallModule
	timer *TimerModule
	http  *HTTPServerModule
	ws    *WebsocketModule
	tasks *TaskModule

	events chan IEvent
}

func (app *App) exec() {
	for event := range app.events {
		event.process()
	}
}

func (app *App) MakeCallTask(url string, timeout int, callback func(string), err func(error)) *APICallTask {
	return app.api.makeCallTask(url, timeout, callback, err)
}

func (app *App) MakeTimerTask(interval int, callback func(int)) *TimerTask {
	return app.timer.makeTimerTask(interval, callback)
}

func (app *App) MakeOneTimeTask(delay int, callback func(int)) *TimerTask {
	return app.timer.makeOneTimeTask(delay, callback)
}

func (app *App) RemoveTimerTask(timerTask *TimerTask) {
	app.timer.removeTimerTask(timerTask)
}

func (app *App) MakeAPIHandler(path string, handler func(*HTTPResponseWriter, *http.Request)) {
	app.http.makeAPIHandler(path, handler)
}

func (app *App) MakeWSHandler(path string, openHandler func(*Session), messageHandler func(*MessageEvent, *Session), closeHandler func(*CloseEvent, *Session) error) {
	app.ws.makeWSHandler(path, openHandler, messageHandler, closeHandler)
}

func (app *App) MakeTask(handler interface{}, callback interface{}, err interface{}) *CustomizedTask {
	return app.tasks.makeTask(handler, callback, err)
}

func (app *App) initModules(events chan IEvent) {
	app.timer = makeTimerModule(events)
	app.api = makeAPICallModule(events)
	app.http = makeHTTPServerModule(events)
	app.ws = makeWebsocketModule(events)
	app.tasks = makeTaskModule(events)
}

func (app *App) startModules() {
	go app.timer.exec()
	go app.api.exec()
	go app.http.exec()
	go app.ws.exec()
	go app.tasks.exec()
}

func (app *App) Exec() {
	app.startModules()
	app.exec()
}

func NewApp() *App {
	events := make(chan IEvent, 1<<16)
	app := &App{events: events}
	app.initModules(events)
	return app
}
