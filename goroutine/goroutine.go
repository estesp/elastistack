package goroutine

import (
	"fmt"
	"time"

	"github.com/maruel/panicparse/stack"
)

type CallInfo struct {
	SourcePath string `json:"source_path"`
	LineNo     int    `json:"line_num"`
	FuncName   string `json:"func_name"`
	ArgList    string `json:"arg_list"`
}

type Stack struct {
	Calls  []CallInfo `json:"calls"`
	Elided bool       `json:"elided"`
}

type GoroutineTrace struct {
	EventTime time.Time `json:"event_time"`
	CreatedBy CallInfo  `json:"created_by"`
	ID        string    `json:"id"`
	State     string    `json:"state"`
	Locked    bool      `json:"locked"`
	SleepMin  int       `json:"sleep_min"`
	CallStack Stack     `json:"call_stack"`
}

func NewGoroutineTrace(routine stack.Goroutine, eventTime time.Time) *GoroutineTrace {
	//translate a goroutine from the stack parsing code
	//into a simpler object for use as JSON input to elasticsearch
	routineEntry := &GoroutineTrace{
		EventTime: eventTime,
		ID:        fmt.Sprintf("%d", routine.ID),
		State:     routine.State,
		SleepMin:  routine.SleepMin,
		Locked:    routine.Locked,
	}

	routineEntry.CreatedBy = CallInfo{
		SourcePath: routine.CreatedBy.SourcePath,
		LineNo:     routine.CreatedBy.Line,
		FuncName:   routine.CreatedBy.Func.String(),
		ArgList:    routine.CreatedBy.Args.String(),
	}

	for _, stackCall := range routine.Stack.Calls {
		callEntry := CallInfo{
			SourcePath: stackCall.SourcePath,
			LineNo:     stackCall.Line,
			FuncName:   stackCall.Func.String(),
			ArgList:    stackCall.Args.String(),
		}
		routineEntry.CallStack.Calls = append(routineEntry.CallStack.Calls, callEntry)
	}
	routineEntry.CallStack.Elided = routine.Stack.Elided

	return routineEntry
}
