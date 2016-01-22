ResDiary API Library
====================

Install
-------

`go get github.com/tellthechef/rdapi`

Usage
-----

```
import (
	"github.com/tellthechef/rdapi"
	"fmt"
	"io/ioutil"
)

func main() {
	api := rdapi.New("consumerKey", "consumerSecret", "secondSecret")
	if err := api.Authenticate(); err != nil {
		return
	}

	// Use api.NewRequest to create an authenticated http.Client
	client, req, _ := api.NewRequest("GET", "/Restaurant/{ID}", nil)

	// do business as usual
	res, _ := client.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
}
```