// Copyright © 2018 James Qualls https://github.com/sourcec0de
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
	"io/ioutil"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// secretsAddCmd represents the create command
var secretsAddCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a new secret to the vault",
	Long:  ``,
	Args: func(cmd *cobra.Command, args []string) error {
		file := viper.GetString("file")
		name := viper.GetString("name")
		usingFile := file != "" && name != ""

		if file != "" && name == "" {
			return fmt.Errorf("--name is required when using --file")
		}

		if len(args) < 1 && !usingFile {
			return fmt.Errorf("must supply at least one KEY=VALUE pair")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {

		file := viper.GetString("file")
		name := viper.GetString("name")

		if file != "" && name != "" {
			bytes, err := ioutil.ReadFile(name)

			if err != nil {
				logger.Fatal(err)
			}

			if err := gvault.SetSecret(name, string(bytes)); err != nil {
				logger.Fatal(err)
			}
		}

		for _, arg := range args {
			pair := strings.Split(arg, "=")
			if len(pair) != 2 {
				logger.Fatalf("%s is not a valid KEY=VALUE pair", arg)
			}

			if err := gvault.SetSecret(pair[0], pair[1]); err != nil {
				logger.Fatal(err)
			}
		}

		if err := gvault.Save(); err != nil {
			logger.Fatal(err)
		}
	},
}

func init() {
	secretsCmd.AddCommand(secretsAddCmd)

	secretsAddCmd.Flags().String("file", "", "The file to be encrypted")
	secretsAddCmd.Flags().String("name", "", "The name of the secret being added to the vault (only works with --file)")
	viper.BindPFlag("file", secretsAddCmd.Flags().Lookup("file"))
	viper.BindPFlag("name", secretsAddCmd.Flags().Lookup("name"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsAddCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsAddCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
