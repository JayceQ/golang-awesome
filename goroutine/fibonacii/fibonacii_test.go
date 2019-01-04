package main

import "testing"

//func TestSubstr(t *testing.T) {
//	f := fibonacii()
//	printFileContents(f)
//}

func BenchmarkFib(b *testing.B){
	f := fibonacii()
	printFileContents(f)
}