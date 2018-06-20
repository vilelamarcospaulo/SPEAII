package main

import (
	"SPEAII/speaii"
	"fmt"
)

func main() {
	ag := speaii.SPEAII{}
	ag.Run(50, 500, 300, 0.02)

	for _, ind := range ag.CurrentPopulation {
		fmt.Println(ind.DNA, "||", ind.Goals)
	}
}
