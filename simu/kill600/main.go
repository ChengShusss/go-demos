package main

import (
	"fmt"
	"math/rand"
)

const (
	TotalPeople = 600
)

func main() {
	result := map[int]int{}

	for i := 0; i < 1000000; i++ {
		if i%10000 == 0 {
			fmt.Printf("%d%%...\n", i/10000)
		}
		result[getAlive()] += 1
	}
	// res := getAlive()
	// fmt.Printf("res: %v\n", res)

	fmt.Printf("%v\n", result)

}

func getAlive() int {
	que := []int{}
	for i := 1; i <= TotalPeople; i += 1 {
		que = append(que, i)
	}

	for i := 1; i < TotalPeople; i += 1 {
		k := rand.Intn(len(que)/2+1) * 2
		switch k {
		case 0:
			que = que[1:]
		case len(que):
			que = que[:len(que)-1]
		default:
			que = append(que[:k], que[k+1:]...)
		}
		// fmt.Printf("k: %v, que: %v\n", k+1, que)
	}

	return que[0]
}
