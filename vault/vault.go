package vault

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/ghodss/yaml"
	"github.com/joho/godotenv"
	"github.com/mitchellh/hashstructure"
	"github.com/pkg/errors"
	"github.com/sourcec0de/gvault/crypter"
	"github.com/sourcec0de/gvault/utils"
	"github.com/spf13/cobra"
)

var (
	gvaultFolder = "gvault"
)

var validationErrMsg = `
When initializing a new vault you must supply
--project
--keyring
--location (defaults to global)
--key
`

// Vault a vault that stores in a json format
type Vault struct {
	Name      string            `json:"-"`
	Version   uint64            `json:"version"`
	Secrets   map[string]string `json:"secrets"`
	Project   string            `json:"project"`
	Keyring   string            `json:"keyring"`
	Location  string            `json:"location"`
	Key       string            `json:"key"`
	Crypter   *crypter.Crypter  `json:"-"`
	isNew     bool
	loaded    bool
	decrypted bool
}

// Config a config for initializing a vault
type Config struct {
	Name     string
	Project  string
	Keyring  string
	Location string
	Key      string
}

// Path returns the path of the current vault instance
func (v *Vault) Path() string {
	return fmt.Sprintf(filepath.Join(utils.CWD(), gvaultFolder, v.Name+".json"))
}

// SetSecret add a secret to the vault
func (v *Vault) SetSecret(key, value string) error {
	encValue, err := v.Crypter.Encrypt([]byte(value))
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

// KmsKeyName name of the KMS resrouce
func (v *Vault) KmsKeyName() string {
	return v.Crypter.KmsKeyName()
}

// GetSecret gets a secret from the vault
func (v *Vault) GetSecret(key string) (string, error) {

	cipherText := v.Secrets[key]

	if cipherText == "" {
		return "", fmt.Errorf("No secret by that name")
	}

	secretBytes, err := v.Crypter.Decrypt(cipherText)
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

	if validationErr := v.validate(); validationErr != nil {
		if v.isNew {
			return errors.Wrap(validationErr, validationErrMsg)
		}
		return validationErr
	}

	if createErr := v.CreateIfNotExists(); createErr != nil {
		return createErr
	}

	if ioWriteErr := ioutil.WriteFile(v.Path(), bytes, os.ModePerm); ioWriteErr != nil {
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
	bytes, ioReadErr := ioutil.ReadFile(v.Path())
	if ioReadErr != nil {
		return errors.Wrap(ioReadErr, "failed to read vault file")
	}

	if unmarshalErr := json.Unmarshal(bytes, v); unmarshalErr != nil {
		return errors.Wrap(unmarshalErr, "failed to unmarshal vault JSON")
	}

	v.loaded = true
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
			bytes, _ := v.Crypter.Decrypt(cipherText)
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
			bytes, _ := v.Crypter.Encrypt([]byte(plainText))
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

// Base64Encode encodes vault strings as base64 only useful when exporting to k8s
func (v *Vault) Base64Encode() map[string]string {
	results := map[string]string{}
	for key, value := range v.Secrets {
		results[key] = base64.StdEncoding.EncodeToString([]byte(value))
	}
	return results
}

// Exists checks if the vault already exists on disk
func (v *Vault) Exists() (bool, error) {
	if _, err := os.Stat(v.Path()); err != nil {
		if os.IsNotExist(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// CreateIfNotExists checks if the vault already exists on disk
func (v *Vault) CreateIfNotExists() error {
	if exists, err := v.Exists(); !exists {

		if err != nil {
			return err
		}

		os.MkdirAll(filepath.Dir(v.Path()), os.ModePerm)

		if _, err := os.Create(v.Path()); err != nil {
			return errors.Wrap(err, "failed to create vault file")
		}

		v.isNew = true
	}
	return nil
}

// LoadOrCreate attempts to find a vault on the filesystem and creates it if it does not exist.
// Otherwise it will load it and populate its values with the contents
func (v *Vault) LoadOrCreate() error {

	if exists, _ := v.Exists(); exists {
		return v.Load()
	}

	v.isNew = true
	return v.Save()
}

// InitCrypter initialize the vaults crypter
func (v *Vault) InitCrypter() error {
	newCrypter, err := crypter.NewCrypter(&v.Project, &v.Location, &v.Keyring, &v.Key)
	if err != nil {
		return errors.Wrap(err, "failed to initialize vault crypter")
	}
	v.Crypter = newCrypter
	return nil
}

func (v *Vault) validate() error {
	if v.Project == "" {
		return errors.New("No Goolge Cloud `Project` was specified")
	}

	if v.Keyring == "" {
		return errors.New("No Google Cloud KMS `Keyring` was specified")
	}

	if v.Location == "" {
		return errors.New("No Google Cloud KMS Keyring `Location` was specified")
	}

	if v.Key == "" {
		return errors.New("No Google Cloud KMS `Key` was specified")
	}

	if _, encryptErr := v.Crypter.Encrypt([]byte("test")); encryptErr != nil {
		return errors.Wrap(encryptErr, "failed to verify cryptoKey settings")
	}

	return nil
}

// New returns a pointer to a Vault
// if the vault does not exist on disk it will be created
// if the vault was not newely created it will attempt to load and unmarshal it
func New(config Config) *Vault {
	return &Vault{
		Name:     config.Name,
		Project:  config.Project,
		Location: config.Location,
		Keyring:  config.Keyring,
		Key:      config.Key,
		Secrets:  map[string]string{},
	}
}

// EsureVaultLoaded ensure that the vault was successfully loaded
func EsureVaultLoaded(v *Vault) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if !v.loaded {
			return fmt.Errorf("This vault doesnt exist. You must first initialize it with `gvault init --vault %s`", v.Name)
		}
		return nil
	}
}
