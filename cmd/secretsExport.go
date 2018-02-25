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
	"log"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// secretsExportCmd represents the export command
var secretsExportCmd = &cobra.Command{
	Use:   "export",
	Short: "Export vault secrets in a secified format",
	Run: func(cmd *cobra.Command, args []string) {

		if viper.GetBool("decrypt") {
			if err := secretsCmd.vault.DecryptAll(); err != nil {
				log.Fatal(err)
			}
		}

		bytes, err := secretsCmd.vault.MarshalAs(viper.GetString("format"))
		if err != nil {
			log.Fatal(err)
		}
		fmt.Print(string(bytes))
	},
}

func init() {
	secretsCmd.AddCommand(secretsExportCmd)

	secretsExportCmd.Flags().Bool("decrypt", false, "Export the vault after decrypting it (default false)")
	secretsExportCmd.Flags().String("format", "", "The format to export the vault as (json, yaml, env, shell)")
	secretsExportCmd.MarkFlagRequired("format")

	viper.BindPFlag("decrypt", secretsExportCmd.Flags().Lookup("decrypt"))
	viper.BindPFlag("format", secretsExportCmd.Flags().Lookup("format"))

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsExportCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
