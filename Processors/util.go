package Processors

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	random "math/rand"
	"os"
	"time"
	"unicode"

	"github.com/aaronangxz/SeaDinner/Common"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//MakeToken Decrypts the encrypted key with AES and appends the token prefix
func MakeToken(key string) string {
	if key == "" {
		log.Println("Key is invalid:", key)
		return ""
	}

	decrypt := DecryptKey(key, os.Getenv("AES_KEY"))
	if len(decrypt) != 40 {
		log.Printf("Key length invalid | length: %v", len(decrypt))
		return ""
	}
	return fmt.Sprint(Common.Config.Prefix.TokenPrefix, decrypt)
}

//MakeURL Returns the corresponding API URL appended with Id, if needed
func MakeURL(opt int, id *int64) string {
	prefix := Common.Config.Prefix.UrlPrefix
	switch opt {
	case int(sea_dinner.URLType_URL_CURRENT):
		return fmt.Sprint(prefix, "/api/current")
	case int(sea_dinner.URLType_URL_MENU):
		if id == nil {
			return ""
		}
		return fmt.Sprint(prefix, "/api/menu/", *id)
	case int(sea_dinner.URLType_URL_ORDER):
		if id == nil {
			return ""
		}
		return fmt.Sprint(prefix, "/api/order/", *id)
	}
	return ""
}

//OutputResults Prints out the total success and failure cases
func OutputResults(resultMap map[int64]int64, service string) {
	var (
		passed int
	)
	for _, m := range resultMap {
		if m == int64(sea_dinner.OrderStatus_ORDER_STATUS_OK) {
			passed++
		}
	}

	fmt.Println(service)
	fmt.Println("*************************")
	fmt.Println("Total Order: ", len(resultMap))
	fmt.Println("Total Success: ", passed)
	fmt.Println("Total Failures: ", len(resultMap)-passed)
	fmt.Println("*************************")
}

//IsNotNumber Verifiy whether a string contains non-numeric characters
func IsNotNumber(a string) bool {
	if a == "" {
		return true
	}

	// Catch "-1", hacky, to be re-done
	if a == "-1" {
		return false
	}

	for _, char := range a {
		if unicode.IsSymbol(char) {
			return true
		}
	}
	for _, char := range a {
		if !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

//EncryptKey Encrypts a string using the provided encryption key
func EncryptKey(stringToEncrypt string, keyString string) (encryptedString string) {

	//Since the key is in string, we need to convert decode it to bytes
	key, _ := hex.DecodeString(keyString)
	plaintext := []byte(stringToEncrypt)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM - https://en.wikipedia.org/wiki/Galois/Counter_Mode
	//https://golang.org/pkg/crypto/cipher/#NewGCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Create a nonce. Nonce should be from GCM
	nonce := make([]byte, aesGCM.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panic(err.Error())
	}

	//Encrypt the data using aesGCM.Seal
	//Since we don't want to save the nonce somewhere else in this case, we add it as a prefix to the encrypted data. The first nonce argument in Seal is the prefix.
	ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)
	return fmt.Sprintf("%x", ciphertext)
}

//DecryptKey Decrypts a string using the provided encryption key
func DecryptKey(encryptedString string, keyString string) (decryptedString string) {

	key, _ := hex.DecodeString(keyString)
	enc, _ := hex.DecodeString(encryptedString)

	//Create a new Cipher Block from the key
	block, err := aes.NewCipher(key)
	if err != nil {
		panic(err.Error())
	}

	//Create a new GCM
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panic(err.Error())
	}

	//Get the nonce size
	nonceSize := aesGCM.NonceSize()

	//Extract the nonce from the encrypted data
	nonce, ciphertext := enc[:nonceSize], enc[nonceSize:]

	//Decrypt the data
	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panic(err.Error())
	}
	return string(plaintext)
}

//MakeKey Creates a random AES key
func MakeKey() string {
	bytes := make([]byte, 16) //generate a random 32 byte key for AES-256
	if _, err := rand.Read(bytes); err != nil {
		panic(err.Error())
	}
	return hex.EncodeToString(bytes) //encode key in bytes to string for saving
}

//RandomFood Returns a random element in the provided menu map, excluding RAND and -1
func RandomFood(m map[string]string) string {
	s := []string{}

	for k := range m {
		if k == "RAND" || k == "-1" {
			continue
		}
		s = append(s, k)
	}

	r := random.New(random.NewSource(time.Now().UnixNano()))
	gen := int64(r.Intn(len(m) - 3))
	log.Println("RandomFood | result:", s[gen])
	return s[gen]
}
