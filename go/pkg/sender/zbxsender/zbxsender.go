package zbxsender

import (
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"log"
	"net"
	"time"
)

const ConnectionTimeout = time.Second * 10

type ZBXSender struct {
	conf Config
}

func NewZBXSender(conf Config) *ZBXSender {
	return &ZBXSender{conf: conf}
}

func (z *ZBXSender) connect() (net.Conn, error) {
	return net.DialTimeout("tcp", z.conf.ServerAddr, ConnectionTimeout)
}

func (z *ZBXSender) Send(items []sender.Item, try int) error {
	con, err := z.connect()
	if err != nil {
		panic(err)
	}
	defer con.Close()

	data := newPacket(z.conf.Hostname, items).GetPayload()
	log.Println("Send data: ", string(data))
	_, err = con.Write(data)
	if err != nil {
		err = fmt.Errorf("Error while sending the data: %s", err.Error())
		return err
	}

	res := make([]byte, 1024)
	_, err = con.Read(res)
	if err != nil {
		err = fmt.Errorf("Error while receiving the data: %s", err.Error())
		return err
	}

	return nil
}
