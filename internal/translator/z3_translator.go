// Package translator содержит реализацию транслятора в Z3
package translator

import (
	"fmt"
	"symbolic-execution-course/internal/symbolic"

	"github.com/ebukreev/go-z3/z3"
)

// Z3Translator транслирует символьные выражения в Z3 формулы
type Z3Translator struct {
	ctx    *z3.Context
	config *z3.Config
	vars   map[string]z3.Value // Кэш переменных
}

// NewZ3Translator создаёт новый экземпляр Z3 транслятора
func NewZ3Translator() *Z3Translator {
	config := &z3.Config{}
	ctx := z3.NewContext(config)

	return &Z3Translator{
		ctx:    ctx,
		config: config,
		vars:   make(map[string]z3.Value),
	}
}

// GetContext возвращает Z3 контекст
func (zt *Z3Translator) GetContext() interface{} {
	return zt.ctx
}

// Reset сбрасывает состояние транслятора
func (zt *Z3Translator) Reset() {
	zt.vars = make(map[string]z3.Value)
}

// Close освобождает ресурсы
func (zt *Z3Translator) Close() {
	// Z3 контекст закрывается автоматически
}

// TranslateExpression транслирует символьное выражение в Z3
func (zt *Z3Translator) TranslateExpression(expr symbolic.SymbolicExpression) (interface{}, error) {
	return expr.Accept(zt), nil
}

// TODO: Реализуйте следующие методы в рамках домашнего задания

// VisitVariable транслирует символьную переменную в Z3
func (zt *Z3Translator) VisitVariable(expr *symbolic.SymbolicVariable) interface{} {
	if value, exists := zt.vars[expr.Name]; exists {
		return value
	}

	value := zt.createZ3Variable(expr.Name, expr.Type())

	zt.vars[expr.Name] = value

	return value
	// Подсказки:
	// - Используйте zt.ctx.IntConst(name) для int переменных
	// - Используйте zt.ctx.BoolConst(name) для bool переменных
	// - Храните переменные в zt.vars для повторного использования
}

// VisitIntConstant транслирует целочисленную константу в Z3
func (zt *Z3Translator) VisitIntConstant(expr *symbolic.IntConstant) interface{} {
	value := zt.ctx.FromInt(
		expr.Value,
		zt.ctx.IntSort(),
	)

	return value.(z3.Int)
}

// VisitBoolConstant транслирует булеву константу в Z3
func (zt *Z3Translator) VisitBoolConstant(expr *symbolic.BoolConstant) interface{} {
	return zt.ctx.FromBool(expr.Value)
}

// VisitBinaryOperation транслирует бинарную операцию в Z3
func (zt *Z3Translator) VisitBinaryOperation(expr *symbolic.BinaryOperation) interface{} {
	left := expr.Left.Accept(zt)
	right := expr.Right.Accept(zt)

	switch expr.Operator {
	case symbolic.ADD:
		return left.(z3.Int).Add(right.(z3.Int))

	case symbolic.SUB:
		return left.(z3.Int).Sub(right.(z3.Int))

	case symbolic.MUL:
		return left.(z3.Int).Mul(right.(z3.Int))

	case symbolic.DIV:
		return left.(z3.Int).Div(right.(z3.Int))

	case symbolic.MOD:
		return left.(z3.Int).Mod(right.(z3.Int))

	case symbolic.LT:
		return left.(z3.Int).LT(right.(z3.Int))

	case symbolic.LE:
		return left.(z3.Int).LE(right.(z3.Int))

	case symbolic.GT:
		return left.(z3.Int).GT(right.(z3.Int))

	case symbolic.GE:
		return left.(z3.Int).GE(right.(z3.Int))

	case symbolic.EQ:
		return translateEquality(
			left,
			right,
			false,
		)

	case symbolic.NE:
		return translateEquality(
			left,
			right,
			true,
		)

	default:
		panic(fmt.Sprintf(
			"unsupported binary operator: %s",
			expr.Operator,
		))
	}

	// Подсказки по операциям в Z3:
	// - Арифметические: left.Add(right), left.Sub(right), left.Mul(right), left.Div(right)
	// - Сравнения: left.Eq(right), left.LT(right), left.LE(right), etc.
	// - Приводите типы: left.(z3.Int), right.(z3.Int) для int операций
}

func equality[T interface {
	Eq(T) z3.Bool
	NE(T) z3.Bool
}](
	left, right T,
	notEqual bool,
) z3.Bool {
	if notEqual {
		return left.NE(right)
	}
	return left.Eq(right)
}

func translateEquality(
	left, right interface{},
	notEqual bool,
) z3.Bool {
	switch l := left.(type) {
	case z3.Int:
		return equality(l, right.(z3.Int), notEqual)

	case z3.Bool:
		return equality(l, right.(z3.Bool), notEqual)

	default:
		panic(fmt.Sprintf("unsupported equality type: %T", left))
	}
}

// VisitLogicalOperation транслирует логическую операцию в Z3
func (zt *Z3Translator) VisitLogicalOperation(expr *symbolic.LogicalOperation) interface{} {
	operands := make([]z3.Bool, len(expr.Operands))

	for i, operand := range expr.Operands {
		value := operand.Accept(zt)

		boolValue, ok := value.(z3.Bool)
		if !ok {
			panic(fmt.Sprintf(
				"logical operand %d is not bool",
				i,
			))
		}

		operands[i] = boolValue
	}

	switch expr.Operator {
	case symbolic.AND:
		return operands[0].And(operands[1:]...)

	case symbolic.OR:
		return operands[0].Or(operands[1:]...)

	case symbolic.NOT:
		return operands[0].Not()

	case symbolic.IMPLIES:
		return operands[0].Implies(operands[1])

	default:
		panic(fmt.Sprintf(
			"unsupported logical operator: %s",
			expr.Operator,
		))
	}

	// Подсказки:
	// - AND: zt.ctx.And(operands...)
	// - OR: zt.ctx.Or(operands...)
	// - NOT: operand.Not() (для единственного операнда)
	// - IMPLIES: antecedent.Implies(consequent)
}

// Вспомогательные методы

// createZ3Variable создаёт Z3 переменную соответствующего типа
func (zt *Z3Translator) createZ3Variable(name string, exprType symbolic.ExpressionType) z3.Value {
	switch exprType {
	case symbolic.IntType:
		return zt.ctx.IntConst(name)

	case symbolic.BoolType:
		return zt.ctx.BoolConst(name)

	default:
		panic(fmt.Sprintf("Unsupported symbolic type: %s", exprType))
	}
}

// castToZ3Type приводит значение к нужному Z3 типу
func (zt *Z3Translator) castToZ3Type(value interface{}, targetType symbolic.ExpressionType) (z3.Value, error) {
	switch targetType {
	case symbolic.IntType:
		v, ok := value.(z3.Int)
		if !ok {
			return nil, fmt.Errorf(
				"expected z3.Int, got %T",
				value,
			)
		}

		return v, nil

	case symbolic.BoolType:
		v, ok := value.(z3.Bool)
		if !ok {
			return nil, fmt.Errorf(
				"expected z3.Bool, got %T",
				value,
			)
		}

		return v, nil

	default:
		return nil, fmt.Errorf(
			"unsupported symbolic type: %s",
			targetType,
		)
	}
}
