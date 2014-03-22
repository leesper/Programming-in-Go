package main

import (
    "fmt"
    "math"
    "strings"
    "sort"
)

func main() {
    uniqueSlice := UniqueInts([]int{9, 1, 9, 5, 4, 4, 2, 1, 5, 4, 8, 8, 4, 3, 6, 9, 5, 7, 5})

    irregularMatrix := [][]int {
        {1, 2, 3, 4},
        {5, 6, 7, 8},
        {9, 10, 11},
        {12, 13, 14, 15},
        {16, 17, 18, 19, 20},
    }
    flattenSlice := Flatten(irregularMatrix)

    matrix := Make2D([]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20}, 3)

    iniData := []string {
        "; Cut down copy of Mozilla application.ini file",
        "",
        "[App]",
        "Vendor=Mozilla",
        "Name=Iceweasel",
        "Profile=mozilla/firefox",
        "Version=3.5.16",
        "[Gecko]",
        "MinVersion=1.9.1",
        "MaxVersion=1.9.1.*",
        "[XRE]",
        "EnableProfileMigrator=0",
        "EnableExtensionManager=1",
    }
    iniMap := ParseIni(iniData)

    fmt.Printf("UniqueInts: %v\n", uniqueSlice)
    fmt.Printf("Flatten: %v\n", flattenSlice)
    fmt.Printf("Make2D: %v\n", matrix)
    fmt.Printf("ParseIni: %v\n", iniMap)
    PrintIni(iniMap)
}

func UniqueInts(slice []int) []int {
    count := map[int]int{}
    result := []int{}
    for _, val := range slice {
        if _, found := count[val]; !found {
            result = append(result, val)
        }
        count[val]++
    }
    return result
}

func Flatten(matrix [][]int) []int {
    result := []int{}
    for i := range matrix {
        for j := range matrix[i] {
            result = append(result, matrix[i][j])
        }
    }
    return result
}

func Make2D(slice []int, ncols int) [][]int {
    matrix := make([][]int, int(math.Ceil(float64(len(slice)) / float64(ncols))))
    i := 0
    for len(slice) > 0 {
        if len(slice) > ncols {
            matrix[i] = append(matrix[i], slice[:ncols]...)
            slice = slice[ncols:]
        } else {
            matrix[i] = append(matrix[i], slice[:]...)
            for j := 0; j < ncols - len(slice); j++ {
                matrix[i] = append(matrix[i], 0)
            }
            slice = slice[len(slice):]
        }
        i++
    }
    return matrix
}

func ParseIni(lines []string) map[string]map[string]string {
    result := map[string]map[string]string{}
    var group string
    for _, line := range lines {
        line = strings.TrimSpace(line)

        if strings.HasPrefix(line, ";") || len(line) == 0 {
            continue
        }

        if strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            group = line[1:len(line)-1]
            result[group] = map[string]string{}
        } else {
            index := strings.Index(line, "=")
            key := line[:index]
            val := line[index+1:]
            result[group][key] = val
        }
    }
    return result
}

func PrintIni(inimap map[string]map[string]string) {
    groups := []string{}
    lines := []string{}

    for group := range inimap {
        groups = append(groups, group)
    }
    sort.Strings(groups)

    for _, group := range groups {
        lines = append(lines, fmt.Sprintf("[%s]", group))
        keys := []string{}
        for key := range inimap[group] {
            keys = append(keys, key)
        }
        sort.Strings(keys)

        for _, key := range keys {
            lines = append(lines, fmt.Sprintf("%s=%s", key, inimap[group][key]))
        }
    }

    for index, line := range lines {
        if index > 0 && strings.HasPrefix(line, "[") && strings.HasSuffix(line, "]") {
            fmt.Printf("\n")
        }
        fmt.Println(line)
    }
}
