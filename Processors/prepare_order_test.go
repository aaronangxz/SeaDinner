package Processors

import (
	"context"
	"testing"
)

func TestGetOrderByUserId(t *testing.T) {
	LoadEnv()
	ConnectTestMySQL()
	type args struct {
		user_id int64
	}
	tests := []struct {
		name  string
		args  args
		want  string
		want1 bool
	}{
		{
			name:  "HappyCase",
			args:  args{12345},
			want:  "I can't find your order 😥 Try to cancel from SeaTalk instead!",
			want1: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := GetOrderByUserId(context.TODO(), tt.args.user_id)
			if got != tt.want {
				t.Errorf("GetOrderByUserId() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("GetOrderByUserId() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
