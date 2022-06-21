package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/log"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"os"
	"testing"
	"time"
)

func TestSendPotentialUsers(t *testing.T) {
	os.Setenv("TEST_DEPLOY", "TRUE")
	test_helper.InitTest()
	ctx := context.TODO()
	u := test_helper.RandomInt(99999999)
	u1 := test_helper.RandomInt(99999999)
	toWrite := fmt.Sprint(u, ":", time.Now().Unix()-common.ONE_MONTH-common.ONE_DAY)
	toWrite1 := fmt.Sprint(u1, ":", time.Now().Unix()-common.ONE_DAY)

	if err := processors.CacheInstance().SAdd(common.POTENTIAL_USER_SET, toWrite).Err(); err != nil {
		log.Error(ctx, "TestSendPotentialUsers | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "TestSendPotentialUsers | Successful | Written %v to potential_user set", toWrite)
	}

	if err := processors.CacheInstance().SAdd(common.POTENTIAL_USER_SET, toWrite1).Err(); err != nil {
		log.Error(ctx, "TestSendPotentialUsers | Error while writing to redis: %v", err.Error())
	} else {
		log.Info(ctx, "TestSendPotentialUsers | Successful | Written %v to potential_user set", toWrite1)
	}

	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
	}{
		{
			"HappyCaseAndNotWithinTimeRange",
			args{ctx},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			SendPotentialUsers(tt.args.ctx)
		})
	}

	if err := processors.CacheInstance().Del(common.POTENTIAL_USER_SET).Err(); err != nil {
		log.Error(ctx, "TestSendPotentialUsers | Error while erasing from redis: %v", err.Error())
	} else {
		log.Info(ctx, "TestSendPotentialUsers | Successful | Deleted potential_user set")
	}
	os.Setenv("TEST_DEPLOY", "FALSE")
}
