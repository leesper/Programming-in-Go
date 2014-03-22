package main

import (
	"os"
	"image"
	_ "image/gif"
	_ "image/png"
	_ "image/jpeg"
	"fmt"
	"log"
	"runtime"
	"path/filepath"
	"strings"
)

type imageInfoType struct {
	filename	string
	height		int
	width		int
}

var workers = runtime.NumCPU()
func main() {
	runtime.GOMAXPROCS(runtime.NumCPU()) // use all cpu cores
	if len(os.Args) == 1 || os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("usage: %s [image files]\n", filepath.Base(os.Args[0]))
		os.Exit(1)
	}
	
	infoChan := make(chan imageInfoType)
	imgfnameChan := make(chan string)
	done := make(chan struct{})
	
	go filterImgFNames(imgfnameChan, os.Args[1:])
	
	for i := 0; i < workers; i++ {
		go procImgInfo(infoChan, imgfnameChan, done)
	}
	
	waitAndOutputTags(infoChan, done)
}

func filterImgFNames(imgfnameChan chan<- string, args []string) {
	for _, fname := range args {
		isJpg := strings.HasSuffix(fname, "jpg")
		isPng := strings.HasSuffix(fname, "png")
		isGif := strings.HasSuffix(fname, "gif")
		if isJpg || isPng || isGif {
			imgfnameChan <- fname
		}
	}
	close(imgfnameChan)
}

func procImgInfo(infoChan chan<- imageInfoType, imgfChan <-chan string, done chan<- struct{}) {
	for f := range imgfChan {
		fh, err := os.Open(f)
		if err != nil {
			log.Println("error: ", err)
		}
		defer fh.Close()
		config, _, err := image.DecodeConfig(fh)
		if err != nil {
			log.Println("error: ", err)
		}
		infoChan <- imageInfoType{filepath.Base(f), config.Height, config.Width}
	}
	done <- struct{}{}
}

func waitAndOutputTags(infoChan <-chan imageInfoType, done <-chan struct{}) {
	for working:= workers; working > 0; {
		select {
			case <-done:
				working--
			case info := <-infoChan:
				fmt.Printf(`<img src="%s" width="%d" height="%d" />`, 
					info.filename, info.width, info.height)
				fmt.Println()
		}
	}
DONE:
	for {
		select {
			case info := <-infoChan:
				fmt.Printf(`<img src="%s" width="%d" height="%d" />`, 
					info.filename, info.width, info.height)
				fmt.Println()
			default:
				break DONE
		}
	}
}