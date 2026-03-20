package network

type ExportPacket struct {
	Data []byte
}

func (p *ExportPacket) Id() ID {
	return IDExport
}

func (p *ExportPacket) Encode(stream *BinaryStream) {}
func (p *ExportPacket) Decode(stream *BinaryStream) error {
	p.Data = stream.RemainingSlice()
	return nil
}
