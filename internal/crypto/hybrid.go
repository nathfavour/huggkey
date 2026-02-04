package crypto

import (
	"crypto/ed25519"
	"fmt"
	"io"

	"github.com/cloudflare/circl/sign/eddilithium2"
)

// HybridPrivateKey represents the combined Ed25519 and Dilithium private key.
type HybridPrivateKey struct {
	priv *eddilithium2.PrivateKey
}

// HybridPublicKey represents the combined Ed25519 and Dilithium public key.
type HybridPublicKey struct {
	pub *eddilithium2.PublicKey
}

// GenerateKeyPair creates a new hybrid Ed25519 + Dilithium2 key pair.
// If r is nil, crypto/rand.Reader is used.
func GenerateKeyPair(r io.Reader) (*HybridPrivateKey, *HybridPublicKey, error) {
	pk, sk, err := eddilithium2.GenerateKey(r)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate hybrid keypair: %w", err)
	}

	return &HybridPrivateKey{priv: sk}, &HybridPublicKey{pub: pk}, nil
}

// Sign creates a hybrid signature for the given message.
func (sk *HybridPrivateKey) Sign(message []byte) ([]byte, error) {
	return sk.priv.Sign(nil, message, nil)
}

// Verify checks the hybrid signature against the message and public key.
func (pk *HybridPublicKey) Verify(message, signature []byte) bool {
	return eddilithium2.Verify(pk.pub, message, signature)
}

// Ed25519PublicKey extracts the classical Ed25519 public key.
func (pk *HybridPublicKey) Ed25519PublicKey() ed25519.PublicKey {
	// eddilithium2.PublicKey is [Ed25519PK || DilithiumPK]
	// Ed25519 PK is 32 bytes
	b := pk.pub.Bytes()
	return ed25519.PublicKey(b[:32])
}
