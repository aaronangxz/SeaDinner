package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"testing"
)

func AdHocTestGetDayID(t *testing.T) {
	common.Config.Adhoc = true
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
			if gotID := GetDayID(context.TODO()); gotID != tt.wantID {
				t.Errorf("GetDayId() = %v, want %v", gotID, tt.wantID)
			}
		})
	}
}
