package cmdsender

import (
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"io"
	"log"
	"os"
	"os/exec"
	"time"
)

const ZabbixSenderTimeout = 15 * time.Second
const MaxTries = 5

type CMDSender struct {
	Bin     string
	Config  string
	Timeout time.Duration

	DummyRun bool
}

func NewCMDSender(bin string, config string) *CMDSender {
	if bin == "" {
		bin = "zabbix_sender"
	}

	if config == "" {
		config = "/etc/zabbix/zabbix_agentd.conf"
	}

	return &CMDSender{Bin: bin, Config: config, Timeout: ZabbixSenderTimeout}
}

func (z *CMDSender) getCommand() *exec.Cmd {
	return exec.Command(z.Bin,
		"-c",
		z.Config,
		"-t",
		fmt.Sprintf("%.0f", z.Timeout.Seconds()),
		"-T",
		"-i",
		"-")
}

func (z *CMDSender) Send(items []sender.Item, errorSink chan<- sender.ItemSendError) {
	var err error
	for i := 0; i < MaxTries; i++ {
		err = z.runCommand(items)
		if err == nil {
			return
		}
	}

	if errorSink != nil {
		errorSink <- sender.ItemSendError{Items: items, Err: err}
	}
}

func (z *CMDSender) runCommand(items []sender.Item) error {
	var senderIn io.WriteCloser
	var err error

	senderCmd := z.getCommand()

	if z.DummyRun {
		log.Println("RUN: ", senderCmd.String())
		senderIn = os.Stdout
	} else if senderIn, err = senderCmd.StdinPipe(); err != nil {
		return err
	}

	for _, v := range items {
		host := v.Host
		if host == "" {
			host = "-"
		}

		fmt.Fprintf(senderIn, "%v %v %v %v\n", host, v.Key, v.Clock, v.Value)
	}

	if !z.DummyRun {
		senderIn.Close()
		_, err = senderCmd.CombinedOutput()
	}

	return err
}
