package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

// const something uint8 = 2
func (g *Generator) op_const(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	length := node.FieldsLength()
	if length == 0 {
		return nil
	}

	multiline := length > 1
	separatorFlag := SeperatorSpace
	if multiline {
		separatorFlag = SeperatorTab
	}

	// keyword
	buf.Write(TokenConst.Bytes())

	if multiline {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenParenLeft.Bytes())
	}

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(golang.ValueFlagPointer) == 0 {
			return fmt.Errorf("const node fields must be pointers")
		}

		// Checks whether the field value's node is valid or not
		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		if multiline {
			buf.Write(TokenNewLine.Bytes())
		}

		// Evaluate constant node
		if err := g.evalNode(buf, target, separatorFlag); err != nil {
			return err
		}
	}

	if multiline {
		buf.Write(TokenNewLine.Bytes())
		buf.Write(TokenParenRight.Bytes())
	}

	return nil
}
