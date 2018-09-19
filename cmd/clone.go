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
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox"
	"github.com/dropbox/dropbox-sdk-go-unofficial/dropbox/files"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	. "github.com/gky360/atcli/client"
	"github.com/gky360/atcli/utils"
	"github.com/gky360/atsrv/models"
)

type CloneOptions struct {
	Out, ErrOut io.Writer
}

const atcoderDropboxRootURL = "https://www.dropbox.com/sh/arnpe0ef5wds8cv/AAAk_SECQ2Nc6SVGii3rHX6Fa"

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
     │   └── samples
     │       ├── 01.in.txt
     │       ├── 01.out.txt
     │       ├── 02.in.txt
     │       ├── 02.out.txt
     │       ├── ...
     ├── d
     │   ├── Main.cpp
     │   └── samples
     │       ├── ...
     ├── e
     │   ├── Main.cpp
     │   └── samples
     │       ├── ...
     └── f
         ├── Main.cpp
         └── samples
             ├── ...
`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := cloneOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(cloneOpt.ErrOut, err)
		}
	},
}

var dbx files.Client

func init() {
	rootCmd.AddCommand(cloneCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")

	dbxToken := viper.GetString("dropboxAccessToken")
	dbxConfig := dropbox.Config{
		Token:    dbxToken,
		LogLevel: dropbox.LogOff,
	}
	dbx = files.New(dbxConfig)
}

func (opt *CloneOptions) Run(cmd *cobra.Command, args []string) (err error) {
	contestID := viper.GetString("contest.id")
	if err = runClone(contestID, opt.Out, opt.ErrOut); err != nil {
		return err
	}

	return nil
}

func runClone(contestID string, out, errOut io.Writer) error {
	contest := new(models.Contest)
	_, err := Client.GetContest(contestID, contest)
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

	fmt.Fprintf(out, "atcli root: %s\n", utils.RootPath())
	if err = utils.CreateFilesForTasks(contest, tasks); err != nil {
		return err
	}

	if err := downloadTestcases(contest.ID); err != nil {
		fmt.Fprintln(cloneOpt.ErrOut, err)
		return err
	}

	return nil
}

func getContestDropboxURL(contestID string) error {
	arg := files.NewListFolderArg("")
	arg.SharedLink = files.NewSharedLink(atcoderDropboxRootURL)
	argJSONStr, _ := json.Marshal(arg)
	fmt.Fprintln(cloneOpt.Out, string(argJSONStr))

	res, err := dbx.ListFolder(arg)
	if err != nil {
		return fmt.Errorf("Could not list folders in %s", atcoderDropboxRootURL)
	}

	normalizedContestID := normalizeContestDropboxFolderName(contestID)
	for _, e := range res.Entries {
		folderMeta, ok := e.(*files.FolderMetadata)
		if !ok {
			// not a folder
			continue
		}
		folderName := folderMeta.Name
		if normalizeContestDropboxFolderName(folderName) == normalizedContestID {
		}
	}

	return nil
}

func downloadTestcases(contestID string) error {

	// resJSONBytes, _ := json.Marshal(res)
	// fmt.Fprintln(cloneOpt.Out, string(resJSONBytes))

	// contestFolderUrl := contestFolderMetadata.Url
	// fmt.Fprintf(cloneOpt.Out, "Downloading testcases from %s ...", contestFolderUrl)

	// zipf, err := ioutil.TempFile("", contest.ID)
	// if err != nil {
	// 	return err
	// }
	// defer zipf.Close()
	//
	// respHead, err := http.Head(contestFolderUrl)
	// if err != nil {
	// 	return err
	// }
	// defer respHead.Body.Close()
	//
	// size, err := strconv.Atoi(respHead.Header.Get("Content-Length"))
	// if err != nil{
	// 	return err
	// }
	//
	// done := make(chan int64)

	// DownloadFile("/Users/inagaki/tmp/testcases", "https://www.slimjet.com/chrome/download-chrome.php?file=files%2F69.0.3497.92%2FChromeStandaloneSetup64.exe")

	return nil
}
