package main

import (
	"fmt"

	"github.com/go-event-loop/eventloop"
	"golang.org/x/sys/unix"
)

func timeout1(interval int) {
	pid := unix.Getpid()
	fmt.Println("Timeout 1 is called: ", interval, pid)
}

func timeout2(interval int) {
	pid := unix.Getpid()
	fmt.Println("Timeout 2 is called: ", interval, pid)
}

func main() {
	eventloop.NewApp()
	pid := unix.Getpid()
	fmt.Println(" MainThread ID: ", pid)
	eventloop.NewTimerTask(3000, timeout1)
	eventloop.NewTimerTask(5000, timeout2)
	eventloop.RunApp()
}
