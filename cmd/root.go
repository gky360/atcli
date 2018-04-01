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

var cfgFile string

var rootOpt = &RootOptions{
	Out:    os.Stdout,
	ErrOut: os.Stderr,
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "atcli",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
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

	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.atcli.yaml)")
	rootCmd.PersistentFlags().StringVarP(&rootOpt.host, "host", "H", "", "atsrv host")
	viper.BindPFlag("host", rootCmd.PersistentFlags().Lookup("host"))
	rootCmd.PersistentFlags().StringVarP(&rootOpt.port, "port", "P", "4700", "atsrv port")
	viper.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
	rootCmd.PersistentFlags().StringVarP(&rootOpt.contestID, "contest", "c", "", "contest id")
	viper.BindPFlag("contest.id", rootCmd.PersistentFlags().Lookup("contest"))
	rootCmd.PersistentFlags().StringVarP(&rootOpt.userID, "user", "u", "", "user id of AtCoder")
	viper.BindPFlag("user.id", rootCmd.PersistentFlags().Lookup("user"))
	rootCmd.PersistentFlags().StringVarP(&rootOpt.token, "auth-token", "a", "", "token for atsrv")
	viper.BindPFlag("auth-token", rootCmd.PersistentFlags().Lookup("auth-token"))

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	// rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
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
