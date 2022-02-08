package tool

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	TryTimes       = 100 //
	ExpireInterval = 1   //s
	SleepInterval  = 10  //ms
)

type redisHand interface {
	Get(string, chan error) (string, bool)
	Delete(string, chan error) bool
	Put(string, string, time.Duration, chan error) bool
}

type Lock struct {
	channel   chan error
	id        string
	uuid      string
	redisHand redisHand
}

func NewLock(id string, redisHand redisHand) *Lock {
	l := &Lock{channel: make(chan error, 2), id: id, uuid: uuid.NewV4().String(), redisHand: redisHand}
	return l
}

func (l *Lock) GetErrChan() chan error {
	return l.channel
}

func (l *Lock) WaitLock() {
	l.tryLock(0)
}

func (l *Lock) ReleaseLock() {
	if cache_id, ok := l.redisHand.Get(l.getKey(), l.channel); ok && (cache_id == l.uuid || cache_id != "") {
		if ok := l.redisHand.Delete(l.getKey(), l.channel); ok {
			l.channel <- nil
		}
	}
}

func (l *Lock) getKey() string {
	return "Lock:" + l.id
}

func (l *Lock) tryLock(times int) {
	if times >= TryTimes {
		l.channel <- fmt.Errorf("TryLock over times please check")
		return
	} else {
		times++
	}
	if cache_id, ok := l.redisHand.Get(l.getKey(), l.channel); !ok {
		return
	} else if cache_id == l.uuid || cache_id == "" {
		if ok := l.redisHand.Put(l.getKey(), l.uuid, time.Duration(ExpireInterval)*time.Second, l.channel); !ok {
			return
		}
		l.channel <- nil
	} else {
		time.Sleep(time.Duration(SleepInterval) * time.Millisecond) //* chock wait next trylock
		l.tryLock(times)
	}
}
