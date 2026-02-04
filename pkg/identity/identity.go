package identity

import (
	"bytes"
	"fmt"

	"github.com/mr-tron/base58"
	"github.com/nathfavour/huggkey/internal/crypto"
	"github.com/tyler-smith/go-bip39"
)

// Identity represents a HuggID and its associated keys.
type Identity struct {
	HuggID     string
	PrivateKey *crypto.HybridPrivateKey
	PublicKey  *crypto.HybridPublicKey
}

// CreateNewIdentity generates a new mnemonic and derived identity.
func CreateNewIdentity() (string, *Identity, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", nil, err
	}
	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", nil, err
	}

	id, err := FromMnemonic(mnemonic)
	if err != nil {
		return "", nil, err
	}

	return mnemonic, id, nil
}

// FromMnemonic derives an identity from a 24-word seed phrase.
func FromMnemonic(mnemonic string) (*Identity, error) {
	if !bip39.IsMnemonicValid(mnemonic) {
		return nil, fmt.Errorf("invalid mnemonic")
	}
	seed := bip39.NewSeed(mnemonic, "") // No passphrase for now as per minimal design
	
	// We use the seed to create a deterministic reader for key generation.
	reader := bytes.NewReader(seed)
	
	priv, pub, err := crypto.GenerateKeyPair(reader)
	if err != nil {
		return nil, err
	}

	return &Identity{
		HuggID:     FormatHuggID(pub),
		PrivateKey: priv,
		PublicKey:  pub,
	}, nil
}

// IsValidMnemonic checks if the mnemonic is valid.
func IsValidMnemonic(mnemonic string) bool {
	return bip39.IsMnemonicValid(mnemonic)
}

// FormatHuggID creates the hugg:v1:<base58> string.
func FormatHuggID(pub *crypto.HybridPublicKey) string {
	pkBytes := pub.Ed25519PublicKey()
	return fmt.Sprintf("hugg:v1:%s", base58.Encode(pkBytes))
}

// VerifyHuggID checks if a HuggID matches a public key.
func VerifyHuggID(huggID string, pub *crypto.HybridPublicKey) bool {
	return huggID == FormatHuggID(pub)
}
