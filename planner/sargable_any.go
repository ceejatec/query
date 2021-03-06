//  Copyright 2014-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included in
//  the file licenses/Couchbase-BSL.txt.  As of the Change Date specified in that
//  file, in accordance with the Business Source License, use of this software will
//  be governed by the Apache License, Version 2.0, included in the file
//  licenses/APL.txt.

package planner

import (
	"github.com/couchbase/query/expression"
	base "github.com/couchbase/query/plannerbase"
)

func (this *sargable) VisitAny(pred *expression.Any) (interface{}, error) {
	if this.defaultSargable(pred) {
		return true, nil
	}

	all, ok := this.key.(*expression.All)
	if !ok {
		return false, nil
	}

	array, ok := all.Array().(*expression.Array)
	if !ok {
		bindings := pred.Bindings()
		return len(bindings) == 1 &&
				!bindings[0].Descend() &&
				bindings[0].Expression().EquivalentTo(all.Array()),
			nil
	}

	if !pred.Bindings().SubsetOf(array.Bindings()) {
		return false, nil
	}

	renamer := expression.NewRenamer(pred.Bindings(), array.Bindings())
	satisfies, err := renamer.Map(pred.Satisfies().Copy())
	if err != nil {
		return nil, err
	}

	if array.When() != nil && !base.SubsetOf(satisfies, array.When()) {
		return false, nil
	}

	mappings := expression.Expressions{array.ValueMapping()}
	min, _, _, _ := SargableFor(satisfies, mappings, this.missing, this.gsi)
	return min > 0, nil
}
