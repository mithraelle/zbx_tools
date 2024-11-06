package cmd

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

const ZabbixSenderInterval = 60 * time.Second
const ZabbixSenderTimeout = 15 * time.Second
const ZabbixSenderMaxValues = 500
const ZabbixSenderMaxTries = 3

type SenderCommand struct {
	Bin       string
	Config    string
	Interval  time.Duration
	Timeout   time.Duration
	MaxValues int
	DummyRun  bool

	inChan chan *CommandValue
}

type CommandValue struct {
	Host  string
	Key   string
	Value string
	TS    string
}

func NewZbxCmdSender(bin string, config string) *SenderCommand {
	if bin == "" {
		bin = "zabbix_sender"
	}

	if config == "" {
		config = "/etc/zabbix/zabbix_agentd.conf"
	}

	return &SenderCommand{Bin: bin, Config: config, Timeout: ZabbixSenderTimeout, Interval: ZabbixSenderInterval, MaxValues: ZabbixSenderMaxValues}
}

func (z *SenderCommand) Run(ctx context.Context, ch chan *CommandValue) {
	z.inChan = ch
	values := make([]*CommandValue, 0)
	timeout := time.After(z.Interval)

	for {
		select {
		case <-ctx.Done():
			return
		case v := <-ch:
			values = append(values, v)
			if len(values) >= z.MaxValues {
				go z.SendValues(values, 0)
				values = make([]*CommandValue, 0)
			}
		case <-timeout:
			if len(values) > 0 {
				go z.SendValues(values, 0)
				values = make([]*CommandValue, 0)
			}
			timeout = time.After(z.Timeout)
		}
	}
}

func (z *SenderCommand) getCommand() *exec.Cmd {
	return exec.Command(z.Bin,
		"-c",
		z.Config,
		"-t",
		fmt.Sprintf("%.0f", z.Timeout.Seconds()),
		"-T",
		"-i",
		"-")
}

func (z *SenderCommand) SendValues(values []*CommandValue, try int) {
	var senderIn io.WriteCloser
	var err error

	senderCmd := z.getCommand()

	if z.DummyRun {
		log.Println("RUN: ", senderCmd.String())
		senderIn = os.Stdout
	} else {
		senderIn, err = senderCmd.StdinPipe()
		if err != nil {
			log.Fatalln(err.Error())
		}
		defer senderIn.Close()
	}

	for _, v := range values {
		host := v.Host
		if host == "" {
			host = "-"
		}

		fmt.Fprintf(senderIn, "%v %v %v %v\n", host, v.Key, v.TS, v.Value)
	}

	if z.DummyRun {
		return
	}
	senderIn.Close()

	if _, err := senderCmd.CombinedOutput(); err != nil {
		if try < ZabbixSenderMaxTries {
			z.SendValues(values, try+1)
		} else {
			for _, v := range values {
				z.inChan <- v
			}
		}
	}
}
