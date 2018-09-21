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
	"strconv"

	. "github.com/gky360/atcli/client"
	"github.com/gky360/atsrv/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetSbmOptions struct {
	Out, ErrOut io.Writer

	taskName string
	status   string
}

var getSbmOpt = &GetSbmOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// getSubmissionCmd represents the getSubmission command
var getSubmissionCmd = &cobra.Command{
	Use:     "submission [submission id]",
	Aliases: []string{"sbm"},
	Short:   "Get submissions from \"atsrv\"",
	Long: `Get submissions from "atsrv"

"atcli get submission" command gets submissions and prints the data
in yaml format.

When you specify submission id, one submission will be returned.
When you specify task name flag or submission status flag, multiple
submissions will be returned filtered by the flags.`,
	Args: cobra.RangeArgs(0, 1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := getSbmOpt.Run(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	getCmd.AddCommand(getSubmissionCmd)

	getSubmissionCmd.Flags().StringVarP(&getSbmOpt.taskName, "task", "t", "", "task name")
	getSubmissionCmd.Flags().StringVarP(&getSbmOpt.status, "status", "s", "", "submission status")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getSubmissionCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getSubmissionCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *GetSbmOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")

	if len(args) >= 1 {
		sbmID, err := strconv.Atoi(args[0])
		if err != nil {
			return err
		}
		sbm := new(models.Submission)
		if _, err = Client.GetSubmission(contestID, sbmID, sbm); err != nil {
			return err
		}

		sbmYaml, err := sbm.ToYaml()
		if err != nil {
			return err
		}

		fmt.Fprintln(opt.Out, sbmYaml)
	} else {
		_, sbms, err := Client.GetSubmissions(contestID, getSbmOpt.taskName, getSbmOpt.status)
		if err != nil {
			return err
		}

		sbmsYaml, err := models.SubmissionsToYaml(sbms)
		if err != nil {
			return err
		}

		fmt.Fprintln(opt.Out, sbmsYaml)
	}

	return nil
}
