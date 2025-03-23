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

package cli

import (
	"github.com/spf13/cobra"
)

type Command struct {
	*cobra.Command
}

// Interface cli interface
type Interface interface {
	RegisterCommand(cmd Command)
	// start the cli with interactive shell. stop is a channel to stop the shell
	StartInteractiveShell(stop chan struct{}) error
}

var defaultCli Interface

func SetDefaultCli(cli Interface) {
	defaultCli = cli
}

func AddCommand(cmd Command) {
	defaultCli.RegisterCommand(cmd)
}

func StartInteractiveShell(stop chan struct{}) error {
	return defaultCli.StartInteractiveShell(stop)
}
