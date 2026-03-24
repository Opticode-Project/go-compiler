package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	"github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_switch(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var body bytes.Buffer
	buf.Write(TokenSwitch.Bytes())

	for i := 0; i < node.FieldsLength(); i++ {
		var field program.NodeValue
		node.Fields(&field, i)

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		// Expression
		case field.Flags()&uint32(golang.ValueFlagSwitchExp) != 0:
			buf.Write(TokenSpace.Bytes())

			if err := g.evalNode(buf, target, 0); err != nil {
				return err
			}

		// Body (cases/default)
		case field.Flags()&uint32(golang.ValueFlagSwitchBody) != 0:
			body.Write(TokenTab.Bytes())

			if err := g.evalNode(&body, target, 0); err != nil {
				return err
			}

			body.Write(TokenNewLine.Bytes())
		}
	}

	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenBraceLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())

	// Write body
	buf.Write(body.Bytes())

	buf.Write(TokenBraceRight.Bytes())

	return nil
}
