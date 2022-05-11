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
	return fmt.Sprint(TokenPrefix, key)
}

func MakeURL(opt int, id *int) string {
	switch opt {
	case URL_CURRENT:
		return fmt.Sprint(UrlPrefix, "/api/current")
	case URL_MENU:
		if id == nil {
			return ""
		}
		return fmt.Sprint(UrlPrefix, "/api/menu/", *id)
	case URL_ORDER:
		if id == nil {
			return ""
		}
		return fmt.Sprint(UrlPrefix, "/api/order/", *id)
	}
	return ""
}

func OutputResults(resultMap map[int64]bool) {
	var (
		passed int
	)
	for _, m := range resultMap {
		if m {
			passed++
		}
	}

	fmt.Println("*************************")
	fmt.Println("Total Order: ", len(resultMap))
	fmt.Println("Total Success: ", passed)
	fmt.Println("Total Failures: ", len(resultMap)-passed)
	fmt.Println("*************************")
}
