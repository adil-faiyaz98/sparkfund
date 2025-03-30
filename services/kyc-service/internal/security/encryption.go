package security

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
)

// EncryptionConfig holds encryption configuration
type EncryptionConfig struct {
	AESKey    []byte
	RSAPublic string
	RSAPrivate string
}

// DefaultEncryptionConfig returns default encryption configuration
func DefaultEncryptionConfig() EncryptionConfig {
	return EncryptionConfig{
		AESKey:    []byte("your-32-byte-aes-key-here"), // Change this in production
		RSAPublic: "your-rsa-public-key",               // Change this in production
		RSAPrivate: "your-rsa-private-key",             // Change this in production
	}
}

// EncryptAES encrypts data using AES-GCM
func EncryptAES(plaintext []byte, config EncryptionConfig) ([]byte, error) {
	block, err := aes.NewCipher(config.AESKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)
	return ciphertext, nil
}

// DecryptAES decrypts data using AES-GCM
func DecryptAES(ciphertext []byte, config EncryptionConfig) ([]byte, error) {
	block, err := aes.NewCipher(config.AESKey)
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	if len(ciphertext) < gcm.NonceSize() {
		return nil, fmt.Errorf("ciphertext too short")
	}

	nonce := ciphertext[:gcm.NonceSize()]
	ciphertext = ciphertext[gcm.NonceSize():]

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

// EncryptRSA encrypts data using RSA
func EncryptRSA(plaintext []byte, config EncryptionConfig) ([]byte, error) {
	block, _ := pem.Decode([]byte(config.RSAPublic))
	if block == nil {
		return nil, fmt.Errorf("failed to parse RSA public key")
	}

	pub, err := x509.ParsePKCS1PublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	label := []byte("")
	encrypted, err := rsa.EncryptOAEP(hash, rand.Reader, pub, plaintext, label)
	if err != nil {
		return nil, err
	}

	return encrypted, nil
}

// DecryptRSA decrypts data using RSA
func DecryptRSA(ciphertext []byte, config EncryptionConfig) ([]byte, error) {
	block, _ := pem.Decode([]byte(config.RSAPrivate))
	if block == nil {
		return nil, fmt.Errorf("failed to parse RSA private key")
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	hash := sha256.New()
	label := []byte("")
	decrypted, err := rsa.DecryptOAEP(hash, rand.Reader, priv, ciphertext, label)
	if err != nil {
		return nil, err
	}

	return decrypted, nil
}

// HashPassword hashes a password using SHA-256
func HashPassword(password string) string {
	hash := sha256.New()
	hash.Write([]byte(password))
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// VerifyPassword verifies a password against its hash
func VerifyPassword(password, hash string) bool {
	passwordHash := HashPassword(password)
	return passwordHash == hash
}

// MaskSensitiveData masks sensitive data in logs
func MaskSensitiveData(data string) string {
	if len(data) <= 4 {
		return "****"
	}
	return data[:2] + "****" + data[len(data)-2:]
}

// GenerateSecureToken generates a secure random token
func GenerateSecureToken(length int) (string, error) {
	bytes := make([]byte, length)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return base64.URLEncoding.EncodeToString(bytes), nil
}

// HashFile generates a SHA-256 hash of a file
func HashFile(data []byte) string {
	hash := sha256.New()
	hash.Write(data)
	return base64.URLEncoding.EncodeToString(hash.Sum(nil))
}

// GenerateKeyPair generates a new RSA key pair
func GenerateKeyPair(bits int) (string, string, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return "", "", err
	}

	publicKey := &privateKey.PublicKey

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	publicKeyBytes := x509.MarshalPKCS1PublicKey(publicKey)

	privateKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	publicKeyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	})

	return string(publicKeyPEM), string(privateKeyPEM), nil
} 