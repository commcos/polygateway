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
	"fmt"

	"github.com/commcos/component-base/cli"
	"github.com/commcos/component-base/cli/shell"
)

func EnterShell() error {
	cfg := shell.Config{
		Type: shell.ShellTypeLocal,
	}
	sl := shell.NewShell(cfg)
	cli.SetDefaultCli(sl)

	stop := make(chan struct{})
	if err := cli.StartInteractiveShell(stop); err != nil {
		return fmt.Errorf("start interactive shell failed: %w", err)
	}
	return fmt.Errorf("shell exited")
}
