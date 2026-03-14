package golang

import (
	"testing"

	flatbuffers "github.com/google/flatbuffers/go"
)

func TestCompile(t *testing.T) {
	server := &TCPServer{
		Port:    27430,
		builder: flatbuffers.NewBuilder(256),
	}

	err := server.Start()
	if err != nil {
		panic(err)
	}

	select {} // keep running
}
