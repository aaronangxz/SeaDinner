package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/go-redis/redis"
	"gorm.io/gorm/clause"
	"os"
	"time"
)

func StoreFoodMappings(ctx context.Context) {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := App.StartTransaction("store_food_mappings")
	defer txn.End()

	year, week := ConvertTimeStampWeekOfYear(time.Now().Unix())
	cacheKey := fmt.Sprint(common.FOOD_MAPPING_PREFIX, year, ":", week)

	val, redisErr := CacheInstance().Get(cacheKey).Result()
	if redisErr != nil {
		if redisErr == redis.Nil {
			log.Warn(ctx, "StoreFoodMappings | Mapping not stored yet. | %v", cacheKey)
		} else {
			log.Error(ctx, "StoreFoodMappings | Error while reading from redis: %v", redisErr.Error())
			return
		}
	} else {
		if val == "1" {
			log.Warn(ctx, "StoreFoodMappings | Skipping, mapping has been stored | %v", cacheKey)
			return
		}
	}

	menu := GetMenu(ctx, key)
	if menu.GetFood() == nil {
		log.Warn(ctx, "StoreFoodMappings | No record to store.")
		return
	}

	mappings := ConvertFoodToFoodMappingByYearAndWeek(ctx, menu.GetFood())

	//Assuming there will be no duplicated food_id
	if err := DbInstance().Clauses(clause.OnConflict{DoNothing: true}).Table(common.DB_FOOD_MAPPING_TAB).Create(&mappings).Error; err != nil {
		log.Error(ctx, fmt.Sprintf("StoreFoodMappings | Failed to insert records | %v", err.Error()))
		return
	}

	if err := CacheInstance().Set(cacheKey, "1", 0).Err(); err != nil {
		log.Error(ctx, "StoreFoodMappings | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "StoreFoodMappings | Successful | Written %v to redis", cacheKey)
	}
	log.Info(ctx, "StoreFoodMappings | Successfully stored food mappings.")
}
