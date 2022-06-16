package handlers

import (
	"github.com/aaronangxz/SeaDinner/log"
	"testing"
)

func TestMain(m *testing.M) {
	log.InitializeLogger()
	m.Run()
}
