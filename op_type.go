package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

// type Struct struct { ... }
func (g *Generator) op_type(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	var field program.NodeValue
	node.Fields(&field, 0)

	name, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}

	// keyword
	buf.Write(TokenType.Bytes())

	buf.Write(TokenSpace.Bytes())
	buf.Write(name)
	buf.Write(TokenSpace.Bytes())

	// Look for the field type
	def, ok := g.LookUpType(field.Type())
	if !ok {
		return fmt.Errorf("type with id %d is undefined", field.Type())
	}

	// Check whether the field is either a structure or a function type
	kind := golang.Kind(def.Base())
	if kind != golang.Kindstruct_ && kind != golang.Kindinterface_ && kind != golang.Kindfunc_ {
		return fmt.Errorf("the operand must be either a structure or a function but got: %d", kind)
	}

	if kind == golang.Kindfunc_ {
		buf.Write(TokenFunc.Bytes())
	}

	// Evaluates the type definition and writes to the buffer
	if err := g.evalType(buf, def); err != nil {
		return err
	}

	return nil
}
