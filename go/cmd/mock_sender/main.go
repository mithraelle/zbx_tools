package main

import (
	"context"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/mocksender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/zbxsender"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	ch := make(chan sender.Item)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	conf := zbxsender.NewConfig("zabbix_sender.conf")

	zbxSender := zbxsender.NewZBXSender(*conf)
	items := []sender.Item{
		sender.Item{Key: "trap", Value: "1", Clock: int(time.Now().Unix())},
	}

	zbxSender.Send(items, 0)
	return

	go mocksender.ThrowDice(ctx, ch, 15)
	iBuffer := sender.NewItemBuffer()
	go iBuffer.Read(ctx, ch, zbxSender)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
