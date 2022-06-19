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

	m := MakeMenuCodeMap(ctx)

	status := map[int64]string{
		int64(sea_dinner.OrderStatus_ORDER_STATUS_OK):     "🟢",
		int64(sea_dinner.OrderStatus_ORDER_STATUS_FAIL):   "🔴",
		int64(sea_dinner.OrderStatus_ORDER_STATUS_CANCEL): "🟡"}

	header := fmt.Sprintf("Your orders from %v to %v\n", processors.ConvertTimeStampMonthDay(start), processors.ConvertTimeStampMonthDay(end))

	table := "<pre>\n    Day     Code  Status\n-------------------------\n"
	for _, r := range record {
		//In the event when menu was changed during the week, and we have no info of that particular food code
		var code string
		if _, ok := m[r.GetFoodId()]; !ok {
			code = "??"
		} else {
			code = m[r.GetFoodId()]
		}
		table += fmt.Sprintf(" %v   %v     %v\n", processors.ConvertTimeStampDayOfWeek(r.GetOrderTime()), code, status[r.GetStatus()])
	}
	table += "</pre>"
	legend := "\n\n🟢 Success\n🟡 Cancelled\n🔴 Failed"
	return fmt.Sprint(header, table, legend)
}
