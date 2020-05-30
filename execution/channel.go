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
	"github.com/couchbase/query/plan"
	"github.com/couchbase/query/value"
)

// Dummy operator that simply wraps an item channel.
type Channel struct {
	base
}

func NewChannel(context *Context) *Channel {
	rv := &Channel{}
	newBase(&rv.base, context)
	rv.dormant()
	rv.output = rv
	return rv
}

func (this *Channel) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitChannel(this)
}

func (this *Channel) Copy() Operator {
	rv := &Channel{}
	this.base.copy(&rv.base)
	return rv
}

func (this *Channel) PlanOp() plan.Operator {
	return nil
}

// This operator is a no-op. It simply provides a shared itemChannel.
func (this *Channel) RunOnce(context *Context, parent value.Value) {
}

func (this *Channel) MarshalJSON() ([]byte, error) {

	// there's no corresponding plan.Channel, so we have a dummy
	return nil, nil
}
