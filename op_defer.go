package golang

import (
	"bytes"
	"fmt"

	schema "github.com/Opticode-Project/go-compiler/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

// defer fmt.Println("Hello, world!")
func (g *Generator) op_defer(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)

	if value == nil {
		return fmt.Errorf("defer value cannot be nil")
	}

	// Checks whether the node value is valid or not
	nodeValue := g.GetNode(value.Value())
	if nodeValue == nil {
		return fmt.Errorf("attempt to access undefined node: %d", value.Value())
	}

	op := schema.Opcode(nodeValue.Opcode())
	if op != schema.OpcodeCall {
		return fmt.Errorf("the operand must be a call but got: %d", op)
	}

	buf.Write(TokenDefer.Bytes())
	buf.Write(TokenSpace.Bytes())

	// Evaluates the call node and writes to the buffer
	return g.evalNode(buf, nodeValue, flags)
}
