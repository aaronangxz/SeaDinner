package Processors

import "testing"

func AdhocTestGetDayId(t *testing.T) {
	Config.Adhoc = true
	Init()
	tests := []struct {
		name   string
		wantID int
	}{
		{
			name:   "HappyCase",
			wantID: 3556,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotID := GetDayId(); gotID != tt.wantID {
				t.Errorf("GetDayId() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
