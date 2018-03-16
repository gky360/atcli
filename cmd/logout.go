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

type LogoutOptions struct {
	Out, ErrOut io.Writer
}

var logoutOpt = &LogoutOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// logoutCmd represents the logout command
var logoutCmd = &cobra.Command{
	Use:   "logout",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := logoutOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(logoutOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(logoutCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// logoutCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// logoutCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *LogoutOptions) Run(cmd *cobra.Command, args []string) (err error) {
	fmt.Fprintln(opt.Out, "logoutCmd")

	user := new(models.User)
	_, err = Client.Logout(user)
	if err != nil {
		return err
	}

	viper.Set("user.token", "")
	if err = viper.WriteConfig(); err != nil {
		return err
	}

	fmt.Fprintln(opt.Out, "Successfully logged out.")

	return nil
}
