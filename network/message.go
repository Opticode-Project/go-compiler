package network

import (
	"fmt"

	flatbuffers "github.com/google/flatbuffers/go"
)

type Message interface {
	Id() ID
	Encode(stream *BinaryStream)
	Decode(stream *BinaryStream) error
}

var packetRegistry = map[ID]func() Message{
	IDVerification:  func() Message { return &VerificationPacket{} },
	IDRVerification: func() Message { return &VerificationResponsePacket{} },
}

func getPayloadBytes(packet *Packet) []byte {
	length := packet.PayloadLength()
	buf := make([]byte, length)

	for i := range length {
		buf[i] = byte(packet.Payload(i))
	}

	return buf
}

func DecodePacket(packet *Packet) (Message, error) {
	constructor, ok := packetRegistry[packet.Id()]
	if !ok {
		return nil, fmt.Errorf("unknown packet %d", packet.Id())
	}

	instance := constructor()
	payload := FromData(getPayloadBytes(packet))

	if err := instance.Decode(payload); err != nil {
		return nil, err
	}

	return instance, nil
}

func EncodePacket(packet Message, builder *flatbuffers.Builder) []byte {
	builder.Reset()

	packetPayload := NewStream(0)
	packet.Encode(packetPayload)

	payload := packetPayload.ToBytes(false)
	payloadOffset := builder.CreateByteVector(payload)

	PacketStart(builder)
	PacketAddId(builder, packet.Id())
	PacketAddPayload(builder, payloadOffset)

	packetOffset := PacketEnd(builder)
	builder.Finish(packetOffset)

	return builder.FinishedBytes()
}
