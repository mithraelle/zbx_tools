package zbxsender

import (
	"errors"
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"github.com/mithraelle/zbx_tools/go/pkg/sender/zbxsender/config"
	"net"
	"time"
)

const ConnectionTimeout = time.Second * 10

type ZBXSender struct {
	conf config.Config
}

func NewZBXSender(conf config.Config) *ZBXSender {
	return &ZBXSender{conf: conf}
}

func (z *ZBXSender) connect() (net.Conn, error) {
	if z.conf.TLSConnect == "cert" {
		return nil, errors.New("the cert based auth isn't supported")
	} else if z.conf.TLSConnect == "psk" {
		return nil, errors.New("the PSK auth isn't supported")
	}
	return net.DialTimeout("tcp", z.conf.ServerAddr, ConnectionTimeout)
}

func (z *ZBXSender) Send(items []sender.Item, errorSink chan<- sender.ItemSendError) {
	con, err := z.connect()
	if err != nil {
		defer con.Close()

		data := newPacket(z.conf.Hostname, items).GetPayload()
		_, err = con.Write(data)
		if err != nil {
			err = fmt.Errorf("error while sending the data: %s", err.Error())
		} else {
			res := make([]byte, 1024)
			_, err = con.Read(res)
			if err != nil {
				err = fmt.Errorf("error while receiving the data: %s", err.Error())
			}
		}
	}

	if err != nil && errorSink != nil {
		errorSink <- sender.ItemSendError{Err: err, Items: items}
	}
}
