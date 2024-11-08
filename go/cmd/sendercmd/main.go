package main

import (
	"context"
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/cmd"
	"strconv"
	"time"
)

func main() {
	var key, value string

	fmt.Println("Zabbix sender command test")

	zbxSender := cmd.NewZbxCmdSender("", "")
	zbxSender.DummyRun = true
	zbxSender.Interval = 10 * time.Second

	senderChan := make(chan *cmd.CommandValue)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go zbxSender.Run(ctx, senderChan)

	fmt.Println("Enter key value pairs")
	for {
		_, err := fmt.Scanf("%v %v", &key, &value)
		if err != nil {
			fmt.Println("Error reading input: ", err.Error())
		} else {
			fmt.Printf("Key: %v, Value: %v\n", key, value)
			senderChan <- &cmd.CommandValue{Key: key, Value: value, TS: strconv.Itoa(int(time.Now().Unix()))}
		}
	}
}
