package myhttp

import (
	"crypto/md5"
	"io/ioutil"
	"net/http"
	"regexp"
	"sync"
)

// MyHTTP
type MyHTTP struct {
	urls []string
}

// Response Result of a request
type Response struct {
	Url string
	MD5 [md5.Size]byte
}

// New initialize MyHTTP
func New(urls []string) *MyHTTP {
	myHTTP := MyHTTP{
		make([]string, len(urls)),
	}
	re := regexp.MustCompile(`^http://|^https://`)
	for key, url := range urls {
		if !re.MatchString(url) {
			url = "http://" + url
		}

		myHTTP.urls[key] = url
	}

	return &myHTTP
}

func (mh MyHTTP) Urls() []string {
	return mh.urls
}

// Process return the response body as md5 hash
func (mh MyHTTP) Process(url string) (*Response, error) {
	res, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	return &Response{url, md5.Sum(body)}, nil
}

//ProcessParallel spawn processes for the given limit
func (mh MyHTTP) ProcessParallel(limit int, done chan struct{})(chan *Response, chan string) {
	wg := sync.WaitGroup{}
	wg.Add(len(mh.Urls()))
	process := make(chan string, len(mh.Urls()))
	result := make(chan *Response, len(mh.Urls()))
	fail := make(chan string)
	for proc := 0; proc < limit; proc++ {
		go func() {
			for {
				select {
				case url := <-process:
					res, err := mh.Process(url)
					if err != nil {
						fail <- err.Error()
					} else {
						result <-res
					}
					wg.Done()
				case <-done:
					return
				}
			}
		}()
	}

	for _, url := range mh.Urls() {
		process <-url
	}

	go func() {
		wg.Wait()
		close(done)
	}()

	return result, fail
}
