package base

import (
	"errors"
	"sync"
	"time"
)

type TimeWheel struct {
	name 			string
	interval 		int64
	scales 			int
	curIndex 		int
	maxCap 			int
	timerQueue 		map[int]map[uint32]*Timer 
	nextTimeWheel 	*TimeWheel
	sync.RWMutex
}

func NewTimeWheel(name string, interval int64, scales int, maxCap int) *TimeWheel {
	tw := &TimeWheel{
		name:       name,
		interval:   interval,
		scales:     scales,
		maxCap:     maxCap,
		timerQueue: make(map[int]map[uint32]*Timer, scales),
	}
	for i := 0; i < scales; i++ {
		tw.timerQueue[i] = make(map[uint32]*Timer, maxCap)
	}
	return tw
}

func (tw *TimeWheel) AddTimeWheel(next *TimeWheel) {
	tw.nextTimeWheel = next
	//Println("Add timerWheel[", tw.name,"]'s next [", next.name,"] is succ!")
}

func (tw *TimeWheel) RemoveTimer(tid uint32) {
	tw.Lock()
	defer tw.Unlock()

	for i := 0; i < tw.scales; i++ {
		if _, ok := tw.timerQueue[i][tid]; ok {
			delete(tw.timerQueue[i], tid)
		}
	}
}

func (tw *TimeWheel) AddTimer(tid uint32, t *Timer) error {
	tw.Lock()
	defer tw.Unlock()

	return tw.addTimer(tid, t, false)
}

func (tw *TimeWheel) addTimer(tid uint32, t *Timer, forceNext bool) error {
	defer func() error {
		if err := recover(); err != nil {
			errStr := Sprintf("addTimer function err : %s", err)
			Println(errStr)
			return errors.New(errStr)
		}
		return nil
	}()

	delayInterval := t.unixMilli - GetUnixMs()

	if delayInterval >= tw.interval {
		dn := delayInterval / tw.interval
		tw.timerQueue[(tw.curIndex+int(dn))%tw.scales][tid] = t
		return nil
	}

	if delayInterval < tw.interval && tw.nextTimeWheel == nil {
		if forceNext == true {
			tw.timerQueue[(tw.curIndex+1) % tw.scales][tid] = t
		} else {
			tw.timerQueue[tw.curIndex][tid] = t
		}
		return nil
	}

	if delayInterval < tw.interval {
		return tw.nextTimeWheel.AddTimer(tid, t)
	}

	return nil
}

func (tw *TimeWheel) run() {
	for {
		time.Sleep(time.Duration(tw.interval) * time.Millisecond)

		tw.Lock()
		curTimers := tw.timerQueue[tw.curIndex]
		tw.timerQueue[tw.curIndex] = make(map[uint32]*Timer, tw.maxCap)
		for tid, timer := range curTimers {
			tw.addTimer(tid, timer, true)
		}

		nextTimers := tw.timerQueue[(tw.curIndex+1) % tw.scales]
		tw.timerQueue[(tw.curIndex+1) % tw.scales] = make(map[uint32]*Timer, tw.maxCap)
		for tid, timer := range nextTimers {
			tw.addTimer(tid, timer, true)
		}

		tw.curIndex = (tw.curIndex+1) % tw.scales

		tw.Unlock()
	}
}

func (tw *TimeWheel) Run() {
	go tw.run()
	//Println("timerWheel name = ", tw.name, " is running...")
}

func (tw *TimeWheel) GetTimerWithIn(duration time.Duration) map[uint32]Timer {
	leafTW := tw
	for leafTW.nextTimeWheel != nil {
		leafTW = leafTW.nextTimeWheel
	}

	leafTW.Lock()
	defer leafTW.Unlock()
	timerList := make(map[uint32]Timer)

	now := GetUnixMs()

	for tid, timer := range leafTW.timerQueue[leafTW.curIndex] {
		if timer.unixMilli-now < int64(duration/1e6) {
			//time out
			timerList[tid] = *timer
			if timer.forever {
				timer.unixMilli += timer.duration
			}else if timer.times <= 1{
				delete(leafTW.timerQueue[leafTW.curIndex], tid)
			}else{
				timer.times--
				timer.unixMilli += timer.duration
			}
		}
	}

	return timerList
}
