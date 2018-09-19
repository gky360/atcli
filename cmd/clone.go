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
	"github.com/gky360/atcli/utils"
	"github.com/gky360/atsrv/models"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type CloneOptions struct {
	Out, ErrOut io.Writer
	isFull      bool
}

var cloneOpt = &CloneOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// cloneCmd represents the clone command
var cloneCmd = &cobra.Command{
	Use:   "clone",
	Short: "Create source code files and download sample cases",
	Long: `Create source code files and download sample cases.

Example:
    atcli clone

"atcli clone" command generates files into the following directory structure.

$ATCLI_ROOT/
├── arc090
     ├── c
     │   ├── Main.cpp
     │   ├── samples
     │   │  ├── in
     │   │  │  ├── 01.txt
     │   │  │  ├── 02.txt
     │   │  │  ├── ...
     │   │  └── out
     │   │      ├── 01.txt
     │   │      ├── 02.txt
     │   │      ├── ...
     │   └── (testcases)
     │       ├── in
     │       │  ├── 01.txt
     │       │  ├── 02.txt
     │       │  ├── ...
     │       └── out
     │           ├── 01.txt
     │           ├── 02.txt
     │           ├── ...
     ├── d
     │   ├── Main.cpp
     │   ├── samples
     │   └── (testcases)
     ├── e
     │   └── ...
     └── f
          └── ...
`,
	Args: cobra.RangeArgs(0, 1),
	Run: func(cmd *cobra.Command, args []string) {
		if err := cloneOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(cloneOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(cloneCmd)
	cloneCmd.Flags().BoolVarP(&cloneOpt.isFull, "full", "", false, "download full testcases used in the contest.")

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *CloneOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")
	if err = runClone(contestID, opt.isFull, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	return nil
}

func runClone(contestID string, isFull bool, out, errOut io.Writer) error {
	contest := new(models.Contest)
	_, err := Client.GetContest(contestID, isFull, contest)
	if err != nil {
		return err
	}
	contestYaml, err := contest.ToYaml()
	if err != nil {
		return err
	}
	fmt.Fprintln(out, contestYaml)

	_, tasks, err := Client.GetTasks(contestID, true)
	if err != nil {
		return err
	}
	tasksYaml, err := models.TasksToYaml(tasks)
	if err != nil {
		return err
	}
	fmt.Fprintln(out, tasksYaml)

	fmt.Fprintf(out, "atcli root: %s\n", utils.RootDir())
	if err = utils.CreateFilesForTasks(contest, tasks); err != nil {
		return err
	}

	if isFull {
		if err := utils.DownloadTestcases(contest, tasks); err != nil {
			return err
		}
	}

	return nil
}
