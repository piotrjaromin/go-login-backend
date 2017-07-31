package accounts

import (
	"github.com/op/go-logging"
	"io"
	"golang.org/x/crypto/scrypt"
	cryptoRand "crypto/rand"
	"fmt"
)

const (
	PW_SALT_BYTES = 32
	PW_HASH_BYTES = 64
)

type Encrypt struct {
	Hash func(pass Password) (Password, string)
	Validate func(pass Password, hashToCompare Password, salt string) bool
}


func CreateEncrypt() Encrypt {

	var log = logging.MustGetLogger("[Encrypt]")

	cryptToHash := func(pass Password, salt string) (Password, error) {
		hash, err := scrypt.Key([]byte(pass), []byte(salt), 1<<14, 8, 1, PW_HASH_BYTES)
		if err != nil {
			return "", err
		}

		return Password(fmt.Sprintf("%x", hash)), nil
	}

	hash := func(pass Password) (Password, string) {
		saltInBytes := make([]byte, PW_SALT_BYTES)

		_, err := io.ReadFull(cryptoRand.Reader, saltInBytes)
		if err != nil {
			log.Error(err)
		}

		salt := fmt.Sprintf("%x", saltInBytes)
		hash, cryptErr := cryptToHash(pass, salt)
		if cryptErr != nil {
			log.Error(cryptErr)
		}

		return hash, salt
	}

	validate := func(pass Password, hashToCompare Password, salt string) bool {

		hash, cryptErr := cryptToHash(pass, salt)
		if cryptErr != nil {
			log.Error("Encryption error. Details: ", cryptErr)
			return false
		}

		return hash == hashToCompare
	}

	return Encrypt{
		Hash: hash,
		Validate: validate,
	}
}
