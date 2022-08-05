package eventloop

import "time"

type ITimerTask interface {
	timeout(interval int)
}

type TimerTask struct {
	interval   int
	callback   func(int)
	latestTime int64
}

func NewTimerTask(interval int, callback func(int)) *TimerTask {
	timerTask := &TimerTask{interval, callback, time.Now().UnixMilli()}
	timerManager.addTask(timerTask)
	return timerTask
}

type TimerTaskResult struct {
}

type TimerManager struct {
	timers map[*TimerTask]struct{}
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

var timerManager *TimerManager

func initTimerManager() {
	timerManager = &TimerManager{make(map[*TimerTask]struct{})}
}
