package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	// "strings"
	"sync"
)

const (
	main_url = "http://ub.ac.id/akademik/fakultas"
	// domain   = "ub.ac.id"
)

func removeDuplicates(elements []string) []string {
	// Use map to record duplicates as we find them.
	encountered := map[string]bool{}
	result := []string{}

	for v := range elements {
		if encountered[elements[v]] == true {
			// Do not add duplicate.
		} else {
			// Record this element as an encountered element.
			encountered[elements[v]] = true
			// Append to result slice.
			result = append(result, elements[v])
		}
	}
	// Return the new slice.
	return result
}

func fetchSite(url string, wg *sync.WaitGroup) {

	// defer close(chanIn)
	defer wg.Done()
	client := &http.Client{}
	resp, err := client.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	respBody, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	//pattern, err := regexp.Compile(`<a\s+(?:[^>]*?\s+)?href="([^"]*)"`)
	pattern := regexp.MustCompile(url + "[a-z]+")
	bodyStr := string(respBody[:])
	found := pattern.FindAllString(bodyStr, -1)

	urls := removeDuplicates(found)
	urls = append(urls, url)

	for _, urlList := range urls {
		fmt.Println(urlList)
	}

}

func getSiteURL(mainURL string, wg *sync.WaitGroup) {

	defer wg.Done()
	client := &http.Client{}
	res, err := client.Get(mainURL)
	if err != nil {
		log.Fatal(err)
	}
	body, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}

	bodyStr := string(body[:])
	var pattern = regexp.MustCompile("http://" + "[a-z]+" + ".ub.ac.id/en")
	var urlStr = pattern.FindAllString(bodyStr, -1)

	for _, linkURL := range urlStr {
		var regexRep = regexp.MustCompile("en")
		var strRep = regexRep.ReplaceAllString(linkURL, "")
		linkURL := strRep

		wg.Add(1)
		go fetchSite(linkURL, wg)
		// fmt.Println(linkURL)
	}

}

func main() {
	//channel := make(chan []string)

	var wg sync.WaitGroup
	getSiteURL(main_url, &wg)

	wg.Wait()

	//var url_list = <-channel

	/*
		for _, urlList := range url_list {
			fmt.Println(urlList)
		}
	*/
}
