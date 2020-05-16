package db

import (
	"errors"
	//"io"
	//"net"
	"strings"
	"sync"
	"sync/atomic"

	"github.com/go-redis/redis"
	base "github.com/shiniu0606/engine/core/base"
)

var (
	scriptMap             = sync.Map{} //map[int]string{}
	scriptCommitMap       = sync.Map{} //map[int]string{}
	scriptHashMap         = sync.Map{} //map[int]string{}
	scriptIndex     int32 = 0

	ErrRedisErr = errors.New("Redis Script Error")
)

func NewRedisScript(commit, str string) int {
	cmd := int(atomic.AddInt32(&scriptIndex, 1))
	scriptMap.Store(cmd, str)
	scriptCommitMap.Store(cmd, commit)
	return cmd
}

func GetRedisScript(cmd int) string {
	if s, ok := scriptMap.Load(cmd); ok {
		return s.(string)
	}
	return ""
}

type RedisConfig struct {
	Addr     string
	Passwd   string
	PoolSize int
}

type RedisManager struct {
	dbs      map[int]*Redis
	subMap   map[string]*Redis
	channels []string
	fun      func(channel, data string)
	lock     sync.RWMutex
}

type Redis struct {
	*redis.Client
	pubsub  *redis.PubSub
	conf    *RedisConfig
	manager *RedisManager
}

func (r *Redis) ScriptStr(cmd int, keys []string, args ...interface{}) (string, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		base.LogError("redis script failed err:%v", err)
		return "", errors.New("Redis ScriptStr Error")
	}

	if data == nil {
		return "", nil
	}
	str, ok := data.(string)
	if !ok {
		return "", errors.New("Redis ScriptStr Data Error")
	}

	return str, nil
}

func (r *Redis) ScriptInt64(cmd int, keys []string, args ...interface{}) (int64, error) {
	data, err := r.Script(cmd, keys, args...)
	if err != nil {
		base.LogError("redis script failed err:%v", err)
		return 0, errors.New("Redis ScriptInt64 Error")
	}
	if data == nil {
		return 0, nil
	}
	code, ok := data.(int64)
	if ok {
		return code, nil
	}
	return 0, errors.New("Redis ScriptInt64 Data Error")
}

func (r *Redis) Script(cmd int, keys []string, args ...interface{}) (interface{}, error) {
	var err error = ErrRedisErr
	var re interface{}
	hashStr, ok := scriptHashMap.Load(cmd)
	if ok {
		re, err = r.EvalSha(hashStr.(string), keys, args...).Result()
	}
	if err != nil {
		scriptStr, ok := scriptMap.Load(cmd)
		if !ok {
			base.LogError("redis script error cmd not found cmd:%v", cmd)
			return nil, err
		}
		cmdStr, _ := scriptCommitMap.Load(cmd)
		if strings.HasPrefix(err.Error(), "NOSCRIPT ") || err == ErrRedisErr {
			base.LogInfo("try reload redis script %v", cmdStr.(string))
			hashStr, err = r.ScriptLoad(scriptStr.(string)).Result()
			if err != nil {
				base.LogError("redis script load cmd:%v errstr:%s", cmdStr.(string), err)
				return nil, err
			}
			scriptHashMap.Store(cmd, hashStr.(string))
			re, err = r.EvalSha(hashStr.(string), keys, args...).Result()
			if err == nil {
				return re, nil
			}
		}
		base.LogError("redis script error cmd:%v errstr:%s", cmdStr.(string), err)
		return nil, ErrRedisErr
	}

	return re, nil
}