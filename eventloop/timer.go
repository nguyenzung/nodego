package eventloop

import (
	"fmt"
	"time"
)

type TimerTask struct {
	interval   int
	callback   func(int)
	latestTime int64
}

func (timerTask *TimerTask) check(currentTime int64) int {
	dt := int(currentTime - timerTask.latestTime)
	if dt > timerTask.interval {
		timerTask.latestTime = currentTime
		return dt
	}
	return 0
}

func NewTimerTask(interval int, callback func(int)) *TimerTask {
	timerTask := &TimerTask{interval, callback, time.Now().UnixMilli()}
	timerManager.addTask(timerTask)
	fmt.Println(timerManager.timers)
	return timerTask
}

func RemoveTimerTask(timerTask *TimerTask) {
	timerManager.removeTask(timerTask)
}

type TimerTaskResult struct {
	timerTask *TimerTask
	dt        int
}

func (timerResult *TimerTaskResult) process() {
	fmt.Println(" process in timer ")
	timerResult.timerTask.callback(timerResult.dt)
}

func MakeTimerResult(timerTask *TimerTask, dt int) *TimerTaskResult {
	return &TimerTaskResult{timerTask, dt}
}

type TimerManager struct {
	timers map[*TimerTask]struct{}
	events chan IResult
}

func (timerManager *TimerManager) addTask(timerTask *TimerTask) {
	timerManager.timers[timerTask] = struct{}{}
}

func (timerManager *TimerManager) removeTask(timerTask *TimerTask) {
	_, ok := timerManager.timers[timerTask]
	if ok {
		delete(timerManager.timers, timerTask)
	}
}

func (timerManager *TimerManager) exec() {
	for {
		for timerTask, _ := range timerManager.timers {
			timerManager.updateTimerTask(timerTask)
		}
		time.Sleep(time.Millisecond)
	}
}

func (timerManager *TimerManager) updateTimerTask(timerTask *TimerTask) {
	// fmt.Println("Update Timer", timerTask.interval)
	currentTime := time.Now().UnixMilli()
	dt := timerTask.check(currentTime)
	if dt > 0 {
		timerResult := MakeTimerResult(timerTask, dt)
		timerManager.events <- timerResult
	}
}

var timerManager *TimerManager

func initTimerManager(events chan IResult) {
	timerManager = &TimerManager{make(map[*TimerTask]struct{}), events}
	go timerManager.exec()
}
