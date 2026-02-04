package storage

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"golang.org/x/crypto/argon2"
)

const (
	VaultVersion = 1
	SaltSize     = 16
	NonceSize    = 12
)

// Vault represents the encrypted storage for the HuggID.
type Vault struct {
	Version  int    `json:"version"`
	Salt     []byte `json:"salt"`
	Nonce    []byte `json:"nonce"`
	Cipher   []byte `json:"cipher"`
}

// SaveVault encrypts and saves the mnemonic to the vault file.
func SaveVault(path, password, mnemonic string) error {
	salt := make([]byte, SaltSize)
	if _, err := io.ReadFull(rand.Reader, salt); err != nil {
		return err
	}

	// Derive key using Argon2id
	// Parameters: time=1, memory=64MB, threads=4, keyLen=32
	key := argon2.IDKey([]byte(password), salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return err
	}

	cipherText := gcm.Seal(nil, nonce, []byte(mnemonic), nil)

	vault := Vault{
		Version: VaultVersion,
		Salt:    salt,
		Nonce:   nonce,
		Cipher:  cipherText,
	}

	data, err := json.MarshalIndent(vault, "", "  ")
	if err != nil {
		return err
	}

	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}

	return os.WriteFile(path, data, 0600)
}

// LoadVault reads and decrypts the mnemonic from the vault file.
func LoadVault(path, password string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}

	var vault Vault
	if err := json.Unmarshal(data, &vault); err != nil {
		return "", err
	}

	key := argon2.IDKey([]byte(password), vault.Salt, 1, 64*1024, 4, 32)

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	plainText, err := gcm.Open(nil, vault.Nonce, vault.Cipher, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt vault: %w", err)
	}

	return string(plainText), nil
}

// GetDefaultVaultPath returns the standard ~/.config/hugg/identity.vault path.
func GetDefaultVaultPath() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return filepath.Join(home, ".config", "hugg", "identity.vault"), nil
}
