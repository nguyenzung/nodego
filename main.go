package main

import (
	"errors"
	"fmt"
	"net/http"

	ev "github.com/nguyenzung/nodego/eventloop"
	"github.com/nguyenzung/nodego/runtimeutils"
)

func sum(arr ...int) (int, error) {
	fmt.Println("Sum function Thread ID", runtimeutils.ThreadID())
	if len(arr) == 0 {
		return -1, errors.New("should have params")
	}
	sum := 0
	for _, value := range arr {
		sum += value
	}
	return sum, nil
}

func main() {
	fmt.Println(" Init ")
	app := ev.NewApp()

	app.MakeAPIHandler("/test", func(hw *ev.HTTPResponseWriter, r *http.Request) {
		hw.SendText("123456")
	})

	app.MakeOneTimeTask(6000, func(i int) {
		fmt.Println(" Callback ")
	})

	task := app.MakeTask(sum,
		func(result int) {
			fmt.Println("Sum =", result, " ||| ThreadID", runtimeutils.ThreadID())
		},
		func(err error) {
			fmt.Println("Error", err, " ||| ThreadID", runtimeutils.ThreadID())
		})
	fmt.Println(task.Exec(1, 2, 3, 5, 6), runtimeutils.ThreadID())
	app.Exec()
}
