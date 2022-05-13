package Processors

import "testing"

func Test_IsWeekDay(t *testing.T) {
	tests := []struct {
		name string
		want bool
	}{
		{
			name: "HappyCase_Weekday",
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsWeekDay(); got != tt.want {
				t.Errorf("isWeekDay() = %v, want %v", got, tt.want)
			}
		})
	}
}
