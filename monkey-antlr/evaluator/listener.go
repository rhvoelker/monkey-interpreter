package evaluator

import (
	"fmt"
	"monkey-antlr/object"
	"monkey-antlr/parser"
	"strconv"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

type ProgramListener struct {
	*parser.BaseMonkeyListener

	stack  []object.Object
	errors []string
}

func (l *ProgramListener) Push(item object.Object) {
	l.stack = append(l.stack, item)
}

func (l *ProgramListener) Pop() object.Object {
	if len(l.stack) == 0 {
		return nil
	}

	result := l.stack[len(l.stack)-1]
	l.stack = l.stack[:len(l.stack)-1]

	return result
}

func (l *ProgramListener) ExitIntegerLiteral(c *parser.IntegerLiteralContext) {
	value, err := strconv.ParseInt(c.GetText(), 0, 64)
	if err != nil {
		panic(fmt.Sprintf("Could not parse '%q' as integer. Check grammar.", c.GetText()))
	}
	l.Push(&object.Integer{Value: value})
}

func (l *ProgramListener) ExitBooleanLiteral(c *parser.BooleanLiteralContext) {
	value, err := strconv.ParseBool(c.GetText())
	if err != nil {
		panic(fmt.Sprintf("Could not parse '%q' as boolean. Check grammar.", c.GetText()))
	}
	l.Push(nativeBoolToBoolean(value))
}

func (l *ProgramListener) ExitUnaryOperatorExpression(c *parser.UnaryOperatorExpressionContext) {
	operator := c.GetOp().GetText()
	right := l.Pop()

	switch operator {
	case "!":
		l.Push(evalBangOperatorExpression(right))
	case "-":
		l.Push(evalNegationOperatorExpression(right))
	default:
		l.Push(newError("unknown operator: %s%s", operator, right))
	}
}

func (l *ProgramListener) ExitMulDivBinaryExpression(c *parser.MulDivBinaryExpressionContext) {
	operator := c.GetOp().GetText()
	right, left := l.Pop(), l.Pop()
	l.Push(evalBinaryOperatorExpression(left, operator, right))
}

func (l *ProgramListener) ExitAddSubBinaryExpression(c *parser.AddSubBinaryExpressionContext) {
	operator := c.GetOp().GetText()
	right, left := l.Pop(), l.Pop()
	l.Push(evalBinaryOperatorExpression(left, operator, right))
}

func evalBinaryOperatorExpression(left object.Object, operator string, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerBinaryOperatorExpression(left, operator, right)
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerBinaryOperatorExpression(left object.Object, operator string, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	default:
		return newError("unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalNegationOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("unknown operator: -%s", right.Type())
	}

	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToBoolean(b bool) object.Object {
	if b {
		return TRUE
	}

	return FALSE
}

func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
