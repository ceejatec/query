//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package planner

import (
	"github.com/couchbase/query/algebra"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/plan"
	base "github.com/couchbase/query/plannerbase"
)

func (this *builder) VisitMerge(stmt *algebra.Merge) (interface{}, error) {
	this.node = stmt
	this.children = make([]plan.Operator, 0, 8)
	this.subChildren = make([]plan.Operator, 0, 8)
	source := stmt.Source()

	this.baseKeyspaces = make(map[string]*base.BaseKeyspace, _MAP_KEYSPACE_CAP)
	sourceKeyspace := base.NewBaseKeyspace(source.Alias(), source.Keyspace())
	this.baseKeyspaces[sourceKeyspace.Name()] = sourceKeyspace
	targetKeyspace := base.NewBaseKeyspace(stmt.KeyspaceRef().Alias(), stmt.KeyspaceRef().Keyspace())
	this.baseKeyspaces[targetKeyspace.Name()] = targetKeyspace
	this.collectKeyspaceNames()

	var left algebra.SimpleFromTerm
	var err error
	outer := false

	if !stmt.IsOnKey() {
		// use outer join if INSERT action is specified
		if stmt.Actions().Insert() != nil {
			outer = true
		} else {
			_, err = this.processPredicate(stmt.On(), true)
			if err != nil {
				return nil, err
			}

			this.pushableOnclause = stmt.On()
		}
	}

	this.initialIndexAdvisor(stmt)
	this.extractPredicates(nil, this.pushableOnclause)

	if source.SubqueryTerm() != nil {
		_, err := source.SubqueryTerm().Accept(this)
		if err != nil && !this.indexAdvisor {
			return nil, err
		}

		left = source.SubqueryTerm()
	} else if source.ExpressionTerm() != nil {
		_, err := source.ExpressionTerm().Accept(this)
		if err != nil && !this.indexAdvisor {
			return nil, err
		}

		left = source.ExpressionTerm()
	} else {
		if source.From() == nil {
			// should have caught in semantics check
			return nil, errors.NewPlanInternalError("VisitMerge: MERGE missing source.")
		}

		_, err := source.From().Accept(this)
		if err != nil && !this.indexAdvisor {
			return nil, err
		}

		left = source.From()
	}

	ksref := stmt.KeyspaceRef()
	ksref.SetDefaultNamespace(this.namespace)

	keyspace, err := this.getNameKeyspace(ksref)
	if err != nil {
		return nil, err
	}

	actions := stmt.Actions()
	var update, delete, insert plan.Operator

	if actions.Update() != nil {
		act := actions.Update()
		ops := make([]plan.Operator, 0, 5)

		if act.Where() != nil {
			filter := this.addMergeFilter(act.Where())
			ops = append(ops, filter)
		}

		ops = append(ops, plan.NewClone(ksref.Alias()))

		if act.Set() != nil {
			ops = append(ops, plan.NewSet(act.Set()))
		}

		if act.Unset() != nil {
			ops = append(ops, plan.NewUnset(act.Unset()))
		}

		ops = append(ops, plan.NewSendUpdate(keyspace, ksref, stmt.Limit()))
		update = plan.NewSequence(ops...)
	}

	if actions.Delete() != nil {
		act := actions.Delete()
		ops := make([]plan.Operator, 0, 4)

		if act.Where() != nil {
			filter := this.addMergeFilter(act.Where())
			ops = append(ops, filter)
		}

		ops = append(ops, plan.NewSendDelete(keyspace, ksref, stmt.Limit()))
		delete = plan.NewSequence(ops...)
	}

	if actions.Insert() != nil {
		act := actions.Insert()
		ops := make([]plan.Operator, 0, 4)

		if act.Where() != nil {
			filter := this.addMergeFilter(act.Where())
			ops = append(ops, filter)
		}

		var keyExpr expression.Expression
		if stmt.IsOnKey() {
			keyExpr = stmt.On()
		} else {
			keyExpr = act.Key()
		}
		ops = append(ops, plan.NewSendInsert(keyspace, ksref, keyExpr, act.Value(),
			act.Options(), stmt.Limit()))
		insert = plan.NewSequence(ops...)
	}

	if stmt.IsOnKey() {
		this.addSubChildren(plan.NewMerge(keyspace, ksref, stmt.On(), update, delete, insert))
	} else {
		// use ANSI JOIN to handle the ON-clause
		right := algebra.NewKeyspaceTermFromPath(ksref.Path(), ksref.As(), nil, stmt.Indexes())
		right.SetAnsiJoin()
		algebra.TransferJoinHint(right, left)

		ansiJoin := algebra.NewAnsiJoin(left, outer, right, stmt.On())
		join, err := this.buildAnsiJoin(ansiJoin)
		if err != nil {
			return nil, err
		}

		switch join := join.(type) {
		case *plan.NLJoin:
			this.addSubChildren(join)
		case *plan.Join, *plan.HashJoin:
			if len(this.subChildren) > 0 {
				this.addChildren(this.addSubchildrenParallel())
			}
			this.addChildren(join)
		}

		this.addSubChildren(plan.NewMerge(keyspace, ksref, nil, update, delete, insert))
	}

	if stmt.Returning() != nil {
		this.subChildren = this.buildDMLProject(stmt.Returning(), this.subChildren)
	}

	this.addChildren(this.addSubchildrenParallel())

	if stmt.Limit() != nil {
		this.addChildren(plan.NewLimit(stmt.Limit(), OPT_COST_NOT_AVAIL, OPT_CARD_NOT_AVAIL))
	}

	if stmt.Returning() == nil {
		cost := OPT_COST_NOT_AVAIL
		cardinality := OPT_CARD_NOT_AVAIL
		lastOp := this.lastOp
		if lastOp != nil {
			cost = lastOp.Cost()
			cardinality = lastOp.Cardinality()
		}
		this.addChildren(plan.NewDiscard(cost, cardinality))
	}

	return plan.NewSequence(this.children...), nil
}

func (this *builder) addMergeFilter(pred expression.Expression) *plan.Filter {
	cost := float64(OPT_COST_NOT_AVAIL)
	cardinality := float64(OPT_CARD_NOT_AVAIL)

	if this.useCBO {
		cost, cardinality = getFilterCost(this.lastOp, pred, this.baseKeyspaces,
			this.keyspaceNames)
	}

	return plan.NewFilter(pred, cost, cardinality)
}
