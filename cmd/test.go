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
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"

	"github.com/gky360/atcli/utils"
	"github.com/spf13/cobra"
)

type TestOptions struct {
	Out, ErrOut io.Writer

	isForce bool
}

var testOpt = &TestOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := testOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(testOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVarP(&testOpt.isForce, "force", "f", false, "Force to build.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// testCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// testCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *TestOptions) Run(cmd *cobra.Command, args []string) (err error) {
	taskName := args[0]
	sampleNum := -1
	if len(args) >= 2 {
		if sampleNum, err = strconv.Atoi(args[1]); err != nil {
			return err
		}
	}

	if err = runBuild(taskName, opt.isForce, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	if sampleNum >= 0 {
		if err = runWithSample(taskName, sampleNum, opt.Out, opt.ErrOut); err != nil {
			return err
		}
	} else {
		if err = runWithSamples(taskName, opt.Out, opt.ErrOut); err != nil {
			return err
		}
	}

	return nil
}

var errSampleInputNotExist = errors.New("Sample input not found.")

func runWithSample(taskName string, sampleNum int, out, errOut io.Writer) error {
	taskPath, err := utils.TaskPath(taskName)
	if err != nil {
		return err
	}
	if err := os.Chdir(taskPath); err != nil {
		return err
	}

	taskInputFilePath, err := utils.TaskInputFilePath(taskName, sampleNum)
	taskInputFilePathRel, err := filepath.Rel(utils.RootPath(), taskInputFilePath)
	if err != nil {
		return err
	}
	fmt.Fprintf(out, "Sample file: %s\n", taskInputFilePathRel)

	if err != nil {
		return err
	}
	if _, err := os.Stat(taskInputFilePath); err != nil {
		if os.IsNotExist(err) {
			return errSampleInputNotExist
		}
		return err
	}

	inBytes, err := ioutil.ReadFile(taskInputFilePath)
	if err != nil {
		return err
	}

	execCmd := exec.Command("./a.out")
	execCmd.Stdout = out
	execCmd.Stderr = errOut
	execCmdIn, err := execCmd.StdinPipe()
	if err != nil {
		return err
	}

	if err = execCmd.Start(); err != nil {
		return err
	}
	if _, err = execCmdIn.Write(inBytes); err != nil {
		return err
	}
	if err = execCmd.Wait(); err != nil {
		return err
	}
	return nil
}

func runWithSamples(taskName string, out, errOut io.Writer) error {
	for sampleNum := 0; sampleNum <= 99; sampleNum++ {
		if err := runWithSample(taskName, sampleNum, out, errOut); err != nil {
			if err != errSampleInputNotExist {
				return err
			}
		}
	}

	return nil
}
