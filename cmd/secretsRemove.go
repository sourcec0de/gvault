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
	"github.com/spf13/cobra"
)

// secretsRemoveCmd represents the remove command
var secretsRemoveCmd = &cobra.Command{
	Use:   "remove",
	Short: "delete a secret from the vault",
	Long:  "gvault secrets remove SECRET_NAME",
	Args:  cobra.MinimumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		gvault.RemoveSecret(args[0])
		gvault.Save()
	},
}

func init() {
	secretsCmd.AddCommand(secretsRemoveCmd)
}
