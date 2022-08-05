package eventloop

import (
	"fmt"
	"time"
)

type App struct {
}

func (app *App) exec() {
	for {
		fmt.Println("event loop is running")
		time.Sleep(time.Millisecond)
	}
}

var app *App

func InitAndRunApp() {
	initModules()
	app = &App{}
	app.exec()
}

func initModules() {
	initTimerManager()
	initApi()
}
