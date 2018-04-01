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

	. "github.com/gky360/atcli/client"
	"github.com/gky360/atcli/utils"
	"github.com/gky360/atsrv/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type SubmitOptions struct {
	Out, ErrOut io.Writer
}

var submitOpt = &SubmitOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// submitCmd represents the submit command
var submitCmd = &cobra.Command{
	Use:   "submit",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Args: cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := submitOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(submitOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(submitCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// submitCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// submitCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *SubmitOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")
	taskName := args[0]
	taskSourceFilePath, err := utils.TaskSourceFilePath(taskName)
	if err != nil {
		return err
	}

	buf, err := ioutil.ReadFile(taskSourceFilePath)
	if err != nil {
		return err
	}
	sbm := new(models.Submission)
	if _, err := Client.PostSubmission(contestID, taskName, string(buf), sbm); err != nil {
		return err
	}

	fmt.Fprintln(opt.Out, "Submit succeeded.")

	sbmYaml, err := sbm.ToYaml()
	if err != nil {
		return err
	}

	fmt.Fprintln(opt.Out, sbmYaml)

	return nil
}
