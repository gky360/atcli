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

	. "github.com/gky360/atcli/client"
	"github.com/gky360/atsrv/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type GetTaskOptions struct {
	Out, ErrOut io.Writer
}

var getTaskOpt = &GetTaskOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// getTaskCmd represents the getTask command
var getTaskCmd = &cobra.Command{
	Use:   "task [task name]",
	Short: "Get tasks from \"atsrv\"",
	Long: `Get tasks from "atsrv"

"atcli get task" command gets tasks from "atsrv" and prints the data
in yaml format.

When you specify task name, one task will be returned.`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := getTaskOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(getTaskOpt.ErrOut, err)
		}
	},
}

func init() {
	getCmd.AddCommand(getTaskCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// getTaskCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// getTaskCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *GetTaskOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")

	if len(args) >= 1 {
		taskName := args[0]
		task := new(models.Task)
		if _, err = Client.GetTask(contestID, taskName, task); err != nil {
			return err
		}

		taskYaml, err := task.ToYaml()
		if err != nil {
			return err
		}

		fmt.Fprintln(opt.Out, taskYaml)
	} else {
		_, tasks, err := Client.GetTasks(contestID, false)
		if err != nil {
			return err
		}

		tasksYaml, err := models.TasksToYaml(tasks)
		if err != nil {
			return err
		}

		fmt.Fprintln(opt.Out, tasksYaml)
	}

	return nil
}
