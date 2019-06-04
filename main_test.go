package main

import "testing"

func Test_checkCode(t *testing.T) {
	type args struct {
		expected Response
		actual   Response
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"empty responses", args{Response{}, Response{}}, true},
		{"same 200 codes", args{Response{Code: 200}, Response{Code: 200}}, true},
		{"same 0 codes", args{Response{Code: 0}, Response{Code: 0}}, true},
		{"different codes 0 and 400", args{Response{Code: 0}, Response{Code: 400}}, false},
		{"different codes 200 and 201", args{Response{Code: 200}, Response{Code: 201}}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := checkCode(tt.args.expected, tt.args.actual); got != tt.want {
				t.Errorf("checkCode() = %v, want %v", got, tt.want)
			}
		})
	}
}
