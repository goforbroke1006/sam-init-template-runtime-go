package main

import (
	"testing"
)

func TestHandler(t *testing.T) {
	//t.Run("Unable to get IP", func(t *testing.T) {
	//	DefaultHTTPGetAddress = "http://127.0.0.1:12345"
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err == nil {
	//		t.Fatal("Error failed to trigger with an invalid request")
	//	}
	//})
	//
	//t.Run("Successful Request", func(t *testing.T) {
	//	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	//		w.WriteHeader(200)
	//		fmt.Fprintf(w, "127.0.0.1")
	//	}))
	//	defer ts.Close()
	//
	//	DefaultHTTPGetAddress = ts.URL
	//
	//	_, err := handler(events.APIGatewayProxyRequest{})
	//	if err != nil {
	//		t.Fatal("Everything should be ok")
	//	}
	//})
}
