package Processors

import (
	"testing"

	"github.com/aaronangxz/SeaDinner/Common"
)

func AdHocTestGetDayId(t *testing.T) {
	Common.Config.Adhoc = true
	Init()
	tests := []struct {
		name   string
		wantID int64
	}{
		{
			name:   "HappyCase",
			wantID: 3654,
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
