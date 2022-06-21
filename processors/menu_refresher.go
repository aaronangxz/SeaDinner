package processors

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"os"
	"time"
)

//MenuRefresher Periodically refreshes cached menu with the updated live menu
func MenuRefresher(ctx context.Context) {
	ticker := time.NewTicker(time.Duration(common.Config.Runtime.MenuRefreshIntervalSeconds) * time.Second)

	for range ticker.C {
		func() {
			if !IsActiveDay() {
				if IsWeekDay() {
					log.Warn(ctx, "MenuRefresher | Weekday but not active | Resumes check later.")
					time.Sleep(3600 * time.Second)
				} else {
					log.Warn(ctx, "MenuRefresher | Inactive day | Resumes check tomorrow.")
					time.Sleep(time.Duration(GetEOD().Unix()-time.Now().Unix()) * time.Second)
				}
				return
			}
			key := os.Getenv("TOKEN")
			log.Info(ctx, "MenuRefresher | Comparing Live and Cached menu.")

			liveMenu := GetMenu(ctx, key)
			cacheMenu := GetMenuUsingCache(ctx, key)

			if !CompareSliceStruct(ctx, liveMenu.GetFood(), cacheMenu.GetFood()) {
				log.Warn(ctx, "MenuRefresher | Live and Cached menu are inconsistent.")
				cacheKey := fmt.Sprint(common.MENU_CACHE_KEY_PREFIX, ConvertTimeStamp(time.Now().Unix()))

				data, err := json.Marshal(liveMenu)
				if err != nil {
					log.Error(ctx, "MenuRefresher | Failed to marshal JSON results: %v\n", err.Error())
				}

				//Use live menu as the source of truth
				if err := RedisClient.Set(cacheKey, data, 0).Err(); err != nil {
					log.Error(ctx, "MenuRefresher | Error while writing to redis: %v", err.Error())
				} else {
					log.Info(ctx, "MenuRefresher | Successful | Written %v to redis", cacheKey)
				}
			}
			log.Info(ctx, "MenuRefresher | Live and Cached menu are consistent.")
		}()
	}
}
