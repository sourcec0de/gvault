// Copyright © 2018 James Qualls
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
	"path/filepath"

	"github.com/sourcec0de/gvault/utils"
	"github.com/sourcec0de/gvault/vault"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type secretsCmdWithVault struct {
	*cobra.Command
	vault *vault.Vault
}

func (s *secretsCmdWithVault) initVault() error {
	vaultPath := fmt.Sprintf(filepath.Join(utils.CWD(), "gvault", viper.GetString("vault")+".json"))
	newVault, err := vault.NewVault(vaultPath)

	if err != nil {
		return err
	}

	s.vault = newVault
	return nil
}

// secretsCmd represents the store command
var secretsCmd = &secretsCmdWithVault{
	Command: &cobra.Command{
		Use:   "secrets",
		Short: "Manage secrets stored in a vault",
	},
}

func init() {
	rootCmd.AddCommand(secretsCmd.Command)

	secretsCmd.PersistentFlags().StringP("vault", "v", "", "name of a local vault")
	viper.BindPFlag("vault", secretsCmd.PersistentFlags().Lookup("vault"))
	viper.SetDefault("vault", "main")

	if err := secretsCmd.initVault(); err != nil {
		panic(err)
	}

	// Here you will define your flags and configuration settings.

	// Cobra supports Persistent Flags which will work for this command
	// and all subcommands, e.g.:
	// secretsCmd.PersistentFlags().String("foo", "", "A help for foo")

	// Cobra supports local flags which will only run when this command
	// is called directly, e.g.:
	// secretsCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
