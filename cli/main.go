package main

import (
	rd "../"
	"fmt"
	"io/ioutil"
)

func main() {
	api := rd.New("consumerKey", "consumerSecret", "secondSecret")
	if err := api.Authenticate(); err != nil {
		return
	}

	fmt.Println("\n\n Getting http://uk.rdbranch.com/WebServices/Epos/v1/Restaurant/4075\n")

	client, req, _ := api.NewRequest("GET", "/Restaurant/4075", nil)
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))

	fmt.Println("\n\n Getting http://uk.rdbranch.com/WebServices/Epos/v1/Restaurant/4075/DiaryData?date=2016-01-20\n")

	client, req, _ = api.NewRequest("GET", "/Restaurant/4075/DiaryData?date=2016-01-21", nil)
	res, _ = client.Do(req)
	body, _ = ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
