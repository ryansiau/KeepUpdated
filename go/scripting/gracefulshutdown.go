package scripting

import (
	"os"
	"os/signal"
	"syscall"
)

type GracefulShutdown interface {
	IsTerminated() bool
	Terminate(signal os.Signal)
}

type gracefulShutdown struct {
	isTerminated bool
	ch           chan os.Signal
}

func NewGracefulShutdown() GracefulShutdown {
	g := gracefulShutdown{
		isTerminated: false,
		ch:           make(chan os.Signal, 1),
	}

	signal.Notify(g.ch, os.Interrupt, os.Kill, syscall.SIGTERM)

	go func() {
		_ = <-g.ch
		g.isTerminated = true
	}()

	return &g
}

func (g *gracefulShutdown) IsTerminated() bool {
	return g.isTerminated
}

func (g *gracefulShutdown) Terminate(signal os.Signal) {
	g.ch <- signal
}
