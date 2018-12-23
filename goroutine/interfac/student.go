package main

import (
	"fmt"
	"math/rand"
	"sort"
)

type Student struct {
	name  string
	age   int
	score float32
}

type StudentList []Student

func (p StudentList) Len() int {
	return len(p)
}

func (p StudentList) Less(i, j int) bool {
	return p[i].age < p[j].age
}

func (p StudentList) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func main() {
	var stus StudentList

	for i := 0; i < 10; i++ {
		stu := Student{
			name:  fmt.Sprintf("stu%d", rand.Intn(100)),
			age:   rand.Intn(100),
			score: rand.Float32() * 100,
		}
		stus = append(stus, stu)
	}

	for _,v := range stus{
		fmt.Println(v)
	}
	fmt.Println("\n\n")
	sort.Sort(stus)

	for _,v := range stus{
		fmt.Println(v)
	}
}
