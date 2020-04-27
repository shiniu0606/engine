package base

import (
	"errors"
	"sync"
	"time"
)

//Snowflake 算法
const (
	twepoch        = int64(1483228800000)             //开始时间截 (2017-01-01)
	workeridBits   = uint(10)                         //机器id所占的位数
	sequenceBits   = uint(12)                         //序列所占的位数
	workeridMax    = int64(-1 ^ (-1 << workeridBits)) //支持的最大机器id数量
	sequenceMask   = int64(-1 ^ (-1 << sequenceBits)) //
	workeridShift  = sequenceBits                     //机器id左移位数
	timestampShift = sequenceBits + workeridBits      //时间戳左移位数
	WorkeridMax    = workeridMax					  //集群自增量
)

type Snowflake struct {
	sequence int64
	workerid int64
	timestamp int64
	sync.Mutex
}

type ISnowflake interface {
	Init(workerid int64)
	UUID() int64
}

// 实例化一个工作节点
func NewSnowflake(workerId int64) (*Snowflake, error) {
    if workerId < 0 || workerId > WorkeridMax {
        return nil, errors.New("Snowflake worker ID excess of quantity")
    }

    return &Snowflake{
        timestamp: 0,
        workerid: workerId,
        sequence: 0,
    }, nil
}

func (this *Snowflake) Init(workerid int64){
	if workerid < 0 || workerid > workeridMax {
		LogError("workerid must be between 0 and 1023")
		return
	}

	this.workerid = workerid
	LogDebug("snowflake [  workid : ", workerid, "]")
}

// Generate creates and returns a unique snowflake ID
func (s *Snowflake) UUID() int64 {
	s.Lock()
	now := time.Now().UnixNano() / 1000000
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask

		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixNano() / 1000000
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now
	r := int64((now-twepoch)<<timestampShift | (s.workerid << workeridShift) | (s.sequence))
	s.Unlock()
	return r
}
