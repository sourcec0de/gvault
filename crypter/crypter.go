package crypter

import (
	"encoding/base64"
	"fmt"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Crypter encrypt and descrypt secrets using KMS
type Crypter struct {
	Project  *string
	Location *string
	Keyring  *string
	Key      *string
	kms      *cloudkms.Service
}

// NewCrypter creates a new Crypter instance
func NewCrypter(project, location, keyring, key *string) (*Crypter, error) {
	crypter := &Crypter{
		Project:  project,
		Location: location,
		Keyring:  keyring,
		Key:      key,
	}

	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		return crypter, err
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		return crypter, err
	}

	crypter.kms = cloudkmsService
	return crypter, nil
}

// KmsKeyName name of the kms key
func (c *Crypter) KmsKeyName() string {
	return fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		*c.Project, *c.Location, *c.Keyring, *c.Key)
}

// Encrypt encrypts a secret using Google KMS
func (c *Crypter) Encrypt(plainText []byte) (string, error) {
	resp, err := c.kms.Projects.Locations.KeyRings.CryptoKeys.
		Encrypt(c.KmsKeyName(), &cloudkms.EncryptRequest{
			Plaintext: base64.StdEncoding.EncodeToString(plainText),
		}).Do()

	if err != nil {
		return "", err
	}

	return resp.Ciphertext, nil
}

// Decrypt decrypts a secret using Google KMS
func (c *Crypter) Decrypt(cipherText string) ([]byte, error) {
	resp, err := c.kms.Projects.Locations.KeyRings.CryptoKeys.
		Decrypt(c.KmsKeyName(), &cloudkms.DecryptRequest{
			Ciphertext: cipherText,
		}).Do()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}
