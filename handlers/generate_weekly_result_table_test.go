package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/common"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
	"github.com/aaronangxz/SeaDinner/test_helper"
	"google.golang.org/protobuf/proto"
	"testing"
	"time"
)

func TestGenerateWeeklyResultTable(t *testing.T) {
	ctx := context.TODO()
	m := test_helper.GetLiveMenuDetails()
	mC := MakeMenuCodeMap(ctx)
	r := []*sea_dinner.OrderRecord{
		{
			Id:        proto.Int64(1),
			UserId:    proto.Int64(12345),
			FoodId:    proto.String(fmt.Sprint(m[0].GetId())),
			OrderTime: proto.Int64(time.Now().Unix()),
			TimeTaken: proto.Int64(100),
			Status:    proto.Int64(int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)),
			ErrorMsg:  nil,
		},
		{
			Id:        proto.Int64(2),
			UserId:    proto.Int64(12345),
			FoodId:    proto.String(fmt.Sprint(m[1].GetId())),
			OrderTime: proto.Int64(time.Now().Unix() - 1*common.ONE_DAY),
			TimeTaken: proto.Int64(100),
			Status:    proto.Int64(int64(sea_dinner.OrderStatus_ORDER_STATUS_CANCEL)),
			ErrorMsg:  nil,
		},
		{
			Id:        proto.Int64(2),
			UserId:    proto.Int64(12345),
			FoodId:    proto.String(fmt.Sprint(m[2].GetId())),
			OrderTime: proto.Int64(time.Now().Unix() - 2*common.ONE_DAY),
			TimeTaken: proto.Int64(100),
			Status:    proto.Int64(int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL)),
			ErrorMsg:  nil,
		},
		{
			Id:        proto.Int64(3),
			UserId:    proto.Int64(12345),
			FoodId:    proto.String("696969"),
			OrderTime: proto.Int64(time.Now().Unix() - 3*common.ONE_DAY),
			TimeTaken: proto.Int64(100),
			Status:    proto.Int64(int64(sea_dinner.OrderStatus_ORDER_STATUS_OK)),
			ErrorMsg:  nil,
		},
	}
	start, end := processors.WeekStartEndDate(time.Now().Unix())
	header := fmt.Sprintf("Your orders from %v to %v\n", processors.ConvertTimeStampMonthDay(start), processors.ConvertTimeStampMonthDay(end))
	table := "<pre>\n    Day     Code  Status\n-------------------------\n"
	table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r[0].GetOrderTime()), mC[r[0].GetFoodId()], "🟢")
	table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r[1].GetOrderTime()), mC[r[1].GetFoodId()], "🟡")
	table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r[2].GetOrderTime()), mC[r[2].GetFoodId()], "🔴")
	table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r[3].GetOrderTime()), "??", "🟢")
	table += "</pre>"
	legend := "\n\n🟢 Success\n🟡 Cancelled\n🔴 Failed"
	expected := header + table + legend

	//Applicable for ListWeeklyResult
	//if !processors.IsWeekDay() {
	//	expected = "We are done for this week! Check again next week 😀"
	//}

	type args struct {
		ctx    context.Context
		record []*sea_dinner.OrderRecord
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			"HappyCase",
			args{ctx: ctx, record: r},
			expected,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GenerateWeeklyResultTable(tt.args.ctx, tt.args.record); got != tt.want {
				t.Errorf("GenerateWeeklyResultTable() = %v, want %v", got, tt.want)
			}
		})
	}
}
