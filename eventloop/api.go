package eventloop

import (
	"io/ioutil"
	"net/http"
	"time"
)

type APIProcessThread struct {
	callModule *APICallModule
}

func (taskThread *APIProcessThread) exec() {
	for {
		callTask := <-taskThread.callModule.tasks
		resp, err := callTask.send()
		callResult := newAPICallTaskResult(callTask, resp, err)
		taskThread.callModule.pushCallResult(callResult)
	}
}

func newAPIProcessThread(callModule *APICallModule) *APIProcessThread {
	return &APIProcessThread{callModule}
}

type APICallTask struct {
	url      string
	timeout  int
	callback func(string)
	err      func(error)
}

func (caller *APICallTask) send() (string, error) {
	client := http.Client{
		Timeout: time.Duration(caller.timeout) * time.Second,
	}
	resp, err := client.Get(caller.url)
	if err == nil {
		body, err := ioutil.ReadAll(resp.Body)
		if err == nil {
			return string(body), err
		} else {
			return "", err
		}
	}
	return "", err
}

func MakeCallTask(url string, timeout int, callback func(string), err func(error)) *APICallTask {
	callTask := &APICallTask{url, timeout, callback, err}
	api.tasks <- callTask
	return callTask
}

type APICallTaskResult struct {
	callTask *APICallTask
	resp     string
	err      error
}

func (callResult *APICallTaskResult) process() {
	if callResult.err == nil {
		callResult.callTask.callback(callResult.resp)
	} else {
		callResult.callTask.err(callResult.err)
	}
}

func newAPICallTaskResult(callTask *APICallTask, resp string, err error) *APICallTaskResult {
	return &APICallTaskResult{callTask, resp, err}
}

type APICallModule struct {
	BaseModule
	// events    chan IEvent
	tasks     chan *APICallTask
	numThread int
}

func (api *APICallModule) pushCallResult(callResult *APICallTaskResult) {
	api.events <- callResult
}

func (api *APICallModule) makeWorkerThread() {
	for i := 0; i < api.numThread; i++ {
		go newAPIProcessThread(api).exec()
	}
}

func (api *APICallModule) exec() {
	api.makeWorkerThread()
}

var api *APICallModule

func initAPICallModule(events chan IEvent) {
	api = &APICallModule{BaseModule: BaseModule{events}, numThread: API_NUM_THREAD, tasks: make(chan *APICallTask, API_NUM_THREAD)}
}

func startAPICallModule() {
	go api.exec()
}
