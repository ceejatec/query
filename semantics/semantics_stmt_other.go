//  Copyright 2018-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included in
//  the file licenses/Couchbase-BSL.txt.  As of the Change Date specified in that
//  file, in accordance with the Business Source License, use of this software will
//  be governed by the Apache License, Version 2.0, included in the file
//  licenses/APL.txt.

package semantics

import (
	"github.com/couchbase/query/algebra"
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/errors"
)

func (this *SemChecker) VisitGrantRole(stmt *algebra.GrantRole) (interface{}, error) {
	return nil, stmt.MapExpressions(this)
}

func (this *SemChecker) VisitRevokeRole(stmt *algebra.RevokeRole) (interface{}, error) {
	return nil, stmt.MapExpressions(this)
}

func (this *SemChecker) VisitExplain(stmt *algebra.Explain) (interface{}, error) {
	return stmt.Statement().Accept(this)
}

func (this *SemChecker) VisitAdvise(stmt *algebra.Advise) (interface{}, error) {
	if !this.hasSemFlag(_SEM_ENTERPRISE) {
		return nil, errors.NewEnterpriseFeature("Advise", "semantics.visit_advise")
	}
	switch stmt.Statement().Type() {
	case "SELECT", "DELETE", "MERGE", "UPDATE":
		return stmt.Statement().Accept(this)
	default:
		return nil, errors.NewAdviseUnsupportedStmtError("semantics.visit_advise")
	}
}

func (this *SemChecker) VisitPrepare(stmt *algebra.Prepare) (interface{}, error) {
	return stmt.Statement().Accept(this)
}

func (this *SemChecker) VisitExecute(stmt *algebra.Execute) (interface{}, error) {
	return nil, stmt.MapExpressions(this)
}

func (this *SemChecker) VisitInferKeyspace(stmt *algebra.InferKeyspace) (interface{}, error) {
	return nil, stmt.MapExpressions(this)
}

func (this *SemChecker) VisitUpdateStatistics(stmt *algebra.UpdateStatistics) (interface{}, error) {
	if !this.hasSemFlag(_SEM_ENTERPRISE) {
		return nil, errors.NewEnterpriseFeature("Update Statistics", "semantics.visit_update_statistics")
	}
	if (stmt.IndexAll() || len(stmt.Indexes()) > 0) &&
		(stmt.Using() != datastore.GSI && stmt.Using() != datastore.DEFAULT) {
		return nil, errors.NewUpdateStatInvalidIndexTypeError()
	}
	if stmt.IndexAll() && !stmt.Keyspace().Path().IsCollection() {
		return nil, errors.NewUpdateStatIndexAllCollectionOnly()
	}
	return nil, stmt.MapExpressions(this)
}
