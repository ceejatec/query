//  Copyright (c) 2014 Couchbase, Inc.
//  Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file
//  except in compliance with the License. You may obtain a copy of the License at
//    http://www.apache.org/licenses/LICENSE-2.0
//  Unless required by applicable law or agreed to in writing, software distributed under the
//  License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
//  either express or implied. See the License for the specific language governing permissions
//  and limitations under the License.

package algebra

import (
	"encoding/json"

	"github.com/couchbase/query/auth"
	"github.com/couchbase/query/datastore"
	"github.com/couchbase/query/errors"
	"github.com/couchbase/query/expression"
	"github.com/couchbase/query/value"
)

type AlterIndex struct {
	statementBase

	keyspace *KeyspaceRef        `json:"keyspace"`
	name     string              `json:"name"`
	using    datastore.IndexType `json:"using"`
	with     value.Value         `json:"with"`
}

func NewAlterIndex(keyspace *KeyspaceRef, name string, using datastore.IndexType, with value.Value) *AlterIndex {
	rv := &AlterIndex{
		keyspace: keyspace,
		name:     name,
		using:    using,
		with:     with,
	}

	rv.stmt = rv
	return rv
}

func (this *AlterIndex) Accept(visitor Visitor) (interface{}, error) {
	return visitor.VisitAlterIndex(this)
}

func (this *AlterIndex) Signature() value.Value {
	return nil
}

func (this *AlterIndex) Formalize() error {
	return nil
}

func (this *AlterIndex) MapExpressions(mapper expression.Mapper) error {
	return nil
}

func (this *AlterIndex) Expressions() expression.Expressions {
	return nil
}

/*
Returns all required privileges.
*/
func (this *AlterIndex) Privileges() (*auth.Privileges, errors.Error) {
	privs := auth.NewPrivileges()
	fullName := this.keyspace.FullName()
	privs.Add(fullName, auth.PRIV_QUERY_ALTER_INDEX, auth.PRIV_PROPS_NONE)
	return privs, nil
}

func (this *AlterIndex) Keyspace() *KeyspaceRef {
	return this.keyspace
}

func (this *AlterIndex) Name() string {
	return this.name
}

func (this *AlterIndex) Using() datastore.IndexType {
	return this.using
}

func (this *AlterIndex) With() value.Value {
	return this.with
}

func (this *AlterIndex) MarshalJSON() ([]byte, error) {
	r := map[string]interface{}{"type": "alterIndex"}
	r["keyspaceRef"] = this.keyspace
	r["name"] = this.name
	r["using"] = this.using
	if this.with != nil {
		r["with"] = this.with
	}
	return json.Marshal(r)
}

func (this *AlterIndex) Type() string {
	return "ALTER_INDEX"
}
