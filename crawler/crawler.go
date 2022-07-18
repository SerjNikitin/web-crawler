package crawler

import (
	"fmt"
	"html"
	"io"
	"log"
	"net/http"
	"net/url"
	"regexp"
)

func ScanUrl() {
	for true {
		fmt.Print("Please write url or 'exit': ")
		var message string
		scan, err := fmt.Scan(&message)
		if err != nil || scan != 0 {
			switch message {
			case "exit":
				log.Print("bay-bay :)")
				return
			default:
				crawlUrl(message)
			}
		} else {
			log.Println(err)
		}
	}
}

func crawlUrl(urlToCrawl string) {
	var (
		err      error
		content  string
		urlToGet *url.URL
		setLinks = make(map[string]int)
	)

	if urlToGet, err = url.Parse(urlToCrawl); err != nil {
		log.Println(err)
		return
	}
	if content, err = getUrlContent(urlToGet.String()); err != nil {
		log.Println(err)
		return
	}
	content = html.UnescapeString(content)
	if setLinks, err = parseLinks(urlToGet, content); err != nil {
		log.Println(err)
		return
	}
	for k := range setLinks {
		log.Println(k)
	}
}

func getUrlContent(urlToGet string) (string, error) {
	var (
		err     error
		content []byte
		resp    *http.Response
	)
	if resp, err = http.Get(urlToGet); err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", err
	}
	if content, err = io.ReadAll(resp.Body); err != nil {
		return "", err
	}
	return string(content), err
}

func parseLinks(urlToGet *url.URL, content string) (map[string]int, error) {
	var (
		err       error
		setLinks  = make(map[string]int)
		matches   [][]string
		findLinks = regexp.MustCompile("<a.*?href=\"(.*?)\"")
	)
	matches = findLinks.FindAllStringSubmatch(content, -1)

	for _, val := range matches {
		var linkUrl *url.URL

		if linkUrl, err = url.Parse(val[1]); err != nil {
			return setLinks, err
		}
		if linkUrl.IsAbs() {
			setLinks[linkUrl.String()] = 0
		} else {
			setLinks[urlToGet.Scheme+"://"+urlToGet.Host+linkUrl.String()] = 0
		}
	}
	return setLinks, err
}
