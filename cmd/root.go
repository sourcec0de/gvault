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

	log "github.com/sirupsen/logrus"
	prefixed "github.com/x-cray/logrus-prefixed-formatter"

	"os"

	"github.com/sourcec0de/gvault/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var logger *log.Logger
var gvault = vault.New(vault.Config{})

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "gvault",
	Short: "Manage secrets for your Google Cloud Platorm projects",
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

	// init cobra
	cobra.OnInitialize(initConfig)
	cobra.OnInitialize(initVault)
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.
	rootCmd.PersistentFlags().Bool("debug", false, "Enable debug statements")
	rootCmd.PersistentFlags().StringP("vault", "v", "", "The name of the vault you want to use (default to main)")
	viper.BindPFlag("vault", rootCmd.PersistentFlags().Lookup("vault"))
	viper.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))
	viper.SetDefault("vault", "main")
	viper.SetEnvPrefix("GVAULT")

	// init logger
	logger = log.New()
	logger.Formatter = new(prefixed.TextFormatter)

	if initCrypterErr := gvault.InitCrypter(); initCrypterErr != nil {
		log.Fatal(initCrypterErr)
	}
}

func initConfig() {
	viper.AutomaticEnv()
	viper.ReadInConfig()

	if viper.GetBool("debug") {
		logger.SetLevel(log.DebugLevel)
	} else {
		logger.SetLevel(log.InfoLevel)
	}
}

func initVault() {

	gvault.Name = viper.GetString("vault")

	if exists, _ := gvault.Exists(); exists {
		if loadErr := gvault.Load(); loadErr != nil {
			logger.Fatal(loadErr)
		}
		logger.Debugf("Using vault (%s)", gvault.Path())
	}
}
