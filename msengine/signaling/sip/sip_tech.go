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

package sip

import (
	"context"
	"fmt"
	"log/slog"
	"net/netip"
	"strings"
	"sync"

	"github.com/livekit/sipgo"

	apisip "github.com/commcos/msengine/apis/sip"
	"github.com/commcos/msengine/media"
	"github.com/commcos/msengine/signaling"
	"github.com/commcos/msengine/signaling/sip/config"
	"github.com/commcos/msengine/signaling/sip/stats"
	"github.com/commcos/msengine/version"
)

type ServiceConfig struct {
	SignalingIP      netip.Addr
	SignalingIPLocal netip.Addr
}

type sipProtocolTech struct {
	conf  *config.Config
	sconf *ServiceConfig
	mon   *stats.Monitor
	cli   *Client
	srv   *Server

	mu               sync.Mutex
	pendingTransfers map[transferKey]chan struct{}
	observer         signaling.Observer
}

type transferKey struct {
	SipCallId  string
	TransferTo string
}

func NewSIPTech(region string, conf *config.Config, mon *stats.Monitor) (signaling.ProtocolTech, error) {
	s := &sipProtocolTech{
		conf:             conf,
		mon:              mon,
		cli:              NewClient(region, conf, mon),
		srv:              NewServer(region, conf, mon),
		pendingTransfers: make(map[transferKey]chan struct{}),
	}
	var err error
	s.sconf, err = GetServiceConfig(s.conf)
	if err != nil {
		return nil, err
	}

	s.conf.SIPHostname = strings.ReplaceAll(
		s.conf.SIPHostname,
		"${IP}",
		strings.NewReplacer(
			".", "-", // IPv4
			"[", "", "]", "", ":", "-", // IPv6
		).Replace(s.sconf.SignalingIP.String()),
	)
	if strings.ContainsAny(s.conf.SIPHostname, "$%{}[]:/| ") {
		return nil, fmt.Errorf("invalid hostname: %q", s.conf.SIPHostname)
	}
	if s.conf.SIPHostname != "" {
		slog.Info("using hostname", "hostname", s.conf.SIPHostname)
	}
	return s, nil
}

func (s *sipProtocolTech) GetProtocolMetadata() signaling.ProtocolMetadata {
	return signaling.ProtocolMetadata{
		ProtocolType: signaling.ProtocolTypeSIP,
	}
}

func (s *sipProtocolTech) SetObserver(observer signaling.Observer) {
	s.observer = observer
}

func (s *sipProtocolTech) GetSession() signaling.Session {
	return s
}

func (s *sipProtocolTech) InitiateSession(ctx context.Context) error {
	request := signaling.GetContext[apisip.NewSessionRequest](ctx, signaling.SIPContextKey)
	if request == nil {
		slog.Error("unable to get SIP request from context")
		return fmt.Errorf("unable to get SIP request from context")
	}

	return s.cli.NewSIPSession(ctx, request)
}

func (s *sipProtocolTech) AcceptSession(ctx context.Context) error {
	return nil
}

func (s *sipProtocolTech) TerminateSession(ctx context.Context) error {
	return nil
}

func (s *sipProtocolTech) ActiveCalls() int {
	s.cli.cmu.Lock()
	activeClientCalls := len(s.cli.activeCalls)
	s.cli.cmu.Unlock()

	s.srv.cmu.Lock()
	activeServerCalls := len(s.srv.activeCalls)
	s.srv.cmu.Unlock()

	return activeClientCalls + activeServerCalls
}

func (s *sipProtocolTech) Stop() error {
	s.cli.Stop()
	s.srv.Stop()
	s.mon.Stop()

	return nil
}

func (s *sipProtocolTech) SetHandler(handler Handler) {
	s.srv.SetHandler(handler)
	s.cli.SetHandler(handler)
}

func (s *sipProtocolTech) Start() error {
	slog.Debug("starting sip service", "version", version.Version)
	for name, enabled := range s.conf.Codecs {
		if enabled {
			slog.Warn("codec enabled", nil, "name", name)
		} else {
			slog.Warn("codec disabled", nil, "name", name)
		}
	}
	media.CodecsSetEnabled(s.conf.Codecs)

	if err := s.mon.Start(s.conf); err != nil {
		return err
	}
	// The UA must be shared between the client and the server.
	// Otherwise, the client will have to listen on a random port, which must then be forwarded.
	//
	// Routers are smart, they usually keep the UDP "session" open for a few moments, and may allow INVITE handshake
	// to pass even without forwarding rules on the firewall. ut it will inevitably fail later on follow-up requests like BYE.
	ua, err := sipgo.NewUA(
		sipgo.WithUserAgent(UserAgent),
	)
	if err != nil {
		return err
	}
	if err := s.cli.Start(ua, s.sconf); err != nil {
		return err
	}
	// Server is responsible for answering all transactions. However, the client may also receive some (e.g. BYE).
	// Thus, all unhandled transactions will be checked by the client.
	if err := s.srv.Start(ua, s.sconf, s.cli.OnRequest); err != nil {
		return err
	}
	slog.Debug("sip service ready")
	return nil
}
