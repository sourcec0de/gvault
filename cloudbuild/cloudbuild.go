package cloudbuild

import (
	"github.com/ghodss/yaml"
)

// Secret a cloudbuild secret
type Secret struct {
	KmsKeyName string            `json:"kmsKeyName"`
	SecretEnv  map[string]string `json:"secretEnv"`
}

// Build a cloudbuild container
type Build struct {
	Secrets []Secret `json:"secrets"`
}

// MarshalToYAML encodes the build to YAML
func (b *Build) MarshalToYAML() ([]byte, error) {
	return yaml.Marshal(b)
}
