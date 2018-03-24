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

	. "github.com/gky360/atcli/constants"
	"github.com/gky360/atsrv/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetOptions struct {
	Out, ErrOut io.Writer
}

var getOpt = &GetOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// getCmd represents the get command
var getCmd = &cobra.Command{
	Use:   "get",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := getOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(getOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(getCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *GetOptions) Run(cmd *cobra.Command, args []string) (err error) {
	// get contest
	contestID := viper.GetString("contest.id")
	contest := new(models.Contest)
	if _, err = Client.GetContest(contestID, contest); err != nil {
		return err
	}

	// get tasks
	_, tasks, err := Client.GetTasks(contestID, false)
	if err != nil {
		return err
	}

	// get submissions
	var sbms []models.Submission
	if _, err = Client.GetSubmissions(contestID, "", &sbms); err != nil {
		return err
	}

	contestYaml, err := contest.ToYaml()
	if err != nil {
		return err
	}
	tasksYaml, err := models.TasksToYamlShort(tasks)
	if err != nil {
		return err
	}
	sbmsYaml, err := models.SubmissionsToYamlShort(sbms)
	if err != nil {
		return err
	}
	fmt.Fprintln(opt.Out, contestYaml)
	fmt.Fprintln(opt.Out, tasksYaml)
	fmt.Fprintln(opt.Out, sbmsYaml)

	return nil
}
