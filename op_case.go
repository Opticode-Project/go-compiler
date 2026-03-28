package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	"github.com/Opticode-Project/go-compiler/program"
)

func (g *Generator) op_case(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var (
		conditions bytes.Buffer
		body       bytes.Buffer
		hasExpr    = false
	)

	for i := range node.FieldsLength() {
		var field program.NodeValue
		node.Fields(&field, i)

		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		case field.Flags()&uint32(golang.ValueFlagCaseExp) != 0:
			if conditions.Len() > 0 {
				conditions.Write(TokenComma.Bytes())
				conditions.Write(TokenSpace.Bytes())
			}

			if err := g.evalValue(&conditions, &field, false); err != nil {
				return err
			}

			hasExpr = true

		case field.Flags()&uint32(golang.ValueFlagCaseBody) != 0:
			body.Write(TokenTab.Bytes())

			if err := g.evalNode(&body, target, 1); err != nil {
				return err
			}

			body.Write(TokenNewLine.Bytes())
		}
	}

	if !hasExpr {
		return fmt.Errorf("Missing case expression")
	}

	buf.Write(TokenCase.Bytes())
	buf.Write(TokenSpace.Bytes())

	buf.Write(conditions.Bytes())
	buf.Write(TokenColon.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(body.Bytes())
	return nil
}
