package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

var (
	gatherIP       = "https://tenant.domain.net/"
	gatherUsername = "username"
	gatherPassword = "password"
	tenantID       = "tenantID"
)

func loginGather() (*http.Cookie, error) {
	jsonData := `username=` + gatherUsername + `&password=` + gatherPassword
	request, _ := http.NewRequest("POST", gatherIP+"api/v1/auth/login", bytes.NewBufferString(jsonData))
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		return nil, err
	}

	if response.Status == "200 OK" {
		fmt.Println("Connected to Gather")
		return response.Cookies()[0], nil
	}

	return nil, errors.New("Did not connect to Gather. Status: " + response.Status)
}

func getPolicies(cookie *http.Cookie) {
	request, _ := http.NewRequest("GET", gatherIP+"api/v2/policies/alerting?policyType=security", bytes.NewBufferString(""))
	request.Header.Set("X-Forwarded-Tenant-Id", tenantID)
	request.AddCookie(cookie)

	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	var jsonObject interface{}

	body, _ := ioutil.ReadAll(response.Body)
	json.Unmarshal(body, &jsonObject)
	data, _ := json.MarshalIndent(jsonObject, "", "    ")
	fmt.Println(string(data))
}

func main() {
	cookie, err := loginGather()

	if err != nil {
		fmt.Println(err.Error())
	} else {
		getPolicies(cookie)
	}
}
