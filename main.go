package main

import (
	ev "github.com/nguyenzung/nodego/eventloop"
)

func main() {
	app := ev.NewApp()
	app.Exec()
}
