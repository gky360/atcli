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

type EnvOptions struct {
	Out, ErrOut io.Writer
}

var envOpt = &EnvOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// envCmd represents the env command
var envCmd = &cobra.Command{
	Use:   "env",
	Short: "Print environment variables related to atcli and atsrv",
	Long: `Print environment variables related to atcli and atsrv.

If you set some global flags, "atcli env" command will show the value
you specified using the flag instead of the value set via envvars.`,
	Args: cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := envOpt.Run(cmd, args); err != nil {
			return err
		}
		return nil
	},
}

func init() {
	rootCmd.AddCommand(envCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// cloneCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// cloneCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *EnvOptions) Run(cmd *cobra.Command, args []string) error {
	runEnv(opt.Out)
	return nil
}

func runEnv(out io.Writer) {
	fmt.Fprintf(out, "ATSRV_HOST=%s\n", viper.GetString("host"))
	fmt.Fprintf(out, "ATSRV_PORT=%s\n", viper.GetString("port"))
	fmt.Fprintf(out, "ATSRV_USER_ID=%s\n", viper.GetString("user.id"))
	fmt.Fprintf(out, "ATSRV_AUTH_TOKEN=%s\n", viper.GetString("auth-token"))
	fmt.Fprintf(out, "ATCLI_ROOT=%s\n", viper.GetString("root"))
	fmt.Fprintf(out, "ATCLI_CONTEST_ID=%s\n", viper.GetString("contest.id"))
	fmt.Fprintf(out, "ATCLI_CPP_TEMPLATE_PATH=%s\n", viper.GetString("cppTemplatePath"))
}
