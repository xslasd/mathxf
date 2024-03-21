package mathxf

import (
	"fmt"
	"github.com/pkg/errors"
	"strconv"
)

type ECodes interface {
	// Error sometimes Error return Status in string form
	// NOTE: don't use Error in monitor report even it also work for now
	Error() string
	// Code Status get error code.
	Code() int
	// Message get code message.
	Message() string
	//Values get formatted message parameter,it may be nil.
	Values() []any
	// Details Detail get error detail,it may be nil.
	Details() []any

	Position() (line int, col int)

	SetCol(col int) ECodes
	SetPosition(line int, col int) ECodes
	SetMessage(msg string) ECodes
	SetMessagef(a ...any) ECodes
	WithDetails(msg any) ECodes
}

type ECode struct {
	id     int
	col    int
	line   int
	msg    string
	f      []any
	detail []any
}

func (e *ECode) Error() string {
	return strconv.Itoa(e.id)
}

func (e *ECode) Code() int {
	return e.id
}

func (e *ECode) Message() string {
	if e.f != nil {
		return fmt.Sprintf(e.msg, e.f...)
	}
	return e.msg
}

func (e *ECode) Values() []any {
	return e.f
}

func (e *ECode) Details() []any {
	return e.detail
}

func (e *ECode) Position() (line int, col int) {
	return e.line, e.col
}
func (e *ECode) SetCol(col int) ECodes {
	return &ECode{id: e.id, col: col, line: e.line, msg: e.msg, f: e.f, detail: e.detail}
}
func (e *ECode) SetPosition(line int, col int) ECodes {
	return &ECode{id: e.id, col: col, line: line, msg: e.msg, f: e.f, detail: e.detail}
}

func (e *ECode) SetMessage(msg string) ECodes {
	return &ECode{id: e.id, col: e.col, line: e.line, msg: msg, f: nil, detail: e.detail}
}

func (e *ECode) SetMessagef(f ...any) ECodes {
	return &ECode{id: e.id, col: e.col, line: e.line, msg: e.msg, f: f, detail: e.detail}
}

func (e *ECode) WithDetails(msg any) ECodes {
	return &ECode{id: e.id, col: e.col, line: e.line, msg: e.msg, f: e.f, detail: append(e.detail, msg)}
}

var (
	_codes = make(map[int]string, 0)
)

func New(e int, msg string) ECodes {
	if _, ok := _codes[e]; ok {
		panic(fmt.Sprintf("ecode: %d already exist", e))
	}
	_codes[e] = msg
	return &ECode{
		id:  e,
		msg: msg,
	}
}

func Cause(e error) ECodes {
	if e == nil {
		return nil
	}
	ec, ok := errors.Cause(e).(ECodes)
	fmt.Println("---------Cause---------", ec, ok)
	if ok {
		return ec
	}
	return ServerErr.WithDetails(e.Error())
}

func As(e error) ECodes {
	if e == nil {
		return nil
	}
	var c *ECode
	if errors.As(e, &c) {
		return c
	}
	return ServerErr.WithDetails(e.Error())
}
