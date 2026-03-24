package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	"github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_default(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var body bytes.Buffer
	buf.Write(TokenDefault.Bytes())

	for i := 0; i < node.FieldsLength(); i++ {
		var field program.NodeValue
		node.Fields(&field, i)

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		if field.Flags()&uint32(golang.ValueFlagCaseBody) == 0 {
			return fmt.Errorf("expected indexed field to be a case body")
		}

		body.Write(TokenTab.Bytes())
		if err := g.evalNode(&body, target, 1); err != nil {
			return err
		}
		body.Write(TokenNewLine.Bytes())
	}

	buf.Write(TokenColon.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(body.Bytes())
	return nil
}
