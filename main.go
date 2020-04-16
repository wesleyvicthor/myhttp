package main

import (
	"flag"
	"fmt"
	"github.com/wesleyvicthor/myhttp/pkg"
	"os"
)

func main() {
	const defaultParallel int = 10
	var limitRequests int
	flag.Usage = help
	flag.IntVar(&limitRequests, "parallel", defaultParallel, "Limit the number of parallel requests")
	flag.Parse()

	if len(flag.Args()) == 0 {
		flag.Usage()
		os.Exit(1)
	}

	myHTTP := myhttp.New(flag.Args())
	done := make(chan struct{})
	result, fail := myHTTP.ProcessParallel(limitRequests, done)
	defer close(result)

	for {
		select {
		case resp := <-result:
			fmt.Printf("%s %x\n", resp.Url, resp.MD5)
		case msg := <-fail:
			fmt.Println(msg)
		case <-done:
			return
		}
	}
}

func help() {
	fmt.Println(`myhttp:
	Display the MD5 hash of the response body for the given url.
	Usage:
	    myhttp [option] [urls ...]
	Options:
	    -parallel       Limit the number of parallel requests
	    -help           Display this help message.`)
}