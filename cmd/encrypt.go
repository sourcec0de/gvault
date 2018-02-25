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
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// encryptCmd represents the encrypt command
var encryptCmd = &cobra.Command{
	Use:   "encrypt",
	Short: "Encrypt a secret",
	Long: `Prints out a base64 encoded version of the secret you provided

$ gvault encrypt SUPER_AWESOME_SECRET
> CiQAuu4Laa3N0AwXlqDy1kTCZm3YdqEtrk/mpnsuHfMEDtNxCxISPQC8LsbdMQ1fjDsiRZn2p+HsXluLGaFG1YyQvahPHDAyXAQT1snON180ODweOIeo1MzoLYYtzHMNzC7vakg=`,
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) != 1 && !viper.GetBool("stdin") {
			return errors.New("requires one value to encrypt")
		}
		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var plainText []byte

		if viper.GetBool("stdin") {
			bytes, err := ioutil.ReadAll(os.Stdin)

			if err != nil {
				log.Fatal(err)
			}

			plainText = bytes
		} else {
			plainText = []byte(args[0])
		}

		encryptedData, err := rootCmd.crypter.Encrypt(plainText)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Print(encryptedData)
	},
}

func init() {
	rootCmd.AddCommand(encryptCmd)
	encryptCmd.Flags().Bool("stdin", false, "Read from stdin instead of arguments")
	viper.BindPFlag("stdin", encryptCmd.Flags().Lookup("stdin"))

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// encryptCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// encryptCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
