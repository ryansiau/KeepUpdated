package looper

import (
	"time"

	"github.com/ryansiau/KeepUpdated/go/pkg/graceful-shutdown"
)

type Looper[T any] interface {
	SetThreadCount(count int)
	SetDelay(delay time.Duration)
	SetGracefulShutdown(gs graceful_shutdown.GracefulShutdown)
	Loop(data []T, f func(data T))
	WaitFinish()
}

type looper[T any] struct {
	threadCount      int
	sm               chan struct{}
	delay            time.Duration
	gracefulShutdown graceful_shutdown.GracefulShutdown
}

func NewLooper[T any]() Looper[T] {
	return &looper[T]{
		threadCount: 1,
	}
}

func (l *looper[T]) SetThreadCount(count int) {
	l.threadCount = count
}

func (l *looper[T]) SetDelay(d time.Duration) {
	l.delay = d
}

func (l *looper[T]) SetGracefulShutdown(gs graceful_shutdown.GracefulShutdown) {
	l.gracefulShutdown = gs
}

func (l *looper[T]) WaitFinish() {
	if l.sm == nil {
		return
	}

	for i := 0; i < l.threadCount; i++ {
		<-l.sm
	}
}

func (l *looper[T]) Loop(data []T, f func(data T)) {
	l.sm = make(chan struct{}, l.threadCount)

	for i := 0; i < l.threadCount; i++ {
		l.sm <- struct{}{}
	}

	for _, d := range data {
		<-l.sm

		go func() {
			defer func() {
				l.sm <- struct{}{}
			}()
			f(d)
		}()

		if l.gracefulShutdown != nil && l.gracefulShutdown.IsTerminated() {
			break
		}
	}

	l.WaitFinish()
}
