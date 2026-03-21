package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	program "github.com/Opticode-Project/go-compiler/program"
)

const (
	bodyGrowthModifer    = 32
	paramsGrowthModifer  = 16
	resultsGrowthModifer = 16
)

// func main() {}
func (g *Generator) op_func(buf *bytes.Buffer, node *program.IndexedNode, flags EvalFlags) error {
	funcName, ok := g.LookUpStr(node.Id())
	if !ok {
		return fmt.Errorf("string with id %d is undefined", node.Id())
	}

	var (
		funcType *golang.Type
		params   bytes.Buffer
		body     bytes.Buffer
	)

	length := node.FieldsLength()
	params.Grow(length * paramsGrowthModifer)
	body.Grow(length * bodyGrowthModifer)

	for i := range length {
		var field program.NodeValue
		node.Fields(&field, i)

		if field.Flags()&uint32(golang.ValueFlagFuncMeta) != 0 {
			def, ok := g.LookUpType(field.Type())
			if !ok {
				return fmt.Errorf("type with id %d is undefined", field.Type())
			}

			funcType = def.Type(nil)
			continue
		}

		if field.Flags()&uint32(golang.ValueFlagPointer) == 0 {
			return fmt.Errorf("func node fields can only be pointers")
		}

		// Checks whether the field value's node is valid or not
		target := g.GetNode(field.Value())
		if target == nil {
			return fmt.Errorf("attempt to access undefined node: %d", field.Value())
		}

		switch {
		case field.Flags()&uint32(golang.ValueFlagFuncParam) != 0:
			if params.Len() > 0 {
				params.Write(TokenComma.Bytes())
				params.Write(TokenSpace.Bytes())
			}

			// Evaluates the parameters
			if err := g.evalNode(&params, target, 0); err != nil {
				return err
			}

		case field.Flags()&uint32(golang.ValueFlagFuncBody) != 0:
			body.Write(TokenTab.Bytes())

			// Evaluates the body of the function
			if err := g.evalNode(&body, target, 0); err != nil {
				return err
			}

			body.Write(TokenNewLine.Bytes())
		}
	}

	resultLength := funcType.ResultsLength()

	var declarationLength = TokenFunc.Len() + (resultLength * resultsGrowthModifer) + 8
	buf.Grow(params.Len() + body.Len() + declarationLength)

	// Function declaration
	buf.Write(TokenFunc.Bytes())
	buf.Write(TokenSpace.Bytes())
	buf.Write(funcName)
	buf.Write(TokenParenLeft.Bytes())

	// Parameters
	if funcType != nil && funcType.ParamsLength() > 0 {
		err := g.writePairList(buf, funcType.ParamsLength(), funcType.Params)
		if err != nil {
			return err
		}
	}

	buf.Write(TokenParenRight.Bytes())

	// Return values
	if funcType != nil && funcType.ResultsLength() > 0 {
		buf.Write(TokenSpace.Bytes())

		// look at first result to decide parentheses
		packed := funcType.Results(0)

		//typeId := uint32(packed & 0xffffffff)
		valueId := uint32(packed >> 32)

		name, ok := g.LookUpStr(valueId)
		if !ok {
			return fmt.Errorf("string with id %d is undefined", valueId)
		}

		needParens := funcType.ResultsLength() > 1 || len(name) > 0
		if needParens {
			buf.Write(TokenParenLeft.Bytes())
		}

		err := g.writePairList(buf, funcType.ResultsLength(), funcType.Results)
		if err != nil {
			return err
		}

		if needParens {
			buf.Write(TokenParenRight.Bytes())
		}
	}

	// Body
	buf.Write(TokenSpace.Bytes())
	buf.Write(TokenBraceLeft.Bytes())
	buf.Write(TokenNewLine.Bytes())
	buf.Write(body.Bytes())
	buf.Write(TokenBraceRight.Bytes())

	return nil
}
