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
	log "github.com/sirupsen/logrus"

	"github.com/sourcec0de/gvault/cloudbuild"
	"github.com/spf13/cobra"
)

var cloudBuildLongExample = `
Export encrypted secrets to your cloudbuild.yaml file

gvault cloudbuild >> cloudbuild.yaml
`

// cloudbuildCmd represents the cloudbuild command
var cloudbuildCmd = &cobra.Command{
	Use:   "cloudbuild",
	Short: "Export your vault in a format accepted by Google Container Builder cloudbuild.yml",
	Long:  cloudBuildLongExample,
	Run: func(cmd *cobra.Command, args []string) {
		build := cloudbuild.Build{}
		secret := cloudbuild.Secret{
			KmsKeyName: secretsCmd.vault.KmsKeyName(),
			SecretEnv:  secretsCmd.vault.Secrets,
		}

		build.Secrets = append(build.Secrets, secret)

		bytes, marshalErr := build.MarshalToYAML()
		if marshalErr != nil {
			log.Fatal(marshalErr)
		}
		fmt.Println(string(bytes))

	},
}

func init() {
	rootCmd.AddCommand(cloudbuildCmd)
}
