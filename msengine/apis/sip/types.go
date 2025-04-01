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
	"time"
)

const (
	// maxCallDuration sets a global max call duration.
	MaxCallDuration = 24 * time.Hour
	// defaultRingingTimeout is a maximal duration which SIP participant will wait to connect.
	//
	// For inbound, the participant will wait this duration for other participant tracks.
	//
	// For outbound, this sets a timeout for the other end to pick up the call.
	DefaultRingingTimeout = 3 * time.Minute
)

// Names of participant attributes for SIP.
const (
	// AttrSIPPrefix is shared for all SIP attributes.
	AttrSIPPrefix = "sip."
	// AttrSIPCallID attribute contains LiveKit SIP call ID.
	AttrSIPCallID = AttrSIPPrefix + "callID"
	// AttrSIPTrunkID attribute contains LiveKit SIP Trunk ID used for the call.
	AttrSIPTrunkID = AttrSIPPrefix + "trunkID"
	// AttrSIPDispatchRuleID attribute contains LiveKit SIP DispatchRule ID used for the inbound call.
	AttrSIPDispatchRuleID = AttrSIPPrefix + "ruleID"
	// AttrSIPTrunkNumber attribute contains number associate with LiveKit SIP Trunk.
	// This attribute will be omitted if HidePhoneNumber is set.
	AttrSIPTrunkNumber = AttrSIPPrefix + "trunkPhoneNumber"
	// AttrSIPPhoneNumber attribute contains number external to LiveKit SIP (caller for inbound and called number for outbound).
	// This attribute will be omitted if HidePhoneNumber is set.
	AttrSIPPhoneNumber = AttrSIPPrefix + "phoneNumber"
	// AttrSIPHostName attribute contains host name external to LiveKit SIP (caller for inbound and called number for outbound).
	AttrSIPHostName = AttrSIPPrefix + "hostname"
	// AttrSIPCallStatus attribute contains current call status for a SIP call associated with the participant.
	//
	// SIP participant is ready when it reaches "active" status.
	AttrSIPCallStatus = AttrSIPPrefix + "callStatus"
	// AttrSIPHeaderPrefix is a prefix for automatically mapped SIP header attributes.
	AttrSIPHeaderPrefix = AttrSIPPrefix + "h."

	// AttrIngressPrefix is shared for all Ingress attributes
	AttrIngressPrefix = "ingress."
	// AttrIngressOutOfNetworkPrefix is shared for all ingress out of network (Ads break) related attributes
	AttrIngressOutOfNetworkPrefix = AttrIngressPrefix + "outOfNetwork."
	// AttrIngressOutOfNetworkEventID contains the event ID of the current Out of network splice
	AttrIngressOutOfNetworkEventID = AttrIngressOutOfNetworkPrefix + "eventID"
)

const (
	AttrSIPCallIDFull = AttrSIPPrefix + "callIDFull"
	AttrSIPCallTag    = AttrSIPPrefix + "callTag"
)

type SIPCallStatus int32

const (
	SIPCallStatus_SCS_CALL_INCOMING      SIPCallStatus = 0 // Incoming call is being handled by the SIP service. The SIP participant hasn't joined a LiveKit room yet
	SIPCallStatus_SCS_PARTICIPANT_JOINED SIPCallStatus = 1 // SIP participant for outgoing call has been created. The SIP outgoing call is being established
	SIPCallStatus_SCS_ACTIVE             SIPCallStatus = 2 // Call is ongoing. SIP participant is active in the LiveKit room
	SIPCallStatus_SCS_DISCONNECTED       SIPCallStatus = 3 // Call has ended
	SIPCallStatus_SCS_ERROR              SIPCallStatus = 4 // Call has ended or never succeeded because of an error
)

type SIPStatusCode int32

const (
	SIPStatusCode_SIP_STATUS_UNKNOWN                          SIPStatusCode = 0
	SIPStatusCode_SIP_STATUS_TRYING                           SIPStatusCode = 100
	SIPStatusCode_SIP_STATUS_RINGING                          SIPStatusCode = 180
	SIPStatusCode_SIP_STATUS_CALL_IS_FORWARDED                SIPStatusCode = 181
	SIPStatusCode_SIP_STATUS_QUEUED                           SIPStatusCode = 182
	SIPStatusCode_SIP_STATUS_SESSION_PROGRESS                 SIPStatusCode = 183
	SIPStatusCode_SIP_STATUS_OK                               SIPStatusCode = 200
	SIPStatusCode_SIP_STATUS_ACCEPTED                         SIPStatusCode = 202
	SIPStatusCode_SIP_STATUS_MOVED_PERMANENTLY                SIPStatusCode = 301
	SIPStatusCode_SIP_STATUS_MOVED_TEMPORARILY                SIPStatusCode = 302
	SIPStatusCode_SIP_STATUS_USE_PROXY                        SIPStatusCode = 305
	SIPStatusCode_SIP_STATUS_BAD_REQUEST                      SIPStatusCode = 400
	SIPStatusCode_SIP_STATUS_UNAUTHORIZED                     SIPStatusCode = 401
	SIPStatusCode_SIP_STATUS_PAYMENT_REQUIRED                 SIPStatusCode = 402
	SIPStatusCode_SIP_STATUS_FORBIDDEN                        SIPStatusCode = 403
	SIPStatusCode_SIP_STATUS_NOTFOUND                         SIPStatusCode = 404
	SIPStatusCode_SIP_STATUS_METHOD_NOT_ALLOWED               SIPStatusCode = 405
	SIPStatusCode_SIP_STATUS_NOT_ACCEPTABLE                   SIPStatusCode = 406
	SIPStatusCode_SIP_STATUS_PROXY_AUTH_REQUIRED              SIPStatusCode = 407
	SIPStatusCode_SIP_STATUS_REQUEST_TIMEOUT                  SIPStatusCode = 408
	SIPStatusCode_SIP_STATUS_CONFLICT                         SIPStatusCode = 409
	SIPStatusCode_SIP_STATUS_GONE                             SIPStatusCode = 410
	SIPStatusCode_SIP_STATUS_REQUEST_ENTITY_TOO_LARGE         SIPStatusCode = 413
	SIPStatusCode_SIP_STATUS_REQUEST_URI_TOO_LONG             SIPStatusCode = 414
	SIPStatusCode_SIP_STATUS_UNSUPPORTED_MEDIA_TYPE           SIPStatusCode = 415
	SIPStatusCode_SIP_STATUS_REQUESTED_RANGE_NOT_SATISFIABLE  SIPStatusCode = 416
	SIPStatusCode_SIP_STATUS_BAD_EXTENSION                    SIPStatusCode = 420
	SIPStatusCode_SIP_STATUS_EXTENSION_REQUIRED               SIPStatusCode = 421
	SIPStatusCode_SIP_STATUS_INTERVAL_TOO_BRIEF               SIPStatusCode = 423
	SIPStatusCode_SIP_STATUS_TEMPORARILY_UNAVAILABLE          SIPStatusCode = 480
	SIPStatusCode_SIP_STATUS_CALL_TRANSACTION_DOES_NOT_EXISTS SIPStatusCode = 481
	SIPStatusCode_SIP_STATUS_LOOP_DETECTED                    SIPStatusCode = 482
	SIPStatusCode_SIP_STATUS_TOO_MANY_HOPS                    SIPStatusCode = 483
	SIPStatusCode_SIP_STATUS_ADDRESS_INCOMPLETE               SIPStatusCode = 484
	SIPStatusCode_SIP_STATUS_AMBIGUOUS                        SIPStatusCode = 485
	SIPStatusCode_SIP_STATUS_BUSY_HERE                        SIPStatusCode = 486
	SIPStatusCode_SIP_STATUS_REQUEST_TERMINATED               SIPStatusCode = 487
	SIPStatusCode_SIP_STATUS_NOT_ACCEPTABLE_HERE              SIPStatusCode = 488
	SIPStatusCode_SIP_STATUS_INTERNAL_SERVER_ERROR            SIPStatusCode = 500
	SIPStatusCode_SIP_STATUS_NOT_IMPLEMENTED                  SIPStatusCode = 501
	SIPStatusCode_SIP_STATUS_BAD_GATEWAY                      SIPStatusCode = 502
	SIPStatusCode_SIP_STATUS_SERVICE_UNAVAILABLE              SIPStatusCode = 503
	SIPStatusCode_SIP_STATUS_GATEWAY_TIMEOUT                  SIPStatusCode = 504
	SIPStatusCode_SIP_STATUS_VERSION_NOT_SUPPORTED            SIPStatusCode = 505
	SIPStatusCode_SIP_STATUS_MESSAGE_TOO_LARGE                SIPStatusCode = 513
	SIPStatusCode_SIP_STATUS_GLOBAL_BUSY_EVERYWHERE           SIPStatusCode = 600
	SIPStatusCode_SIP_STATUS_GLOBAL_DECLINE                   SIPStatusCode = 603
	SIPStatusCode_SIP_STATUS_GLOBAL_DOES_NOT_EXIST_ANYWHERE   SIPStatusCode = 604
	SIPStatusCode_SIP_STATUS_GLOBAL_NOT_ACCEPTABLE            SIPStatusCode = 606
)

// SIPStatus is returned as an error detail in CreateSIPParticipant.
type SIPStatus struct {
	Code   SIPStatusCode
	Status string
}

type DisconnectReason int32

const (
	DisconnectReason_UNKNOWN_REASON DisconnectReason = 0
	// the client initiated the disconnect
	DisconnectReason_CLIENT_INITIATED DisconnectReason = 1
	// another participant with the same identity has joined the room
	DisconnectReason_DUPLICATE_IDENTITY DisconnectReason = 2
	// the server instance is shutting down
	DisconnectReason_SERVER_SHUTDOWN DisconnectReason = 3
	// RoomService.RemoveParticipant was called
	DisconnectReason_PARTICIPANT_REMOVED DisconnectReason = 4
	// RoomService.DeleteRoom was called
	DisconnectReason_ROOM_DELETED DisconnectReason = 5
	// the client is attempting to resume a session, but server is not aware of it
	DisconnectReason_STATE_MISMATCH DisconnectReason = 6
	// client was unable to connect fully
	DisconnectReason_JOIN_FAILURE DisconnectReason = 7
	// Cloud-only, the server requested Participant to migrate the connection elsewhere
	DisconnectReason_MIGRATION DisconnectReason = 8
	// the signal websocket was closed unexpectedly
	DisconnectReason_SIGNAL_CLOSE DisconnectReason = 9
	// the room was closed, due to all Standard and Ingress participants having left
	DisconnectReason_ROOM_CLOSED DisconnectReason = 10
	// SIP callee did not respond in time
	DisconnectReason_USER_UNAVAILABLE DisconnectReason = 11
	// SIP callee rejected the call (busy)
	DisconnectReason_USER_REJECTED DisconnectReason = 12
	// SIP protocol failure or unexpected response
	DisconnectReason_SIP_TRUNK_FAILURE DisconnectReason = 13
)

type SIPCallDirection int32

const (
	SIPCallDirection_SCD_UNKNOWN  SIPCallDirection = 0
	SIPCallDirection_SCD_INBOUND  SIPCallDirection = 1
	SIPCallDirection_SCD_OUTBOUND SIPCallDirection = 2
)

type SIPTransport int32

const (
	SIPTransport_SIP_TRANSPORT_AUTO SIPTransport = 0
	SIPTransport_SIP_TRANSPORT_UDP  SIPTransport = 1
	SIPTransport_SIP_TRANSPORT_TCP  SIPTransport = 2
	SIPTransport_SIP_TRANSPORT_TLS  SIPTransport = 3
)

type SIPHeaderOptions int32

const (
	SIPHeaderOptions_SIP_NO_HEADERS  SIPHeaderOptions = 0 // do not map any headers, except ones mapped explicitly
	SIPHeaderOptions_SIP_X_HEADERS   SIPHeaderOptions = 1 // map all X-* headers to sip.h.x-* attributes
	SIPHeaderOptions_SIP_ALL_HEADERS SIPHeaderOptions = 2 // map all headers to sip.h.* attributes
)

type SIPFeature int32

const (
	SIPFeature_NONE          SIPFeature = 0
	SIPFeature_KRISP_ENABLED SIPFeature = 1
)

type SIPMediaEncryption int32

const (
	SIPMediaEncryption_SIP_MEDIA_ENCRYPT_DISABLE SIPMediaEncryption = 0 // do not enable encryption
	SIPMediaEncryption_SIP_MEDIA_ENCRYPT_ALLOW   SIPMediaEncryption = 1 // use encryption if available
	SIPMediaEncryption_SIP_MEDIA_ENCRYPT_REQUIRE SIPMediaEncryption = 2 // require encryption
)

type CallStatus int

func (v CallStatus) Attribute() string {
	switch v {
	default:
		return "" // no attribute for these statuses
	case CallDialing:
		return "dialing"
	case CallRinging:
		return "ringing"
	case CallAutomation:
		return "automation"
	case CallActive:
		return "active"
	case CallHangup:
		return "hangup"
	}
}

func (v CallStatus) DisconnectReason() DisconnectReason {
	switch v {
	default:
		return DisconnectReason_UNKNOWN_REASON
	case CallHangup:
		// It's the default that LK sets, but map it here explicitly to show the assumption.
		return DisconnectReason_CLIENT_INITIATED
	case CallUnavailable:
		return DisconnectReason_USER_UNAVAILABLE
	case CallRejected:
		return DisconnectReason_USER_REJECTED
	}
}

const (
	CallDropped = CallStatus(iota)
	CallFlood
	CallDialing
	CallRinging
	CallAutomation
	CallActive
	CallHangup
	CallUnavailable
	CallRejected
)

type SIPUri struct {
	User      string
	Host      string
	Ip        string
	Port      uint32
	Transport SIPTransport
}

type SIPCallInfo struct {
	CallId           string
	Region           string
	FromUri          *SIPUri
	ToUri            *SIPUri
	CallDirection    SIPCallDirection
	CallStatus       SIPCallStatus
	CreatedAtNs      int64
	StartedAtNs      int64
	EndedAtNs        int64
	DisconnectReason DisconnectReason
	Error            string
	CallStatusCode   *SIPStatus
	AudioCodec       string
	MediaEncryption  string
}

type SipDTMF struct {
	Code  uint32
	Digit string
}

type NewSessionRequest struct {
	CallID  string
	Address string
	// Hostname for the 'From' SIP address in INVITE
	Hostname string
	// Number used to make the call
	FromNumber string
	// Number to call to
	ToNumber  string
	Transport SIPTransport
	Username  string
	Password  string
	// Optionally send following DTMF digits (extension codes) when making a call.
	// Character 'w' can be used to add a 0.5 sec delay.
	Dtmf string
	// Optionally play dialtone in the room as an audible indicator for existing participants
	PlayDialtone        bool
	Headers             map[string]string
	HeadersToAttributes map[string]string
	// Map LiveKit attributes to SIP X-* headers when sending BYE or REFER requests.
	// Keys are the names of attributes and values are the names of X-* headers they will be mapped to.
	AttributesToHeaders map[string]string
	// Map SIP headers from 200 OK to sip.h.* participant attributes automatically.
	//
	// When the names of required headers is known, using headers_to_attributes is strongly recommended.
	//
	// When mapping 200 OK headers to follow-up request headers with attributes_to_headers map,
	// lowercase header names should be used, for example: sip.h.x-custom-header.
	IncludeHeaders  SIPHeaderOptions
	EnabledFeatures []SIPFeature
	// Max time for the callee to answer the call.
	RingingTimeout *time.Duration
	// Max call duration.
	MaxCallDuration *time.Duration
	MediaEncryption SIPMediaEncryption
	// Wait for the answer for the call before returning.
	WaitUntilAnswered bool
}

type SIPCall struct {
	LkCallId string
	SourceIp string    // source ip (without port)
	Address  *SIPUri   // address in the request line (INVITE)
	From     *SIPUri   // From header
	To       *SIPUri   // To header
	Via      []*SIPUri // Via headers
}
