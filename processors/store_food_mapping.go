package processors

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"gorm.io/gorm/clause"
	"os"
)

func StoreFoodMappings(ctx context.Context) {
	var (
		key = os.Getenv("TOKEN")
	)
	txn := App.StartTransaction("store_food_mappings")
	defer txn.End()

	menu := GetMenu(ctx, key)
	if menu.GetFood() == nil {
		log.Warn(ctx, "StoreFoodMappings | No record to store.")
		return
	}

	mappings := ConvertFoodToFoodMapping(ctx, menu.GetFood())

	//Assuming there will be no duplicated food_id
	if err := DB.Clauses(clause.OnConflict{DoNothing: true}).Table(common.DB_FOOD_MAPPING_TAB).Create(&mappings).Error; err != nil {
		log.Error(ctx, fmt.Sprintf("StoreFoodMappings | Failed to insert records | %v", err.Error()))
		return
	}
	log.Info(ctx, fmt.Sprintf("StoreFoodMappings | Successfully stored food mappings | size: %v", len(mappings)))
}
