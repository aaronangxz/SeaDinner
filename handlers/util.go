package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"google.golang.org/protobuf/proto"
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
		menu = processors.GetMenu(ctx, key)
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
		menu = processors.GetMenu(ctx, key)
	}
	for _, m := range menu.GetFood() {
		menuMap[fmt.Sprint(m.GetId())] = m.GetCode()
	}
	menuMap["RAND"] = "Random"
	return menuMap
}

//MakeFoodMapping Returns the archived food_id:food_code mapping of current menu
//Food Code may be inconsistent
func MakeFoodMapping(ctx context.Context) map[int64]map[int64]map[string]string {
	var (
		mappings []*sea_dinner.FoodMappingByYearAndWeek
	)
	txn := processors.App.StartTransaction("make_menu_code_map")
	defer txn.End()
	mapped := make(map[int64]map[int64]map[string]string)
	if err := processors.DbInstance().Raw("SELECT * FROM food_mapping_tab").Scan(&mappings).Error; err != nil {
		log.Error(ctx, err.Error())
		return nil
	}

	for _, m := range mappings {
		if mapped[m.GetYear()] == nil {
			mapped[m.GetYear()] = make(map[int64]map[string]string)
		}

		if mapped[m.GetYear()][m.GetWeek()] == nil {
			mapped[m.GetYear()][m.GetWeek()] = make(map[string]string)
		}
		food := sea_dinner.FoodMappings{}
		err := proto.Unmarshal(m.GetFoodMapping(), &food)
		if err != nil {
			log.Error(ctx, "MakeFoodMapping | Failed | %v", err.Error())
			return nil
		}
		for _, f := range food.GetFoodMapping() {
			mapped[m.GetYear()][m.GetWeek()][fmt.Sprint(f.GetFoodId())] = f.GetFoodCode()
		}
	}
	log.Info(ctx, "MakeFoodMapping | Success")
	return mapped
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
