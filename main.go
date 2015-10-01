package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

// getResponses GETs a slice of strings, maps them to server responses
func getResponses(urlStrings []string) map[string]string {
	responses := make(map[string]string)
	for _, urlString := range urlStrings {
		resp, err := http.Get(urlString)
		if err != nil {
			panic(err)
		}
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		responses[urlString] = string(body)
	}
	return responses
}

// formatRepsonseMap formats map[string]string nicely
func formatResponseMap(responses map[string]string) string {
	var formattedResponses []string
	for k, v := range responses {
		formattedResponses = append(formattedResponses, fmt.Sprintf("%s -> %s", k, v))
	}
	return strings.Join(formattedResponses, "\n")
}

// for randstring
const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ1234567890"
const (
	letterIdxBits = 6                    // 6 bits to represent a letter index
	letterIdxMask = 1<<letterIdxBits - 1 // All 1-bits, as many as letterIdxBits
	letterIdxMax  = 63 / letterIdxBits   // # of letter indices fitting in 63 bits
)

var src = rand.NewSource(time.Now().UnixNano())

func RandStringBytesMaskImprSrc(n int) string {
	b := make([]byte, n)
	// A src.Int63() generates 63 random bits, enough for letterIdxMax characters!
	for i, cache, remain := n-1, src.Int63(), letterIdxMax; i >= 0; {
		if remain == 0 {
			cache, remain = src.Int63(), letterIdxMax
		}
		if idx := int(cache & letterIdxMask); idx < len(letterBytes) {
			b[i] = letterBytes[idx]
			i--
		}
		cache >>= letterIdxBits
		remain--
	}
	return string(b)
}

func main() {
	toAppend := flag.String("append", "", "string to append to each url")
	toPrepend := flag.String("prepend", "", "string to prepend to each url")
	fuzz := flag.Bool("fuzz", false, "add fuzz string (letter only) after appended string")
	sampleSize := flag.Int("samplesize", 1, "number of times to get each url, fuzz is regenerated each time")
	flag.Parse()

	// all the different url strings
	args := flag.Args()
	var urls []string

	// append, prepend all url strings
	for _, url := range args {
		urls = append(urls, fmt.Sprintf("%s%s%s", *toPrepend, url, *toAppend))
	}

	// take samples
	var responses []map[string]string
	for i := 0; i < *sampleSize; i++ {
		// if fuzz is enabled, append fuzz to toAppend
		if *fuzz {
			var fuzzedUrls []string
			fuzzString := RandStringBytesMaskImprSrc(rand.Intn(50))
			for _, url := range urls {
				// rand string 0-50 chars
				fuzzedUrls = append(fuzzedUrls, fmt.Sprintf("%s%s", url, fuzzString))
			}
			responses = append(responses, getResponses(fuzzedUrls))
			continue
		}
		responses = append(responses, getResponses(urls))
	}

	// print formatted responses
	inconsistentTally := 0
	for _, responsesMap := range responses {
		// tally responses
		responsesTally := make(map[string]int)
		for _, response := range responsesMap {
			responsesTally[response]++
		}
		if len(responsesTally) > 1 {
			fmt.Println("Inconsistent!")
			inconsistentTally++
		}
		fmt.Println(formatResponseMap(responsesMap) + "\n")
	}
	fmt.Println(fmt.Sprintf("%v total inconsistent responses in %v trials. Rate: %v%%.",
		inconsistentTally, *sampleSize, float32(inconsistentTally)/float32(*sampleSize)*100))
}
