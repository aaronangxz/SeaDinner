package common

import (
	"bytes"
	"context"
	"image"
	"io/ioutil"
	"os"
	//Required for QR decode
	_ "image/jpeg"
	_ "image/png"

	"github.com/aaronangxz/SeaDinner/log"
	"github.com/liyue201/goqr"
)

//GetTGToken Get default Telegram token
func GetTGToken(ctx context.Context) string {
	if os.Getenv("TEST_DEPLOY") == "TRUE" || Config.Adhoc {
		log.Info(ctx, "Running Test Telegram handlers Instance")
		return os.Getenv("TELEGRAM_TEST_APITOKEN")
	}
	return os.Getenv("TELEGRAM_APITOKEN")
}

//IsInGrayScale Verify userID is within grayscale percentage
func IsInGrayScale(userID int64) bool {
	return userID%100 >= Config.GrayScale.Percentage
}

//DecodeQR decodes QR image
func DecodeQR() (string, error) {
	filepath := "../common/resource/DinnerQR.jpg"
	if os.Getenv("HEROKU_DEPLOY") == "TRUE" || os.Getenv("TEST_DEPLOY") == "TRUE" {
		filepath = "common/resource/DinnerQR.jpg"
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

//recognizeFile Recognize QR image
func recognizeFile(path string) ([]*goqr.QRData, error) {
	log.Info(ctx, "recognize file: %v", path)
	imgData, err := ioutil.ReadFile(path)
	if err != nil {
		log.Error(ctx, "%v\n", err)
		return nil, err
	}
	img, _, err := image.Decode(bytes.NewReader(imgData))
	if err != nil {
		log.Error(ctx, "image.Decode error: %v", err)
		return nil, err
	}
	qrCodes, err := goqr.Recognize(img)
	if err != nil {
		log.Error(ctx, "Recognize failed: %v", err)
		return nil, err
	}
	return qrCodes, nil
}
