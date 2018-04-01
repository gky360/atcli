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
	"path/filepath"

	. "github.com/gky360/atcli/client"
	"github.com/gky360/atcli/utils"
	homedir "github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type RootOptions struct {
	Out, ErrOut io.Writer

	cfgFile   string
	host      string
	port      string
	contestID string
	userID    string
	token     string
}

const (
	Version = "v0.0.1"
)

var (
	banner = fmt.Sprintf(`
        __           ___
       /\ \__       /\_ \    __
   __  \ \ ,_\   ___\//\ \  /\_\
 /'__'\ \ \ \/  /'___\\ \ \ \/\ \
/\ \L\.\_\ \ \_/\ \__/ \_\ \_\ \ \
\ \__/.\_\\ \__\ \____\/\____\\ \_\
 \/__/\/_/ \/__/\/____/\/____/ \/_/
%38s
`, Version)
)

var cfgFile string

var rootOpt = &RootOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atcli",
	Short: "\"atcli\" is a command line interface for AtCoder (unofficial).",
	Long: fmt.Sprintf("%s\n%s", banner,
		`"atcli" is a command line interface for AtCoder (unofficial).

A basic flow of the usage of this command is as follows.

    1. Start "atsrv" and get auth token.
       (see https://github.com/gky360/atsrv for details)
    2. Set user, contest and the auth token to the config file
       (~/.atcli.yaml) using "atcli config" command.
    3. Join contest using "atcli join" command.
    4. Generate empty source code file and download sample cases
       from AtCoder using "atcli clone" command.
    5. Write your code to the generated source code file (Main.cpp).
    6. Test your code with downloaded sample cases
       using "atcli test" command.
    7. Submit you code to AtCoder using "atcli submit" command.
    8. Check your submission status
       using "atcli get submission" command.`),
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize(initConfig)

	// global flags
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.atcli.yaml)")

	rootCmd.PersistentFlags().StringVarP(&rootOpt.host, "host", "H", "", "atsrv host")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))

	rootCmd.PersistentFlags().StringVarP(&rootOpt.port, "port", "P", "4700", "atsrv port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))

	rootCmd.PersistentFlags().StringVarP(&rootOpt.token, "auth-token", "a", "", "auth token for atsrv")
	viper.BindPFlag("auth-token", rootCmd.PersistentFlags().Lookup("auth-token"))

	rootCmd.PersistentFlags().StringVarP(&rootOpt.token, "root", "r", utils.DefaultRootPath(), "root directory where atcli create files")
	viper.BindPFlag("root", rootCmd.PersistentFlags().Lookup("root"))

	rootCmd.PersistentFlags().StringVarP(&rootOpt.contestID, "contest", "c", "", "contest id of AtCoder")
	viper.BindPFlag("contest.id", rootCmd.PersistentFlags().Lookup("contest"))

	rootCmd.PersistentFlags().StringVarP(&rootOpt.userID, "user", "u", "", "user id of AtCoder")
	viper.BindPFlag("user.id", rootCmd.PersistentFlags().Lookup("user"))
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".atcli" (without extension).
		viper.AddConfigPath(home)
		viper.SetConfigName(".atcli")
		viper.SetConfigType("yaml")
		cfgFile = filepath.Join(home, ".atcli.yaml")
		if _, err = os.OpenFile(cfgFile, os.O_RDONLY|os.O_CREATE, 0644); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	viper.SetEnvPrefix("atcli")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// set access token to http client
	host := viper.GetString("host")
	port := viper.GetString("port")
	userID := viper.GetString("user.id")
	authToken := viper.GetString("auth-token")
	Client = NewClient(host, port)
	if userID != "" && authToken != "" {
		Client.SetBasicAuth(userID, authToken)
	}
}
