package handlers

import (
	"context"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper/user_log"
	"testing"
)

func TestBatchGetLatestResult(t *testing.T) {
	ctx := context.TODO()
	lunch := processors.GetLunchTime().Unix()
	uOk := user_log.New().SetStatus(int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)).SetOrderTime(lunch).Build()
	uTooOld := user_log.New().SetStatus(int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)).SetOrderTime(lunch - 600).Build()
	uTooNew := user_log.New().SetStatus(int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)).SetOrderTime(lunch + 600).Build()
	defer func() {
		uOk.TearDown()
		uTooOld.TearDown()
		uTooNew.TearDown()
	}()
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want *sea_dinner.OrderRecord
	}{
		{
			"HappyCase",
			args{ctx: ctx},
			user_log.ConvertUserLogToOrderRecord(uOk),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := BatchGetLatestResult(tt.args.ctx)
			has := false
			for _, e := range got {
				if e.GetUserId() == tt.want.GetUserId() {
					has = true
				}
			}
			if !has {
				t.Errorf("BatchGetUsersChoiceWithKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
