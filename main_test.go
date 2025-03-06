package main

import (
	"fmt"
	"net/http"
	"testing"
)

// Test to ensure http response from NWS API is 200

func mainTesting() int {

	url := fmt.Sprintf("https://api.weather.gov")

	resp, _ := http.Get(url)

	return resp.StatusCode
}

func TestMain(t *testing.T) {
	expected := mainTesting()
	want := 200

	if mainTesting() != want {
		t.Errorf("got %q want %q // CHECK THAT PARAMETERS BEING PASSED RESULT IN A 200 RESPONSE CODE FROM NWS DB", expected, want)
	}
}
