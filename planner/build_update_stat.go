//  Copyright (c) 2018 Couchbase, Inc.
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
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/plan"
)

func (this *builder) VisitUpdateStatistics(stmt *algebra.UpdateStatistics) (interface{}, error) {
	ksref := stmt.Keyspace()
	keyspace, err := this.getNameKeyspace(ksref, false)
	if err != nil {
		return nil, err
	}

	var indexes []datastore.Index
	if len(stmt.Indexes()) > 0 {
		gsiIndexer, err := keyspace.Indexer(datastore.GSI)
		if err != nil {
			return nil, err
		}

		indexes = make([]datastore.Index, 0, len(stmt.Indexes()))
		for _, idxRef := range stmt.Indexes() {
			// only GSI indexes supported, check done in semantics
			index, err := gsiIndexer.IndexByName(idxRef.Name())
			if err != nil {
				return nil, err
			}
			indexes = append(indexes, index)
		}
	}

	return plan.NewUpdateStatistics(keyspace, indexes, stmt), nil
}
