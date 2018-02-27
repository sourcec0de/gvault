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

	"github.com/sourcec0de/gvault/utils"
	"github.com/sourcec0de/gvault/vault"
	"github.com/spf13/cobra"
)

var encryptLongExample = `
Encrypt the provided secret

$ gvault encrypt SUPER_AWESOME_SECRET
> BASE64_ENCODED_CIPHER_TEXT

Alternatively, you can pass data in via STDIN
**Note** the usage of "-" as the first argument

$ cat service.json | gvault encrypt -
> BASE64_ENCODED_CIPHER_TEXT
`

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:     "encrypt",
	Short:   "Encrypt a secret",
	Long:    encryptLongExample,
	Args:    cobra.MinimumNArgs(1),
	PreRunE: vault.EsureVaultLoaded(gvault),
	Run: func(cmd *cobra.Command, args []string) {
		var plainText []byte

		if args[0] == "-" {
			plainText = utils.ReadAllStdin()
		} else {
			plainText = []byte(args[0])
		}

		encryptedData, err := gvault.Crypter.Encrypt(plainText)

		if err != nil {
			logger.Fatal(err)
		}

		fmt.Print(encryptedData)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	// encryptCmd.Flags().Bool("stdin", false, "Read from stdin instead of arguments")
	// viper.BindPFlag("stdin", encryptCmd.Flags().Lookup("stdin"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
