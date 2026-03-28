package compiler

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

// <identifier>
func (g *Generator) op_reference(buf *bytes.Buffer, node *program.UnaryNode, flags EvalFlags) error {
	value := node.Value(nil)
	if value == nil {
		return fmt.Errorf("unary operands cannot be nil")
	}

	// Evaluate unary node
	return g.evalValue(buf, value, false)
}
