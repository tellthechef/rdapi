package main

import (
	rd "../"
	"fmt"
	"io/ioutil"
	"net/http"
)

func main() {
	// - redacted api setup
	// if err := api.Authenticate(); err != nil {
	// 	return
	// }

	req, _ := api.NewRequest("GET", "/Restaurant/4075", nil)
	req, _ := api.NewRequest("GET", "/Restaurant/4075/DiaryData?date=2016-01-19", nil)
	client := &http.Client{}

	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)

	fmt.Println(string(body))
}
