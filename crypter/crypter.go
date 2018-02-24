package crypter

import (
	"encoding/base64"
	"fmt"
	"log"

	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	cloudkms "google.golang.org/api/cloudkms/v1"
)

// Encrypt encrypts a secret using Google KMS
func Encrypt(projectID, keyRing, cryptoKey string, plainText []byte) (string, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, "global", keyRing, cryptoKey)

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.
		Encrypt(parentName, &cloudkms.EncryptRequest{
			Plaintext: base64.StdEncoding.EncodeToString(plainText),
		}).Do()
	if err != nil {
		return "", err
	}

	return resp.Ciphertext, nil
}

// Decrypt decrypts a secret using Google KMS
func Decrypt(projectID, keyRing, cryptoKey string, cipherText string) ([]byte, error) {
	ctx := context.Background()
	client, err := google.DefaultClient(ctx, cloudkms.CloudPlatformScope)
	if err != nil {
		log.Fatal(err)
	}

	cloudkmsService, err := cloudkms.New(client)
	if err != nil {
		log.Fatal(err)
	}

	parentName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s",
		projectID, "global", keyRing, cryptoKey)

	resp, err := cloudkmsService.Projects.Locations.KeyRings.CryptoKeys.
		Decrypt(parentName, &cloudkms.DecryptRequest{
			Ciphertext: cipherText,
		}).Do()
	if err != nil {
		return nil, err
	}
	return base64.StdEncoding.DecodeString(resp.Plaintext)
}
