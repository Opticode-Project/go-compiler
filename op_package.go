package compiler

import (
	"bytes"
	"fmt"

	program "github.com/Opticode-Project/go-compiler/program"
)

// package main
func (g *Generator) op_package(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	id, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}

	buf.Write(TokenPackage.Bytes()) // keyword
	buf.Write(TokenSpace.Bytes())
	buf.Write(id)
	return nil
}
