package provider

import (
	"context"
	"testing"
)

func TestRegisterAndResolve(t *testing.T) {
	want := []byte("cek")
	Register("test-scheme", stubProvider{key: want})

	unwrap, err := Resolve(context.Background(), "test-scheme://some-key-id")
	if err != nil {
		t.Fatalf("Resolve: %v", err)
	}
	got, err := unwrap(nil)
	if err != nil {
		t.Fatalf("unwrap: %v", err)
	}
	if string(got) != string(want) {
		t.Fatalf("got %q, want %q", got, want)
	}
}

func TestResolveRejectsUnsupportedScheme(t *testing.T) {
	if _, err := Resolve(context.Background(), "gcp-kms://projects/p/cryptoKeys/k"); err == nil {
		t.Fatal("expected error for unsupported provider")
	}
}

func TestResolveRejectsMalformedURI(t *testing.T) {
	if _, err := Resolve(context.Background(), "not-a-uri"); err == nil {
		t.Fatal("expected error for malformed key URI")
	}
}

type stubProvider struct {
	key []byte
}

func (s stubProvider) Unwrapper(_ context.Context, _ string) (UnwrapFunc, error) {
	return func(_ []byte) ([]byte, error) { return s.key, nil }, nil
}
