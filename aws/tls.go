package aws

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"math/big"
	"strings"
	"time"
)

const (
	pemBlockTypeCertificate   = `CERTIFICATE`
	pemBlockTypeRsaPrivateKey = `RSA PRIVATE KEY`
)

var tlsX509CertificateSerialNumberLimit = new(big.Int).Lsh(big.NewInt(1), 128)

// tlsRsaPrivateKeyPem generates a RSA private key PEM string.
// Wrap with tlsPemEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: private_key_pem = "%[1]s"
func tlsRsaPrivateKeyPem(bits int) string {
	key, err := rsa.GenerateKey(rand.Reader, bits)

	if err != nil {
		panic(err)
	}

	block := &pem.Block{
		Bytes: x509.MarshalPKCS1PrivateKey(key),
		Type:  pemBlockTypeRsaPrivateKey,
	}

	return string(pem.EncodeToMemory(block))
}

// tlsRsaX509SelfSignedCertificatePem generates a x509 certificate PEM string.
// Wrap with tlsPemEscapeNewlines() to allow simple fmt.Sprintf()
// configurations such as: private_key_pem = "%[1]s"
func tlsRsaX509SelfSignedCertificatePem(keyPem, commonName string) string {
	keyBlock, _ := pem.Decode([]byte(keyPem))

	key, err := x509.ParsePKCS1PrivateKey(keyBlock.Bytes)

	if err != nil {
		panic(err)
	}

	serialNumber, err := rand.Int(rand.Reader, tlsX509CertificateSerialNumberLimit)

	if err != nil {
		panic(err)
	}

	certificate := &x509.Certificate{
		BasicConstraintsValid: true,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		NotAfter:              time.Now().Add(24 * time.Hour),
		NotBefore:             time.Now(),
		SerialNumber:          serialNumber,
		Subject: pkix.Name{
			CommonName:   commonName,
			Organization: []string{"ACME Examples, Inc"},
		},
	}

	certificateBytes, err := x509.CreateCertificate(rand.Reader, certificate, certificate, &key.PublicKey, key)

	if err != nil {
		panic(err)
	}

	certificateBlock := &pem.Block{
		Bytes: certificateBytes,
		Type:  pemBlockTypeCertificate,
	}

	return string(pem.EncodeToMemory(certificateBlock))
}

func tlsPemEscapeNewlines(pem string) string {
	return strings.ReplaceAll(pem, "\n", "\\n")
}
