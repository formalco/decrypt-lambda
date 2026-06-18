// Package provider routes a JWE "kid" key URI to the KMS that can unwrap the
// record's content key. A key URI is "<scheme>://<keyID>"; the scheme selects a
// registered provider and the keyID is passed through. Add a provider by
// implementing Provider and calling Register from its package init.
package provider

import (
	"context"
	"fmt"
	"strings"
)

// UnwrapFunc unwraps the RSA-OAEP-256 wrapped content key (CEK), returning the raw
// key bytes for go-jose to run the A256GCM content decryption.
type UnwrapFunc func(encryptedKey []byte) ([]byte, error)

// Provider unwraps a wrapped CEK with a specific KMS.
type Provider interface {
	Unwrapper(ctx context.Context, keyID string) (UnwrapFunc, error)
}

var registry = map[string]Provider{}

// Register binds a key-URI scheme (e.g. "aws-kms") to a provider. It panics on a
// duplicate scheme so a misconfigured build fails at startup, not at decrypt time.
func Register(scheme string, p Provider) {
	if _, dup := registry[scheme]; dup {
		panic(fmt.Sprintf("provider: scheme %q already registered", scheme))
	}
	registry[scheme] = p
}

func Resolve(ctx context.Context, keyURI string) (UnwrapFunc, error) {
	scheme, keyID, ok := strings.Cut(keyURI, "://")
	if !ok {
		return nil, fmt.Errorf("invalid KMS key URI %q: want <scheme>://<key-id>", keyURI)
	}
	p, ok := registry[scheme]
	if !ok {
		return nil, fmt.Errorf("unsupported KMS key URI scheme %q", scheme)
	}
	return p.Unwrapper(ctx, keyID)
}
