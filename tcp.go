package compiler

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
	"slices"
	"time"

	"github.com/Opticode-Project/go-compiler/go/golang"
	"github.com/Opticode-Project/go-compiler/network"
	"github.com/Opticode-Project/go-compiler/program"
	flatbuffers "github.com/google/flatbuffers/go"
)

const OpticodeProtocol uint16 = 1

type TCPServer struct {
	generator *Generator
	listener  net.Listener
	builder   *flatbuffers.Builder
}

func (s *TCPServer) Start(addr string) error {
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	s.listener = ln
	fmt.Println("Server listening on address", ln.Addr())

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

		case *network.SetNodePacket:
			node := program.GetRootAsNode(p.Data, 0)
			if node == nil {
				break
			}

			s.generator.SetNode(node)

		case *network.ExportPacket:
			app := golang.GetRootAsApp(p.Data, 0)

			keys := make([]uint64, 0, len(s.generator.nodes))
			for id := range s.generator.nodes {
				keys = append(keys, id)
			}

			slices.Sort(keys)

			var path = []uint64{0}
			for _, id := range keys {
				node := s.generator.nodes[id]
				log.Println(id, golang.Opcode(node.Opcode()), golang.NodeFlag(node.Flags()), int64(node.Next()))

				if id == path[len(path)-1] {
					path = append(path, node.Next())
				}
			}

			log.Printf("Path: %v", path)
			now := time.Now()

			s.generator.Compile(app, path)
			gf, err := s.generator.Export("main", path)
			if err != nil {
				panic(err)
			}

			log.Printf("Time elapse: %dms", time.Since(now).Milliseconds())

			for _, g := range gf {
				log.Println(string(*g.Content))
			}

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
