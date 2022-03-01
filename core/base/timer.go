package base

import (
	"time"
)

const (
	HOUR_NAME     = "HOUR"
	HOUR_INTERVAL = 60 * 60 * 1e3
	HOUR_SCALES   = 12

	MINUTE_NAME     = "MINUTE"
	MINUTE_INTERVAL = 60 * 1e3
	MINUTE_SCALES   = 60

	SECOND_NAME     = "SECOND"
	SECOND_INTERVAL = 1e3
	SECOND_SCALES   = 60

	TIMERS_MAX_CAP = 2048
)

/*
   	time.Second = time.Millisecond * 1e3
	time.Millisecond = time.Microsecond * 1e3
	time.Microsecond = time.Nanosecond * 1e3

	time.Now().UnixNano() ==> time.Nanosecond
*/

type Timer struct {
	delayFunc *DelayFunc
	times     int
	duration  int64
	unixMilli int64
	forever   bool
}

func NewTimerAt(unixNano int64, f func(v ...interface{}), args []interface{}) *Timer {
	df := NewDelayFunc(f, args)
	return &Timer{
		delayFunc: df,
		duration:  unixNano / 1e6,
		times:     1,
		unixMilli: unixNano / 1e6,
		forever:   false,
	}
}

func NewTimerAfter(duration time.Duration, f func(v ...interface{}), args []interface{}) *Timer {
	return NewTimerAt(time.Now().UnixNano()+int64(duration), f, args)
}

func NewTimerInterval(duration time.Duration, times int, f func(v ...interface{}), args []interface{}) *Timer {
	df := NewDelayFunc(f, args)
	now := time.Now().UnixNano()
	unixMilli := (now + int64(duration)) / 1e6
	fe := false
	LogInfo("============")
	if times == -1 {
		times = 1
		fe = true
	}
	return &Timer{
		delayFunc: df,
		duration:  int64(duration / 1e6),
		times:     times,
		unixMilli: unixMilli,
		forever:   fe,
	}
}

func (t *Timer) running() {
	now := GetUnixMs()
	if t.unixMilli > now {
		time.Sleep(time.Duration(t.unixMilli-now) * time.Millisecond)
	}
	t.delayFunc.Call()
	if t.forever {
		t.unixMilli += t.duration
		t.running()
	} else if t.times > 1 {
		t.unixMilli += t.duration
		t.times--
		t.running()
	}
}

func (t *Timer) Run() {
	go func() {
		t.running()
	}()
}
