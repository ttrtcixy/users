package app

import (
	"log"
	"os"
	"os/signal"
)

var globalCloser = NewCloser()

type Close interface {
	Close() error
}

type Closer struct {
	closers []Close
}

func NewCloser() *Closer {
	closer := &Closer{}
	closer.start()
	return closer
}

func (c *Closer) start() {
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt)
		<-ch
		signal.Stop(ch)
		c.CloseAll()
	}()
}

func (c *Closer) Add(closer Close) {
	c.closers = append(c.closers, closer)
}

func (c *Closer) CloseAll() {
	for _, v := range c.closers {
		err := v.Close()
		if err != nil {
			log.Println(err)
		}
	}
}
