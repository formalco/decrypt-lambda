package main

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"errors"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/kms"
)

type FormalEncryptedData struct {
	EncryptedData string `json:"encryptedData"`
	EncryptedKey  string `json:"encryptedKey"`
	KmsKeyId      string `json:"kmsKeyId"`
	KmsKeyRegion  string `json:"kmsKeyRegion"`
}

func decryptString(encrypted string, key []byte) (string, error) {
	enc, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}
	nonceSize := aesGCM.NonceSize()
	if len(enc) < nonceSize {
		return "", errors.New("data not encrypted or not encrypted with a nonce")
	}
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}
	return string(plaintext), nil
}

func decryptDataKey(ctx context.Context, kmsKeyRegion, kmsKeyId string, encryptedKey []byte) ([]byte, error) {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(kmsKeyRegion))
	if err != nil {
		return nil, err
	}
	if os.Getenv("DEV_AWS_ENDPOINT") != "" {
		cfg.BaseEndpoint = aws.String(os.Getenv("DEV_AWS_ENDPOINT"))
	}

	svc := kms.NewFromConfig(cfg)
	input := &kms.DecryptInput{
		CiphertextBlob: encryptedKey,
		KeyId:          aws.String(string(kmsKeyId)),
	}

	result, err := svc.Decrypt(ctx, input)
	if err != nil {
		return nil, err
	}

	return result.Plaintext, nil
}

func parseFormalEncryptedData(encryptedData string) (FormalEncryptedData, error) {
	split := strings.Split(encryptedData, ":")
	if len(split) < 5 {
		return FormalEncryptedData{}, errors.New("Invalid input. Expected format: <encrypted data>:<encrypted key>:<kms key id>:<kms key region>")
	}

	encryptedKey, err := base64.StdEncoding.DecodeString(split[2])
	if err != nil {
		return FormalEncryptedData{}, err
	}

	kmsKeyId, err := base64.StdEncoding.DecodeString(split[3])
	if err != nil {
		return FormalEncryptedData{}, err
	}

	kmsKeyRegion, err := base64.StdEncoding.DecodeString(split[4])
	if err != nil {
		return FormalEncryptedData{}, err
	}

	return FormalEncryptedData{
		EncryptedData: split[1],
		EncryptedKey:  string(encryptedKey),
		KmsKeyId:      string(kmsKeyId),
		KmsKeyRegion:  string(kmsKeyRegion),
	}, nil
}
