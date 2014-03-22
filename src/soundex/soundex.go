package main

import (
    "fmt"
    "net/http"
    "strings"
    "bytes"
    "strconv"
    "os"
    "log"
    "bufio"
    "io"
)

const (
    pageTop = "<html><title>%s</title><body>"
    pageBot = "</body></html>"
    homeBody = `<h1>Soundex</h1>
                <p> Compute soundex codes for a list of names.</p>
                <form action="/" method="POST">
                    Names(comma or space-separated): <br>
                    <input type="text" name="words"> <br>
                    <input type="submit" name="compute" value="Submit">
                </form>`
    homeTableTop = `<table border="1"><tr><th>Name</th><th>Soundex</th></tr>`
    homeTableRow = "<tr><td>%s</td><td>%s</td></tr>"
    homeTableBot = "</table>"
    testTableTop = `<table border="1"><tr><th>Name</th><th>Soundex</th><th>Expected</th><th>Test</th></tr>`
    testTableRow = "<tr><td>%s</td><td>%s</td><td>%s</td><td>%s</td></tr>"
    testTableBot = "</table>"
    errorInfo = "<p>%s</p>"
)

type SoundexData struct {
    name string
    value string
}

func main() {
    http.HandleFunc("/", homeHandler)
    http.HandleFunc("/test", testHandler)
    if err := http.ListenAndServe(":8080", nil); err != nil {
        log.Fatal("failed to start server", err)
    }
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    fmt.Fprintf(w, pageTop, "Soundex")
    if err != nil {
        fmt.Fprintf(w, errorInfo)
    } else {
        fmt.Fprintf(w, homeBody)
        names := processRequest(r);
        if len(names) != 0 {
            fmt.Fprintf(w, homeTableTop)
            dataSet := calculateSoundexData(names)
            for _, item := range dataSet {
                fmt.Fprintf(w, homeTableRow, item.name, item.value)
            }
            fmt.Fprintf(w, homeTableBot)
        }
    }
    fmt.Fprintf(w, pageBot)
}

func processRequest(r *http.Request) []string {
    var names []string
    if slice, found := r.Form["words"]; found && len(slice) > 0 {
        text := strings.Replace(slice[0], ",", " ", -1)
        for _, field := range strings.Fields(text) {
            names = append(names, field)
        }
        return names
    }
    return names
}

func calculateSoundexData(names []string) (sdata []SoundexData) {
    // calculate soundex value for every name, return struct SoundexData
    for _, name := range names {
        var item SoundexData
        item.name = name
        item.value = soundex(name)
        sdata = append(sdata, item)
    }
    return sdata
}

func soundex(word string) string {
    /* Python implementation from Rosetta Code website
       def soundex(word):
           codes = ("bfpv","cgjkqsxz", "dt", "l", "mn", "r")
           soundDict = dict((ch, str(ix+1)) for ix,cod in enumerate(codes) for ch in cod)
           cmap2 = lambda kar: soundDict.get(kar, '9')
           sdx =  ''.join(cmap2(kar) for kar in word.lower())
           sdx2 = word[0].upper() + ''.join(k for k,g in list(groupby(sdx))[1:] if k!='9')
           sdx3 = sdx2[0:4].ljust(4,'0')
           return sdx3
    */
    codes := []string{"bfpv", "cgjkqsxz", "dt", "l", "mn", "r"}
    var soundDict = make(map[rune]int)
    for ix, cod := range codes {
        for _, c := range cod {
            soundDict[c] = ix + 1
        }
    }
    cmap2 := func(kar rune) int {
        if value, found := soundDict[kar]; found {
            return value
        } else {
            return 9
        }
    }
    var buffer bytes.Buffer
    for _, ch := range strings.ToLower(word) {
        buffer.WriteString(strconv.Itoa(cmap2(ch)))
    }
    sdx := buffer.String()
    buffer.Reset()
    var lastChar rune
    for _, ch := range sdx[1:] {
        if string(ch) != "9" && lastChar != ch {
            buffer.WriteString(string(ch))
        }
        lastChar = ch
    }
    sdx2 := strings.ToUpper(string(word[0])) + buffer.String()
    sdx3 := sdx2[0:]
    buffer.Reset()
    if len(sdx3) < 4 {
        for i := 0; i < 4 - len(sdx3); i++ {
            buffer.WriteString("0")
        }
    } else {
        sdx3 = sdx3[:4]
    }
    sdx3 = sdx3 + buffer.String()
    return sdx3
}

const test_file_name = "soundex-test-data.txt"
func testHandler(w http.ResponseWriter, r *http.Request) {
    err := r.ParseForm()
    fmt.Fprintf(w, pageTop, "Soundex Test")
    fmt.Fprintf(w, testTableTop)
    test_file, err := os.Open(test_file_name)
    if err != nil {
        log.Fatal(err)
    }
    defer test_file.Close()

    reader := bufio.NewReader(test_file)
    var line string
    var sdata []SoundexData
    eof := false
    for !eof {
        line, err = reader.ReadString('\n')
        if err == io.EOF {
            eof = true
        }
        fields := strings.Fields(line)
        if len(fields) == 2 {
            var item SoundexData
            item.name = fields[1]
            item.value = fields[0]
            sdata = append(sdata, item)
        }
    }

    for _, iter := range sdata {
        sdx := soundex(iter.name)
        var pass string
        if iter.value == sdx {
            pass = "PASS"
        } else {
            pass = "FAIL"
        }
        fmt.Fprintf(w, testTableRow, iter.name, sdx, iter.value, pass)
    }

    fmt.Fprintf(w, testTableBot)
}
