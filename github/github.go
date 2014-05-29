package github

import (
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/dieyushi/golang-brutedict"
)

const (
	reqURL    = "https://github.com/signup_check/username"
	minLength = 1
	maxLength = 39
	isNum     = false
	isLow     = true
	isCap     = false
)

type username string

type result struct {
	username username
	err      error
	duration time.Duration
	valid    bool
}

type GitHub struct {
	ConcurrencyLevel int
	jobs             chan username
	results          chan *result
}

var bruteDict *brutedict.BruteDict

func (g *GitHub) Run() {
	g.results = make(chan *result)

	bruteDict = brutedict.New(isNum, isLow, isCap, minLength, maxLength)

	var wg sync.WaitGroup
	wg.Add(g.ConcurrencyLevel)

	g.jobs = make(chan username, g.ConcurrencyLevel)

	for i := 0; i < g.ConcurrencyLevel; i++ {
		go func() {
			g.worker(g.jobs, bruteDict)
			wg.Done()
		}()
	}

	go func() {
		for {
			// We must trim \u0000 because of a bug in the brutedict package.
			// It normally generates strings that look like aab\x00\x00\x00...
			g.jobs <- username(strings.Trim(bruteDict.Id(), "\u0000"))
		}
	}()

	go func() {
		completed, validCount := 0, 0
		for result := range g.results {
			completed++

			if result.valid {
				validCount++
			}

			fmt.Printf("\rCOMPLETED: %d, VALID: %d, LAST TRY: %s", completed, validCount, result.username)
		}
	}()

	wg.Wait()

	bruteDict.Close()
}

func (g *GitHub) worker(ch chan username, dict *brutedict.BruteDict) {
	for user := range ch {
		var valid bool
		start := time.Now()

		resp, err := http.PostForm(reqURL, url.Values{"value": {string(user)}})
		if resp != nil {
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				valid = true
			}
		}

		g.results <- &result{
			username: user,
			err:      err,
			duration: time.Now().Sub(start),
			valid:    valid,
		}
	}
}
