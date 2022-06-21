package handlers

import (
	"context"
	"fmt"
	"github.com/aaronangxz/SeaDinner/processors"
	"github.com/aaronangxz/SeaDinner/sea_dinner.pb"
)

//GenerateResultTable Outputs pre-formatted order status.
func GenerateResultTable(ctx context.Context, record []*sea_dinner.OrderRecord, start int64, end int64) string {
	txn := processors.App.StartTransaction("generate_weekly_result_table")
	defer txn.End()

	m := MakeFoodMapping(ctx)
	status := map[int64]string{
		int64(sea_dinner.OrderStatus_ORDER_STATUS_OK):     "🟢",
		int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL):   "🔴",
		int64(sea_dinner.OrderStatus_ORDER_STATUS_CANCEL): "🟡"}

	header := fmt.Sprintf("Your orders from %v to %v\n", processors.ConvertTimeStampMonthDay(start), processors.ConvertTimeStampMonthDay(end))

	table := "<pre>\n    Day     Code  Status\n-------------------------\n"
	for _, r := range record {
		year, week := processors.ConvertTimeStampWeekOfYear(r.GetOrderTime())
		code, ok := m[year][week][r.GetFoodId()]
		if !ok {
			code = "??"
		}
		table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r.GetOrderTime()), code, status[r.GetStatus()])
	}
	table += "</pre>"
	legend := "\n\n🟢 Successful\n🟡 Cancelled\n🔴 Failed\n ?? Dish removed"
	return fmt.Sprint(header, table, legend)
}
