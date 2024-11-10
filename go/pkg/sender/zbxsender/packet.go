package zbxsender

import (
	"encoding/binary"
	"encoding/json"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"time"
)

var dataHeader = []byte("ZBXD\x01")

type Packet struct {
	Request string        `json:"request"`
	Data    []sender.Item `json:"data"`
	Clock   int           `json:"clock"`
}

func newPacket(hostname string, items []sender.Item) *Packet {
	for i := range items {
		if items[i].Host == "" || items[i].Host == "-" {
			items[i].Host = hostname
		}
	}
	return &Packet{Request: "sender data", Data: items, Clock: int(time.Now().Unix())}
}

func (p *Packet) Marshal() []byte {
	s, _ := json.Marshal(p)
	return s
}

func (p *Packet) GetPayload() []byte {
	payload := make([]byte, 0)
	payload = append(payload, dataHeader...)

	pData := p.Marshal()

	pSize := len(pData)
	dataLen := make([]byte, 8)
	binary.LittleEndian.PutUint32(dataLen, uint32(pSize))
	payload = append(payload, dataLen...)

	return append(payload, pData...)
}
