package network

type VerificationPacket struct {
	Protocol uint16
	Language string
}

func (p *VerificationPacket) Id() ID {
	return IDVerification
}

func (p *VerificationPacket) Encode(stream *BinaryStream) {}
func (p *VerificationPacket) Decode(stream *BinaryStream) error {
	// Protocol version
	protocol, err := stream.ReadUInt16(false)
	if err != nil {
		return err
	}

	p.Protocol = protocol

	// Language used for code generation
	language, err := stream.ReadString(false)
	if err != nil {
		return err
	}

	p.Language = language
	return nil
}
