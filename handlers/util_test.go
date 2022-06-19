package handlers

import "testing"

func TestIsContainsSpace(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"NoSpace",
			args{"IHaveNoSpace"},
			false,
		},
		{
			"HasSpace",
			args{"I Have Space"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContainsSpace(tt.args.a); got != tt.want {
				t.Errorf("IsContainsSpace() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIsContainsSpecialChar(t *testing.T) {
	type args struct {
		a string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{
			"NoSpecialChar",
			args{"abcdef12345"},
			false,
		},
		{
			"HasSpecialChar",
			args{"abcdef12345!!"},
			true,
		},
		{
			"Emoji",
			args{"😍"},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsContainsSpecialChar(tt.args.a); got != tt.want {
				t.Errorf("IsContainsSpecialChar() = %v, want %v", got, tt.want)
			}
		})
	}
}
