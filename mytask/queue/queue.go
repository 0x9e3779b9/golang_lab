package queue

import (
	"redis"
	"time"
)

var defaultQ *redis.Redis

func Init() (err error) {
	defaultQ, err = redis.NewRedis("127.0.0.1:6379")
	return
}

func Pop(qName string) (v interface{}, err error) {
	return defaultQ.Rpop(qName)
}

func Push(qName string, data interface{}) (err error) {
	return defaultQ.Lpush(qName, data)
}
