package main

import (
	"crypto/ecdh"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"
)

func generatePublicPrivateRsaKeys() {
	fmt.Println("cryptographic testing ")
	pvtKey, err := rsa.GenerateKey(rand.Reader, 2048)

	if err != nil {
		fmt.Println("failed to complete the operation", err)
		os.Exit(0)
	}

	der, err := x509.MarshalPKCS8PrivateKey(pvtKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling RSA private key: %s", err)
		return
	}

	fmt.Printf("%s", pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: der,
	}))
	der = x509.MarshalPKCS1PublicKey(&pvtKey.PublicKey)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error marshalling RSA private key: %s", err)
		return
	}

	fmt.Printf("%s", pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: der,
	}))

}

func generateEcdhPublicPrivatePair() {

	alicePrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating the private key: %s", err)
		return
	}
	alicePublicKey := alicePrivateKey.PublicKey()

	bobPrivateKey, err := ecdh.P256().GenerateKey(rand.Reader)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating the private key: %s", err)
		return
	}
	bobPublicKey := bobPrivateKey.PublicKey()

	sharedSecretAlice, err := alicePrivateKey.ECDH(bobPublicKey)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating the public key: %s", err)
		return
	}

	sharedSecretBob, err := bobPrivateKey.ECDH(alicePublicKey)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating the public key: %s", err)
		return
	}

	fmt.Printf("%x\n", sharedSecretAlice)
	fmt.Printf("%x\n", sharedSecretBob)

	if string(sharedSecretAlice) == string(sharedSecretBob) {
		fmt.Println("we got same secrets today")
	}

	// encrypt the text

	// cipher, err := aes.NewCipher(sharedSecretAlice)
	// if err != nil {
	// 	fmt.Fprintf(os.Stderr, "failed in deriving key: %s", err)
	// 	return
	// }

	// cipher.Encrypt("this is our secret")
}

func main() {
	generatePublicPrivateRsaKeys()
	generateEcdhPublicPrivatePair()
}
