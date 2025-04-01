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

package msengine

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/commcos/msengine/signaling"
	"github.com/commcos/msengine/signaling/sip"
	sipconfig "github.com/commcos/msengine/signaling/sip/config"
	"github.com/commcos/msengine/signaling/sip/stats"
)

type CoreEngine struct {
	sipTech signaling.ProtocolTech
}

func NewCoreEngine() *CoreEngine {
	ce := &CoreEngine{}

	if err := ce.newSIPTech(); err != nil {
		return nil
	}

	return ce
}

func (ce *CoreEngine) newSIPTech() error {
	conf := &sipconfig.Config{}

	mon, err := stats.NewMonitor(conf)
	if err != nil {
		slog.Error("failed to create monitor", "error", err)
		return err
	}

	ce.sipTech, err = sip.NewSIPTech("", conf, mon)
	if err != nil {
		slog.Error("failed to create SIP tech", "error", err)
		return err
	}

	return nil
}

func (ce *CoreEngine) Start() error {
	if err := ce.sipTech.Start(); err != nil {
		return err
	}

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGTERM, syscall.SIGQUIT)

	killChan := make(chan os.Signal, 1)
	signal.Notify(killChan, syscall.SIGINT)

	go func() {
		select {
		case sig := <-stopChan:
			slog.Info("exit requested, finishing all SIP then shutting down", "signal", sig)
			ce.Stop()
		case sig := <-killChan:
			slog.Info("exit requested, stopping all SIP and shutting down", "signal", sig)
			ce.Stop()
		}
	}()

	return nil
}

func (ce *CoreEngine) Stop() error {
	return ce.shutdown()
}

func (ce *CoreEngine) shutdown() error {
	slog.Info("Shutting down the core engine")
	ce.sipTech.Stop()

	return nil
}

func (ce *CoreEngine) session(typ signaling.ProtocolType) signaling.Session {
	var session signaling.Session
	switch typ {
	case signaling.ProtocolTypeSIP:
		session = ce.sipTech.GetSession()
	}
	return session
}

func (ce *CoreEngine) InitiateSession(ctx context.Context, typ signaling.ProtocolType) error {

	session := ce.session(typ)
	if session == nil {
		slog.Error("session is nil")
		return fmt.Errorf("session not found for type %d", typ)
	}

	if err := session.InitiateSession(ctx); err != nil {
		slog.Error("failed to initiate SIP session", "error", err)
		return err
	}

	return nil
}

func (ce *CoreEngine) AcceptSession(ctx context.Context, typ signaling.ProtocolType) error {
	session := ce.session(typ)
	if session == nil {
		slog.Error("session is nil")
		return fmt.Errorf("session not found for type %d", typ)
	}
	if err := session.AcceptSession(ctx); err != nil {
		slog.Error("failed to accept SIP session", "error", err)
		return err
	}

	return nil
}

func (ce *CoreEngine) TerminateSession(ctx context.Context, typ signaling.ProtocolType) error {
	session := ce.session(typ)
	if session == nil {
		slog.Error("session is nil")
		return fmt.Errorf("session not found for type %d", typ)
	}
	if err := session.TerminateSession(ctx); err != nil {
		slog.Error("failed to terminate SIP session", "error", err)
		return err
	}

	return nil
}
