package storage

import (
	"fmt"

	"github.com/zalando/go-keyring"
)

const (
	KeyringService = "huggkey"
	KeyringUser    = "master-vault"
)

// StorePassphraseInKeyring saves the vault passphrase to the system keyring.
func StorePassphraseInKeyring(passphrase string) error {
	err := keyring.Set(KeyringService, KeyringUser, passphrase)
	if err != nil {
		return fmt.Errorf("failed to store passphrase in keyring: %w", err)
	}
	return nil
}

// GetPassphraseFromKeyring retrieves the vault passphrase from the system keyring.
func GetPassphraseFromKeyring() (string, error) {
	pass, err := keyring.Get(KeyringService, KeyringUser)
	if err != nil {
		return "", fmt.Errorf("failed to get passphrase from keyring: %w", err)
	}
	return pass, nil
}

// DeletePassphraseFromKeyring removes the passphrase from the system keyring.
func DeletePassphraseFromKeyring() error {
	err := keyring.Delete(KeyringService, KeyringUser)
	if err != nil {
		return fmt.Errorf("failed to delete passphrase from keyring: %w", err)
	}
	return nil
}
