package linkutil_test

import (
	"os"
	"fmt"
	"strings"
	"testing"
	"io/ioutil"
	"linkcheck/linkutil"
)

func TestLinksFromURL(t *testing.T) {
	links, err := linkutil.LinksFromURL("http://www.qtrac.eu")
	if err != nil {
		t.Error("can't get from www.qtrac.eu: ", err)
	}
	fmt.Printf("links: %v\n", links)
}

func TestLinksFromReader(t *testing.T) {
	htmlReader, err := os.Open("index.html")
	if err != nil {
		t.Error("can't open file index.html")
	}
	defer htmlReader.Close()
	
	linksReader, err := os.Open("index.links")
	if err != nil {
		t.Error("can't open file index.links")
	}
	defer linksReader.Close()
	
	links, err := linkutil.LinksFromReader(htmlReader)
	if err != nil {
		t.Error("can't read index.html")
	}
	
	content, err := ioutil.ReadAll(linksReader)
	if err != nil {
		t.Error("can't read index.links")
	}
	
	for _, link := range links {
		if !strings.Contains(string(content), link) {
			t.Errorf("%s not exists", link)
		}
	}
	
}