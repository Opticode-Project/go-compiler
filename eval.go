package compiler

import (
	"bytes"
	"fmt"

	"github.com/Opticode-Project/go-compiler/go/golang"
	"github.com/Opticode-Project/go-compiler/program"
	fb "github.com/google/flatbuffers/go"
)

type EvalFlags uint16

const (
	SeperatorTab EvalFlags = 1 << iota
	SeperatorSpace
)

func (g *Generator) Eval(node *program.Node, evalFlags EvalFlags) ([]byte, error) {
	var buf bytes.Buffer

	// Evaluates a node and writes to the buffer
	err := g.evalNode(&buf, node, evalFlags)

	return buf.Bytes(), err
}

func EvalNode(node *program.Node) (any, error) {
	var unionTable fb.Table
	if !node.Node(&unionTable) {
		return nil, fmt.Errorf("failed to access union of node: %d", node.Id())
	}

	switch node.NodeType() {
	case program.NodeUnionIndexedNode:
		n := new(program.IndexedNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		return n, nil

	case program.NodeUnionBinaryNode:
		n := new(program.BinaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		return n, nil

	case program.NodeUnionUnaryNode:
		n := new(program.UnaryNode)
		n.Init(unionTable.Bytes, unionTable.Pos)

		return n, nil

	default:
		return nil, fmt.Errorf("unknown node union type %d", node.NodeType())
	}
}

func (g *Generator) evalNode(buf *bytes.Buffer, node *program.Node, evalFlags EvalFlags) error {
	if node == nil {
		return fmt.Errorf("node is nil")
	}

	opcode := golang.Opcode(node.Opcode())
	if node.NodeType() == program.NodeUnionNONE {
		return fmt.Errorf("node %d has no union payload", node.Id())
	}

	v, err := EvalNode(node)
	if err != nil {
		return err
	}

	switch n := v.(type) {
	case *program.IndexedNode:
		err := g.EvalIndexed(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}

	case *program.BinaryNode:
		err := g.EvalBinary(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}

	case *program.UnaryNode:
		err := g.EvalUnary(buf, opcode, n, evalFlags)
		if err != nil {
			return err
		}
	}

	return nil
}

func (g *Generator) EvalIndexed(buf *bytes.Buffer, opcode golang.Opcode, node *program.IndexedNode, evalFlags EvalFlags) error {
	switch opcode {
	case golang.OpcodePackage:
		return g.op_package(buf, node, evalFlags)
	case golang.OpcodeImport:
		return g.op_import(buf, node, evalFlags)
	/*case golang.OpcodeConst:
		return g.op_const(buf, node, evalFlags)
	case golang.OpcodeVar:
		return g.op_var(buf, node, evalFlags)*/
	case golang.OpcodeIf:
		return g.op_if(buf, node, evalFlags)
		/*case golang.OpcodeFunc:
			return g.op_func(buf, node, evalFlags)
		case golang.OpcodeCall:
			return g.op_call(buf, node, evalFlags)
		case golang.OpcodeType:
			return g.op_type(buf, node, evalFlags)
		case golang.OpcodeReturn:
			return g.op_return(buf, node, evalFlags)*/
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalBinary(buf *bytes.Buffer, opcode golang.Opcode, node *program.BinaryNode, evalFlags EvalFlags) error {
	switch opcode {
	case golang.OpcodeImportValue:
		return g.op_importValue(buf, node, evalFlags)
	/*case golang.OpcodeConstValue:
		return g.op_constValue(buf, node, evalFlags)
	case golang.OpcodeVarValue:
		return g.op_varValue(buf, node, evalFlags)*/

	case golang.OpcodeEqual:
		return g.op_binary(buf, node, TokenCompare, evalFlags)
	case golang.OpcodeNotEqual:
		return g.op_binary(buf, node, TokenNotEqual, evalFlags)
	case golang.OpcodeLess:
		return g.op_binary(buf, node, TokenLess, evalFlags)
	case golang.OpcodeLessEqual:
		return g.op_binary(buf, node, TokenLessEqual, evalFlags)
	case golang.OpcodeGreater:
		return g.op_binary(buf, node, TokenGreater, evalFlags)
	case golang.OpcodeGreaterEqual:
		return g.op_binary(buf, node, TokenGreaterEqual, evalFlags)
	case golang.OpcodeAnd:
		return g.op_binary(buf, node, TokenAnd, evalFlags)
	case golang.OpcodeOr:
		return g.op_binary(buf, node, TokenOr, evalFlags)

	case golang.OpcodeAdd:
		return g.op_binary(buf, node, TokenPlus, evalFlags)
	case golang.OpcodeSub:
		return g.op_binary(buf, node, TokenMinus, evalFlags)
	case golang.OpcodeMul:
		return g.op_binary(buf, node, TokenStar, evalFlags)
	case golang.OpcodeDiv:
		return g.op_binary(buf, node, TokenSlash, evalFlags)
	case golang.OpcodeMod:
		return g.op_binary(buf, node, TokenModulus, evalFlags)
	case golang.OpcodeAssign:
		return g.op_binary(buf, node, TokenEqual, evalFlags)
	case golang.OpcodeAddAssign:
		return g.op_binary(buf, node, TokenAddAssign, evalFlags)
	case golang.OpcodeSubAssign:
		return g.op_binary(buf, node, TokenSubAssign, evalFlags)
	case golang.OpcodeMulAssign:
		return g.op_binary(buf, node, TokenMulAssign, evalFlags)
	case golang.OpcodeDivAssign:
		return g.op_binary(buf, node, TokenDivAssign, evalFlags)
	case golang.OpcodeModAssign:
		return g.op_binary(buf, node, TokenModAssign, evalFlags)

	case golang.OpcodeBitAndAssign:
		return g.op_binary(buf, node, TokenBitAndAssign, evalFlags)
	case golang.OpcodeBitOrAssign:
		return g.op_binary(buf, node, TokenBitOrAssign, evalFlags)
	case golang.OpcodeBitXorAssign:
		return g.op_binary(buf, node, TokenBitXorAssign, evalFlags)
	case golang.OpcodeBitClearAssign:
		return g.op_binary(buf, node, TokenBitClearAssign, evalFlags)
	case golang.OpcodeLeftShiftAssign:
		return g.op_binary(buf, node, TokenShiftLeftAssign, evalFlags)
	case golang.OpcodeRightShiftAssign:
		return g.op_binary(buf, node, TokenShiftRightAssign, evalFlags)
	case golang.OpcodeBitAnd:
		return g.op_binary(buf, node, TokenBitAnd, evalFlags)
	case golang.OpcodeBitOr:
		return g.op_binary(buf, node, TokenBitOr, evalFlags)
	case golang.OpcodeBitXor:
		return g.op_binary(buf, node, TokenBitXor, evalFlags)
	case golang.OpcodeBitClear:
		return g.op_binary(buf, node, TokenBitClear, evalFlags)
	case golang.OpcodeLeftShift:
		return g.op_binary(buf, node, TokenShiftLeft, evalFlags)
	case golang.OpcodeRightShift:
		return g.op_binary(buf, node, TokenShiftRight, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func (g *Generator) EvalUnary(buf *bytes.Buffer, opcode golang.Opcode, node *program.UnaryNode, evalFlags EvalFlags) error {
	switch opcode {
	case golang.OpcodeNot:
		return g.op_unaryPrefix(buf, node, TokenNot, evalFlags)
	/*case golang.OpcodeDefer:
		return g.op_defer(buf, node, evalFlags)
	case golang.OpcodeGoRoutine:
		return g.op_goRoutine(buf, node, evalFlags)*/

	case golang.OpcodeInc:
		return g.op_unarySuffix(buf, node, TokenIncrement, evalFlags)
	case golang.OpcodeDec:
		return g.op_unarySuffix(buf, node, TokenDecrement, evalFlags)
	case golang.OpcodeAddrOf:
		return g.op_unaryPrefix(buf, node, TokenBitAnd, evalFlags)
	case golang.OpcodeDeref:
		return g.op_unaryPrefix(buf, node, TokenStar, evalFlags)

	case golang.OpcodeReceive:
		return g.op_unaryPrefix(buf, node, TokenArrowLeft, evalFlags)
	}

	return fmt.Errorf("invalid opcode on node with opcode of %s", opcode)
}

func isConstOperator(op golang.Opcode) bool {
	switch {
	// arithmetic
	case op >= golang.OpcodeAdd && op <= golang.OpcodeMod:

	// comparisons
	case op >= golang.OpcodeEqual && op <= golang.OpcodeGreaterEqual:

	// logical
	case op >= golang.OpcodeAnd && op <= golang.OpcodeNot:

	// bitwise
	case op >= golang.OpcodeBitAnd && op <= golang.OpcodeRightShift:
		return true
	}

	return false
}

func (g *Generator) isConstValue(v *program.NodeValue) bool {
	if v == nil {
		return false
	}

	// literal value
	if v.Flags()&uint32(golang.ValueFlagPointer) == 0 {
		return true
	}

	// pointer -> must recurse
	node := g.GetNode(v.Value())
	if node == nil {
		return false
	}

	return g.isConstantExpression(node)
}

func (g *Generator) isConstantExpression(node *program.Node) bool {
	if node == nil {
		return false
	}

	op := golang.Opcode(node.Opcode())

	// only allow const-safe operators
	if !isConstOperator(op) {
		return false
	}

	v, err := EvalNode(node)
	if err != nil {
		return false
	}

	switch n := v.(type) {
	case *program.UnaryNode:
		value := n.Value(nil)

		return g.isConstValue(value)

	case *program.BinaryNode:
		left := n.Left(nil)
		right := n.Right(nil)

		return g.isConstValue(left) && g.isConstValue(right)

	default:
		return false
	}
}

// Evaluates a pointer of a node value and writes the result to the buffer
func (g *Generator) evalValue(buf *bytes.Buffer, nodeValue *program.NodeValue, isConst bool) error {
	if nodeValue == nil {
		return fmt.Errorf("node value is null")
	}

	// if value is a pointer
	if nodeValue.Flags()&uint32(golang.ValueFlagPointer) != 0 {
		if isConst {
			return fmt.Errorf("const value cannot reference runtime expression")
		}

		// Checks whether the node value is valid or not
		node := g.GetNode(nodeValue.Value())
		if node == nil {
			return fmt.Errorf("attempt to access undefined node: %d", nodeValue.Value())
		}

		// Evaluates the node and writes to the buffer
		err := g.evalNode(buf, node, 0)
		if err != nil {
			return err
		}

		return nil
	}

	// literal value
	value, ok := g.LookUpStr(uint32(nodeValue.Value()))
	if !ok {
		return fmt.Errorf("string with id %d is undefined", nodeValue.Value())
	}

	if nodeValue.Flags()&uint32(golang.ValueFlagQuotation) != 0 {
		buf.Write(TokenQuotation.Bytes())
		buf.Write(value)
		buf.Write(TokenQuotation.Bytes())

		return nil
	}

	buf.Write(value)
	return nil
}
