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

package shell

import (
	"context"
	"fmt"
	"hash/fnv"
	"log/slog"
	"os"
	"strings"
	"text/template"
	"time"

	"github.com/chzyer/readline"
	"github.com/commcos/component-base/cli"
	"github.com/commcos/component-base/cli/invoker"
	localinvoke "github.com/commcos/component-base/cli/invoker/local"
	"github.com/spf13/cobra"
)

type ShellType int

const (
	ShellTypeLocal ShellType = 0
)

const (
	cmdTemplate = `
	Usage: {{.UseLine}}
	
	{{if .HasAvailableSubCommands}}Commands:
	{{range .Commands}}{{if .IsAvailableCommand}}
		{{printf "%-15s" .Name}}  {{.Short}}
	{{end}}{{end}}{{end}}
	{{if .HasExample}}Examples:
	{{.Example}}{{end}}
	`

	helpTemplate = `
	{{range $cmd, $desc := .}}
		{{printf "%-20s" $cmd}} - {{$desc}}
	{{end}}
	`

	resultTemplate = `
	Execution Results:
	{{range $key, $value := .}}
		{{printf "%-20s" $key}} : {{$value}}
	{{end}}
	`
)

type Config struct {
	Type ShellType
}

type shell struct {
	config     *Config
	cmdSet     map[uint32]*cobra.Command
	defaultCmd *cobra.Command

	//command as key and desc as value
	helperTmpl *template.Template
	helper     map[string]string

	cliResultTmpl *template.Template
	rlCompleter   *readline.PrefixCompleter

	invoker invoker.CommandInvoker
}

// NewInteractiveCli local interactive
func NewShell(config Config) cli.Interface {

	sl := &shell{
		helper:      make(map[string]string),
		cmdSet:      make(map[uint32]*cobra.Command),
		rlCompleter: readline.NewPrefixCompleter(),
		config:      &config,
	}

	sl.defaultCmd = &cobra.Command{
		Short:              "cli",
		Long:               "Launch a interactive cli",
		DisableFlagParsing: true,
	}

	exitCmd := &cobra.Command{
		Use:   "exit",
		Short: "exit shell",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("############ GoodBye!!! ###########")
			os.Exit(0)
		},
	}

	sl.RegisterCommand(cli.Command{Command: exitCmd})

	sl.helperTmpl = template.Must(template.New("helper").Parse(helpTemplate))
	sl.defaultCmd.SetHelpFunc(func(cmd *cobra.Command, input []string) {
		if err := sl.helperTmpl.Execute(os.Stdout, sl.helper); err != nil {
			slog.Error("cli help", "error", err)
		}
	})
	sl.defaultCmd.SetHelpTemplate(cmdTemplate)

	sl.cliResultTmpl = template.Must(template.New("shellrst").Parse(resultTemplate))

	sl.invoker = localinvoke.NewInvoker()

	return sl
}

func (sl *shell) RegisterCommand(cmd cli.Command) {
	newCmd := cmd.Command
	cmd.Command.DisableFlagParsing = true
	inputHashKey := strings.ReplaceAll(newCmd.Use, " ", "")

	key := HashKey(inputHashKey)
	sl.cmdSet[key] = newCmd
	sl.pcFromCommands(sl.rlCompleter, newCmd)
	sl.helper[newCmd.Use] = newCmd.Short

	slog.Info("register command", "cmd", newCmd.Use)
}

func (sl *shell) checkUserInput(args []string) bool {
	if len(args) == 0 {
		return false
	}

	if len(args) > 0 {
		switch args[0] {
		case "exit":
			return false
		}
	}

	return true
}

func (sl *shell) pcFromCommands(parent readline.PrefixCompleterInterface, c *cobra.Command) {
	pc := readline.PcItem(c.Use)
	parent.SetChildren(append(parent.GetChildren(), pc))
}

func (sl *shell) findcmd(args []string) (*cobra.Command, error) {
	var cmdName string

	for _, item := range args {
		cmdName += item

		key := HashKey(strings.TrimSpace(cmdName))
		cmd, ok := sl.cmdSet[key]
		slog.Info("find command", "cmd", cmdName, "key", key, "ok", ok)

		if ok {
			return cmd, nil
		}
	}

	return nil, fmt.Errorf("command not found")
}

// DoExcute excute command
func (sl *shell) DoExcute(info *InputLine) error {

	tmpCtx, cancel := context.WithTimeout(context.Background(), time.Second*30)
	defer cancel()

	var err error
	var cmd *cobra.Command
	cmd, err = sl.findcmd(info.Argument)
	if err != nil {
		return err
	}
	cmd.SetArgs(info.Argument)

	sl.invoker.Invoke(tmpCtx, &cli.Command{Command: cmd})

	return nil
}

func (sl *shell) StartInteractiveShell(stop chan struct{}) error {

	invoker := func(args []string) {
		if !sl.checkUserInput(args) {
			return
		}

		info := &InputLine{
			Argument: append([]string{}, args...),
		}
		sl.DoExcute(info)
	}

	shell, err := readline.NewEx(&readline.Config{
		Prompt:       "cmd> ",
		AutoComplete: sl.rlCompleter,
		EOFPrompt:    "exit",
	})
	if err != nil {
		panic(err)
	}
	defer shell.Close()

	for {
		select {
		case <-stop:
			return nil
		default:
			l, err := shell.Readline()
			if err != nil {

				if err == readline.ErrInterrupt {
					os.Exit(0)
				}

				fmt.Println("Find comand error=", err)
				continue
			}

			args := strings.Fields(l)
			invoker(args)
		}
	}
}

func HashKey(key string) uint32 {
	h := fnv.New32a()
	if _, err := h.Write([]byte(key)); err != nil {
		return 0
	}

	return h.Sum32()
}
