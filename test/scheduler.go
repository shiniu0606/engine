package main

import (
	"time"
	base "github.com/shiniu0606/engine/core/base"
)

func foo(args ...interface{}){
	base.Printf("I am No. %d function, delay %d s\n", args[0].(int), args[1].(int))
}

func main() {

	timerScheduler := base.NewTimerScheduler()
	timerScheduler.Start()

	//for i := 1; i < 3; i ++ {
		tid, err := timerScheduler.NewTimerInterval(time.Duration(1)*time.Second,-1, foo, []interface{}{0, 2})
		//tid, err := timerScheduler.NewTimerAfter(time.Duration(3*i) * time.Millisecond, foo, []interface{}{i, i*3})
		if err != nil {
			base.Println("create timer error", tid, err)
			return
		}
	//}

	go func() {
		delayFuncChan := timerScheduler.GetTriggerChan()
		for df := range delayFuncChan {
			df.Call()
		}
	}()

	select{}
}