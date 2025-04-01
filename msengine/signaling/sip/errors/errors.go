// Copyright 2023 LiveKit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package errors

import (
	"fmt"

	"google.golang.org/protobuf/proto"
)

var (
	ErrNoConfig    = NewErrorf(InvalidArgument, "missing config")
	ErrUnavailable = NewErrorf(Unavailable, "cpu exhausted")
)

type Error interface {
	error
	Code() ErrorCode
	Details() []any
}

func ErrCouldNotParseConfig(err error) Error {
	return NewErrorf(InvalidArgument, "could not parse config: %v", err)
}

func NewError(code ErrorCode, err error, details ...proto.Message) Error {
	if err == nil {
		panic("error is nil")
	}
	var protoDetails []any
	for _, e := range details {
		protoDetails = append(protoDetails, e)
	}
	return &sipError{
		error:   err,
		code:    code,
		details: protoDetails,
	}
}

func NewErrorf(code ErrorCode, msg string, args ...interface{}) Error {
	return &sipError{
		error: fmt.Errorf(msg, args...),
		code:  code,
	}
}
