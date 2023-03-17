// Echo2 prints its command-line arguments.
package main

import (
	"fmt"
	"os"
	"strings"
)

func main() {
	mod1()
}

func version1() {
	var s, sep string
	for i := 1; i < len(os.Args); i++ {
		s += sep + os.Args[i]
		sep = " "
	}

	fmt.Println(s)
}

func version2() {
	var s, sep string
	for _, arg := range os.Args[1:] {
		s += sep + arg
		sep = " "
	}
	fmt.Println(s)
}

func version3() {
	fmt.Println(strings.Join(os.Args[1:], " "))
}

func mod1() {
	for i := 1; i < len(os.Args); i++ {
		fmt.Println(i, " : ", os.Args[i])
	}
}
