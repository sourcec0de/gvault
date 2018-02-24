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
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/chzyer/readline"
	"github.com/sourcec0de/gvault/utils"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new gvault",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if viper.ConfigFileUsed() != "" {
			return fmt.Errorf("This project is already initialized. gvaultrc.json already exists: %s", viper.ConfigFileUsed())
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		rl, err := readline.NewEx(&readline.Config{
			UniqueEditLine: true,
		})

		if err != nil {
			log.Fatal(err)
		}

		defer rl.Close()

		project, _ := utils.Ask("Google Cloud ProjectID: ", rl)
		keyring, _ := utils.Ask("Google KMS Keyring: ", rl)
		key, _ := utils.Ask("Google KMS Key: ", rl)

		if project == "" || keyring == "" || key == "" {
			log.Fatal(fmt.Errorf("must provide an answer to all questions"))
		}

		bytes, _ := json.MarshalIndent(map[string]string{
			"project": project,
			"keyring": keyring,
			"key":     key,
		}, "", "  ")

		if err := ioutil.WriteFile(filepath.Join(utils.CWD(), "gvaultrc.json"), bytes, 0744); err != nil {
			panic(err)
		}

		fmt.Println("Created gvaultrc.json")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// initCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// initCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
