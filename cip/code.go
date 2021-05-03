package cip

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
)

func GenerateRsaKeyPair() *rsa.PrivateKey {
	priKey, _ := rsa.GenerateKey(rand.Reader, 2048)
	return priKey
}

func GenerateRsaKeyPairPem() (priPem string, pubPem string, err error) {

	pKey := GenerateRsaKeyPair()
	pubKey := pKey.PublicKey

	priPem = PrivateKeyToPem(pKey)
	pubPem, err = PublicKeyToPem(&pubKey)
	if err != nil {
		pubPem = "ERROR_PUB_TO_PEM"
	}

	return priPem, pubPem, err
}

func PrivateKeyToPem(priKey *rsa.PrivateKey) string {
	keyBytes := x509.MarshalPKCS1PrivateKey(priKey)
	keyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: keyBytes,
		},
	)
	return string(keyPem)
}

func PemToPrivateKey(privateKeyPem string) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode([]byte(privateKeyPem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	priKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	return priKey, nil
}

func PublicKeyToPem(pubKey *rsa.PublicKey) (string, error) {
	bytes, err := x509.MarshalPKIXPublicKey(pubKey)
	if err != nil {
		return "", err
	}
	pem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: bytes,
		},
	)

	return string(pem), nil
}

func PemToPublicKey(pubPem string) (*rsa.PublicKey, error) {
	block, _ := pem.Decode([]byte(pubPem))
	if block == nil {
		return nil, errors.New("failed to parse PEM block containing the key")
	}

	pub, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	switch pub := pub.(type) {
	case *rsa.PublicKey:
		return pub, nil
	default:
		break // fall through
	}
	return nil, errors.New("Key type is not RSA")
}

func EncryptData(pubKey *rsa.PublicKey, data []byte) []byte {

	encryptedBytes, err := rsa.EncryptOAEP(
		sha256.New(),
		rand.Reader,
		pubKey,
		data,
		nil)

	if err != nil {
		panic(err)
	}

	fmt.Println("encrypted bytes: ", encryptedBytes)

	return encryptedBytes
}

func DecryptData(priKey *rsa.PrivateKey, data []byte) []byte {
	// The first argument is an optional random data generator (the rand.Reader we used before)
	// we can set this value as nil
	// The OAEPOptions in the end signify that we encrypted the data using OAEP, and that we used
	// SHA256 to hash the input.
	decryptedBytes, err := priKey.Decrypt(nil, data, &rsa.OAEPOptions{Hash: crypto.SHA256})
	if err != nil {
		panic(err)
	}

	// We get back the original information in the form of bytes, which we
	// the cast to a string and print
	fmt.Println("decrypted message: ", string(decryptedBytes))

	return decryptedBytes
}
