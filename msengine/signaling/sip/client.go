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

package sip

import (
	"context"
	"log/slog"
	"net/netip"
	"strings"
	"sync"
	"time"

	"github.com/frostbyte73/core"
	"golang.org/x/exp/maps"

	"github.com/livekit/sipgo"
	"github.com/livekit/sipgo/sip"

	apisip "github.com/commcos/msengine/apis/sip"
	"github.com/commcos/msengine/signaling/sip/config"
	siperrors "github.com/commcos/msengine/signaling/sip/errors"
	"github.com/commcos/msengine/signaling/sip/stats"
)

type Client struct {
	conf   *config.Config
	sconf  *ServiceConfig
	region string
	mon    *stats.Monitor

	sipCli *sipgo.Client

	closing     core.Fuse
	cmu         sync.Mutex
	activeCalls map[LocalTag]*outboundCall
	byRemote    map[RemoteTag]*outboundCall

	handler Handler
}

func NewClient(region string, conf *config.Config, mon *stats.Monitor) *Client {

	c := &Client{
		conf:        conf,
		region:      region,
		mon:         mon,
		activeCalls: make(map[LocalTag]*outboundCall),
		byRemote:    make(map[RemoteTag]*outboundCall),
	}
	return c
}

func (c *Client) Start(agent *sipgo.UserAgent, sc *ServiceConfig) error {
	c.sconf = sc
	slog.Info("client starting", "local", c.sconf.SignalingIPLocal, "external", c.sconf.SignalingIP)

	if agent == nil {
		ua, err := sipgo.NewUA(
			sipgo.WithUserAgent(UserAgent),
		)
		if err != nil {
			return err
		}
		agent = ua
	}

	var err error
	c.sipCli, err = sipgo.NewClient(agent,
		sipgo.WithClientHostname(c.sconf.SignalingIP.String()),
		sipgo.WithClientLogger(slog.Default()),
	)
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Stop() {
	c.closing.Break()
	c.cmu.Lock()
	calls := maps.Values(c.activeCalls)
	c.activeCalls = make(map[LocalTag]*outboundCall)
	c.byRemote = make(map[RemoteTag]*outboundCall)
	c.cmu.Unlock()
	for _, call := range calls {
		call.Close()
	}
	if c.sipCli != nil {
		c.sipCli.Close()
		c.sipCli = nil
	}
}

func (c *Client) SetHandler(handler Handler) {
	c.handler = handler
}

func (c *Client) ContactURI(tr Transport) URI {
	return getContactURI(c.conf, c.sconf.SignalingIP, tr)
}

func (c *Client) NewSIPSession(ctx context.Context, req *apisip.NewSessionRequest) (retErr error) {
	if !c.mon.CanAccept() {
		return siperrors.ErrUnavailable
	}
	if req.ToNumber == "" {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "call-to number must be set")
	} else if req.Address == "" {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "trunk adresss must be set")
	} else if req.FromNumber == "" {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "trunk outbound number must be set")
	}

	if strings.Contains(req.ToNumber, "@") {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "call_to should be a phone number or SIP user, not a full SIP URI")
	}
	if strings.HasPrefix(req.Address, "sip:") || strings.HasPrefix(req.Address, "sips:") {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "address must be a hostname without 'sip:' prefix")
	}
	if strings.Contains(req.Address, "transport=") {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "address must not contain parameters; use transport field")
	}
	if strings.ContainsAny(req.Address, ";=") {
		return siperrors.NewErrorf(siperrors.InvalidArgument, "address must not contain parameters")
	}

	slog := slog.With(
		"callID", req.CallID,
		"fromHost", req.Hostname,
		"fromUser", req.FromNumber,
		"toHost", req.Address,
		"toUser", req.ToNumber,
	)

	state := NewCallState(c.createSIPCallInfo(req))

	defer func() {
		state.Update(ctx, func(info *apisip.SIPCallInfo) {

			switch retErr {
			case nil:
				info.CallStatus = apisip.SIPCallStatus_SCS_PARTICIPANT_JOINED
			default:
				info.CallStatus = apisip.SIPCallStatus_SCS_ERROR
				info.DisconnectReason = apisip.DisconnectReason_UNKNOWN_REASON
				info.Error = retErr.Error()
			}
		})
	}()

	sipConf := sipOutboundConfig{
		address:         req.Address,
		transport:       req.Transport,
		host:            req.Hostname,
		from:            req.FromNumber,
		to:              req.ToNumber,
		user:            req.Username,
		pass:            req.Password,
		dtmf:            req.Dtmf,
		dialtone:        req.PlayDialtone,
		headers:         req.Headers,
		includeHeaders:  req.IncludeHeaders,
		headersToAttrs:  req.HeadersToAttributes,
		attrsToHeaders:  req.AttributesToHeaders,
		enabledFeatures: req.EnabledFeatures,
	}
	if req.RingingTimeout != nil {
		sipConf.ringingTimeout = *req.RingingTimeout
	}
	if req.MaxCallDuration != nil {
		sipConf.maxCallDuration = *req.MaxCallDuration
	}

	slog.Info("Creating SIP participant")
	call, err := c.newCall(ctx, c.conf, *slog, LocalTag(req.CallID), sipConf, state)
	if err != nil {
		return err
	}

	if !req.WaitUntilAnswered {
		call.DialAsync(ctx)
		return nil
	}
	if err := call.Dial(ctx); err != nil {
		return err
	}
	go call.WaitClose(context.WithoutCancel(ctx))
	return nil
}

func (c *Client) createSIPCallInfo(req *apisip.NewSessionRequest) *apisip.SIPCallInfo {
	toUri := CreateURIFromUserAndAddress(req.ToNumber, req.Address, TransportFrom(req.Transport))
	fromiUri := URI{
		User: req.FromNumber,
		Host: req.Hostname,
		Addr: netip.AddrPortFrom(c.sconf.SignalingIP, uint16(c.conf.SIPPort)),
	}

	callInfo := &apisip.SIPCallInfo{
		CallId:        req.CallID,
		Region:        c.region,
		CallDirection: apisip.SIPCallDirection_SCD_OUTBOUND,
		ToUri:         toUri.ToSIPUri(),
		FromUri:       fromiUri.ToSIPUri(),
		CreatedAtNs:   time.Now().UnixNano(),
	}

	return callInfo
}

func (c *Client) OnRequest(req *sip.Request, tx sip.ServerTransaction) bool {
	switch req.Method {
	default:
		return false
	case "BYE":
		return c.onBye(req, tx)
	case "NOTIFY":
		return c.onNotify(req, tx)
	}
}

func (c *Client) onBye(req *sip.Request, tx sip.ServerTransaction) bool {
	tag, _ := getFromTag(req)
	c.cmu.Lock()
	call := c.byRemote[tag]
	c.cmu.Unlock()
	if call == nil {
		if tag != "" {
			slog.Info("BYE for non-existent call", "sipTag", tag)
		}
		_ = tx.Respond(sip.NewResponseFromRequest(req, sip.StatusCallTransactionDoesNotExists, "Call does not exist", nil))
		return false
	}
	call.log.Info("BYE")
	go func(call *outboundCall) {
		call.cc.AcceptBye(req, tx)
		call.CloseWithReason(apisip.CallHangup, "bye", apisip.DisconnectReason_CLIENT_INITIATED)
	}(call)
	return true
}

func (c *Client) onNotify(req *sip.Request, tx sip.ServerTransaction) bool {
	tag, _ := getFromTag(req)
	c.cmu.Lock()
	call := c.byRemote[tag]
	c.cmu.Unlock()
	if call == nil {
		return false
	}
	call.log.Info("NOTIFY")
	go func() {
		err := call.cc.handleNotify(req, tx)

		code, msg := sipCodeAndMessageFromError(err)

		tx.Respond(sip.NewResponseFromRequest(req, code, msg, nil))
	}()
	return true
}

func (c *Client) RegisterTransferSIPParticipant(sipCallID string, o *outboundCall) error {
	return c.handler.RegisterTransferSIPParticipantTopic(sipCallID)
}

func (c *Client) DeregisterTransferSIPParticipant(sipCallID string) {
	c.handler.DeregisterTransferSIPParticipantTopic(sipCallID)
}
