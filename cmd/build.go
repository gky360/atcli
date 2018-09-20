// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/gky360/atcli/utils"
	"github.com/spf13/cobra"
)

type BuildOptions struct {
	Out, ErrOut io.Writer
}

var buildOpt = &BuildOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// buildCmd represents the build command
var buildCmd = &cobra.Command{
	Use:   "build [task name]",
	Short: "Build your source code for a task",
	Long: `Build your source code for a task.

Example:
    # Build your source code for task 'D'.
    atcli build d

"atcli" supports following languages and compiler commands.

- C++14 (GCC 5.4.1)
  "g++ -std=gnu++1y -O2 -o a.out Main.cpp"

You must install the compiler for your language listed above
and add the compiler to PATH in advance.

Your source code file is supposed to be in
"$ATCLI_ROOT/arc090/d/Main.cpp" , for example.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := buildOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(buildOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(buildCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// buildCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// buildCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *BuildOptions) Run(cmd *cobra.Command, args []string) (err error) {
	taskName := args[0]
	if err = runBuild(taskName, true, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	return nil
}

func runBuild(taskName string, isForce bool, out, errOut io.Writer) error {
	taskDir, err := utils.TaskDir(taskName)
	if err != nil {
		return err
	}

	if err := os.Chdir(taskDir); err != nil {
		return err
	}
	if _, err := os.Stat("a.out"); err == nil {
		// file exists
		if !isForce {
			return nil
		}
	} else if !os.IsNotExist(err) {
		// error
		return err
	}

	execCmd := exec.Command(
		"g++", "-std=gnu++1y", "-O2", "-o", "a.out", "Main.cpp",
	)
	execCmd.Stdout = out
	execCmd.Stderr = errOut
	execCmd.Run()
	if err != nil {
		return err
	}

	fmt.Fprintln(errOut, "Build succeeded.")
	return nil
}
