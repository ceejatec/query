//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package expression

import (
	"github.com/couchbase/query/value"
)

type IsValued struct {
	UnaryFunctionBase
}

func NewIsValued(operand Expression) Function {
	rv := &IsValued{
		*NewUnaryFunctionBase("isvalued", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsValued) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitIsValued(this)
}

func (this *IsValued) Type() value.Type { return value.BOOLEAN }

func (this *IsValued) Evaluate(item value.Value, context Context) (value.Value, error) {
	arg, err := this.operands[0].Evaluate(item, context)
	if err != nil {
		return nil, err
	}
	switch arg.Type() {
	case value.NULL, value.MISSING:
		return value.FALSE_VALUE, nil
	default:
		return value.TRUE_VALUE, nil
	}
}

func (this *IsValued) PropagatesMissing() bool {
	return false
}

func (this *IsValued) PropagatesNull() bool {
	return false
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For IsValued, simply list this expression.
*/
func (this *IsValued) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

/*
Factory method pattern.
*/
func (this *IsValued) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsValued(operands[0])
	}
}

type IsNotValued struct {
	UnaryFunctionBase
}

func NewIsNotValued(operand Expression) Function {
	rv := &IsNotValued{
		*NewUnaryFunctionBase("isnotvalued", operand),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *IsNotValued) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitIsNotValued(this)
}

func (this *IsNotValued) Type() value.Type { return value.BOOLEAN }

func (this *IsNotValued) Evaluate(item value.Value, context Context) (value.Value, error) {
	arg, err := this.operands[0].Evaluate(item, context)
	if err != nil {
		return nil, err
	}
	switch arg.Type() {
	case value.NULL, value.MISSING:
		return value.TRUE_VALUE, nil
	default:
		return value.FALSE_VALUE, nil
	}
}

func (this *IsNotValued) PropagatesMissing() bool {
	return false
}

func (this *IsNotValued) PropagatesNull() bool {
	return false
}

/*
If this expression is in the WHERE clause of a partial index, lists
the Expressions that are implicitly covered.

For IsNotValued, simply list this expression.
*/
func (this *IsNotValued) FilterCovers(covers map[string]value.Value) map[string]value.Value {
	covers[this.String()] = value.TRUE_VALUE
	return covers
}

/*
Factory method pattern.
*/
func (this *IsNotValued) Constructor() FunctionConstructor {
	return func(operands ...Expression) Function {
		return NewIsNotValued(operands[0])
	}
}
