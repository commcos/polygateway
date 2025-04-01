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

package config

import (
	"fmt"
	"log/slog"
	"net"
	"net/netip"
	"os"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/commcos/component-base/uuid"

	"github.com/commcos/msengine/apis"
	siperrors "github.com/commcos/msengine/signaling/sip/errors"
)

const (
	DefaultSIPPort    int = 5060
	DefaultSIPPortTLS int = 5061
)

var (
	DefaultRTPPortRange = apis.PortRange{Start: 10000, End: 20000}
)

type TLSCert struct {
	CertFile string `yaml:"cert_file"`
	KeyFile  string `yaml:"key_file"`
}

type TLSConfig struct {
	Port       int       `yaml:"port"`        // announced SIP signaling port
	ListenPort int       `yaml:"port_listen"` // SIP signaling port to listen on
	Certs      []TLSCert `yaml:"certs"`
}

type Config struct {
	ApiKey    string `yaml:"api_key"`    // required (env LIVEKIT_API_KEY)
	ApiSecret string `yaml:"api_secret"` // required (env LIVEKIT_API_SECRET)
	WsUrl     string `yaml:"ws_url"`     // required (env LIVEKIT_WS_URL)

	HealthPort        int            `yaml:"health_port"`
	PrometheusPort    int            `yaml:"prometheus_port"`
	PProfPort         int            `yaml:"pprof_port"`
	SIPPort           int            `yaml:"sip_port"`        // announced SIP signaling port
	SIPPortListen     int            `yaml:"sip_port_listen"` // SIP signaling port to listen on
	SIPHostname       string         `yaml:"sip_hostname"`
	TLS               *TLSConfig     `yaml:"tls"`
	RTPPort           apis.PortRange `yaml:"rtp_port"`
	ClusterID         string         `yaml:"cluster_id"` // cluster this instance belongs to
	MaxCpuUtilization float64        `yaml:"max_cpu_utilization"`

	UseExternalIP bool   `yaml:"use_external_ip"`
	LocalNet      string `yaml:"local_net"` // local IP net to use, e.g. 192.168.0.0/24
	NAT1To1IP     string `yaml:"nat_1_to_1_ip"`
	ListenIP      string `yaml:"listen_ip"`

	MediaTimeout        time.Duration   `yaml:"media_timeout"`
	MediaTimeoutInitial time.Duration   `yaml:"media_timeout_initial"`
	Codecs              map[string]bool `yaml:"codecs"`

	// HideInboundPort controls how SIP endpoint responds to unverified inbound requests.
	// Setting it to true makes SIP server silently drop INVITE requests if it gets a negative Auth or Dispatch response.
	// Doing so hides our SIP endpoint from (a low effort) port scanners.
	HideInboundPort bool `yaml:"hide_inbound_port"`

	// AudioDTMF forces SIP to generate audio DTMF tones in addition to digital.
	AudioDTMF          bool `yaml:"audio_dtmf"`
	EnableJitterBuffer bool `yaml:"enable_jitter_buffer"`

	// internal
	ServiceName string `yaml:"-"`
	NodeID      string // Do not provide, will be overwritten
}

func NewConfig(confString string) (*Config, error) {
	conf := &Config{
		ApiKey:      os.Getenv("LIVEKIT_API_KEY"),
		ApiSecret:   os.Getenv("LIVEKIT_API_SECRET"),
		WsUrl:       os.Getenv("LIVEKIT_WS_URL"),
		ServiceName: "sip",
	}
	if confString != "" {
		if err := yaml.Unmarshal([]byte(confString), conf); err != nil {
			return nil, siperrors.ErrCouldNotParseConfig(err)
		}
	}

	// if conf.Redis == nil {
	// 	return nil, siperrors.NewErrorf(siperrors.InvalidArgument, "redis configuration is required")
	// }

	return conf, nil
}

func (c *Config) Init() error {
	c.NodeID = fmt.Sprintf("NE_%s", uuid.NewUUID().String())

	if c.SIPPort == 0 {
		c.SIPPort = DefaultSIPPort
	}
	if c.SIPPortListen == 0 {
		c.SIPPortListen = c.SIPPort
	}
	if tc := c.TLS; tc != nil {
		if tc.Port == 0 {
			tc.Port = DefaultSIPPortTLS
		}
		if tc.ListenPort == 0 {
			tc.ListenPort = tc.Port
		}
	}
	if c.RTPPort.Start == 0 {
		c.RTPPort.Start = DefaultRTPPortRange.Start
	}
	if c.RTPPort.End == 0 {
		c.RTPPort.End = DefaultRTPPortRange.End
	}
	if c.MaxCpuUtilization <= 0 || c.MaxCpuUtilization > 1 {
		c.MaxCpuUtilization = 0.9
	}

	if err := c.InitLogger(); err != nil {
		return err
	}

	if c.UseExternalIP && c.NAT1To1IP != "" {
		return fmt.Errorf("use_external_ip and nat_1_to_1_ip can not both be set")
	}

	return nil
}

func (c *Config) InitLogger(values ...interface{}) error {
	// zl, err := logger.NewZapLogger(&c.Logging)
	// if err != nil {
	// 	return err
	// }

	// values = append(c.GetLoggerValues(), values...)
	// l := zl.WithValues(values...)
	// logger.SetLogger(l, c.ServiceName)

	return nil
}

// To use with zap logger
func (c *Config) GetLoggerValues() []interface{} {
	if c.NodeID == "" {
		return nil
	}
	return []interface{}{"nodeID", c.NodeID}
}

func GetLocalIP() (netip.Addr, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return netip.Addr{}, nil
	}
	type Iface struct {
		Name string
		Addr netip.Addr
	}
	var candidates []Iface
	for _, ifc := range ifaces {
		if ifc.Flags&net.FlagUp == 0 || ifc.Flags&net.FlagRunning == 0 {
			continue
		}
		if ifc.Flags&(net.FlagPointToPoint|net.FlagLoopback) != 0 {
			continue
		}
		addrs, err := ifc.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			ipnet, ok := addr.(*net.IPNet)
			if !ok {
				continue
			}
			if ip4 := ipnet.IP.To4(); ip4 != nil {
				ip, _ := netip.AddrFromSlice(ip4)
				candidates = append(candidates, Iface{
					Name: ifc.Name, Addr: ip,
				})
				slog.Debug("considering interface", "iface", ifc.Name, "ip", ip)
			}
		}
	}
	if len(candidates) == 0 {
		return netip.Addr{}, fmt.Errorf("No local IP found")
	}
	return candidates[0].Addr, nil
}
