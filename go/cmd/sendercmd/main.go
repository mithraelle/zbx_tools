package main

import (
	"context"
	"flag"
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/cmdsender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/mocksender"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var dummyRun = flag.Bool("dummy", false, "help message for flag int")
var zbxSenderBin = flag.String("sender-bin", "", "help message for flag int")
var zbxSenderConfig = flag.String("config", "", "help message for flag int")

func main() {
	fmt.Println("Zabbix sender command test")
	flag.Parse()

	fmt.Println("Command: ", *zbxSenderBin)
	fmt.Println("Config: ", *zbxSenderConfig)
	fmt.Println("Dummy Run: ", *dummyRun)

	senderChan := make(chan sender.Item)
	errorChan := make(chan sender.ItemSendError)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go mocksender.ThrowDice(ctx, senderChan, 3)

	cmdSender := cmdsender.NewCMDSender(*zbxSenderBin, *zbxSenderConfig)
	cmdSender.DummyRun = *dummyRun
	iCollector := sender.NewItemCollector()
	iCollector.Timeout = 10 * time.Second

	go iCollector.Read(ctx, senderChan, cmdSender, errorChan)
	go mocksender.LogErrors(ctx, errorChan)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	<-sigs
}
