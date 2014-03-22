package linkutil

import (
	"io"
	"net/http"
	"io/ioutil"
	"regexp"
)

//`<a[^>]+href=['"]?([^'">]+)['"]?`

var lkrx *regexp.Regexp

func init() {
	lkrx = regexp.MustCompile(`<a[^>]+href=['"]?([^'">]+)['"]?`)
}

func LinksFromURL(url string) ([]string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	links, err := LinksFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	return links, nil
}

func LinksFromReader(reader io.Reader) ([]string, error) {
	html, err := ioutil.ReadAll(reader)
	if err != nil {
		return nil, err
	}
	link2count := map[string]int{}
	for _, match := range lkrx.FindAllSubmatch(html, -1) {
		link2count[string(match[1])]++
	}
	links := []string{}
	for k, _ := range link2count {
		links = append(links, k)
	}
	return links, nil
}