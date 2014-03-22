package main

import (
    "fmt"
    "unicode/utf8"
)

func main() {
    fmt.Printf("PULLUP: %t\n", IsPalinDrome("PULLUP"))
    fmt.Printf("ROTOR: %t\n", IsPalinDrome("ROTOR"))
    fmt.Printf("DECIDED: %t\n", IsPalinDrome("DECIDED"))
}

func IsPalinDrome(word string) bool {
    left := 0
    right := utf8.RuneCountInString(word)
    for len(word[left:right]) > 1 {
        first, fsize := utf8.DecodeRuneInString(word[left:])
        last, lsize := utf8.DecodeLastRuneInString(word[:right])
        if first != last {
            return false
        }
        left += fsize
        right -= lsize
    }
    return true
}
