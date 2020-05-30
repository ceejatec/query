//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package execution

import (
	"encoding/json"

	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

type Explain struct {
	base
	plan plan.Operator
}

func NewExplain(plan plan.Operator, context *Context) *Explain {
	rv := &Explain{
		plan: plan,
	}

	newRedirectBase(&rv.base)
	rv.output = rv
	return rv
}

func (this *Explain) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitExplain(this)
}

func (this *Explain) Copy() Operator {
	rv := &Explain{plan: this.plan}
	this.base.copy(&rv.base)
	return rv
}

func (this *Explain) PlanOp() plan.Operator {
	return this.plan
}

func (this *Explain) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer context.Recover(&this.base) // Recover from any panic
		active := this.active()
		defer this.close(context)
		this.switchPhase(_EXECTIME)
		defer this.switchPhase(_NOTIME)
		defer this.notify() // Notify that I have stopped
		if !active {
			return
		}

		bytes, err := this.plan.MarshalJSON()
		if err != nil {
			context.Fatal(errors.NewExplainError(err, "EXPLAIN: Error marshaling JSON."))
			return
		}

		value := value.NewAnnotatedValue(bytes)
		this.sendItem(value)

	})
}

func (this *Explain) MarshalJSON() ([]byte, error) {
	r := this.plan.MarshalBase(func(r map[string]interface{}) {
		this.marshalTimes(r)
		r["plan"] = this.plan
	})
	return json.Marshal(r)
}

func (this *Explain) Done() {
	this.baseDone()
	this.plan = nil
}
