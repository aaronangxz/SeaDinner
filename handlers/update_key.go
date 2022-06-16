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
	"time"
)

//UpdateKey Creates record to store user's key if not exists, or update the existing record.
//With basic parameter verifications
func UpdateKey(ctx context.Context, id int64, s string) (string, bool) {
	hashedKey := processors.EncryptKey(s, os.Getenv("AES_KEY"))

	var (
		cacheKey       = fmt.Sprint(common.USER_KEY_PREFIX, id)
		existingRecord sea_dinner.UserKey
		r              = &sea_dinner.UserKey{
			UserId:  proto.Int64(id),
			UserKey: proto.String(hashedKey),
			Ctime:   proto.Int64(time.Now().Unix()),
			Mtime:   proto.Int64(time.Now().Unix()),
			IsMute:  proto.Int64(int64(sea_dinner.MuteStatus_MUTE_STATUS_NO)),
		}
	)
	txn := processors.App.StartTransaction("update_key")
	defer txn.End()

	if id <= 0 {
		log.Error(ctx, "UpdateKey | Id must be > 1.")
		return "", false
	}

	if s == "" {
		log.Error(ctx, "UpdateKey | Key cannot be empty.")
		return "Key cannot be empty 😟", false
	}

	if len(s) != 40 {
		log.Error(ctx, "UpdateKey | Key length invalid | length: %v", len(s))
		return "Are you sure this is a valid key? 😟", false
	}

	if err := processors.DB.Raw("SELECT * FROM user_key_tab WHERE user_id = ?", id).Scan(&existingRecord).Error; err != nil {
		log.Error(ctx, "UpdateKey | %v", err.Error())
		return err.Error(), false
	}
	if existingRecord.UserId == nil {
		if err := processors.DB.Table(common.DB_USER_KEY_TAB).Create(&r).Error; err != nil {
			log.Error(ctx, "UpdateKey | Failed to insert DB | %v", err.Error())
			return err.Error(), false
		}
		return "Okay got it. I remember your key now! 😙\n Disclaimer: I will never disclose your key. Your key is safely encrypted.", true
	}
	//Update key if user_id exists
	if err := processors.DB.Exec("UPDATE user_key_tab SET user_key = ?, mtime = ? WHERE user_id = ?", hashedKey, time.Now().Unix(), id).Error; err != nil {
		log.Error(ctx, "UpdateKey | Failed to insert DB | %v", err.Error())
		return err.Error(), false
	}

	//Invalidate cache after successful update
	if _, err := processors.RedisClient.Del(cacheKey).Result(); err != nil {
		log.Error(ctx, "UpdateKey | Failed to invalidate cache: %v. %v", cacheKey, err)
	}
	log.Info(ctx, "UpdateKey | Successfully invalidated cache: %v", cacheKey)

	return "Okay got it. I will take note of your new key 😙", true
}
