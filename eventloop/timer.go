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

func MakeTimerTask(interval int, callback func(int)) *TimerTask {
	timerTask := &TimerTask{interval, callback, time.Now().UnixMilli()}
	timerModule.addTask(timerTask)
	return timerTask
}

func MakeOneTimeTask(delay int, callback func(int)) *TimerTask {
	var timerTask *TimerTask
	wrapper := func(delay int) {
		callback(delay)
		timerModule.removeTask(timerTask)
	}
	timerTask = &TimerTask{delay, wrapper, time.Now().UnixMilli()}
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

func makeTimerResult(timerTask *TimerTask, dt int) *TimerTaskResult {
	return &TimerTaskResult{timerTask, dt}
}

type TimerModule struct {
	timerLock sync.Mutex
	timers    map[*TimerTask]struct{}
	BaseModule
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
		for timerTask := range timerModule.timers {
			timerModule.updateTimerTask(timerTask)
		}
		timerModule.timerLock.Unlock()
		time.Sleep(time.Millisecond)
	}
}

func (timerModule *TimerModule) updateTimerTask(timerTask *TimerTask) {
	currentTime := time.Now().UnixMilli()
	dt := timerTask.check(currentTime)
	if dt > 0 {
		timerResult := makeTimerResult(timerTask, dt)
		timerModule.events <- timerResult
	}
}

func GetTimerLength() int {
	timerModule.timerLock.Lock()
	defer timerModule.timerLock.Unlock()
	return len(timerModule.timers)
}

var timerModule *TimerModule

func initTimerModule(events chan IEvent) {
	timerModule = &TimerModule{timers: make(map[*TimerTask]struct{}), BaseModule: BaseModule{events}}
}

func startTimerModule() {
	go timerModule.exec()
}
