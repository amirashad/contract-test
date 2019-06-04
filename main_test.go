package main

import (
	"net/http"
	"testing"
)

func Test_checkCode(t *testing.T) {
	type args struct {
		expected Response
		actual   Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty responses", args{Response{}, Response{}}, false},
		{"same 200 codes", args{Response{Code: 200}, Response{Code: 200}}, false},
		{"same 0 codes", args{Response{Code: 0}, Response{Code: 0}}, false},
		{"different codes 0 and 400", args{Response{Code: 0}, Response{Code: 400}}, true},
		{"different codes 200 and 201", args{Response{Code: 200}, Response{Code: 201}}, true}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkCode(tt.args.expected, tt.args.actual); (err != nil) != tt.wantErr {
				t.Errorf("checkCode() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkHeaders(t *testing.T) {
	type args struct {
		expected Response
		actual   Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty responses", args{Response{}, Response{}}, false},
		{"empty expected and with some actual responses", args{
			expected: Response{},
			actual:   Response{Headers: map[string][]string{"Some-Header": []string{"1", "2"}}},
		}, false},
		{"same headers", args{
			expected: Response{Headers: map[string]interface{}{"Some-Header": "Some-Value"}},
			actual:   Response{Headers: http.Header{"Some-Header": []string{"Some-Value"}}},
		}, false},
		{"same more headers", args{
			expected: Response{Headers: map[string]interface{}{"Some-Header": "Some-Value", "Other-Header1": "other-value1", "Other-Header2": ""}},
			actual:   Response{Headers: http.Header{"Some-Header": []string{"Some-Value"}, "Other-Header1": []string{"other-value1"}, "Other-Header2": []string{""}}},
		}, false},
		{"some headers with different order", args{
			expected: Response{Headers: map[string]interface{}{"Some-Header": "Some-Value", "Other-Header1": "other-value1", "Other-Header2": ""}},
			actual:   Response{Headers: http.Header{"Other-Header2": []string{""}, "Some-Header": []string{"Some-Value"}, "Other-Header1": []string{"other-value1"}}},
		}, false},
		{"different headers", args{
			expected: Response{Headers: map[string]interface{}{"Header1": "Some-Value", "Header2": ""}},
			actual:   Response{Headers: http.Header{"Some-Header": []string{"Some-Value"}, "Other-Header1": []string{"other-value1"}, "Other-Header2": []string{""}}},
		}, true},
		{"expected headers with empty actual", args{
			expected: Response{Headers: map[string]interface{}{"Some-Header": "Some-Value"}},
			actual:   Response{},
		}, true}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkHeaders(tt.args.expected, tt.args.actual); (err != nil) != tt.wantErr {
				t.Errorf("checkHeaders() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_checkBody(t *testing.T) {
	type args struct {
		expected Response
		actual   Response
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{"empty responses", args{Response{}, Response{}}, false},
		{"no any expected body", args{
			expected: Response{Headers: map[string]interface{}{"Content-Type": "text/plain"}},
			actual:   Response{Headers: http.Header{"Content-Type": []string{"text/plain"}}},
		}, false},
		{"same text bodies", args{
			expected: Response{
				Headers: map[string]interface{}{"Content-Type": "text/plain"},
				Body:    string("Lorem ipsum"),
			},
			actual: Response{
				Headers: http.Header{"Content-Type": []string{"text/plain"}},
				Body:    []byte("Lorem ipsum"),
			},
		}, false},
		{"same json bodies", args{
			expected: Response{
				Headers: map[string]interface{}{"Content-Type": "application/json"},
				Body:    map[string]interface{}{"id": 12, "height": 12.2, "exists": true, "name": "some name", "surname": nil},
			},
			actual: Response{
				Headers: http.Header{"Content-Type": []string{"application/json"}},
				Body:    []byte("{ \"name\": \"some name\", \"id\":12,\"height\":12.2, \"exists\":true, \"surname\": null }"),
			},
		}, false},
		{"expected body exists but actual body is nil", args{
			expected: Response{
				Headers: map[string]interface{}{"Content-Type": "application/json"},
				Body:    map[string]interface{}{"id": 12, "height": 12.2, "exists": true, "name": "some name", "surname": nil},
			},
			actual: Response{
				Headers: http.Header{"Content-Type": []string{"application/json"}},
				Body:    nil,
			},
		}, true},
		{"expected body exists but actual body is empty", args{
			expected: Response{
				Headers: map[string]interface{}{"Content-Type": "application/json"},
				Body:    map[string]interface{}{"id": 12, "height": 12.2, "exists": true, "name": "some name", "surname": nil},
			},
			actual: Response{
				Headers: http.Header{"Content-Type": []string{"application/json"}},
				Body:    []byte("{  }"),
			},
		}, true},
		{"different bodies", args{
			expected: Response{
				Headers: map[string]interface{}{"Content-Type": "application/json"},
				Body:    map[string]interface{}{"id": 1, "height": 12.3, "exists": false, "name": "some other name", "surname": "other"},
			},
			actual: Response{
				Headers: http.Header{"Content-Type": []string{"application/json"}},
				Body:    []byte("{ \"name\": \"some name\", \"id\":12,\"height\":12.2, \"exists\":true, \"surname\": null }"),
			},
		}, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := checkBody(tt.args.expected, tt.args.actual); (err != nil) != tt.wantErr {
				t.Errorf("checkBody() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_isNumber(t *testing.T) {
	type args struct {
		a interface{}
	}
	tests := []struct {
		name  string
		args  args
		want  bool
		want1 float64
	}{
		{"empty responses", args{1.1}, true, 1.1},
		{"empty responses", args{1}, true, 1},
		{"empty responses", args{1}, true, 1.0},
		{"empty responses", args{int8(1)}, true, 1.0},
		{"empty responses", args{int8(1)}, true, 1.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1 := isNumber(tt.args.a)
			if got != tt.want {
				t.Errorf("isNumber() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("isNumber() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
