package Common

import (
	"bytes"
	"image"
	"io/ioutil"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/liyue201/goqr"
)

func GetTGToken() string {
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Config.Adhoc {
		log.Println("Running Test Telegram Bot Instance")
		return os.Getenv("TELEGRAM_TEST_APITOKEN")
	}
	return os.Getenv("TELEGRAM_APITOKEN")
}

func IsInGrayScale(userId int64) bool {
	return userId%100 >= Config.GrayScale.Percentage
}

func DecodeQR() (string, error) {
	//absPath, _ := filepath.Abs("../Common/resource/DinnerQR.jpg")
	filepath := "../resource/DinnerQR.jpg"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		filepath = "../Common/resource/DinnerQR.jpg"
	}
	qr, err := recognizeFile(filepath)
	if err != nil {
		return "DecodeQR | Failed to recognize file.", err
	}

	if len(qr) == 0 {
		return "DecodeQR | Unable to find QR URL.", nil
	}

	return string(qr[0].Payload), nil
}

func recognizeFile(path string) ([]*goqr.QRData, error) {
	log.Printf("recognize file: %v\n", path)
	imgdata, err := ioutil.ReadFile(path)
	if err != nil {
		log.Printf("%v\n", err)
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		log.Printf("image.Decode error: %v\n", err)
		return nil, err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Printf("Recognize failed: %v\n", err)
		return nil, err
	}
	return qrCodes, nil
}
