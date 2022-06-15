package processors

import (
	"context"
	"github.com/aaronangxz/SeaDinner/common"
	"os"
	"reflect"
	"testing"
)

func AdhocTestGetMenuUsingCache(t *testing.T) {
	LoadEnv()
	common.LoadConfig()
	ConnectTestRedis()
	ConnectTestMySQL()
	key := os.Getenv("TOKEN")
	type args struct {
		key string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "HappyCase",
			args: args{key: key},
			want: 12,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetMenuUsingCache(context.TODO(), tt.args.key); !reflect.DeepEqual(len(got.GetFood()), tt.want) {
				t.Errorf("GetMenu() = %v, want %v", len(got.GetFood()), tt.want)
			}
		})
	}
}
