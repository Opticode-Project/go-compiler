package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

// go worker()
func (g *Generator) op_goRoutine(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)

	if value == nil {
		return fmt.Errorf("go routine value cannot be nil")
	}

	// Checks whether the value is valid or not
	nodeValue := g.GetNode(value.Value())
	if nodeValue == nil {
		return fmt.Errorf("attempt to access undefined node: %d", value.Value())
	}

	op := golang.Opcode(nodeValue.Opcode())
	if op != golang.OpcodeCall {
		return fmt.Errorf("the operand must be a call but got: %d", op)
	}

	buf.Write(TokenGo.Bytes())
	buf.Write(TokenSpace.Bytes())

	// Evaluates the call node and writes to the buffer
	return g.evalNode(buf, nodeValue, flags)
}
