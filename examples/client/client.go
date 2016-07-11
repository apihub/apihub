package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

func main() {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(network, addr string) (net.Conn, error) {
				// return net.DialTimeout("unix", "/tmp/apihub.sock", 2*time.Second)
				return net.DialTimeout("tcp", ":8000", 2*time.Second)
			},
		},
	}
	request, _ := http.NewRequest(http.MethodGet, "http://api/", nil)
	resp, err := client.Do(request)
	if err != nil {
		panic(err)
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%+v", string(body))
}
