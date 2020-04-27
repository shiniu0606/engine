package base

import (
	"time"
)

const(
	INTERVAL_DAY = iota		//当前时间的第二天
	INTERVAL_WEEK = iota		//当前时间的下周一
	INTERVAL_MONTH = iota	//当前时间的下一月第一天)
	INTERVAL_YEAR   = iota   //当前时间的下一年第一天)
	TIME_SET_MAX_VAL = iota //类型最大值
)

var StartTick int64
var NowTick int64
var Timestamp int64

func init() {

}

func timerTick() {
	StartTick = time.Now().UnixNano() / 1000000
	NowTick = StartTick
	Timestamp = NowTick / 1000
	var ticker = time.NewTicker(time.Millisecond)
	Go(func() {
		for {
			select {
			case <-ticker.C:
				//LogInfo("time tick")
				NowTick = time.Now().UnixNano() / 1000000
				Timestamp = NowTick / 1000
			}
		}
	})
}

func GetTimestamp() int64 {
	return Timestamp
}

func GetNextTime(intervalType int) time.Time{
	t := time.Now()
	if intervalType == INTERVAL_YEAR{
		t = t.AddDate(1, 0, 0)
	}else if intervalType == INTERVAL_MONTH{
		t = t.AddDate(0, 1, 0)
	}else if intervalType == INTERVAL_WEEK{
		if t.Weekday() != time.Sunday{
			t = t.AddDate(0, 0, (8-int(t.Weekday())))
		}else{
			t = t.AddDate(0,0, 1)
		}
	}else{
		t = t.AddDate(0,0, 1)
	}

	DefaultTimeLoc := time.Local
	return  time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, DefaultTimeLoc)
}

func GetDate() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func GetUnixTime(sec, nsec int64) time.Time {
	return time.Unix(sec, nsec)
}

func GetUnixMs() int64 {
	return time.Now().UnixNano() / 1000000
}

func GetUnixNano() int64 {
	return time.Now().UnixNano()
}

func GetNowTime() time.Time {
	return time.Now()
}

func GetNextHourIntervalS(timestamp int64) int {
	return int(3600 - (timestamp % 3600))
}

func GetHour24(timestamp int64) int {
	hour := (int(timestamp%86400)/3600)
	if hour > 24 {
		return hour - 24
	}
	return hour
}

func GetHour(timestamp int64) int {
	hour := GetHour24(timestamp)
	if hour == 24 {
		return 0 //24点就是0点
	}
	return hour
}

func GetYearMonthDay(timestamp int64) (int32, int32, int32) {
	year, month, day := time.Unix(timestamp, 0).UTC().Date()
	return int32(year), int32(month), int32(day)
}

func IsDiffDay(now, old int64) int {
	return int((now / 86400) - (old / 86400))
}

func NewTimer(ms int) *time.Timer {
	return time.NewTimer(time.Millisecond * time.Duration(ms))
}

func NewTicker(ms int) *time.Ticker {
	return time.NewTicker(time.Millisecond * time.Duration(ms))
}

func ParseTime(str string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", str)
}

func After(ms int) <-chan time.Time {
	return time.After(time.Millisecond * time.Duration(ms))
}

func Tick(ms int) <-chan time.Time {
	return time.Tick(time.Millisecond * time.Duration(ms))
}

func Sleep(ms int) {
	time.Sleep(time.Millisecond * time.Duration(ms))
}

