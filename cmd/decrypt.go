// Copyright © 2018 NAME HERE qualls.james@gmail.com
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

	"github.com/sourcec0de/gvault/utils"
	"github.com/spf13/cobra"
)

var decryptLongExample = `
Decrypt the provided secret

$ gvault decrypt BASE64_ENCODED_CIPHER_TEXT
> SUPER_AWESOME_SECRET

Alternatively, you can pass data in via STDIN
**Note** the usage of "-" as the first argument

$ echo BASE64_ENCODED_CIPHER_TEXT | gvault decrypt -
> SUPER_AWESOME_SECRET
`

// decryptCmd represents the decrypt command
var decryptCmd = &cobra.Command{
	Use:   "decrypt",
	Short: "Decrypt a secret",
	Long:  decryptLongExample,
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {

		var cipherText string

		if args[0] == "-" {
			cipherText = string(utils.ReadAllStdin())
		} else {
			cipherText = args[0]
		}

		decrypted, err := rootCmd.crypter.Decrypt(cipherText)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(string(decrypted))
	},
}

func init() {
	rootCmd.AddCommand(decryptCmd)

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// decryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// decryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
