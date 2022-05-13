package TestHelper

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/aaronangxz/SeaDinner/Processors"
)

func RandomString(length int) string {
	rand.Seed(time.Now().UnixNano())
	b := make([]byte, length)
	rand.Read(b)
	return fmt.Sprintf("%x", b)[:length]
}

func RandomInt(max int) int64 {
	rand.Seed(time.Now().UnixNano())
	return int64(rand.Intn(max))
}

func InitTest() {
	Processors.ConnectTestDataBase()
}
