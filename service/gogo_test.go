package service

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/unrolled/render"
)

var (
	formatter = render.New(render.Options{
		IndentJSON: true,
	})
)

const (
	fakeMatchLocationResult = "/matches/5a003b78-409e-4452-b456-a6f0dcee05bd"
)

func TestCreateMatch(t *testing.T) {
	client := &http.Client{}
	server := httptest.NewServer(http.HandlerFunc(createMatchHandler(formatter)))
	defer server.Close()
	body := []byte("{\n \"gridsize\": 19,\n \"players\":[\n {\n \"color\": \"white\",\n \"name\":\"bob\"\n },\n {\n \"color\":\"black\",\n \"name\": \"alfred\"\n }\n ]\n}")
	req, err := http.NewRequest("POST", server.URL, bytes.NewBuffer(body))
	if err != nil {
		t.Errorf("Error in creating post request for createMatchHandler: %v", err)
	}
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		t.Errorf("Error in POST to createMatchHandler: %v", err)
	}
	defer res.Body.Close()
	payload, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Errorf("Error in parsing response body: %v", err)
	}
	if res.StatusCode != http.StatusCreated {
		t.Errorf("Expected response status 201, but %s", res.Status)
	}
	fmt.Printf("Payload: %s", string(payload))
	loc, ok := res.Header["Location"]
	if !ok {
		t.Error("Location header is not set")
	} else if !strings.Contains(loc[0], "/matches/") {
		t.Error("Location not contain matchers")
	} else if len(loc[0]) != len(fakeMatchLocationResult) {
		t.Error("Location value does not contain guid")
	}
}
