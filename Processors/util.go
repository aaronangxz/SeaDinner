package Processors

import "fmt"

func MakeToken(key string) string {
	return fmt.Sprint(TokenPrefix, key)
}

func MakeURL(opt int, id *int) string {
	switch opt {
	case URL_CURRENT:
		return fmt.Sprint(UrlPrefix, "/api/current")
	case URL_MENU:
		return fmt.Sprint(UrlPrefix, "/api/menu/", id)
	case URL_ORDER:
		return fmt.Sprint(UrlPrefix, "/api/order/", id)
	}
	return ""
}
