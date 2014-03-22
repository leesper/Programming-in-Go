package main

import (
	"os"
	"fmt"
	"net/url"
	"strings"
	"path/filepath"
	"linkcheck/linkutil"
)

var (
	addCh	chan string
	qryCh	chan string
	resCh	chan bool
)

func init() {
	addCh = make(chan string)
	qryCh = make(chan string)
	resCh = make(chan bool)
}

func main() {
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s url\n", filepath.Base(os.Args[0]))
		return
	}
	
	uri := os.Args[1]
	
	if !strings.HasPrefix(uri, "http://") {
		uri = "http://" + uri
	}
	
	host, err := url.Parse(uri)
	if err != nil {
		fmt.Println("url.Parse: ", err)
		return
	}
	
	fmt.Println("host: ", host.Host)
	
	go runSharedMap()
	
	checkURL(uri, "http://" + host.Host)
	
}

func runSharedMap() {
	checked := map[string]bool{}
	for {
		select {
		case url := <-addCh:
			checked[url] = true
		case url := <-qryCh:
			_, found := checked[url]
			resCh<- found
		}
	}
}

func alreadyChecked(uri string) bool {
	qryCh <- uri
	if <-resCh {
		return true
	}
	addCh <- uri
	return false
}

func processLink(link, site string, pInfo *[]string, done chan bool) int {
	if strings.HasPrefix(link, "mailto:") || strings.HasPrefix(link, "ftp://") {
		info := fmt.Sprintf("  - can't check non-http link: %s", link)
		*pInfo = append(*pInfo, info)
		return 0
	}
	
	if strings.HasPrefix(link, "http://") {	// outer links ignored
		info := fmt.Sprintf("  + checked %s", link)
		*pInfo = append(*pInfo, info)
		return 0
	}
	
	if !strings.HasSuffix(link, ".html") && !strings.HasSuffix(link, ".htm"){ // non html files ignored
		info := fmt.Sprintf("  + checked %s", link)
		*pInfo = append(*pInfo, info)
		return 0
	}
	
	go func() {
		checkURL(site + "/" + link, site)
		done <- true
	}()
	return 1
}

func checkURL(uri string, site string) {
	if alreadyChecked(uri) {
		return
	}
	
	uris, err := linkutil.LinksFromURL(uri)
	if err != nil {
		fmt.Println("* can't get http URL: ", uri)
		return
	}
	
	pending := 0
	infos := []string{}
	done := make(chan bool)
	
	fmt.Println("+ read ", uri)
	for _, link := range uris {
		pending += processLink(link, site, &infos, done)
	}
	
	fmt.Println("+ links on ", uri)
	for _, info := range infos {
		fmt.Println(info)
	}
	
	for i := 0; i < pending; i++ {
		<-done
	}
	
}