package main

import (
	"os"
	"regexp"
	"path/filepath"
	"fmt"
	"image"
	"runtime"
	"strings"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"log"
	"bufio"
	"io"
)

var workers = runtime.NumCPU()
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s [html files...]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	
	htmlFiles := make(chan string)
	done := make(chan struct{})
	results := make(chan string)
	
	go addHtmlFiles(htmlFiles, os.Args[1:])
	
	for i := 0; i < workers; i++ {
		go parseAndReplace(done, results, htmlFiles)
	}
	
	waitAndOutputResult(done, results)
}

func waitAndOutputResult(done <-chan struct{}, results <-chan string) {
	for working := workers; working > 0; {
		select {
			case <-done:
				working--
			case result := <-results:
				fmt.Printf("%s\n", result)
		}
	}
DONE:
	for {
		select {
			case result := <-results:
				fmt.Printf("%s\n", result)
			default:
				break DONE
		}
	}
}

func parseAndReplace(done chan<- struct{}, results chan<- string, htmls <-chan string) {
	imgtagrx := regexp.MustCompile(`<[iI][mM][gG][^>]+>`)
	
	for fname := range htmls {
		dir := filepath.Dir(fname)
		frh, err := os.Open(fname)
		if err != nil {
			log.Println("error: ", err)
		}
		defer frh.Close()
		
		fwh, err := os.Create(filepath.Base(fname) + "_final.html")
		if err != nil {
			log.Println("error: ", err)
		}
		defer fwh.Close()
		
		reader := bufio.NewReader(frh)
		writer := bufio.NewWriter(fwh)

		for lino := 1; ; lino++ {
			line, err := reader.ReadBytes('\n')
			if err != nil {
				if err != io.EOF {
					log.Printf("error:%d: %s\n", lino, err)
				}
				break
			}
			
			line = []byte(imgtagrx.ReplaceAllStringFunc(string(line), makeReplacer(dir)))
			
			writer.Write(line)
		}
		writer.Flush()
		results <- fmt.Sprintf("file %s finished", fname)
	}
	done <- struct{}{}
}

func makeReplacer(dir string) func(string) string {
	return func(originTag string) string {
		if strings.Contains(originTag, "width=") && strings.Contains(originTag, "height=") {
			return originTag
		}
		imgsrcrx := regexp.MustCompile(`src=["']([^"']+)["']`)
		match := imgsrcrx.FindStringSubmatch(originTag)
		if match == nil {
			log.Println("no src attr in img tag: ", originTag)
			return originTag
		}
		
		fmt.Printf("processing img tag: %s\n", originTag)
		src := match[1]
		var srcFile string
		if !filepath.IsAbs(src) {
			srcFile = dir + "/" + src
		}
		
		fh, err := os.Open(srcFile)
		if err != nil {
			log.Println("error: ", err)
			return originTag
		}
		defer fh.Close()
		
		config, _, err := image.DecodeConfig(fh)
		if err != nil {
			log.Println("error: ", err)
			return originTag
		}
		
		return fmt.Sprintf(`<img src="%s" width="%d" height="%d" />`, src, config.Width, config.Height)
	}
}
func addHtmlFiles(htmls chan<- string, flist []string) {
	for _, f := range flist {
		fname := strings.TrimSpace(f)
		if strings.HasSuffix(fname, "html") || strings.HasSuffix(fname, "htm") {
			htmls <- fname
		}
	}
	close(htmls)
}