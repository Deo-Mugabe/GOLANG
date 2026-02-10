package validparentheses

func main() {
	input := "([)]"
	println(isValid(input))
}

func isValid(s string) bool {
	myMap := map[rune]rune{']': '[', ')': '(', '}': '{'}
	var stack []rune

	for _, ch := range s {
		if closing, found := myMap[ch]; found {

			if len(stack) == 0 || closing != stack[len(stack)-1] {
				return false
			}
			stack = stack[:len(stack)-1]
		} else {
			stack = append(stack, ch)
		}
	}
	return len(stack) == 0
}
