package vault

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/mitchellh/hashstructure"
	"github.com/sourcec0de/gvault/crypter"
)

// Vault a vault that stores in a json format
type Vault struct {
	filePath  string
	crypter   *crypter.Crypter
	Version   uint64            `json:"version"`
	Secrets   map[string]string `json:"secrets"`
	isNew     bool
	decrypted bool
}

// SetSecret add a secret to the vault
func (v *Vault) SetSecret(key, value string) error {
	encValue, err := v.crypter.Encrypt([]byte(value))
	if err != nil {
		return err
	}
	v.Secrets[key] = encValue
	return nil
}

// RemoveSecret removes a secret from the vault
func (v *Vault) RemoveSecret(key string) {
	delete(v.Secrets, key)
}

// GetSecret gets a secret from the vault
func (v *Vault) GetSecret(key string) (string, error) {

	cipherText := v.Secrets[key]

	if cipherText == "" {
		return "", fmt.Errorf("No secret by that name")
	}

	secretBytes, err := v.crypter.Decrypt(cipherText)
	if err != nil {
		return "", err
	}

	return string(secretBytes), nil
}

func (v *Vault) toJSON() ([]byte, error) {
	return json.MarshalIndent(v.Secrets, "", "  ")
}

func (v *Vault) toYAML() ([]byte, error) {
	return yaml.Marshal(v.Secrets)
}

func (v *Vault) toENV() ([]byte, error) {
	env, err := godotenv.Marshal(v.Secrets)
	return []byte(env), err
}

func (v *Vault) toSHELL() ([]byte, error) {
	var output string
	env, err := godotenv.Marshal(v.Secrets)
	for _, line := range strings.Split(strings.TrimSuffix(env, "\n"), "\n") {
		output += ("export " + line + "\n")
	}
	return []byte(output), err
}

// MarshalAs marshals the vault secrets as the supplid format
func (v *Vault) MarshalAs(format string) ([]byte, error) {
	if format == "json" {
		return v.toJSON()
	}

	if format == "yml" || format == "yaml" {
		return v.toYAML()
	}

	if format == "env" {
		return v.toENV()
	}

	if format == "shell" {
		return v.toSHELL()
	}

	return nil, fmt.Errorf("%s is not a supported vault export format", format)
}

// Save writes the vault to it's storage location
func (v *Vault) Save() error {

	version, hashErr := v.HashSecrets()
	if hashErr != nil {
		return hashErr
	}

	v.Version = version

	bytes, jsonSaveErr := json.MarshalIndent(v, "", "  ")

	if jsonSaveErr != nil {
		return jsonSaveErr
	}

	if ioWriteErr := ioutil.WriteFile(v.filePath, bytes, os.ModePerm); ioWriteErr != nil {
		return ioWriteErr
	}

	return nil
}

// HashSecrets generates a unique hash of the encrypted secrets
// this is indended to be used as a version when syncronizing this with a secret store
// like kubernetes secrets
func (v *Vault) HashSecrets() (uint64, error) {
	return hashstructure.Hash(v.Secrets, nil)
}

// Load a vault from a filePath
func (v *Vault) Load() error {
	bytes, ioReadErr := ioutil.ReadFile(v.filePath)
	if ioReadErr != nil {
		return ioReadErr
	}

	json.Unmarshal(bytes, v)
	return nil
}

type kmsAPIResult struct {
	key   string
	value string
}

// DecryptAll decrypts all secrets in this vault
func (v *Vault) DecryptAll() error {
	totalReq := 0
	resultsChan := make(chan kmsAPIResult)

	for name, cipherText := range v.Secrets {
		totalReq++
		go func(name, cipherText string) {
			bytes, _ := v.crypter.Decrypt(cipherText)
			resultsChan <- kmsAPIResult{
				key:   name,
				value: string(bytes),
			}
		}(name, cipherText)
	}

	for {
		result := <-resultsChan
		v.Secrets[result.key] = result.value
		totalReq--

		if totalReq == 0 {
			break
		}
	}

	v.decrypted = true

	return nil
}

// EncryptEnvMap encrypts all secrets in a given envMap
func (v *Vault) EncryptEnvMap(envMap map[string]string) (map[string]string, error) {
	totalReq := 0
	resultsChan := make(chan kmsAPIResult)

	for key, plainText := range envMap {
		totalReq++
		go func(key, plainText string) {
			bytes, _ := v.crypter.Encrypt([]byte(plainText))
			resultsChan <- kmsAPIResult{
				key:   key,
				value: string(bytes),
			}
		}(key, plainText)
	}

	encryptedEnvMap := map[string]string{}

	for {
		result := <-resultsChan
		encryptedEnvMap[result.key] = result.value
		totalReq--

		if totalReq == 0 {
			break
		}
	}

	return encryptedEnvMap, nil
}

// MergeEncryptedEnvMap merges an encryptedEnvMap into the secrets
func (v *Vault) MergeEncryptedEnvMap(encryptedEnvMap map[string]string) {
	for key, value := range encryptedEnvMap {
		v.Secrets[key] = value
	}
}

// NewVault returns a pointer to a Vault
// if the vault does not exist on disk it will be created
// if the vault was not newely created it will attempt to load and unmarshal it
func NewVault(filePath string, crypter *crypter.Crypter) (*Vault, error) {

	vault := &Vault{
		filePath: filePath,
		crypter:  crypter,
		Secrets:  map[string]string{},
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(filePath), os.ModePerm)
		if _, err := os.Create(filePath); err != nil {
			return vault, err
		}
		vault.isNew = true
		vault.Save()
	}

	if !vault.isNew {
		if vaultLoadErr := vault.Load(); vaultLoadErr != nil {
			return vault, vaultLoadErr
		}
	}

	return vault, nil
}
