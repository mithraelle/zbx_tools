package sender

import (
	"context"
	"log"
	"time"
)

const BufferTimeout = time.Second * 60

type CollectorOption interface {
	set(*Collector)
}

type flushLimitOption int

func WithFlushLimit(limit int) CollectorOption {
	return flushLimitOption(limit)
}

func (f flushLimitOption) set(c *Collector) {
	c.FlushLimit = int(f)
}

type flushTimeoutOption time.Duration

func WithFlushTimeout(timeout time.Duration) CollectorOption {
	return flushTimeoutOption(timeout)
}

func (f flushTimeoutOption) set(c *Collector) {
	c.FlushTimeout = time.Duration(f)
}

type errorSinkOption chan<- ItemSendError

func WithErrorSink(errorSink chan<- ItemSendError) CollectorOption {
	return errorSinkOption(errorSink)
}

func (errorSink errorSinkOption) set(c *Collector) {
	c.ErrorSink = errorSink
}

type Collector struct {
	FlushLimit   int
	FlushTimeout time.Duration
	ErrorSink    chan<- ItemSendError
	Sender       Sender
}

func NewCollector(sender Sender, opts ...CollectorOption) *Collector {
	collector := &Collector{FlushLimit: 100, FlushTimeout: BufferTimeout, ErrorSink: nil, Sender: sender}

	for _, o := range opts {
		o.set(collector)
	}

	return collector
}

func (s *Collector) Read(ctx context.Context, fanIn <-chan Item) {
	items := make([]Item, 0)
	timeout := time.After(s.FlushTimeout)

	for {
		select {
		case <-ctx.Done():
			return
		case item := <-fanIn:
			items = append(items, item)
			log.Println("Got item. Items: ", len(items))
			if len(items) >= s.FlushLimit {
				go s.Sender.Send(items, s.ErrorSink)
				items = make([]Item, 0)
				timeout = time.After(s.FlushTimeout)
			}
		case <-timeout:
			log.Println("Buffer timeout")
			if len(items) > 0 {
				go s.Sender.Send(items, s.ErrorSink)
				items = make([]Item, 0)
				timeout = time.After(s.FlushTimeout)
			}
		}
	}
}
