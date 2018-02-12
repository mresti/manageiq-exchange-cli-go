package api

import (
	"reflect"
	"testing"
)

func TestApi_Init(t *testing.T) {
	var server Api
	inputServer := "localhost"
	inputPort := 3000
	wantServer := "localhost"
	wantPort := 3000
	server.Init(inputServer, inputPort)
	if !reflect.DeepEqual(server.Server, wantServer) || !reflect.DeepEqual(server.Port, wantPort) {
		t.Errorf("Api.Init(%v, %v) returned server=%v and port=%v, want server=%v and port=%v", inputServer, inputPort, server.Server, server.Port, wantServer, wantPort)
	}
}

func TestApi_CheckConnectionServer(t *testing.T) {
	t.Skip("TODO")
}

func TestApi_URL(t *testing.T) {
	var tests = []struct {
		inputServer string
		inputPort   int
		wantURL     string
	}{
		{"localhost", 0, "localhost"},
		{"localhost", 3000, "localhost:3000"},
	}
	for _, tt := range tests {
		t.Run(tt.inputServer, func(t *testing.T) {
			var server Api
			server.Init(tt.inputServer, tt.inputPort)
			gotURL := server.URL()
			if !reflect.DeepEqual(server.Server, tt.wantURL) {
				t.Fatalf("Api.URL() returned %v, want %v", gotURL, tt.wantURL)
			}
		})
	}
}
