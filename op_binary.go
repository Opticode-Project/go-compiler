package golang

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

// x + b, ... etc
func (g *Generator) op_binary(buf *bytes.Buffer, node *program.BinaryNode, t TokenKind, flags EvalFlags) error {
	// Get the left and right values
	left := node.Left(nil)
	right := node.Right(nil)

	if left == nil || right == nil {
		return fmt.Errorf("binary value operands cannot be nil")
	}

	// Evaluate left node
	if err := g.evalValue(buf, left, false); err != nil {
		return err
	}

	buf.Write(TokenSpace.Bytes())
	buf.Write(t.Bytes()) // Write the token to the buffer
	buf.Write(TokenSpace.Bytes())

	// Evaluate right node
	return g.evalValue(buf, right, false)
}
