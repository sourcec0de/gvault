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

	"github.com/chzyer/readline"
	"github.com/sourcec0de/gvault/utils"
	"github.com/spf13/cobra"
)

// initCmd represents the init command
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize a new gvault",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		if exists, err := gvault.Exists(); exists || err != nil {
			if exists {
				return fmt.Errorf("A gvault with the name (%s) already exists: %s", gvault.Name, gvault.Path())
			}
			return err
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
		location, _ := utils.Ask("Google KMS Keyring Location (defaults to global): ", rl)
		key, _ := utils.Ask("Google KMS Key: ", rl)

		if location == "" {
			location = "global"
		}

		gvault.Project = project
		gvault.Keyring = keyring
		gvault.Location = location
		gvault.Key = key

		if saveErr := gvault.Save(); saveErr != nil {
			logger.Fatal(saveErr)
		}

		fmt.Printf("Created %s", gvault.Path())
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
