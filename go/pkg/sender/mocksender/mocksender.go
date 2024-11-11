package mocksender

import (
	"context"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"log"
	"math/rand"
	"strconv"
	"time"
)

func ThrowDice(ctx context.Context, ch chan<- sender.Item, interval int) {
	for {
		select {
		case <-ctx.Done():
			return
		case <-time.After(time.Duration(interval) * time.Second):
			ch <- sender.Item{Key: "test.dice[A]", Value: strconv.Itoa(rand.Intn(100)), Clock: int(time.Now().Unix()), Ns: time.Now().Nanosecond()}
			ch <- sender.Item{Key: "test.dice[B]", Value: strconv.Itoa(rand.Intn(100)), Clock: int(time.Now().Unix()), Ns: time.Now().Nanosecond()}
		}
	}
}

type MockSender struct{}

func (m *MockSender) Send(items []sender.Item, errorSink chan<- sender.ItemSendError) {
	log.Println("Send items ", len(items))
	for _, item := range items {
		println(item.Key, item.Value, item.Clock, item.Ns)
	}
}

func LogErrors(ctx context.Context, ch <-chan sender.ItemSendError) {
	for {
		select {
		case <-ctx.Done():
			return
		case e := <-ch:
			log.Println("Sender error: ", e.Err.Error())
		}
	}
}
