package tool

import (
	"context"
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	TryTimes       = 100  //
	ExpireInterval = 1000 //ms
	SleepInterval  = 10   //ms
)

type redisHand interface {
	Get(string, chan error) (string, bool)
	Delete(string, chan error) bool
	Set(string, string, time.Duration, chan error) bool
}

type config struct {
	tryTimes       int
	expireInterval int
	sleepInterval  int
}

func newConfig() config {
	return config{
		tryTimes:       TryTimes,
		expireInterval: ExpireInterval,
		sleepInterval:  SleepInterval,
	}
}

type Lock struct {
	channel   chan error
	id        string
	uuid      string
	redisHand redisHand
	config    config
}

func NewLock(id string, redisHand redisHand) *Lock {
	l := &Lock{channel: make(chan error, 2), id: id, uuid: uuid.NewV4().String(), redisHand: redisHand, config: newConfig()}
	return l
}

//default TryTimes       = 100 ExpireInterval = 1000   //ms SleepInterval  = 10  //ms
func (l *Lock) SetConfig(tryTimes, expireInterva, sleepInterval int) *Lock {
	l.config.tryTimes = tryTimes
	l.config.expireInterval = expireInterva
	l.config.sleepInterval = sleepInterval
	return l
}

func (l *Lock) Result(ctx context.Context) error {
	select {
	case err := <-l.channel:
		return err
	case <-ctx.Done():
		return fmt.Errorf("ctx has done")
	}
}

func (l *Lock) GetErrChan() chan error {
	return l.channel
}

func (l *Lock) WaitLock() {
	l.lock(0, true)
}

func (l *Lock) TryLock() {
	l.lock(l.config.tryTimes-1, false)
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

func (l *Lock) lock(times int, wait bool) {
	if times >= l.config.tryTimes {
		l.channel <- fmt.Errorf("TryLock over times please check")
		return
	} else {
		times++
	}
	if cache_id, ok := l.redisHand.Get(l.getKey(), l.channel); !ok {
		return
	} else if cache_id == l.uuid || cache_id == "" {
		if ok := l.redisHand.Set(l.getKey(), l.uuid, time.Duration(l.config.expireInterval)*time.Millisecond, l.channel); !ok {
			return
		}
		l.channel <- nil
	} else {
		if wait {
			time.Sleep(time.Duration(l.config.sleepInterval) * time.Millisecond) //* chock wait next trylock
		}
		l.lock(times, wait)
	}
}
