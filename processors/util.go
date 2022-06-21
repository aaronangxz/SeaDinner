package processors

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"google.golang.org/protobuf/proto"
	"io"
	random "math/rand"
	"os"
	"reflect"
	"strings"
	"time"
	"unicode"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//MakeToken Decrypts the encrypted key with AES and appends the token prefix
func MakeToken(ctx context.Context, key string) string {
	if key == "" {
		log.Error(ctx, "Key is invalid: %v", key)
		return ""
	}

	decrypt := DecryptKey(key, os.Getenv("AES_KEY"))
	if len(decrypt) != 40 {
		log.Error(ctx, "Key length invalid | length: %v", len(decrypt))
		return ""
	}
	return fmt.Sprint(common.Config.Prefix.TokenPrefix, decrypt)
}

//MakeURL Returns the corresponding API URL appended with Id, if needed
func MakeURL(opt int, id *int64) string {
	prefix := common.Config.Prefix.URLPrefix
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
func OutputResults(ctx context.Context, resultMap map[int64]int64, service string) {
	var (
		passed int
	)
	for _, m := range resultMap {
		if m == int64(sea_dinner.OrderStatus_ORDER_STATUS_OK) {
			passed++
		}
	}

	log.Info(ctx, fmt.Sprintf("\n%v\n*************************\nTotal Order: %v\nTotal Success: %v\nTotal Failures: %v\n*************************",
		service,
		len(resultMap),
		passed,
		len(resultMap)-passed))
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
func RandomFood(ctx context.Context, m map[string]string) string {
	var s []string
	count := 0
	for k, v := range m {
		if k == "RAND" ||
			k == "-1" ||
			strings.Contains(v, "Vegetarian") ||
			strings.Contains(v, "vegetarian") {
			count++
			continue
		}
		s = append(s, k)
	}

	if count >= len(m) {
		log.Error(ctx, "RandomFood | Count exceeds Index")
		return ""
	}

	r := random.New(random.NewSource(time.Now().UnixNano()))
	gen := int64(r.Intn(len(m) - count))

	if gen >= int64(len(s)) {
		log.Error(ctx, "RandomFood | Index exceeds len")
		return ""
	}
	log.Info(ctx, "RandomFood | result: %v", s[gen])
	return s[gen]
}

//CompareSliceStruct Compares two struct slices and outputs the difference.
func CompareSliceStruct(ctx context.Context, a interface{}, b interface{}) bool {
	same := true

	if reflect.TypeOf(a).Kind() != reflect.TypeOf(b).Kind() {
		log.Error(ctx, "CompareSliceStruct | Slices must be the same type.")
		return false
	}

	switch reflect.TypeOf(a).Kind() {
	case reflect.Slice:
		first := reflect.Indirect(reflect.ValueOf(a))
		second := reflect.Indirect(reflect.ValueOf(b))

		if first.Len() == 0 || second.Len() == 0 {
			return false
		}

		for i, j := 0, 0; i < first.Len() && j < second.Len(); i, j = i+1, j+1 {
			if !reflect.DeepEqual(first.Index(i).Interface(), second.Index(j).Interface()) {
				log.Warn(ctx, "CompareSliceStruct | diff | \nfirst: %v | \nsecond: %v", first.Index(i).Interface(), second.Index(j).Interface())
				same = false
			}
		}
	default:
		log.Error(ctx, "CompareSliceStruct | Only slice is supported.")
		return false
	}
	return same
}

func ConvertFoodToFoodMapping(ctx context.Context, food []*sea_dinner.Food) []*sea_dinner.FoodMapping {
	var (
		mappings []*sea_dinner.FoodMapping
	)

	for _, f := range food {
		mappings = append(mappings, &sea_dinner.FoodMapping{
			FoodId:    f.Id,
			FoodCode:  f.Code,
			FoodName:  f.Name,
			FoodImage: f.ImageUrl,
		})
	}
	log.Info(ctx, "ConvertFoodToFoodMapping | Success | size:%v", len(mappings))
	return mappings
}

func ConvertFoodToFoodMappings(ctx context.Context, food []*sea_dinner.Food) *sea_dinner.FoodMappings {
	return &sea_dinner.FoodMappings{FoodMapping: ConvertFoodToFoodMapping(ctx, food)}
}

func ConvertFoodToFoodMappingByYearAndWeek(ctx context.Context, food []*sea_dinner.Food) *sea_dinner.FoodMappingByYearAndWeek {
	year, week := ConvertTimeStampWeekOfYear(time.Now().Unix())
	mapping := ConvertFoodToFoodMappings(ctx, food)
	foodBytes, err := proto.Marshal(mapping)
	if err != nil {
		log.Error(ctx, "ConvertFoodToFoodMapping | Failed | %v", err.Error())
		return nil
	}
	log.Info(ctx, "ConvertFoodToFoodMapping | Success")
	return &sea_dinner.FoodMappingByYearAndWeek{
		Year:        proto.Int64(year),
		Week:        proto.Int64(week),
		FoodMapping: foodBytes,
	}
}

//MakeMenuMap Returns food_id:food_name mapping of current menu
func MakeMenuMap(ctx context.Context) map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := App.StartTransaction("make_menu_map")
	defer txn.End()

	menuMap := make(map[string]string)
	menu := GetMenuUsingCache(ctx, key)
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "*NOTHING*" // to be renamed
	menuMap["RAND"] = "RAND"
	return menuMap
}
