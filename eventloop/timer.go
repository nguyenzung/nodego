package eventloop

import (
	"sync"
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
	timerModule.addTask(timerTask)
	return timerTask
}

func RemoveTimerTask(timerTask *TimerTask) {
	timerModule.removeTask(timerTask)
}

type TimerTaskResult struct {
	timerTask *TimerTask
	dt        int
}

func (timerResult *TimerTaskResult) process() {
	timerResult.timerTask.callback(timerResult.dt)
}

func MakeTimerResult(timerTask *TimerTask, dt int) *TimerTaskResult {
	return &TimerTaskResult{timerTask, dt}
}

type TimerModule struct {
	timerLock sync.Mutex
	timers    map[*TimerTask]struct{}
	events    chan IResult
}

func (timerModule *TimerModule) addTask(timerTask *TimerTask) {
	timerModule.timerLock.Lock()
	timerModule.timers[timerTask] = struct{}{}
	timerModule.timerLock.Unlock()
}

func (timerModule *TimerModule) removeTask(timerTask *TimerTask) {
	timerModule.timerLock.Lock()
	_, ok := timerModule.timers[timerTask]
	if ok {
		delete(timerModule.timers, timerTask)
	}
	timerModule.timerLock.Unlock()
}

func (timerModule *TimerModule) exec() {
	for {
		timerModule.timerLock.Lock()
		for timerTask, _ := range timerModule.timers {
			timerModule.updateTimerTask(timerTask)
		}
		timerModule.timerLock.Unlock()
		time.Sleep(time.Millisecond)
	}
}

func (timerModule *TimerModule) updateTimerTask(timerTask *TimerTask) {
	// fmt.Println("Update Timer", timerTask.interval)
	currentTime := time.Now().UnixMilli()
	dt := timerTask.check(currentTime)
	if dt > 0 {
		timerResult := MakeTimerResult(timerTask, dt)
		timerModule.events <- timerResult
	}
}

var timerModule *TimerModule

func initTimerModule(events chan IResult) {
	timerModule = &TimerModule{timers: make(map[*TimerTask]struct{}), events: events}
	go timerModule.exec()
}
