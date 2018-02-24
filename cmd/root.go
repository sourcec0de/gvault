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
	"log"
	"os"

	"github.com/sourcec0de/gvault/crypter"
	"github.com/sourcec0de/gvault/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile     string
	cfgFileName = ".gvaultrc"
	rcFile      = fmt.Sprintf("%s.json", cfgFileName)
)

type rootCmdWithCrypter struct {
	*cobra.Command
	crypter *crypter.Crypter
}

func (r *rootCmdWithCrypter) initCrypter() error {
	newCrypter, err := crypter.NewCrypter(viper.GetString("project"),
		viper.GetString("location"), viper.GetString("keyring"), viper.GetString("key"))

	if err != nil {
		return err
	}

	r.crypter = newCrypter
	return nil
}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &rootCmdWithCrypter{
	Command: &cobra.Command{
		Use:   "gvault",
		Short: "Manage secrets for your Google Cloud Platorm projects",
	},
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
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $PWD/.gvault.json)")
	rootCmd.PersistentFlags().StringP("project", "p", "", "Google Cloud ProjectID")
	rootCmd.PersistentFlags().StringP("keyring", "k", "", "Google KMS Keyring")
	rootCmd.PersistentFlags().StringP("location", "l", "", "Google KMS Keyring Location (defaults to global)")
	rootCmd.PersistentFlags().StringP("key", "", "", "Google KMS Key")

	viper.BindPFlag("project", rootCmd.PersistentFlags().Lookup("project"))
	viper.BindPFlag("keyring", rootCmd.PersistentFlags().Lookup("keyring"))
	viper.BindPFlag("location", rootCmd.PersistentFlags().Lookup("location"))
	viper.BindPFlag("key", rootCmd.PersistentFlags().Lookup("key"))

	viper.SetDefault("location", "global")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Search config in home directory with name ".gvault" (without extension).
		viper.AddConfigPath(utils.CWD())
		viper.SetConfigName(cfgFileName)
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match
	viper.ReadInConfig()

	// If a config file is found, read it in.
	// if err := viper.ReadInConfig(); err == nil {
	// 	fmt.Println("Using config file:", viper.ConfigFileUsed())
	// }
	// if err := viper.ReadInConfig(); err != nil {
	// 	fmt.Println(err)
	// }
	if err := rootCmd.initCrypter(); err != nil {
		log.Fatal(err)
	}
}
