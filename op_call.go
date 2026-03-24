package compiler

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

// fmt.Println("hello, world!")
func (g *Generator) op_call(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()

	name, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}
	buf.Write(name)

	buf.Write(TokenParenLeft.Bytes())

	// Parameters
	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if i > 0 {
			buf.Write(TokenComma.Bytes())
			buf.Write(TokenSpace.Bytes())
		}

		if err := g.evalValue(buf, &field, false); err != nil {
			return err
		}
	}

	buf.Write(TokenParenRight.Bytes())
	return nil
}
