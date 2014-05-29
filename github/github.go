package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/dieyushi/golang-brutedict"
)

const (
	reqUrl    = "https://github.com/signup_check/username"
	minLength = 1
	maxLength = 39
	isNum     = false
	isLow     = true
	isCap     = false
)

var bruteDict *brutedict.BruteDict

func Run() {
	bruteDict = brutedict.New(isNum, isLow, isCap, minLength, maxLength)

	c := make(chan string)

	const maxOpen = 500
	openConnections := 0
	completed := 0
	found := 0
	lastUsername := ""

	go func(success chan string) {

		for {
			if openConnections < maxOpen {
				openConnections++
				go func() {
					// We must trim \u0000 because of a bug in the brutedict package.
					// It normally generates strings that look like aab\x00\x00\x00...
					username := strings.Trim(bruteDict.Id(), "\u0000")

					defer func() {
						openConnections--
						completed++
						lastUsername = username

						fmt.Printf("\rCompleted: %d, Open: %d, Found: %d, Last: %s", completed, openConnections, found, lastUsername)
					}()

					resp, err := http.PostForm(reqUrl, url.Values{"value": {username}})
					if err != nil {
						fmt.Println(err.Error())
					}
					defer resp.Body.Close()

					if resp.StatusCode == http.StatusOK {
						found++
						success <- username
					}
				}()
			} else {
				// Give some time to the CPU
				time.Sleep(time.Millisecond * 100)
			}
		}
	}(c)

	for {
		select {
		case _ = <-c:
			//fmt.Println(username)
		}
	}

	bruteDict.Close()
}
