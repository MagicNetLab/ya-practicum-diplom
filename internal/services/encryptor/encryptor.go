package encryptor

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/sha256"
	"encoding/base64"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/config"
	"github.com/MagicNetLab/ya-practicum-diplom/internal/logger"
)

// EncryptData - шифрует текст с помощью ключа
func EncryptData(data string) (string, error) {
	cnf := config.GetSecretConfig()
	secretKey, err := cnf.GetSecret()

	aesGCM, nonce, err := prepareEncryptor(secretKey)
	if err != nil {
		return "", err
	}

	encryptedData := aesGCM.Seal(nil, nonce, []byte(data), nil)
	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

// DecryptData - дешифрует текст с помощью ключа
func DecryptData(data string) (string, error) {
	cnf := config.GetSecretConfig()
	secretKey, err := cnf.GetSecret()

	d, err := base64.StdEncoding.DecodeString(data)
	if err != nil {
		return "", err
	}

	aesGCM, iv, err := prepareEncryptor(secretKey)
	if err != nil {
		return "", err
	}

	plainData, err := aesGCM.Open(nil, iv, d, nil)
	if err != nil {
		logger.Error("failed to decrypt data", logger.StrArg("error", err.Error()))
		return "", err
	}

	return string(plainData), nil
}

// EncryptPassword - шифрует пароль с помощью ключа
func EncryptPassword(password string) (string, error) {
	hash := sha256.Sum256([]byte(password))
	hashString := base64.StdEncoding.EncodeToString(hash[:])

	return hashString, nil
}

func prepareEncryptor(secretPassword string) (cipher.AEAD, []byte, error) {
	hash := sha256.New()

	_, err := hash.Write([]byte(secretPassword))
	if err != nil {
		logger.Error("failed to generate hash for password", logger.StrArg("error", err.Error()))
		return nil, nil, err
	}

	key := hash.Sum(nil)
	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error("failed to generate block for password", logger.StrArg("error", err.Error()))
		return nil, nil, err
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error("failed to generate AES GCM for password", logger.StrArg("error", err.Error()))
		return nil, nil, err
	}

	// создаем вектор инициализации из последних байт ключа
	iv := []byte(secretPassword[len(secretPassword)-aesGCM.NonceSize():])
	return aesGCM, iv, nil
}
