package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

// if x > y { ... }
func (g *Generator) op_if(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var (
		condition bytes.Buffer
		thenBody  bytes.Buffer
		elseBody  bytes.Buffer
	)
	var elseIfNode *program.Node

	length := node.FieldsLength()
	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(golang.ValueFlagPointer) == 0 {
			return fmt.Errorf("if node fields can only be pointers")
		}

		// Checks whether the field value's node is valid or not
		node := g.GetNode(field.Value())
		if node == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		case field.Flags()&uint32(golang.ValueFlagIfConditon) != 0:
			// Evaluate the if statement's condition
			if err := g.evalNode(&condition, node, 0); err != nil {
				return err
			}

		case field.Flags()&uint32(golang.ValueFlagIfBody) != 0:
			thenBody.Write(TokenTab.Bytes())

			if err := g.evalNode(&thenBody, node, 0); err != nil {
				return err
			}

			thenBody.Write(TokenNewLine.Bytes())

		case field.Flags()&uint32(golang.ValueFlagIfElse) != 0:
			if golang.Opcode(node.Opcode()) == golang.OpcodeIf {
				elseIfNode = node // else if statement
			} else {
				elseBody.Write(TokenTab.Bytes())

				// Evaluate the else body statement
				if err := g.evalNode(&elseBody, node, 0); err != nil {
					return err
				}

				elseBody.Write(TokenNewLine.Bytes())
			}
		}
	}

	buf.Write(TokenIf.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(condition.Bytes())
	buf.Write(TokenSpace.Bytes())

	buf.Write(TokenBraceLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(thenBody.Bytes())
	buf.Write(TokenBraceRight.Bytes())

	if elseIfNode != nil {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenElse.Bytes())
		buf.Write(TokenSpace.Bytes())

		// Evaluates the else if node and writes to the buffer
		if err := g.evalNode(buf, elseIfNode, flags); err != nil {
			return err
		}
	} else if elseBody.Len() > 0 {
		buf.Write(TokenSpace.Bytes())
		buf.Write(TokenElse.Bytes())
		buf.Write(TokenSpace.Bytes())

		buf.Write(TokenBraceLeft.Bytes())
		buf.Write(TokenNewLine.Bytes())
		buf.Write(elseBody.Bytes())
		buf.Write(TokenBraceRight.Bytes())
	}

	return nil
}
