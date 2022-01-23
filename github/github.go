package github

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type Github struct {
	TOKEN string
}

func (g *Github) SendReq(method string, url string, body map[string]interface{}) (map[string]interface{}, error) {
	payloadBuf := new(bytes.Buffer)
	err := json.NewEncoder(payloadBuf).Encode(body)

	if err != nil {
		return nil, err
	}

	req, _ := http.NewRequest(method, url, payloadBuf)

	req.Header.Add("Authorization", "token "+g.TOKEN)
	req.Header.Add("Content-Type", "application/json")

	return perform(req)
}

func perform(req *http.Request) (map[string]interface{}, error) {
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	if res.StatusCode > 400 {
		buf := new(bytes.Buffer)
		buf.ReadFrom(res.Body)
		return nil, fmt.Errorf("API error: %s", buf.String())
	}

	var r map[string]interface{}
	if err := json.NewDecoder(res.Body).Decode(&r); err != nil {
		return nil, err
	}

	return r, err
}
