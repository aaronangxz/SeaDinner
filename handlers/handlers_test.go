package handlers

import (
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"testing"
)

func TestMain(m *testing.M) {
	log.InitializeLogger()
	test_helper.InitTest()
	m.Run()
}
