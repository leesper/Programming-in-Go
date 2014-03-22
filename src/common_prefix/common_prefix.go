package main

import (
    "fmt"
    "strings"
    "math"
)

func main() {
    pathgroups := [][]string {
        {"/home/user/goeg", "/home/user/goeg/prefix", "/home/user/goeg/prefix/extra"},
        {"/home/user/goeg", "/home/user/goeg/prefix", "/home/user/prefix/extra"},
        {"/home/user/goeg", "/home/users/goeg", "/home/userspace/goeg"},
        {"/home/user/goeg", "/tmp/user", "/var/log"},
        {"/home/mark/goeg", "/home/user/goeg"},
        {"home/user/goeg", "/tmp/user", "/var/log"},
    }
    for _, group := range pathgroups {
        charPref := CommonPrefix(group)
        pathPref := CommonPathPrefix(group)
        if charPref == pathPref {
            fmt.Printf("char*path prefix: %q == %q\n", charPref, pathPref)
        } else {
            fmt.Printf("char*path prefix: %q != %q\n", charPref, pathPref)
        }
    }
}

func CommonPrefix(texts []string) string {
    // return common prefix of all strings in slice
    strs := [][]rune{}
    for _, text := range texts {
        strs = append(strs, []rune(text))
    }
    shortest := findShortestStr(strs)
    for len(shortest) > 0 {
        isOK := true
        for _, str := range strs {
            if !strings.HasPrefix(string(str), shortest) {
                isOK = false
            }
        }
        if isOK {
            return shortest
        }
        shortest = shortest[:len(shortest)-1]
    }
    return shortest
}

func findShortestStr(strs [][]rune) string {
    min := math.MaxInt32
    var shortest string
    for _, str := range strs {
        if len(str) < min {
            min = len(str)
            shortest = string(str)
        }
    }
    return shortest
}

func CommonPathPrefix(paths []string) string {
    // return common prefix of all paths
    min := math.MaxInt32
    series := [][]string{}
    for _, path := range paths {
        subdirs := separatePath(path)
        series = append(series, subdirs)
        if len(subdirs) < min {
            min = len(subdirs)
        }
    }

    commons := []string{}
    for i := 0; i < min; i++ {
        for j := 1; j < len(series); j++ {
            if series[0][i] != series[j][i] {
                goto FINAL
            }
        }
        commons = append(commons, series[0][i])
    }
    FINAL:
    if len(commons) > 1 {
        return strings.Replace(strings.Join(commons, "/"), ".", "", 1)
    }
    return strings.Replace(strings.Join(commons, "/"), ".", "/", 1)
}

func separatePath(path string) []string {
    subdirs := []string{}
    if strings.HasPrefix(path, "/") {
        subdirs = append(subdirs, ".")
        path = path[1:]
    }
    strs := strings.Split(path, "/")
    subdirs = append(subdirs, strs...)
    return subdirs
}
