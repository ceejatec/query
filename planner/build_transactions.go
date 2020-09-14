//  Copyright (c) 2020 Couchbase, Inc.
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
	"github.com/couchbase/query/plan"
)

func (this *builder) VisitStartTransaction(stmt *algebra.StartTransaction) (interface{}, error) {
	this.maxParallelism = 1
	return plan.NewSequence(plan.NewStartTransaction(stmt)), nil
}

func (this *builder) VisitCommitTransaction(stmt *algebra.CommitTransaction) (interface{}, error) {
	this.maxParallelism = 1
	return plan.NewSequence(plan.NewCommitTransaction(stmt)), nil
}

func (this *builder) VisitRollbackTransaction(stmt *algebra.RollbackTransaction) (interface{}, error) {
	this.maxParallelism = 1
	return plan.NewSequence(plan.NewRollbackTransaction(stmt)), nil
}

func (this *builder) VisitTransactionIsolation(stmt *algebra.TransactionIsolation) (interface{}, error) {
	this.maxParallelism = 1
	return plan.NewSequence(plan.NewTransactionIsolation(stmt)), nil
}

func (this *builder) VisitSavepoint(stmt *algebra.Savepoint) (interface{}, error) {
	this.maxParallelism = 1
	return plan.NewSequence(plan.NewSavepoint(stmt)), nil
}
