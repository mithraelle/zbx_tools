package zbxsender

import (
	"bytes"
	"fmt"
	"github.com/mithraelle/zbx_tools/go/pkg/sender"
	"testing"
)

func TestPacket_GetPayload(t *testing.T) {
	items := []sender.Item{
		sender.Item{Key: "testA", Value: "1", Clock: 2},
		sender.Item{Key: "testB", Value: "3", Clock: 4},
	}
	//JSON length is 184
	dataLen := make([]byte, 8)
	dataLen[0] = byte('\xAB')

	p := newPacket("host", items)

	expectedPayload := []byte("ZBXD\x01")
	expectedPayload = append(expectedPayload, dataLen...)
	expectedPayload = append(expectedPayload, p.Marshal()...)

	pData := p.GetPayload()
	if !bytes.Equal(pData, expectedPayload) {
		log := fmt.Sprintf("Data length: %v/%v\n", len(expectedPayload), len(pData))
		if len(pData) == len(expectedPayload) {
			for i, expectedV := range expectedPayload {
				v := pData[i]
				if v != expectedV {
					log += fmt.Sprintf("Position %v: %v/%v \n", i, expectedV, v)
				}
			}
		}
		t.Fatalf("Expected\t%v\nGot\t%v\n%v", string(expectedPayload), string(pData), log)
	}
}
