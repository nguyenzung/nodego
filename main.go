package main

import (
	"fmt"
	"net/http"

	ev "github.com/nguyenzung/nodego/eventloop"
	"github.com/nguyenzung/nodego/threadutils"
)

func main() {
	app := ev.NewApp()
	app.Exec()
}
