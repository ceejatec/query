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

/*
Represents multiplication for arithmetic expressions. Type Mult is a
struct that implements CommutativeFunctionBase.
*/
type Mult struct {
	CommutativeFunctionBase
}

func NewMult(operands ...Expression) Function {
	rv := &Mult{
		*NewCommutativeFunctionBase("mult", operands...),
	}

	rv.expr = rv
	return rv
}

/*
Visitor pattern.
*/
func (this *Mult) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitMult(this)
}

func (this *Mult) Type() value.Type { return value.NUMBER }

/*
Range over input arguments, if the type is a number multiply it to
the product. If the value is missing, return a missing value. For
all other types return a null value. Return the final product.
*/
func (this *Mult) Evaluate(item value.Value, context Context) (value.Value, error) {
	null := false
	prod := value.ONE_NUMBER

	for _, op := range this.operands {
		arg, err := op.Evaluate(item, context)
		if err != nil {
			return nil, err
		} else if arg.Type() == value.MISSING {
			return value.MISSING_VALUE, nil
		} else if !null && arg.Type() == value.NUMBER {
			prod = prod.Mult(value.AsNumberValue(arg))
		} else {
			null = true
		}
	}

	if null {
		return value.NULL_VALUE, nil
	}

	return prod, nil
}

/*
Factory method pattern.
*/
func (this *Mult) Constructor() FunctionConstructor {
	return NewMult
}
