// To lazy to port the TS BinaryStream to Go so I asked AI. I do hope this works
package network

import (
	"errors"
	"math"
	"math/bits"
)

const MaxStringSize = 1 << 16

type BinaryStream struct {
	buf []byte
	r   int
}

func NewStream(cap uint32) *BinaryStream {
	if cap <= 0 {
		cap = 64
	}

	return &BinaryStream{
		buf: make([]byte, 0, cap),
	}
}

func FromData(data []byte) *BinaryStream {
	return &BinaryStream{buf: data}
}

func (s *BinaryStream) Size() int {
	return len(s.buf)
}

func (s *BinaryStream) Capacity() int {
	return cap(s.buf)
}

func (s *BinaryStream) Remaining() int {
	return len(s.buf) - s.r
}

func (s *BinaryStream) Reset() {
	s.r = 0
}

func (s *BinaryStream) ensure(n int) error {
	if s.r+n > len(s.buf) {
		return errors.New("read past end")
	}
	return nil
}

func (s *BinaryStream) PeekByte() (byte, error) {
	if err := s.ensure(1); err != nil {
		return 0, err
	}
	return s.buf[s.r], nil
}

func (s *BinaryStream) Skip(n int) error {
	if err := s.ensure(n); err != nil {
		return err
	}
	s.r += n
	return nil
}

func (s *BinaryStream) RemainingSlice() []byte {
	return s.buf[s.r:]
}

func (s *BinaryStream) ReadByte() (byte, error) {
	if err := s.ensure(1); err != nil {
		return 0, err
	}

	b := s.buf[s.r]
	s.r++
	return b, nil
}

func (s *BinaryStream) ReadSignedByte() (int8, error) {
	b, err := s.ReadByte()
	return int8(b), err
}

func (s *BinaryStream) ReadBytes(n int) ([]byte, error) {
	if err := s.ensure(n); err != nil {
		return nil, err
	}

	b := s.buf[s.r : s.r+n]
	s.r += n
	return b, nil
}

func (s *BinaryStream) WriteByte(b byte) error {
	s.buf = append(s.buf, b)
	return nil
}

func (s *BinaryStream) WriteSignedByte(v int8) {
	s.WriteByte(byte(v))
}

func (s *BinaryStream) WriteBytes(b []byte) {
	s.buf = append(s.buf, b...)
}

func (s *BinaryStream) ReadUInt16(le bool) (uint16, error) {
	if err := s.ensure(2); err != nil {
		return 0, err
	}

	b := s.buf[s.r : s.r+2]
	s.r += 2

	v := uint16(b[0])<<8 | uint16(b[1])

	if le {
		v = bits.ReverseBytes16(v)
	}

	return v, nil
}

func (s *BinaryStream) WriteUInt16(v uint16, le bool) {
	if le {
		v = bits.ReverseBytes16(v)
	}

	s.buf = append(s.buf,
		byte(v>>8),
		byte(v),
	)
}

func (s *BinaryStream) ReadUInt32(le bool) (uint32, error) {
	if err := s.ensure(4); err != nil {
		return 0, err
	}

	b := s.buf[s.r : s.r+4]
	s.r += 4

	v := uint32(b[0])<<24 |
		uint32(b[1])<<16 |
		uint32(b[2])<<8 |
		uint32(b[3])

	if le {
		v = bits.ReverseBytes32(v)
	}

	return v, nil
}

func (s *BinaryStream) WriteUInt32(v uint32, le bool) {
	if le {
		v = bits.ReverseBytes32(v)
	}

	s.buf = append(s.buf,
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (s *BinaryStream) ReadUInt64(le bool) (uint64, error) {
	if err := s.ensure(8); err != nil {
		return 0, err
	}

	b := s.buf[s.r : s.r+8]
	s.r += 8

	v := uint64(b[0])<<56 |
		uint64(b[1])<<48 |
		uint64(b[2])<<40 |
		uint64(b[3])<<32 |
		uint64(b[4])<<24 |
		uint64(b[5])<<16 |
		uint64(b[6])<<8 |
		uint64(b[7])

	if le {
		v = bits.ReverseBytes64(v)
	}

	return v, nil
}

func (s *BinaryStream) WriteUInt64(v uint64, le bool) {
	if le {
		v = bits.ReverseBytes64(v)
	}

	s.buf = append(s.buf,
		byte(v>>56),
		byte(v>>48),
		byte(v>>40),
		byte(v>>32),
		byte(v>>24),
		byte(v>>16),
		byte(v>>8),
		byte(v),
	)
}

func (s *BinaryStream) ReadFloat32(le bool) (float32, error) {
	v, err := s.ReadUInt32(le)
	return math.Float32frombits(v), err
}

func (s *BinaryStream) WriteFloat32(v float32, le bool) {
	s.WriteUInt32(math.Float32bits(v), le)
}

func (s *BinaryStream) ReadFloat64(le bool) (float64, error) {
	v, err := s.ReadUInt64(le)
	return math.Float64frombits(v), err
}

func (s *BinaryStream) WriteFloat64(v float64, le bool) {
	s.WriteUInt64(math.Float64bits(v), le)
}

func (s *BinaryStream) ReadUnsignedVarInt() (uint32, error) {
	var value uint32
	var shift uint

	for i := 0; i < 5; i++ {
		b, err := s.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= uint32(b&0x7F) << shift

		if b&0x80 == 0 {
			return value, nil
		}

		shift += 7
	}

	return 0, errors.New("VarInt too big")
}

func (s *BinaryStream) WriteUnsignedVarInt(v uint32) {
	for v > 0x7F {
		s.WriteByte(byte(v&0x7F | 0x80))
		v >>= 7
	}
	s.WriteByte(byte(v))
}

func (s *BinaryStream) ReadSignedVarInt() (int32, error) {
	u, err := s.ReadUnsignedVarInt()
	if err != nil {
		return 0, err
	}

	return int32((u >> 1) ^ uint32(-(int32(u & 1)))), nil
}

func (s *BinaryStream) WriteSignedVarInt(v int32) {
	zigzag := uint32(v<<1) ^ uint32(v>>31)
	s.WriteUnsignedVarInt(zigzag)
}

func (s *BinaryStream) ReadUnsignedVarLong() (uint64, error) {
	var value uint64

	for i := 0; i < 10; i++ {
		b, err := s.ReadByte()
		if err != nil {
			return 0, err
		}

		value |= uint64(b&0x7F) << (7 * i)

		if b&0x80 == 0 {
			return value, nil
		}
	}

	return 0, errors.New("VarLong too big")
}

func (s *BinaryStream) WriteUnsignedVarLong(v uint64) {
	for v > 0x7F {
		s.WriteByte(byte(v&0x7F | 0x80))
		v >>= 7
	}
	s.WriteByte(byte(v))
}

func (s *BinaryStream) ReadBoolean() (bool, error) {
	b, err := s.ReadByte()
	return b != 0, err
}

func (s *BinaryStream) WriteBoolean(v bool) {
	if v {
		s.WriteByte(1)
	} else {
		s.WriteByte(0)
	}
}

func (s *BinaryStream) ReadString(le bool) (string, error) {
	length, err := s.ReadUInt32(le)
	if err != nil {
		return "", err
	}

	if length > MaxStringSize {
		return "", errors.New("string too large")
	}

	b, err := s.ReadBytes(int(length))
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (s *BinaryStream) WriteString(v string, le bool) {
	s.WriteUInt32(uint32(len(v)), le)
	s.WriteBytes([]byte(v))
}

func (s *BinaryStream) ReadStream(n int) (*BinaryStream, error) {
	b, err := s.ReadBytes(n)
	if err != nil {
		return nil, err
	}

	return FromData(b), nil
}

func (s *BinaryStream) ToBytes(includeSize bool) []byte {
	if !includeSize {
		return s.buf
	}

	size := len(s.buf)
	out := make([]byte, 4+size)

	out[0] = byte(size >> 24)
	out[1] = byte(size >> 16)
	out[2] = byte(size >> 8)
	out[3] = byte(size)

	copy(out[4:], s.buf)

	return out
}
