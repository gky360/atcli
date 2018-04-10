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
	"io/ioutil"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
	"github.com/gky360/atcli/utils"
	"github.com/spf13/cobra"
)

type TestOptions struct {
	Out, ErrOut io.Writer

	isSkip bool
}

var testOpt = &TestOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [task name [sample number]]",
	Short: "Build, run and test your source code",
	Long: `Build, run and test your source code.

"atcli test" command builds your source code, executes it with the
downloaded sample inputs passed as stdin, prints the stdout and stderr,
and check if your source code is correct by comparing the stdout and
downloaded sample outputs.

If you specify a sample number, this command only runs for the
specified sample input and output.

This command ignores leading and trailing spaces and line breaks when
it compares the stdout and sample outputs.

Note that this command only compares the outputs as strings, thus it
can not give correct judges for tasks that accept multiple answers.`,
	Args: cobra.RangeArgs(1, 2),
	Run: func(cmd *cobra.Command, args []string) {
		if err := testOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(testOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVarP(&testOpt.isSkip, "skip-build", "s", false, "Skip build if possible.")

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

	if err := runBuild(taskName, !opt.isSkip, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	if sampleNum >= 0 {
		if _, err := testWithSample(taskName, sampleNum, opt.Out, opt.ErrOut); err != nil {
			return err
		}
	} else {
		if err := testWithSamples(taskName, opt.Out, opt.ErrOut); err != nil {
			return err
		}
	}

	return nil
}

type SampleInputNotExistError struct {
	msg string
}

func (e SampleInputNotExistError) Error() string {
	return e.msg
}

func NewSampleInputNotExistError(msg string) *SampleInputNotExistError {
	return &SampleInputNotExistError{msg}
}

func testWithSample(taskName string, sampleNum int, out, errOut io.Writer) (bool, error) {
	res, err := runWithSample(taskName, sampleNum, out, errOut)
	if err != nil {
		return false, err
	}

	taskOutputFilePath, err := utils.TaskOutputFilePath(taskName, sampleNum)
	if err != nil {
		return false, err
	}
	sampleOutByte, err := ioutil.ReadFile(taskOutputFilePath)
	if err != nil {
		return false, err
	}
	sampleOut := string(sampleOutByte)

	isPass := false
	if strings.Compare(strings.TrimSpace(res), strings.TrimSpace(sampleOut)) == 0 {
		successColor := color.New(color.FgGreen)
		successColor.Fprintln(out, "Test: pass")
		isPass = true
	} else {
		failureColor := color.New(color.FgRed)
		failureColor.Fprintln(out, "Test: fail")
		failureColor.Fprintln(out, "Correct output:")
		fmt.Fprintln(out, sampleOut)
	}

	return isPass, nil
}

func testWithSamples(taskName string, out, errOut io.Writer) error {
	totalCount := 0
	passCount := 0
	for sampleNum := 0; sampleNum <= 99; sampleNum++ {
		isPass, err := testWithSample(taskName, sampleNum, out, errOut)
		switch err := err.(type) {
		case nil:
			totalCount++
			if isPass {
				passCount++
			}
		case *SampleInputNotExistError:
			// task did not have sample with this sampleNum
		default:
			return err
		}
	}

	var reportColor *color.Color
	statusStr := ""
	if passCount == totalCount {
		// passed all sample cases
		reportColor = color.New(color.FgBlack, color.BgHiGreen)
		statusStr = "samples AC"
	} else {
		reportColor = color.New(color.FgBlack, color.BgHiRed)
		statusStr = "samples WA"
	}
	fmt.Fprintln(out)
	reportColor.Fprintf(out, "%s (pass: %d, fail: %d, total: %d)\n",
		statusStr, passCount, totalCount-passCount, totalCount)

	return nil
}
