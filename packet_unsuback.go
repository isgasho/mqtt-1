package mqtt

import (
	"bytes"
	"errors"
	"fmt"
)

////////////////////Interface//////////////////////////////

type UnsubackPacket interface {
	Packet

	//Variable Header
	GetPacketId() uint16
	SetPacketId(id uint16)
}

////////////////////Implementation////////////////////////

type packet_unsuback struct {
	packet

	packetId uint16
}

func NewPacketUnsuback() *packet_unsuback {
	this := packet_unsuback{}

	this.IBytizer = &this
	this.IParser = &this

	this.packetType = PACKET_UNSUBACK
	this.packetFlag = 0

	return &this
}

func (this *packet_unsuback) IBytize() []byte {
	var buffer bytes.Buffer

	buffer.WriteByte((byte(this.packetType) << 4) | (this.packetFlag & 0x0F))
	buffer.WriteByte(2)
	buffer.WriteByte(byte(this.packetId >> 8))
	buffer.WriteByte(byte(this.packetId & 0xFF))

	return buffer.Bytes()
}

func (this *packet_unsuback) IParse(buffer []byte) error {
	if buffer == nil || len(buffer) != 4 {
		return errors.New("Invalid Control Packet Size")
	}

	if packetType := PacketType((buffer[0] >> 4) & 0x0F); packetType != this.packetType {
		return fmt.Errorf("Invalid Control Packet Type %d\n", packetType)
	}
	if packetFlag := buffer[0] & 0x0F; packetFlag != this.packetFlag {
		return fmt.Errorf("Invalid Control Packet Flags %d\n", packetFlag)
	}
	if buffer[1] != 2 {
		return fmt.Errorf("Invalid Control Packet Remaining Length %d\n", buffer[1])
	}

	this.packetId = ((uint16(buffer[2])) << 8) | uint16(buffer[3])

	return nil
}

//Variable Header
func (this *packet_unsuback) GetPacketId() uint16 {
	return this.packetId
}
func (this *packet_unsuback) SetPacketId(id uint16) {
	this.packetId = id
}
