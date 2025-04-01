/********************************************************************
* Copyright (c) All Rights Reserved.
*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
*         http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*******************************************************************/

package errors

type ErrorCode string

func (e ErrorCode) Error() string {
	return string(e)
}

const (
	OK ErrorCode = ""

	// Request Canceled by client
	Canceled ErrorCode = "canceled"
	// Could not unmarshal request
	MalformedRequest ErrorCode = "malformed_request"
	// Could not unmarshal result
	MalformedResponse ErrorCode = "malformed_result"
	// Request timed out
	DeadlineExceeded ErrorCode = "deadline_exceeded"
	// Service unavailable due to load and/or affinity constraints
	Unavailable ErrorCode = "unavailable"
	// Unknown (server returned non-psrpc error)
	Unknown ErrorCode = "unknown"

	// Invalid argument in request
	InvalidArgument ErrorCode = "invalid_argument"
	// Entity not found
	NotFound ErrorCode = "not_found"
	// Cannot produce and entity matching requested format
	NotAcceptable ErrorCode = "not_acceptable"
	// Duplicate creation attempted
	AlreadyExists ErrorCode = "already_exists"
	// Caller does not have required permissions
	PermissionDenied ErrorCode = "permission_denied"
	// Some resource has been exhausted, e.g. memory or quota
	ResourceExhausted ErrorCode = "resource_exhausted"
	// Inconsistent state to carry out request
	FailedPrecondition ErrorCode = "failed_precondition"
	// Request aborted
	Aborted ErrorCode = "aborted"
	// Operation was out of range
	OutOfRange ErrorCode = "out_of_range"
	// Operation is not implemented by the server
	Unimplemented ErrorCode = "unimplemented"
	// Operation failed due to an internal error
	Internal ErrorCode = "internal"
	// Irrecoverable loss or corruption of data
	DataLoss ErrorCode = "data_loss"
	// Similar to PermissionDenied, used when the caller is unidentified
	Unauthenticated ErrorCode = "unauthenticated"
)

type sipError struct {
	error
	code    ErrorCode
	details []any
}

func (e sipError) Code() ErrorCode {
	return e.code
}

func (e sipError) Details() []any {
	return e.details
}

func (e sipError) Unwrap() []error {
	return []error{e.error, e.code}
}
