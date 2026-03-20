package network

type SetNodePacket struct {
	Data []byte
}

func (p *SetNodePacket) Id() ID {
	return IDSetNode
}

func (p *SetNodePacket) Encode(stream *BinaryStream) {}
func (p *SetNodePacket) Decode(stream *BinaryStream) error {
	p.Data = stream.RemainingSlice()
	return nil
}
