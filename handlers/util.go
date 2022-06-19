package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/processors"
	"os"
	"unicode"
)

//MakeMenuNameMap Returns food_id:food_name mapping of current menu
func MakeMenuNameMap(ctx context.Context) map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := processors.App.StartTransaction("make_menu_name_map")
	defer txn.End()

	menuMap := make(map[string]string)
	menu := processors.GetMenuUsingCache(ctx, key)

	if common.Config.UnitTest {
		menu = processors.GetMenu(ctx, processors.Client, key)
	}
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "<b>NOTHING</b>" // to be renamed
	menuMap["RAND"] = "Random"
	return menuMap
}

//MakeMenuCodeMap Returns food_id:food_code mapping of current menu
func MakeMenuCodeMap(ctx context.Context) map[string]string {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := processors.App.StartTransaction("make_menu_code_map")
	defer txn.End()

	menuMap := make(map[string]string)
	menu := processors.GetMenuUsingCache(ctx, key)

	if common.Config.UnitTest {
		menu = processors.GetMenu(ctx, processors.Client, key)
	}
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetCode()
	}
	menuMap["RAND"] = "Random"
	return menuMap
}

func IsContainsSpecialChar(a string) bool {
	for _, char := range a {
		if unicode.IsSymbol(char) {
			return true
		}
	}
	for _, char := range a {
		if !unicode.IsLetter(char) && !unicode.IsNumber(char) {
			return true
		}
	}
	return false
}

func IsContainsSpace(a string) bool {
	for _, char := range a {
		if unicode.IsSpace(char) {
			return true
		}
	}
	return false
}
