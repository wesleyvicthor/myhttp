package myhttp

import (
	"reflect"
	"runtime"
	"testing"
)

func TestProcess(t *testing.T) {
	type urls []string
	tests := []struct {
		name    string
		args    urls
		want    *Response
		wantErr bool
	}{
		{"bad url format", urls{"adjust."}, nil, true},
		{"md5 response body", urls{"example.com"}, &Response{"http://example.com", [16]byte{132,35,141,252,128,146,229,217,192,218,200,239,147,55,26,7}}, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := New(tt.args)
			got, err := mh.Process(mh.Urls()[0])
			if (err != nil) != tt.wantErr {
				t.Errorf("Process() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Process() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type urls []string
	tests := []struct {
		name string
		args []string
		want *MyHTTP
	}{
		{"no schema urls", []string{"adjust.com", "twitter.com"}, &MyHTTP{[]string{"http://adjust.com", "http://twitter.com"}}},
		{"mixed with and no schema", []string{"http://adjust.com", "twitter.com"}, &MyHTTP{[]string{"http://adjust.com", "http://twitter.com"}}},
		{"with https schema", []string{"https://adjust.com", "twitter.com"}, &MyHTTP{[]string{"https://adjust.com", "http://twitter.com"}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestProcessParallel(t *testing.T) {
	done := make(chan struct{})

	tests := []struct {
		name   string
		urls []string
		want   *Response
	}{
		{"running processes", []string{"example.com"}, &Response{"http://example.com", [16]byte{132,35,141,252,128,146,229,217,192,218,200,239,147,55,26,7}}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mh := New(tt.urls)
			result, _ := mh.ProcessParallel(10, done)
			defer close(result)
			if runtime.NumGoroutine() < 10 {
				t.Errorf("failed expected amount of parallel processes = %v, want %v", runtime.NumGoroutine(), 10)
			}
			if got := <-result; !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}