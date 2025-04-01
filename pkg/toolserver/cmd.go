/********************************************************************
 * Copyright (c) 2025. All Rights Reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 *********************************************************************/

package toolserver

import (
	"context"
	"strings"

	"github.com/commcos/component-base/cli"
	"github.com/commcos/component-base/uuid"
	"github.com/commcos/msengine/signaling"
	"github.com/spf13/cobra"
	"k8s.io/klog/v2"

	apisip "github.com/commcos/msengine/apis/sip"
)

func (ts *ToolServer) installCmd() {
	ts.cmdOriginCall()
}

func (ts *ToolServer) cmdOriginCall() {
	originalCall := cobra.Command{
		Use:   "origincall",
		Short: "origincall start a new call to destination",
		Long: `origincall start a new call to destination. 
		This command is used to initiate a call to a specified destination from the internal.`,
		Args: cobra.MatchAll(cobra.ExactArgs(4), cobra.OnlyValidArgs),
		Run: func(cmd *cobra.Command, args []string) {
			// This is a placeholder for the original call logic
			// You can replace this with the actual implementation
			klog.V(2).Infof("origincall called with args: %v", strings.Join(args, " "))

			ctx := context.Background()
			sessReq := &apisip.NewSessionRequest{
				CallID:     uuid.NewUUID().String(),
				Address:    args[0],
				FromNumber: args[1],
				ToNumber:   args[2],
			}
			ctx = signaling.WithContext(ctx, signaling.SIPContextKey, sessReq)
			ts.msEngine.InitiateSession(ctx, signaling.ProtocolTypeSIP)
		},
	}
	ts.cliHandle.RegisterCommand(cli.Command{
		Command: &originalCall})

}
