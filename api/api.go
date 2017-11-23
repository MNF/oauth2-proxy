package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/bitly/go-simplejson"
)

func Request(req *http.Request) (*simplejson.Json, error) {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error Connecting Request %s %s", req.Method, req.URL)
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("Request Status %d", resp.StatusCode)
	}
	data, err := simplejson.NewJson(body)
	if err != nil {
		return nil, fmt.Errorf("Request Status %d", resp.StatusCode)
	}
	return data, nil
}

func RequestJson(req *http.Request, v interface{}) error {
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Printf("Error Connecting RequestJson %s %s %s", req.Method, req.URL, err)
		return err
	}
	body, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return err
	}
	if resp.StatusCode != 200 {
		return fmt.Errorf("got %d %s", resp.StatusCode, req.URL)
	}
	return json.Unmarshal(body, v)
}

func RequestUnparsedResponse(url string, header http.Header) (resp *http.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header = header

	return http.DefaultClient.Do(req)
}
