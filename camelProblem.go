package main

import (
	"fmt"
	"unicode"
)

func main(){
	var str string
	fmt.Scanf("%s\n", &str)
	ans := 1

	for _, ch := range str {
		if unicode.IsUpper(ch) {
			ans += 1
		}
	}

	fmt.Printf("String contains %v words\n", ans)
}