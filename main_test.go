package compiler

import (
	"testing"
)

func TestCompile(t *testing.T) {
	generator := NewGenerator()

	generator.Listen("127.0.0.1:27430")
}
