package Processors

import (
	"fmt"
	"log"
)

func MakeToken(key string) string {
	if key == "" {
		log.Println("Key is invalid:", key)
		return ""
	}
	return fmt.Sprint(Config.Prefix.TokenPrefix, key)
}

func MakeURL(opt int, id *int) string {
	prefix := Config.Prefix.UrlPrefix
	switch opt {
	case URL_CURRENT:
		return fmt.Sprint(prefix, "/api/current")
	case URL_MENU:
		if id == nil {
			return ""
		}
		return fmt.Sprint(prefix, "/api/menu/", *id)
	case URL_ORDER:
		if id == nil {
			return ""
		}
		return fmt.Sprint(prefix, "/api/order/", *id)
	}
	return ""
}

func OutputResults(resultMap map[int64]int) {
	var (
		passed int
	)
	for _, m := range resultMap {
		if m == ORDER_STATUS_OK {
			passed++
		}
	}

	fmt.Println("*************************")
	fmt.Println("Total Order: ", len(resultMap))
	fmt.Println("Total Success: ", passed)
	fmt.Println("Total Failures: ", len(resultMap)-passed)
	fmt.Println("*************************")
}
