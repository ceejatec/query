//  Copyright 2020-Present Couchbase, Inc.
//
//  Use of this software is governed by the Business Source License included in
//  the file licenses/Couchbase-BSL.txt.  As of the Change Date specified in that
//  file, in accordance with the Business Source License, use of this software will
//  be governed by the Apache License, Version 2.0, included in the file
//  licenses/APL.txt.

package errors

import (
	"fmt"
)

const DICT_INTERNAL = 18010

func NewDictInternalError(msg string, e error) Error {
	return &err{level: EXCEPTION, ICode: DICT_INTERNAL, IKey: "dictionary.internal", ICause: e,
		InternalMsg: "Unexpected error in dictionary: " + msg, InternalCaller: CallerN(1)}
}

const DICT_INVALID_INDEXER = 18020

func NewInvalidGSIIndexerError(s string) Error {
	return &err{level: EXCEPTION, ICode: DICT_INVALID_INDEXER, IKey: "dictionary.invalid_indexer",
		InternalMsg: "GSI Indexer does not support collections - " + s, InternalCaller: CallerN(1)}
}

const DICT_INVALID_INDEX = 18030

func NewInvalidGSIIndexError(s string) Error {
	return &err{level: EXCEPTION, ICode: DICT_INVALID_INDEX, IKey: "dictionary.invalid_index",
		InternalMsg: "GSI Index " + s + " does not support collections", InternalCaller: CallerN(1)}
}

const DICT_SYS_COLLECTION = 18040

func NewSystemCollectionError(s string, e error) Error {
	return &err{level: EXCEPTION, ICode: DICT_SYS_COLLECTION, IKey: "dictionary.system_collection", ICause: e,
		InternalMsg: "Error accessing system collection - " + s, InternalCaller: CallerN(1)}
}

const DICT_ENCODING_ERROR = 18050

func NewDictionaryEncodingError(what string, name string, reason error) Error {
	return &err{level: EXCEPTION, ICode: DICT_ENCODING_ERROR, IKey: "dictionary.encoding_error", ICause: reason,
		InternalMsg:    fmt.Sprintf("Cound not %s data dictionary entry for %s due to %v", what, name, reason),
		InternalCaller: CallerN(1)}
}

const DICT_KEYSPACE_MISMATCH_ERROR = 18060

func NewDictKeyspaceMismatchError(ks1, ks2 string) Error {
	return &err{level: EXCEPTION, ICode: DICT_KEYSPACE_MISMATCH_ERROR, IKey: "dictionary.keyspace_mismatch_error",
		InternalMsg:    fmt.Sprintf("Decoded dictionary entry for keyspace %s does not match %s", ks2, ks1),
		InternalCaller: CallerN(1)}
}
