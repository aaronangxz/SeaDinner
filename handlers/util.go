package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/processors"
	"os"
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
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetName()
	}
	// Store -1 hash to menuMap
	menuMap["-1"] = "*NOTHING*" // to be renamed
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
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetCode()
	}
	menuMap["RAND"] = "Random"
	return menuMap
}
