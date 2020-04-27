package base

import (
	"math"
	"sync"
	"time"
)

const (
	MAX_CHAN_BUFF = 2048
	MAX_TIME_DELAY = 100
)

type TimerScheduler struct {
	tw *TimeWheel
	idGen uint32
	triggerChan chan *DelayFunc
	sync.RWMutex
}

func NewTimerScheduler() *TimerScheduler {

	secondTW := NewTimeWheel(SECOND_NAME, SECOND_INTERVAL, SECOND_SCALES, TIMERS_MAX_CAP)
	minuteTW := NewTimeWheel(MINUTE_NAME, MINUTE_INTERVAL, MINUTE_SCALES, TIMERS_MAX_CAP)
	hourTW := NewTimeWheel(HOUR_NAME, HOUR_INTERVAL, HOUR_SCALES, TIMERS_MAX_CAP)

	hourTW.AddTimeWheel(minuteTW)
	minuteTW.AddTimeWheel(secondTW)

	secondTW.Run()
	minuteTW.Run()
	hourTW.Run()

	return &TimerScheduler{
		tw:          hourTW,
		idGen:		 1,
		triggerChan: make(chan *DelayFunc, MAX_CHAN_BUFF),
	}
}

func (this *TimerScheduler) CreateTimerAt(unixNano int64, f func(v ...interface{}), args []interface{})(uint32, error) {
	this.Lock()
	defer this.Unlock()

	this.idGen++

	return this.idGen, this.tw.AddTimer(this.idGen, NewTimerAt(unixNano, f, args))
}

func (this *TimerScheduler) NewTimerAfter(duration time.Duration, f func(v ...interface{}), args []interface{})(uint32, error) {
	this.Lock()
	defer this.Unlock()

	this.idGen++
	return this.idGen, this.tw.AddTimer(this.idGen, NewTimerAfter(duration, f, args))
}

func  (this *TimerScheduler) NewTimerInterval(duration time.Duration, times int, f func(v ...interface{}), args []interface{})(uint32, error){
	this.Lock()
	defer this.Unlock()

	this.idGen++
	return this.idGen, this.tw.AddTimer(this.idGen, NewTimerInterval(duration, times, f, args))
}

func(this *TimerScheduler) CancelTimer(tid uint32) {
	this.Lock()
	this.Unlock()

	//Println("CancelTimer:",tid)
	this.tw.RemoveTimer(tid)
}

func (this *TimerScheduler) GetTriggerChan() chan *DelayFunc {
	return this.triggerChan
}

func (this *TimerScheduler) Start() {
	go func() {
		for {
			now := GetUnixMs()
			timerList := this.tw.GetTimerWithIn(MAX_TIME_DELAY * time.Millisecond)
			for _, timer := range timerList {
				if math.Abs(float64(now-timer.unixMilli)) > MAX_TIME_DELAY {
					LogWarn("want call at ", timer.unixMilli, "; real call at", now, "; delay ", now-timer.unixMilli)
				}
				this.triggerChan <- timer.delayFunc
			}

			time.Sleep(MAX_TIME_DELAY/2 * time.Millisecond)
		}
	}()
}

func NewAutoExecTimerScheduler() *TimerScheduler{
	autoExecScheduler := NewTimerScheduler()
	autoExecScheduler.Start()

	go func() {
		delayFuncChan := autoExecScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			go df.Call()
		}
	}()

	return autoExecScheduler
}