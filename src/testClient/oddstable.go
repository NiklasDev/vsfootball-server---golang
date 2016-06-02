package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

// "strings"
// "time"
// "crypto/sha1"
)

func main() {
	for i := 0; i < 100000; i++ {
		for defensive := 24; defensive < 36; defensive++ {
			for offensive := 1; offensive < 23; offensive++ {
				if offensive == 13 || offensive == 14 {
					continue
				}
				offensive := 20
				values := make(url.Values)
				values.Set("offensivePlayId", strconv.Itoa(offensive))
				values.Set("defensivePlayId", strconv.Itoa(defensive))
				client, _ := http.PostForm("http://localhost:8082/odds", values)
				defer client.Body.Close()

				body, _ := ioutil.ReadAll(client.Body)
				fmt.Println(string(body))
			}
		}
	}
}
