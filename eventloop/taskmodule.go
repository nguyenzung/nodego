package eventloop

import (
	"errors"
	"fmt"
	"reflect"
)

type TaskProcessThread struct {
	taskModule *TaskModule
}

func (taskThread *TaskProcessThread) exec() {
	for {
		task := <-taskThread.taskModule.taskChannel
		fmt.Println(task)
		resp, err := task.handlerExec()
		callResult := makeTaskResult(task, resp, err)
		taskThread.taskModule.pushCallResult(callResult)
	}
}

func newTaskProcessThread(taskModule *TaskModule) *TaskProcessThread {
	return &TaskProcessThread{taskModule}
}

type CustomizedTask struct {
	taskModule *TaskModule
	handler    interface{}
	data       []interface{}
	callback   interface{}
	err        interface{}
}

func (task *CustomizedTask) handlerExec() (interface{}, interface{}) {
	params := make([]reflect.Value, len(task.data))
	for i, value := range task.data {
		params[i] = reflect.ValueOf(value)
	}
	value := reflect.ValueOf(task.handler).Call(params)
	return value[0].Interface(), value[1].Interface()
}

func (task *CustomizedTask) Exec(a ...interface{}) error {
	if len(task.data) == 0 {
		task.data = a
		task.taskModule.addTask(task)
		return nil
	} else {
		return errors.New("[ERROR] have executed")
	}
}

type CustomizedTaskResult struct {
	task *CustomizedTask
	resp interface{}
	err  interface{}
}

func (taskResult *CustomizedTaskResult) callbackExec() {
	reflect.ValueOf(taskResult.task.callback).Call([]reflect.Value{reflect.ValueOf(taskResult.resp)})
}

func (taskResult *CustomizedTaskResult) errorExec() {
	reflect.ValueOf(taskResult.task.err).Call([]reflect.Value{reflect.ValueOf(taskResult.err)})
}

func (taskResult *CustomizedTaskResult) process() {
	if taskResult.err == nil {
		taskResult.callbackExec()
	} else {
		taskResult.errorExec()
	}
}

func makeTaskResult(task *CustomizedTask, resp interface{}, err interface{}) *CustomizedTaskResult {
	return &CustomizedTaskResult{task: task, resp: resp, err: err}
}

type TaskModule struct {
	BaseModule
	numThread   int
	taskChannel chan *CustomizedTask
}

func (taskModule *TaskModule) makeTask(handler interface{}, callback interface{}, err interface{}) *CustomizedTask {
	return &CustomizedTask{taskModule: taskModule, handler: handler, callback: callback, err: err}
}

func (taskModule *TaskModule) pushCallResult(taskResult *CustomizedTaskResult) {
	taskModule.events <- taskResult
}

func (taskModule *TaskModule) addTask(task *CustomizedTask) {
	taskModule.taskChannel <- task
}

func (taskModule *TaskModule) makeWorkerThread() {
	for i := 0; i < taskModule.numThread; i++ {
		go newTaskProcessThread(taskModule).exec()
	}
}

func (taskModule *TaskModule) exec() {
	taskModule.makeWorkerThread()
}

func makeTaskModule(events chan IEvent) *TaskModule {
	return &TaskModule{BaseModule: BaseModule{events}, numThread: TASK_NUM_THREAD, taskChannel: make(chan *CustomizedTask, TASK_NUM_THREAD)}
}
