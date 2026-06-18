//go:build integration

// Exercises the full decrypt path (provider registry -> awskms -> kms.Decrypt)
// against a real KMS. Point it at a local LocalStack KMS or any KMS the default
// AWS credential chain can reach:
//
//	go test -tags integration -run TestDecryptViaKMS ./...
//
// The endpoint defaults to http://localhost:4566; override with TEST_KMS_ENDPOINT.
package main

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/kms"
)

func TestDecryptViaKMS(t *testing.T) {
	endpoint := os.Getenv("TEST_KMS_ENDPOINT")
	if endpoint == "" {
		endpoint = "http://localhost:4566"
	}
	if os.Getenv("AWS_ACCESS_KEY_ID") == "" {
		t.Setenv("AWS_ACCESS_KEY_ID", "test")
		t.Setenv("AWS_SECRET_ACCESS_KEY", "test")
	}
	// awskms reads DEV_AWS_ENDPOINT to target a non-AWS KMS.
	t.Setenv("DEV_AWS_ENDPOINT", endpoint)

	ctx := context.Background()
	svc := newKMSClient(t, endpoint)

	createOut, err := svc.CreateKeyWithContext(ctx, &kms.CreateKeyInput{
		KeySpec:  aws.String(kms.KeySpecRsa2048),
		KeyUsage: aws.String(kms.KeyUsageTypeEncryptDecrypt),
	})
	if err != nil {
		t.Fatalf("CreateKey (is LocalStack KMS up at %s?): %v", endpoint, err)
	}
	keyARN := aws.StringValue(createOut.KeyMetadata.Arn)

	pub := fetchPublicKey(t, svc, keyARN)

	const plaintext = "super secret value"
	jwe := sealJWE(t, plaintext, pub, "aws-kms://"+keyARN)

	got, err := decryptValue(ctx, jwe)
	if err != nil {
		t.Fatalf("decryptValue: %v", err)
	}
	if got != plaintext {
		t.Fatalf("got %q, want %q", got, plaintext)
	}
}

// decryptValue runs the same parse + provider resolve + decrypt the Lambda handler does.
func decryptValue(ctx context.Context, value string) (string, error) {
	obj, err := parseJWE(value)
	if err != nil {
		return "", err
	}
	return decryptJWE(ctx, obj)
}

func newKMSClient(t *testing.T, endpoint string) *kms.KMS {
	t.Helper()
	sess, err := session.NewSession(&aws.Config{
		Region:   aws.String("us-east-1"),
		Endpoint: aws.String(endpoint),
	})
	if err != nil {
		t.Fatalf("new session: %v", err)
	}
	return kms.New(sess)
}

func fetchPublicKey(t *testing.T, svc *kms.KMS, keyID string) *rsa.PublicKey {
	t.Helper()
	out, err := svc.GetPublicKey(&kms.GetPublicKeyInput{KeyId: aws.String(keyID)})
	if err != nil {
		t.Fatalf("GetPublicKey: %v", err)
	}
	parsed, err := x509.ParsePKIXPublicKey(out.PublicKey)
	if err != nil {
		t.Fatalf("parse public key: %v", err)
	}
	pub, ok := parsed.(*rsa.PublicKey)
	if !ok {
		t.Fatalf("public key is %T, want *rsa.PublicKey", parsed)
	}
	return pub
}
