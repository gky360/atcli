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

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type NewOptions struct {
	Out, ErrOut io.Writer
}

var newOpt = &NewOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// newCmd represents the new command
var newCmd = &cobra.Command{
	Use:   "new",
	Short: "Join and clone a contest",
	Long: `Join and clone a contest.

"atcli new" command is a combination of "atcli join" and "atcli clone"
command.

If you already have joined the contest (i.e. the "Register" button is
not displayed), this command will fail. You just need to run "atcli clone"
command instead.`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := newOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(newOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(newCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// newCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// newCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *NewOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")
	if err = runJoin(contestID, opt.Out, opt.ErrOut); err != nil {
		return err
	}
	if err = runClone(contestID, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	return nil
}
