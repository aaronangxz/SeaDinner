package Common

import (
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"log"
	"os"

	_ "image/jpeg"
	_ "image/png"

	"github.com/aaronangxz/SeaDinner/Log"
	"github.com/liyue201/goqr"
)

var (
	Ctx = context.TODO()
)

func GetTGToken(ctx context.Context) string {
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Config.Adhoc {
		Log.Info(ctx, "Running Test Telegram Bot Instance")
		// log.Println("Running Test Telegram Bot Instance")
		return os.Getenv("TELEGRAM_TEST_APITOKEN")
	}
	return os.Getenv("TELEGRAM_APITOKEN")
}

func IsInGrayScale(userId int64) bool {
	return userId%100 >= Config.GrayScale.Percentage
}

func DecodeQR() (string, error) {
	filepath := "../Common/resource/DinnerQR.jpg"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		filepath = "Common/resource/DinnerQR.jpg"
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
	Log.Info(Ctx, "recognize file: %v", path)
	// log.Printf("recognize file: %v\n", path)
	imgdata, err := ioutil.ReadFile(path)
	if err != nil {
		Log.Error(Ctx, "%v\n", err)
		log.Printf("%v\n", err)
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgdata))
	if err != nil {
		Log.Error(Ctx, "image.Decode error: %v", err)
		// log.Printf("image.Decode error: %v\n", err)
		return nil, err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		Log.Error(Ctx, "Recognize failed: %v", err)
		// log.Printf("Recognize failed: %v\n", err)
		return nil, err
	}
	return qrCodes, nil
}
