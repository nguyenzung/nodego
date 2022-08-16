package main

import (
	"errors"
	"fmt"

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
	app := ev.NewApp()

	app.MakeCallTask("http://localhost:8080", 9, func(s string) { fmt.Println(s) }, func(err error) { fmt.Println(err) })

	task := app.MakeTask(sum,
		func(result int) {
			fmt.Println("Sum =", result, " ||| ThreadID", runtimeutils.ThreadID())
		},
		func(err error) {
			fmt.Println("Error", err, " ||| ThreadID", runtimeutils.ThreadID())
		})
	fmt.Println(task.Exec(), runtimeutils.ThreadID())
	app.Exec()
}
