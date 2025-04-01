// Copyright 2025 LiveKit, Inc.
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

package sip

import (
	"context"
	"sync"

	apisip "github.com/commcos/msengine/apis/sip"
)

func NewCallState(initial *apisip.SIPCallInfo) *CallState {
	if initial == nil {
		initial = &apisip.SIPCallInfo{}
	}
	s := &CallState{
		info:  initial,
		dirty: true,
	}
	return s
}

type CallState struct {
	mu    sync.Mutex
	info  *apisip.SIPCallInfo
	dirty bool
}

func (s *CallState) DeferUpdate(update func(info *apisip.SIPCallInfo)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = true
	update(s.info)
}

func (s *CallState) Update(ctx context.Context, update func(info *apisip.SIPCallInfo)) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.dirty = true
	update(s.info)
	s.flush(ctx)
}

func (s *CallState) flush(ctx context.Context) {

}

func (s *CallState) Flush(ctx context.Context) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.dirty {
		return
	}
	s.flush(ctx)
}
