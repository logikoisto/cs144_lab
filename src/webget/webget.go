package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
)

// lab0 - webget
func main() {
	if len(os.Args) != 2 {
		fmt.Println("wrong parameters")
		return
	}
	addr := os.Args[1]
	resp, err := http.Get("http://" + addr)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		fmt.Println("error code: ", resp.StatusCode)
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(body))
}
