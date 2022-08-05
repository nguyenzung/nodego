package eventloop

type IResult interface {
	process()
}

type Result struct {
	status       bool
	resultObject IResult
}

func (result *Result) Status() bool {
	return result.status
}

func (result *Result) ResultObject() IResult {
	return result.resultObject
}

func MakeResult(status bool, resultObject IResult) *Result {
	return &Result{status, resultObject}
}
