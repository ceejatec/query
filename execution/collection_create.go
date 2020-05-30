//  Copyright (c) 2020 Couchbase, Inc.
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

	//	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

type CreateCollection struct {
	base
	plan *plan.CreateCollection
}

func NewCreateCollection(plan *plan.CreateCollection, context *Context) *CreateCollection {
	rv := &CreateCollection{
		plan: plan,
	}

	newRedirectBase(&rv.base)
	rv.output = rv
	return rv
}

func (this *CreateCollection) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitCreateCollection(this)
}

func (this *CreateCollection) Copy() Operator {
	rv := &CreateCollection{plan: this.plan}
	this.base.copy(&rv.base)
	return rv
}

func (this *CreateCollection) PlanOp() plan.Operator {
	return this.plan
}

func (this *CreateCollection) RunOnce(context *Context, parent value.Value) {
	this.once.Do(func() {
		defer context.Recover(&this.base) // Recover from any panic
		active := this.active()
		defer this.close(context)
		this.switchPhase(_EXECTIME)
		defer this.switchPhase(_NOTIME)
		defer this.notify() // Notify that I have stopped

		if !active || context.Readonly() {
			return
		}

		// Actually create collection
		this.switchPhase(_SERVTIME)
		err := this.plan.Scope().CreateCollection(this.plan.Node().Name())
		if err != nil {
			context.Error(err)
		}
	})
}

func (this *CreateCollection) MarshalJSON() ([]byte, error) {
	r := this.plan.MarshalBase(func(r map[string]interface{}) {
		this.marshalTimes(r)
	})
	return json.Marshal(r)
}
