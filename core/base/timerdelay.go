package base

import (
	"reflect"
)

type DelayFunc struct {
	f func(...interface{}) 
	args []interface{} 
}

func NewDelayFunc(f func(v ...interface{}), args []interface{}) *DelayFunc {
	return &DelayFunc{
		f:f,
		args:args,
	}
}

func (df *DelayFunc) String() string {
	return Sprintf("{DelayFun:%s, args:%v}", reflect.TypeOf(df.f).Name(), df.args)
}

func (df *DelayFunc) Call() {
	defer func() {
		if err := recover(); err != nil {
			LogError(df.String(), "Call err: ", err)
		}
	}()

	//call func
	df.f(df.args...)
}
