package mqtt

import (
	"testing"
)

func TestPacketPublish(t *testing.T) {
	input := []byte{0x30, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F, 0xF0}

	pkt := NewPacketPublish()
	if err := pkt.Parse(input); err != nil {
		t.Errorf(err.Error())
		return
	}

	output := pkt.Bytes()
	if len(input) == len(output) {
		for i := 0; i < len(input); i++ {
			if input[i] != output[i] {
				t.Errorf("Mismatch %02x vs %02x\n", input[i], output[i])
			}
		}
	} else {
		t.Errorf("Mismatch length %x vs %x\n", len(input), len(output))
	}

	invalids := [][]byte{{0x32},
		{0x32, 0x07},
		{0x32, 0x07, 0x00},
		{0x32, 0x07, 0x00, 0x03},
		{0x32, 0x07, 0x00, 0x03, 0x61},
		{0x32, 0x07, 0x00, 0x03, 0x61, 0x2F},
		{0x32, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62},
		{0x32, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F},
		{0x32, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F, 0xF0, 0x90},
		{0xF2, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F, 0xF0},
		{0x36, 0x07, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F, 0xF0},
		{0x32, 0x02, 0x00, 0x03, 0x61, 0x2F, 0x62, 0x0F, 0xF0},
		{0x32, 0x07, 0x00, 0x04, 0x61, 0x2F, 0x62, 0x0F, 0xF0},
		{0x32, 0x07, 0x00, 0x05, 0x61, 0x2F, 0x62, 0x0F, 0xF0},
		{0x32, 0x07, 0x00, 0x06, 0x61, 0x2F, 0x62, 0x0F, 0xF0}}
	for i := 0; i < len(invalids); i++ {
		if err := pkt.Parse(invalids[i]); err != nil {
			t.Logf(err.Error())
		} else {
			t.Logf("%v", invalids[i])
		}
	}
}
