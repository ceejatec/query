//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package plan

import (
	"encoding/json"

	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/expression/parser"
)

type Limit struct {
	readonly
	optEstimate
	expr expression.Expression
}

func NewLimit(expr expression.Expression, cost, cardinality float64) *Limit {
	rv := &Limit{
		expr: expr,
	}
	setOptEstimate(&rv.optEstimate, cost, cardinality)
	return rv
}

func (this *Limit) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitLimit(this)
}

func (this *Limit) New() Operator {
	return &Limit{}
}

func (this *Limit) Expression() expression.Expression {
	return this.expr
}

func (this *Limit) MarshalJSON() ([]byte, error) {
	return json.Marshal(this.MarshalBase(nil))
}

func (this *Limit) MarshalBase(f func(map[string]interface{})) map[string]interface{} {
	r := map[string]interface{}{"#operator": "Limit"}
	r["expr"] = expression.NewStringer().Visit(this.expr)
	if optEstimate := marshalOptEstimate(&this.optEstimate); optEstimate != nil {
		r["optimizer_estimates"] = optEstimate
	}
	if f != nil {
		f(r)
	}
	return r
}

func (this *Limit) UnmarshalJSON(body []byte) error {
	var _unmarshalled struct {
		_           string             `json:"#operator"`
		Expr        string             `json:"expr"`
		OptEstimate map[string]float64 `json:"optimizer_estimates"`
	}

	err := json.Unmarshal(body, &_unmarshalled)
	if err != nil {
		return err
	}

	this.expr, err = parser.Parse(_unmarshalled.Expr)
	if err != nil {
		return err
	}

	unmarshalOptEstimate(&this.optEstimate, _unmarshalled.OptEstimate)

	return nil
}
