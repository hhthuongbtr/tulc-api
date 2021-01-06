package utils

import (
	"io/ioutil"
	"net/http"
)

func HttpGet(url string) (responseBody string, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	//We Read the response body on the line below.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	//Convert the body to type string
	sb := string(body)
	return sb, nil
}
