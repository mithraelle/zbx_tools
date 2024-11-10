package sender

import (
	"context"
	"log"
	"time"
)

const BufferTimeout = time.Second * 60
const BufferSendTries = 5

type ItemBuffer struct {
	Size    int
	Timeout time.Duration
}

func NewItemBuffer() *ItemBuffer {
	return &ItemBuffer{Size: 100, Timeout: BufferTimeout}
}

func (s *ItemBuffer) Read(ctx context.Context, in <-chan Item, send Sender) {
	items := make([]Item, 0)
	timeout := time.After(s.Timeout)

	for {
		select {
		case <-ctx.Done():
			return
		case item := <-in:
			items = append(items, item)
			log.Println("Got item. Items: ", len(items))
			if len(items) >= s.Size {
				go send.Send(items, BufferSendTries)
				items = make([]Item, 0)
				timeout = time.After(s.Timeout)
			}
		case <-timeout:
			log.Println("Buffer timeout")
			if len(items) > 0 {
				go send.Send(items, BufferSendTries)
				items = make([]Item, 0)
				timeout = time.After(s.Timeout)
			}
		}
	}
}
