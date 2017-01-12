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
	"fmt"
	"math"
	"time"

	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

// Send to keyspace
type SendUpdate struct {
	base
	plan  *plan.SendUpdate
	limit int64
}

func NewSendUpdate(plan *plan.SendUpdate) *SendUpdate {
	rv := &SendUpdate{
		base:  newBase(),
		plan:  plan,
		limit: -1,
	}

	rv.output = rv
	return rv
}

func (this *SendUpdate) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitSendUpdate(this)
}

func (this *SendUpdate) Copy() Operator {
	return &SendUpdate{this.base.copy(), this.plan, this.limit}
}

func (this *SendUpdate) RunOnce(context *Context, parent value.Value) {
	this.phaseTimes = func(d time.Duration) { context.AddPhaseTime(UPDATE, d) }
	this.runConsumer(this, context, parent)
}

func (this *SendUpdate) processItem(item value.AnnotatedValue, context *Context) bool {
	rv := this.limit != 0 && this.enbatch(item, this, context)

	if this.limit > 0 {
		this.limit--
	}

	return rv
}

func (this *SendUpdate) beforeItems(context *Context, parent value.Value) bool {
	if this.plan.Limit() == nil {
		return true
	}

	limit, err := this.plan.Limit().Evaluate(parent, context)
	if err != nil {
		context.Error(errors.NewEvaluationError(err, "LIMIT clause"))
		return false
	}

	l := limit.ActualForIndex() // Exact number
	switch l := l.(type) {
	case int64:
		this.limit = l
		return true
	case float64:
		if math.Trunc(l) == l {
			this.limit = int64(l)
			return true
		}
	}

	context.Error(errors.NewInvalidValueError(fmt.Sprintf("Invalid LIMIT %v of type %T.", l, l)))
	return false
}

func (this *SendUpdate) afterItems(context *Context) {
	this.flushBatch(context)
}

func (this *SendUpdate) flushBatch(context *Context) bool {
	defer this.releaseBatch()

	if len(this.batch) == 0 {
		return true
	}

	var pairs []value.Pair
	if _UPDATE_POOL.Size() >= len(this.batch) {
		pairs = _UPDATE_POOL.Get()
		defer _UPDATE_POOL.Put(pairs)
	} else {
		pairs = make([]value.Pair, 0, len(this.batch))
	}

	for i, item := range this.batch {
		uv, ok := item.Field(this.plan.Alias())
		if !ok {
			context.Error(errors.NewUpdateAliasMissingError(this.plan.Alias()))
			return false
		}

		av, ok := uv.(value.AnnotatedValue)
		if !ok {
			context.Error(errors.NewUpdateAliasMetadataError(this.plan.Alias()))
			return false
		}

		key, ok := this.requireKey(av, context)
		if !ok {
			return false
		}

		pairs = pairs[0 : i+1]
		pairs[i].Name = key

		clone := item.GetAttachment("clone")
		switch clone := clone.(type) {
		case value.AnnotatedValue:
			cv, ok := clone.Field(this.plan.Alias())
			if !ok {
				context.Error(errors.NewUpdateAliasMissingError(this.plan.Alias()))
				return false
			}

			cav := value.NewAnnotatedValue(cv)
			cav.SetAnnotations(av)
			pairs[i].Value = cav
			item.SetField(this.plan.Alias(), cav)
		default:
			context.Error(errors.NewInvalidValueError(fmt.Sprintf(
				"Invalid UPDATE value of type %T.", clone)))
			return false
		}
	}

	this.switchPhase(_SERVTIME)

	pairs, e := this.plan.Keyspace().Update(pairs)

	this.switchPhase(_EXECTIME)

	// Update mutation count with number of updated docs
	context.AddMutationCount(uint64(len(pairs)))

	if e != nil {
		context.Error(e)
	}

	for _, item := range this.batch {
		if !this.sendItem(item) {
			return false
		}
	}

	return true
}

func (this *SendUpdate) readonly() bool {
	return false
}

func (this *SendUpdate) MarshalJSON() ([]byte, error) {
	r := this.plan.MarshalBase(func(r map[string]interface{}) {
		this.marshalTimes(r)
	})
	return json.Marshal(r)
}

var _UPDATE_POOL = value.NewPairPool(_BATCH_SIZE)
