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

func (z *CMDSender) Send(items []sender.Item, try int) error {
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

	for _, v := range items {
		host := v.Host
		if host == "" {
			host = "-"
		}

		fmt.Fprintf(senderIn, "%v %v %v %v\n", host, v.Key, v.Clock, v.Value)
	}

	if z.DummyRun {
		return err
	}
	senderIn.Close()

	if _, err := senderCmd.CombinedOutput(); err != nil {
		if try > 0 {
			z.Send(items, try-1)
		}
	}

	return err
}
