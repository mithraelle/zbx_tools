package sender

import (
	"context"
	"log"
	"time"
)

const BufferTimeout = time.Second * 60

type ItemCollector struct {
	Size    int
	Timeout time.Duration
}

func NewItemCollector() *ItemCollector {
	return &ItemCollector{Size: 100, Timeout: BufferTimeout}
}

func (s *ItemCollector) Read(ctx context.Context, fanIn <-chan Item, send Sender, errorSink <-chan ItemSendError) {
	items := make([]Item, 0)
	timeout := time.After(s.Timeout)

	for {
		select {
		case <-ctx.Done():
			return
		case item := <-fanIn:
			items = append(items, item)
			log.Println("Got item. Items: ", len(items))
			if len(items) >= s.Size {
				go send.Send(items, errorSink)
				items = make([]Item, 0)
				timeout = time.After(s.Timeout)
			}
		case <-timeout:
			log.Println("Buffer timeout")
			if len(items) > 0 {
				go send.Send(items, errorSink)
				items = make([]Item, 0)
				timeout = time.After(s.Timeout)
			}
		}
	}
}
