package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/gofuzz"
)

// parseUrls appends a query to a slice of url strings, creates slice of *url.URLs
func parseUrls(strings []string, query string) (urls []*url.URL) {
	for _, str := range strings {
		parsed, err := url.Parse(fmt.Sprintf("%s/?uuid=%s", str, query))
		if err != nil {
			panic(err)
		}
		urls = append(urls, parsed)
	}
	return
}

// makeGetRequests takes a slice of *url.URLs, gets each of them, and makes a map of *url.URLs to responses
func makeGetRequests(urls []*url.URL) map[*url.URL]string {
	responses := make(map[*url.URL]string)
	for _, get := range urls {
		resp, err := http.Get(get.String())
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		responses[get] = string(body)
	}
	return responses
}

// formatRepsonseMap takes a map of *url.URLs to response strings and formats them nicely
func formatResponseMap(responses map[*url.URL]string) string {
	var formattedResponses []string
	for k, v := range responses {
		formattedResponses = append(formattedResponses, fmt.Sprintf("%s -> %s", k, v))
	}
	return strings.Join(formattedResponses, "\n")
}

func main() {
	// query string to be run on all urls
	query := flag.String("query", "", "the query string to make on each url")
	flag.Parse()

	// all the different url strings
	args := flag.Args()

	// if there is no query, fuzz up something random
	if *query == "" {
		fuzzer := fuzz.New()
		fuzzer.Fuzz(query)
	}

	// parse all the urls
	urls := parseUrls(args, *query)

	fmt.Println(formatResponseMap(makeGetRequests(urls)))
}
