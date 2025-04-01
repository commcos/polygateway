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

package signaling

import (
	"context"
)

type ContextKey int

const (
	// SIPContextKey is the context key for SIP related data
	SIPContextKey ContextKey = iota
)

func WithContext[T any](ctx context.Context, key ContextKey, requestCtx *T) context.Context {
	return context.WithValue(ctx, key, requestCtx)
}

func GetContext[T any](ctx context.Context, key ContextKey) *T {
	if requestCtx, ok := ctx.Value(key).(*T); ok {
		return requestCtx
	}
	return nil
}
