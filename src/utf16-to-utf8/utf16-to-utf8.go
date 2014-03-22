package main

import (
	"io"
	"os"
	"fmt"
	"bufio"
	"encoding/binary"
	"path/filepath"
	"unicode/utf16"
)

const (
	SUCCESS = iota
	WRONG_ARGS
	ERROR_FILE
	ERROR_READ
	ERROR_WRITE
	ERROR_STAT
)

func main() {
	if len(os.Args) <= 1 || len(os.Args) >= 4 || 
		os.Args[1] == "-h" || os.Args[1] == "--help" {
		
		fmt.Printf("usage: %s infile [outfile]\n", filepath.Base(os.Args[0]))
		os.Exit(-WRONG_ARGS)
	}
	
	var infile *os.File
	infile, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println("1 error: ", err)
		os.Exit(-ERROR_FILE)
	}
	defer infile.Close()
	
	var outfile *os.File = os.Stdout
	if len(os.Args) == 3 {
		if outfile, err = os.Create(os.Args[2]); err != nil {
			fmt.Println("2 error: ", err)
			os.Exit(-ERROR_FILE)
		}
		defer outfile.Close()
	}
	
	if err = utf162utf8(outfile, infile); err != nil {
		fmt.Println("error: ", err)
		return
	}
	
}

func utf162utf8(outfile, infile *os.File) error {
	var boTag uint16
	err := binary.Read(infile, binary.LittleEndian, &boTag)
	if err != nil {
		fmt.Println("3 error: ", err)
		return err
	}
	
	var byteOrder binary.ByteOrder = binary.LittleEndian 
	if boTag == 0xFFFE {
		byteOrder = binary.BigEndian
	}
	
	fileWriter := bufio.NewWriter(outfile)
	defer fileWriter.Flush()	

	for {
		var ui16 uint16
		if err = binary.Read(infile, byteOrder, &ui16); err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}
		
		_, err = fileWriter.WriteString(string(utf16.Decode([]uint16{ui16})))
		if err != nil {
			return err
		}		
		
	}
	
	return nil
}