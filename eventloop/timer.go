package eventloop

import (
	"sync"
	"time"
)

type TimerTask struct {
	interval   int
	callback   func(int)
	latestTime int64
	isOneTime  bool
}

func (timerTask *TimerTask) check(currentTime int64) int {
	dt := int(currentTime - timerTask.latestTime)
	if dt > timerTask.interval {
		timerTask.latestTime = currentTime
		return dt
	}
	return 0
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

func (timerModule *TimerModule) makeTimerTask(interval int, callback func(int)) *TimerTask {
	timerTask := &TimerTask{interval, callback, time.Now().UnixMilli(), false}
	timerModule.addTask(timerTask)
	return timerTask
}

func (timerModule *TimerModule) makeOneTimeTask(interval int, callback func(int)) *TimerTask {
	timerTask := &TimerTask{interval, callback, time.Now().UnixMilli(), true}
	timerModule.addTask(timerTask)
	return timerTask
}

func (timerModule *TimerModule) removeTimerTask(timerTask *TimerTask) {
	timerModule.removeTask(timerTask)
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

func (timerModule *TimerModule) updateTimerTask(timerTask *TimerTask) {
	currentTime := time.Now().UnixMilli()
	dt := timerTask.check(currentTime)
	if dt > 0 {
		timerResult := makeTimerResult(timerTask, dt)
		timerModule.events <- timerResult
		if timerTask.isOneTime {
			timerModule.removeTimerTask(timerTask)
		}
	}
}

func (timerModule *TimerModule) GetTimerLength() int {
	timerModule.timerLock.Lock()
	defer timerModule.timerLock.Unlock()
	return len(timerModule.timers)
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

func makeTimerModule(events chan IEvent) *TimerModule {
	return &TimerModule{timers: make(map[*TimerTask]struct{}), BaseModule: BaseModule{events}}
}
