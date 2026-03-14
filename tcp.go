package golang

import (
	"encoding/binary"
	"fmt"
	"io"
	"net"

	"github.com/Opticode-Project/go-compiler/network"
	flatbuffers "github.com/google/flatbuffers/go"
)

const OpticodeProtocol uint16 = 1

type TCPServer struct {
	Port     uint16
	listener net.Listener
	builder  *flatbuffers.Builder
}

func (s *TCPServer) Start() error {
	addr := fmt.Sprintf("127.0.0.1:%d", s.Port)

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = ln
	fmt.Println("Server listening on port", s.Port)

	go s.acceptLoop()
	return nil
}

func (s *TCPServer) acceptLoop() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			continue
		}

		go s.handle(conn)
	}
}

func (s *TCPServer) handle(conn net.Conn) {
	defer conn.Close()

	for {
		packet, err := s.readPacket(conn)
		if err != nil {
			return
		}

		instance, err := network.DecodePacket(packet)
		if err != nil {
			fmt.Println("Failed to decode packet", err)
			continue
		}

		switch p := instance.(type) {

		case *network.VerificationPacket:
			response := &network.VerificationResponsePacket{}

			if p.Protocol != OpticodeProtocol {
				response.Accepted = false
				response.Message = "Protocol version is not supported."

				s.writePacket(conn, network.EncodePacket(response, s.builder))
				return
			}

			response.Accepted = true
			response.Message = ""

			s.writePacket(conn, network.EncodePacket(response, s.builder))

		default:
			fmt.Println("unknown packet:", p)
		}
	}
}

func (s *TCPServer) readPacket(conn net.Conn) (*network.Packet, error) {
	lengthBuf := make([]byte, 4)

	_, err := io.ReadFull(conn, lengthBuf)
	if err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthBuf)

	data := make([]byte, length)
	_, err = io.ReadFull(conn, data)
	if err != nil {
		return nil, err
	}

	root := network.GetRootAsPacket(data, 0)
	return root, nil
}

func (s *TCPServer) writePacket(conn net.Conn, buf []byte) error {
	stream := network.NewStream(0)
	stream.WriteBytes(buf)

	payload := stream.ToBytes(true)

	_, err := conn.Write(payload)
	return err
}

func (s *TCPServer) Stop() error {
	if s.listener != nil {
		return s.listener.Close()
	}

	return nil
}
