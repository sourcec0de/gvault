package vault

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/sourcec0de/gvault/crypter"
)

// Vault a vault that stores in a json format
type Vault struct {
	filePath string
	crypter  *crypter.Crypter
	Version  string            `json:"version"`
	Secrets  map[string]string `json:"secrets"`
	isNew    bool
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
func (v *Vault) GetSecret(key string) string {
	return v.Secrets[key]
}

// Save writes the vault to it's storage location
func (v *Vault) Save() error {
	bytes, jsonSaveErr := json.MarshalIndent(v, "", "  ")

	if jsonSaveErr != nil {
		return jsonSaveErr
	}

	if ioWriteErr := ioutil.WriteFile(v.filePath, bytes, os.ModePerm); ioWriteErr != nil {
		return ioWriteErr
	}

	return nil
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
