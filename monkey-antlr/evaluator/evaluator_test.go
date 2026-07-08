package evaluator

import (
	"monkey-antlr/object"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Eval_Integers(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
		{"-5", -5},
		{"-10", -10},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"-50 + 100 - 50", 0},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"20 + 2 * -10", 0},
		{"50 / 2 * 2 + 10", 60},
		{"2 * (5 + 10)", 30},
		{"3 * 3 * 3 + 10", 37},
		{"3 * (3 * 3) + 10", 37},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func Test_Eval_Booleans(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},
		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},
		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test_Eval_Bang_Operator(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

func Test_Eval_Error_Handling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{"5 + true;", "type mismatch: INTEGER + BOOLEAN"},
		//{"5 + true; 5;", "type mismatch: INTEGER + BOOLEAN"},
		{"-true", "unknown operator: -BOOLEAN"},
		{"true + false", "unknown operator: BOOLEAN + BOOLEAN"},
		//{"5; true + false; 5", "unknown operator: BOOLEAN + BOOLEAN"},
		//		{"if (10 > 1) { true + false; }", "unknown operator: BOOLEAN + BOOLEAN"},
		//		{
		//			`
		//if (10 > 1) {
		//  if (10 > 1) {
		//    return true + false;
		//  }
		//
		//  return 1;
		//}
		//`,
		//			"unknown operator: BOOLEAN + BOOLEAN",
		//		},
		//{"foobar", "identifier not found: foobar"},
		//{`"Hello" - "World"`, "unknown operator: STRING - STRING"},
		//{`"Hello" * "World"`, "unknown operator: STRING * STRING"},
		//{`"Hello" / "World"`, "unknown operator: STRING / STRING"},
		//{`{"name": "Monkey"}[fn(x) { x }];`, "unusable as hash key: FUNCTION"},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testErrorObject(t, evaluated, tt.expectedMessage)
	}
}

func testEval(input string) object.Object {
	eval := New(input, object.NewEnvironment())
	return eval.Eval()
}

func testIntegerObject(t *testing.T, obj object.Object, expected int64) {
	assert.IsType(t, &object.Integer{}, obj)
	result := obj.(*object.Integer).Value
	assert.Equal(t, expected, result)
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) {
	assert.IsType(t, &object.Boolean{}, obj)
	result := obj.(*object.Boolean).Value
	assert.Equal(t, expected, result)
}

func testErrorObject(t *testing.T, obj object.Object, expectedMessage string) {
	assert.IsType(t, &object.Error{}, obj)
	message := obj.(*object.Error).Message
	assert.Equal(t, expectedMessage, message)
}
