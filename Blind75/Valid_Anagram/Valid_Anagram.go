package main

import "fmt"

func main() {
	s := "racecar"
	t := "carrace"

	fmt.Println(isAnagram(s, t))
}

func isAnagram(s string, t string) bool {
	if len(s) != len(t) {
		return false
	}

	myMap := make(map[rune]int)

	for _, ch := range s {
		myMap[ch]++
	}

	for _, ch := range t {
		myMap[ch]--
		if myMap[ch] < 0 {
			return false
		}
	}

	return true

}
