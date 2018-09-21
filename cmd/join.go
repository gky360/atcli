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

type JoinOptions struct {
	Out, ErrOut io.Writer
}

var joinOpt = &JoinOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// joinCmd represents the join command
var joinCmd = &cobra.Command{
	Use:   "join",
	Short: "Join a contest on AtCoder",
	Long: `Join a contest on AtCoder.

"atcli join" command will let you join a contest when a "Register" button
is displayed on the contest page.

If you already have joined the contest (i.e. the "Register" button is not
displayed), this command will fail.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := joinOpt.Run(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(joinCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// joinCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// joinCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *JoinOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")
	if err = runJoin(contestID, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	return nil
}

func runJoin(contestID string, out, errOut io.Writer) error {
	contest := new(models.Contest)
	if _, err := Client.Join(contestID, contest); err != nil {
		return err
	}

	contestYaml, err := contest.ToYaml()
	if err != nil {
		return err
	}

	fmt.Fprintln(out, contestYaml)

	return nil
}
