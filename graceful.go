package grace

import (
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	_defaultOutTime  time.Duration = time.Second * 5
	_defaultOsSignal               = []os.Signal{syscall.SIGINT, syscall.SIGTERM, syscall.SIGTSTP}
)

type (
	Option       func(o *gracefulShut)
	OverHandler  func() error
	gracefulShut struct {
		chanSign chan os.Signal
		wg       *sync.WaitGroup
		outTime  time.Duration
		signal   []os.Signal
		handlers []OverHandler
	}
)

// WithHandlers add shutdown handler
func WithHandlers(h ...OverHandler) Option {
	return func(o *gracefulShut) {
		if len(h) > 0 {
			o.handlers = append(o.handlers, h...)
		}
	}
}

// WithOutTime define timeout. default use _defaultOutTime
func WithOutTime(d time.Duration) Option {
	return func(o *gracefulShut) {
		if d > 0 {
			o.outTime = d
		}
	}
}

// WithOsSignal custom exit signal
func WithOsSignal(sig []os.Signal) Option {
	return func(o *gracefulShut) {
		if len(sig) > 0 {
			o.signal = sig
		}
	}
}

// New graceful shutdown
func New(opts ...Option) *gracefulShut {
	var grace = &gracefulShut{
		chanSign: make(chan os.Signal, 1),
		wg:       new(sync.WaitGroup),
		outTime:  _defaultOutTime,
		signal:   _defaultOsSignal,
	}
	// set opt
	for _, opt := range opts {
		opt(grace)
	}

	// notify
	signal.Notify(grace.chanSign, grace.signal...)
	grace.wg.Add(1)
	go grace.startMonitorSignalNotify()
	return grace
}

func (m *gracefulShut) startMonitorSignalNotify() {
	for {
		select {
		case v, ok := <-m.chanSign:
			log.Printf("startMonitorSignalNotify chan signal return %+v. ok:%v", v, ok)
			m.wg.Done()
			return
		}
	}
}

func (m *gracefulShut) Wait() error {
	// wait
	m.wg.Wait()
	// time.out
	go func() {
		<-time.After(m.outTime)
		log.Printf("gracefull outTime over...")
		os.Exit(-1)
	}()
	// do action
	log.Printf("start do exit handlers...")
	for _, h := range m.handlers {
		if err := h(); err != nil {
			log.Printf("gracefulShut.Wait doHandler return err . err:%v", err)
			return err
		}
	}
	log.Printf("has over...")
	return nil
}

// Wait quick start
func Wait(opts ...Option) error {
	g := New(opts...)
	return g.Wait()
}
