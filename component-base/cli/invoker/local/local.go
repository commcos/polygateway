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

package local

import (
	"context"
	"log/slog"

	"github.com/commcos/component-base/cli"
	"github.com/commcos/component-base/cli/invoker"
)

type Config struct {
}

type localInvoker struct {
	config *Config
}

func NewInvoker() invoker.CommandInvoker {
	inv := &localInvoker{}

	return inv
}

func (inv *localInvoker) Invoke(ctx context.Context, inCmd *cli.Command) (result invoker.CMDInvokeResutl) {
	result.Suceess = false

	cmd := inCmd.Command
	slog.Info("exec cmd", "cmd", cmd.Use)

	if err := cmd.ExecuteContext(ctx); err != nil {
		slog.Error("execute command error",
			"cmd", cmd.Use,
			"error", err.Error())

		result.Rearson = err.Error()
		return
	}
	result.Suceess = true

	return
}
