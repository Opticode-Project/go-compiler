package network

import "fmt"

type VerificationResponsePacket struct {
	Accepted bool
	Message  string
}

func (p *VerificationResponsePacket) Id() ID {
	return IDRVerification
}

func (p *VerificationResponsePacket) Encode(stream *BinaryStream) {
	stream.WriteBoolean(p.Accepted)
	stream.WriteString(p.Message, false)
}

func (p *VerificationResponsePacket) Decode(stream *BinaryStream) error {
	return fmt.Errorf("packet not handled by server")
}
