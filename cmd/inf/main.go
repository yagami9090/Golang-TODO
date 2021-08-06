package main

import "fmt"

func main() {
	// "abcd" => []string{"ab","cd"}
	// "abcde" => []string{"ab","cd","e_"}
	fmt.Println(couple("abcd"))
	fmt.Println(couple("abcde"))
	fmt.Println(couple("abcdef"))
	fmt.Println(couple("abcdefg"))
}

type Intner interface {
	Intn(int) int
}

func couple(s string) (ss []string) {
	for s += "_"; len(s) > 1; s = s[2:] {
		ss = append(ss, s[:2])
	}
	return
}
