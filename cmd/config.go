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
	"gopkg.in/yaml.v2"
)

type ConfigOptions struct {
	Out, ErrOut io.Writer
}

var configOpt = &ConfigOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// configCmd represents the config command
var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Get and set global options",
	Long: `Get and set global options.

"atcli config" command reads current config from environment variables
and "~/.atcli.yaml", modify it according to the passed flags, save it
to "~/.atcli.yaml", and then prints the new config.

The priority of the config is
command flags > env vars > config file .`,
	Args: cobra.NoArgs,
	Run: func(cmd *cobra.Command, args []string) {
		if err := configOpt.Run(cmd, args); err != nil {
			fmt.Fprintln(configOpt.ErrOut, err)
		}
	},
}

func init() {
	rootCmd.AddCommand(configCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// configCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// configCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}

func (opt *ConfigOptions) Run(cmd *cobra.Command, args []string) (err error) {
	fmt.Fprintln(opt.Out, "confingCmd")

	if err = viper.WriteConfig(); err != nil {
		return err
	}

	var config interface{}
	if err = viper.Unmarshal(&config); err != nil {
		return err
	}
	d, err := yaml.Marshal(&config)
	if err != nil {
		return err
	}

	fmt.Fprintf(opt.Out, "---\n%s\n", string(d))

	return nil
}
