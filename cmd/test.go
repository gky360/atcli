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
	"strings"

	"github.com/fatih/color"
	"github.com/gky360/atcli/utils"
	"github.com/spf13/cobra"
)

type TestOptions struct {
	Out, ErrOut io.Writer

	isSkip  bool
	isFull  bool
	isQuiet bool
}

var testOpt = &TestOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test [task name [sample name]]",
	Short: "Build, run and test your source code",
	Long: `Build, run and test your source code.

"atcli test" command builds your source code, executes it with the
downloaded sample inputs passed as stdin, prints the stdout and stderr,
and check if your source code is correct by comparing the stdout and
downloaded sample outputs.

If you specify a sample name, this command only runs for the
specified sample input and output.

This command ignores leading and trailing spaces and line breaks when
it compares the stdout and sample outputs.

Note that this command only compares the outputs as strings, thus it
can not make correct judgements for tasks that accept multiple answers.`,
	Args: cobra.RangeArgs(1, 2),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := testOpt.Run(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.Flags().BoolVarP(&testOpt.isSkip, "skip-build", "s", false, "skip build if possible.")
	testCmd.Flags().BoolVarP(&testOpt.isFull, "full", "", false, "execute with full testcases inputs.")
	testCmd.Flags().BoolVarP(&testOpt.isQuiet, "quiet", "q", false, "don't show stdout of your program.")

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
	sampleName := ""
	if len(args) >= 2 {
		sampleName = args[1]
	}

	if err := runBuild(taskName, !opt.isSkip, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	if sampleName == "" {
		if err := testWithSamples(taskName, opt.isFull, opt.isQuiet, opt.Out, opt.ErrOut); err != nil {
			return err
		}
	} else {
		if _, err := testWithSample(taskName, sampleName, opt.isFull, opt.isQuiet, opt.Out, opt.ErrOut); err != nil {
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

func convNewline(str, nlcode string) string {
	return strings.NewReplacer(
		"\r\n", nlcode,
		"\r", nlcode,
		"\n", nlcode,
	).Replace(str)
}

func normalizeStr(str string) string {
	return convNewline(strings.TrimSpace(str), "\n")
}

func testWithSample(taskName string, sampleName string, isFull, isQuiet bool, out, errOut io.Writer) (bool, error) {
	res, err := runWithSample(taskName, sampleName, isFull, isQuiet, out, errOut)
	if err != nil {
		return false, err
	}

	taskOutputFilePath, err := utils.TaskOutputFilePath(taskName, sampleName, isFull)
	if err != nil {
		return false, err
	}
	sampleOutByte, err := ioutil.ReadFile(taskOutputFilePath)
	if err != nil {
		return false, err
	}
	sampleOut := string(sampleOutByte)

	isPass := false
	if strings.Compare(normalizeStr(res), normalizeStr(sampleOut)) == 0 {
		successColor := color.New(color.FgGreen)
		successColor.Fprintln(out, "Test: pass")
		isPass = true
	} else {
		failureColor := color.New(color.FgRed)
		if !isQuiet {
			failureColor.Fprintln(out, "Correct output:")
			fmt.Fprintln(out, sampleOut)
		}
		failureColor.Fprintln(out, "Test: fail")
	}

	return isPass, nil
}

func testWithSamples(taskName string, isFull, isQuiet bool, out, errOut io.Writer) error {
	totalCount := 0
	passCount := 0

	sampleNames, err := utils.GetSampleNames(taskName, isFull)
	if err != nil {
		return err
	}
	for _, sampleName := range sampleNames {
		isPass, err := testWithSample(taskName, sampleName, isFull, isQuiet, out, errOut)
		if err != nil {
			return err
		}
		totalCount++
		if isPass {
			passCount++
		}
	}

	var reportColor *color.Color
	statusStr := ""
	if isFull {
		statusStr = "testcases "
	} else {
		statusStr = "samples "
	}
	if passCount == totalCount {
		// passed all sample cases
		reportColor = color.New(color.FgBlack, color.BgHiGreen)
		statusStr += "AC"
	} else {
		reportColor = color.New(color.FgBlack, color.BgHiRed)
		statusStr += "WA"
	}
	fmt.Fprintln(out)
	reportColor.Fprintf(out, "%s (pass: %d, fail: %d, total: %d)\n",
		statusStr, passCount, totalCount-passCount, totalCount)

	return nil
}
